package schema

import (
	"errors"
	"fmt"
	"ring/schema/entitytype"
	"strings"
)

// Language is an ISO 639-1 language with code, name and native name.
type Language struct {
	id             int16
	code           string
	name           string
	nativeName     string
	currentCountry *country
	valid          bool
}

type country struct {
	id   int16
	code string
	name string
}

const (
	languageErrorEmpty           string = "empty code"
	languageErrorInvalidCountry  string = "invalid country code '%s'"
	languageErrorInvalidLang     string = "invalid language code '%s'"
	languageErrorInvalidCounLang string = "the combination of language code and country '%s' does not exist"
	languageCountrySeparator     string = "-"
	languageToStringFormat       string = "id: %d; code: %s; name: %s; description: %s"
)

// [10000, 20000] - language
var languages = map[string]*Language{
	"aa": {id: 11111, code: "aa", name: "Afar", nativeName: "Afaraf", valid: true},
	"ab": {id: 11121, code: "ab", name: "Abkhaz", nativeName: "аҧсуа бызшәа", valid: true},
	"ae": {id: 11131, code: "ae", name: "Avestan", nativeName: "avesta", valid: true},
	"af": {id: 11141, code: "af", name: "Afrikaans", nativeName: "Afrikaans", valid: true},
	"ak": {id: 11151, code: "ak", name: "Akan", nativeName: "Akan", valid: true},
	"am": {id: 11161, code: "am", name: "Amharic", nativeName: "አማርኛ", valid: true},
	"an": {id: 11171, code: "an", name: "Aragonese", nativeName: "aragonés", valid: true},
	"ar": {id: 11181, code: "ar", name: "Arabic", nativeName: "اللغة العربية", valid: true},
	"as": {id: 11191, code: "as", name: "Assamese", nativeName: "অসমীয়া", valid: true},
	"av": {id: 11201, code: "av", name: "Avaric", nativeName: "авар мацӀ", valid: true},
	"ay": {id: 11211, code: "ay", name: "Aymara", nativeName: "aymar aru", valid: true},
	"az": {id: 11221, code: "az", name: "Azerbaijani", nativeName: "azərbaycan dili", valid: true},
	"ba": {id: 11232, code: "ba", name: "Bashkir", nativeName: "башҡорт теле", valid: true},
	"be": {id: 11241, code: "be", name: "Belarusian", nativeName: "беларуская мова", valid: true},
	"bg": {id: 11251, code: "bg", name: "Bulgarian", nativeName: "български език", valid: true},
	"bh": {id: 11261, code: "bh", name: "Bihari", nativeName: "भोजपुरी", valid: true},
	"bi": {id: 11271, code: "bi", name: "Bislama", nativeName: "Bislama", valid: true},
	"bm": {id: 11281, code: "bm", name: "Bambara", nativeName: "bamanankan", valid: true},
	"bn": {id: 11291, code: "bn", name: "Bengali", nativeName: "বাংলা", valid: true},
	"bo": {id: 11301, code: "bo", name: "Tibetan Standard", nativeName: "བོད་ཡིག", valid: true},
	"br": {id: 11311, code: "br", name: "Breton", nativeName: "brezhoneg", valid: true},
	"bs": {id: 11321, code: "bs", name: "Bosnian", nativeName: "bosanski jezik", valid: true},
	"ca": {id: 11331, code: "ca", name: "Catalan", nativeName: "català", valid: true},
	"ce": {id: 11341, code: "ce", name: "Chechen", nativeName: "нохчийн мотт", valid: true},
	"ch": {id: 11351, code: "ch", name: "Chamorro", nativeName: "Chamoru", valid: true},
	"co": {id: 11361, code: "co", name: "Corsican", nativeName: "corsu", valid: true},
	"cr": {id: 11371, code: "cr", name: "Cree", nativeName: "ᓀᐦᐃᔭᐍᐏᐣ", valid: true},
	"cs": {id: 11381, code: "cs", name: "Czech", nativeName: "čeština", valid: true},
	"cu": {id: 11391, code: "cu", name: "Old Church Slavonic", nativeName: "ѩзыкъ словѣньскъ", valid: true},
	"cv": {id: 11401, code: "cv", name: "Chuvash", nativeName: "чӑваш чӗлхи", valid: true},
	"cy": {id: 11411, code: "cy", name: "Welsh", nativeName: "Cymraeg", valid: true},
	"da": {id: 11421, code: "da", name: "Danish", nativeName: "dansk", valid: true},
	"de": {id: 11431, code: "de", name: "German", nativeName: "Deutsch", valid: true},
	"dv": {id: 11441, code: "dv", name: "Divehi", nativeName: "Dhivehi", valid: true},
	"dz": {id: 11451, code: "dz", name: "Dzongkha", nativeName: "རྫོང་ཁ", valid: true},
	"ee": {id: 11461, code: "ee", name: "Ewe", nativeName: "Eʋegbe", valid: true},
	"el": {id: 11471, code: "el", name: "Greek", nativeName: "Ελληνικά", valid: true},
	"en": {id: 11481, code: "en", name: "English", nativeName: "English", valid: true},
	"eo": {id: 11491, code: "eo", name: "Esperanto", nativeName: "Esperanto", valid: true},
	"es": {id: 11501, code: "es", name: "Spanish", nativeName: "Español", valid: true},
	"et": {id: 11511, code: "et", name: "Estonian", nativeName: "eesti", valid: true},
	"eu": {id: 11521, code: "eu", name: "Basque", nativeName: "euskara", valid: true},
	"fa": {id: 11531, code: "fa", name: "Persian", nativeName: "فارسی", valid: true},
	"ff": {id: 11541, code: "ff", name: "Fula", nativeName: "Fulfulde", valid: true},
	"fi": {id: 11551, code: "fi", name: "Finnish", nativeName: "suomi", valid: true},
	"fj": {id: 11561, code: "fj", name: "Fijian", nativeName: "Vakaviti", valid: true},
	"fo": {id: 11571, code: "fo", name: "Faroese", nativeName: "føroyskt", valid: true},
	"fr": {id: 11581, code: "fr", name: "French", nativeName: "Français", valid: true},
	"fy": {id: 11591, code: "fy", name: "Western Frisian", nativeName: "Frysk", valid: true},
	"ga": {id: 11601, code: "ga", name: "Irish", nativeName: "Gaeilge", valid: true},
	"gd": {id: 11611, code: "gd", name: "Scottish Gaelic", nativeName: "Gàidhlig", valid: true},
	"gl": {id: 11621, code: "gl", name: "Galician", nativeName: "galego", valid: true},
	"gn": {id: 11631, code: "gn", name: "Guaraní", nativeName: "Avañeẽ", valid: true},
	"gu": {id: 11641, code: "gu", name: "Gujarati", nativeName: "ગુજરાતી", valid: true},
	"gv": {id: 11651, code: "gv", name: "Manx", nativeName: "Gaelg", valid: true},
	"ha": {id: 11661, code: "ha", name: "Hausa", nativeName: "هَوُسَ", valid: true},
	"he": {id: 11671, code: "he", name: "Hebrew", nativeName: "עברית", valid: true},
	"hi": {id: 11681, code: "hi", name: "Hindi", nativeName: "हिन्दी", valid: true},
	"ho": {id: 11691, code: "ho", name: "Hiri Motu", nativeName: "Hiri Motu", valid: true},
	"hr": {id: 11701, code: "hr", name: "Croatian", nativeName: "hrvatski jezik", valid: true},
	"ht": {id: 11711, code: "ht", name: "Haitian", nativeName: "Kreyòl ayisyen", valid: true},
	"hu": {id: 11721, code: "hu", name: "Hungarian", nativeName: "magyar", valid: true},
	"hy": {id: 11731, code: "hy", name: "Armenian", nativeName: "Հայերեն", valid: true},
	"hz": {id: 11741, code: "hz", name: "Herero", nativeName: "Otjiherero", valid: true},
	"ia": {id: 11751, code: "ia", name: "Interlingua", nativeName: "Interlingua", valid: true},
	"id": {id: 11761, code: "id", name: "Indonesian", nativeName: "Indonesian", valid: true},
	"ie": {id: 11771, code: "ie", name: "Interlingue", nativeName: "Interlingue", valid: true},
	"ig": {id: 11781, code: "ig", name: "Igbo", nativeName: "Asụsụ Igbo", valid: true},
	"ii": {id: 11791, code: "ii", name: "Nuosu", nativeName: "ꆈꌠ꒿ Nuosuhxop", valid: true},
	"ik": {id: 11801, code: "ik", name: "Inupiaq", nativeName: "Iñupiaq", valid: true},
	"io": {id: 11811, code: "io", name: "Ido", nativeName: "Ido", valid: true},
	"is": {id: 11821, code: "is", name: "Icelandic", nativeName: "Íslenska", valid: true},
	"it": {id: 11831, code: "it", name: "Italian", nativeName: "Italiano", valid: true},
	"iu": {id: 11841, code: "iu", name: "Inuktitut", nativeName: "ᐃᓄᒃᑎᑐᑦ", valid: true},
	"ja": {id: 11851, code: "ja", name: "Japanese", nativeName: "日本語", valid: true},
	"jv": {id: 11861, code: "jv", name: "Javanese", nativeName: "basa Jawa", valid: true},
	"ka": {id: 11871, code: "ka", name: "Georgian", nativeName: "ქართული", valid: true},
	"kg": {id: 11881, code: "kg", name: "Kongo", nativeName: "Kikongo", valid: true},
	"ki": {id: 11891, code: "ki", name: "Kikuyu", nativeName: "Gĩkũyũ", valid: true},
	"kj": {id: 11901, code: "kj", name: "Kwanyama", nativeName: "Kuanyama", valid: true},
	"kk": {id: 11911, code: "kk", name: "Kazakh", nativeName: "қазақ тілі", valid: true},
	"kl": {id: 11921, code: "kl", name: "Kalaallisut", nativeName: "kalaallisut", valid: true},
	"km": {id: 11931, code: "km", name: "Khmer", nativeName: "ខេមរភាសា", valid: true},
	"kn": {id: 11941, code: "kn", name: "Kannada", nativeName: "ಕನ್ನಡ", valid: true},
	"ko": {id: 11951, code: "ko", name: "Korean", nativeName: "한국어", valid: true},
	"kr": {id: 11961, code: "kr", name: "Kanuri", nativeName: "Kanuri", valid: true},
	"ks": {id: 11971, code: "ks", name: "Kashmiri", nativeName: "कश्मीरी", valid: true},
	"ku": {id: 11981, code: "ku", name: "Kurdish", nativeName: "Kurdî", valid: true},
	"kv": {id: 11991, code: "kv", name: "Komi", nativeName: "коми кыв", valid: true},
	"kw": {id: 12001, code: "kw", name: "Cornish", nativeName: "Kernewek", valid: true},
	"ky": {id: 12011, code: "ky", name: "Kyrgyz", nativeName: "Кыргызча", valid: true},
	"la": {id: 12021, code: "la", name: "Latin", nativeName: "latine", valid: true},
	"lb": {id: 12031, code: "lb", name: "Luxembourgish", nativeName: "Lëtzebuergesch", valid: true},
	"lg": {id: 12041, code: "lg", name: "Ganda", nativeName: "Luganda", valid: true},
	"li": {id: 12051, code: "li", name: "Limburgish", nativeName: "Limburgs", valid: true},
	"ln": {id: 12061, code: "ln", name: "Lingala", nativeName: "Lingála", valid: true},
	"lo": {id: 12071, code: "lo", name: "Lao", nativeName: "ພາສາ", valid: true},
	"lt": {id: 12081, code: "lt", name: "Lithuanian", nativeName: "lietuvių kalba", valid: true},
	"lu": {id: 12091, code: "lu", name: "Luba-Katanga", nativeName: "Tshiluba", valid: true},
	"lv": {id: 12101, code: "lv", name: "Latvian", nativeName: "latviešu valoda", valid: true},
	"mg": {id: 12111, code: "mg", name: "Malagasy", nativeName: "fiteny malagasy", valid: true},
	"mh": {id: 12121, code: "mh", name: "Marshallese", nativeName: "Kajin M̧ajeļ", valid: true},
	"mi": {id: 12131, code: "mi", name: "Māori", nativeName: "te reo Māori", valid: true},
	"mk": {id: 12141, code: "mk", name: "Macedonian", nativeName: "македонски јазик", valid: true},
	"ml": {id: 12151, code: "ml", name: "Malayalam", nativeName: "മലയാളം", valid: true},
	"mn": {id: 12161, code: "mn", name: "Mongolian", nativeName: "Монгол хэл", valid: true},
	"mr": {id: 12171, code: "mr", name: "Marathi", nativeName: "मराठी", valid: true},
	"ms": {id: 12181, code: "ms", name: "Malay", nativeName: "هاس ملايو‎", valid: true},
	"mt": {id: 12191, code: "mt", name: "Maltese", nativeName: "Malti", valid: true},
	"my": {id: 12201, code: "my", name: "Burmese", nativeName: "ဗမာစာ", valid: true},
	"na": {id: 12211, code: "na", name: "Nauru", nativeName: "Ekakairũ Naoero", valid: true},
	"nb": {id: 12221, code: "nb", name: "Norwegian Bokmål", nativeName: "Norsk bokmål", valid: true},
	"nd": {id: 12231, code: "nd", name: "Northern Ndebele", nativeName: "isiNdebele", valid: true},
	"ne": {id: 12241, code: "ne", name: "Nepali", nativeName: "नेपाली", valid: true},
	"ng": {id: 12251, code: "ng", name: "Ndonga", nativeName: "Owambo", valid: true},
	"nl": {id: 12261, code: "nl", name: "Dutch", nativeName: "Nederlands", valid: true},
	"nn": {id: 12271, code: "nn", name: "Norwegian Nynorsk", nativeName: "Norsk nynorsk", valid: true},
	"no": {id: 12281, code: "no", name: "Norwegian", nativeName: "Norsk", valid: true},
	"nr": {id: 12291, code: "nr", name: "Southern Ndebele", nativeName: "isiNdebele", valid: true},
	"nv": {id: 12301, code: "nv", name: "Navajo", nativeName: "Diné bizaad", valid: true},
	"ny": {id: 12311, code: "ny", name: "Chichewa", nativeName: "chiCheŵa", valid: true},
	"oc": {id: 12321, code: "oc", name: "Occitan", nativeName: "occitan", valid: true},
	"oj": {id: 12331, code: "oj", name: "Ojibwe", nativeName: "ᐊᓂᔑᓈᐯᒧᐎᓐ", valid: true},
	"om": {id: 12341, code: "om", name: "Oromo", nativeName: "Afaan Oromoo", valid: true},
	"or": {id: 12351, code: "or", name: "Oriya", nativeName: "ଓଡ଼ିଆ", valid: true},
	"os": {id: 12361, code: "os", name: "Ossetian", nativeName: "ирон æвзаг", valid: true},
	"pa": {id: 12371, code: "pa", name: "Panjabi", nativeName: "ਪੰਜਾਬੀ", valid: true},
	"pi": {id: 12381, code: "pi", name: "Pāli", nativeName: "पाऴि", valid: true},
	"pl": {id: 12391, code: "pl", name: "Polish", nativeName: "język polski", valid: true},
	"ps": {id: 12401, code: "ps", name: "Pashto", nativeName: "پښتو", valid: true},
	"pt": {id: 12411, code: "pt", name: "Portuguese", nativeName: "Português", valid: true},
	"qu": {id: 12421, code: "qu", name: "Quechua", nativeName: "Runa Simi", valid: true},
	"rm": {id: 12431, code: "rm", name: "Romansh", nativeName: "rumantsch grischun", valid: true},
	"rn": {id: 12441, code: "rn", name: "Kirundi", nativeName: "Ikirundi", valid: true},
	"ro": {id: 12451, code: "ro", name: "Romanian", nativeName: "Română", valid: true},
	"ru": {id: 12461, code: "ru", name: "Russian", nativeName: "Русский", valid: true},
	"rw": {id: 12471, code: "rw", name: "Kinyarwanda", nativeName: "Ikinyarwanda", valid: true},
	"sa": {id: 12481, code: "sa", name: "Sanskrit", nativeName: "संस्कृतम्", valid: true},
	"sc": {id: 12491, code: "sc", name: "Sardinian", nativeName: "sardu", valid: true},
	"sd": {id: 12501, code: "sd", name: "Sindhi", nativeName: "सिन्धी", valid: true},
	"se": {id: 12511, code: "se", name: "Northern Sami", nativeName: "Davvisámegiella", valid: true},
	"sg": {id: 12521, code: "sg", name: "Sango", nativeName: "yângâ tî sängö", valid: true},
	"si": {id: 12531, code: "si", name: "Sinhala", nativeName: "සිංහල", valid: true},
	"sk": {id: 12541, code: "sk", name: "Slovak", nativeName: "slovenčina", valid: true},
	"sl": {id: 12551, code: "sl", name: "Slovene", nativeName: "slovenski jezik", valid: true},
	"sm": {id: 12561, code: "sm", name: "Samoan", nativeName: "gagana faa Samoa", valid: true},
	"sn": {id: 12571, code: "sn", name: "Shona", nativeName: "chiShona", valid: true},
	"so": {id: 12581, code: "so", name: "Somali", nativeName: "Soomaaliga", valid: true},
	"sq": {id: 12591, code: "sq", name: "Albanian", nativeName: "Shqip", valid: true},
	"sr": {id: 12601, code: "sr", name: "Serbian", nativeName: "српски језик", valid: true},
	"ss": {id: 12611, code: "ss", name: "Swati", nativeName: "SiSwati", valid: true},
	"st": {id: 12621, code: "st", name: "Southern Sotho", nativeName: "Sesotho", valid: true},
	"su": {id: 12631, code: "su", name: "Sundanese", nativeName: "Basa Sunda", valid: true},
	"sv": {id: 12641, code: "sv", name: "Swedish", nativeName: "svenska", valid: true},
	"sw": {id: 12651, code: "sw", name: "Swahili", nativeName: "Kiswahili", valid: true},
	"ta": {id: 12661, code: "ta", name: "Tamil", nativeName: "தமிழ்", valid: true},
	"te": {id: 12671, code: "te", name: "Telugu", nativeName: "తెలుగు", valid: true},
	"tg": {id: 12681, code: "tg", name: "Tajik", nativeName: "тоҷикӣ", valid: true},
	"th": {id: 12691, code: "th", name: "Thai", nativeName: "ไทย", valid: true},
	"ti": {id: 12701, code: "ti", name: "Tigrinya", nativeName: "ትግርኛ", valid: true},
	"tk": {id: 12711, code: "tk", name: "Turkmen", nativeName: "Türkmen", valid: true},
	"tl": {id: 12721, code: "tl", name: "Tagalog", nativeName: "Wikang Tagalog", valid: true},
	"tn": {id: 12731, code: "tn", name: "Tswana", nativeName: "Setswana", valid: true},
	"to": {id: 12741, code: "to", name: "Tonga", nativeName: "faka Tonga", valid: true},
	"tr": {id: 12751, code: "tr", name: "Turkish", nativeName: "Türkçe", valid: true},
	"ts": {id: 12761, code: "ts", name: "Tsonga", nativeName: "Xitsonga", valid: true},
	"tt": {id: 12771, code: "tt", name: "Tatar", nativeName: "татар теле", valid: true},
	"tw": {id: 12781, code: "tw", name: "Twi", nativeName: "Twi", valid: true},
	"ty": {id: 12791, code: "ty", name: "Tahitian", nativeName: "Reo Tahiti", valid: true},
	"ug": {id: 12801, code: "ug", name: "Uyghur", nativeName: "ئۇيغۇرچە‎", valid: true},
	"uk": {id: 12811, code: "uk", name: "Ukrainian", nativeName: "Українська", valid: true},
	"ur": {id: 12821, code: "ur", name: "Urdu", nativeName: "اردو", valid: true},
	"uz": {id: 12831, code: "uz", name: "Uzbek", nativeName: "Ўзбек", valid: true},
	"ve": {id: 12841, code: "ve", name: "Venda", nativeName: "Tshivenḓa", valid: true},
	"vi": {id: 12851, code: "vi", name: "Vietnamese", nativeName: "Tiếng Việt", valid: true},
	"vo": {id: 12861, code: "vo", name: "Volapük", nativeName: "Volapük", valid: true},
	"wa": {id: 12871, code: "wa", name: "Walloon", nativeName: "walon", valid: true},
	"wo": {id: 12881, code: "wo", name: "Wolof", nativeName: "Wollof", valid: true},
	"xh": {id: 12891, code: "xh", name: "Xhosa", nativeName: "isiXhosa", valid: true},
	"yi": {id: 12901, code: "yi", name: "Yiddish", nativeName: "ייִדיש", valid: true},
	"yo": {id: 12911, code: "yo", name: "Yoruba", nativeName: "Yorùbá", valid: true},
	"za": {id: 12921, code: "za", name: "Zhuang", nativeName: "Saɯ cueŋƅ", valid: true},
	"zh": {id: 12931, code: "zh", name: "Chinese", nativeName: "中文", valid: true},
	"zu": {id: 12941, code: "zu", name: "Zulu", nativeName: "isiZulu", valid: true},
}

