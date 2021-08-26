package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/kristianvalind/batch-dicom-extract/pkg/bde"
)

var Version string = "v0.0.1"
var outputFileName, tagList string
var oneFilePerSeries, recurseIntoDirectories, stopOnError bool

func init() {
	flag.Usage = func() {
		fmt.Printf("batch-dicom-extract %v. by Kristian Valind.\nExtracts tag data from any number of DICOM files to an XLSX-file.\n", Version)
		fmt.Fprintf(flag.CommandLine.Output(), "Usage:\n%s [flags] files\n%s [flags] -r directories\n", os.Args[0], os.Args[0])
		flag.PrintDefaults()
	}
	flag.StringVar(&outputFileName, "o", "dicom-batch.xlsx", "name of `file` to write results to.")
	flag.BoolVar(&oneFilePerSeries, "1", false, "Consider only 1 file per series.")
	flag.BoolVar(&recurseIntoDirectories, "r", false, "Recurse into directories.")
	flag.BoolVar(&stopOnError, "s", false, "Stop parsing when encountering an error, rather than skipping the file.")
	flag.StringVar(&tagList, "t", "PatientID, StudyDescription", "Comma separated list of DICOM tag keywords. Spaces are stripped.")

	flag.Parse()
}

func main() {
	p, err := bde.NewParser(&bde.ParserInput{
		InputFiles:       flag.Args(),
		RecursiveMode:    recurseIntoDirectories,
		OneFilePerSeries: oneFilePerSeries,
		OutputFileName:   outputFileName,
		TagList:          tagList,
		StopOnError:      stopOnError,
	})
	if err != nil {
		log.Fatal(err)
	}

	err = p.Parse()
	if err != nil {
		log.Fatal(err)
	}
}
