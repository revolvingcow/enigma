package main

import "testing"

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
