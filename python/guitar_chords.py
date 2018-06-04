#!/usr/bin/env python3
# generates random guitar tabs from user input.
import sys
import random
import re
from pprint import pprint

STRINGS = ['e', 'B', 'G', 'D', 'A', 'E']
SOLONOTES = 24

def main():
    song = instantiateSong()
    parseFiles(song)
    getInteract(song)
    getSongTitle(song)
    getChordProgression(song)
    generateSolo(song)
    generateOutput(song)
    sys.exit(0)

def instantiateSong():
    # creates our data structure
    song = {
        'chords': {},
        'notes': {},
        'words': [],
    }
    return(song)

def parseFiles(song):
    # parses each of our 3 files and loads them into our song structure
    for category in ['chords', 'notes', 'words']:
        with open('deps/' + category + '.txt') as fh:
            for line in fh:
                # chords.txt and notes.txt have a different format than words.txt, so we have to parse them separately
                if category == 'chords' or category == 'notes':
                    elements = line.rsplit(' ')
                    # the first entry in the line is the key, so we need to shift that off and use it as the key in our dict
                    key = elements.pop(0)
                    song[category][key] = []
                    for element in elements:
                        element = element.strip()
                        song[category][key].append(element)
                else:
                    # this is the parsing for words.txt
                    song[category].append(line.strip())
        fh.close
    return(song)

def getInteract(song):
    _getKey(song)
    _getUnique(song)
    return(song)

def getSongTitle(song):
    # generates our random song title, idk extra credit
    song['title'] = []
    i = 0
    while i < 3:
        song['title'].append(random.choice(song['words']))
        i += 1
    return(song)

def getChordProgression(song):
    # generates the chord progression
    song['progression'] = []
    # while the length of the song['progression'] list is lower than the unique number of chords we want
    while len(song['progression']) < song['unique']:
        chord = random.choice(song['chords'][song['key']])
        # if we already have that chord, stop and generate a new one
        if chord in song['progression']:
            continue
        song['progression'].append(chord)
    return(song)

def generateSolo(song):
    i = 0
    song['solo'] = {}
    # setup the list of frets for each string
    for string in STRINGS:
        song['solo'][string] = []
    # stop once we finished our SOLONOTES value
    while i <= SOLONOTES:
        i += 1
        # get a random note, splitting the string and fret
        string, fret = _getNote(song)
        # loop each string in order
        for s in STRINGS:
            j = 0
            # the markup uses 3 -'s to denote a pause every beat, so we count that up to that so we fill with the right number of pauses
            while j < 3:
                j += 1
                # if the note we've generated belongs on this string
                if string == s:
                    # if our fret is higher than 10, we need to increment x by one to make sure we dont generate 4 characters for the beat
                    if isinstance(fret, int) and int(fret) >= 10:
                        j += 1
                    # add our note to the solo
                    song['solo'][s].append(fret)
                    # reset fret to x so we fill the remaining 1 or 2 slots with x, if we dont do this we'll just print the fret 3 times for this beat
                    fret = '-'
                else:
                    # if our generated note doesn't belong on this string, just put the pause
                    song['solo'][s].append('-')

def generateOutput(song):
    # this could also be written as:
    # print("\n\nSong Title: " + ' '.join(song['title']))
    print('\n\nSong Title: {}'.format(' '.join(song['title'])))
    print('Chord Progression: {}'.format(' '.join(song['progression'])))
    for string in STRINGS:
        print('{}|'.format(string), end='')
        for note in song['solo'][string]:
            print('{}'.format(note), end='')
        print()

# private functions
def _getKey(song):
    print('What key would you like to play in?\nOptions: {}'.format(str.join(', ', song['chords'].keys())))
    song['key'] = _queryUser('> ')
    while not song['key'] in song['chords'].keys():
        _getKey(song)
    return(song)
    
def _getUnique(song):
    print('How many unique chords? (1 through 6)')
    song['unique'] = int(_queryUser('> '))
    while not 1 <= song['unique'] <= 6:
        _getUnique(song)
    return(song)
        
def _queryUser(query):
    response = input(query)
    return(response)

def _getNote(song):
    note = random.choice(song['notes'][song['key']])
    r = re.match(r'([a-zA-Z])([0-9]+)', note)
    string = r.group(1)
    fret = int(r.group(2))
    return(string, fret)

if __name__ == '__main__':
    main()
