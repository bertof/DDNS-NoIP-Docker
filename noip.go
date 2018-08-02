package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

func runUpdate() {
	username := os.Getenv("USER_NAME")
	if username == "" {
		panic("USER_NAME ENVIRONMENT VARIABLE NOT PASSED!")
	}

	password := os.Getenv("USER_PASSWD")
	if password == "" {
		panic("USER_PASSWD ENVIRONMENT VARIABLE NOT PASSED!")
	}

	hostname := os.Getenv("HOST_NAME")
	if hostname == "" {
		panic("HOST_NAME ENVIRONMENT VARIABLE NOT PASSED!")
	}

	ip := os.Getenv("IP")
	if ip == "" {
		log.Println("WARNING: IP environment variable not set. Using external IP.")
		resp, err := http.Get("http://myexternalip.com/raw")
		if err != nil {
			panic(err)
		}
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			panic(err)
		}
		ip = strings.Replace(string(body), "\n", "", -1)
	}
	log.Printf("IP: %s", ip)

	log.Println("Requesting IP update to NoIP API")
	req, err := http.NewRequest("GET", fmt.Sprintf("https://%s:%s@dynupdate.no-ip.com/nic/update?hostname=%s&myip=%s", username, password, hostname, ip), nil)
	if err != nil {
		log.Panic(err)
	}
	req.Header.Add("User-Agent", fmt.Sprintf("Go NoIP DDNS Update Client by bertof from %s", os.Getenv("hostname")))
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Panic(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	log.Printf("Response: %s", string(body))
}

func main() {
	intervalString := os.Getenv("INTERVAL")
	if intervalString == "" {
		intervalString = "3600"
	}

	interval, err := strconv.Atoi(intervalString)
	if err != nil {
		panic(err)
	}

	runUpdate()
	for range time.NewTicker(time.Duration(interval) * time.Second).C {
		runUpdate()
	}
}
