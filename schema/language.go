package schema

import (
	"errors"
	"fmt"
	"ring/schema/entitytype"
	"strconv"
	"strings"
)

// Language is an ISO 639-1 language with code, name and native name.
type Language struct {
	id             int32
	code           string
	name           string
	nativeName     string
	currentCountry *country
	isValid        bool
}

type country struct {
	id   uint16
	code string
	name string
}

const (
	languageName                string = "@language"
	languageErrorEmpty          string = "empty code"
	languageErrorInvalidCountry string = "invalid country code '%s'"
	languageErrorInvalidLang    string = "invalid language code '%s'"
)

//TODO add countries
var languages = map[string]Language{
	"aa": {id: 11110, code: "aa", name: "Afar", nativeName: "Afaraf"},
	"ab": {id: 11210, code: "ab", name: "Abkhaz", nativeName: "аҧсуа бызшәа"},
	"ae": {id: 11310, code: "ae", name: "Avestan", nativeName: "avesta"},
	"af": {id: 11410, code: "af", name: "Afrikaans", nativeName: "Afrikaans"},
	"ak": {id: 11510, code: "ak", name: "Akan", nativeName: "Akan"},
	"am": {id: 11610, code: "am", name: "Amharic", nativeName: "አማርኛ"},
	"an": {id: 11710, code: "an", name: "Aragonese", nativeName: "aragonés"},
	"ar": {id: 11810, code: "ar", name: "Arabic", nativeName: "اللغة العربية"},
	"as": {id: 11910, code: "as", name: "Assamese", nativeName: "অসমীয়া"},
	"av": {id: 12010, code: "av", name: "Avaric", nativeName: "авар мацӀ"},
	"ay": {id: 12110, code: "ay", name: "Aymara", nativeName: "aymar aru"},
	"az": {id: 12210, code: "az", name: "Azerbaijani", nativeName: "azərbaycan dili"},
	"ba": {id: 12310, code: "ba", name: "Bashkir", nativeName: "башҡорт теле"},
	"be": {id: 12410, code: "be", name: "Belarusian", nativeName: "беларуская мова"},
	"bg": {id: 12510, code: "bg", name: "Bulgarian", nativeName: "български език"},
	"bh": {id: 12610, code: "bh", name: "Bihari", nativeName: "भोजपुरी"},
	"bi": {id: 12710, code: "bi", name: "Bislama", nativeName: "Bislama"},
	"bm": {id: 12810, code: "bm", name: "Bambara", nativeName: "bamanankan"},
	"bn": {id: 12910, code: "bn", name: "Bengali", nativeName: "বাংলা"},
	"bo": {id: 13010, code: "bo", name: "Tibetan Standard", nativeName: "བོད་ཡིག"},
	"br": {id: 13110, code: "br", name: "Breton", nativeName: "brezhoneg"},
	"bs": {id: 13210, code: "bs", name: "Bosnian", nativeName: "bosanski jezik"},
	"ca": {id: 13310, code: "ca", name: "Catalan", nativeName: "català"},
	"ce": {id: 13410, code: "ce", name: "Chechen", nativeName: "нохчийн мотт"},
	"ch": {id: 13510, code: "ch", name: "Chamorro", nativeName: "Chamoru"},
	"co": {id: 13610, code: "co", name: "Corsican", nativeName: "corsu"},
	"cr": {id: 13710, code: "cr", name: "Cree", nativeName: "ᓀᐦᐃᔭᐍᐏᐣ"},
	"cs": {id: 13810, code: "cs", name: "Czech", nativeName: "čeština"},
	"cu": {id: 13910, code: "cu", name: "Old Church Slavonic", nativeName: "ѩзыкъ словѣньскъ"},
	"cv": {id: 14010, code: "cv", name: "Chuvash", nativeName: "чӑваш чӗлхи"},
	"cy": {id: 14110, code: "cy", name: "Welsh", nativeName: "Cymraeg"},
	"da": {id: 14210, code: "da", name: "Danish", nativeName: "dansk"},
	"de": {id: 14310, code: "de", name: "German", nativeName: "Deutsch"},
	"dv": {id: 14410, code: "dv", name: "Divehi", nativeName: "Dhivehi"},
	"dz": {id: 14510, code: "dz", name: "Dzongkha", nativeName: "རྫོང་ཁ"},
	"ee": {id: 14610, code: "ee", name: "Ewe", nativeName: "Eʋegbe"},
	"el": {id: 14710, code: "el", name: "Greek", nativeName: "Ελληνικά"},
	"en": {id: 14810, code: "en", name: "English", nativeName: "English"},
	"eo": {id: 14910, code: "eo", name: "Esperanto", nativeName: "Esperanto"},
	"es": {id: 15010, code: "es", name: "Spanish", nativeName: "Español"},
	"et": {id: 15110, code: "et", name: "Estonian", nativeName: "eesti"},
	"eu": {id: 15210, code: "eu", name: "Basque", nativeName: "euskara"},
	"fa": {id: 15310, code: "fa", name: "Persian", nativeName: "فارسی"},
	"ff": {id: 15410, code: "ff", name: "Fula", nativeName: "Fulfulde"},
	"fi": {id: 15510, code: "fi", name: "Finnish", nativeName: "suomi"},
	"fj": {id: 15610, code: "fj", name: "Fijian", nativeName: "Vakaviti"},
	"fo": {id: 15710, code: "fo", name: "Faroese", nativeName: "føroyskt"},
	"fr": {id: 15810, code: "fr", name: "French", nativeName: "Français"},
	"fy": {id: 15910, code: "fy", name: "Western Frisian", nativeName: "Frysk"},
	"ga": {id: 16010, code: "ga", name: "Irish", nativeName: "Gaeilge"},
	"gd": {id: 16110, code: "gd", name: "Scottish Gaelic", nativeName: "Gàidhlig"},
	"gl": {id: 16210, code: "gl", name: "Galician", nativeName: "galego"},
	"gn": {id: 16310, code: "gn", name: "Guaraní", nativeName: "Avañeẽ"},
	"gu": {id: 16410, code: "gu", name: "Gujarati", nativeName: "ગુજરાતી"},
	"gv": {id: 16510, code: "gv", name: "Manx", nativeName: "Gaelg"},
	"ha": {id: 16610, code: "ha", name: "Hausa", nativeName: "هَوُسَ"},
	"he": {id: 16710, code: "he", name: "Hebrew", nativeName: "עברית"},
	"hi": {id: 16810, code: "hi", name: "Hindi", nativeName: "हिन्दी"},
	"ho": {id: 16910, code: "ho", name: "Hiri Motu", nativeName: "Hiri Motu"},
	"hr": {id: 17010, code: "hr", name: "Croatian", nativeName: "hrvatski jezik"},
	"ht": {id: 17110, code: "ht", name: "Haitian", nativeName: "Kreyòl ayisyen"},
	"hu": {id: 17210, code: "hu", name: "Hungarian", nativeName: "magyar"},
	"hy": {id: 17310, code: "hy", name: "Armenian", nativeName: "Հայերեն"},
	"hz": {id: 17410, code: "hz", name: "Herero", nativeName: "Otjiherero"},
	"ia": {id: 17510, code: "ia", name: "Interlingua", nativeName: "Interlingua"},
	"id": {id: 17610, code: "id", name: "Indonesian", nativeName: "Indonesian"},
	"ie": {id: 17710, code: "ie", name: "Interlingue", nativeName: "Interlingue"},
	"ig": {id: 17810, code: "ig", name: "Igbo", nativeName: "Asụsụ Igbo"},
	"ii": {id: 17910, code: "ii", name: "Nuosu", nativeName: "ꆈꌠ꒿ Nuosuhxop"},
	"ik": {id: 18010, code: "ik", name: "Inupiaq", nativeName: "Iñupiaq"},
	"io": {id: 18110, code: "io", name: "Ido", nativeName: "Ido"},
	"is": {id: 18210, code: "is", name: "Icelandic", nativeName: "Íslenska"},
	"it": {id: 18310, code: "it", name: "Italian", nativeName: "Italiano"},
	"iu": {id: 18410, code: "iu", name: "Inuktitut", nativeName: "ᐃᓄᒃᑎᑐᑦ"},
	"ja": {id: 18510, code: "ja", name: "Japanese", nativeName: "日本語"},
	"jv": {id: 18610, code: "jv", name: "Javanese", nativeName: "basa Jawa"},
	"ka": {id: 18710, code: "ka", name: "Georgian", nativeName: "ქართული"},
	"kg": {id: 18810, code: "kg", name: "Kongo", nativeName: "Kikongo"},
	"ki": {id: 18910, code: "ki", name: "Kikuyu", nativeName: "Gĩkũyũ"},
	"kj": {id: 19010, code: "kj", name: "Kwanyama", nativeName: "Kuanyama"},
	"kk": {id: 19110, code: "kk", name: "Kazakh", nativeName: "қазақ тілі"},
	"kl": {id: 19210, code: "kl", name: "Kalaallisut", nativeName: "kalaallisut"},
	"km": {id: 19310, code: "km", name: "Khmer", nativeName: "ខេមរភាសា"},
	"kn": {id: 19410, code: "kn", name: "Kannada", nativeName: "ಕನ್ನಡ"},
	"ko": {id: 19510, code: "ko", name: "Korean", nativeName: "한국어"},
	"kr": {id: 19610, code: "kr", name: "Kanuri", nativeName: "Kanuri"},
	"ks": {id: 19710, code: "ks", name: "Kashmiri", nativeName: "कश्मीरी"},
	"ku": {id: 19810, code: "ku", name: "Kurdish", nativeName: "Kurdî"},
	"kv": {id: 19910, code: "kv", name: "Komi", nativeName: "коми кыв"},
	"kw": {id: 20010, code: "kw", name: "Cornish", nativeName: "Kernewek"},
	"ky": {id: 20110, code: "ky", name: "Kyrgyz", nativeName: "Кыргызча"},
	"la": {id: 20210, code: "la", name: "Latin", nativeName: "latine"},
	"lb": {id: 20310, code: "lb", name: "Luxembourgish", nativeName: "Lëtzebuergesch"},
	"lg": {id: 20410, code: "lg", name: "Ganda", nativeName: "Luganda"},
	"li": {id: 20510, code: "li", name: "Limburgish", nativeName: "Limburgs"},
	"ln": {id: 20610, code: "ln", name: "Lingala", nativeName: "Lingála"},
	"lo": {id: 20710, code: "lo", name: "Lao", nativeName: "ພາສາ"},
	"lt": {id: 20810, code: "lt", name: "Lithuanian", nativeName: "lietuvių kalba"},
	"lu": {id: 20910, code: "lu", name: "Luba-Katanga", nativeName: "Tshiluba"},
	"lv": {id: 21010, code: "lv", name: "Latvian", nativeName: "latviešu valoda"},
	"mg": {id: 21110, code: "mg", name: "Malagasy", nativeName: "fiteny malagasy"},
	"mh": {id: 21210, code: "mh", name: "Marshallese", nativeName: "Kajin M̧ajeļ"},
	"mi": {id: 21310, code: "mi", name: "Māori", nativeName: "te reo Māori"},
	"mk": {id: 21410, code: "mk", name: "Macedonian", nativeName: "македонски јазик"},
	"ml": {id: 21510, code: "ml", name: "Malayalam", nativeName: "മലയാളം"},
	"mn": {id: 21610, code: "mn", name: "Mongolian", nativeName: "Монгол хэл"},
	"mr": {id: 21710, code: "mr", name: "Marathi", nativeName: "मराठी"},
	"ms": {id: 21810, code: "ms", name: "Malay", nativeName: "هاس ملايو‎"},
	"mt": {id: 21910, code: "mt", name: "Maltese", nativeName: "Malti"},
	"my": {id: 22010, code: "my", name: "Burmese", nativeName: "ဗမာစာ"},
	"na": {id: 22110, code: "na", name: "Nauru", nativeName: "Ekakairũ Naoero"},
	"nb": {id: 22210, code: "nb", name: "Norwegian Bokmål", nativeName: "Norsk bokmål"},
	"nd": {id: 22310, code: "nd", name: "Northern Ndebele", nativeName: "isiNdebele"},
	"ne": {id: 22410, code: "ne", name: "Nepali", nativeName: "नेपाली"},
	"ng": {id: 22510, code: "ng", name: "Ndonga", nativeName: "Owambo"},
	"nl": {id: 22610, code: "nl", name: "Dutch", nativeName: "Nederlands"},
	"nn": {id: 22710, code: "nn", name: "Norwegian Nynorsk", nativeName: "Norsk nynorsk"},
	"no": {id: 22810, code: "no", name: "Norwegian", nativeName: "Norsk"},
	"nr": {id: 22910, code: "nr", name: "Southern Ndebele", nativeName: "isiNdebele"},
	"nv": {id: 23010, code: "nv", name: "Navajo", nativeName: "Diné bizaad"},
	"ny": {id: 23110, code: "ny", name: "Chichewa", nativeName: "chiCheŵa"},
	"oc": {id: 23210, code: "oc", name: "Occitan", nativeName: "occitan"},
	"oj": {id: 23310, code: "oj", name: "Ojibwe", nativeName: "ᐊᓂᔑᓈᐯᒧᐎᓐ"},
	"om": {id: 23410, code: "om", name: "Oromo", nativeName: "Afaan Oromoo"},
	"or": {id: 23510, code: "or", name: "Oriya", nativeName: "ଓଡ଼ିଆ"},
	"os": {id: 23610, code: "os", name: "Ossetian", nativeName: "ирон æвзаг"},
	"pa": {id: 23710, code: "pa", name: "Panjabi", nativeName: "ਪੰਜਾਬੀ"},
	"pi": {id: 23810, code: "pi", name: "Pāli", nativeName: "पाऴि"},
	"pl": {id: 23910, code: "pl", name: "Polish", nativeName: "język polski"},
	"ps": {id: 24010, code: "ps", name: "Pashto", nativeName: "پښتو"},
	"pt": {id: 24110, code: "pt", name: "Portuguese", nativeName: "Português"},
	"qu": {id: 24210, code: "qu", name: "Quechua", nativeName: "Runa Simi"},
	"rm": {id: 24310, code: "rm", name: "Romansh", nativeName: "rumantsch grischun"},
	"rn": {id: 24410, code: "rn", name: "Kirundi", nativeName: "Ikirundi"},
	"ro": {id: 24510, code: "ro", name: "Romanian", nativeName: "Română"},
	"ru": {id: 24610, code: "ru", name: "Russian", nativeName: "Русский"},
	"rw": {id: 24710, code: "rw", name: "Kinyarwanda", nativeName: "Ikinyarwanda"},
	"sa": {id: 24810, code: "sa", name: "Sanskrit", nativeName: "संस्कृतम्"},
	"sc": {id: 24910, code: "sc", name: "Sardinian", nativeName: "sardu"},
	"sd": {id: 25010, code: "sd", name: "Sindhi", nativeName: "सिन्धी"},
	"se": {id: 25110, code: "se", name: "Northern Sami", nativeName: "Davvisámegiella"},
	"sg": {id: 25210, code: "sg", name: "Sango", nativeName: "yângâ tî sängö"},
	"si": {id: 25310, code: "si", name: "Sinhala", nativeName: "සිංහල"},
	"sk": {id: 25410, code: "sk", name: "Slovak", nativeName: "slovenčina"},
	"sl": {id: 25510, code: "sl", name: "Slovene", nativeName: "slovenski jezik"},
	"sm": {id: 25610, code: "sm", name: "Samoan", nativeName: "gagana faa Samoa"},
	"sn": {id: 25710, code: "sn", name: "Shona", nativeName: "chiShona"},
	"so": {id: 25810, code: "so", name: "Somali", nativeName: "Soomaaliga"},
	"sq": {id: 25910, code: "sq", name: "Albanian", nativeName: "Shqip"},
	"sr": {id: 26010, code: "sr", name: "Serbian", nativeName: "српски језик"},
	"ss": {id: 26110, code: "ss", name: "Swati", nativeName: "SiSwati"},
	"st": {id: 26210, code: "st", name: "Southern Sotho", nativeName: "Sesotho"},
	"su": {id: 26310, code: "su", name: "Sundanese", nativeName: "Basa Sunda"},
	"sv": {id: 26410, code: "sv", name: "Swedish", nativeName: "svenska"},
	"sw": {id: 26510, code: "sw", name: "Swahili", nativeName: "Kiswahili"},
	"ta": {id: 26610, code: "ta", name: "Tamil", nativeName: "தமிழ்"},
	"te": {id: 26710, code: "te", name: "Telugu", nativeName: "తెలుగు"},
	"tg": {id: 26810, code: "tg", name: "Tajik", nativeName: "тоҷикӣ"},
	"th": {id: 26910, code: "th", name: "Thai", nativeName: "ไทย"},
	"ti": {id: 27010, code: "ti", name: "Tigrinya", nativeName: "ትግርኛ"},
	"tk": {id: 27110, code: "tk", name: "Turkmen", nativeName: "Türkmen"},
	"tl": {id: 27210, code: "tl", name: "Tagalog", nativeName: "Wikang Tagalog"},
	"tn": {id: 27310, code: "tn", name: "Tswana", nativeName: "Setswana"},
	"to": {id: 27410, code: "to", name: "Tonga", nativeName: "faka Tonga"},
	"tr": {id: 27510, code: "tr", name: "Turkish", nativeName: "Türkçe"},
	"ts": {id: 27610, code: "ts", name: "Tsonga", nativeName: "Xitsonga"},
	"tt": {id: 27710, code: "tt", name: "Tatar", nativeName: "татар теле"},
	"tw": {id: 27810, code: "tw", name: "Twi", nativeName: "Twi"},
	"ty": {id: 27910, code: "ty", name: "Tahitian", nativeName: "Reo Tahiti"},
	"ug": {id: 28010, code: "ug", name: "Uyghur", nativeName: "ئۇيغۇرچە‎"},
	"uk": {id: 28110, code: "uk", name: "Ukrainian", nativeName: "Українська"},
	"ur": {id: 28210, code: "ur", name: "Urdu", nativeName: "اردو"},
	"uz": {id: 28310, code: "uz", name: "Uzbek", nativeName: "Ўзбек"},
	"ve": {id: 28410, code: "ve", name: "Venda", nativeName: "Tshivenḓa"},
	"vi": {id: 28510, code: "vi", name: "Vietnamese", nativeName: "Tiếng Việt"},
	"vo": {id: 28610, code: "vo", name: "Volapük", nativeName: "Volapük"},
	"wa": {id: 28710, code: "wa", name: "Walloon", nativeName: "walon"},
	"wo": {id: 28810, code: "wo", name: "Wolof", nativeName: "Wollof"},
	"xh": {id: 28910, code: "xh", name: "Xhosa", nativeName: "isiXhosa"},
	"yi": {id: 29010, code: "yi", name: "Yiddish", nativeName: "ייִדיש"},
	"yo": {id: 29110, code: "yo", name: "Yoruba", nativeName: "Yorùbá"},
	"za": {id: 29210, code: "za", name: "Zhuang", nativeName: "Saɯ cueŋƅ"},
	"zh": {id: 29310, code: "zh", name: "Chinese", nativeName: "中文"},
	"zu": {id: 29410, code: "zu", name: "Zulu", nativeName: "isiZulu"},
}

