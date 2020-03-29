package csvget

import (
	"bytes"
	"encoding/csv"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

// Table is the base operational data type in this package
type Table struct {
	Data []Row
}

// Row is a data structure to store Rows, mostly for my sanity
type Row struct {
	Values []string
}

// ReadCSV reads a local csv file and returns a Table
func ReadCSV(filename string) Table {
	file, err := os.Open(filename)
	if err != nil {
		log.Fatalln("Could not read csv files: '", filename, "': ", err)
	}
	r := csv.NewReader(file)
	var data []Row
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		data = append(data, Row{Values: record})
	}
	return Table{
		Data: data,
	}
}

// GetCurrentFromUID gets the latest situational data from a CSV file(specific to my usecase)
func (t Table) GetCurrentFromUID(uid string) Row {
	var curRecord Row
	for _, r := range t.Data {
		if r.Values[0] == uid {
			curRecord = r
		}
	}
	return curRecord
}

// ReadCSVFromURL reads a CSV file from a remote URL
func ReadCSVFromURL(resURL string) Table {
	// res, e := http.Get("https://health-infobase.canada.ca/src/data/covidLive/covid19.csv")
	res, e := http.Get(resURL)
	if e != nil {
		log.Fatalf("Error reading CSV from URL:{%v}: %v", resURL, e)
	}
	d, e := ioutil.ReadAll(res.Body)
	if e != nil {
		log.Fatalf("Error reading CSV from URL:{%v}: %v", resURL, e)
	}
	r := bytes.NewReader(d)
	cr := csv.NewReader(r)
	var data []Row
	for {
		record, err := cr.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Error reading CSV from URL:{%v}: %v", resURL, err)
		}
		data = append(data, Row{Values: record})
	}
	return Table{
		Data: data,
	}
}
