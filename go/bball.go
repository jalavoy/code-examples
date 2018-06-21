// jalavoy - 06.20.2018
// this is another assignment borrowed from the University of Arizona
// The purpose of this assignment is to practice using classes and instances, and use them to produce each conferences win ratio avg
// some of the requirements of the assignment don't make a lot of practical sense, but I tried to follow them anyway. Any less-than-ideal behavior in this program is likely due to those contraints
// input looks something like this
/*
	# School (Conference)	Wins	Losses
	UConn (AAC)	38		0
	Baylor (Big 12)	36		2
	Notre Dame (Atlantic Coast)	33		2
	South Carolina (Southeastern)	33		2
	Colorado St. (Mountain West)	31		2
	Army West Point (Patriot)	29		3
	Maryland (Big Ten)	31		4
*/
package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

const inputFile string = "deps/2015-16.txt"

// Team object
type Team struct {
	Name       string
	Conference string
	Wins       float64
	Losses     float64
}

// Conference object
type Conference struct {
	Name  string
	Teams []Team
}

// ConferenceSet object
type ConferenceSet struct {
	Conferences []Conference
}

// custom sort type
type byWinRatio []Conference

// custom sort functions
func (c byWinRatio) Len() int           { return len(c) }
func (c byWinRatio) Swap(i, j int)      { c[i], c[j] = c[j], c[i] }
func (c byWinRatio) Less(i, j int) bool { return c[i].winRatioAvg() > c[j].winRatioAvg() }

func main() {
	fh, err := os.Open(inputFile)
	if err != nil {
		log.Panic("Unable to open input file")
	}
	var regexp = regexp.MustCompile(`^(.*)\s+\((.*)\)\s+(\d+)\s+(\d+)$`)
	scanner := bufio.NewScanner(fh)
	cs := &ConferenceSet{}
	for scanner.Scan() {
		// skip comments
		if strings.HasPrefix(scanner.Text(), "#") {
			continue
		}
		// pull our fields out of each line
		groups := regexp.FindStringSubmatch(scanner.Text())
		teamName := groups[1]
		conferenceName := groups[2]
		teamWins, _ := strconv.ParseFloat(groups[3], 64)
		teamLosses, _ := strconv.ParseFloat(groups[4], 64)
		// instantiate our team object
		t := &Team{
			teamName,
			conferenceName,
			teamWins,
			teamLosses,
		}
		// instantiate our Conference object, by pulling it from ConferenceSet if it already exists, and making a new instance if not
		c := cs.getConference(conferenceName)
		// add the Team object to the Conference object
		c.add(*t)
		// add the Conference object to the ConferenceSet object
		cs.add(*c)
	}
	// output the top 10
	for i, c := range cs.best() {
		fmt.Printf("%d) %s : %f\n", i+1, c.name(), c.winRatioAvg())
		if i > 10 {
			break
		}
	}

}

// Team methods
func (t *Team) name() string {
	return t.Name
}

func (t *Team) conf() string {
	return t.Conference
}

func (t *Team) winRatio() float64 {
	return (t.Wins / (t.Wins + t.Losses))
}

// Conference methods
func (c *Conference) name() string {
	return c.Name
}

func (c *Conference) add(t Team) error {
	c.Teams = append(c.Teams, t)
	return nil
}

func (c *Conference) winRatioAvg() float64 {
	var total float64
	var i float64
	for _, t := range c.Teams {
		total += t.winRatio()
		i++
	}
	return (total / i)
}

// ConferenceSet methods
func (cs *ConferenceSet) add(c Conference) {
	// determine if the Conference already exists in the set,
	// if it does: replace the old data with the new updated data
	// if it doesn't: append the data to the Conferences slice
	found := -1
	for i, each := range cs.Conferences {
		if each.Name == c.Name {
			found = i
			break
		}
	}
	if found >= 0 {
		cs.Conferences[found] = c
	} else {
		cs.Conferences = append(cs.Conferences, c)
	}
}

func (cs *ConferenceSet) getConference(conferenceName string) *Conference {
	// if the Conference already exists, return it,
	// else, instantiate a new one and return that
	for _, each := range cs.Conferences {
		if each.Name == conferenceName {
			return &each
		}
	}
	c := &Conference{}
	c.Name = conferenceName
	return c
}

func (cs *ConferenceSet) best() []Conference {
	// sorts the ConferenceSet by win ratios, descending
	sort.Sort(byWinRatio(cs.Conferences))
	return cs.Conferences
}
