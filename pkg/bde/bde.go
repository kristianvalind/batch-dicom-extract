// Package bde contains the recursion, parsing and output logic for batch-dicom-extract
package bde

import (
	"fmt"
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
}

type Parser struct {
	input     ParserInput
	tagList   []*dicomtag.Info
	excelFile *excelize.File
}

func NewParser(pi *ParserInput) (*Parser, error) {
	tagList := []string{}

	// Parse tagList and make sure included tags are valid
	tags := strings.Split(pi.TagList, ",")
	for _, t := range tags {
		tagList = append(tagList, strings.TrimSpace(t))
	}

	tagInfo := []*dicomtag.Info{}
	for _, t := range tagList {
		dicomTag, err := dicomtag.FindByName(t)
		if err != nil {
			return nil, fmt.Errorf("error when parsing tag list: %w", err)
		}
		tagInfo = append(tagInfo, &dicomTag)
	}

	return &Parser{
		input:     *pi,
		tagList:   tagInfo,
		excelFile: excelize.NewFile(),
	}, nil
}

func (p *Parser) Parse() error {

	return nil
}
