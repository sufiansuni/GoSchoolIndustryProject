package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"

	geo "github.com/kellydunn/golang-geo"
)

// OneMAP API Documentation: https://www.onemap.gov.sg/docs/
// OneMap API testsite: https://app.swaggerhub.com/apis/onemap-sg/new-onemap-api/1.0.4
// Structs formed with the assistance of https://mholt.github.io/json-to-go/

var OneMapEmail string = "email@email.com"
var OneMapPassword string = "password"

type OneMapSearchResult struct {
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

type OneMapGetTokenResult struct {
	Access_Token     string `json:"access_token"`
	Expiry_Timestamp string `json:"expiry_timestamp"`
}

type OneMapErrorResult struct {
	Error string `json:"error"`
}

// Function sends GET request to OneMap Search API and returns the unmarshaled json response
func OneMapSearch(search_val string) (OneMapSearchResult, error) {

	var result OneMapSearchResult

	// make string safe for http query
	search_val = url.QueryEscape(search_val)

	my_url := "https://developers.onemap.sg/commonapi/search?searchVal=" + search_val + "&returnGeom=Y&getAddrDetails=Y"
	if resp, err := http.Get(my_url); err == nil {
		defer resp.Body.Close()
		if body, err := ioutil.ReadAll(resp.Body); err == nil {
			if resp.StatusCode == http.StatusOK {
				json.Unmarshal(body, &result)
				return result, err
			} else {
				var error_result OneMapErrorResult
				json.Unmarshal(body, &error_result)
				err_msg := error_result.Error + " Status Code: " + strconv.Itoa(resp.StatusCode)
				return result, errors.New(err_msg)
			}
		} else {
			return result, err
		}
	} else {
		return result, err
	}
}

// Function sends POST request to OneMap GetToken API and returns the unmarshaled json response.
// Requires email and password in .env file
func OneMapGetToken() (OneMapGetTokenResult, error) {

	var result OneMapGetTokenResult

	url := "https://developers.onemap.sg/privateapi/auth/post/getToken"
	values := map[string]string{
		"email":    OneMapEmail,
		"password": OneMapPassword,
	}
	jsonValue, _ := json.Marshal(values)
	if resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonValue)); err == nil {
		defer resp.Body.Close()
		if body, err := ioutil.ReadAll(resp.Body); err == nil {
			if resp.StatusCode == http.StatusOK {
				json.Unmarshal(body, &result)
				return result, err
			} else {
				var error_result OneMapErrorResult
				json.Unmarshal(body, &error_result)
				err_msg := error_result.Error + " Status Code: " + strconv.Itoa(resp.StatusCode)
				return result, errors.New(err_msg)
			}
		} else {
			return result, err
		}
	} else {
		return result, err
	}
}

// Function makes use of TomTomRouting() to processes co-ords and return an image link (OneMap)
// Returns a suitable request string to OneMAP Static Map API.
// Note that the result the above request is of PNG format.
// Pass it to a html template. Eg: <img src = {{.}} alt="Map">
func OneMapGenerateMapPNGTwoPoints(start_lat string, start_lng string, end_lat string, end_lng string) (MapPNG string, Distance int) {

	route, err := TomTomRouting(start_lat, start_lng, end_lat, end_lng)
	if err != nil {
		fmt.Println(err)
	}

	mid_lat, mid_lng, err := find_mid(start_lat, start_lng, end_lat, end_lng)
	if err != nil {
		fmt.Println(err)
	}

	var lines string
	lines = "["
	lines += "[" + start_lat + "," + start_lng + "]"
	for _, v := range route.Routes[0].Legs[0].Points {
		lines += ",[" + strconv.FormatFloat(v.Latitude, 'f', -1, 64) + "," + strconv.FormatFloat(v.Longitude, 'f', -1, 64) + "]"
	}
	lines += ",[" + end_lat + "," + end_lng + "]"
	lines += "]:177,0,0:3" //Line R,G,B,Thickness

	var points string
	points += "[" + start_lat + "," + start_lng + ",%22" + "175,50,0" + "%22,%22" + "A" + "%22]" //R,G,B,Label
	points += "|"
	points += "[" + end_lat + "," + end_lng + ",%22" + "255,255,178" + "%22,%22" + "B" + "%22]" //R,G,B,Label

	MapPNG = "https://developers.onemap.sg/commonapi/staticmap/getStaticImage?layerchosen=default&" +
		"&lat=" + mid_lat +
		"&lng=" + mid_lng +
		"&zoom=16" +
		"&height=512" +
		"&width=400" +
		"&lines=" + lines +
		"&points=" + points +
		"&color=" +
		"&fillColor="

	return MapPNG, route.Routes[0].Summary.Lengthinmeters
}

// Convert string inputs into Point type and return a mid-point
func find_mid(start_lat string, start_lng string, end_lat string, end_lng string) (string, string, error) {

	start_lat_conv, err := strconv.ParseFloat(start_lat, 64)
	if err != nil {
		return "", "", err
	}

	start_lng_conv, err := strconv.ParseFloat(start_lng, 64)
	if err != nil {
		return "", "", err
	}

	end_lat_conv, err := strconv.ParseFloat(end_lat, 64)
	if err != nil {
		return "", "", err
	}

	end_lng_conv, err := strconv.ParseFloat(end_lng, 64)
	if err != nil {
		return "", "", err
	}

	start := geo.NewPoint(start_lat_conv, start_lng_conv)
	end := geo.NewPoint(end_lat_conv, end_lng_conv)

	mid := start.MidpointTo(end)
	mid_lat := strconv.FormatFloat(mid.Lat(), 'f', -1, 64)
	mid_lng := strconv.FormatFloat(mid.Lng(), 'f', -1, 64)

	return mid_lat, mid_lng, err

}

// Function processes co-ords and return an image link (OneMap)
func OneMapGenerateMapPNGSingle(lat string, lng string) string {
	points := "[" + lat + "," + lng + ",%22" + "175,50,0" + "%22,%22" + "A" + "%22]" //R,G,B,Label

	MapPNG := "https://developers.onemap.sg/commonapi/staticmap/getStaticImage?layerchosen=default&" +
		"&lat=" + lat +
		"&lng=" + lng +
		"&zoom=17" +
		"&height=512" +
		"&width=400" +
		"&points=" + points +
		"&color=" +
		"&fillColor="

	return MapPNG
}
