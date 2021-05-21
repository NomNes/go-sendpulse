package go_sendpulse

import (
	"strings"
	"time"
)

type JsonDateTime time.Time

func (jdt *JsonDateTime) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), "\"")
	t, err := time.Parse("2006-01-02 15:04:05", s)
	if err != nil {
		return err
	}
	*jdt = JsonDateTime(t)
	return nil
}

type ListOptions struct {
	Limit  *int
	Offset *int
}

type ResultResponse struct {
	Result bool `json:"result"`
}
