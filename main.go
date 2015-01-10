package main

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

type Site struct {
	Host                      string `json:host`
	MinimumLength             int    `json:minimumLength`
	MaximumLength             int    `json:maximumLength`
	SpecialCharacters         string `json:specialCharacters`
	NumberOfSpecialCharacters int    `json:numberOfSpecialCharacters`
	NumberOfUpperCase         int    `json:numberOfUpperCase`
	NumberOfDigits            int    `json:numberOfDigits`
}

func main() {
	http.HandleFunc("/api", func(w http.ResponseWriter, r *http.Request) {
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		path := strings.TrimPrefix(r.URL.Path, "/")
		if path == "" {
			// Index
			page, err := ioutil.ReadFile("index.html")
			if err != nil {
				http.NotFound(w, r)
				return
			}

			fmt.Fprintf(w, string(page))
		} else {
			// A passphrase has been entered
		}
	})

	log.Fatal(http.ListenAndServe("localhost:8080", nil))

	//execDir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	//if err != nil {
	//	log.Fatal(err)
	//}

	//file := filepath.Join(execDir, "enigma.safe")
	//sites, err := Read(file)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//log.Println(sites)

	//err = Save(file, sites)
	//if err != nil {
	//	log.Fatal(err)
	//}
}

func Save(file string, sites []Site) error {
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
	buffer := bytes.NewBuffer(b)
	gzipWriter := gzip.NewWriter(buffer)
	if err != nil {
		return err
	}
	defer gzipWriter.Close()

	// Write to the file
	fi, err := os.OpenFile(file, os.O_WRONLY, 0666)
	if err != nil {
		return err
	}
	defer fi.Close()
	_, err = fi.Write(buffer.Bytes())
	if err != nil {
		return err
	}

	return nil
}

func Read(file string) ([]Site, error) {
	// If the file doesn't exist yet no worries
	if _, err := os.Stat(file); os.IsNotExist(err) {
		return []Site{}, nil
	}

	// Bring in the compressed data
	log.Println("Reading compressed file")
	compressed, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}

	// Decompress the file contents
	log.Println("Decompressing")
	buffer := bytes.NewBuffer(compressed)
	gzipReader, err := gzip.NewReader(buffer)
	if err != nil {
		return nil, err
	}
	defer gzipReader.Close()

	// Unmarshal the JSON information
	log.Println("Unmarshal")
	var sites []Site
	err = json.Unmarshal(buffer.Bytes(), &sites)
	if err != nil {
		return nil, err
	}

	return sites, nil
}
