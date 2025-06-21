package main

import (
	"html/template"
	"net/http"
	"strconv"
	"sync"
	"time"
)

var _ http.Handler = &msgCacheHandler{}

type datedMessage struct {
	message   string
	timestamp time.Time
}

type msgCacheHandler struct {
	messages        []datedMessage
	messageDuration time.Duration
	mutex           sync.RWMutex
	tmpl            *template.Template
}

func (c *msgCacheHandler) Receive(msgChan <-chan string) {
	for msg := range msgChan {
		datedMsg := datedMessage{
			message:   msg,
			timestamp: time.Now(),
		}

		c.mutex.Lock()
		c.messages = append(c.messages, datedMsg)
		c.mutex.Unlock()
	}
}

func (c *msgCacheHandler) GetCurrentMessages() []string {
	c.CleanOldMessages()

	c.mutex.RLock()
	defer c.mutex.RUnlock()

	res := make([]string, len(c.messages))
	for _, dm := range c.messages {
		res = append(res, dm.message)
	}
	return res
}

func (c *msgCacheHandler) CleanOldMessages() {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	recentIndex := 0
	for i, dm := range c.messages {
		recentIndex = i
		if time.Since(dm.timestamp) < c.messageDuration {
			break
		}
	}
	c.messages = c.messages[recentIndex:]
}

func (c *msgCacheHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c.tmpl.Execute(w, c.GetCurrentMessages())
}

func startDisplayServer(msgChan <-chan string, port int64, refreshRate int64, tmplPath string) {
	cache := msgCacheHandler{
		messageDuration: time.Duration(refreshRate) * time.Second,
		tmpl:            template.Must(template.ParseFiles(tmplPath)),
	}

	go cache.Receive(msgChan)

	http.Handle("/", &cache)

	http.ListenAndServe(convertPort(port), nil)
}

func convertPort(port int64) string {
	if port <= 0 || port > 65535 {
		return ":8080"
	}

	return ":" + strconv.FormatInt(port, 10)
}
