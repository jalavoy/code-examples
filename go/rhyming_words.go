// jalavoy - 06.04.2018
// this program takes input from the user and attempts to find words that perfect rhyme with the specified input.
// a perfect rhyme is defined as a rhyme where the following conditions are met:
// 1) The stressed vowel sound (the primary stress) in both words must be identical, as well as any subsequent sounds.
// 2) The sound that precedes the stressed vowel in the words must not be the same.
package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

// Target keeping track of our target
type Target struct {
	word          string
	pronunciation []string
	primaryStress string
	stressIndex   int
	rhymes        []string
}

// declare our target
var target = &Target{}

func main() {
	// usage check
	if len(os.Args) < 2 {
		usage()
	}
	fmt.Println("[*] This program uses a pronunciation dictionary to find words that rhyme with the input provided")
	// parse our pronunctiation dictionary
	words := parseDictionary()
	// check to make sure our target word is defined in the pronunciation map, as well as if it has more than one pronunciation
	getTarget(target, &words)
	// start looking for rhymes
	findRhymes(target, &words)
	outputResults(target)
	os.Exit(0)
}

func usage() {
	fmt.Println("Usage:", os.Args[0], "<word to rhyme>")
	os.Exit(1)
}

func parseDictionary() map[string][][]string {
	// dictionary format is roughly HIGHFLIER  HH AY1 F L AY2 ER0
	m := make(map[string][][]string)

	var filename = "deps/PronunciationDictionary.txt"
	fh, err := os.Open(filename)
	if err != nil {
		panic("Unable to open dictionary file")
	}
	defer fh.Close()

	scanner := bufio.NewScanner(fh)

	for scanner.Scan() {
		line := scanner.Text()
		elements := strings.Fields(line)
		key := elements[0]
		elements = elements[1:]
		m[key] = append(m[key], elements)
	}

	return m
}

func getTarget(target *Target, words *map[string][][]string) error {
	// start populating our target object
	target.populate()

	// check if our target exists in the dictionary so we have a base pronunciation to go off of
	target.checkValid(&(*words))

	// find our primary stress and it's index
	target.stressIndex, target.primaryStress = getStress(target.pronunciation)
	if target.stressIndex == -1 {
		panic("[!] Unable to determine the primary stress for your target")
	}

	return nil
}

func (target *Target) populate() error {
	target.word = os.Args[1]
	var match = regexp.MustCompile(`^[a-zA-Z]+$`)
	if !match.MatchString(target.word) {
		usage()
	}
	target.word = strings.ToUpper(target.word)
	return nil
}

func (target *Target) checkValid(words *map[string][][]string) error {
	p, ok := (*words)[target.word]
	if !ok {
		panic("Specified input not found in pronunciation dictionary")
	}
	if len(p) > 1 {
		// if our input has more than one possible pronunciation, query the user for which they meant
		target.pronunciation = checkDoubles(p)
	} else {
		target.pronunciation = p[0]
	}
	return nil
}

func checkDoubles(p [][]string) []string {
	var a []string
	for {
		fmt.Println("[*] The target specified has more than one pronuciation, please specify which one you'd like to use?")
		// display the list of pronunciations we found
		for i := 0; i < len(p); i++ {
			fmt.Printf("%d) %s\n", i+1, p[i])
		}
		fmt.Printf("> ")
		// get input from user
		reader := bufio.NewReader(os.Stdin)
		response, _ := reader.ReadString('\n')
		response = strings.TrimSpace(response)
		// make sure what the user gave us makes sense
		var match = regexp.MustCompile(`^[0-9]+$`)
		if !match.MatchString(response) {
			fmt.Println("[!] Invalid Input. Try again.")
			continue
		}
		x, _ := strconv.Atoi(response)
		if x > len(p) {
			fmt.Println("[!] Input out of range, please try again")
			continue
		}
		x--
		a = p[x]
		break
	}
	return a
}

func getStress(pronunciation []string) (int, string) {
	// the primary stress is always notated as <SOUND>1
	match := regexp.MustCompile(`^[A-Z]+1$`)
	for i, phoneme := range pronunciation {
		if match.MatchString(phoneme) {
			return i, phoneme
		}
	}
	return -1, "" // this is probably not the right way to do this
}

func findRhymes(target *Target, words *map[string][][]string) error {
	for key, value := range *words {
		// obviously our target will match itself, so we want to skip that
		if key == target.word {
			continue
		}
		// loop through the values in case there's more than one, many of the words in the dictionary have more than one pronunciation
		for _, phonemes := range value {
			if checkConditions(target, phonemes) {
				target.rhymes = append(target.rhymes, fmt.Sprintf("%s [%s]", key, strings.Join(phonemes, " ")))
			}
		}
	}
	return nil
}

func checkConditions(target *Target, phonemes []string) bool {
	index, stress := getStress(phonemes)
	// if the primary stresses do not match - condition 1
	if target.primaryStress != stress {
		return false
	}
	// if the subsequent phoneme length is different - condition 1
	if (len(target.pronunciation) - target.stressIndex) != (len(phonemes) - index) {
		return false
	}
	// if the subsequent sounds do not match - condition 1
	for i := 0; i < (len(target.pronunciation) - target.stressIndex); i++ {
		if target.pronunciation[target.stressIndex+i] != phonemes[index+i] {
			return false
		}
	}
	// if the stress is the start of the word - condition 2
	if target.stressIndex == 0 && index == 0 {
		return false
	}
	// if the stress is at the beginning of the word and not the other
	if (target.stressIndex == 0 && index != 0) || (target.stressIndex != 0 && index == 0) {
		return true
	}
	// if the proceeding phoneme matches - condition 2
	if target.pronunciation[target.stressIndex-1] != phonemes[index-1] {
		return true
	}
	return false
}

func outputResults(target *Target) error {
	if len(target.rhymes) == 0 {
		fmt.Println("[!] I found no rhymes that match the target")
	} else {
		fmt.Println("[*] I found the following rhymes:")
		for i, rhyme := range target.rhymes {
			fmt.Printf("%d) %s\n", i, rhyme)
		}
	}
	return nil
}
