package main

import (
	"bytes"
	"compress/gzip"
	"crypto/md5"
	"crypto/sha512"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
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
	Revision                  int    `json:revision`
}

func main() {
	http.HandleFunc("/api/generate", func(w http.ResponseWriter, r *http.Request) {
		profile := r.FormValue("profile")
		passphrase := r.FormValue("p")
		host := r.FormValue("host")
		minimumLength, _ := strconv.Atoi(r.FormValue("minimumLength"))
		maximumLength, _ := strconv.Atoi(r.FormValue("maximumLength"))
		minimumDigits, _ := strconv.Atoi(r.FormValue("minimumDigits"))
		minimumUppercase, _ := strconv.Atoi(r.FormValue("minimumUppercase"))
		minimumSpecialCharacters, _ := strconv.Atoi(r.FormValue("minimumSpecialCharacters"))
		specialCharacters := r.FormValue("specialCharacters")

		if profile == "" || passphrase == "" || host == "" {
			http.Error(w, "Missing credentials", http.StatusUnauthorized)
			return
		}

		site := Site{
			Host:                      host,
			MinimumLength:             minimumLength,
			MaximumLength:             maximumLength,
			SpecialCharacters:         specialCharacters,
			NumberOfSpecialCharacters: minimumSpecialCharacters,
			NumberOfDigits:            minimumDigits,
			NumberOfUpperCase:         minimumUppercase,
			Revision:                  0,
		}

		book := getBookname(profile)
		sites, err := Read(book, passphrase)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		sites = append(sites, site)
		err = Save(book, passphrase, sites)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})
	http.HandleFunc("/api/update", func(w http.ResponseWriter, r *http.Request) {
		profile := r.FormValue("profile")
		passphrase := r.FormValue("p")
		newPassphrase := r.FormValue("newPassphrase")
		confirmPassphrase := r.FormValue("confirmPassphrase")
		cmd := r.FormValue("cmd")

		if profile == "" || passphrase == "" || newPassphrase == "" || confirmPassphrase == "" || cmd == "" {
			http.Error(w, "Missing credentials", http.StatusUnauthorized)
			return
		}

		book := getBookname(profile)

		if cmd == "delete" {
			err := os.Remove(book)
			if err != nil {
				// Return an error
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		} else if cmd == "update" {
			if newPassphrase != confirmPassphrase {
			}
			sites, err := Read(book, passphrase)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			err = Save(book, newPassphrase, sites)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}
	})
	http.HandleFunc("/api/refresh", func(w http.ResponseWriter, r *http.Request) {
		profile := r.FormValue("profile")
		passphrase := r.FormValue("p")
		host := r.FormValue("host")

		if profile == "" || passphrase == "" || host == "" {
			http.Error(w, "Missing credentials", http.StatusUnauthorized)
			return
		}

		// Update the revision number and generate a new password
		book := getBookname(profile)
		sites, err := Read(book, passphrase)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		for _, site := range sites {
			if site.Host == host {
				site.Revision++
				break
			}
		}
		err = Save(book, passphrase, sites)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})
	http.HandleFunc("/api/remove", func(w http.ResponseWriter, r *http.Request) {
		profile := r.FormValue("profile")
		passphrase := r.FormValue("p")
		host := r.FormValue("host")

		if profile == "" || passphrase == "" || host == "host" {
			http.Error(w, "Missing credentials", http.StatusUnauthorized)
			return
		}

		// Remove the site from our book and save it
		book := getBookname(profile)
		sites, err := Read(book, passphrase)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		for i, site := range sites {
			if site.Host == host {
				sites = append(sites[:i], sites[i+1:]...)
				break
			}
		}
		err = Save(book, passphrase, sites)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		profile := r.FormValue("profile")
		passphrase := r.FormValue("p")
		//path := strings.TrimPrefix(r.URL.Path, "/")

		if profile == "" || passphrase == "" {
			// Index
			page, err := ioutil.ReadFile("index.html")
			if err != nil {
				http.NotFound(w, r)
				return
			}

			fmt.Fprintf(w, string(page))
		} else {
			// A passphrase has been entered
			page, err := ioutil.ReadFile("book.html")
			if err != nil {
				http.NotFound(w, r)
				return
			}

			fmt.Fprintf(w, string(page))
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

func Save(file, passphrase string, sites []Site) error {
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

// Read the password book
func Read(file, passphrase string) ([]Site, error) {
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

// Get the book name
func getBookname(profile string) string {
	hash := md5.New()
	hash.Write([]byte(profile))
	return fmt.Sprintf("%s", string(hash.Sum(nil)))
}

// Encrypt the password book
func encrypt(clearText, profile, passphrase string) ([]byte, error) {
	return nil, nil
}

// Decrypt the password book
func decrypt(encryptedText, profile, passphrase string) ([]byte, error) {
	return nil, nil
}

// Generate the passphrase
func generatePassphrase(profile, passphrase string, settings Site) ([]byte, error) {
	clearText := fmt.Sprintf(
		"%s-%s-%s-%s",
		strings.ToLower(profile),
		strings.ToLower(passphrase),
		strings.ToLower(settings.Host),
		settings.Revision)

	sha := sha512.New()
	sha.Write([]byte(clearText))
	hash := sha.Sum(nil)
	hash = []byte(fmt.Sprintf("%x", hash))

	// Apply site criteria
	applySiteSettings(hash, settings)

	// If there is a maximum length truncate the hash
	if settings.MaximumLength > -1 {
		hash = hash[:settings.MaximumLength]
	}

	// Ensure the length is adequate
	if !validateLength(hash, settings.MinimumLength, settings.MaximumLength) {
		log.Println("Does not meed the length requirements")
	}

	return hash, nil
}

// Apply site settings to the hashed value
func applySiteSettings(source []byte, settings Site) []byte {
	if !containsUppercase(source, settings.NumberOfUpperCase) {
		i := 0
		r := regexp.MustCompile(`[a-z]+`)

		var matches [][]int
		if matches = r.FindAllIndex(source, -1); matches != nil {
			for _, v := range matches {
				if i < settings.NumberOfUpperCase {
					c := strings.ToUpper(string(source[v[0]]))
					source[v[0]] = []byte(c)[0]
					i += 1
				}
			}
		}
	}

	if !containsDigits(source, settings.NumberOfDigits) {
		i := 0
		r := regexp.MustCompile(`[a-z]+`)

		var matches [][]int
		if matches = r.FindAllIndex(source, -1); matches != nil {
			for _, v := range matches {
				if i < settings.NumberOfDigits {
					source[v[0]] = byte(i)
					i += 1
				}
			}
		}
	}

	if !containsSpecialCharacters(source, settings.SpecialCharacters, settings.NumberOfSpecialCharacters) {
		i := 0
		r := regexp.MustCompile(`[a-z]+`)

		var matches [][]int
		if matches = r.FindAllIndex(source, -1); matches != nil {
			for _, v := range matches {
				if i < settings.NumberOfSpecialCharacters {
					i += 1
					source[v[0]] = []byte(settings.SpecialCharacters)[len(settings.SpecialCharacters)-i]
				}
			}
		}
	}

	return source
}

// Determine if the hash currently contains the appropriate amount of digits
func containsDigits(source []byte, minOccurrences int) bool {
	r := regexp.MustCompile(`\d`)

	var matches [][]byte
	if matches = r.FindAll(source, -1); matches == nil {
		return false
	}

	return len(matches) >= minOccurrences
}

// Determine if the hash currently contains the appropriate amount of uppercase characters
func containsUppercase(source []byte, minOccurrences int) bool {
	r := regexp.MustCompile(`[A-Z]+`)

	var matches [][]byte
	if matches = r.FindAll(source, -1); matches == nil {
		return false
	}

	return len(matches) >= minOccurrences
}

// Determine if the hash currently contains the appropriate amount of special characters from the allowed
// character set
func containsSpecialCharacters(source []byte, specialCharacters string, minOccurrences int) bool {
	s := specialCharacters
	s = strings.Replace(s, "\\", "\\\\", -1)
	s = strings.Replace(s, ".", "\\.", -1)
	s = strings.Replace(s, " ", "\\s", -1)
	s = strings.Replace(s, "-", "\\-", -1)
	s = strings.Replace(s, "[", "\\[", -1)
	s = strings.Replace(s, "]", "\\]", -1)

	r := regexp.MustCompile(`[` + s + `]+`)

	var matches [][]byte
	if matches = r.FindAll(source, -1); matches == nil {
		return false
	}

	return len(matches) >= minOccurrences
}

// Determine if the hash currently abides by the length restrictions
func validateLength(source []byte, minimum, maximum int) bool {
	if minimum > -1 && len(source) < minimum {
		return false
	}

	if maximum > -1 && len(source) > maximum {
		return false
	}

	return true
}
