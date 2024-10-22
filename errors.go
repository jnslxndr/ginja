package ginja

import "fmt"

/* From: jsonapi.org:
An error object MAY have the following members:

    id: a unique identifier for this particular occurrence of the problem.
    links: a links object containing the following members:
        about: a link that leads to further details about this particular occurrence of the problem.
    status: the HTTP status code applicable to this problem, expressed as a string value.
    code: an application-specific error code, expressed as a string value.
    title: a short, human-readable summary of the problem that SHOULD NOT change from occurrence to occurrence of the problem, except for purposes of localization.
    detail: a human-readable explanation specific to this occurrence of the problem.
    source: an object containing references to the source of the error, optionally including any of the following members:
        pointer: a JSON Pointer [RFC6901] to the associated entity in the request document [e.g. "/data" for a primary data object, or "/data/attributes/title" for a specific attribute].
        parameter: a string indicating which query parameter caused the error.
    meta: a meta object containing non-standard meta-information about the error.

*/

type Error struct {
	Id      string `json:"id,omitempty"`
	Status  int    `json:"status,omitempty"`
	Code    string `json:"code,omitempty"`
	Title   string `json:"title,omitempty"`
	Details string `json:"detail,omitempty"`
}

func (e Error) Error() string {
	return e.Title
}

func NewError(err interface{}) Error {
	var title string
	switch err.(type) {
	case error:
		title = err.(error).Error()
	case fmt.Stringer:
		title = err.(fmt.Stringer).String()
	case string:
		title = err.(string)
	default:
		title = "Unknown error occurred"
	}
	return Error{
		Title: title,
	}
}
