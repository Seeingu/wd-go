package main

import (
	"fmt"
	"github.com/tmaxmax/go-sse"
	"net/http"
	"time"
)

const (
	topicRandomNumbers = "numbers"
	topicMetrics       = "metrics"
)

var sseHandler = &sse.Server{
	Provider: &sse.Joe{
		ReplayProvider: &sse.ValidReplayProvider{
			TTL:        time.Minute * 5,
			GCInterval: time.Minute,
			AutoIDs:    true,
		},
	},
	OnSession: func(s *sse.Session) (sse.Subscription, bool) {
		topics := s.Req.URL.Query()["topic"]
		for _, topic := range topics {
			if topic != topicRandomNumbers && topic != topicMetrics {
				fmt.Fprintf(s.Res, "invalid topic %q; supported are %q, %q", topic, topicRandomNumbers, topicMetrics)
				s.Res.WriteHeader(http.StatusBadRequest)
				return sse.Subscription{}, false
			}
		}
		if len(topics) == 0 {
			// Provide default topics, if none are given.
			topics = []string{topicRandomNumbers, topicMetrics}
		}

		return sse.Subscription{
			Client:      s,
			LastEventID: s.LastEventID,
			Topics:      append(topics, sse.DefaultTopic), // the shutdown message is sent on the default topic
		}, true
	},
}

func SSEReload() {
	m := &sse.Message{}
	m.AppendData("__RELOAD__")
	sseHandler.Publish(m)
}
