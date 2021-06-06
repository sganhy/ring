package tabletype

type TableType int8

const (
	Business             TableType = 1
	Meta                 TableType = 3
	MetaId               TableType = 4
	Fake                 TableType = 5
	Mtm                  TableType = 6
	Log                  TableType = 7
	Lexicon              TableType = 8
	LexiconItem          TableType = 9
	SchemaDictionary     TableType = 22
	TableDictionary      TableType = 24
	TableSpaceDictionary TableType = 25
	Logical              TableType = 27
)

func (tableType *TableType) String() string {
	tabletype := *tableType
	switch tabletype {
	case Business:
		return "business"
	case Meta:
		return "meta"
	case MetaId:
		return "metaId"
	case Fake:
		return "fake"
	case Mtm:
		return "mtm"
	case Log:
		return "log"
	case Lexicon:
		return "lexicon"
	case LexiconItem:
		return "lexiconItem"
	case SchemaDictionary:
		return "schemaDictionary"
	case TableDictionary:
		return "tableDictionary"
	case Logical:
		return "logical"
	}

	return ""
}
