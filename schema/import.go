package schema

import (
	"errors"
	"ring/schema/databaseprovider"
	"ring/schema/documenttype"
	"ring/schema/sourcetype"
	"runtime"
	"strings"
)

type Import struct {
	id                int
	fileName          string
	schemaId          int32
	schemaName        string
	source            sourcetype.SourceType
	initialized       bool
	loaded            bool
	jobId             int64
	errorCount        int32
	metaList          []*meta
	newSchema         *Schema
	logger            *log
	provider          databaseprovider.DatabaseProvider
	metaChangeCount   int64
	metaIdChangeCount int64
	doc               *document
}

const (
	errorImportFileNotInitialized string = "Import File object is not initialized."
	errorImportInvalidId          string = "Invalid attribute {id}"
	errorImportFieldSize          string = "Invalid field size"
	importTableTag1               string = "object"
	importTableTag2               string = "table"
	importFieldTag                string = "field"
	importTablespaceTag           string = "tablespace"
	importSchemaTag               string = "schema"
	importDescriptionTag          string = "description"
	importRelationTag             string = "relation"
	importIndexTag                string = "index"
	importIndexFieldTag           string = "index_field"
	importMinId                   int64  = -2147483648
	importMaxId                   int64  = 2147483647
	baseErrorId                   int32  = 11
	maxXmlFileSize                int64  = 33554432 // 32 mega byte
)

var (
	currentSchemaImportId int = 101
)

func (importFile *Import) Init(source sourcetype.SourceType, fileName string) {
	currentSchemaImportId++
	importFile.id = currentSchemaImportId
	importFile.fileName = fileName
	importFile.initialized = true
	importFile.source = source
	importFile.logger = new(log)
	importFile.logger.Init(schemaUndefined, 0, false)
	importFile.loaded = false
}

//******************************
// getters and setters
//******************************
func (importFile *Import) GetJobId() int64 {
	return importFile.jobId
}

func (importFile *Import) GetId() int {
	return importFile.id
}

func (importFile *Import) GetFile() string {
	return importFile.fileName
}

func (importFile *Import) GetSchemaName() string {
	return importFile.schemaName
}

func (importFile *Import) GetSchemaId() int32 {
	return importFile.schemaId
}

func (importFile *Import) GetSourceType() sourcetype.SourceType {
	return importFile.source
}

func (importFile *Import) GetSchema() *Schema {
	return importFile.newSchema
}

func (importFile *Import) GetDatabaseProvider() databaseprovider.DatabaseProvider {
	return importFile.provider
}

//******************************
// public methods
//******************************
func (importFile *Import) Load() error {
	var err error
	var metaSchema = GetSchemaByName(metaSchemaName)
	var fileName = strings.ReplaceAll(importFile.fileName, "\\", "/")

	importFile.provider = getDefaultDbProvider()
	importFile.metaList = nil
	importFile.errorCount = 0
	importFile.jobId = metaSchema.getJobIdNextValue()
	importFile.logInfo("Load schema", "import_file: "+fileName)
	importFile.doc = new(document)

	//importFile.doc.Init(

	if importFile.initialized == false {
		return errors.New(errorImportFileNotInitialized)
	}

	if importFile.source == sourcetype.XmlDocument {

		importFile.doc.Init(importFile.fileName, documenttype.Xml, importFile.jobId, importFile.provider)
		importFile.doc.loadSchemaInfo()

		var isValid = false
		isValid, err = importFile.doc.isXmlValid()
		if isValid {
			importFile.doc.loadXml()
			valid := new(validator)
			valid.Init()
			valid.ValidateImport(importFile)
		}

	}
	importFile.loaded = true
	return err
}

func (importFile *Import) Upgrade() {
	importFile.metaChangeCount = 0
	importFile.metaIdChangeCount = 0

	if importFile.IsValid() == true {
		var err error
		err = importFile.saveMetaList()
		if err != nil {
			importFile.logError(err)
			return
		}

		err = importFile.saveMetaIdList()
		if err != nil {
			importFile.logError(err)
			return
		}
		//TODO rollback !!
		// no error
		importFile.loaded = false
		importFile.metaList = nil
		importFile.newSchema = getSchemaById(importFile.schemaId, true, true)

		// database.go
		upgradeSchema(importFile.jobId, importFile.newSchema)

		// at the end perform vacuum
		go importFile.analyzeMetaTables()

		// release ressources
		runtime.GC()
	}
}

func (importFile *Import) IsValid() bool {
	return importFile.errorCount == 0 && importFile.loaded
}

//******************************
// private methods
//******************************

//go:noinline
func (importFile *Import) logInfo(message string, description string) {
	importFile.logger.writePartialLog(9, levelInfo, importFile.jobId, message, description)
}

//go:noinline
func (importFile *Import) logError(err error) {
	importFile.logger.writePartialLog(baseErrorId+importFile.errorCount, levelError, importFile.jobId, err)
	importFile.errorCount++
}

//go:noinline
func (importFile *Import) logErrorStr(id int32, message string, description string) {
	importFile.logger.writePartialLog(id, levelError, importFile.jobId, message, description)
	importFile.errorCount++
}

func (importFile *Import) saveMetaList() error {
	metaData := new(meta)
	var count int64
	var err error

	count, err = metaData.saveMetaList(importFile.schemaId, importFile.metaList)
	importFile.metaChangeCount = count

	return err
}

func (importFile *Import) saveMetaIdList() error {
	var count int64
	var err error

	metaId := new(metaId)
	count, err = metaId.saveMetaIdList(importFile.schemaId, importFile.metaList)
	importFile.metaIdChangeCount = count
	return err
}

func (importFile *Import) analyzeMetaTables() {
	var table *Table
	var metaSchema = GetSchemaByName(metaSchemaName)

	if importFile.metaIdChangeCount > 0 {
		table = metaSchema.GetTableByName(metaIdTableName)
		table.Vacuum(importFile.jobId, false)
		table.Analyze(importFile.jobId)
	}
	if importFile.metaChangeCount > 10 {
		table = metaSchema.GetTableByName(metaTableName)
		table.Vacuum(importFile.jobId, true)
		table.Analyze(importFile.jobId)
	}
}
