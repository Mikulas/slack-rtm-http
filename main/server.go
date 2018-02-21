package main

import (
	"context"
	"fmt"
	"github.com/nlopes/slack"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
	"encoding/json"
)

const (
	Username = "Deploy"
)

type Server struct {
	port    int
	channel string
	mutex   *sync.Mutex

	rtm                     *slack.RTM
	respChan, respTimestamp string
}

func MakeServer(port int, channel string, rtm *slack.RTM) *Server {
	srv := &Server{
		port:    port,
		channel: channel,
		rtm:     rtm,
		mutex:   &sync.Mutex{},
	}
	return srv
}

func (srv *Server) sendMessage(payload Payload) error {
	var err error
	srv.respChan, srv.respTimestamp, err = srv.rtm.PostMessage(srv.channel, payload.Message, slack.PostMessageParameters{
		Username:   Username,
		AsUser:     true,
		EscapeText: false,
	})
	return err
}

func (srv *Server) updateMessage(payload Payload) error {
	options := make([]slack.MsgOption, 0)

	parse := slack.MsgOptionPostMessageParameters(slack.PostMessageParameters{
		LinkNames: 1,
		UnfurlMedia: true,
		UnfurlLinks: true,
	})
	options = append(options, parse)

	update := slack.MsgOptionUpdate(srv.respTimestamp)
	options = append(options, update)

	text := slack.MsgOptionText(payload.Message, false)
	options = append(options, text)

	if payload.RollbackValue != "" {
		btn := slack.MsgOptionAttachments(slack.Attachment{
			CallbackID: "deploy",
			Color:      "#FFFFFF",
			Actions: []slack.AttachmentAction{{
				Name:  "rollback",
				Text:  "Rollback to this version",
				Style: "default",
				Type:  "button",
				Value: payload.RollbackValue,
				Confirm: &slack.ConfirmationField{
					Title:       "Are you sure?",
					Text:        "Migrations will NOT be rolled back!",
					OkText:      "Yes, rollback",
					DismissText: "Cancel",
				},
			}},
		})
		options = append(options, btn)
	}

	_, _, _, err := srv.rtm.SendMessageContext(context.Background(), srv.respChan, options...)
	return err
}

func (srv *Server) message(payload Payload) error {
	if srv.respChan == "" {
		log.Println("sending new message")
		return srv.sendMessage(payload)
	} else {
		log.Println("updating message")
		return srv.updateMessage(payload)
	}
}

func (srv *Server) handler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(405)
		fmt.Fprintln(w, "Invalid method")
		return
	}

	defer r.Body.Close()
	body, _ := ioutil.ReadAll(r.Body)

	payload := Payload{}
	err := json.Unmarshal(body, &payload)
	if err != nil {
		w.WriteHeader(500)
		log.Println(err.Error())
		return
	}

	err = srv.message(payload)
	if err != nil {
		w.WriteHeader(500)
		log.Println(err.Error())
		return
	}

	fmt.Fprintln(w, string(body))
}

func (srv *Server) lockingHandler(w http.ResponseWriter, r *http.Request) {
	srv.mutex.Lock()
	srv.handler(w, r)
	srv.mutex.Unlock()
}

func (srv *Server) Start() {
	log.Printf("listening on :%v\n", srv.port)
	http.HandleFunc("/", srv.lockingHandler)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", srv.port), nil))
}
