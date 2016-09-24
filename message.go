package main

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

type Message struct {
	Data      map[string]string
	seperator string
}

func (m *Message) FromString(message string) error {
	if m.Data == nil {
		m.Data = make(map[string]string)
	}
	bodySegments := strings.SplitN(message, "|", 2)

	if len(bodySegments) >= 2 {
		m.seperator = bodySegments[0]
		messageSegments := strings.Split(bodySegments[1], m.seperator)
		if len(messageSegments) <= 0 {
			return errors.New("<!> No message salvaged from string")
		}
		for _, v := range messageSegments {
			//fmt.Println(v)
			segmentSegments := strings.SplitN(v, ":", 2)
			if len(segmentSegments) >= 2 {
				m.Data[segmentSegments[0]] = segmentSegments[1]
			}
		}
	} else {
		return errors.New("<!> Message format not recognized")
		fmt.Printf("[%v]<!> ERROR Message format not recognized\n", time.Now().Unix())
	}
	return nil

}

/*
func (m *Message) FromCompressedString(message) {

}
*/
