package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"net/http"
)

type connection struct {
	ws   *websocket.Conn
	sc   chan []byte
	data *Data
}

var wu = &websocket.Upgrader{ReadBufferSize: 512,
	WriteBufferSize: 512,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

func myWs(w http.ResponseWriter, r *http.Request) {
	ws, err := wu.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	c := &connection{sc: make(chan []byte, 256), ws: ws, data: &Data{}}
	h.r <- c
	go c.writer()
	c.reader()
	defer func() {
		c.data.Type = "logout"
		userList = del(userList, c.data.User)
		c.data.UserList = userList
		c.data.Content = c.data.User
		dataB, _ := json.Marshal(c.data)
		h.b <- dataB
		h.r <- c
	}()
}

func (c *connection) writer() {
	for message := range c.sc {
		c.ws.WriteMessage(websocket.TextMessage, message)
	}
	c.ws.Close()
}

var userList []string

func (c *connection) reader() {
	for {
		_, message, err := c.ws.ReadMessage()
		if err != nil {
			h.r <- c
			break
		}
		json.Unmarshal(message, &c.data)
		switch c.data.Type {
		case "login":
			c.data.User = c.data.Content
			c.data.From = c.data.User
			userList = append(userList, c.data.User)
			c.data.UserList = userList
			dataB, _ := json.Marshal(c.data)
			h.b <- dataB
		case "user":
			c.data.Type = "user"
			dataB, _ := json.Marshal(c.data)
			h.b <- dataB
		case "logout":
			c.data.Type = "logout"
			userList = del(userList, c.data.User)
			dataB, _ := json.Marshal(c.data)
			h.b <- dataB
			h.r <- c
		default:
			fmt.Print("========default================")
		}
	}
}

func del(slice []string, user string) []string {
	count := len(slice)
	if count == 0 {
		return slice
	}
	if count == 1 && slice[0] == user {
		return []string{}
	}
	var nSlice []string
	for i := range slice {
		if slice[i] == user && i == count {
			return slice[:count]
		} else if slice[i] == user {
			nSlice = append(slice[:i], slice[i+1:]...)
			break
		}
	}
	fmt.Println(nSlice)
	return nSlice
}
