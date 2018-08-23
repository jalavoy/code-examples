#!/usr/local/bin/python3
# jalavoy 08.23.2018
# parses a given file with the following format and displays the top 10 IP's seen in the log
# format: timestamp:user:ip
import sys
import re

data = {}
pattern = re.compile(r'^[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}$')

def main(file):
    with open(file) as fh:
        for line in fh:
            ip = line.split(':')[2].strip()
            if not pattern.match(ip):
                continue
            if ip in data:
                data[ip] += 1
            else:
                data[ip] = 1
    i = 1
    for ip, count in sorted(data.items(), key=lambda x: x[1], reverse=True):
        print("{}) {}: {}".format(i, ip, count))
        i += 1
        if i > 10:
            break

if __name__ == '__main__':
    main('deps/parse.txt')