// [10000, 20000] - country
var countries = map[string]country{
	"AF": {id: 11102, code: "AF", name: "Afghanistan"},
	"AX": {id: 11113, code: "AX", name: "Åland Islands"},
	"AL": {id: 11123, code: "AL", name: "Albania"},
	"DZ": {id: 11133, code: "DZ", name: "Algeria"},
	"AS": {id: 11143, code: "AS", name: "American Samoa"},
	"AD": {id: 11153, code: "AD", name: "Andorra"},
	"AO": {id: 11163, code: "AO", name: "Angola"},
	"AI": {id: 11173, code: "AI", name: "Anguilla"},
	"AQ": {id: 11183, code: "AQ", name: "Antarctica"},
	"AG": {id: 11193, code: "AG", name: "Antigua and Barbuda"},
	"AR": {id: 11202, code: "AR", name: "Argentina"},
	"AM": {id: 11213, code: "AM", name: "Armenia"},
	"AW": {id: 11223, code: "AW", name: "Aruba"},
	"AU": {id: 11233, code: "AU", name: "Australia"},
	"AT": {id: 11243, code: "AT", name: "Austria"},
	"AZ": {id: 11253, code: "AZ", name: "Azerbaijan"},
	"BH": {id: 11263, code: "BH", name: "Bahrain"},
	"BS": {id: 11273, code: "BS", name: "Bahamas"},
	"BD": {id: 11283, code: "BD", name: "Bangladesh"},
	"BB": {id: 11293, code: "BB", name: "Barbados"},
	"BY": {id: 11303, code: "BY", name: "Belarus"},
	"BE": {id: 11313, code: "BE", name: "Belgium"},
	"BZ": {id: 11323, code: "BZ", name: "Belize"},
	"BJ": {id: 11333, code: "BJ", name: "Benin"},
	"BM": {id: 11343, code: "BM", name: "Bermuda"},
	"BT": {id: 11353, code: "BT", name: "Bhutan"},
	"BO": {id: 11363, code: "BO", name: "Bolivia, Plurinational State of"},
	"BQ": {id: 11373, code: "BQ", name: "Bonaire, Sint Eustatius and Saba"},
	"BA": {id: 11383, code: "BA", name: "Bosnia and Herzegovina"},
	"BW": {id: 11393, code: "BW", name: "Botswana"},
	"BV": {id: 11403, code: "BV", name: "Bouvet Island"},
	"BR": {id: 11413, code: "BR", name: "Brazil"},
	"IO": {id: 11423, code: "IO", name: "British Indian Ocean Territory"},
	"BN": {id: 11433, code: "BN", name: "Brunei Darussalam"},
	"BG": {id: 11443, code: "BG", name: "Bulgaria"},
	"BF": {id: 11453, code: "BF", name: "Burkina Faso"},
	"BI": {id: 11463, code: "BI", name: "Burundi"},
	"KH": {id: 11473, code: "KH", name: "Cambodia"},
	"CM": {id: 11483, code: "CM", name: "Cameroon"},
	"CA": {id: 11495, code: "CA", name: "Canada"},
	"CV": {id: 11503, code: "CV", name: "Cape Verde"},
	"KY": {id: 11513, code: "KY", name: "Cayman Islands"},
	"CF": {id: 11523, code: "CF", name: "Central African Republic"},
	"TD": {id: 11533, code: "TD", name: "Chad"},
	"CL": {id: 11543, code: "CL", name: "Chile"},
	"CN": {id: 11553, code: "CN", name: "China"},
	"CX": {id: 11563, code: "CX", name: "Christmas Island"},
	"CC": {id: 11573, code: "CC", name: "Cocos (Keeling) Islands"},
	"CO": {id: 11584, code: "CO", name: "Colombia"},
	"KM": {id: 11593, code: "KM", name: "Comoros"},
	"CG": {id: 11603, code: "CG", name: "Congo"},
	"CD": {id: 11613, code: "CD", name: "Congo"},
	"CK": {id: 11623, code: "CK", name: "Cook Islands"},
	"CR": {id: 11633, code: "CR", name: "Costa Rica"},
	"CI": {id: 11643, code: "CI", name: "Côte d'Ivoire"},
	"HR": {id: 11654, code: "HR", name: "Croatia"},
	"CU": {id: 11663, code: "CU", name: "Cuba"},
	"CW": {id: 11673, code: "CW", name: "Curaçao"},
	"CY": {id: 11683, code: "CY", name: "Cyprus"},
	"CZ": {id: 11693, code: "CZ", name: "Czech Republic"},
	"DK": {id: 11703, code: "DK", name: "Denmark"},
	"DJ": {id: 11713, code: "DJ", name: "Djibouti"},
	"DM": {id: 11723, code: "DM", name: "Dominica"},
	"DO": {id: 11733, code: "DO", name: "Dominican Republic"},
	"EC": {id: 11744, code: "EC", name: "Ecuador"},
	"EG": {id: 11753, code: "EG", name: "Egypt"},
	"SV": {id: 11763, code: "SV", name: "El Salvador"},
	"GQ": {id: 11773, code: "GQ", name: "Equatorial Guinea"},
	"ER": {id: 11783, code: "ER", name: "Eritrea"},
	"EE": {id: 11793, code: "EE", name: "Estonia"},
	"ET": {id: 11803, code: "ET", name: "Ethiopia"},
	"FK": {id: 11813, code: "FK", name: "Falkland Islands (Malvinas)"},
	"FO": {id: 11823, code: "FO", name: "Faroe Islands"},
	"FJ": {id: 11833, code: "FJ", name: "Fiji"},
	"FI": {id: 11844, code: "FI", name: "Finland"},
	"FR": {id: 11853, code: "FR", name: "France"},
	"GF": {id: 11863, code: "GF", name: "French Guiana"},
	"PF": {id: 11873, code: "PF", name: "French Polynesia"},
	"TF": {id: 11883, code: "TF", name: "French Southern Territories"},
	"GA": {id: 11893, code: "GA", name: "Gabon"},
	"GM": {id: 11903, code: "GM", name: "Gambia"},
	"GE": {id: 11914, code: "GE", name: "Georgia"},
	"DE": {id: 11923, code: "DE", name: "Germany"},
	"GH": {id: 11933, code: "GH", name: "Ghana"},
	"GI": {id: 11943, code: "GI", name: "Gibraltar"},
	"GR": {id: 11953, code: "GR", name: "Greece"},
	"GL": {id: 11963, code: "GL", name: "Greenland"},
	"GD": {id: 11973, code: "GD", name: "Grenada"},
	"GP": {id: 11983, code: "GP", name: "Guadeloupe"},
	"GU": {id: 11993, code: "GU", name: "Guam"},
	"GT": {id: 12003, code: "GT", name: "Guatemala"},
	"GG": {id: 12013, code: "GG", name: "Guernsey"},
	"GN": {id: 12023, code: "GN", name: "Guinea"},
	"GW": {id: 12033, code: "GW", name: "Guinea-Bissau"},
	"GY": {id: 12043, code: "GY", name: "Guyana"},
	"HT": {id: 12053, code: "HT", name: "Haiti"},
	"HM": {id: 12063, code: "HM", name: "Heard Island and McDonald Islands"},
	"VA": {id: 12073, code: "VA", name: "Holy See (Vatican City State)"},
	"HN": {id: 12083, code: "HN", name: "Honduras"},
	"HK": {id: 12093, code: "HK", name: "Hong Kong"},
	"HU": {id: 12104, code: "HU", name: "Hungary"},
	"IS": {id: 12113, code: "IS", name: "Iceland"},
	"IN": {id: 12120, code: "IN", name: "India"},
	"ID": {id: 12134, code: "ID", name: "Indonesia"},
	"IR": {id: 12143, code: "IR", name: "Iran"},
	"IQ": {id: 12153, code: "IQ", name: "Iraq"},
	"IE": {id: 12163, code: "IE", name: "Ireland"},
	"IM": {id: 12173, code: "IM", name: "Isle of Man"},
	"IL": {id: 12183, code: "IL", name: "Israel"},
	"IT": {id: 12193, code: "IT", name: "Italy"},
	"JM": {id: 12203, code: "JM", name: "Jamaica"},
	"JP": {id: 12213, code: "JP", name: "Japan"},
	"JE": {id: 12223, code: "JE", name: "Jersey"},
	"JO": {id: 12233, code: "JO", name: "Jordan"},
	"KZ": {id: 12243, code: "KZ", name: "Kazakhstan"},
	"KE": {id: 12253, code: "KE", name: "Kenya"},
	"KI": {id: 12263, code: "KI", name: "Kiribati"},
	"KP": {id: 12273, code: "KP", name: "Korea"},
	"KR": {id: 12283, code: "KR", name: "Korea, Republic of"},
	"KW": {id: 12293, code: "KW", name: "Kuwait"},
	"KG": {id: 12303, code: "KG", name: "Kyrgyzstan"},
	"LA": {id: 12310, code: "LA", name: "Lao People's Democratic Republic"},
	"LV": {id: 12323, code: "LV", name: "Latvia"},
	"LB": {id: 12333, code: "LB", name: "Lebanon"},
	"LS": {id: 12343, code: "LS", name: "Lesotho"},
	"LR": {id: 12353, code: "LR", name: "Liberia"},
	"LY": {id: 12363, code: "LY", name: "Libya"},
	"LI": {id: 12373, code: "LI", name: "Liechtenstein"},
	"LT": {id: 12383, code: "LT", name: "Lithuania"},
	"LU": {id: 12395, code: "LU", name: "Luxembourg"},
	"MO": {id: 12403, code: "MO", name: "Macao"},
	"MK": {id: 12413, code: "MK", name: "North Macedonia"},
	"MG": {id: 12423, code: "MG", name: "Madagascar"},
	"MW": {id: 12433, code: "MW", name: "Malawi"},
	"MY": {id: 12443, code: "MY", name: "Malaysia"},
	"MV": {id: 12453, code: "MV", name: "Maldives"},
	"ML": {id: 12463, code: "ML", name: "Mali"},
	"MT": {id: 12473, code: "MT", name: "Malta"},
	"MH": {id: 12483, code: "MH", name: "Marshall Islands"},
	"MQ": {id: 12493, code: "MQ", name: "Martinique"},
	"MR": {id: 12503, code: "MR", name: "Mauritania"},
	"MU": {id: 12513, code: "MU", name: "Mauritius"},
	"YT": {id: 12523, code: "YT", name: "Mayotte"},
	"MX": {id: 12533, code: "MX", name: "Mexico"},
	"FM": {id: 12543, code: "FM", name: "Micronesia"},
	"MD": {id: 12553, code: "MD", name: "Moldova"},
	"MC": {id: 12563, code: "MC", name: "Monaco"},
	"MN": {id: 12573, code: "MN", name: "Mongolia"},
	"ME": {id: 12583, code: "ME", name: "Montenegro"},
	"MS": {id: 12593, code: "MS", name: "Montserrat"},
	"MA": {id: 12605, code: "MA", name: "Morocco"},
	"MZ": {id: 12613, code: "MZ", name: "Mozambique"},
	"MM": {id: 12623, code: "MM", name: "Myanmar"},
	"NA": {id: 12633, code: "NA", name: "Namibia"},
	"NR": {id: 12643, code: "NR", name: "Nauru"},
	"NP": {id: 12653, code: "NP", name: "Nepal"},
	"NL": {id: 12664, code: "NL", name: "Netherlands"},
	"NC": {id: 12673, code: "NC", name: "New Caledonia"},
	"NZ": {id: 12683, code: "NZ", name: "New Zealand"},
	"NI": {id: 12693, code: "NI", name: "Nicaragua"},
	"NE": {id: 12703, code: "NE", name: "Niger"},
	"NG": {id: 12713, code: "NG", name: "Nigeria"},
	"NU": {id: 12723, code: "NU", name: "Niue"},
	"NF": {id: 12733, code: "NF", name: "Norfolk Island"},
	"MP": {id: 12743, code: "MP", name: "Northern Mariana Islands"},
	"NO": {id: 12754, code: "NO", name: "Norway"},
	"OM": {id: 12763, code: "OM", name: "Oman"},
	"PK": {id: 12773, code: "PK", name: "Pakistan"},
	"PW": {id: 12783, code: "PW", name: "Palau"},
	"PS": {id: 12793, code: "PS", name: "Palestine"},
	"PA": {id: 12803, code: "PA", name: "Panama"},
	"PG": {id: 12813, code: "PG", name: "Papua New Guinea"},
	"PY": {id: 12823, code: "PY", name: "Paraguay"},
	"PE": {id: 12834, code: "PE", name: "Peru"},
	"PH": {id: 12844, code: "PH", name: "Philippines"},
	"PN": {id: 12853, code: "PN", name: "Pitcairn"},
	"PL": {id: 12863, code: "PL", name: "Poland"},
	"PT": {id: 12873, code: "PT", name: "Portugal"},
	"PR": {id: 12883, code: "PR", name: "Puerto Rico"},
	"QA": {id: 12893, code: "QA", name: "Qatar"},
	"RE": {id: 12903, code: "RE", name: "Réunion"},
	"RO": {id: 12913, code: "RO", name: "Romania"},
	"RU": {id: 12923, code: "RU", name: "Russian Federation"},
	"RW": {id: 12933, code: "RW", name: "Rwanda"},
	"BL": {id: 12943, code: "BL", name: "Saint Barthélemy"},
	"SH": {id: 12953, code: "SH", name: "Saint Helena, Ascension and Tristan da Cunha"},
	"KN": {id: 12963, code: "KN", name: "Saint Kitts and Nevis"},
	"LC": {id: 12973, code: "LC", name: "Saint Lucia"},
	"MF": {id: 12983, code: "MF", name: "Saint Martin (French part)"},
	"PM": {id: 12993, code: "PM", name: "Saint Pierre and Miquelon"},
	"VC": {id: 13003, code: "VC", name: "Saint Vincent and the Grenadines"},
	"WS": {id: 13013, code: "WS", name: "Samoa"},
	"SM": {id: 13023, code: "SM", name: "San Marino"},
	"ST": {id: 13033, code: "ST", name: "Sao Tome and Principe"},
	"SA": {id: 13043, code: "SA", name: "Saudi Arabia"},
	"SN": {id: 13053, code: "SN", name: "Senegal"},
	"RS": {id: 13063, code: "RS", name: "Serbia"},
	"SC": {id: 13073, code: "SC", name: "Seychelles"},
	"SL": {id: 13083, code: "SL", name: "Sierra Leone"},
	"SG": {id: 13093, code: "SG", name: "Singapore"},
	"SX": {id: 13103, code: "SX", name: "Sint Maarten (Dutch part)"},
	"SK": {id: 13113, code: "SK", name: "Slovakia"},
	"SI": {id: 13123, code: "SI", name: "Slovenia"},
	"SB": {id: 13133, code: "SB", name: "Solomon Islands"},
	"SO": {id: 13143, code: "SO", name: "Somalia"},
	"ZA": {id: 13154, code: "ZA", name: "South Africa"},
	"GS": {id: 13163, code: "GS", name: "South Georgia and the South Sandwich Islands"},
	"SS": {id: 13173, code: "SS", name: "South Sudan"},
	"ES": {id: 13183, code: "ES", name: "Spain"},
	"LK": {id: 13193, code: "LK", name: "Sri Lanka"},
	"SD": {id: 13203, code: "SD", name: "Sudan"},
	"SR": {id: 13213, code: "SR", name: "Suriname"},
	"SJ": {id: 13223, code: "SJ", name: "Svalbard and Jan Mayen"},
	"SZ": {id: 13233, code: "SZ", name: "Swaziland"},
	"SE": {id: 13243, code: "SE", name: "Sweden"},
	"CH": {id: 13254, code: "CH", name: "Switzerland"},
	"SY": {id: 13263, code: "SY", name: "Syrian Arab Republic"},
	"TW": {id: 13273, code: "TW", name: "Taiwan"},
	"TJ": {id: 13283, code: "TJ", name: "Tajikistan"},
	"TZ": {id: 13293, code: "TZ", name: "Tanzania"},
	"TH": {id: 13303, code: "TH", name: "Thailand"},
	"TL": {id: 13313, code: "TL", name: "Timor-Leste"},
	"TG": {id: 13323, code: "TG", name: "Togo"},
	"TK": {id: 13333, code: "TK", name: "Tokelau"},
	"TO": {id: 13343, code: "TO", name: "Tonga"},
	"TT": {id: 13353, code: "TT", name: "Trinidad and Tobago"},
	"TN": {id: 13363, code: "TN", name: "Tunisia"},
	"TR": {id: 13373, code: "TR", name: "Turkey"},
	"TM": {id: 13383, code: "TM", name: "Turkmenistan"},
	"TC": {id: 13393, code: "TC", name: "Turks and Caicos Islands"},
	"TV": {id: 13403, code: "TV", name: "Tuvalu"},
	"UG": {id: 13413, code: "UG", name: "Uganda"},
	"UA": {id: 13423, code: "UA", name: "Ukraine"},
	"AE": {id: 13433, code: "AE", name: "United Arab Emirates"},
	"GB": {id: 13443, code: "GB", name: "United Kingdom"},
	"US": {id: 13453, code: "US", name: "United States"},
	"UM": {id: 13463, code: "UM", name: "United States Minor Outlying Islands"},
	"UY": {id: 13473, code: "UY", name: "Uruguay"},
	"UZ": {id: 13483, code: "UZ", name: "Uzbekistan"},
	"VU": {id: 13493, code: "VU", name: "Vanuatu"},
	"VE": {id: 13503, code: "VE", name: "Venezuela"},
	"VN": {id: 13513, code: "VN", name: "Viet Nam"},
	"VG": {id: 13523, code: "VG", name: "Virgin Islands, British"},
	"VI": {id: 13533, code: "VI", name: "Virgin Islands, U.S."},
	"WF": {id: 13543, code: "WF", name: "Wallis and Futuna"},
	"EH": {id: 13553, code: "EH", name: "Western Sahara"},
	"YE": {id: 13563, code: "YE", name: "Yemen"},
	"ZM": {id: 13573, code: "ZM", name: "Zambia"},
	"ZW": {id: 13583, code: "ZW", name: "Zimbabwe"},
}

