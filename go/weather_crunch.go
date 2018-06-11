// jalavoy 06.10.2018
// this program reads 48 years of historical weather data from the NCDC/NOAA and calculates averages per day
// output looks something like this:
/*
	January/1 - Chance of Snow: 15%, Chance of Rain: 11%, High Temp: 17 F, Low Temp: 4 F, Average Temperature: 11 F
	January/2 - Chance of Snow: 13%, Chance of Rain: 14%, High Temp: 17 F, Low Temp: 5 F, Average Temperature: 11 F
	January/3 - Chance of Snow: 10%, Chance of Rain: 11%, High Temp: 20 F, Low Temp: 6 F, Average Temperature: 13 F
	January/4 - Chance of Snow: 24%, Chance of Rain: 20%, High Temp: 19 F, Low Temp: 5 F, Average Temperature: 12 F
	January/5 - Chance of Snow: 32%, Chance of Rain: 29%, High Temp: 19 F, Low Temp: 4 F, Average Temperature: 12 F
	January/6 - Chance of Snow: 28%, Chance of Rain: 26%, High Temp: 20 F, Low Temp: 5 F, Average Temperature: 12 F
	January/7 - Chance of Snow: 6%, Chance of Rain: 6%, High Temp: 20 F, Low Temp: 5 F, Average Temperature: 13 F
	January/8 - Chance of Snow: 23%, Chance of Rain: 21%, High Temp: 20 F, Low Temp: 6 F, Average Temperature: 13 F
	January/9 - Chance of Snow: 10%, Chance of Rain: 6%, High Temp: 21 F, Low Temp: 7 F, Average Temperature: 14 F
*/
package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"reflect"
	"sort"
	"strconv"
	"time"
)

const csvFile string = "deps/weather.csv"

// Months ints to names - there's probably a better way to do this with the time module
var Months = [12]string{
	"January",
	"Febuary",
	"March",
	"April",
	"May",
	"June",
	"July",
	"August",
	"September",
	"October",
	"November",
	"December",
}

// WeatherDay holds the total weather information for a particular month and day (excluding year)
type WeatherDay struct {
	snow     []float64
	rain     []float64
	avgtemp  []int
	hightemp []int
	lowtemp  []int
}

func main() {
	rawdata := getCSV()
	data := parseData(rawdata)
	outputData(data)
}

// imports our data from our csv
func getCSV() [][]string {
	fh, err := os.Open(csvFile)
	if err != nil {
		log.Fatal(err)
	}
	reader := csv.NewReader(bufio.NewReader(fh))
	rawdata, err := reader.ReadAll()
	return rawdata
}

// parses our CSV data and puts it into a more managable format
func parseData(rawdata [][]string) map[int]map[int]WeatherDay {
	data := make(map[int]map[int]WeatherDay)
	// initialize inner map - there's probably a better way to do this
	for i := 1; i <= 12; i++ {
		data[i] = make(map[int]WeatherDay)
	}
	skip := true
	for _, each := range rawdata {
		// skip the first entry since it's just the column titles from the csv
		if skip == true {
			skip = false
			continue
		}
		month, day := parseDate(each[2])
		w := data[month][day]
		w.appendSnow(each[8])
		w.appendRain(each[7])
		w.appendTemp(each[10], each[11])
		data[month][day] = w
	}
	return data
}

// does the final averaging and outputs our data
func outputData(data map[int]map[int]WeatherDay) {
	for _, month := range sortKeys(data) {
		for _, day := range sortKeys(data[month]) {
			w := data[month][day]
			avgSnow := getAverageFloat(w.snow...)
			avgRain := getAverageFloat(w.rain...)
			avgHigh := getAverageInt(w.hightemp...)
			avgLow := getAverageInt(w.lowtemp...)
			avgTemp := getAverageInt(w.avgtemp...)
			fmt.Printf("%s/%d - Chance of Snow: %d%%, Chance of Rain: %d%%, High Temp: %d F, Low Temp: %d F, Average Temperature: %d F\n", Months[(month-1)], day, int(avgSnow*100), int(avgRain*100), avgHigh, avgLow, avgTemp)
		}
	}
}

// converts the date into month/day as ints - go has the coolest date parsing i've ever seen
func parseDate(date string) (int, int) {
	// Mon Jan 2 15:04:05 -0700 MST 2006
	time, err := time.Parse("2006-01-02", date)
	if err != nil {
		log.Panic(err)
	}
	return int(time.Month()), int(time.Day())
}

// calculates averages between N number of ints
func getAverageInt(inputs ...int) int {
	var total int
	var i int
	var each int
	for i, each = range inputs {
		total += each
	}
	avg := (total / (i + 1))
	return avg
}

// calculates averages between N number of floats
func getAverageFloat(inputs ...float64) float64 {
	var total float64
	var i int
	var each float64
	for i, each = range inputs {
		total += each
	}
	avg := (total / (float64(i) + 1.0))
	return avg
}

func (w *WeatherDay) appendSnow(snow string) {
	s, _ := strconv.ParseFloat(snow, 32)
	if s > 0.0 {
		w.snow = append(w.snow, 1)
	} else {
		w.snow = append(w.snow, 0)
	}
}

func (w *WeatherDay) appendRain(rain string) {
	r, _ := strconv.ParseFloat(rain, 32)
	if r > 0.0 {
		w.rain = append(w.rain, 1)
	} else {
		w.rain = append(w.rain, 0)
	}
}

func (w *WeatherDay) appendTemp(max string, min string) {
	maxInt, _ := strconv.Atoi(max)
	minInt, _ := strconv.Atoi(min)
	w.hightemp = append(w.hightemp, maxInt)
	w.lowtemp = append(w.lowtemp, minInt)
	w.avgtemp = append(w.avgtemp, getAverageInt(maxInt, minInt))
}

// sorts map keys by int
func sortKeys(m interface{}) []int {
	var l []int
	keys := reflect.ValueOf(m).MapKeys()
	for _, key := range keys {
		l = append(l, key.Interface().(int))
	}
	sort.Ints(l)
	return l
}
