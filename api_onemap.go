package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

type OneMapAPI_SearchResult struct {
	Found         int
	TotalNumPages int
	PageNum       int
	Results       []map[string]string
}

func OneAPI_Search(search_val string) OneMapAPI_SearchResult {
	var result OneMapAPI_SearchResult
	url := "https://developers.onemap.sg/commonapi/search?searchVal=" + search_val + "&returnGeom=Y&getAddrDetails=Y"
	if resp, err := http.Get(url); err == nil {
		defer resp.Body.Close()
		if body, err := ioutil.ReadAll(resp.Body); err == nil {
			json.Unmarshal(body, &result)
		} else {
			log.Fatal(err)
		}
	} else {
		log.Fatal(err)
	}
	return result
}
