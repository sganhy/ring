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
	User                 TableType = 10
	SchemaDictionary     TableType = 22
	TableDictionary      TableType = 24
	TableSpaceDictionary TableType = 25
)
