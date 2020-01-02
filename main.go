package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"net/http"
	"sort"
	"strconv"
)

const (
	data = "https://raw.githubusercontent.com/GoesToEleven/lynda/master/cc/01_parse_weather-data/Environmental_Data_Deep_Moor_2015.txt"
)

func main() {
	res, err := http.Get(data)
	if err != nil {
		log.Fatalf("HTTP error: %v", err)
	}
	defer res.Body.Close()

	rdr := csv.NewReader(res.Body)
	rdr.Comma = '\t'
	rdr.TrimLeadingSpace = true

	rows, err := rdr.ReadAll()
	if err != nil {
		log.Fatalf("unable to read data: %v", err)
	}

	// Don't count the header row in len(rows)
	fmt.Println("Total Records: ", len(rows)-1)

	fmt.Println("Air Temp:\t", mean(rows, 1), median(rows, 1))
	fmt.Println("Barometric:\t", mean(rows, 2), median(rows, 2))
	fmt.Println("Wind Speed:\t", mean(rows, 7), median(rows, 7))

}

// 2015_01_01 00:02:43	19.50	30.62	14.78	81.60	159.78	14.00	 9.20
func mean(rows [][]string, column int) float64 {
	var accumulator float64
	for i, r := range rows {
		if i == 0 {
			continue // Skip header row
		}
		if val, err := strconv.ParseFloat(r[column], 64); err == nil {
			accumulator += val
		} else {
			log.Printf("unable to parse string value: %v", err)
		}
	}
	return accumulator / float64(len(rows)-1)
}

func median(rows [][]string, column int) (median float64) {
	var sorted []float64
	for i, r := range rows {
		if i == 0 {
			continue // Skip header row
		}
		if val, err := strconv.ParseFloat(r[column], 64); err == nil {
			sorted = append(sorted, val)
		} else {
			log.Printf("unable to parse string value: %v", err)
		}
	}
	sort.Float64s(sorted)
	if (len(sorted) % 2) == 0 {
		median = (sorted[len(sorted)-1] + sorted[len(sorted)]) / 2
	} else {
		median = sorted[len(sorted)/2]
	}
	return median
}
