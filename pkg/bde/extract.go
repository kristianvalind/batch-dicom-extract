package bde

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/suyashkumar/dicom"
)

func (p *Parser) extractAndWrite(dataset *dicom.Dataset) error {
	results := make([]interface{}, 0)

	for _, tag := range p.tagList {
		elem, err := dataset.FindElementByTagNested(tag.Tag)
		if err != nil {
			if errors.Is(err, dicom.ErrorElementNotFound) {
				results = append(results, "")
			} else {
				return fmt.Errorf("error in finding element %v: %w", tag.Name, err)
			}
		} else {
			value := elem.Value.GetValue()

			// If value is a slice
			if reflect.TypeOf(value).Kind() == reflect.Slice {
				sliceLength := reflect.ValueOf(value).Len()

				// If slice contains only one value, write it out
				if sliceLength == 1 {
					sliceValue := reflect.ValueOf(value).Index(0)
					results = append(results, sliceValue)
				} else { // Otherwise write the entire slice
					results = append(results, value)
				}
			} else {
				// Non-slice values are written directly
				results = append(results, value)
			}
		}
	}

	err := p.excelFile.SetSheetRow(p.sheetName, fmt.Sprintf("A%v", p.rowCounter), &results)
	if err != nil {
		return fmt.Errorf("error in writing result row: %w", err)
	}

	p.rowCounter++
	return nil
}
