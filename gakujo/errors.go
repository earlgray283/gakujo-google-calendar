package gakujo

import "net/http"

type ErrUnexpectedStatus struct {
	code   int
	expect int
}

func (e ErrUnexpectedStatus) Error() string {
	return "unexpected status code: " + http.StatusText(e.code) + ", expected: " + http.StatusText(e.expect)
}