/* Missing:
-en-CB      English (Caribbean)
-en-029 	English (Caribbean)
-hsb-DE     Upper Sorbian (Germany)
-iu-Latn-CA Inuktitut (Latin) (Canada)
-kok-IN	    Konkani (India)
-moh-CA     Mohawk (Canada)
-iu-Cans-CA Inuktitut (Syllabics) (Canada)
-nso-ZA     Sesotho sa Leboa (South Africa)
-ns-za      ????
-sah-RU	    Yakut (Russia)
-qut-GT     K'iche (Guatemala)
-quz-BO     Quechua (Bolivia)
-quz-EC     Quechua (Ecuador)
-quz-PE     Quechua (Peru)
*/
var languageCountry = map[string]bool{
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
	"az-az": true,
	"as-in": true,
	"ba-ru": true,
	"be-by": true,
	"bg-bg": true,
	"bn-bd": true,
	"bn-in": true,
	"bo-cn": true,
	"br-fr": true,
	"bs-ba": true,
	"ca-es": true,
	"co-fr": true,
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
	"en-in": true,
	"en-my": true,
	"en-nz": true,
	"en-ph": true,
	"en-sg": true,
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
	"fy-nl": true,
	"ga-ie": true,
	"gd-gb": true,
	"gl-es": true,
	"gu-in": true,
	"ha-ng": true,
	"he-il": true,
	"hi-in": true,
	"hr-ba": true,
	"hr-hr": true,
	"hu-hu": true,
	"hy-am": true,
	"id-id": true,
	"ig-ng": true,
	"ii-cn": true,
	"is-is": true,
	"it-ch": true,
	"it-it": true,
	"ja-jp": true,
	"ka-ge": true,
	"kk-kz": true,
	"kl-gl": true,
	"km-kh": true,
	"kn-in": true,
	"ko-kr": true,
	"ky-kg": true,
	"lb-lu": true,
	"lo-la": true,
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
	"ne-np": true,
	"nl-be": true,
	"nl-nl": true,
	"nn-no": true,
	"oc-fr": true,
	"or-in": true,
	"pa-in": true,
	"pl-pl": true,
	"pt-br": true,
	"pt-pt": true,
	"ps-af": true,
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
	"tk-tm": true,
	"tn-za": true,
	"tr-tr": true,
	"tt-ru": true,
	"ug-cn": true,
	"uk-ua": true,
	"ur-pk": true,
	"vi-vn": true,
	"wo-sn": true,
	"xh-za": true,
	"yo-ng": true,
	"zh-cn": true,
	"zh-hk": true,
	"zh-mo": true,
	"zh-sg": true,
	"zh-tw": true,
	"zu-za": true,
}

