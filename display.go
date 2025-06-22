package main

import (
	"html/template"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
)

var _ http.Handler = &msgCacheHandler{}

type datedMessage struct {
	message   string
	timestamp time.Time
}

type msgCacheHandler struct {
	messages []datedMessage
	mutex    sync.RWMutex

	messageDuration time.Duration
	messageFilter   func(string) string
	tmpl            *template.Template
}

func (c *msgCacheHandler) Receive(msgChan <-chan string) {
	for msg := range msgChan {
		datedMsg := datedMessage{
			message:   c.messageFilter(msg),
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
	for i, dm := range c.messages {
		res[i] = dm.message
	}
	return res
}

func (c *msgCacheHandler) CleanOldMessages() {
	now := time.Now()
	c.mutex.Lock()
	defer c.mutex.Unlock()

	recentIndex := len(c.messages) // when no messages are recent, clean all
	for i, dm := range c.messages {
		if now.Sub(dm.timestamp) < c.messageDuration {
			recentIndex = i
			break
		}
	}
	c.messages = c.messages[recentIndex:]
}

func (c *msgCacheHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c.tmpl.Execute(w, c.GetCurrentMessages())
}

func startDisplayServer(msgChan <-chan string, conf *Config) {
	msgHandler := msgCacheHandler{
		messageDuration: time.Duration(conf.refreshRate) * time.Second,
		messageFilter:   filterFromConf(conf),
		tmpl:            template.Must(template.ParseFiles(conf.tmplPath)),
	}

	go msgHandler.Receive(msgChan)

	http.Handle("/", &msgHandler)

	http.ListenAndServe(convertPort(conf.port), nil)
}

func convertPort(port int) string {
	return ":" + strconv.Itoa(port)
}

func filterFromConf(conf *Config) func(string) string {
	cutUntil := conf.cutUntil
	if cutUntil == "" {
		return noFilter
	}

	return func(msg string) string {
		if _, after, found := strings.Cut(msg, cutUntil); found {
			return after
		}
		return msg
	}
}

func noFilter(msg string) string {
	return msg
}