var countries = map[string]country{
	"AF": {id: 1110, code: "AF", name: "Afghanistan"},
	"AX": {id: 1120, code: "AX", name: "Åland Islands"},
	"AL": {id: 1130, code: "AL", name: "Albania"},
	"DZ": {id: 1140, code: "DZ", name: "Algeria"},
	"AS": {id: 1150, code: "AS", name: "American Samoa"},
	"AD": {id: 1160, code: "AD", name: "Andorra"},
	"AO": {id: 1170, code: "AO", name: "Angola"},
	"AI": {id: 1180, code: "AI", name: "Anguilla"},
	"AQ": {id: 1190, code: "AQ", name: "Antarctica"},
	"AG": {id: 1200, code: "AG", name: "Antigua and Barbuda"},
	"AR": {id: 1210, code: "AR", name: "Argentina"},
	"AM": {id: 1220, code: "AM", name: "Armenia"},
	"AW": {id: 1230, code: "AW", name: "Aruba"},
	"AU": {id: 1240, code: "AU", name: "Australia"},
	"AT": {id: 1250, code: "AT", name: "Austria"},
	"AZ": {id: 1260, code: "AZ", name: "Azerbaijan"},
	"BH": {id: 1270, code: "BH", name: "Bahrain"},
	"BS": {id: 1280, code: "BS", name: "Bahamas"},
	"BD": {id: 1290, code: "BD", name: "Bangladesh"},
	"BB": {id: 1300, code: "BB", name: "Barbados"},
	"BY": {id: 1310, code: "BY", name: "Belarus"},
	"BE": {id: 1320, code: "BE", name: "Belgium"},
	"BZ": {id: 1330, code: "BZ", name: "Belize"},
	"BJ": {id: 1340, code: "BJ", name: "Benin"},
	"BM": {id: 1350, code: "BM", name: "Bermuda"},
	"BT": {id: 1360, code: "BT", name: "Bhutan"},
	"BO": {id: 1370, code: "BO", name: "Bolivia, Plurinational State of"},
	"BQ": {id: 1380, code: "BQ", name: "Bonaire, Sint Eustatius and Saba"},
	"BA": {id: 1390, code: "BA", name: "Bosnia and Herzegovina"},
	"BW": {id: 1400, code: "BW", name: "Botswana"},
	"BV": {id: 1410, code: "BV", name: "Bouvet Island"},
	"BR": {id: 1420, code: "BR", name: "Brazil"},
	"IO": {id: 1430, code: "IO", name: "British Indian Ocean Territory"},
	"BN": {id: 1440, code: "BN", name: "Brunei Darussalam"},
	"BG": {id: 1450, code: "BG", name: "Bulgaria"},
	"BF": {id: 1460, code: "BF", name: "Burkina Faso"},
	"BI": {id: 1470, code: "BI", name: "Burundi"},
	"KH": {id: 1480, code: "KH", name: "Cambodia"},
	"CM": {id: 1490, code: "CM", name: "Cameroon"},
	"CA": {id: 1500, code: "CA", name: "Canada"},
	"CV": {id: 1510, code: "CV", name: "Cape Verde"},
	"KY": {id: 1520, code: "KY", name: "Cayman Islands"},
	"CF": {id: 1530, code: "CF", name: "Central African Republic"},
	"TD": {id: 1540, code: "TD", name: "Chad"},
	"CL": {id: 1550, code: "CL", name: "Chile"},
	"CN": {id: 1560, code: "CN", name: "China"},
	"CX": {id: 1570, code: "CX", name: "Christmas Island"},
	"CC": {id: 1580, code: "CC", name: "Cocos (Keeling) Islands"},
	"CO": {id: 1590, code: "CO", name: "Colombia"},
	"KM": {id: 1600, code: "KM", name: "Comoros"},
	"CG": {id: 1610, code: "CG", name: "Congo"},
	"CD": {id: 1620, code: "CD", name: "Congo"},
	"CK": {id: 1630, code: "CK", name: "Cook Islands"},
	"CR": {id: 1640, code: "CR", name: "Costa Rica"},
	"CI": {id: 1650, code: "CI", name: "Côte d'Ivoire"},
	"HR": {id: 1660, code: "HR", name: "Croatia"},
	"CU": {id: 1670, code: "CU", name: "Cuba"},
	"CW": {id: 1680, code: "CW", name: "Curaçao"},
	"CY": {id: 1690, code: "CY", name: "Cyprus"},
	"CZ": {id: 1700, code: "CZ", name: "Czech Republic"},
	"DK": {id: 1710, code: "DK", name: "Denmark"},
	"DJ": {id: 1720, code: "DJ", name: "Djibouti"},
	"DM": {id: 1730, code: "DM", name: "Dominica"},
	"DO": {id: 1740, code: "DO", name: "Dominican Republic"},
	"EC": {id: 1750, code: "EC", name: "Ecuador"},
	"EG": {id: 1760, code: "EG", name: "Egypt"},
	"SV": {id: 1770, code: "SV", name: "El Salvador"},
	"GQ": {id: 1780, code: "GQ", name: "Equatorial Guinea"},
	"ER": {id: 1790, code: "ER", name: "Eritrea"},
	"EE": {id: 1800, code: "EE", name: "Estonia"},
	"ET": {id: 1810, code: "ET", name: "Ethiopia"},
	"FK": {id: 1820, code: "FK", name: "Falkland Islands (Malvinas)"},
	"FO": {id: 1830, code: "FO", name: "Faroe Islands"},
	"FJ": {id: 1840, code: "FJ", name: "Fiji"},
	"FI": {id: 1850, code: "FI", name: "Finland"},
	"FR": {id: 1860, code: "FR", name: "France"},
	"GF": {id: 1870, code: "GF", name: "French Guiana"},
	"PF": {id: 1880, code: "PF", name: "French Polynesia"},
	"TF": {id: 1890, code: "TF", name: "French Southern Territories"},
	"GA": {id: 1900, code: "GA", name: "Gabon"},
	"GM": {id: 1910, code: "GM", name: "Gambia"},
	"GE": {id: 1920, code: "GE", name: "Georgia"},
	"DE": {id: 1930, code: "DE", name: "Germany"},
	"GH": {id: 1940, code: "GH", name: "Ghana"},
	"GI": {id: 1950, code: "GI", name: "Gibraltar"},
	"GR": {id: 1960, code: "GR", name: "Greece"},
	"GL": {id: 1970, code: "GL", name: "Greenland"},
	"GD": {id: 1980, code: "GD", name: "Grenada"},
	"GP": {id: 1990, code: "GP", name: "Guadeloupe"},
	"GU": {id: 2000, code: "GU", name: "Guam"},
	"GT": {id: 2010, code: "GT", name: "Guatemala"},
	"GG": {id: 2020, code: "GG", name: "Guernsey"},
	"GN": {id: 2030, code: "GN", name: "Guinea"},
	"GW": {id: 2040, code: "GW", name: "Guinea-Bissau"},
	"GY": {id: 2050, code: "GY", name: "Guyana"},
	"HT": {id: 2060, code: "HT", name: "Haiti"},
	"HM": {id: 2070, code: "HM", name: "Heard Island and McDonald Islands"},
	"VA": {id: 2080, code: "VA", name: "Holy See (Vatican City State)"},
	"HN": {id: 2090, code: "HN", name: "Honduras"},
	"HK": {id: 2100, code: "HK", name: "Hong Kong"},
	"HU": {id: 2110, code: "HU", name: "Hungary"},
	"IS": {id: 2120, code: "IS", name: "Iceland"},
	"IN": {id: 2130, code: "IN", name: "India"},
	"ID": {id: 2140, code: "ID", name: "Indonesia"},
	"IR": {id: 2150, code: "IR", name: "Iran"},
	"IQ": {id: 2160, code: "IQ", name: "Iraq"},
	"IE": {id: 2170, code: "IE", name: "Ireland"},
	"IM": {id: 2180, code: "IM", name: "Isle of Man"},
	"IL": {id: 2190, code: "IL", name: "Israel"},
	"IT": {id: 2200, code: "IT", name: "Italy"},
	"JM": {id: 2210, code: "JM", name: "Jamaica"},
	"JP": {id: 2220, code: "JP", name: "Japan"},
	"JE": {id: 2230, code: "JE", name: "Jersey"},
	"JO": {id: 2240, code: "JO", name: "Jordan"},
	"KZ": {id: 2250, code: "KZ", name: "Kazakhstan"},
	"KE": {id: 2260, code: "KE", name: "Kenya"},
	"KI": {id: 2270, code: "KI", name: "Kiribati"},
	"KP": {id: 2280, code: "KP", name: "Korea"},
	"KR": {id: 2290, code: "KR", name: "Korea, Republic of"},
	"KW": {id: 2300, code: "KW", name: "Kuwait"},
	"KG": {id: 2310, code: "KG", name: "Kyrgyzstan"},
	"LA": {id: 2320, code: "LA", name: "Lao People's Democratic Republic"},
	"LV": {id: 2330, code: "LV", name: "Latvia"},
	"LB": {id: 2340, code: "LB", name: "Lebanon"},
	"LS": {id: 2350, code: "LS", name: "Lesotho"},
	"LR": {id: 2360, code: "LR", name: "Liberia"},
	"LY": {id: 2370, code: "LY", name: "Libya"},
	"LI": {id: 2380, code: "LI", name: "Liechtenstein"},
	"LT": {id: 2390, code: "LT", name: "Lithuania"},
	"LU": {id: 2400, code: "LU", name: "Luxembourg"},
	"MO": {id: 2410, code: "MO", name: "Macao"},
	"MK": {id: 2420, code: "MK", name: "North Macedonia"},
	"MG": {id: 2430, code: "MG", name: "Madagascar"},
	"MW": {id: 2440, code: "MW", name: "Malawi"},
	"MY": {id: 2450, code: "MY", name: "Malaysia"},
	"MV": {id: 2460, code: "MV", name: "Maldives"},
	"ML": {id: 2470, code: "ML", name: "Mali"},
	"MT": {id: 2480, code: "MT", name: "Malta"},
	"MH": {id: 2490, code: "MH", name: "Marshall Islands"},
	"MQ": {id: 2500, code: "MQ", name: "Martinique"},
	"MR": {id: 2510, code: "MR", name: "Mauritania"},
	"MU": {id: 2520, code: "MU", name: "Mauritius"},
	"YT": {id: 2530, code: "YT", name: "Mayotte"},
	"MX": {id: 2540, code: "MX", name: "Mexico"},
	"FM": {id: 2550, code: "FM", name: "Micronesia"},
	"MD": {id: 2560, code: "MD", name: "Moldova"},
	"MC": {id: 2570, code: "MC", name: "Monaco"},
	"MN": {id: 2580, code: "MN", name: "Mongolia"},
	"ME": {id: 2590, code: "ME", name: "Montenegro"},
	"MS": {id: 2600, code: "MS", name: "Montserrat"},
	"MA": {id: 2610, code: "MA", name: "Morocco"},
	"MZ": {id: 2620, code: "MZ", name: "Mozambique"},
	"MM": {id: 2630, code: "MM", name: "Myanmar"},
	"NA": {id: 2640, code: "NA", name: "Namibia"},
	"NR": {id: 2650, code: "NR", name: "Nauru"},
	"NP": {id: 2660, code: "NP", name: "Nepal"},
	"NL": {id: 2670, code: "NL", name: "Netherlands"},
	"NC": {id: 2680, code: "NC", name: "New Caledonia"},
	"NZ": {id: 2690, code: "NZ", name: "New Zealand"},
	"NI": {id: 2700, code: "NI", name: "Nicaragua"},
	"NE": {id: 2710, code: "NE", name: "Niger"},
	"NG": {id: 2720, code: "NG", name: "Nigeria"},
	"NU": {id: 2730, code: "NU", name: "Niue"},
	"NF": {id: 2740, code: "NF", name: "Norfolk Island"},
	"MP": {id: 2750, code: "MP", name: "Northern Mariana Islands"},
	"NO": {id: 2760, code: "NO", name: "Norway"},
	"OM": {id: 2770, code: "OM", name: "Oman"},
	"PK": {id: 2780, code: "PK", name: "Pakistan"},
	"PW": {id: 2790, code: "PW", name: "Palau"},
	"PS": {id: 2800, code: "PS", name: "Palestine"},
	"PA": {id: 2810, code: "PA", name: "Panama"},
	"PG": {id: 2820, code: "PG", name: "Papua New Guinea"},
	"PY": {id: 2830, code: "PY", name: "Paraguay"},
	"PE": {id: 2840, code: "PE", name: "Peru"},
	"PH": {id: 2850, code: "PH", name: "Philippines"},
	"PN": {id: 2860, code: "PN", name: "Pitcairn"},
	"PL": {id: 2870, code: "PL", name: "Poland"},
	"PT": {id: 2880, code: "PT", name: "Portugal"},
	"PR": {id: 2890, code: "PR", name: "Puerto Rico"},
	"QA": {id: 2900, code: "QA", name: "Qatar"},
	"RE": {id: 2910, code: "RE", name: "Réunion"},
	"RO": {id: 2920, code: "RO", name: "Romania"},
	"RU": {id: 2930, code: "RU", name: "Russian Federation"},
	"RW": {id: 2940, code: "RW", name: "Rwanda"},
	"BL": {id: 2950, code: "BL", name: "Saint Barthélemy"},
	"SH": {id: 2960, code: "SH", name: "Saint Helena, Ascension and Tristan da Cunha"},
	"KN": {id: 2970, code: "KN", name: "Saint Kitts and Nevis"},
	"LC": {id: 2980, code: "LC", name: "Saint Lucia"},
	"MF": {id: 2990, code: "MF", name: "Saint Martin (French part)"},
	"PM": {id: 3000, code: "PM", name: "Saint Pierre and Miquelon"},
	"VC": {id: 3010, code: "VC", name: "Saint Vincent and the Grenadines"},
	"WS": {id: 3020, code: "WS", name: "Samoa"},
	"SM": {id: 3030, code: "SM", name: "San Marino"},
	"ST": {id: 3040, code: "ST", name: "Sao Tome and Principe"},
	"SA": {id: 3050, code: "SA", name: "Saudi Arabia"},
	"SN": {id: 3060, code: "SN", name: "Senegal"},
	"RS": {id: 3070, code: "RS", name: "Serbia"},
	"SC": {id: 3080, code: "SC", name: "Seychelles"},
	"SL": {id: 3090, code: "SL", name: "Sierra Leone"},
	"SG": {id: 3100, code: "SG", name: "Singapore"},
	"SX": {id: 3110, code: "SX", name: "Sint Maarten (Dutch part)"},
	"SK": {id: 3120, code: "SK", name: "Slovakia"},
	"SI": {id: 3130, code: "SI", name: "Slovenia"},
	"SB": {id: 3140, code: "SB", name: "Solomon Islands"},
	"SO": {id: 3150, code: "SO", name: "Somalia"},
	"ZA": {id: 3160, code: "ZA", name: "South Africa"},
	"GS": {id: 3170, code: "GS", name: "South Georgia and the South Sandwich Islands"},
	"SS": {id: 3180, code: "SS", name: "South Sudan"},
	"ES": {id: 3190, code: "ES", name: "Spain"},
	"LK": {id: 3200, code: "LK", name: "Sri Lanka"},
	"SD": {id: 3210, code: "SD", name: "Sudan"},
	"SR": {id: 3220, code: "SR", name: "Suriname"},
	"SJ": {id: 3230, code: "SJ", name: "Svalbard and Jan Mayen"},
	"SZ": {id: 3240, code: "SZ", name: "Swaziland"},
	"SE": {id: 3250, code: "SE", name: "Sweden"},
	"CH": {id: 3260, code: "CH", name: "Switzerland"},
	"SY": {id: 3270, code: "SY", name: "Syrian Arab Republic"},
	"TW": {id: 3280, code: "TW", name: "Taiwan"},
	"TJ": {id: 3290, code: "TJ", name: "Tajikistan"},
	"TZ": {id: 3300, code: "TZ", name: "Tanzania"},
	"TH": {id: 3310, code: "TH", name: "Thailand"},
	"TL": {id: 3320, code: "TL", name: "Timor-Leste"},
	"TG": {id: 3330, code: "TG", name: "Togo"},
	"TK": {id: 3340, code: "TK", name: "Tokelau"},
	"TO": {id: 3350, code: "TO", name: "Tonga"},
	"TT": {id: 3360, code: "TT", name: "Trinidad and Tobago"},
	"TN": {id: 3370, code: "TN", name: "Tunisia"},
	"TR": {id: 3380, code: "TR", name: "Turkey"},
	"TM": {id: 3390, code: "TM", name: "Turkmenistan"},
	"TC": {id: 3400, code: "TC", name: "Turks and Caicos Islands"},
	"TV": {id: 3410, code: "TV", name: "Tuvalu"},
	"UG": {id: 3420, code: "UG", name: "Uganda"},
	"UA": {id: 3430, code: "UA", name: "Ukraine"},
	"AE": {id: 3440, code: "AE", name: "United Arab Emirates"},
	"GB": {id: 3450, code: "GB", name: "United Kingdom"},
	"US": {id: 3460, code: "US", name: "United States"},
	"UM": {id: 3470, code: "UM", name: "United States Minor Outlying Islands"},
	"UY": {id: 3480, code: "UY", name: "Uruguay"},
	"UZ": {id: 3490, code: "UZ", name: "Uzbekistan"},
	"VU": {id: 3500, code: "VU", name: "Vanuatu"},
	"VE": {id: 3510, code: "VE", name: "Venezuela"},
	"VN": {id: 3520, code: "VN", name: "Viet Nam"},
	"VG": {id: 3530, code: "VG", name: "Virgin Islands, British"},
	"VI": {id: 3540, code: "VI", name: "Virgin Islands, U.S."},
	"WF": {id: 3550, code: "WF", name: "Wallis and Futuna"},
	"EH": {id: 3560, code: "EH", name: "Western Sahara"},
	"YE": {id: 3570, code: "YE", name: "Yemen"},
	"ZM": {id: 3580, code: "ZM", name: "Zambia"},
	"ZW": {id: 3590, code: "ZW", name: "Zimbabwe"},
}

