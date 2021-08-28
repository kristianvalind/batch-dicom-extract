# batch-dicom-extract
A small utility for extracting DICOM tag data to an Excel file.

```
Usage:
batch-dicom-extract [flags] files
batch-dicom-extract [flags] -r directories
  -1    Consider only 1 file per series.
  -e    Stop parsing when encountering an error, rather than skipping the file.
  -o file
        name of file to write results to. (default "dicom-batch.xlsx")
  -r    Recurse into directories.
  -s suffix
        File name suffix of files to consider. (default ".dcm")
  -t list
        Comma separated list of DICOM tag keywords. Spaces are stripped. (default "PatientID, StudyDescription")
```