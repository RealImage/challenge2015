import requests, os
from pprint import pprint
from functools import reduce

"""
    Movies Data API (C) Real Image Media Technologies & MovieBuff 
"""

__author__ = "Karthik M A M ( https://www.github.com/KarthikMAM )"


def getData(url, more_data = False):
    """
    Download the data and store it in local storage 
    Also cache the processed data
    """
    if not hasattr(getData, 'cache'): getData.cache = {}

    def dwnlData(url):
        if os.path.isfile('data\\' + url + '.json'): 
            with open('data\\' + url + '.json', 'r') as file:
                data = eval(file.read())
        else:
            data = requests.get('http://data.moviebuff.com/' + url).json()

            with open('data\\' + url + '.json', 'w') as file:
                file.write(str(data))

        if not more_data:
            data = data['cast'] if data['type'] == 'Movie' else data['movies']
            
            data = set(list(map(lambda x: x['url'] , data)))

        return data

    return getData.cache.get(url, dwnlData(url))

def matchFunc(start, end):
    """
        Match function made as a tree and traversed like a BFS from two roots so as to reduce number of computations

        1. If it were a tree with single root the growth for 10 levels would be 1, 2, 4, 8, 16, 32, 64, 128, 256, 512 = "1023 Computations"
        2. By this solution we reduce the overall computations to 1 -> 2 -> 4 -> 8 -> 16 <- 8 <- 4 <- 2 <- 1 = "46 Computations"
        3. Also since there would be a large repitation of actors, this seems a viable solution to use a modified two way BFS using sets
    """
    extract = lambda x: sorted(list(x))[0]
    aggregate = lambda x: reduce(lambda x, y: x | getData(y), x, set())
    
    if len(start & end) > 0:
        return [extract(start & end)]
    else:
        next_start, next_end = aggregate(start), aggregate(end)

        if next_start == start and next_end == end:
            print("This scenario is invalid")
            exit(-1)

        result = matchFunc(next_start, next_end)

        match_start, match_end = extract(start & getData(result[0])), extract(end & getData(result[-1]))

        return [match_start] + result + [match_end]

if __name__ == "__main__":
    if int(input()) == 1:
        pprint(matchFunc({input()}, {input()}))
    else:
        pprint(getData(input(), True))

"""
Sample Input:
    1. If you want to find the minimum degree of seperation
        1
        vijay
        suriya-sivakumar
    2. If you want to find the details of a movie / actor
        2
        vijay
"""