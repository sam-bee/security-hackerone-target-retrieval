package csvfiles

import (
	"encoding/csv"
	"fmt"
	"os"
)

type OutputFile struct {
	Path      string
	handle    *os.File
	csvWriter *csv.Writer
}

type OutputDestination interface {
	Open() error
	Close()
	Write([]string) error
}

func (out *OutputFile) Open() error {
	handle, err := os.OpenFile(out.Path, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("error opening file %s: %s", out.Path, err)
	}

	out.handle = handle
	out.csvWriter = csv.NewWriter(out.handle)

	return nil
}

func (out *OutputFile) Write(record []string) error {
	if out.csvWriter == nil {
		return fmt.Errorf("csvWriter not initialised. Did you call OutputFile.Open()?")
	}
	err := out.csvWriter.Write(record)
	if err != nil {
		return fmt.Errorf("error writing to CSV file. Record affected: '%s'", record[0])
	}
	return nil
}

func (out *OutputFile) Close() {
	out.csvWriter.Flush()
	out.handle.Close()
}
