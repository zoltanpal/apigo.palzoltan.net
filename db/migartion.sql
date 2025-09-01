/* ======================================================================
   POWER OF WORDS — DB B -> DB A MIGRATION (PostgreSQL 14)
   Run this in the TARGET database (DB A).
   ====================================================================== */

/* ----------------------------------------------------------------------
   0) Session speed knobs (use per step or keep as-is)
   ---------------------------------------------------------------------- */
-- Wrap heavy steps in transactions with fast-but-safe defaults:
-- BEGIN;
-- SET LOCAL synchronous_commit = off;
-- SET LOCAL statement_timeout  = '30min';
-- COMMIT;

/* If you ever see 25P02 (transaction aborted), run:
ROLLBACK;
-- then re-run the failing statement alone to see the real error.
*/

/* ----------------------------------------------------------------------
   1) FDW wiring to Source DB (DB B)  — requires superuser or proper grants
   ---------------------------------------------------------------------- */
CREATE EXTENSION IF NOT EXISTS postgres_fdw;

-- Drop/recreate the foreign server to be safe
DROP SERVER IF EXISTS src_db CASCADE;

CREATE SERVER src_db
  FOREIGN DATA WRAPPER postgres_fdw
  OPTIONS (host 'localhost', dbname 'YOUR_SOURCE_DB_NAME', port '5432');

-- Map current user to a source DB user (adjust credentials if needed)
DROP USER MAPPING IF EXISTS FOR CURRENT_USER SERVER src_db;

CREATE USER MAPPING FOR CURRENT_USER
  SERVER src_db
  OPTIONS (user 'postgres', password 'YOUR_PASSWORD');

-- Dedicated schema for foreign tables
CREATE SCHEMA IF NOT EXISTS src;

-- Import only the needed tables from source's `public` schema
IMPORT FOREIGN SCHEMA public
  LIMIT TO (feeds, feed_sentiments, feed_categories, sources)
  FROM SERVER src_db INTO src;

-- Sanity checks
SELECT COUNT(*) AS src_feeds_count FROM src.feeds;
SELECT COUNT(*) AS src_sents_count FROM src.feed_sentiments;

/* ----------------------------------------------------------------------
   2) (Optional) FULL REPLACE: truncate target tables
   ---------------------------------------------------------------------- */
-- If you want a clean slate (NOT required for incremental loads):
-- Recommended: keep categories if they are already correct in target.
-- BEGIN;
-- TRUNCATE TABLE public.feed_sentiments RESTART IDENTITY;
-- TRUNCATE TABLE public.feeds RESTART IDENTITY;
-- -- If you also want to reload categories, uncomment this:
-- -- TRUNCATE TABLE public.feed_categories RESTART IDENTITY;
-- COMMIT;

/* ----------------------------------------------------------------------
   3) Partitions
   ---------------------------------------------------------------------- */
-- If your target partitions already cover the months you’ll import,
-- you can SKIP this section. Otherwise, create missing monthly partitions.

-- FEEDS monthly partitions for the date range present in source:
DO $$
DECLARE d date; mind date; maxd date;
BEGIN
  SELECT date_trunc('month', MIN(COALESCE(feed_date, published::date)))::date,
         date_trunc('month', MAX(COALESCE(feed_date, published::date)))::date
    INTO mind, maxd
  FROM src.feeds;

  IF mind IS NULL THEN
    RAISE NOTICE 'No rows in src.feeds; skipping feed partitions.';
    RETURN;
  END IF;

  d := mind;
  WHILE d <= maxd LOOP
    EXECUTE format(
      'CREATE TABLE IF NOT EXISTS feeds_%s PARTITION OF feeds
         FOR VALUES FROM (%L) TO (%L)',
      to_char(d,'YYYY_MM'), d, (d + interval '1 month')::date
    );
    d := (d + interval '1 month')::date;
  END LOOP;
END$$;

-- FEED_SENTIMENTS monthly partitions for range present in source:
DO $$
DECLARE d date; mind date; maxd date;
BEGIN
  SELECT date_trunc('month', MIN(COALESCE(s.feed_date, f.published::date)))::date,
         date_trunc('month', MAX(COALESCE(s.feed_date, f.published::date)))::date
    INTO mind, maxd
  FROM src.feed_sentiments s
  JOIN src.feeds f ON f.id = s.feed_id;

  IF mind IS NULL THEN
    RAISE NOTICE 'No rows in src.feed_sentiments; skipping sentiment partitions.';
    RETURN;
  END IF;

  d := mind;
  WHILE d <= maxd LOOP
    EXECUTE format(
      'CREATE TABLE IF NOT EXISTS feed_sentiments_%s PARTITION OF feed_sentiments
         FOR VALUES FROM (%L) TO (%L)',
      to_char(d,'YYYY_MM'), d, (d + interval '1 month')::date
    );
    d := (d + interval '1 month')::date;
  END LOOP;
END$$;

/* ----------------------------------------------------------------------
   4) (Optional) Dimension sync
   ---------------------------------------------------------------------- */
