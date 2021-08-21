package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

type API_OneMap_Search_Result struct {
	Found         int
	TotalNumPages int
	PageNum       int
	Results       []map[string]string
}

type API_OneMap_GetToken_Result struct {
	Access_Token     string
	Expiry_Timestamp string
}

func API_OneMap_Search(search_val string) API_OneMap_Search_Result {
	var result API_OneMap_Search_Result
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

func API_OneMap_GetToken() API_OneMap_GetToken_Result {
	var result API_OneMap_GetToken_Result
	url := "https://developers.onemap.sg/privateapi/auth/post/getToken"
	values := map[string]string{"email": goDotEnvVariable("API_ONEMAP_EMAIL"), "password": goDotEnvVariable("API_ONEMAP_PASSWORD")}
	jsonValue, _ := json.Marshal(values)
	if resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonValue)); err == nil {
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
