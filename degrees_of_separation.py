import requests
import json
import sys
import time
import os

finalList = []
visitedMovieList = []
visitedPersonList = []
degressOfSeperation = 0
try:
    os.remove('degressOfSeperation.txt')
except OSError:
    pass

degressOfSeperationFile = open(r'degressOfSeperation.txt','a')

def getData(urlName):
    url = 'http://data.moviebuff.com/{0}'.format(urlName)
    try:
        response = requests.get(url)
    except Exception, e:
        raise e
    else:
        return json.loads(response.text)

def validatePeople(peopleName):
    try:
        data = getData(peopleName)
    except ValueError:
        print "peopleName doesn't exist"
        return False
    else:
        if "movies" in data.keys():
            print "It is a People"
            return True

def actorMovie(actorMovieList): 
    degressOfSeperationFile.write("\nActor: "+ actorMovieList[0])
    degressOfSeperationFile.write("\nRole: "+ actorMovieList[1])
    degressOfSeperationFile.write("\nMovie: "+ actorMovieList[2])

    
def movieActor(movieActorList): 
    degressOfSeperationFile.write("\nMovie: "+ movieActorList[0])
    degressOfSeperationFile.write("\nRole: "+ movieActorList[1])
    degressOfSeperationFile.write("\nActor: "+ movieActorList[2])
    
def makeSubList(subList,eachList):
    if len(eachList)>=3:
        subList.append(eachList[0:3])
        makeSubList(subList,eachList[2:len(eachList)])
    degressOfSeperationFile.write("\n")
    return subList   

def printResult (globalList):
    for eachList in globalList: 
        subList = []
        degrees = makeSubList(subList,eachList)
        for i in range (0,len(degrees)): 
            if i%2 == 0: 
                actorMovie(degrees[i])
            i += 1
            if i%2 != 0: 
                movieActor(degrees[i])
            i += 1
        
def flatten(lst):
    for element in lst:
        if isinstance(element, list):
            for i in flatten(element):
                yield i
        else:
            yield element

            
def makeSingleList(finalList):
    globalList = []
    for lis in finalList:
        fL = []
        for x in list(flatten(lis[0])):
            fL.append(x)
        fL.append(lis[1])
        tempLists = list(flatten(lis[2]))
        tempLists.reverse()
        for tempList in tempLists:
            fL.append(tempList)
        globalList.append(fL)
    printResult(globalList)
    
    
def constructFinalList (degressOfSeperation,firstList,secondList,commonList): 
    finalList1 = []
    finalList2 = []
    for commonItem in commonList: 
        finalListforFirstPerson = next((item for item in firstList if commonItem in item[-1]), None) #picks the list which contain commonItem as last element of the list 
        finalListforSecondPerson = next((item for item in secondList if commonItem in item[-1]), None)
        finalList1.append(finalListforFirstPerson)
        finalList2.append(finalListforSecondPerson)
    for i in range(0,len(finalList1)): 
        final =[finalList1[i][0:-1],finalList1[i][-1],finalList2[i][0:-1]]
        finalList.append(final)
    makeSingleList(finalList)
    degressOfSeperationFile.write("\n")
    degressOfSeperationFile.write("\nDegrees of seperation is: " + str(degressOfSeperation))
    
    
def finalVisitedMovieList (movieListOfPerson): 
    actorListOfMovie = []
    for item in movieListOfPerson:
        if item[-1] not in visitedMovieList:
            visitedMovieList.append(item[-1])
            try: 
                MovieDetails = getData(item[-1])
            except: 
                continue 
            for cast in MovieDetails['cast']:
                if cast['url'] not in visitedPersonList:
                    actorListOfMovie.append([item,str(cast['role']),str(cast['url'])])
                else: 
                    continue 
            for crew in MovieDetails['crew']:
                if crew['url'] not in visitedPersonList: 
                    actorListOfMovie.append([item,str(crew['role']),str(crew['url'])])
                else: 
                    continue
        else: 
            continue 
    return actorListOfMovie

def finalVisitedPersonList(listOfPerson): 
    movieListOfPerson = []
    for item in listOfPerson:
        if item[-1] not in visitedPersonList:
            visitedPersonList.append(item[-1])
            try: 
                personDetails = getData(item[-1])
            except: 
                continue
            for movieDetails in personDetails['movies']:
                movieListOfPerson.append([item,str(movieDetails['role']),str(movieDetails['url'])])
        else: 
            continue
    return movieListOfPerson

def getDataUsingMovieList(degressOfSeperation,movieListofFirstPerson,movieListofSecondPerson):
    if degressOfSeperation < 7: 
        degressOfSeperation += 1
        actorListOfFirstMovie = finalVisitedMovieList (movieListofFirstPerson)
        actorListOfSecondMovie = finalVisitedMovieList (movieListofSecondPerson)
        commonActor = list(set(i[-1] for i in actorListOfFirstMovie) & set(i[-1] for i in actorListOfSecondMovie)) 
        if commonActor: 
            constructFinalList (degressOfSeperation,actorListOfFirstMovie,actorListOfSecondMovie,commonActor)
        else: 
            getDataUsingActorList(degressOfSeperation,actorListOfFirstMovie,actorListOfSecondMovie)
    else: 
        degressOfSeperationFile.write("Exceeds six degrees of Seperation")


def getDataUsingActorList(degressOfSeperation,firstListOfPerson,secondListOfPerson):
    if degressOfSeperation < 7:
        degressOfSeperation += 1
        movieListofFirstPerson = finalVisitedPersonList(firstListOfPerson)
        movieListofSecondPerson = finalVisitedPersonList(secondListOfPerson)
        commonMovies = list(set(i[-1] for i in movieListofFirstPerson) & set(i[-1] for i in movieListofSecondPerson))
        if commonMovies: 
            constructFinalList (degressOfSeperation,movieListofFirstPerson,movieListofSecondPerson,commonMovies)
        else: 
            getDataUsingMovieList(degressOfSeperation,movieListofFirstPerson,movieListofSecondPerson)
    else: 
        degressOfSeperationFile.write("Exceeds six degrees of Seperation")


if __name__ == "__main__":
    if len(sys.argv) != 3:
        print "Please input two people url"
    else:
        if not (validatePeople(sys.argv[1]) and validatePeople(sys.argv[2])):
            print "One of the url is not people"
        else:
            degressOfSeperation = 0
            degressOfSeperationFile.write("\nstartedTime:" + str(time.strftime("%I:%M:%S")))
            firstPersonList = [[sys.argv[1]]]
            secondPersonList = [[sys.argv[2]]]
            getDataUsingActorList(degressOfSeperation,firstPersonList,secondPersonList)
            degressOfSeperationFile.write("\nendTime:" +str(time.strftime("%I:%M:%S")))
            degressOfSeperationFile.write("\nVisited Person: " + str(len(visitedPersonList)))
            degressOfSeperationFile.write("\nvisitedMovieList: "+str(len(visitedMovieList)))
            print "Calculation for degressOfSeperation is done. Please look at file in degressOfSeperation.txt"

