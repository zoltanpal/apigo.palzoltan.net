package utils

import "strings"

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
	"ezen": {}, "ezzel": {}, "ezért": {}, "ezúttal": {}, "fel": {}, "felé": {}, "hanem": {},
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
	"vele": {}, "vissza": {}, "volna": {}, "volt": {}, "voltak": {}, "voltunk": {},
}

// IsStopword returns true if word is in STOPWORDS list
func IsStopword(word string) bool {
	_, exists := stopwords[strings.ToLower(word)]
	return exists
}
