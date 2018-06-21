// jalavoy 06.12.2018
// This program idea was borrowed from a computer science course at The University of Arizona.
// The course is in python, but I'm doing it in go. It tells the user what type of Pokemon would be best given the stat they specify wanting
// Looks something like this:
/*
	# ./pokemon_thing
	[!] Usage: ./pokemon_thing <data file> <power|hp|attack|defense|specialatk|specialdef|speed>
	[!] This program takes your desired stat and determines which type would be the strongest for that stat
	# ./pokemon_thing deps/Pokemon.csv defense
	[*] Results: The best class for defense would be Electric with a total of 64
*/
package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

// Class holds all the data for each class as it parses the csv
type Class struct {
	power      []int
	hp         []int
	attack     []int
	defense    []int
	specialatk []int
	specialdef []int
	speed      []int
}

func main() {
	csvFile, query := checkInput()
	data := parseData(csvFile)
	findHighest(&data, query)
}

// checkInput checks our input to make sure it's what we want
func checkInput() (string, string) {
	if len(os.Args) < 3 {
		usage()
	}
	if _, err := os.Stat(os.Args[1]); os.IsNotExist(err) {
		fmt.Println("Specified input file does not exist")
		usage()
	}
	match := regexp.MustCompile("^(power|hp|attack|defense|specialatk|specialdef|speed)$")
	if !match.MatchString(strings.ToLower(os.Args[2])) {
		usage()
	}
	return os.Args[1], strings.ToLower(os.Args[2])
}

// parseData parses the csv into our data structure. I originally wrote this to use the CSV module but since I've used that before, I decided to do it manually
func parseData(csvFile string) map[string]*Class {
	fh, err := os.Open(csvFile)
	if err != nil {
		log.Fatal(err)
	}
	defer fh.Close()

	data := make(map[string]*Class)

	scanner := bufio.NewScanner(fh)
	for scanner.Scan() {
		if strings.HasPrefix(scanner.Text(), "#") {
			continue
		}
		elements := strings.Split(scanner.Text(), ",")
		class := elements[2]
		if _, exists := data[class]; !exists {
			data[class] = &Class{}
		}
		c := data[class]
		c.addPokemon(elements)
	}
	return data
}

// findHighest loops the data structure to find the best option depending on the users input
func findHighest(data *map[string]*Class, query string) {
	watermark := ""
	last := 0
	for pType, class := range *data {
		if avg := class.getAverage(query); avg > last {
			watermark = pType
			last = avg
		}
	}
	fmt.Println("[*] Results: The best class for", query, "would be", watermark, "with a total of", last)
}

// addPokemon just fills our data
func (c *Class) addPokemon(elements []string) {
	// #,Name,Type 1,Type 2,Total,HP,Attack,Defense,Sp. Atk,Sp. Def,Speed,Generation,Legendary
	power, _ := strconv.Atoi(elements[4])
	hp, _ := strconv.Atoi(elements[5])
	attack, _ := strconv.Atoi(elements[6])
	defense, _ := strconv.Atoi(elements[7])
	specialatk, _ := strconv.Atoi(elements[8])
	specialdef, _ := strconv.Atoi(elements[9])
	speed, _ := strconv.Atoi(elements[10])

	c.power = append(c.power, power)
	c.hp = append(c.hp, hp)
	c.attack = append(c.attack, attack)
	c.defense = append(c.defense, defense)
	c.specialatk = append(c.specialatk, specialatk)
	c.specialdef = append(c.specialdef, specialdef)
	c.speed = append(c.speed, speed)
}

// getAverage uses reflection to find the value the user specified in our data.
// Reflection probably isn't the best way to do this but I wanted to learn how it worked
func (c *Class) getAverage(query string) int {
	r := reflect.ValueOf(c).Elem()
	var count int64
	var total int64
	for i := 0; i < r.NumField(); i++ {
		if r.Type().Field(i).Name == query {
			for j := 0; j < r.Field(i).Len(); j++ {
				total += r.Field(i).Index(j).Int()
				count++
			}
		}
	}
	avg := (total / count)
	return int(avg)
}

func usage() {
	fmt.Println("[!] Usage:", os.Args[0], "<data file>", "<power|hp|attack|defense|specialatk|specialdef|speed>")
	fmt.Println("[!] This program takes your desired stat and determines which type would be the strongest for that stat")
	os.Exit(1)
}