var countryLanguage = map[string]bool{
	"af-za": true,
	"ar-ae": true,
	"ar-bh": true,
	"ar-dz": true,
	"ar-eg": true,
	"ar-iq": true,
	"ar-jo": true,
	"ar-kw": true,
	"ar-lb": true,
	"ar-ly": true,
	"ar-ma": true,
	"ar-om": true,
	"ar-qa": true,
	"ar-sa": true,
	"ar-sy": true,
	"ar-tn": true,
	"ar-ye": true,
	"be-by": true,
	"bg-bg": true,
	"ca-es": true,
	"cs-cz": true,
	"cy-gb": true,
	"da-dk": true,
	"de-at": true,
	"de-de": true,
	"de-ch": true,
	"de-li": true,
	"de-lu": true,
	"dv-mv": true,
	"el-gr": true,
	"en-au": true,
	"en-bz": true,
	"en-ca": true,
	"en-gb": true,
	"en-ie": true,
	"en-jm": true,
	"en-nz": true,
	"en-ph": true,
	"en-tt": true,
	"en-us": true,
	"en-za": true,
	"en-zw": true,
	"es-ar": true,
	"es-bo": true,
	"es-cl": true,
	"es-co": true,
	"es-cr": true,
	"es-do": true,
	"es-ec": true,
	"es-es": true,
	"es-gt": true,
	"es-hn": true,
	"es-mx": true,
	"es-ni": true,
	"es-pa": true,
	"es-pe": true,
	"es-pr": true,
	"es-py": true,
	"es-sv": true,
	"es-uy": true,
	"es-ve": true,
	"et-ee": true,
	"eu-es": true,
	"fa-ir": true,
	"fi-fi": true,
	"fo-fo": true,
	"fr-be": true,
	"fr-ca": true,
	"fr-fr": true,
	"fr-ch": true,
	"fr-lu": true,
	"fr-mc": true,
	"gl-es": true,
	"gu-in": true,
	"he-il": true,
	"hi-in": true,
	"hr-ba": true,
	"hr-hr": true,
	"hu-hu": true,
	"hy-am": true,
	"id-id": true,
	"is-is": true,
	"it-ch": true,
	"it-it": true,
	"ja-jp": true,
	"ka-ge": true,
	"kk-kz": true,
	"kn-in": true,
	"ko-kr": true,
	"ky-kg": true,
	"lt-lt": true,
	"lv-lv": true,
	"mi-nz": true,
	"mk-mk": true,
	"mn-mn": true,
	"mr-in": true,
	"ms-bn": true,
	"ms-my": true,
	"mt-mt": true,
	"nb-no": true,
	"nl-be": true,
	"nl-nl": true,
	"nn-no": true,
	"ns-za": true,
	"pa-in": true,
	"pl-pl": true,
	"pt-br": true,
	"pt-pt": true,
	"qu-bo": true,
	"qu-ec": true,
	"qu-pe": true,
	"ro-ro": true,
	"ru-ru": true,
	"sa-in": true,
	"se-fi": true,
	"se-no": true,
	"se-se": true,
	"sk-sk": true,
	"sl-si": true,
	"sq-al": true,
	"sv-fi": true,
	"sv-se": true,
	"sw-ke": true,
	"ta-in": true,
	"te-in": true,
	"th-th": true,
	"tn-za": true,
	"tr-tr": true,
	"tt-ru": true,
	"uk-ua": true,
	"ur-pk": true,
	"vi-vn": true,
	"xh-za": true,
	"zh-cn": true,
	"zh-hk": true,
	"zh-mo": true,
	"zh-sg": true,
	"zh-tw": true,
	"zu-za": true,
}

