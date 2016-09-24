package main

import (
	"bufio"
	"compress/flate"
	"fmt"
	"net"
	"strings"
	"time"

	bp "github.com/nexustix/boilerplate"
)

/*
<;>|saction:uplink
<;>|saction:uplink<;>encrypt:true<;>compress:true
<;>|saction:subscribe<;>channel:caketalk
<;>|saction:broadcast<;>message:the cake is a lie
*/

//Remote represents (and handles) a [remote] connection
type Remote struct {
	Channels    []string          // cannels remote is listening to
	Buffer      []string          // message buffer
	ID          string            // ID/IP of remote
	Name        string            // Name of remote (if provided)
	seperator   string            //segment seperator
	doEncrypt   bool              // should transmission be encrypted
	doCompress  bool              // should transmission be compressed
	isAlive     bool              // is connection still alive
	isLegalized bool              // if remote is trusted
	server      *Server           // The server
	conn        *net.Conn         // connnection to remote
	remo        *bufio.ReadWriter // uncompressed read/write access
	remoc       *bufio.ReadWriter // uncompressed read/write access
}

//Run represents the main "Loop" of the Remote
func (r *Remote) Run(connection net.Conn) {
	defer connection.Close()
	r.ID = connection.RemoteAddr().String()
	r.conn = &connection

	tmpReader := bufio.NewReader(*r.conn)
	tmpWriter := bufio.NewWriter(*r.conn)
	r.remo = bufio.NewReadWriter(tmpReader, tmpWriter)

	tmpCRead := flate.NewReader(*r.conn)
	tmpCWrite, err := flate.NewWriter(*r.conn, -1)
	if bp.GotError(err) {
		fmt.Printf("[%v](%v)<-> WARNING FAILED to create compressed writer\n", time.Now().Unix(), r.ID)
	}
	tmpCReader := bufio.NewReader(tmpCRead)
	tmpCWriter := bufio.NewWriter(tmpCWrite)

	r.remoc = bufio.NewReadWriter(tmpCReader, tmpCWriter)

	fmt.Printf("[%v](%v)<-> NEW Connection\n", time.Now().Unix(), r.ID)
	r.Receive()
	for _, v := range r.Buffer {
		tmpMessage := Message{}
		tmpMessage.FromString(v)

		switch tmpMessage.Data["saction"] {
		case "uplink":
			if strings.ToLower(tmpMessage.Data["encrypt"]) == "true" {
				r.doEncrypt = true
			}
			if strings.ToLower(tmpMessage.Data["compress"]) == "true" {
				r.doCompress = true
			}
			r.isAlive = true
		}
	}

	if r.isAlive {
		fmt.Printf("[%v](%v)<-> SUCCESS Uplink\n", time.Now().Unix(), r.ID)
		fmt.Printf("[%v](%v)<-> <E>%v< <C>%v<\n", time.Now().Unix(), r.ID, r.doEncrypt, r.doCompress)
	} else {
		fmt.Printf("[%v](%v)<-> FAILED Uplink\n", time.Now().Unix(), r.ID)
	}

	//r.isAlive = true
	for r.isAlive {
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
		//HACK
		r.Buffer = []string{}
	}
	fmt.Printf("[%v](%v)<-> CLOSED Connection\n", time.Now().Unix(), r.ID)

}

//Receive receives one message (blocking) and adds it to the buffer
func (r *Remote) Receive() {
	var tmpBytes []byte
	var isPrefix bool
	var err error

	if r.doCompress {
		tmpBytes, isPrefix, err = r.remoc.ReadLine()
	} else {
		tmpBytes, isPrefix, err = r.remo.ReadLine()
	}

	tmpLine := string(tmpBytes)
	if bp.GotError(err) {
		fmt.Printf("[%v](%v)<!> ERR receiving\n", time.Now().Unix(), r.ID)
		r.isAlive = false
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

//Send sends a string via the connetion of the Remote
func (r *Remote) Send(message string) {
	if r.doCompress {
		r.remoc.WriteString(message + "\n")
		r.remoc.Flush()
	} else {
		r.remo.WriteString(message + "\n")
		r.remo.Flush()
	}
}

func (r *Remote) Broadcast(message string) {

}
