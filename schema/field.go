package schema

import (
	"errors"
	"fmt"
	"ring/schema/databaseprovider"
	"ring/schema/entitytype"
	"ring/schema/fieldtype"
	"ring/schema/searchabletype"
	"ring/schema/tabletype"
	"strconv"
	"strings"
	"time"
	"unicode"

	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

type Field struct {
	id            int32
	name          string
	description   string
	fieldType     fieldtype.FieldType
	size          uint16
	defaultValue  string
	baseline      bool
	notNull       bool
	caseSensitive bool
	multilingual  bool
	active        bool
}

const (
	defaultNumberValue         string = "0"
	doubleQuotes               string = "\""
	defaultBooleanValue        string = "false"
	defaultDateTimeValue       string = "0001-01-01T00:00:00.000"
	defaultShortDateTimeValue  string = "0001-01-01"
	defaultLongDateTimeValue   string = "0001-01-01T00:00:00Z"
	postgreSqlByteConstraint   string = " CHECK (%s >= '-128' AND %s <= 127)"
	maxInt08                   string = "127"
	maxInt16                   string = "32767"
	maxInt32                   string = "2147483647"
	maxInt64                   string = "9223372036854775807"
	minInt08                   string = "-128"
	minInt16                   string = "-32768"
	minInt32                   string = "-2147483648"
	minInt64                   string = "-9223372036854775808"
	defaultTimeFormat          string = "2006-01-02T15:04:05.000" // rfc3339
	defaultShortTimeFormat     string = "2006-01-02"              // rfc3339
	unknownFieldDataType       string = ""
	invalidValue               string = ""
	primaryKeyFieldName        string = "id"
	primaryKeyDesc             string = "Internal record number"
	errorInvalidValueType      string = "Invalid value type"
	errorInvalidDateTimeFormat string = "Invalid Date/Time format"
	postgreVarcharMaxSize      uint16 = 65535
	mySqlVarcharMaxSize        uint16 = 65535
	sqliteVarcharMaxSize       int64  = 1000000000
	fieldToStringFormat        string = "name=%s; description=%s; type=%s; defaultValue=%s; baseline=%t; notNull=%t; caseSensitive=%t; active=%t"
)

var defaultPrimaryKeyInt64 *Field = nil
var defaultPrimaryKeyInt32 *Field = nil
var defaultPrimaryKeyInt16 *Field = nil

var postgreDataType = map[fieldtype.FieldType]string{
	fieldtype.String:        "varchar",
	fieldtype.Double:        "float8",
	fieldtype.Float:         "float4",
	fieldtype.Long:          "int8",
	fieldtype.Int:           "int4",
	fieldtype.Short:         "int2",
	fieldtype.Byte:          "int2",
	fieldtype.Boolean:       "bool",
	fieldtype.ShortDateTime: "date",
	fieldtype.DateTime:      "timestamp without time zone",
	fieldtype.LongDateTime:  "timestamp with time zone"}

// sql key words ==> SQL:2003,SQL:1999,SQL-92, and PostgreSQL
var sqlKeyWords = map[string]bool{
	"A":                                true,
	"ABORT":                            true,
	"ABS":                              true,
	"ABSOLUTE":                         true,
	"ACCESS":                           true,
	"ACTION":                           true,
	"ADA":                              true,
	"ADD":                              true,
	"ADMIN":                            true,
	"AFTER":                            true,
	"AGGREGATE":                        true,
	"ALIAS":                            true,
	"ALL":                              true,
	"ALLOCATE":                         true,
	"ALSO":                             true,
	"ALTER":                            true,
	"ALWAYS":                           true,
	"ANALYSE":                          true,
	"ANALYZE":                          true,
	"AND":                              true,
	"ANY":                              true,
	"ARE":                              true,
	"ARRAY":                            true,
	"AS":                               true,
	"ASC":                              true,
	"ASENSITIVE":                       true,
	"ASSERTION":                        true,
	"ASSIGNMENT":                       true,
	"ASYMMETRIC":                       true,
	"AT":                               true,
	"ATOMIC":                           true,
	"ATTRIBUTE":                        true,
	"ATTRIBUTES":                       true,
	"AUTHORIZATION":                    true,
	"AVG":                              true,
	"BACKWARD":                         true,
	"BEFORE":                           true,
	"BEGIN":                            true,
	"BERNOULLI":                        true,
	"BETWEEN":                          true,
	"BIGINT":                           true,
	"BINARY":                           true,
	"BIT":                              true,
	"BITVAR":                           true,
	"BIT_LENGTH":                       true,
	"BLOB":                             true,
	"BOOLEAN":                          true,
	"BOTH":                             true,
	"BREADTH":                          true,
	"BY":                               true,
	"C":                                true,
	"CACHE":                            true,
	"CALL":                             true,
	"CALLED":                           true,
	"CARDINALITY":                      true,
	"CASCADE":                          true,
	"CASCADED":                         true,
	"CASE":                             true,
	"CAST":                             true,
	"CATALOG":                          true,
	"CATALOG_NAME":                     true,
	"CEIL":                             true,
	"CEILING":                          true,
	"CHAIN":                            true,
	"CHAR":                             true,
	"CHARACTER":                        true,
	"CHARACTERISTICS":                  true,
	"CHARACTERS":                       true,
	"CHARACTER_LENGTH":                 true,
	"CHARACTER_SET_CATALOG":            true,
	"CHARACTER_SET_NAME":               true,
	"CHARACTER_SET_SCHEMA":             true,
	"CHAR_LENGTH":                      true,
	"CHECK":                            true,
	"CHECKED":                          true,
	"CHECKPOINT":                       true,
	"CLASS":                            true,
	"CLASS_ORIGIN":                     true,
	"CLOB":                             true,
	"CLOSE":                            true,
	"CLUSTER":                          true,
	"COALESCE":                         true,
	"COBOL":                            true,
	"COLLATE":                          true,
	"COLLATION":                        true,
	"COLLATION_CATALOG":                true,
	"COLLATION_NAME":                   true,
	"COLLATION_SCHEMA":                 true,
	"COLLECT":                          true,
	"COLUMN":                           true,
	"COLUMN_NAME":                      true,
	"COMMAND_FUNCTION":                 true,
	"COMMAND_FUNCTION_CODE":            true,
	"COMMENT":                          true,
	"COMMIT":                           true,
	"COMMITTED":                        true,
	"COMPLETION":                       true,
	"CONDITION":                        true,
	"CONDITION_NUMBER":                 true,
	"CONNECT":                          true,
	"CONNECTION":                       true,
	"CONNECTION_NAME":                  true,
	"CONSTRAINT":                       true,
	"CONSTRAINTS":                      true,
	"CONSTRAINT_CATALOG":               true,
	"CONSTRAINT_NAME":                  true,
	"CONSTRAINT_SCHEMA":                true,
	"CONSTRUCTOR":                      true,
	"CONTAINS":                         true,
	"CONTINUE":                         true,
	"CONVERSION":                       true,
	"CONVERT":                          true,
	"COPY":                             true,
	"CORR":                             true,
	"CORRESPONDING":                    true,
	"COUNT":                            true,
	"COVAR_POP":                        true,
	"COVAR_SAMP":                       true,
	"CREATE":                           true,
	"CREATEDB":                         true,
	"CREATEROLE":                       true,
	"CREATEUSER":                       true,
	"CROSS":                            true,
	"CSV":                              true,
	"CUBE":                             true,
	"CUME_DIST":                        true,
	"CURRENT":                          true,
	"CURRENT_DATE":                     true,
	"CURRENT_DEFAULT_TRANSFORM_GROUP":  true,
	"CURRENT_PATH":                     true,
	"CURRENT_ROLE":                     true,
	"CURRENT_TIME":                     true,
	"CURRENT_TIMESTAMP":                true,
	"CURRENT_TRANSFORM_GROUP_FOR_TYPE": true,
	"CURRENT_USER":                     true,
	"CURSOR":                           true,
	"CURSOR_NAME":                      true,
	"CYCLE":                            true,
	"DATA":                             true,
	"DATABASE":                         true,
	"DATE":                             true,
	"DATETIME_INTERVAL_CODE":           true,
	"DATETIME_INTERVAL_PRECISION":      true,
	"DAY":                              true,
	"DEALLOCATE":                       true,
	"DEC":                              true,
	"DECIMAL":                          true,
	"DECLARE":                          true,
	"DEFAULT":                          true,
	"DEFAULTS":                         true,
	"DEFERRABLE":                       true,
	"DEFERRED":                         true,
	"DEFINED":                          true,
	"DEFINER":                          true,
	"DEGREE":                           true,
	"DELETE":                           true,
	"DELIMITER":                        true,
	"DELIMITERS":                       true,
	"DENSE_RANK":                       true,
	"DEPTH":                            true,
	"DEREF":                            true,
	"DERIVED":                          true,
	"DESC":                             true,
	"DESCRIBE":                         true,
	"DESCRIPTOR":                       true,
	"DESTROY":                          true,
	"DESTRUCTOR":                       true,
	"DETERMINISTIC":                    true,
	"DIAGNOSTICS":                      true,
	"DICTIONARY":                       true,
	"DISABLE":                          true,
	"DISCONNECT":                       true,
	"DISPATCH":                         true,
	"DISTINCT":                         true,
	"DO":                               true,
	"DOMAIN":                           true,
	"DOUBLE":                           true,
	"DROP":                             true,
	"DYNAMIC":                          true,
	"DYNAMIC_FUNCTION":                 true,
	"DYNAMIC_FUNCTION_CODE":            true,
	"EACH":                             true,
	"ELEMENT":                          true,
	"ELSE":                             true,
	"ENABLE":                           true,
	"ENCODING":                         true,
	"ENCRYPTED":                        true,
	"END":                              true,
	"END-EXEC":                         true,
	"EQUALS":                           true,
	"ESCAPE":                           true,
	"EVERY":                            true,
	"EXCEPT":                           true,
	"EXCEPTION":                        true,
	"EXCLUDE":                          true,
	"EXCLUDING":                        true,
	"EXCLUSIVE":                        true,
	"EXEC":                             true,
	"EXECUTE":                          true,
	"EXISTING":                         true,
	"EXISTS":                           true,
	"EXP":                              true,
	"EXPLAIN":                          true,
	"EXTERNAL":                         true,
	"EXTRACT":                          true,
	"FALSE":                            true,
	"FETCH":                            true,
	"FILTER":                           true,
	"FINAL":                            true,
	"FIRST":                            true,
	"FLOAT":                            true,
	"FLOOR":                            true,
	"FOLLOWING":                        true,
	"FOR":                              true,
	"FORCE":                            true,
	"FOREIGN":                          true,
	"FORTRAN":                          true,
	"FORWARD":                          true,
	"FOUND":                            true,
	"FREE":                             true,
	"FREEZE":                           true,
	"FROM":                             true,
	"FULL":                             true,
	"FUNCTION":                         true,
	"FUSION":                           true,
	"G":                                true,
	"GENERAL":                          true,
	"GENERATED":                        true,
	"GET":                              true,
	"GLOBAL":                           true,
	"GO":                               true,
	"GOTO":                             true,
	"GRANT":                            true,
	"GRANTED":                          true,
	"GREATEST":                         true,
	"GROUP":                            true,
	"GROUPING":                         true,
	"HANDLER":                          true,
	"HAVING":                           true,
	"HEADER":                           true,
	"HIERARCHY":                        true,
	"HOLD":                             true,
	"HOST":                             true,
	"HOUR":                             true,
	"IDENTITY":                         true,
	"IGNORE":                           true,
	"ILIKE":                            true,
	"IMMEDIATE":                        true,
	"IMMUTABLE":                        true,
	"IMPLEMENTATION":                   true,
	"IMPLICIT":                         true,
	"IN":                               true,
	"INCLUDING":                        true,
	"INCREMENT":                        true,
	"INDEX":                            true,
	"INDICATOR":                        true,
	"INFIX":                            true,
	"INHERIT":                          true,
	"INHERITS":                         true,
	"INITIALIZE":                       true,
	"INITIALLY":                        true,
	"INNER":                            true,
	"INOUT":                            true,
	"INPUT":                            true,
	"INSENSITIVE":                      true,
	"INSERT":                           true,
	"INSTANCE":                         true,
	"INSTANTIABLE":                     true,
	"INSTEAD":                          true,
	"INT":                              true,
	"INTEGER":                          true,
	"INTERSECT":                        true,
	"INTERSECTION":                     true,
	"INTERVAL":                         true,
	"INTO":                             true,
	"INVOKER":                          true,
	"IS":                               true,
	"ISNULL":                           true,
	"ISOLATION":                        true,
	"ITERATE":                          true,
	"JOIN":                             true,
	"K":                                true,
	"KEY":                              true,
	"KEY_MEMBER":                       true,
	"KEY_TYPE":                         true,
	"LANCOMPILER":                      true,
	"LANGUAGE":                         true,
	"LARGE":                            true,
	"LAST":                             true,
	"LATERAL":                          true,
	"LEADING":                          true,
	"LEAST":                            true,
	"LEFT":                             true,
	"LENGTH":                           true,
	"LESS":                             true,
	"LEVEL":                            true,
	"LIKE":                             true,
	"LIMIT":                            true,
	"LISTEN":                           true,
	"LN":                               true,
	"LOAD":                             true,
	"LOCAL":                            true,
	"LOCALTIME":                        true,
	"LOCALTIMESTAMP":                   true,
	"LOCATION":                         true,
	"LOCATOR":                          true,
	"LOCK":                             true,
	"LOGIN":                            true,
	"LOWER":                            true,
	"M":                                true,
	"MAP":                              true,
	"MATCH":                            true,
	"MATCHED":                          true,
	"MAX":                              true,
	"MAXVALUE":                         true,
	"MEMBER":                           true,
	"MERGE":                            true,
	"MESSAGE_LENGTH":                   true,
	"MESSAGE_OCTET_LENGTH":             true,
	"MESSAGE_TEXT":                     true,
	"METHOD":                           true,
	"MIN":                              true,
	"MINUTE":                           true,
	"MINVALUE":                         true,
	"MOD":                              true,
	"MODE":                             true,
	"MODIFIES":                         true,
	"MODIFY":                           true,
	"MODULE":                           true,
	"MONTH":                            true,
	"MORE":                             true,
	"MOVE":                             true,
	"MULTISET":                         true,
	"MUMPS":                            true,
	"NAME":                             true,
	"NAMES":                            true,
	"NATIONAL":                         true,
	"NATURAL":                          true,
	"NCHAR":                            true,
	"NCLOB":                            true,
	"NESTING":                          true,
	"NEW":                              true,
	"NEXT":                             true,
	"NO":                               true,
	"NOCREATEDB":                       true,
	"NOCREATEROLE":                     true,
	"NOCREATEUSER":                     true,
	"NOINHERIT":                        true,
	"NOLOGIN":                          true,
	"NONE":                             true,
	"NORMALIZE":                        true,
	"NORMALIZED":                       true,
	"NOSUPERUSER":                      true,
	"NOT":                              true,
	"NOTHING":                          true,
	"NOTIFY":                           true,
	"NOTNULL":                          true,
	"NOWAIT":                           true,
	"NULL":                             true,
	"NULLABLE":                         true,
	"NULLIF":                           true,
	"NULLS":                            true,
	"NUMBER":                           true,
	"NUMERIC":                          true,
	"OBJECT":                           true,
	"OCTETS":                           true,
	"OCTET_LENGTH":                     true,
	"OF":                               true,
	"OFF":                              true,
	"OFFSET":                           true,
	"OIDS":                             true,
	"OLD":                              true,
	"ON":                               true,
	"ONLY":                             true,
	"OPEN":                             true,
	"OPERATION":                        true,
	"OPERATOR":                         true,
	"OPTION":                           true,
	"OPTIONS":                          true,
	"OR":                               true,
	"ORDER":                            true,
	"ORDERING":                         true,
	"ORDINALITY":                       true,
	"OTHERS":                           true,
	"OUT":                              true,
	"OUTER":                            true,
	"OUTPUT":                           true,
	"OVER":                             true,
	"OVERLAPS":                         true,
	"OVERLAY":                          true,
	"OVERRIDING":                       true,
	"OWNER":                            true,
	"PAD":                              true,
	"PARAMETER":                        true,
	"PARAMETERS":                       true,
	"PARAMETER_MODE":                   true,
	"PARAMETER_NAME":                   true,
	"PARAMETER_ORDINAL_POSITION":       true,
	"PARAMETER_SPECIFIC_CATALOG":       true,
	"PARAMETER_SPECIFIC_NAME":          true,
	"PARAMETER_SPECIFIC_SCHEMA":        true,
	"PARTIAL":                          true,
	"PARTITION":                        true,
	"PASCAL":                           true,
	"PASSWORD":                         true,
	"PATH":                             true,
	"PERCENTILE_CONT":                  true,
	"PERCENTILE_DISC":                  true,
	"PERCENT_RANK":                     true,
	"PLACING":                          true,
	"PLI":                              true,
	"POSITION":                         true,
	"POSTFIX":                          true,
	"POWER":                            true,
	"PRECEDING":                        true,
	"PRECISION":                        true,
	"PREFIX":                           true,
	"PREORDER":                         true,
	"PREPARE":                          true,
	"PREPARED":                         true,
	"PRESERVE":                         true,
	"PRIMARY":                          true,
	"PRIOR":                            true,
	"PRIVILEGES":                       true,
	"PROCEDURAL":                       true,
	"PROCEDURE":                        true,
	"PUBLIC":                           true,
	"QUOTE":                            true,
	"RANGE":                            true,
	"RANK":                             true,
	"READ":                             true,
	"READS":                            true,
	"REAL":                             true,
	"RECHECK":                          true,
	"RECURSIVE":                        true,
	"REF":                              true,
	"REFERENCES":                       true,
	"REFERENCING":                      true,
	"REGR_AVGX":                        true,
	"REGR_AVGY":                        true,
	"REGR_COUNT":                       true,
	"REGR_INTERCEPT":                   true,
	"REGR_R2":                          true,
	"REGR_SLOPE":                       true,
	"REGR_SXX":                         true,
	"REGR_SXY":                         true,
	"REGR_SYY":                         true,
	"REINDEX":                          true,
	"RELATIVE":                         true,
	"RELEASE":                          true,
	"RENAME":                           true,
	"REPEATABLE":                       true,
	"REPLACE":                          true,
	"RESET":                            true,
	"RESTART":                          true,
	"RESTRICT":                         true,
	"RESULT":                           true,
	"RETURN":                           true,
	"RETURNED_CARDINALITY":             true,
	"RETURNED_LENGTH":                  true,
	"RETURNED_OCTET_LENGTH":            true,
	"RETURNED_SQLSTATE":                true,
	"RETURNS":                          true,
	"REVOKE":                           true,
	"RIGHT":                            true,
	"ROLE":                             true,
	"ROLLBACK":                         true,
	"ROLLUP":                           true,
	"ROUTINE":                          true,
	"ROUTINE_CATALOG":                  true,
	"ROUTINE_NAME":                     true,
	"ROUTINE_SCHEMA":                   true,
	"ROW":                              true,
	"ROWS":                             true,
	"ROW_COUNT":                        true,
	"ROW_NUMBER":                       true,
	"RULE":                             true,
	"SAVEPOINT":                        true,
	"SCALE":                            true,
	"SCHEMA":                           true,
	"SCHEMA_NAME":                      true,
	"SCOPE":                            true,
	"SCOPE_CATALOG":                    true,
	"SCOPE_NAME":                       true,
	"SCOPE_SCHEMA":                     true,
	"SCROLL":                           true,
	"SEARCH":                           true,
	"SECOND":                           true,
	"SECTION":                          true,
	"SECURITY":                         true,
	"SELECT":                           true,
	"SELF":                             true,
	"SENSITIVE":                        true,
	"SEQUENCE":                         true,
	"SERIALIZABLE":                     true,
	"SERVER_NAME":                      true,
	"SESSION":                          true,
	"SESSION_USER":                     true,
	"SET":                              true,
	"SETOF":                            true,
	"SETS":                             true,
	"SHARE":                            true,
	"SHOW":                             true,
	"SIMILAR":                          true,
	"SIMPLE":                           true,
	"SIZE":                             true,
	"SMALLINT":                         true,
	"SOME":                             true,
	"SOURCE":                           true,
	"SPACE":                            true,
	"SPECIFIC":                         true,
	"SPECIFICTYPE":                     true,
	"SPECIFIC_NAME":                    true,
	"SQL":                              true,
	"SQLCODE":                          true,
	"SQLERROR":                         true,
	"SQLEXCEPTION":                     true,
	"SQLSTATE":                         true,
	"SQLWARNING":                       true,
	"SQRT":                             true,
	"STABLE":                           true,
	"START":                            true,
	"STATE":                            true,
	"STATEMENT":                        true,
	"STATIC":                           true,
	"STATISTICS":                       true,
	"STDDEV_POP":                       true,
	"STDDEV_SAMP":                      true,
	"STDIN":                            true,
	"STDOUT":                           true,
	"STORAGE":                          true,
	"STRICT":                           true,
	"STRUCTURE":                        true,
	"STYLE":                            true,
	"SUBCLASS_ORIGIN":                  true,
	"SUBLIST":                          true,
	"SUBMULTISET":                      true,
	"SUBSTRING":                        true,
	"SUM":                              true,
	"SUPERUSER":                        true,
	"SYMMETRIC":                        true,
	"SYSID":                            true,
	"SYSTEM":                           true,
	"SYSTEM_USER":                      true,
	"TABLE":                            true,
	"TABLESAMPLE":                      true,
	"TABLESPACE":                       true,
	"TABLE_NAME":                       true,
	"TEMP":                             true,
	"TEMPLATE":                         true,
	"TEMPORARY":                        true,
	"TERMINATE":                        true,
	"THAN":                             true,
	"THEN":                             true,
	"TIES":                             true,
	"TIME":                             true,
	"TIMESTAMP":                        true,
	"TIMEZONE_HOUR":                    true,
	"TIMEZONE_MINUTE":                  true,
	"TO":                               true,
	"TOAST":                            true,
	"TOP_LEVEL_COUNT":                  true,
	"TRAILING":                         true,
	"TRANSACTION":                      true,
	"TRANSACTIONS_COMMITTED":           true,
	"TRANSACTIONS_ROLLED_BACK":         true,
	"TRANSACTION_ACTIVE":               true,
	"TRANSFORM":                        true,
	"TRANSFORMS":                       true,
	"TRANSLATE":                        true,
	"TRANSLATION":                      true,
	"TREAT":                            true,
	"TRIGGER":                          true,
	"TRIGGER_CATALOG":                  true,
	"TRIGGER_NAME":                     true,
	"TRIGGER_SCHEMA":                   true,
	"TRIM":                             true,
	"TRUE":                             true,
	"TRUNCATE":                         true,
	"TRUSTED":                          true,
	"TYPE":                             true,
	"UESCAPE":                          true,
	"UNBOUNDED":                        true,
	"UNCOMMITTED":                      true,
	"UNDER":                            true,
	"UNENCRYPTED":                      true,
	"UNION":                            true,
	"UNIQUE":                           true,
	"UNKNOWN":                          true,
	"UNLISTEN":                         true,
	"UNNAMED":                          true,
	"UNNEST":                           true,
	"UNTIL":                            true,
	"UPDATE":                           true,
	"UPPER":                            true,
	"USAGE":                            true,
	"USER":                             true,
	"USER_DEFINED_TYPE_CATALOG":        true,
	"USER_DEFINED_TYPE_CODE":           true,
	"USER_DEFINED_TYPE_NAME":           true,
	"USER_DEFINED_TYPE_SCHEMA":         true,
	"USING":                            true,
	"VACUUM":                           true,
	"VALID":                            true,
	"VALIDATOR":                        true,
	"VALUE":                            true,
	"VALUES":                           true,
	"VARCHAR":                          true,
	"VARIABLE":                         true,
	"VARYING":                          true,
	"VAR_POP":                          true,
	"VAR_SAMP":                         true,
	"VERBOSE":                          true,
	"VIEW":                             true,
	"VOLATILE":                         true,
	"WHEN":                             true,
	"WHENEVER":                         true,
	"WHERE":                            true,
	"WIDTH_BUCKET":                     true,
	"WINDOW":                           true,
	"WITH":                             true,
	"WITHIN":                           true,
	"WITHOUT":                          true,
	"WORK":                             true,
	"WRITE":                            true,
	"YEAR":                             true,
	"ZONE":                             true,
}

func init() {
	//64
	defaultPrimaryKeyInt64 = new(Field)
	defaultPrimaryKeyInt64.Init(0, primaryKeyFieldName, primaryKeyDesc, fieldtype.Long, 0, "", false, true, true, false, true)
	//32
	defaultPrimaryKeyInt32 = new(Field)
	defaultPrimaryKeyInt32.Init(0, primaryKeyFieldName, primaryKeyDesc, fieldtype.Int, 0, "", false, true, true, false, true)
	//16
	defaultPrimaryKeyInt16 = new(Field)
	defaultPrimaryKeyInt16.Init(0, primaryKeyFieldName, primaryKeyDesc, fieldtype.Short, 0, "", false, true, true, false, true)

}

func (field *Field) Init(id int32, name string, description string, fieldType fieldtype.FieldType, size uint32,
	defaultValue string, baseline bool, notNull bool, caseSensitive bool, multilingual bool, active bool) {
	field.id = id
	field.name = name
	field.description = description
	field.fieldType = fieldType
	if size > 65535 {
		field.size = 0
	} else {
		field.size = uint16(size)
	}
	field.baseline = baseline
	field.notNull = notNull
	field.active = active
	field.multilingual = multilingual
	field.caseSensitive = caseSensitive
	//!!! at the end only
	field.defaultValue = field.getDefaultValue(defaultValue)
}

//******************************
// getters
//******************************
func (field *Field) GetId() int32 {
	return field.id
}

func (field *Field) GetName() string {
	return field.name
}

func (field *Field) GetDescription() string {
	return field.description
}

func (field *Field) GetType() fieldtype.FieldType {
	return field.fieldType
}

func (field *Field) GetSize() uint32 {
	return uint32(field.size)
}

func (field *Field) GetDefaultValue() string {
	return field.defaultValue
}

func (field *Field) IsBaseline() bool {
	return field.baseline
}

func (field *Field) IsNotNull() bool {
	return field.notNull
}

func (field *Field) IsCaseSensitive() bool {
	return field.caseSensitive
}

func (field *Field) IsActive() bool {
	return field.active
}

func (field *Field) IsMultilingual() bool {
	return field.multilingual
}

//******************************
// public methods
//******************************
func (field *Field) IsValid() bool {
	// compare addresses
	if defaultPrimaryKeyInt64 == field || defaultPrimaryKeyInt32 == field || defaultPrimaryKeyInt16 == field {
		return true
	}
	return field.id > 0
}

func (field *Field) IsPrimaryKey() bool {
	// compare addresses
	return defaultPrimaryKeyInt64 == field || defaultPrimaryKeyInt32 == field || defaultPrimaryKeyInt16 == field
}

func (field *Field) IsNumeric() bool {
	return field.fieldType == fieldtype.Long || field.fieldType == fieldtype.Int ||
		field.fieldType == fieldtype.Short || field.fieldType == fieldtype.Byte ||
		field.fieldType == fieldtype.Float || field.fieldType == fieldtype.Double

}

func (field *Field) IsDateTime() bool {
	return field.fieldType == fieldtype.DateTime || field.fieldType == fieldtype.LongDateTime || field.fieldType == fieldtype.ShortDateTime
}

///
/// Calculate searchable field value (remove diacritic characters and value.ToUpper())
///
func (field *Field) GetSearchableValue(value string, searchableType searchabletype.SearchableType) string {
	//TODO specific treatment by language
	t := transform.Chain(norm.NFD, runes.Remove(runes.In(unicode.Mn)), norm.NFC)
	output, _, err := transform.String(t, value)
	//output = fmt.Sprint(output)
	// log here
	if err != nil {
		return strings.ToUpper(value)
	} else {
		return strings.ToUpper(output)
	}
}

func (field *Field) ToMeta(tableId int32) *Meta {
	// we cannot have error here
	var result = new(Meta)

	// key
	result.id = field.id
	result.refId = tableId
	result.objectType = int8(entitytype.Field)

	// others
	result.dataType = int32(field.fieldType)
	result.name = field.name // max length 30 !! must be validated before
	result.description = field.description
	result.value = field.defaultValue

	// flags
	result.flags = 0
	result.setFieldNotNull(field.notNull)
	result.setFieldCaseSensitive(field.caseSensitive)
	result.setFieldMultilingual(field.multilingual)
	result.setEntityBaseline(field.baseline)
	result.setFieldSize(uint32(field.size))

	result.enabled = field.active
	return result
}

func (field *Field) GetDdl(provider databaseprovider.DatabaseProvider, tableType tabletype.TableType) string {
	datatype := field.getSqlDataType(provider)
	if datatype == unknownFieldDataType {
		return unknownFieldDataType
	}
	return strings.TrimSpace(field.GetPhysicalName(provider) + " " + field.getSqlDataType(provider) + " " +
		field.getSqlConstraint(provider, tableType))
}

// reformat value for records
func (field *Field) GetValue(value string) (string, error) {
	switch field.fieldType {
	case fieldtype.String:
		if field.size > 0 && len(value) > int(field.size) {
			//truncate or error
			return value[0:field.size], nil
		}
		return value, nil
	case fieldtype.Long, fieldtype.Int, fieldtype.Short, fieldtype.Byte:
		if field.isValidInteger(value) == true {
			return value, nil
		}
		break
	case fieldtype.DateTime, fieldtype.ShortDateTime, fieldtype.LongDateTime:
		// must support iso-8601
		t, e := field.getDateTimeIso8601(value)
		if e == nil {
			return field.GetDateTimeString(*t), nil
		}
		return value, e
	default:
		if field.IsValueValid(value) {
			return value, nil
		}
		break
	}
	return invalidValue, errors.New(errorInvalidValueType)
}

func (field *Field) IsValueValid(value string) bool {
	// nullable field?
	if field.notNull == false && value == "" {
		return true
	}
	switch field.fieldType {
	case fieldtype.Long, fieldtype.Int, fieldtype.Short, fieldtype.Byte:
		return field.isValidInteger(value)
	case fieldtype.Double:
		_, err := strconv.ParseFloat(value, 64)
		return err == nil
	case fieldtype.Float:
		_, err := strconv.ParseFloat(value, 32)
		return err == nil
	case fieldtype.String:
		return true
	case fieldtype.DateTime, fieldtype.LongDateTime, fieldtype.ShortDateTime:
		_, err := field.getDateTimeIso8601(value)
		return err == nil
	case fieldtype.Boolean:
		return strings.ToLower(value) == "true" || strings.ToLower(value) == "false"
	}
	return false
}

func (field *Field) Clone() *Field {
	newField := new(Field)
	newField.Init(field.id, field.name, field.description, field.fieldType, uint32(field.size), field.defaultValue, field.baseline,
		field.notNull, field.caseSensitive, field.multilingual, field.active)
	return newField
}

func (field *Field) String() string {
	var fieldTyp = field.fieldType.String()
	if field.fieldType == fieldtype.String && field.size > 0 {
		fieldTyp += fmt.Sprintf("(%d)", field.size)
	}
	return fmt.Sprintf(fieldToStringFormat, field.name, field.description, fieldTyp, field.defaultValue, field.baseline,
		field.notNull, field.caseSensitive, field.active)
}

func (field *Field) Equal(fieldB *Field) bool {
	if field == nil && fieldB == nil {
		return true
	}
	if (field == nil && fieldB != nil) || (field != nil && fieldB == nil) {
		return false
	}
	if field.String() == fieldB.String() {
		return true
	}
	return false
}

func (field *Field) GetDateTimeString(t time.Time) string {
	switch field.fieldType {
	case fieldtype.DateTime:
		return t.UTC().Format(defaultTimeFormat)
	case fieldtype.ShortDateTime:
		return t.Format(defaultShortTimeFormat)
	case fieldtype.LongDateTime:
		return t.Format(time.RFC3339Nano)
	case fieldtype.String:
		return t.UTC().Format(defaultTimeFormat)
	}
	return ""
}

// calculate
func (field *Field) GetParameterValue(value string) interface{} {
	switch field.fieldType {
	case fieldtype.Long, fieldtype.Int, fieldtype.Byte, fieldtype.Short:
		val, _ := strconv.ParseInt(value, 10, 64)
		return val
	case fieldtype.Double, fieldtype.Float:
		val, _ := strconv.ParseFloat(value, 64)
		return val
	case fieldtype.Boolean:
		if value == "true" {
			return true
		}
		return false
	case fieldtype.DateTime:
		val, _ := time.Parse(defaultTimeFormat, value)
		return val
	}
	return value
}

func (field *Field) GetPhysicalName(provider databaseprovider.DatabaseProvider) string {
	return field.getPhysicalName(provider)
}

//******************************
// private methods
//******************************
func (field *Field) getPhysicalName(provider databaseprovider.DatabaseProvider) string {
	switch provider {
	case databaseprovider.PostgreSql:
		var modifyFieldName = false
		if strings.Contains(field.name, "@") == true {
			modifyFieldName = true
		} else {
			if _, ok := sqlKeyWords[strings.ToUpper(field.name)]; ok {
				modifyFieldName = true
			}
		}
		if modifyFieldName == true {
			var sb strings.Builder
			sb.Grow(len(field.name) + 2)
			sb.WriteString(doubleQuotes)
			sb.WriteString(field.name)
			sb.WriteString(doubleQuotes)
			return sb.String()
		}
	}
	return field.name
}

func (field *Field) isValidInteger(value string) bool {
	var sign, size = 0, len(value)
	for _, v := range value {
		if v >= '0' && v <= '9' {
			sign++
			continue
		} else if v == '-' && sign == 0 {
			sign -= size
			continue
		} else {
			return false
		}
	}
	// it's a digit
	switch field.fieldType {
	case fieldtype.Long:
		return field.int64Condition(value, size, sign)
	case fieldtype.Int:
		return field.int32Condition(value, size, sign)
	case fieldtype.Short:
		return field.int16Condition(value, size, sign)
	case fieldtype.Byte:
		return field.int08Condition(value, size, sign)
	}
	return false
}

func (field *Field) int08Condition(value string, size int, sign int) bool {
	return (size > 0 && size < 3) ||
		(size == 3 && value <= maxInt08 && sign > 0) ||
		(size == 3 && sign == -1) ||
		(size == 4 && value <= minInt08 && sign == -1)
}

func (field *Field) int16Condition(value string, size int, sign int) bool {
	return (size > 0 && size < 5) ||
		(size == 5 && value <= maxInt16 && sign > 0) ||
		(size == 5 && sign == -1) ||
		(size == 6 && value <= minInt16 && sign == -1)
}
func (field *Field) int32Condition(value string, size int, sign int) bool {
	return (size > 0 && size < 10) ||
		(size == 10 && value <= maxInt32 && sign > 0) ||
		(size == 10 && sign == -1) ||
		(size == 11 && value <= minInt32 && sign == -1)
}
func (field *Field) int64Condition(value string, size int, sign int) bool {
	return (size > 0 && size < 19) ||
		(size == 19 && value <= maxInt64 && sign > 0) ||
		(size == 19 && sign == -1) ||
		(size == 20 && value <= minInt64 && sign == -1)
}

func (field *Field) getDefaultValue(defaultValue string) string {
	var newDefaultValue = ""
	if field.IsValueValid(defaultValue) == true {
		newDefaultValue = defaultValue
	}
	if newDefaultValue == "" {
		if field.notNull {

			if field.IsNumeric() {
				return defaultNumberValue
			}
			switch field.fieldType {
			case fieldtype.Boolean:
				return defaultBooleanValue
			case fieldtype.DateTime:
				return defaultDateTimeValue
			case fieldtype.ShortDateTime:
				return defaultShortDateTimeValue
			case fieldtype.LongDateTime:
				return defaultLongDateTimeValue
			}
		}
	}
	return newDefaultValue
}

func (field *Field) getDefaultPrimaryKey() *Field {
	switch field.fieldType {
	case fieldtype.Int:
		return defaultPrimaryKeyInt32
	case fieldtype.Long:
		return defaultPrimaryKeyInt64
	case fieldtype.Short:
		return defaultPrimaryKeyInt16
	}
	return defaultPrimaryKeyInt64
}

func (field *Field) getSqlDataType(provider databaseprovider.DatabaseProvider) string {
	var result = unknownFieldDataType
	switch provider {
	case databaseprovider.MySql:
		if val, ok := postgreDataType[field.fieldType]; ok {
			result = val
		}
		break
	case databaseprovider.PostgreSql:
		if val, ok := postgreDataType[field.fieldType]; ok {
			result = val
			if field.fieldType == fieldtype.String {
				if field.size <= 0 || field.size > postgreVarcharMaxSize {
					result = "text"
				} else {
					// long text in postgresql
					result += fmt.Sprintf("(%d)", field.size)
				}
			}
		}
		break
	}
	return result
}

func (field *Field) getSqlConstraint(provider databaseprovider.DatabaseProvider, tableType tabletype.TableType) string {
	// postgresql
	var result = "NULL"
	if tableType != tabletype.Business && field.notNull == true {
		result = "NOT " + result
	}

	if field.fieldType == fieldtype.Byte && provider == databaseprovider.PostgreSql {
		var physicalName = field.getPhysicalName(provider)
		result += fmt.Sprintf(postgreSqlByteConstraint, physicalName, physicalName)
	}
	return result
}

func (field *Field) getDateTimeIso8601(value string) (*time.Time, error) {
	var t time.Time
	var e error

	//short format
	// YYYY-MM-DD
	if len(value) >= 10 {
		//currentTime := time.Now
		// 2019-10-12T14:20:50.52+07:00
		// get currentTimeZone
		t, e = time.Parse(time.RFC3339, value[0:10]+"T00:00:00Z")

		//no information about time zone, take the local one (can be cached)
		if e != nil {
			return nil, e
		}
		if strings.Contains(value, "Z") || strings.Contains(value, "-") || strings.Contains(value, "+") {
			_, offset := time.Now().Zone()
			t = t.Add(time.Second * time.Duration(offset*-1))
		}

		if len(value) >= 19 {
			// where is the timezone?
			// detect timezone
			hour, errH := strconv.Atoi(value[11:13])
			minute, errM := strconv.Atoi(value[14:16])
			second, errS := strconv.Atoi(value[17:19])
			if errH == nil && errM == nil && errS == nil {
				t = t.Add((time.Hour * time.Duration(hour)) + (time.Minute * time.Duration(minute)) +
					(time.Second * time.Duration(second)))
			}
		}
		return &t, e
	} else {
		return nil, errors.New(errorInvalidDateTimeFormat)
	}
}
