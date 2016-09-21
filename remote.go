package main

import (
	"bufio"
	"fmt"
	"net"
	"time"

	bp "github.com/nexustix/boilerplate"
)

/*
saction:subscribe<;>channel:caketalk
saction:broadcast<;>message:the cake is a lie
*/

type Remote struct {
	Channels  []string
	Buffer    []string
	ID        string
	Name      string
	Legalized bool
	alive     bool
	server    *Server
	conn      *net.Conn
	remo      *bufio.ReadWriter
}

func (r *Remote) Run(connection net.Conn) {
	defer connection.Close()
	r.ID = connection.RemoteAddr().String()
	r.conn = &connection
	tmpReader := bufio.NewReader(*r.conn)
	tmpWriter := bufio.NewWriter(*r.conn)
	//tmpReadWriter := bufio.NewReadWriter(tmpReader, tmpWriter)
	r.remo = bufio.NewReadWriter(tmpReader, tmpWriter)
	r.alive = true
	fmt.Printf("[%v](%v)<-> NEW Connection\n", time.Now().Unix(), r.ID)
	for r.alive {
		r.Receive()

		for _, v := range r.Buffer {
			tmpMessage := Message{}
			tmpMessage.FromString(v)

			switch tmpMessage.Data["saction"] {
			case "subscribe":
				if tmpMessage.Data["channel"] != "" {
					r.Channels = append(r.Channels, tmpMessage.Data["channel"])
					r.Channels = bp.EliminateDuplicates(r.Channels)
					fmt.Printf("[%v](%v)<-> Subscribed %v\n", time.Now().Unix(), r.ID, tmpMessage.Data["channel"])
				}
			case "unsubscribe":
				if tmpMessage.Data["channel"] != "" {
					r.Channels = bp.EliminateStringInSlice(r.Channels, tmpMessage.Data["channel"])
					fmt.Printf("[%v](%v)<-> UN-Subscribed %v\n", time.Now().Unix(), r.ID, tmpMessage.Data["channel"])
				}
			case "broadcast":
				if tmpMessage.Data["message"] != "" {
					fmt.Printf("[%v](%v)<-> Broadcasting >%v<\n", time.Now().Unix(), r.ID, tmpMessage.Data["message"])
					r.server.Bridgecast(r, v)
				}
			}
		}
	}
	fmt.Printf("[%v](%v)<-> CLOSED Connection\n", time.Now().Unix(), r.ID)

}

func (r *Remote) Receive() {
	tmpBytes, isPrefix, err := r.remo.ReadLine()
	tmpLine := string(tmpBytes)
	if bp.GotError(err) {
		fmt.Printf("[%v](%v)<!> ERR receiving\n", time.Now().Unix(), r.ID)
		r.alive = false
		return
	}
	if (tmpLine != "") && (tmpLine != "nil") {
		if isPrefix {
			fmt.Printf("[%v](%v)<~> WARNING long line sent by client\n", time.Now().Unix(), r.ID)
		} else {
			fmt.Printf("[%v](%v)</> >%v<\n", time.Now().Unix(), r.ID, tmpLine)
			r.Buffer = append(r.Buffer, tmpLine)
		}
	}
}

func (r *Remote) Send(message string) {
	r.remo.WriteString(message + "\n")
	r.remo.Flush()
}

func (r *Remote) Broadcast(message string) {

}
