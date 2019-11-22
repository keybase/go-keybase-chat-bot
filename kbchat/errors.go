package kbchat

import "fmt"

type APIError struct {
	err error
}

func (e APIError) Error() string {
	return fmt.Sprintf("failed to call keybase kvstore api: %v", e.err)
}

type UnmarshalError struct {
	err error
}

func (e UnmarshalError) Error() string {
	return fmt.Sprintf("failed to parse output from keybase kvstore api: %v", e.err)
}

type ResponseError struct {
	msg string
}

func (e ResponseError) Error() string {
	return fmt.Sprintf("received error from keybase kvstore api: %s", e.msg)
}
