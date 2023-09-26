package util

import "regexp"

var validURLPattern = regexp.MustCompile(`^[a-z]+(?:-[a-z]+)*$`)

func IsValidURL(url string) bool {
    return validURLPattern.MatchString(url)
}