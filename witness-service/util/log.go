// Copyright (c) 2019 IoTeX
// This is an alpha (internal) release and is not suitable for production. This source code is provided 'as is' and no
// warranties are given as to title or non-infringement, merchantability or fitness for purpose and, to the extent
// permitted by law, all liability for your use of the code is disclaimed. This source code is governed by Apache
// License 2.0 that can be found in the LICENSE file.

package util

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"time"

	lru "github.com/hashicorp/golang-lru"
)

// Payload is the message struct for slack
type Payload struct {
	Text string `json:"text,omitempty"`
}

var (
	larkURL  string
	slackURL string
	prefix   string
	cache    *lru.Cache
)

func init() {
	var err error
	cache, err = lru.New(10)
	if err != nil {
		panic(err)
	}
}

// SetLarkURL sets the lark post url
func SetLarkURL(url string) {
	larkURL = url
}

// SetSlackURL sets the slack post url
func SetSlackURL(url string) {
	slackURL = url
}

func SetPrefix(s string) {
	prefix = s
}

// Alert sends alert to
func Alert(msg string) {
	value, ok := cache.Get(msg)
	if ok {
		ts, ok := value.(time.Time)
		if !ok {
			panic(value)
		}
		if ts.After(time.Now().Add(-time.Hour)) {
			return
		}
	}
	cache.Add(msg, time.Now())
	SendSlackAlert(msg)
	SendLarkAlert(msg)
}

func SendLarkAlert(msg string) {
	if larkURL == "" {
		return
	}
	if prefix != "" {
		msg = prefix + ":" + msg
	}
	msgBytes, err := json.Marshal(struct {
		MsgType string  `json:"msg_type"`
		Content Payload `json:"content"`
	}{
		MsgType: "text",
		Content: Payload{Text: msg},
	})
	if err != nil {
		log.Printf("failed to construct lark message %+v\n", err)
		return
	}
	_, err = http.Post(larkURL, "application/json", bytes.NewReader(msgBytes))
	if err != nil {
		log.Printf("failed to send lark message %+v", err)
	}
}

func SendSlackAlert(msg string) {
	if slackURL == "" {
		return
	}
	if prefix != "" {
		msg = prefix + ":" + msg
	}
	msgBytes, err := json.Marshal(Payload{Text: msg})
	if err != nil {
		log.Printf("failed to construct slack message %+v\n", err)
		return
	}
	_, err = http.Post(slackURL, "application/json", bytes.NewReader(msgBytes))
	if err != nil {
		log.Printf("failed to send slack message %+v", err)
	}
}

// LogErr logs error
func LogErr(err error) {
	log.Println(err)
	Alert(err.Error())
}
