package errors

import (
	"fmt"
	"net/http"
)

// Internal Error Code
type CustomError struct {
	InternalCode string
	HttpCode     int
	Message      string
	Resolution   string
}

var (
	//Resource Not Found
	UnableToBindResourceBody = CustomError{
		InternalCode: "100.100",
		HttpCode:     http.StatusBadRequest,
		Message:      "Failed reading the request body",
		Resolution:   "Check Post Resource format",
	}
)

func PackError(errorStruct *CustomError, err ...error) string {
	if len(err) != 0 {
		if len(errorStruct.Resolution) == 0 {
			return fmt.Sprintf("%s : %v.", errorStruct.Message, err[0])
		}
		return fmt.Sprintf("%s : %v. Possible Solution: %s ", errorStruct.Message, err[0], errorStruct.Resolution)
	} else {
		if len(errorStruct.Resolution) == 0 {
			return fmt.Sprintf(errorStruct.Message)
		}
		return fmt.Sprintf("%s. Possible Solution: %s", errorStruct.Message, errorStruct.Resolution)
	}
}
