package data

import (
	"errors"
	"fmt"
	"ring/schema"
)

type Import struct {
	insertCount   int
	updateCount   int
	jobId         int64
	currentSchema *schema.Schema
}

const (
	errorUnknownImpSchema string = "Schema name %s does not exist."
	errorNotDefinedSchema string = "Schema is not defined call!"
)

//******************************
// getters / setters
//******************************
func (impFile *Import) GetJobId() int64 {
	return impFile.jobId
}

//******************************
// public methods
//******************************
func (impFile *Import) SetSchema(schemaName string) error {
	impFile.currentSchema = schema.GetSchemaByName(schemaName)
	if impFile.currentSchema == nil {
		return errors.New(fmt.Sprintf(errorUnknownImpSchema, schemaName))
	}
	return nil
}

func (impFile *Import) ParseFile(file string) error {
	if impFile.currentSchema == nil {
		return errors.New(errorNotDefinedSchema)
	}
	var err error
	impFile.insertCount = 0
	impFile.updateCount = 0
	doc := new(document)
	err = doc.Load(file)
	if err != nil {
		return err
	}

	return nil
}

//******************************
// private methods
//******************************
