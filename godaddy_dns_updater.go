package main

import (
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

	res, _ := http.Get("https://api.ipify.org?format=json")

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

func main() {
	var ipfile = &IPFile{fileName: "/tmp/ip.json"}
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
}
