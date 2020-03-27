package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

type resource struct {
	name     string
	document *goquery.Document
	selector string
}

type body struct {
	Data []string
}

func getDocument(URL string) *goquery.Document {
	log.Println("Creating Document for URL: ", URL)
	// Make HTTP request
	response, err := http.Get(URL)
	if err != nil {
		log.Fatal("Error making HTTP request to: { ", URL, " }: ", err)
	}
	defer response.Body.Close()

	// Create a goquery document from the HTTP response
	document, err := goquery.NewDocumentFromReader(response.Body)
	if err != nil {
		log.Fatal("Error loading HTTP response body. ", err)
	}
	log.Println("Successfully Created Document")
	return document
}

func main() {
	values := make([]string, 0, 8)

	s := time.Now().Format("02/01/2006")
	values = append(values, s)

	docs := make([]*goquery.Document, 2)

	log.Println("Attempting to create goquery documents...")
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
	log.Println("Attempting to scrape goquery documents...")
	for _, res := range resources {
		// for each resource in resources array

		// Find all counts and push to values
		res.document.Find(res.selector).Each(func(index int, element *goquery.Selection) {
			t := element.Text()
			t = strings.ReplaceAll(t, ",", "")
			values = append(values, t)
		})
	}
	log.Println("Successfully scraped Document")
	b := body{
		Data: values,
	}
	// fmt.Printf("b= %v\n", b)
	j, e := json.Marshal(b)
	// fmt.Printf("j= %v\ne= %v", string(j), e)

	sHostEnv := os.Getenv("SERVICENAME_HOST")

	var serviceHostname string

	if sHostEnv == "" {
		log.Println("Please use environment variables to designate updateListener service in a cloud deployment")
		serviceHostname = "localhost:8080"
	} else {
		serviceHostname = sHostEnv
	}

	log.Println("Attempting to send data to the update service...")
	r, e := http.Post(fmt.Sprintf("http://%v", serviceHostname), "application/json", bytes.NewBuffer(j))
	if e != nil {
		log.Fatalln("Failed to send data to the update service", e)
	} else {
		log.Println("Successfully sent data to the update service")
	}
	defer r.Body.Close()
}
