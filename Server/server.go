package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"
)

type siteStatus interface {
	statuschecker()
	updateStatus()
	addURLS()
}

type websites struct {
	urls      []string `json:"urls"`
	statusMap map[string]int
}

func (siteObj *websites) statuschecker(w http.ResponseWriter, r *http.Request) {
	log.Println("Getting Status .....")
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err)
	}
	params := make(map[string][]string)
	err = json.Unmarshal(body, &params)

	if err == nil || len(params["website"]) > 0 {
		w.Write([]byte("Getting status for requested urls\n"))
		for _, val := range params["website"] {
			status, Ok := siteObj.statusMap[val]
			if Ok {
				w.Write([]byte("Last checked status for " + val + ": " + strconv.Itoa(status)))
			} else {
				log.Printf("Url: %v not found.... Adding to memory", val)
				log.Println("Getting Current status of the requested Url ......")
				res, err := http.Get(val)
				if err != nil {
					log.Fatal(err)
				}
				siteObj.statusMap[val] = res.StatusCode
				w.Write([]byte("Current status for " + val + ": " + strconv.Itoa(res.StatusCode) + "\n"))

			}
		}

	} else {
		w.Write([]byte("Getting status for all urls"))
		for _, val := range siteObj.urls {
			status, Ok := siteObj.statusMap[val]
			if Ok {
				w.Write([]byte("Last checked status for " + val + ": " + strconv.Itoa(status) + "\n"))
			} else {
				w.Write([]byte("Status for URL " + val + " not found in Memory"))
				w.Write([]byte("Getting Current Status ... "))
				res, err := http.Get(val)
				if err != nil {
					log.Fatal(err)
				}
				siteObj.statusMap[val] = res.StatusCode
				w.Write([]byte("Current status for " + val + ": " + strconv.Itoa(res.StatusCode) + "\n"))
			}
		}
	}
}

func (siteObj *websites) addURLS(w http.ResponseWriter, r *http.Request) {
	log.Println("Adding Requested URLS to the Memory .... !")
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println("Error while adding URLS ....!")
	}
	params := make(map[string][]string)
	err = json.Unmarshal(body, &params)
	if err != nil {
		log.Println("Error while Unmarshalling Body ....!")
	}
	for _, url := range params["website"] {
		_, Ok := siteObj.statusMap[url]
		if Ok {
			log.Printf("Tried adding url %v, url already present in Memory!!!\n", url)
			continue
		}
		fmt.Println(url)
		siteObj.urls = append(siteObj.urls, url)

		/////////////////////////////////////////////////////////////////////////////////////////////////
		/////  Alternate Method is to Store Status for URL while adding to Memory using Go Routine //////
		/////////////////////////////////////////////////////////////////////////////////////////////////

		// go func(url string) {
		// 	res, err := http.Get(url)
		// 	if err != nil {
		// 		log.Fatal(err)
		// 	}
		// 	siteObj.statusMap[url] = res.StatusCode
		// }(url)
	}
	w.Write([]byte("Added Requested URLS to Memory ....\n"))
}

func (siteObj *websites) updateStatus() {
	log.Println("Updating status for all urls")
	for _, val := range siteObj.urls {
		res, err := http.Get(val)
		if err != nil {
			log.Fatal(err)
		}
		siteObj.statusMap[val] = res.StatusCode
	}
}

func handleRequest() {
	siteObj := websites{[]string{}, make(map[string]int)}
	go func() {
		for {
			siteObj.updateStatus()
			time.Sleep(time.Minute)
		}
	}()
	http.HandleFunc("/addToUrlList", siteObj.addURLS)
	http.HandleFunc("/getStatus", siteObj.statuschecker)
	http.ListenAndServe("127.0.0.1:3421", nil)
}

func main() {
	fmt.Println("Server Started ........")
	handleRequest()
}
