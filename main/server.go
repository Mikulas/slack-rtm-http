package main

import (
	"context"
	"fmt"
	"github.com/nlopes/slack"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
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

func (srv *Server) sendMessage(message string) error {
	var err error
	srv.respChan, srv.respTimestamp, err = srv.rtm.PostMessage(srv.channel, message, slack.PostMessageParameters{
		Username:   Username,
		AsUser:     true,
		EscapeText: false,
	})
	return err
}

func (srv *Server) updateMessage(message string) error {
	_, _, _, err := srv.rtm.SendMessageContext(context.Background(), srv.respChan, slack.MsgOptionUpdate(srv.respTimestamp), slack.MsgOptionText(message, false))
	return err
}

func (srv *Server) message(message string) error {
	if srv.respChan == "" {
		return srv.sendMessage(message)
	} else {
		return srv.updateMessage(message)
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
	srv.message(string(body))

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
