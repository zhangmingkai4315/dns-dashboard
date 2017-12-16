package utils

import (
	"testing"
)

func TestGetTLDDomain(t *testing.T) {
	testCase := []struct {
		domain string
		expect string
	}{{"abc.example.com", "com"},
		{".", "."},
		{"", ""},
		{"abc.com", "com"},
		{"abc.com.", "com"}}

	for _, test := range testCase {
		result := GetTLDDomain(test.domain)
		if test.expect != result {
			t.Fatalf("RemoveSubDomain(%s) = %s != %s[expect]", test.domain, result, test.expect)
		}
	}
}

func TestRemoveSubDomain(t *testing.T) {
	testCase := []struct {
		domain string
		expect string
	}{{"abc.example.com", "example.com"},
		{".", "."},
		{"", ""},
		{"test.abc.com", "abc.com"},
		{"test.abc.com.", "abc.com"}}

	for _, test := range testCase {
		result := RemoveSubDomain(test.domain)
		if test.expect != result {
			t.Fatalf("RemoveSubDomain(%s) = %s != %s[expect]", test.domain, result, test.expect)
		}
	}
}