-- If your target `feed_categories` is already correct, skip this.
-- To reload/align categories (IDs preserved):
-- INSERT INTO public.feed_categories (id, name, alias, created)
-- OVERRIDING SYSTEM VALUE
-- SELECT id, name, alias, created
-- FROM src.feed_categories
-- ON CONFLICT (id) DO UPDATE
--   SET name = EXCLUDED.name,
--       alias = EXCLUDED.alias;

/* ----------------------------------------------------------------------
   5) UPSERT feeds (deduplicated by (source_id, link, feed_date))
   ---------------------------------------------------------------------- */
BEGIN;
SET LOCAL synchronous_commit = off;

/* OPTIONAL: test window — uncomment and set dates to batch by month
WITH date_window AS (
  SELECT DATE '2025-08-01' AS from_date, DATE '2025-09-01' AS to_date
)
*/

WITH src_dedup AS (
  SELECT
    f.title,
    f.link,
    f.source_id,
    f.words,
    f.published,
    COALESCE(f.feed_date, f.published::date) AS fd,
    f.created,
    f.updated,
    f.search_vector,
    f.category_id,
    ROW_NUMBER() OVER (
      PARTITION BY f.source_id, f.link, COALESCE(f.feed_date, f.published::date)
      ORDER BY f.updated DESC NULLS LAST, f.created DESC NULLS LAST, f.id DESC
    ) AS rn
  FROM src.feeds f
  /* OPTIONAL limit by date window:
  JOIN date_window dw ON COALESCE(f.feed_date, f.published::date) >= dw.from_date
                     AND COALESCE(f.feed_date, f.published::date) <  dw.to_date
  */
)
INSERT INTO public.feeds AS t
  (title, link, source_id, words, published, feed_date, created, updated, search_vector, category_id)
SELECT
  title, link, source_id, words, published, fd, created, updated, search_vector, category_id
FROM src_dedup
WHERE rn = 1
ON CONFLICT (source_id, link, feed_date) DO UPDATE
SET title       = EXCLUDED.title,
    words       = EXCLUDED.words,
    published   = EXCLUDED.published,
    updated     = EXCLUDED.updated,
    category_id = EXCLUDED.category_id;

COMMIT;

/* ----------------------------------------------------------------------
   6) Prepare upsert key for feed_sentiments (PG14 style)
   ---------------------------------------------------------------------- */
-- One-time (or safe to re-run). Ensures we can ON CONFLICT cleanly:
CREATE UNIQUE INDEX IF NOT EXISTS uq_feed_sent
  ON public.feed_sentiments (feed_id, model_id, feed_date);

/* If this errors due to prior duplicates in target, de-dup like so:
WITH d AS (
  SELECT ctid,
         ROW_NUMBER() OVER (PARTITION BY feed_id, model_id, feed_date
                            ORDER BY updated DESC NULLS LAST, created DESC NULLS LAST) rn
  FROM public.feed_sentiments
)
DELETE FROM public.feed_sentiments p
USING d
WHERE p.ctid = d.ctid AND d.rn > 1;
*/

/* ----------------------------------------------------------------------
   7) UPSERT feed_sentiments (deduped; joined to target feeds)
   ---------------------------------------------------------------------- */
BEGIN;
SET LOCAL synchronous_commit = off;

/* OPTIONAL: same date window as feeds
WITH date_window AS (
  SELECT DATE '2025-08-01' AS from_date, DATE '2025-09-01' AS to_date
)
*/

WITH s_join AS (
  SELECT
    f_dst.id AS feed_id,
    s.model_id,
    s.sentiments,
    s.sentiment_key,
    s.sentiment_value,
    s.updated,
    s.created,
    s.sentiment_compound,
    COALESCE(s.feed_date, f_dst.feed_date) AS fd,
    ROW_NUMBER() OVER (
      PARTITION BY f_dst.id, s.model_id, COALESCE(s.feed_date, f_dst.feed_date)
      ORDER BY s.updated DESC NULLS LAST, s.created DESC NULLS LAST, s.id DESC
    ) AS rn
  FROM src.feed_sentiments s
  JOIN src.feeds f_src
    ON f_src.id = s.feed_id
  JOIN public.feeds f_dst
    ON f_dst.source_id = f_src.source_id
   AND f_dst.link      = f_src.link
   AND f_dst.feed_date = COALESCE(f_src.feed_date, f_src.published::date)
  /* OPTIONAL date window:
  JOIN date_window dw ON COALESCE(s.feed_date, f_src.published::date) >= dw.from_date
                     AND COALESCE(s.feed_date, f_src.published::date) <  dw.to_date
  */
)
INSERT INTO public.feed_sentiments AS tgt
  (feed_id, model_id, sentiments, sentiment_key, sentiment_value, updated, created, sentiment_compound, feed_date)
SELECT
  feed_id, model_id, sentiments, sentiment_key, sentiment_value, updated, created, sentiment_compound, fd