func init() {
	tempLang := new(Language)
	errorMessage := "Language.init() - invalid key: "
	// Complete language list with combination of language-country
	//-> languages
	for key, _ := range languageCountry {
		lang := tempLang.getLanguage(key, true)
		if lang == nil {
			panic(errorMessage + key)
		}
		lang = lang.Clone()
		lang.currentCountry = lang.getCountry(key)
		lang.code += languageCountrySeparator + lang.currentCountry.code

		if lang.currentCountry == nil {
			panic(errorMessage + key)
		}
		lang.id += lang.currentCountry.id
		languages[key] = lang
	}
}

func (language *Language) Init(code string) {
	lang := language.getLanguage(code, false)
	if lang != nil {
		lang.copyTo(language)
	}
}

//******************************
// getters and setters
//******************************
func (language *Language) GetId() int32 {
	return int32(language.id)
}

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
func (language *Language) GetLanguageCode() string {
	if language.currentCountry != nil && strings.Index(language.code, languageCountrySeparator) >= 0 {
		return language.code[:strings.Index(language.code, languageCountrySeparator)]
	}
	return language.code
}

func (language *Language) GetCountryCode() string {
	if language.currentCountry != nil {
		return language.currentCountry.code
	}
	return ""
}

func (language *Language) GetDescription() string {
	var result = ""
	result += language.name
	if language.currentCountry != nil {
		//eg. English (United States)
		result += " (" + language.currentCountry.name + ")"
	}
	return result
}

