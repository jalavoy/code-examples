// jalavoy - 03.05.2018
// this program plays a quick guessing game with the user. It picks a number between the globals defined below, and then prompts the user to guess.
package main

import "fmt"
import "time"
import "math/rand"
import "bufio"
import "os"
import "strconv"
import "strings"
import "sort"
import "regexp"

// stats struct
type Stats struct {
	games   int
	guesses int
	scores  []int
}

// globals for randomness and our statistics
var min = 1
var max = 100
var stats = &Stats{
	games:   0,
	guesses: 0,
}

func main() {
	// main game loop -- run until player says otherwise
	for {
		// run game
		doGame(stats)
		// query player to see if they want to keep playing -- if they don't, break the loop
		response := queryPlayer("Do you want to play again? ")
		var regex = regexp.MustCompile(`^(Y|y)`)
		if !regex.MatchString(response) {
			break
		}
	}
	// print statistics
	printStats(stats)
	// exit and return true to OS
	os.Exit(0)
}

func doGame(stats *Stats) bool {
	// predefine our games stats, guess is defined outside of the range of randomness so it wont trigger the loop on it's own
	var (
		guess      int64
		guesscount int
	)
	// generating our randomness target
	target := _random(min, max)
	fmt.Println("\nI'm thinking of a number between", min, "and", max)
	// while the players guess isn't correct
	for guess != target {
		guesscount++
		response := queryPlayer("Your guess? ")
		response = strings.TrimSpace(response)
		guess, _ = strconv.ParseInt(response, 10, 0)
		if guess == target {
			if guesscount == 1 {
				fmt.Println("You got it right in 1 guess!")
			} else {
				fmt.Println("You got it right in", guesscount, "guesses!")
			}
		} else {
			if guess < target {
				fmt.Println("Higher!")
			} else {
				fmt.Println("Lower!")
			}
		}
	}
	// update stats
	stats.recordStats(guesscount)
	return true
}

func queryPlayer(input string) string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print(input)
	response, _ := reader.ReadString('\n')
	return response
}

func (stats *Stats) recordStats(guesscount int) bool {
	stats.games++
	stats.guesses += guesscount
	stats.scores = append(stats.scores, guesscount)
	return true
}

func printStats(stats *Stats) bool {
	fmt.Println("\nOverall Results:")
	fmt.Println("Total Games:", stats.games)
	fmt.Println("Total Guesses:", stats.guesses)
	fmt.Println("Average guesses per game:", _getAvg(stats.scores))
	fmt.Println("Best game:", _getBest(stats.scores))
	return true
}

func _random(min, max int) int64 {
	rand.Seed(time.Now().Unix())
	return int64(rand.Intn(max-min) + min)
}

func _getAvg(scores []int) int {
	total := 0
	for _, value := range scores {
		total += value
	}
	avg := total / len(scores)
	return avg
}

func _getBest(scores []int) int {
	sort.Ints(scores)
	return scores[0]
}
