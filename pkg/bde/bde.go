// Package bde contains the recursion, parsing and output logic for batch-dicom-extract
package bde

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"

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
}

type Parser struct {
	input         ParserInput
	tagList       []*dicomtag.Info
	fileList      []string
	excelFile     *excelize.File
	visitedSeries []string
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

	return &Parser{
		input:         *pi,
		tagList:       tagList,
		excelFile:     excelize.NewFile(),
		fileList:      make([]string, 0),
		visitedSeries: make([]string, 0),
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
							p.fileList = append(p.fileList, path)
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
			p.fileList = append(p.fileList, fileName)
		}
	}

	return nil
}
