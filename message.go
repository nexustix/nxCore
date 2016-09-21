package main

import (
	"errors"
	"fmt"
	"strings"
)

type Message struct {
	Data map[string]string
}

func (m *Message) FromString(message string) error {
	if m.Data == nil {
		m.Data = make(map[string]string)
	}
	messageSegments := strings.Split(message, "<;>")
	if len(messageSegments) <= 0 {
		return errors.New("<!> No message salvaged from string")
	}
	for _, v := range messageSegments {
		fmt.Println(v)
		segmentSegments := strings.SplitN(v, ":", 2)
		if len(segmentSegments) >= 2 {
			m.Data[segmentSegments[0]] = segmentSegments[1]
		}
	}
	return nil
}

/*
func (m *Message) FromCompressedString(message) {

}
*/
