package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

const SlackURL = "https://hooks.slack.com/services/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"

type SmsMessage struct {
	MessageAddress string `json:"msgAddr"`
	MessageBody    string `json:"msgBody"`
}

type SlackMessage struct {
	Text string `json:"text"`
}

func welcome(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprint(w, "Just remember, the Mailman doesn't deliver on Sundays\n")
}

func sendToSlack(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	decoder := json.NewDecoder(r.Body)
	var smsMsg SmsMessage
	err := decoder.Decode(&smsMsg)
	if err != nil {
		http.Error(w, "500 Internal Server Error: "+err.Error(), 500)
	}
	r.Body.Close()

	log.Println("SMS Received:", smsMsg.MessageAddress, "-", smsMsg.MessageBody)

	slackMsg := &SlackMessage{}
	slackMsg.Text = fmt.Sprintf("%s\n---------------\n%s\n",
		smsMsg.MessageAddress, smsMsg.MessageBody)
	jsonBytes, err := json.Marshal(slackMsg)
	if err != nil {
		http.Error(w, "500 Internal Server Error: "+err.Error(), 500)
	}

	log.Println("Send request to Slack app")

	res, err := http.Post(SlackURL, "application/json", bytes.NewBuffer(jsonBytes))
	if err != nil {
		http.Error(w, "500 Internal Server Error: "+err.Error(), 500)
	}

	log.Println("Response from Slack app:", res.StatusCode)

	fmt.Fprint(w, "ok")
}

func main() {
	router := httprouter.New()
	router.GET("/", welcome)
	router.POST("/sms", sendToSlack)

	log.Fatal(http.ListenAndServe(":8080", router))
}
