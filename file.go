package main

import (
	"bytes"
	"compress/gzip"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

func save(file, passphrase string, sites []Site) error {
	// If the file doesn't exist then create it
	if _, err := os.Stat(file); os.IsNotExist(err) {
		_, err = os.Create(file)
		if err != nil {
			return err
		}
	}

	// Marshal the JSON
	b, err := json.Marshal(sites)
	if err != nil {
		return err
	}

	// Compress the contents
	var buffer bytes.Buffer
	gzip := gzip.NewWriter(&buffer)
	if err != nil {
		return err
	}
	gzip.Write(b)
	gzip.Close()

	// Write to the file
	fi, err := os.OpenFile(file, os.O_WRONLY, 0666)
	if err != nil {
		return err
	}
	_, err = fi.Write(buffer.Bytes())
	if err != nil {
		return err
	}
	fi.Close()

	return nil
}

// Read the password book
func read(file, passphrase string) ([]Site, error) {
	// If the file doesn't exist yet no worries
	if _, err := os.Stat(file); os.IsNotExist(err) {
		return []Site{}, nil
	}

	// Bring in the compressed data
	fi, err := os.Open(file)
	if err != nil {
		return nil, err
	}

	// Decompress the file contents
	gzip, err := gzip.NewReader(fi)
	if err != nil {
		return nil, err
	}
	decompressed, err := ioutil.ReadAll(gzip)
	gzip.Close()

	// Unmarshal the JSON information
	var sites []Site
	err = json.Unmarshal(decompressed, &sites)
	if err != nil {
		return nil, err
	}
	fi.Close()

	return sites, nil
}

// Get the book name
func getBookname(profile string) string {
	hash := md5.New()
	hash.Write([]byte(profile))
	return fmt.Sprintf("%x", string(hash.Sum(nil)))
}
