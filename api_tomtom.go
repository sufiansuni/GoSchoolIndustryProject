package main

// TomTom: https://www.tomtom.com/
// TomTom (Developers): https://developer.tomtom.com/
// TomTom Routing API: https://developer.tomtom.com/routing-api/routing-api-documentation

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

type API_TomTom_Routing_Result struct {
	Formatversion string `json:"formatVersion"`
	Routes        []struct {
		Summary struct {
			Lengthinmeters        int       `json:"lengthInMeters"`
			Traveltimeinseconds   int       `json:"travelTimeInSeconds"`
			Trafficdelayinseconds int       `json:"trafficDelayInSeconds"`
			Trafficlengthinmeters int       `json:"trafficLengthInMeters"`
			Departuretime         time.Time `json:"departureTime"`
			Arrivaltime           time.Time `json:"arrivalTime"`
		} `json:"summary"`
		Legs []struct {
			Summary struct {
				Lengthinmeters        int       `json:"lengthInMeters"`
				Traveltimeinseconds   int       `json:"travelTimeInSeconds"`
				Trafficdelayinseconds int       `json:"trafficDelayInSeconds"`
				Trafficlengthinmeters int       `json:"trafficLengthInMeters"`
				Departuretime         time.Time `json:"departureTime"`
				Arrivaltime           time.Time `json:"arrivalTime"`
			} `json:"summary"`
			Points []struct {
				Latitude  float64 `json:"latitude"`
				Longitude float64 `json:"longitude"`
			} `json:"points"`
		} `json:"legs"`
		Sections []struct {
			Startpointindex int    `json:"startPointIndex"`
			Endpointindex   int    `json:"endPointIndex"`
			Sectiontype     string `json:"sectionType"`
			Travelmode      string `json:"travelMode"`
		} `json:"sections"`
	} `json:"routes"`
}

func API_TomTom_Routing(start_lat string, start_lng string, end_lat string, end_lng string) (API_TomTom_Routing_Result, error) {
	var result API_TomTom_Routing_Result

	my_url := "https://api.tomtom.com/routing/1/calculateRoute/" +
		start_lat + url.QueryEscape(",") +
		start_lng + url.QueryEscape(":") +
		end_lat + url.QueryEscape(",") +
		end_lng +
		"/json?routeType=fastest&traffic=true&avoid=unpavedRoads&travelMode=pedestrian&key=" +
		goDotEnvVariable("API_TOMTOM_TOKEN")

	if resp, err := http.Get(my_url); err == nil {
		defer resp.Body.Close()
		if body, err := ioutil.ReadAll(resp.Body); err == nil {
			if resp.StatusCode == http.StatusOK {
				json.Unmarshal(body, &result)
				return result, err
			} else {
				err_msg := "Error, Status Code: " + strconv.Itoa(resp.StatusCode)
				return result, errors.New(err_msg)
			}
		} else {
			return result, err
		}
	} else {
		return result, err
	}
}
