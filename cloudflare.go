package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

//Global variable declarations
//var templateFile string = "/go/src/cloudflare/templates/index.html"
var templateFile string = "./templates/index.html"
var interval time.Duration = 30
var tableData []table = make([]table, 10)
var setTime []string = make([]string, 10)
var numOfRecords int
var email string
var gapik string
var zone string

//For html table
type table struct {
	Name  string
	IP    string
	Proxy bool
	Time  string
}

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
	if instruction == nil {
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

//Reads credentials via shell variables
func getCredentials() {
	email = os.Getenv("CF_EMAIL")
	gapik = os.Getenv("CF_KEY")
	zone = os.Getenv("CF_ZONE")

	if email == "" {
		fmt.Println("Account email is not set")
		os.Exit(1)
	} else if gapik == "" {
		fmt.Println("Account API key is not set")
		os.Exit(1)
	} else if zone == "" {
		fmt.Println("Cloudflare zone is not set")
		os.Exit(1)
	}
}

func update(client *http.Client) {
	fmt.Println("Starting program successfully... Output will be displayed only when records are updated.")
	//Infinite loop to update records over time
	for {
		//GETS current record information
		url := "https://api.cloudflare.com/client/v4/zones/" + zone + "/dns_records"
		body := httpRequest(client, "GET", url, nil, email, gapik, zone)
		jsonData := unjsonify(body)

		numOfRecords = len(jsonData.Result)
		publicIP := getIP()
		for i := 0; i < numOfRecords; i++ {
			recordType := jsonData.Result[i].Type
			recordIP := jsonData.Result[i].Content
			recordIdentifier := jsonData.Result[i].Identifier
			recordName := jsonData.Result[i].Name
			recordProxied := jsonData.Result[i].Proxied

			//Proceeds if is an A Record, AND current IP differs from recorded one
			if recordType == "A" && recordIP != publicIP {
				jsonData := jsonify(recordType, recordName, publicIP, recordProxied) //Creates JSON payload
				//PUTS new record information
				recordURL := url + "/" + recordIdentifier
				httpRequest(client, "PUT", recordURL, jsonData, email, gapik, zone)

				//Prints after successful update
				setTime[i] = time.Now().Format("2006-01-02 3:4:5 PM")

				tableData[i] = table{
					Name:  recordName,
					IP:    publicIP,
					Proxy: recordProxied,
					Time:  setTime[i],
				}
				fmt.Println("Current Time: " + setTime[i] + "\nUpdated Record: " + recordName + "\nUpdated IP: " + publicIP + "\n")
			} else {
				tableData[i] = table{
					Name:  recordName,
					IP:    publicIP,
					Proxy: recordProxied,
					Time:  setTime[i],
				}
			}
		}
		time.Sleep(interval * time.Second) //Sleeping for n seconds
	}
}

//Creates a new mux with handler(s)
func setupHandlers() *http.ServeMux {
	mux := http.NewServeMux()

	//Dynamic handlers
	mux.HandleFunc("/", indexHandler)

	return mux
}
func indexHandler(w http.ResponseWriter, r *http.Request) {
	indexTmpl, err := template.ParseFiles(templateFile)
	if err != nil {
		log.Fatalln(err)
	}
	indexTmpl.Execute(w, tableData[0:numOfRecords])
}

func main() {
	getCredentials()

	timeout := time.Duration(120 * time.Second)
	client := &http.Client{
		Timeout: timeout,
	}

	go update(client)

	mux := setupHandlers()
	log.Fatalln(http.ListenAndServe(":8080", mux))
}
