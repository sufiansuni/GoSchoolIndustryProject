package main

import (
	"fmt"
	"net/http"
)

func testmap(res http.ResponseWriter, req *http.Request) {

	search1, err := API_OneMap_Search("13 Marsiling Lane")
	if err != nil {
		fmt.Println(err)
	}
	search2, err := API_OneMap_Search("Woodlands MRT NS9")
	if err != nil {
		fmt.Println(err)
	}

	start_lat := search1.Results[0].Latitude
	start_lng := search1.Results[0].Longitude

	end_lat := search2.Results[0].Latitude
	end_lng := search2.Results[0].Longitude
	
	MapPNG := API_OneMap_GenerateMapPNG(start_lat, start_lng, end_lat, end_lng)
	tpl.ExecuteTemplate(res, "testmap.html", MapPNG)
}
