// Package bde contains the recursion, parsing and output logic for batch-dicom-extract
package bde

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/suyashkumar/dicom"
	dicomtag "github.com/suyashkumar/dicom/pkg/tag"
	"github.com/xuri/excelize/v2"
)

type ParserInput struct {
	InputFiles       []string
	RecursiveMode    bool
	OneFilePerSeries bool
	OutputFileName   string
	TagList          string
	StopOnError      bool
	DICOMSuffix      string
}

type Parser struct {
	input         ParserInput
	tagList       []*dicomtag.Info
	fileList      []string
	excelFile     *excelize.File
	sheetName     string
	visitedSeries map[string]bool
}

func NewParser(pi *ParserInput) (*Parser, error) {
	tagStrings := []string{}

	// Parse tagList and make sure included tags are valid
	tags := strings.Split(pi.TagList, ",")
	for _, t := range tags {
		tagStrings = append(tagStrings, strings.TrimSpace(t))
	}

	tagList := []*dicomtag.Info{}
	for _, t := range tagStrings {
		dicomTag, err := dicomtag.FindByName(t)
		if err != nil {
			return nil, fmt.Errorf("error when parsing tag list: %w", err)
		}
		tagList = append(tagList, &dicomTag)
	}

	excelFile := excelize.NewFile()

	sheetName := "Results"

	// Write headers to Excel file
	excelFile.SetSheetName("Sheet1", sheetName)
	excelFile.SetSheetRow(sheetName, "A1", &tagStrings)

	return &Parser{
		input:         *pi,
		tagList:       tagList,
		excelFile:     excelFile,
		sheetName:     sheetName,
		fileList:      make([]string, 0),
		visitedSeries: make(map[string]bool),
	}, nil
}

func (p *Parser) Parse() error {

	// Collect all files to visit
	for _, fileName := range p.input.InputFiles {
		fi, err := os.Stat(fileName)
		if err != nil {
			if p.input.StopOnError {
				return fmt.Errorf("error when getting file info: %w", err)
			} else {
				log.Printf("error when getting file info: %v", err.Error())
			}
		}

		if fi.IsDir() {
			// File is directory

			// Skip directories when not in recursive mode
			if !p.input.RecursiveMode {
				continue
			} else {
				err := filepath.Walk(fileName, func(path string, info fs.FileInfo, err error) error {
					if !info.IsDir() {
						// Skip hidden files and directories
						if !strings.HasPrefix(path, ".") {
							// Check suffix matching
							if strings.HasSuffix(strings.ToLower(fileName), strings.ToLower(p.input.DICOMSuffix)) {
								p.fileList = append(p.fileList, path)
							}
						}
					}
					return nil
				})
				if err != nil {
					if p.input.StopOnError {
						return fmt.Errorf("error when recursing into directory: %w", err)
					} else {
						log.Printf("error when recursing into directory: %v", err.Error())
					}
				}
			}

		} else {
			// File is not directory

			// Check suffix matching
			if strings.HasSuffix(strings.ToLower(fileName), strings.ToLower(p.input.DICOMSuffix)) {
				p.fileList = append(p.fileList, fileName)
			}
		}
	}

	// Visit files
	for _, fileName := range p.fileList {
		err := p.parseDicom(fileName)
		if err != nil {
			if p.input.StopOnError {
				return fmt.Errorf("error when parsing dicom file %v: %w", fileName, err)
			} else {
				log.Printf("error when parsing dicom file %v: %v", fileName, err)
			}
		}
	}

	// Create and write output file
	f, err := os.Create(p.input.OutputFileName)
	if err != nil {
		return fmt.Errorf("error in creating output file: %w", err)
	}
	defer f.Close()

	err = p.excelFile.Write(f)
	if err != nil {
		return fmt.Errorf("error in writing output file: %w", err)
	}

	return nil
}

func (p *Parser) parseDicom(fileName string) error {
	dataset, err := dicom.ParseFile(fileName, nil)
	if err != nil {
		return fmt.Errorf("error parsing dicom file: %w", err)
	}

	// Check if series has been visited
	siUIDElement, err := dataset.FindElementByTag(dicomtag.SeriesInstanceUID)
	if err != nil {
		return fmt.Errorf("error getting series instance uid: %w", err)
	}

	siUIDSlice, ok := siUIDElement.Value.GetValue().([]string)
	if !ok {
		return fmt.Errorf("series instance uid is not string slice")
	}

	siUID := siUIDSlice[0]

	// If we have marked series as visited and only check once, we stop here
	if p.visitedSeries[siUID] && p.input.OneFilePerSeries {
		return nil
	}

	p.visitedSeries[siUID] = true

	err = p.extractAndWrite(&dataset)
	if err != nil {
		return fmt.Errorf("error in extracting dicom tag data: %w", err)
	}

	return nil
}
