package main

import (
	"encoding/gob"
	"os"
	"path/filepath"
)

// WriteDataToFile - writes data to local db
func WriteDataToFile(data map[MergedTableMapKey][]IfMergedTable, path string) error {

	var rerr, err error

	// create local db
	f, err := os.Create(filepath.Clean(path))
	if err != nil {
		return err
	}

	// close db on exit
	defer func(f *os.File) {
		err = f.Close()
		if err != nil {
			rerr = err
		}
	}(f)

	// prepare encoder for input data
	enc := gob.NewEncoder(f)

	// encode and write input data to local db
	err = enc.Encode(data)
	if err != nil {
		return err
	}

	return rerr
}

// ReadDataFromFile - reads data to local db
func ReadDataFromFile(path string) (map[MergedTableMapKey][]IfMergedTable, error) {

	var rerr, err error

	// declare output variable
	data := make(map[MergedTableMapKey][]IfMergedTable)

	// open local db
	f, err := os.Open(filepath.Clean(path))
	if err != nil {
		return nil, err
	}

	// close db on exit
	defer func(f *os.File) {
		err = f.Close()
		if err != nil {
			rerr = err
		}
	}(f)

	// prepare decoder for data read from local db
	dec := gob.NewDecoder(f)

	// decode local db data to internal represintation
	err = dec.Decode(&data)
	if err != nil {
		return nil, err
	}

	return data, rerr
}
