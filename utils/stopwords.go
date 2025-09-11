package utils

import "strings"

var StopWordsList = []string{
	"a", "ahogy", "ahol", "aki", "akik", "akkor",
	"alatt", "által", "általában", "amely", "amelyek",
	"amelyekben", "amelyeket", "amelyet", "amelynek",
	"ami", "amit", "amolyan", "amíg", "amikor", "át",
	"abban", "ahhoz", "annak", "arra", "arról", "az",
	"azok", "azon", "azt", "azzal", "azért", "aztán",
	"azután", "azonban", "bár", "be", "belül", "benne",
	"cikk", "cikkek", "cikkeket", "csak", "de", "e",
	"eddig", "egész", "egy", "egyes", "egyetlen",
	"egyéb", "egyik", "egyre", "ekkor", "el", "elég",
	"ellen", "elõ", "elõször", "elõtt", "elsõ", "én",
	"éppen", "ebben", "ehhez", "emilyen", "ennek", "erre",
	"ez", "ezt", "ezek", "ezen", "ezzel", "ezért", "és",
	"fel", "felé", "hanem", "hiszen", "hogy", "hogyan",
	"igen", "így", "illetve", "ill.", "ill", "ilyen", "ilyenkor",
	"ison", "ismét", "itt", "jó", "jól", "jobban", "kell",
	"kellett", "keresztül", "keressünk", "ki", "kívül",
	"között", "közül", "legalább", "lehet", "lehetett",
	"legyen", "lenne", "lenni", "lesz", "lett", "maga",
	"magát", "majd", "majd", "már", "más", "másik", "meg",
	"még", "mellett", "mert", "mely", "melyek", "mi", "mit",
	"míg", "miért", "milyen", "mikor", "minden", "mindent",
	"mindenki", "mindig", "mint", "mintha", "mivel", "most",
	"nagy", "nagyobb", "nagyon", "ne", "néha", "nekem", "neki",
	"nem", "néhány", "nélkül", "nincs", "olyan", "ott", "össze",
	"õ", "õk", "õket", "pedig", "persze", "rá", "s", "saját",
	"sem", "semmi", "sok", "sokat", "sokkal", "számára",
	"szemben", "szerint", "szinte", "talán", "tehát",
	"teljes", "tovább", "továbbá", "több", "úgy", "ugyanis",
	"új", "újabb", "újra", "után", "utána", "utolsó", "vagy",
	"vagyis", "valaki", "valami", "valamint", "való", "vagyok",
	"van", "vannak", "volt", "voltam", "voltak", "voltunk",
	"vissza", "vele", "viszont", "volna",
}

var stopwords = map[string]struct{}{
	"a": {}, "ahogy": {}, "ahol": {}, "aki": {}, "akik": {}, "akkor": {}, "alatt": {},
	"amely": {}, "amelyek": {}, "amelyekben": {}, "amelyeket": {}, "amelyet": {}, "amelynek": {},
	"ami": {}, "amit": {}, "amolyan": {}, "amíg": {}, "amikor": {}, "át": {}, "abban": {},
	"ahhoz": {}, "annak": {}, "arra": {}, "arról": {}, "az": {}, "azok": {}, "azon": {},
	"azt": {}, "azzal": {}, "azért": {}, "aztán": {}, "azután": {}, "azóta": {},
	"bár": {}, "be": {}, "belül": {}, "benne": {}, "bizonyos": {}, "cikk": {}, "cikkek": {},
	"cikkeket": {}, "csak": {}, "de": {}, "eddig": {}, "egész": {}, "egy": {}, "egyes": {},
	"egyetlen": {}, "egyik": {}, "egyre": {}, "ekkor": {}, "el": {}, "elég": {}, "ellen": {},
	"elő": {}, "először": {}, "előtt": {}, "első": {}, "én": {}, "éppen": {}, "ebben": {},
	"ehhez": {}, "emilyen": {}, "ennek": {}, "erre": {}, "ez": {}, "ezt": {}, "ezek": {},
	"ezen": {}, "ezzel": {}, "ezért": {}, "ezúttal": {}, "és": {}, "fel": {}, "felé": {}, "hanem": {},
	"hiszen": {}, "hogy": {}, "hogyan": {}, "ide": {}, "igen": {}, "így": {}, "illetve": {},
	"ill": {}, "ilyen": {}, "ilyenkor": {}, "is": {}, "ismét": {}, "itt": {}, "jó": {}, "jól": {},
	"jobban": {}, "kell": {}, "kellett": {}, "keresztül": {}, "keressünk": {}, "ki": {}, "kívül": {},
	"között": {}, "közül": {}, "legalább": {}, "lehet": {}, "lehetett": {}, "legyen": {}, "lenne": {},
	"lenni": {}, "lesz": {}, "lett": {}, "maga": {}, "magát": {}, "majd": {}, "meg": {},
	"még": {}, "mellé": {}, "mellett": {}, "mert": {}, "mi": {}, "mit": {}, "míg": {}, "miért": {},
	"milyen": {}, "minden": {}, "mindent": {}, "mindenki": {}, "mindig": {}, "mint": {}, "mintha": {},
	"mivel": {}, "most": {}, "már": {}, "más": {}, "másik": {}, "másként": {}, "másnap": {}, "mások": {},
	"megint": {}, "mégis": {}, "nagy": {}, "nagyobb": {}, "nagyon": {}, "ne": {}, "néha": {},
	"nekem": {}, "neki": {}, "nem": {}, "néhány": {}, "nélkül": {}, "nincs": {}, "nézzük": {},
	"oda": {}, "olyan": {}, "ott": {}, "össze": {}, "őt": {}, "őket": {}, "ő": {}, "ők": {},
	"pedig": {}, "persze": {}, "rá": {}, "s": {}, "saját": {}, "sem": {}, "semmi": {}, "sok": {},
	"sokat": {}, "sokkal": {}, "számára": {}, "szemben": {}, "szerint": {}, "szinte": {},
	"talán": {}, "tehát": {}, "teljes": {}, "tovább": {}, "továbbá": {}, "több": {}, "túl": {},
	"úgy": {}, "új": {}, "újabb": {}, "újra": {}, "után": {}, "utána": {}, "utolsó": {}, "vagy": {},
	"vagyis": {}, "valaki": {}, "valami": {}, "valamint": {}, "való": {}, "van": {}, "vannak": {},
	"vele": {}, "vissza": {}, "volna": {}, "volt": {}, "voltak": {}, "voltunk": {}, "magyar": {}, "videó": {},
}

// IsStopword returns true if word is in STOPWORDS list
func IsStopword(word string) bool {
	_, exists := stopwords[strings.ToLower(word)]
	return exists
}