func (language *Language) Init(id int32, code string) {
	language.currentCountry = language.getCountry(code)
	lang := language.getLanguage(code)
	language.id = id
	if lang != nil {
		if id <= 0 {
			language.id = lang.id
		}
		language.code = lang.code
		language.name = lang.name
		language.nativeName = lang.nativeName
	}
}

//******************************
// getters and setters
//******************************
func (language *Language) GetCode() string {
	return language.code
}

func (language *Language) GetName() string {
	return language.name
}

func (language *Language) GetNativeName() string {
	return language.nativeName
}

func (language *Language) GetEntityType() entitytype.EntityType {
	return entitytype.Language
}

//******************************
// public methods
//******************************
func (language *Language) GetDescription() string {
	var result = ""
	result += language.name
	if language.currentCountry != nil {
		//eg. English (United States)
		result += " (" + language.currentCountry.name + ")"
	}
	return result
}

func (language *Language) DisplayValue(code string) string {
	var result = languages[string(code[0:2])].id
	return strconv.Itoa(int(result))
}

func (language *Language) IsCodeValid(code string) (bool, error) {
	var codeFormat = strings.ReplaceAll(code, " ", "")
	if len(code) == 0 {
		return false, errors.New(languageErrorEmpty)
	}

	if strings.Contains(codeFormat, "-") == true {
		var country = language.getCountry(code)
		if country == nil {
			return false, errors.New(fmt.Sprintf(languageErrorInvalidCountry, code))
		}
		// lang country combination exists
		// ...
	} else {
		var lang = language.getLanguage(code)
		// just language
		if lang == nil {
			return false, errors.New(fmt.Sprintf(languageErrorInvalidLang, code))
		}
	}
	return true, nil
}

