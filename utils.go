package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

func ParseURLsFile(filename string, bytes []byte) (*URL, error) {
	if filename != "" {
		file, err := os.Open(filename)
		if err != nil {
			log.Fatalf("Error in opening urls file. err : %+v", err)
			return nil, err
		}
		defer file.Close()

		bytes, err = io.ReadAll(file)
		if err != nil {
			log.Fatalf("Error in reading urls file. err : %+v", err)
			return nil, err
		}
	}

	urls := &URL{}
	err := json.Unmarshal(bytes, urls)
	if err != nil {
		log.Fatalf("Error in unmarshaling urls file. buff : %+s err : %+v", bytes, err)
		return nil, err
	}
	return urls, nil
}

func FetchData(url string) ([]byte, error) {
	ctx, cancelFunc := context.WithTimeout(context.Background(), time.Duration(3)*time.Second)
	defer cancelFunc()

	request, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		log.Fatalf("Error in creating request object with timeout. : err : %+v", err)
		return nil, err
	}
	reponse, err := http.DefaultClient.Do(request)
	if err != nil {
		log.Fatalf("Error in gettting data from url. : err : %+v", err)
		return nil, err
	}
	bb, err := io.ReadAll(reponse.Body)
	if err != nil {
		log.Fatalf("Error in reading data from response body. : err : %+v", err)
		return nil, err
	}
	return bb, nil
}

func FetchActorsData(urls []string) ([]*ActorData, error) {
	result := make([]*ActorData, 0)
	for _, url := range urls {
		bb, err := FetchData(url)
		if err != nil {
			return nil, err
		}
		data := &ActorData{}
		err = json.Unmarshal(bb, data)
		if err != nil {
			log.Fatalf("Error in unmarshaling data from url. buff : %+s err : %+v", bb, err)
			return nil, err
		}
		result = append(result, data)
	}
	return result, nil
}

func FetchMoviesData(urls []string) ([]*MoiveData, error) {
	result := make([]*MoiveData, 0)
	for _, url := range urls {
		bb, err := FetchData(url)
		if err != nil {
			return nil, err
		}
		data := &MoiveData{}
		err = json.Unmarshal(bb, data)
		if err != nil {
			log.Fatalf("Error in unmarshaling data from url. buff : %+s err : %+v", bb, err)
			return nil, err
		}
		result = append(result, data)
	}
	return result, nil
}

/*
[actor --> {{movie:role}, {movie:role}, {movie:role}}, actor --> {{movie:role}, {movie:role}, {movie:role}}]
We are using maps so the finding whether a particular actor belongs to a movie will take O(1)
*/
func MakeActorMap(ac []*ActorData, mv []*MoiveData) map[string]*PairSet {
	var result = make(map[string]*PairSet)
	for _, actor := range ac {
		result[actor.Name] = New()
		for _, movie := range actor.Movies {
			result[actor.Name].Insert(movie.Name, movie.Role)
		}
	}

	for _, movie := range mv {
		for _, cast := range movie.Casts {
			_, ok := result[cast.Name]
			if !ok {
				result[cast.Name] = New()
			}
			result[cast.Name].Insert(movie.Name, cast.Role)
		}
	}
	return result
}

/*
[movie --> {{cast:role}, {cast:role}, {cast:role}}, movie --> {{cast:role}, {cast:role}, {cast:role}}]
We are using maps so the finding whether a particular cast belongs to a movie will take O(1)
*/
func MakeMovieMap(ac []*ActorData, mv []*MoiveData) map[string]*PairSet {
	var result = make(map[string]*PairSet)
	for _, movie := range mv {
		result[movie.Name] = New()
		for _, cast := range movie.Casts {
			result[movie.Name].Insert(cast.Name, cast.Role)
		}
	}
	for _, actor := range ac {
		for _, movie := range actor.Movies {
			_, ok := result[movie.Name]
			if !ok {
				result[movie.Name] = New()
			}
			result[movie.Name].Insert(actor.Name, movie.Role)
		}
	}
	return result
}

func MakeAdjecencyList(movieMap, actorMap map[string]*PairSet) map[string]*PairSet {
	var graph = make(map[string]*PairSet)
	for actor, setOfMovies := range actorMap {
		graph[actor] = New()
		allMovies := setOfMovies.GetAllKeys()
		for _, singleMovie := range allMovies {
			casts := movieMap[singleMovie].GetAllKeys()
			for _, cast := range casts {
				graph[actor].Insert(cast, singleMovie)
			}
		}
	}
	return graph
}

func CleanArgs(arg string) string {
	s := strings.Split(arg, "-")
	for i, v := range s {
		s[i] = string(v[0]-32) + v[1:] // Making first letter of string capital.
	}
	return strings.Join(s, " ")
}

func GetCmdLineArgs() []string {
	args := os.Args
	if len(args) != 3 {
		log.Fatalf("You have provided invalid number of arguments.")
	}
	args[1] = CleanArgs(args[1])
	args[2] = CleanArgs(args[2])
	return args[1:]
}

func Reverse(arr []MovieNameCastName) []MovieNameCastName {
	var rev []MovieNameCastName
	for i := len(arr) - 1; i >= 0; i-- {
		rev = append(rev, arr[i])
	}
	return rev
}

func ReverseGenric[T any](arr []T) []T {
	var rev []T
	for i := len(arr) - 1; i >= 0; i-- {
		rev = append(rev, arr[i])
	}
	return rev
}

func PrintResult(result []MovieNameCastName, actorData, movieData map[string]*PairSet, n1, n2 string) {
	fmt.Println("Degrees of Separation: ", len(result))
	fmt.Println()
	for i, v := range result {
		fmt.Printf("%d. %s\n", i+1, v.Movie)
		fmt.Printf("%s: %s\n", actorData[v.Cast1].Get(v.Movie), v.Cast1) // Role_of_cast:Name_of_cast
		fmt.Printf("%s: %s\n", actorData[v.Cast2].Get(v.Movie), v.Cast2) // Role_of_cast:Name_of_cast
		fmt.Println()
	}
}
