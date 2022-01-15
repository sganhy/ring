package data

import (
	"archive/zip"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"os"
	"ring/data/documenttype"
	"strconv"
	"strings"
	"time"
)

const (
	errorUnknownDocumentType string = "Unknown document type."
	xmlSufix                 string = ".xml"
	xmlSheetElement          string = "sheet"
	xmlWorkBookSufix         string = "/workbook.xml"
	xmlShareddStringsSufix   string = "/sharedStrings.xml"
	xmlSheetIdElement        string = "sheetId"
	xmlSheetNameElement      string = "name"
	xmlSharedStringsElement  string = "t"
)

type document struct {
	filePath     string
	creator      string
	creationTime *time.Time
	updateTime   *time.Time
	sheets       []*sheetItem
	documentType documenttype.DocumentType
}

//******************************
// getters / setters
//******************************
func (doc *document) GetCreationTime() *time.Time {
	return doc.creationTime
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
	doc.sheets = make([]*sheetItem, 0, 2)
	var sharedStrings []string

	if err != nil {
		return err
	}
	defer archive.Close()

	// pass 1 : load sheets & shared strings
	for _, f := range archive.File {
		if f.FileInfo().IsDir() || strings.HasSuffix(f.Name, xmlSufix) == false {
			continue
		}

		if strings.HasSuffix(f.Name, xmlWorkBookSufix) == true ||
			strings.HasSuffix(f.Name, xmlShareddStringsSufix) == true {
			fileInArchive, err := f.Open()
			if err != nil {
				panic(err)
			}
			doc.loadXmlSheets(fileInArchive, f.Name)
			doc.loadXmlSharedString(fileInArchive, f, &sharedStrings)
			fileInArchive.Close()
		}
	}
	return nil
}

func (doc *document) loadXmlSharedString(reader io.ReadCloser, file *zip.File, result *[]string) error {
	if !strings.HasSuffix(file.Name, xmlShareddStringsSufix) {
		return nil
	}
	d := xml.NewDecoder(reader)
	count := 0
	for {
		tok, err := d.Token()
		if tok == nil || err == io.EOF {
			break
		}
		switch ty := tok.(type) {
		case xml.StartElement:
			if xmlSharedStringsElement == ty.Name.Local {
				count++
			}
			break
		}
	}
	fmt.Println(count)
	// allow it once
	sharedStrings := make([]string, count, count)
	err := doc.loadXmlSharedStringData(file, result)
	result = &sharedStrings
	return err
}

func (doc *document) loadXmlSharedStringData(file *zip.File, result *[]string) error {
	// reopen file to rewind reader

	reader, err := file.Open()
	if err != nil {
		return err
	}
	d := xml.NewDecoder(reader)
	defer reader.Close()

	i := 0
	count := len(*result)
	loadData := false

	for {
		tok, err := d.Token()
		if tok == nil || err == io.EOF {
			break
		}
		switch ty := tok.(type) {
		case xml.StartElement:
			if xmlSharedStringsElement == ty.Name.Local {
				loadData = true
			} else {
				loadData = false
			}
			break
		case xml.CharData:
			if loadData && i < count {
				(*result)[i] = string([]byte(xml.CharData(ty)))
				i++
			}
		}
	}
	return nil
}

func (doc *document) loadXmlSheets(reader io.ReadCloser, filePath string) error {
	if !strings.HasSuffix(filePath, xmlWorkBookSufix) {
		return nil
	}
	d := xml.NewDecoder(reader)
	for {
		tok, err := d.Token()
		if tok == nil || err == io.EOF {
			// EOF means we're done.
			break
		} else if err != nil {
			// log here
			return err
		}
		switch ty := tok.(type) {
		case xml.StartElement:
			if xmlSheetElement == ty.Name.Local {
				sheetItm := new(sheetItem)
				name := doc.getXmlAttribute(&ty.Attr, xmlSheetNameElement)
				id := int(doc.getXmlIntAttribute(&ty.Attr, xmlSheetIdElement))
				sheetItm.Init(id, name)
				doc.sheets = append(doc.sheets, sheetItm)
			}
			break
		}
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

func (doc *document) getXmlAttribute(attributes *[]xml.Attr, attributeName string) string {
	if attributes != nil {
		count := len(*attributes)
		for i := 0; i < count; i++ {
			var attribute = (*attributes)[i]
			if strings.EqualFold(attribute.Name.Local, attributeName) {
				return attribute.Value
			}
		}
	}
	return ""
}

func (doc *document) getXmlIntAttribute(attributes *[]xml.Attr, attributeName string) int64 {
	value := doc.getXmlAttribute(attributes, attributeName)
	result, err := strconv.ParseInt(value, 10, 64)
	if err == nil {
		return result
	}
	return -1
}
