package data

import (
	"archive/zip"
	"errors"
	"fmt"
	"os"
	"ring/data/documenttype"
	"strings"
	"time"
)

const (
	errorUnknownDocumentType string = "Unknown document type."
	xmlSuffix                string = ".xml"
)

type document struct {
	filePath      string
	creator       string
	creationtTime *time.Time
	updateTime    *time.Time
	sheets        []*sheetItem
	documentType  documenttype.DocumentType
}

//******************************
// getters / setters
//******************************
func (doc *document) GetCreationTime() *time.Time {
	return doc.creationtTime
}
func (doc *document) GetUpdateTime() *time.Time {
	return doc.updateTime
}
func (doc *document) GetDocumentType() documenttype.DocumentType {
	return doc.documentType
}

//******************************
// public methods
//******************************
func (doc *document) Load(filePath string) error {
	var err error
	doc.filePath = filePath
	doc.sheets = nil
	doc.documentType = documenttype.GetDocumentType(filePath)
	err = doc.exists()
	if err != nil {
		return err
	}
	if doc.documentType == documenttype.Undefined {
		return errors.New(errorUnknownDocumentType)
	}
	switch doc.documentType {
	case documenttype.Xlsx:
		err = doc.loadXslx()
		break
	case documenttype.Json:
		break
	}

	return nil
}

//******************************
// private methods
//******************************
func (doc *document) loadXslx() error {
	archive, err := zip.OpenReader(doc.filePath)
	if err != nil {
		return err
	}
	defer archive.Close()

	for _, f := range archive.File {
		filePath := f.Name
		if strings.HasSuffix(f.Name, xmlSuffix) {
			fmt.Println("unzipping file ", filePath)
		}
		if f.FileInfo().IsDir() {
			continue
		}
		fileInArchive, err := f.Open()
		if err != nil {
			panic(err)
		}

		fileInArchive.Close()
	}
	return nil
}

func (doc *document) exists() error {
	_, err := os.Stat(doc.filePath)
	if err == nil {
		return nil
	}
	return err
}
