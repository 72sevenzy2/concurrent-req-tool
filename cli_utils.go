package main

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

func (h *HeaderFlags) String() string { // gets called internally by go's flag pkg, (type flag.Value expects a String() and Set() func)
	return fmt.Sprint(*h)
}

func (h *HeaderFlags) Set(value string) error {
	*h = append(*h, value)
	return nil
}

func Validate(args []string) error {
	if len(args) < 1 {
		UsageMsg := InvalidInputErr
		return UsageMsg
	}
	return nil
}

// small helper func to return integer pointers.
func IntPtr(i int) *int { return &i } // returns *int

// func for adding headers
func AddHeaders(req *http.Request, args HeaderFlags) error {
	for _, h := range args {
		parts := strings.SplitN(h, ":", 2)
		if len(parts) != 2 {
			return InvalidInputErr
		}

		// appending headers
		req.Header.Set(strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1]))
	}
	return nil
}

// normalize key types to string (for storage.go)
func Normalize(keyname any) string {
	switch v := keyname.(type) {
	case int:
		return strconv.Itoa(v)
	case string:
		return v
	default:
		return ""
	}
}
