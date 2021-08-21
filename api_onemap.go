package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

//OneMAP API Documentation: https://www.onemap.gov.sg/docs/
//OneMap API testsite: https://app.swaggerhub.com/apis/onemap-sg/new-onemap-api/1.0.4
//structs formed with the assistance of https://mholt.github.io/json-to-go/

type API_OneMap_Search_Result struct {
	Found         int `json:"found"`
	Totalnumpages int `json:"totalNumPages"`
	Pagenum       int `json:"pageNum"`
	Results       []struct {
		Searchval  string `json:"SEARCHVAL"`
		BlkNo      string `json:"BLK_NO"`
		RoadName   string `json:"ROAD_NAME"`
		Building   string `json:"BUILDING"`
		Address    string `json:"ADDRESS"`
		Postal     string `json:"POSTAL"`
		X          string `json:"X"`
		Y          string `json:"Y"`
		Latitude   string `json:"LATITUDE"`
		Longitude  string `json:"LONGITUDE"`
		Longtitude string `json:"LONGTITUDE"`
	} `json:"results"`
}

type API_OneMap_GetToken_Result struct {
	Access_Token     string `json:"access_token"`
	Expiry_Timestamp string `json:"expiry_timestamp"`
}

type API_OneMap_Error_Result struct {
	Error string `json:"error"`
}

func API_OneMap_Search(search_val string) (interface{}, error) {
	var result API_OneMap_Search_Result
	var error_result API_OneMap_Error_Result
	url := "https://developers.onemap.sg/commonapi/search?searchVal=" + search_val + "&returnGeom=Y&getAddrDetails=Y"
	if resp, err := http.Get(url); err == nil {
		defer resp.Body.Close()
		if body, err := ioutil.ReadAll(resp.Body); err == nil {
			if resp.StatusCode == http.StatusOK {
				json.Unmarshal(body, &result)
				return result, err
			} else {
				json.Unmarshal(body, &error_result)
				error_result.Error += " Status Code: " + string(rune(resp.StatusCode))
				return error_result, err
			}
		} else {
			return error_result, err
		}
	} else {
		return error_result, err
	}
}

func API_OneMap_GetToken() (interface{}, error) {
	var result API_OneMap_GetToken_Result
	var error_result API_OneMap_Error_Result
	url := "https://developers.onemap.sg/privateapi/auth/post/getToken"
	values := map[string]string{
		"email":    goDotEnvVariable("API_ONEMAP_EMAIL"),
		"password": goDotEnvVariable("API_ONEMAP_PASSWORD"),
	}
	jsonValue, _ := json.Marshal(values)
	if resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonValue)); err == nil {
		defer resp.Body.Close()
		if body, err := ioutil.ReadAll(resp.Body); err == nil {
			if resp.StatusCode == http.StatusOK {
				json.Unmarshal(body, &result)
				return result, err
			} else {
				json.Unmarshal(body, &error_result)
				error_result.Error += " Status Code: " + string(rune(resp.StatusCode))
				return error_result, err
			}
		} else {
			return error_result, err
		}
	} else {
		return error_result, err
	}
}
