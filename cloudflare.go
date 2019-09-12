package main

import (
	"net/http"
	"io/ioutil"
	"log"
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"bufio"
)

//JSON response struct
type response struct {
	Result []struct {
		Identifier string `json:"id"`
		Type       string `json:"type"`
		Name       string `json:"name"`
		Proxied    bool   `json:"proxied"`
		Content    string `json:"content"`
	} `json:"result"`
}
//JSON PUT struct
type sendme struct {
	RecordType string `json:"type"`
	Name       string `json:"name"`
	Content    string `json:"content"`
	Proxied    bool   `json:"proxied"`
}

//Gets data from Cloudflare API's DNS zone records and returns a structure of relevant data
func getData(email string, gapik string, zone string) response {
	//Set up HTTP client and a new http GET request
	url := "https://api.cloudflare.com/client/v4/zones/" + zone + "/dns_records"
	client := &http.Client{
	}
	req, err := http.NewRequest("GET", url, nil)
	req.Header.Set("X-Auth-Email", email)
	req.Header.Set("X-Auth-Key", gapik)
	req.Header.Set("Content-type", "application/json")

	//Do GET request
	resp, err := client.Do(req)
	if err != nil {
        log.Fatalln(err)
    }
	//Closing request
	defer resp.Body.Close()
	//Reading response body and converting to usable variable
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	//Uses JSON response structs to grab relevant key values
	var jsonData response
	err = json.Unmarshal([]byte(body), &jsonData)
    if err != nil {
        log.Fatalln(err)
    }
   	return jsonData
}

//Finds the computer's current public IP address and returns it
func getIP() string {
    resp, err := http.Get("http://checkip.amazonaws.com")
	if err != nil {
		log.Fatalln(err)
	}
    defer resp.Body.Close()
    body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
    return string(bytes.TrimSpace(body))
}

//Puts data to Cloudflare API's DNS zone records
func putData(email string, gapik string, zone string, id string, recordType string, name string, ip string, proxied bool) {
	url := "https://api.cloudflare.com/client/v4/zones/" + zone + "/dns_records/" + id
	//Initialize struct and marshal to JSON
	data := sendme{
		RecordType: recordType,
		Name:       name,
		Content:    ip,
		Proxied:    proxied,
	}
	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Fatalln(err)
	}
	//Set up HTTP client and a new http PUT request
	client := &http.Client{
	}
	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(jsonData))
	req.Header.Set("X-Auth-Email", email)
	req.Header.Set("X-Auth-Key", gapik)
	req.Header.Set("Content-type", "application/json")

	//Do PUT request
	resp, err := client.Do(req)
	if err != nil {
        log.Fatalln(err)
    }
	//Closing request
	defer resp.Body.Close()
	//Reading response body and converting to usable variable
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(string(body))
}

//Reads credentials file and returns string slice
func readLines() ([]string, error) {
    file, err := os.Open("credfile")
	if err != nil {
		log.Fatalln(err)
	}
    defer file.Close()

    var lines []string
    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        lines = append(lines, scanner.Text())
    }
    return lines, scanner.Err()
}

func main() {
	//Get authentication variables from file
	credLines, err := readLines()
	if err != nil {
		log.Fatalln(err)
	}
	email := credLines[0]
	gapik := credLines[1]
	zone := credLines[2]

	jsonData := getData(email, gapik, zone)
	numOfRecords := len(jsonData.Result)
	publicIP := getIP()
	for i := 0; i < numOfRecords; i++ {
		recordType := jsonData.Result[i].Type
		recordIP := jsonData.Result[i].Content
		//Filter for only A records to update
		if(recordType == "A") {
			//Only proceed if the record's IP address on file is different from the current one
			if(recordIP != publicIP) {
				//Define the other required variables for this record
				recordIdentifier := jsonData.Result[i].Identifier
				recordName := jsonData.Result[i].Name
				recordProxied := jsonData.Result[i].Proxied
				
				putData(email, gapik, zone, recordIdentifier, recordType, recordName, publicIP, recordProxied)
			}
		}
	}
}