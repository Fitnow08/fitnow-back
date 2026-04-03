package utils

import "errors"

var (
	EventSavedCandle = "SavedCandle"
	EventSavedColor  = "SavedColor"
	EventSavedOrder  = "SavedOrder"
)

var (
	ErrorCreateQueryString = errors.New("Error creating query string")
	ErrorNotFoundRows      = errors.New("Error finding rows")
)
