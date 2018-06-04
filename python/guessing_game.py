#!/usr/bin/env python3
# jalavoy - 02.27.2018
# this program plays a quick guessing game with the user. It picks a number between the globals defined below, and then prompts the user to guess.
import sys
import random
import re

# globals for randomness and our statistics struct (there aren't constants, python doesn't have constants)
MIN = 1
MAX = 100
STATS = {
    'games': 0,
    'guesses': 0,
    'scores': [],
}

def main():
    # our sweet haiku
    print("it's not dns\nit cannot be dns\nit was dns")
    # main game loop -- run until the player says otherwise
    while True:
        # run game, passing STATS struct
        doGame(STATS)
        # query player to see if they want to keep playing -- if they don't, break the loop
        if not re.match("^(Y|y)", queryPlayer("Do you want to play again? ")):
            break
    # print statistics
    printStats(STATS)
    # return true exit status to OS
    sys.exit(1)

def doGame(STATS):
    # predefining our games stats, guess is defined outside of the range of the randomness so it wont trigger on it's own
    guess = 0
    guesscount = 0
    # generating randomness target
    target = random.randint(MIN, MAX)
    print("\nI'm thinking of a number between {} and {}...".format(MIN, MAX))
    # while the players guess isn't correct -- we don't have to use continue when they get it right because this will evaluate to true automatically
    while int(guess) != target:
        # count guesses
        guesscount += 1
        # query the player for their guess
        guess = int(queryPlayer("Your guess? "))
        # if the guess is correct
        if guess == target:
            # if it is the first game
            if guesscount == 1:
                print("You got it right in {} guess!\n".format(guesscount)) 
            else:
                print("You got it right in {} guesses!\n".format(guesscount))
            # append our score to the scores list
            STATS['scores'].append(guesscount)
        else:
            # tell them they were wrong but which direction they need to go
            print("Higher!") if guess < target else print("Lower!")
    # update stats
    recordStats(STATS, guesscount)
    return 1

def queryPlayer(query):
    # query user for input
    response = input(query)
    return(response)

def recordStats(STATS, guesscount):
    # increment games by one
    STATS['games'] += 1
    # increment the guess count by how many we tried this game
    STATS['guesses'] += guesscount
    return 1

def printStats(STATS):
    print("\nOverall results:")
    print("Total Games: {}".format(STATS['games']))
    print("Total Guesses: {}".format(STATS['guesses']))
    print("Average guesses per game: {}".format(int(sum(STATS['scores']) / len(STATS['scores']))))
    print("Best Game: {}".format(sorted(STATS['scores'])[0]))
    return 1

if __name__ == '__main__':
    main()
