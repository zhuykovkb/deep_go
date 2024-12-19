package main

import (
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// go test -v homework_test.go

type MultiError struct {
	errors []error
}

func (e *MultiError) Error() string {
	var b strings.Builder
	_, err := fmt.Fprintf(&b, "%d errors occured:\n", len(e.errors))
	if err != nil {
		return ""
	}

	for _, err := range e.errors {
		b.WriteString("\t* ")
		b.WriteString(err.Error())
	}
	b.WriteString("\n")

	return b.String()
}

func Append(err error, errs ...error) *MultiError {
	mErr, ok := err.(*MultiError)

	if !ok {
		mErr = &MultiError{}
		if err != nil {
			mErr.errors = append(mErr.errors, err)
		}
	}

	for _, er := range errs {
		mErr.errors = append(mErr.errors, er)
	}

	return mErr
}

func TestMultiError(t *testing.T) {
	var err error
	err = Append(err, errors.New("error 1"))
	err = Append(err, errors.New("error 2"))

	expectedMessage := "2 errors occured:\n\t* error 1\t* error 2\n"
	assert.EqualError(t, err, expectedMessage)
}