//******************************
// private methods
//******************************
func (language *Language) exists(schema *Schema) bool {
	query := new(metaQuery)
	query.setSchema(metaSchemaName)
	query.setTable(metaTableName)

	query.addFilter(metaId, operatorEqual, language.id)
	query.addFilter(metaSchemaId, operatorEqual, schema.id)
	query.addFilter(metaObjectType, operatorEqual, int8(entitytype.Language))
	query.addFilter(metaReferenceId, operatorEqual, schema.id)

	result, _ := query.exists()
	return result
}

func (language *Language) create(schema *Schema) error {
	query := new(metaQuery)
	query.setSchema(metaSchemaName)
	query.setTable(metaTableName)
	return query.insertMeta(language.toMeta(), schema.id)
}

func (language *Language) toMeta() *Meta {
	var metaTable = new(Meta)

	// key
	metaTable.id = language.id
	metaTable.name = languageName
	metaTable.description = language.GetDescription()
	metaTable.objectType = int8(entitytype.Language)
	metaTable.dataType = 0
	metaTable.setEntityBaseline(true)
	metaTable.value = language.code
	if language.currentCountry != nil {
		metaTable.value += "-" + language.currentCountry.code
	}
	metaTable.enabled = true

	return metaTable
}

func (language *Language) getCountry(code string) *country {
	var formattedCode = strings.ReplaceAll(code, " ", "")
	var index = strings.Index(formattedCode, "-")

	if index > 0 && index+1 < len(formattedCode) {
		var countryCode = formattedCode[index+1:]
		if val, ok := countries[strings.ToUpper(countryCode)]; ok {
			return &val
		}
	}
	return nil
}

func (language *Language) getLanguage(code string) *Language {
	var formattedCode = strings.ReplaceAll(code, " ", "")
	var index = strings.Index(formattedCode, "-")

	if index > 0 && index+1 < len(formattedCode) {
		var languageCode = formattedCode[:index]
		if val, ok := languages[strings.ToLower(languageCode)]; ok {
			return &val
		}
	}
	return nil
}

func (language *Language) getDefaultLanguage(code string) *Language {
	return nil
}
