package main

import (
	"fmt"
	"log"
	"testing"
)

func TestGeneratePassphrase(t *testing.T) {
	profile := "guest"
	//passphrase := "secret"
	site := Site{
		Host:                      "google.com",
		MinimumLength:             6,
		MaximumLength:             18,
		SpecialCharacters:         " !@#$%^&*()_+-=<>,./|\\",
		NumberOfSpecialCharacters: 2,
		NumberOfDigits:            2,
		NumberOfUpperCase:         1,
	}
	i := 0

	expected := "15898B4\\3738|f656d"
	for i = 0; i < 2; i++ {
		b, _ := generatePassphrase(profile, site)
		actual := fmt.Sprintf("%s", string(b))
		if actual != expected {
			log.Println(actual)
			t.FailNow()
		}
	}
}

func TestContainsDigits(t *testing.T) {
	expected := true
	actual := containsDigits([]byte("blah0blah2"), 2)
	if actual != expected {
		t.FailNow()
	}

	expected = false
	actual = containsDigits([]byte("blah0blah"), 2)
	if actual != expected {
		t.FailNow()
	}

	expected = true
	actual = containsDigits([]byte("blah0blah12"), 2)
	if actual != expected {
		t.FailNow()
	}
}

func TestContainsUppercase(t *testing.T) {
	expected := true
	actual := containsUppercase([]byte("Blah0Blah2"), 2)
	if actual != expected {
		t.FailNow()
	}

	expected = false
	actual = containsUppercase([]byte("Blah0blah"), 2)
	if actual != expected {
		t.FailNow()
	}

	expected = true
	actual = containsUppercase([]byte("Blah0BLaH12"), 2)
	if actual != expected {
		t.FailNow()
	}
}

func TestContainsSpecialCharacters(t *testing.T) {
	special := "~!@#$%^&*()_+-= []{};':\",./<>?\\|"
	expected := true
	actual := containsSpecialCharacters([]byte("Blah!Blah&"), special, 2)
	if actual != expected {
		t.FailNow()
	}

	expected = false
	actual = containsSpecialCharacters([]byte("Blah!blah"), special, 2)
	if actual != expected {
		t.FailNow()
	}

	expected = true
	actual = containsSpecialCharacters([]byte("Blah!BLa l."), special, 2)
	if actual != expected {
		t.FailNow()
	}
}

func TestLength(t *testing.T) {
	expected := true
	actual := validateLength([]byte("eightormore"), 8, -1)
	if actual != expected {
		t.FailNow()
	}

	expected = true
	actual = validateLength([]byte("eighttoten"), 8, 10)
	if actual != expected {
		t.FailNow()
	}

	expected = false
	actual = validateLength([]byte("eighttoten----"), 8, 10)
	if actual != expected {
		t.FailNow()
	}

	expected = true
	actual = validateLength([]byte("maxof8--"), -1, 8)
	if actual != expected {
		t.FailNow()
	}

	expected = false
	actual = validateLength([]byte("ei"), 8, -1)
	if actual != expected {
		t.FailNow()
	}
}