FROM s_join
WHERE rn = 1
ON CONFLICT (feed_id, model_id, feed_date) DO UPDATE
SET sentiments         = EXCLUDED.sentiments,
    sentiment_key      = EXCLUDED.sentiment_key,
    sentiment_value    = EXCLUDED.sentiment_value,
    updated            = EXCLUDED.updated,
    created            = COALESCE(tgt.created, EXCLUDED.created),
    sentiment_compound = EXCLUDED.sentiment_compound;

COMMIT;

/* ----------------------------------------------------------------------
   8) Validate & analyze
   ---------------------------------------------------------------------- */
SELECT COUNT(*) AS dst_feeds FROM public.feeds;
SELECT COUNT(*) AS dst_sents FROM public.feed_sentiments;

ANALYZE public.feeds;
ANALYZE public.feed_sentiments;

/* ----------------------------------------------------------------------
   9) Optional — keep FDW for future syncs, or clean up
   ---------------------------------------------------------------------- */
-- Keep:
--   (nothing to do)
-- Or drop foreign tables & server if this was a one-off:
-- DROP SCHEMA src CASCADE;
-- DROP SERVER src_db CASCADE;


/* ----------------------------------------------------------------------
   10) IF you have all the data and you only want to import a specific month
       (e.g. August 2025), you can run the steps below.
       Adjust the date window in both steps as needed.
   ---------------------------------------------------------------------- */
BEGIN;
SET LOCAL synchronous_commit = off;

WITH date_window AS (
  SELECT DATE '2025-08-01' AS from_date, DATE '2025-09-01' AS to_date
),
src_dedup AS (
  SELECT
    f.title,
    f.link,
    f.source_id,
    f.words,
    f.published,
    COALESCE(f.feed_date, f.published::date) AS fd,
    f.created,
    f.updated,
    f.search_vector,
    f.category_id,
    ROW_NUMBER() OVER (
      PARTITION BY f.source_id, f.link, COALESCE(f.feed_date, f.published::date)
      ORDER BY f.updated DESC NULLS LAST, f.created DESC NULLS LAST, f.id DESC
    ) AS rn
  FROM src.feeds f
  JOIN date_window dw
    ON COALESCE(f.feed_date, f.published::date) >= dw.from_date
   AND COALESCE(f.feed_date, f.published::date) <  dw.to_date
)
INSERT INTO public.feeds AS t
  (title, link, source_id, words, published, feed_date, created, updated, search_vector, category_id)
SELECT
  title, link, source_id, words, published, fd, created, updated, search_vector, category_id
FROM src_dedup
WHERE rn = 1
ON CONFLICT (source_id, link, feed_date) DO UPDATE
SET title       = EXCLUDED.title,
    words       = EXCLUDED.words,
    published   = EXCLUDED.published,
    updated     = EXCLUDED.updated,
    category_id = EXCLUDED.category_id;

COMMIT;


BEGIN;
SET LOCAL synchronous_commit = off;

WITH date_window AS (
  SELECT DATE '2025-08-01' AS from_date, DATE '2025-09-01' AS to_date
),
s_join AS (
  SELECT
    f_dst.id AS feed_id,
    s.model_id,
    s.sentiments,
    s.sentiment_key,
    s.sentiment_value,
    s.updated,
    s.created,
    s.sentiment_compound,
    COALESCE(s.feed_date, f_dst.feed_date) AS fd,
    ROW_NUMBER() OVER (
      PARTITION BY f_dst.id, s.model_id, COALESCE(s.feed_date, f_dst.feed_date)
      ORDER BY s.updated DESC NULLS LAST, s.created DESC NULLS LAST, s.id DESC
    ) AS rn
  FROM src.feed_sentiments s
  JOIN src.feeds f_src
    ON f_src.id = s.feed_id
  JOIN public.feeds f_dst
    ON f_dst.source_id = f_src.source_id
   AND f_dst.link      = f_src.link
   AND f_dst.feed_date = COALESCE(f_src.feed_date, f_src.published::date)
  JOIN date_window dw
    ON COALESCE(s.feed_date, f_src.feed_date, f_src.published::date) >= dw.from_date
   AND COALESCE(s.feed_date, f_src.feed_date, f_src.published::date) <  dw.to_date
)
INSERT INTO public.feed_sentiments AS tgt
  (feed_id, model_id, sentiments, sentiment_key, sentiment_value, updated, created, sentiment_compound, feed_date)
SELECT
  feed_id, model_id, sentiments, sentiment_key, sentiment_value, updated, created, sentiment_compound, fd
FROM s_join
WHERE rn = 1
ON CONFLICT (feed_id, model_id, feed_date) DO UPDATE
SET sentiments         = EXCLUDED.sentiments,
    sentiment_key      = EXCLUDED.sentiment_key,
    sentiment_value    = EXCLUDED.sentiment_value,
    updated            = EXCLUDED.updated,
    created            = COALESCE(tgt.created, EXCLUDED.created),
    sentiment_compound = EXCLUDED.sentiment_compound;

COMMIT;
