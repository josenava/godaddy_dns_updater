package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"sync"
)

type IPFile struct {
	fileName string
	mux      sync.Mutex
}

func needsToUpdateGodaddyDNS(previousIP, currentIP string) bool {
	return strings.Compare(previousIP, currentIP) != 0
}

func getPreviousIPFromFile(ipfile *IPFile, ch chan string) {
	ipfile.mux.Lock()
	_, err := os.Stat(ipfile.fileName)

	// create file if not exists
	if os.IsNotExist(err) {
		file, _ := os.Create(ipfile.fileName)
		defer file.Close()
	}

	previousStoredIPBytes, _ := ioutil.ReadFile(ipfile.fileName)
	// finished reading
	ipfile.mux.Unlock()

	if len(previousStoredIPBytes) == 0 {
		fmt.Println("Previous ip address was empty")
		ch <- ""
		return
	}

	var previousIP map[string]interface{}
	err = json.Unmarshal(previousStoredIPBytes, &previousIP)

	ch <- previousIP["ip"].(string)
}

func getCurrentIPFromAPIAndStoreInFile(ipfile *IPFile, ch chan string) {
	var currentIP map[string]interface{}

	res, _ := http.Get(os.Getenv("ip_finder_url"))

	ip, _ := ioutil.ReadAll(res.Body)
	copyIP := make([]byte, len(ip))
	copy(copyIP, ip)

	err := json.Unmarshal(ip, &currentIP)

	if err != nil {
		panic(err)
	}

	ch <- currentIP["ip"].(string)

	ipfile.mux.Lock()
	//update ip
	err = ioutil.WriteFile(ipfile.fileName, copyIP, 664)
	defer ipfile.mux.Unlock()

	if err != nil {
		fmt.Println(err.Error())
	}
}

func updateGodaddyDNS(url, domain, apiKey, apiSecret, currentIP string) {
	client := &http.Client{}

	fmt.Println(url)
	fmt.Println(domain)
	fmt.Println(apiKey)
	fmt.Println(apiSecret)

	reqData := []byte(fmt.Sprintf(`[{"ttl": 600, "data": "%s" }]`, currentIP))

	req, _ := http.NewRequest("PUT", fmt.Sprintf("%s/v1/domains/%s/records/A/@", url, domain), bytes.NewBuffer(reqData))
	req.Header.Add("Authorization", fmt.Sprintf("sso-key %s:%s", apiKey, apiSecret))
	req.Header.Add("Content-type", "application/json")

	
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err.Error())
	}

	defer res.Body.Close()

    fmt.Println("response Status:", res.Status)
    fmt.Println("response Headers:", res.Header)
    body, _ := ioutil.ReadAll(res.Body)
    fmt.Println("response Body:", string(body))
}

func main() {
	apiUrl := os.Getenv("godaddy_api_url")
	domain := os.Getenv("domain_url")
	apiKey := os.Getenv("godaddy_api_key")
	apiSecret := os.Getenv("godaddy_api_secret")

	var ipfile = &IPFile{fileName: os.Getenv("ip_file_path")}
	ch1 := make(chan string)
	ch2 := make(chan string)

	go getPreviousIPFromFile(ipfile, ch1)
	go getCurrentIPFromAPIAndStoreInFile(ipfile, ch2)

	previousIP, currentIP := <-ch1, <-ch2

	fmt.Printf("Previous IP %v\n Current IP %v\n", previousIP, currentIP)

	if needsToUpdateGodaddyDNS(previousIP, currentIP) == false {
		fmt.Println("Nothing to do here, exiting")
		os.Exit(0)
	}

	fmt.Println("We need to update the dns!!!")
	updateGodaddyDNS(apiUrl, domain, apiKey, apiSecret, currentIP)
}
