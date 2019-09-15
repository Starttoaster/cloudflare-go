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
	"time"
)

//Global variable declaration
var interval time.Duration = 20

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
func getData(client *http.Client, email string, gapik string, zone string) response {
	//Set up new http GET request
	url := "https://api.cloudflare.com/client/v4/zones/" + zone + "/dns_records"
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
func putData(client *http.Client, email string, gapik string, zone string, id string, recordType string, name string, ip string, proxied bool) {
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
	//Set up new http PUT request
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

	//Formatting response
	currentTime := time.Now()
	fmt.Println("Current Time: " + currentTime.Format("2006-01-02 3:4:5 PM") + "\nUpdated Record: " + name + "\nUpdated IP: " + ip + "\n")
}

//Reads credentials file and returns string slice
func readLines() ([]string, error) {
	file, err := os.Open("/config/credfile")
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
	//Set up HTTP client with timeout
	timeout := time.Duration(120 * time.Second)
	client := &http.Client{
		Timeout: timeout,
	}
	//Get authentication variables from file
	credLines, err := readLines()
	if err != nil {
		log.Fatalln(err)
	}
	email := credLines[0]
	gapik := credLines[1]
	zone := credLines[2]

	//Infinite loop to update records over time
	for {
		jsonData := getData(client, email, gapik, zone)
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

					putData(client, email, gapik, zone, recordIdentifier, recordType, recordName, publicIP, recordProxied)
				}
			}
		}
		//Sleeping for n seconds
		time.Sleep(interval * time.Second)
	}
}