package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

type resource struct {
	name     string
	document *goquery.Document
	selector string
}

func getDocument(URL string) *goquery.Document {
	// Make HTTP request
	response, err := http.Get(URL)
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()

	// Create a goquery document from the HTTP response
	document, err := goquery.NewDocumentFromReader(response.Body)
	if err != nil {
		log.Fatal("Error loading HTTP response body. ", err)
	}

	return document
}

// This will get called for each HTML element found
func processElement(index int, element *goquery.Selection) {
	// See if the href attribute exists on the element
	fmt.Println(element.Text())
}

func main() {
	values := make([]string, 0, 8)

	s := time.Now().Format("02/01/2006")
	values = append(values, s)

	docs := make([]*goquery.Document, 2)

	// get documents
	docs[0] = getDocument("https://www.canada.ca/en/public-health/services/diseases/2019-novel-coronavirus-infection.html")
	docs[1] = getDocument("https://www.quebec.ca/en/health/health-issues/a-z/2019-coronavirus/situation-coronavirus-in-quebec/")

	resources := [7]resource{
		resource{
			name:     "probQuebec",
			document: docs[0],
			selector: ".table-striped > tbody:nth-child(3) > tr:nth-child(6) > td:nth-child(3)",
		},
		resource{
			name:     "probCanada",
			document: docs[0],
			selector: ".table-striped > tbody:nth-child(3) > tr:nth-child(15) > td:nth-child(3) > strong:nth-child(1)",
		},
		resource{
			name:     "confMontreal",
			document: docs[1],
			selector: ".contenttable > tbody:nth-child(3) > tr:nth-child(6) > td:nth-child(2) > p:nth-child(1)",
		},
		resource{
			name:     "confQuebec",
			document: docs[0],
			selector: ".table-striped > tbody:nth-child(3) > tr:nth-child(6) > td:nth-child(2)",
		},
		resource{
			name:     "confCanada",
			document: docs[0],
			selector: ".table-striped > tbody:nth-child(3) > tr:nth-child(15) > td:nth-child(2) > strong:nth-child(1)",
		},
		resource{
			name:     "deathQuebec",
			document: docs[0],
			selector: ".table-striped > tbody:nth-child(3) > tr:nth-child(6) > td:nth-child(4)",
		},
		resource{
			name:     "deathCanada",
			document: docs[0],
			selector: ".table-striped > tbody:nth-child(3) > tr:nth-child(15) > td:nth-child(4) > strong:nth-child(1)",
		},
	}
	for _, res := range resources {
		// for each resource in resources array

		// Find all links and process them with the function
		// defined earlier
		res.document.Find(res.selector).Each(func(index int, element *goquery.Selection) {
			t := element.Text()
			t = strings.ReplaceAll(t, ",", "")
			values = append(values, t)
		})
	}
	fmt.Printf("\n%v\n", values)
}
