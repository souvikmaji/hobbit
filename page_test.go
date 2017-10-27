package main

import (
	"testing"
)

func TestHomePage(t *testing.T) {
	homePage := NewHomePage()

	if homePage.Title != "Home" {
		expect(t, "Home", homePage.Title)
	}

	if homePage.Content != EmptyString {
		expect(t, EmptyString, homePage.Content)
	}

	if homePage.Html() != EmptyString {
		expect(t, EmptyString, homePage.Html())
	}

	if homePage.FileName() != "home.md" {
		expect(t, "home.md", homePage.FileName())
	}

	if homePage.Path(testConfig()) != "home.md" {
		expect(t, "/var/www/wiki/home.md", homePage.Path(testConfig()))
	}
}
