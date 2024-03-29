package sqlfmt

import (
	"ring/schema/databaseprovider"
	"strings"
	"unicode"
)

const (
	doubleQuotes   string = "\""
	mysqlQuotes    string = "`"
	snakeSeparator rune   = '_'
	kebabSeparator rune   = '-'
	spaceSeparator rune   = ' '
)

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

func PadLeft(input string, padString string, repeat int) string {
	var count = repeat - len(input)
	if count > 0 {
		return strings.Repeat(padString, count) + input
	} else {
		return input
	}
}

func FormatEntityName(provider databaseprovider.DatabaseProvider, entityName string) string {
	var fieldSeparator = ""

	var modifyFieldName = false
	if strings.Contains(entityName, "@") == true {
		modifyFieldName = true
	} else {
		if _, ok := sqlKeyWords[strings.ToUpper(entityName)]; ok {
			modifyFieldName = true
		}
	}

	switch provider {
	case databaseprovider.PostgreSql:
		fieldSeparator = doubleQuotes
		break
	case databaseprovider.MySql:
		fieldSeparator = mysqlQuotes
		break
	}

	if modifyFieldName == true {
		var sb strings.Builder
		sb.Grow(len(entityName) + 2)
		sb.WriteString(fieldSeparator)
		sb.WriteString(entityName)
		sb.WriteString(fieldSeparator)
		return sb.String()
	}
	return entityName
}

func UnFormatEntityName(provider databaseprovider.DatabaseProvider, entityName string) string {

	// fast version
	switch provider {
	case databaseprovider.PostgreSql:
		return strings.ReplaceAll(entityName, doubleQuotes, "")

	case databaseprovider.MySql:
		return strings.ReplaceAll(entityName, mysqlQuotes, "")
	}

	return entityName
}

func ToSnakeCase(name string) string {
	var result strings.Builder
	var letterBefore = false
	var prevRune rune

	result.Grow(len(name) + 10)
	prevRune = snakeSeparator
	name = strings.ReplaceAll(name, " ", string(snakeSeparator))

	for _, chr := range name {
		if letterBefore == true {
			if chr == snakeSeparator && prevRune != snakeSeparator {
				result.WriteRune(snakeSeparator)
			} else {
				splitSnake(&result, chr, prevRune, &letterBefore)
			}
		} else {
			if unicode.IsUpper(chr) || unicode.IsLower(chr) {
				result.WriteRune(unicode.ToLower(chr))
				letterBefore = true
			}
		}
		prevRune = chr
	}
	return result.String()
}

func ToCamelCase(name string) string {
	var result strings.Builder
	var letterBefore = false
	var toUpperCaseMode = false

	result.Grow(len(name))

	for _, chr := range name {
		if letterBefore == false && (chr == snakeSeparator || chr == kebabSeparator || chr == spaceSeparator) {
			continue
		}
		if chr == snakeSeparator || chr == kebabSeparator || chr == spaceSeparator {
			toUpperCaseMode = true
			continue
		}
		if toUpperCaseMode == true {
			result.WriteRune(unicode.ToUpper(chr))
		} else {
			result.WriteRune(unicode.ToLower(chr))
		}
		toUpperCaseMode = false
		letterBefore = true
	}

	return result.String()
}

func ToPascalCase(name string) string {
	var result strings.Builder
	nameCamel := ToCamelCase(name)

	if len(nameCamel) > 0 {
		result.Grow(len(nameCamel))
		result.WriteRune(unicode.ToUpper(rune(nameCamel[0])))
		result.WriteString(nameCamel[1:])
	}

	return result.String()
}

func ToKebabCase(name string) string {
	return ""
}

func splitSnake(result *strings.Builder, currentChr rune, previousChr rune, letterBefore *bool) {
	if unicode.IsUpper(currentChr) {
		if unicode.IsUpper(previousChr) == false && previousChr != snakeSeparator {
			result.WriteRune(snakeSeparator)
		}
		result.WriteRune(unicode.ToLower(currentChr))
		*letterBefore = true
	}
	if unicode.IsLower(currentChr) {
		result.WriteRune(unicode.ToLower(currentChr))
		*letterBefore = true
	}
}
