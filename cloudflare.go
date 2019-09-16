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

//Global variable declarations
var interval time.Duration = 20
var credentialsFile string = "/config/credfile"

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

//Uniform http.NewRequest template for mutliple operations
func httpRequest(client *http.Client, reqType string, url string, instruction []byte, email string, gapik string, zone string) []byte {
	var req *http.Request
	var err error
	if(instruction == nil) {
		req, err = http.NewRequest(reqType, url, nil)
	} else {
		req, err = http.NewRequest(reqType, url, bytes.NewBuffer(instruction))
	}
	if err != nil {
		log.Fatalln(err)
	}
	req.Header.Set("X-Auth-Email", email)
	req.Header.Set("X-Auth-Key", gapik)
	req.Header.Set("Content-type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		log.Fatalln(err)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	return body
}

//Unmarshals the JSON response into the 'response' struct type
func unjsonify(body []byte) response {
	var jsonData response
	err := json.Unmarshal([]byte(body), &jsonData)
	if err != nil {
		log.Fatalln(err)
	}
	return jsonData
}

//Creates a JSON payload of type 'sendme'
func jsonify(recordType string, name string, ip string, proxied bool) []byte {
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

//Reads credentials file and returns string slice
func readLines(credPath string) ([]string, error) {
	file, err := os.Open(credPath)
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
	credLines, err := readLines(credentialsFile)
	if err != nil {
		log.Fatalln(err)
	}
	email := credLines[0]
	gapik := credLines[1]
	zone := credLines[2]

	fmt.Println("Starting program successfully... Output will be displayed only when records are updated.")
	//Infinite loop to update records over time
	for {
		//GETS current record information
		url := "https://api.cloudflare.com/client/v4/zones/" + zone + "/dns_records"
		body := httpRequest(client, "GET", url, nil, email, gapik, zone)
		jsonData := unjsonify(body)

		numOfRecords := len(jsonData.Result)
		publicIP := getIP()
		for i := 0; i < numOfRecords; i++ {
			recordType := jsonData.Result[i].Type
			recordIP := jsonData.Result[i].Content
			//Only proceeds if is an A Record, and current IP differs from recorded one
			if(recordType == "A" && recordIP != publicIP) {
				//Fetches more required variables
				recordIdentifier := jsonData.Result[i].Identifier
				recordName := jsonData.Result[i].Name
				recordProxied := jsonData.Result[i].Proxied

				//Creates JSON payload
				jsonData := jsonify(recordType, recordName, publicIP, recordProxied)
				//PUTS new record information
				recordUrl := url + "/" + recordIdentifier
				httpRequest(client, "PUT", recordUrl, jsonData, email, gapik, zone)

				//Prints after successful update
				currentTime := time.Now()
				fmt.Println("Current Time: " + currentTime.Format("2006-01-02 3:4:5 PM") + "\nUpdated Record: " + recordName + "\nUpdated IP: " + publicIP + "\n")
			}
		}
		//Sleeping for n seconds
		time.Sleep(interval * time.Second)
	}
}