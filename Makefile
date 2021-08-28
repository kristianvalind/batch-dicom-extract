SHELL := /bin/bash
.PHONY: all bde-windows

all: bde-windows 

bde-windows:
	$(eval GIT_TAG := $(shell (git tag)))
	cd cmd/batch-dicom-extract; \
	export GOARCH=amd64 GOOS=windows && go build -o ../../batch-dicom-extract.exe -ldflags="-X 'main.Version=$(GIT_TAG)'"

release: all
	zip batch-dicom-extract-windows.zip batch-dicom-extract.exe