func (language *Language) IsCodeValid(code string) (bool, error) {
	var codeFormat = strings.ReplaceAll(code, " ", "")
	if len(code) <= 0 {
		return false, errors.New(languageErrorEmpty)
	}

	if strings.Contains(codeFormat, languageCountrySeparator) == true {
		var country = language.getCountry(code)
		if country == nil {
			return false, errors.New(fmt.Sprintf(languageErrorInvalidCountry, code))
		}
		// language, country combination exists
		if _, ok := languageCountry[strings.ToLower(codeFormat)]; !ok {
			return false, errors.New(fmt.Sprintf(languageErrorInvalidCounLang, code))
		}
	}
	var lang = language.getLanguage(code, true)
	// just language
	if lang == nil {
		return false, errors.New(fmt.Sprintf(languageErrorInvalidLang, code))
	}
	return true, nil
}

func (language *Language) GetList() []Language {
	return language.getList()
}

func (language *Language) Clone() *Language {
	result := new(Language)
	language.copyTo(result)
	return result
}

func (language *Language) String() string {
	// languageToStringFormat: "id: %d; code: %s; name: %s; nativeName: %s"
	return fmt.Sprintf(languageToStringFormat, language.id, language.code, language.name, language.GetLanguageCode())
}

//******************************
// private methods
//******************************
func (source *Language) copyTo(target *Language) {
	target.id = source.id
	target.code = source.code
	target.name = source.name
	target.nativeName = source.nativeName
	if source.currentCountry != nil {
		target.currentCountry = new(country)
		target.currentCountry.id = source.currentCountry.id
		target.currentCountry.name = source.currentCountry.name
		target.currentCountry.code = source.currentCountry.code
	} else {
		target.currentCountry = nil
	}
	target.valid = source.valid
}

