// jalavoy 06.10.2018
// this program reads 48 years of historical weather data from the NCDC/NOAA and calculates averages per day
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

// Months ints to names
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

// WeatherDay holds the weather information for a particular day in time
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
	//spew.Dump(data)
}

func getCSV() [][]string {
	fh, err := os.Open(csvFile)
	if err != nil {
		log.Fatal(err)
	}
	reader := csv.NewReader(bufio.NewReader(fh))
	rawdata, err := reader.ReadAll()
	return rawdata
}

func parseData(rawdata [][]string) map[int]map[int]WeatherDay {
	data := make(map[int]map[int]WeatherDay)
	// initialize inner loop - there's probably a better way to do this
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

func parseDate(date string) (int, int) {
	// Mon Jan 2 15:04:05 -0700 MST 2006
	time, err := time.Parse("2006-01-02", date)
	if err != nil {
		log.Panic(err)
	}
	return int(time.Month()), int(time.Day())
}

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

func sortKeys(m interface{}) []int {
	var l []int
	keys := reflect.ValueOf(m).MapKeys()
	for _, key := range keys {
		l = append(l, key.Interface().(int))
	}
	sort.Ints(l)
	return l
}
