package main

import (
	"crypto/sha512"
	"fmt"
	"log"
	"regexp"
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
	Password                  string `json:",omitempty"`
}

// Generate the passphrase
func (s *Site) generatePassphrase(profile, passphrase string) []byte {
	clearText := fmt.Sprintf(
		"%s-%s-%s-%s",
		strings.ToLower(profile),
		strings.ToLower(passphrase),
		strings.ToLower(s.Host),
		s.Revision)

	sha := sha512.New()
	sha.Write([]byte(clearText))

	return s.applyCriteria(sha.Sum(nil))
}

func (s *Site) applyCriteria(sha []byte) []byte {
	hash := []byte(fmt.Sprintf("%x", sha))

	if !containsUppercase(hash, s.NumberOfUpperCase) {
		i := 0
		r := regexp.MustCompile(`[a-z]+`)

		var matches [][]int
		if matches = r.FindAllIndex(hash, -1); matches != nil {
			for _, v := range matches {
				if i < s.NumberOfUpperCase {
					c := strings.ToUpper(string(hash[v[0]]))
					hash[v[0]] = []byte(c)[0]
					i += 1
				}
			}
		}
	}

	if !containsDigits(hash, s.NumberOfDigits) {
		i := 0
		r := regexp.MustCompile(`[a-z]+`)

		var matches [][]int
		if matches = r.FindAllIndex(hash, -1); matches != nil {
			for _, v := range matches {
				if i < s.NumberOfDigits {
					hash[v[0]] = byte(i)
					i += 1
				}
			}
		}
	}

	if !containsSpecialCharacters(hash, s.SpecialCharacters, s.NumberOfSpecialCharacters) {
		i := 0
		r := regexp.MustCompile(`[a-z]+`)

		var matches [][]int
		if matches = r.FindAllIndex(hash, -1); matches != nil {
			for _, v := range matches {
				if i < s.NumberOfSpecialCharacters {
					i += 1
					hash[v[0]] = []byte(s.SpecialCharacters)[len(s.SpecialCharacters)-i]
				}
			}
		}
	}

	// If there is a maximum length truncate the hash
	if s.MaximumLength > -1 {
		hash = hash[:s.MaximumLength]
	}

	// Ensure the length is adequate
	if !validateLength(hash, s.MinimumLength, s.MaximumLength) {
		log.Println("Does not meed the length requirements")
	}

	return hash
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
