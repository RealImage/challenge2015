package httpget

import (
	"encoding/json"
	"time"
	"net/http"
	"io/ioutil"
	moviebuffDatatype "degreesOfSeparation/datatype"
)

var TotalRequest uint

// Fetch Actor and Movie data from https://data.moviebuff.com/{moviebuff_url}
func FetchMoviebuffData(url string) (*moviebuffDatatype.MoviebuffRes, error) {
	time.Sleep(100 * time.Millisecond)

	resp, _ := http.Get(moviebuffDatatype.DataUri + url)
	defer resp.Body.Close()

	result, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var i moviebuffDatatype.MoviebuffRes
	err = json.Unmarshal(result, &i)
	if err != nil {
		return nil, err
	}
	TotalRequest++
	// print_response
	// fmt.Printf("%v\n\n", i)
	return &i, nil
}