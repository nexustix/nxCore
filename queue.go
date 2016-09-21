package main

import "errors"

//XXX Threadsafe ?
type Queue struct {
	data    []string
	topItem string
}

func (q *Queue) Push(item string) {
	q.data = append(q.data, item)
}

func (q *Queue) Pop() (string, error) {
	if len(q.data) <= 0 {
		return "", errors.New("<!> Queue empty")
	} else {
		q.topItem = q.data[0]
		q.data = q.data[1:]
	}
	return q.topItem, nil
}