func (language *Language) getList() []Language {
	var result = make([]Language, 0, len(languages))
	for _, value := range languages {
		result = append(result, *value)
	}
	return result
}

func (language *Language) exists(schema *Schema) bool {
	query := new(metaQuery)
	query.setSchema(metaSchemaName)
	query.setTable(metaTableName)

	query.addFilter(metaFieldId, operatorEqual, language.id)
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

func (language *Language) toMeta() *meta {
	var metaTable = new(meta)

	// key
	metaTable.id = language.GetId()
	metaTable.name = language.code
	metaTable.description = language.GetDescription()
	metaTable.objectType = int8(entitytype.Language)
	metaTable.dataType = 0
	metaTable.setEntityBaseline(true)
	metaTable.enabled = true

	return metaTable
}

func (language *Language) getCountry(code string) *country {
	var formattedCode = strings.ReplaceAll(code, " ", "")
	var index = strings.Index(formattedCode, languageCountrySeparator)

	if index > 0 && index+1 < len(formattedCode) {
		var countryCode = formattedCode[index+1:]
		if val, ok := countries[strings.ToUpper(countryCode)]; ok {
			return &val
		}
	}
	return nil
}

func (language *Language) getLanguage(code string, languageOnly bool) *Language {
	var formattedCode = strings.ReplaceAll(code, " ", "")
	var languageCode string = formattedCode

	if languageOnly {
		if len(languageCode) >= 2 {
			languageCode = formattedCode[:2]
		}
	}
	if val, ok := languages[strings.ToLower(languageCode)]; ok {
		return val
	}
	return nil
}
