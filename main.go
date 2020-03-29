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

	"github.com/ssd71/ex8j/csvget"

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

type data struct {
	Prob  string
	Conf  string
	Death string
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

func cleanStrings(arr []string) {
	for idx := range arr {
		arr[idx] = strings.ReplaceAll(arr[idx], ",", "")
		arr[idx] = strings.ReplaceAll(arr[idx], " ", "")
	}
}

func main() {
	dt := csvget.ReadCSVFromURL("https://health-infobase.canada.ca/src/data/covidLive/covid19.csv")
	canadaRow := dt.GetCurrentFromUID("1")
	canadaData := data{
		Conf:  canadaRow.Values[4],
		Prob:  canadaRow.Values[5],
		Death: canadaRow.Values[6],
	}
	quebecRow := dt.GetCurrentFromUID("24")
	quebecData := data{
		Conf:  quebecRow.Values[4],
		Prob:  quebecRow.Values[5],
		Death: quebecRow.Values[6],
	}
	fmt.Println(canadaData, quebecData)

	var confMont string
	mondoc := getDocument("https://www.quebec.ca/en/health/health-issues/a-z/2019-coronavirus/situation-coronavirus-in-quebec/")
	mondoc.Find(".contenttable > tbody:nth-child(3) > tr:nth-child(6) > td:nth-child(2) > p:nth-child(1)").Each(func(i int, e *goquery.Selection) {
		t := e.Text()
		t = strings.ReplaceAll(t, ",", "")
		t = strings.ReplaceAll(t, " ", "")
		confMont = t
	})

	values := []string{time.Now().Format("2 Jan 2006 15:04:05"), quebecData.Prob, canadaData.Prob, confMont, quebecData.Conf, canadaData.Conf, quebecData.Death, canadaData.Death}
	cleanStrings(values)
	b := body{
		Data: values,
	}

	j, e := json.Marshal(b)
	fmt.Printf("j= %v\ne= %v", string(j), e)

	sHostEnv := fmt.Sprintf("%v:%v", os.Getenv("UPDATE_LISTENER_SERVICE_HOST"), os.Getenv("UPDATE_LISTENER_SERVICE_PORT"))

	var serviceHostname string

	if sHostEnv == ":" {
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
