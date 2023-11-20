package moviebuffclient

import (
	"degrees/datastructs"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

var c = &http.Client{
	Timeout: 30 * time.Second,
}

func MakeHttpReq(suffix string) (*http.Response, error) {
	httpUrl := fmt.Sprintf("http://data.moviebuff.com/%s", suffix)
	req, err := http.NewRequest("GET", httpUrl, nil)
	if err != nil {
		return nil, err
	}
	res, err := c.Do(req)
	return res, err
}

func GetPersonInfo(personUrl string) (*datastructs.Person, error) {
	res, err := MakeHttpReq(personUrl)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	var data datastructs.Person
	if res.StatusCode != 200 {
		return &data, nil
	}

	err = json.NewDecoder(res.Body).Decode(&data)
	return &data, err
}

func GetMovieInfo(movieUrl string) (*datastructs.Movie, error) {
	res, err := MakeHttpReq(movieUrl)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	var data datastructs.Movie
	if res.StatusCode != 200 {
		return &data, nil
	}
	err = json.NewDecoder(res.Body).Decode(&data)
	return &data, err
}
