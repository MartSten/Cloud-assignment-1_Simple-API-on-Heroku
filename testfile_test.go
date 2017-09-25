package main

import (
	"testing"
)

func Test_getContributor(t *testing.T) {

	testURL := "https://api.github.com/repos/apache/camel/contributors"
	x := getContributor(testURL)
	if len(x) < 1 {
		t.Error("Faild to get contributors")
	}

}

func Test_getLanguages(t *testing.T) {

	testURL := "https://api.github.com/repos/apache/camel/languages"
	x := getLanguage(testURL)
	if len(x.Language) < 1 {
		t.Error("Faild to get languages")
	}

}
