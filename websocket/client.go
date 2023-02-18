// Package websocket
/*
Copyright © 2023 runtimeracer@gmail.com

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package websocket

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"nhooyr.io/websocket"
	"strconv"
	"sync"
	"sync/atomic"
	"time"
)

// MessageHandlerFunc describes a function interface to be bound to the websocket client
type MessageHandlerFunc func(message *KajiwotoWebSocketMessage) error

// MessageHandler is used to bind handler functions to the websocket client
type MessageHandler struct {
	handlerKey      string // uuid to identify handler, internal
	handleFunc      MessageHandlerFunc
	removeOnSuccess bool
}

// KajiwotoWebSocketClient is a custom websocket client for kajiwoto reqeusts using the websocket API
type KajiwotoWebSocketClient struct {
	// Params
	endpoint string
	apiKey   string
	// WS Handling
	wsConn        *websocket.Conn
	options       *websocket.DialOptions
	socketID      string
	listen        atomic.Bool
	listenCtx     context.Context
	listenCtxStop context.CancelFunc
	handlers      map[string]*MessageHandler
	handlerMtx    sync.RWMutex
}

func GetKajiwotoWebSocketClient(endpoint, apiKey string) *KajiwotoWebSocketClient {
	// Init WebSocket Client
	return &KajiwotoWebSocketClient{
		endpoint: endpoint,
		apiKey:   apiKey,
		options:  &websocket.DialOptions{},
		handlers: make(map[string]*MessageHandler),
	}
}

func (c *KajiwotoWebSocketClient) Connect() error {
	if c.wsConn != nil {
		return errors.New("client is already connected")
	}
	// Dial Backend using client config
	conn, _, errClient := websocket.Dial(context.Background(), c.endpoint, c.options)
	if errClient != nil {
		return errClient
	}
	// Update conn reference
	c.wsConn = conn

	// Get Welcome message for initial handshake
	msgType, data, errWelcome := c.wsConn.Read(context.Background())
	if errWelcome != nil {
		return errWelcome
	}

	// Check if Server responded with welcome message
	if strconv.Itoa(int(msgType)) != DataFrameText {
		return fmt.Errorf("server did not respond with text frame. Message was: (%v)[%v]", msgType, string(data))
	}

	// Add Default Handlers
	c.AddDefaultHandlers()
	// Set WS Connection to listen for incoming messages
	c.StartListeningToMessages()

	// Add Auth Response Handler and wit for Auth to be confirmed
	authChannel := make(chan *KaiwotoWebSocketAuthResponse, 1)
	c.AddMessageHandler(NewKajiwotoWebSocketAuthResponseHandler(c, authChannel), true)

	// Send API Key to authenticate against the Kajiwoto websocket backend
	authMessage := &KajiwotoWebSocketMessage{
		MessageCode: SocketCodeMessageConnect,
		MessageContent: &KaiwotoWebSocketAuthRequest{
			ApiKey: c.apiKey,
		},
	}
	if errAuth := c.SendMessage(authMessage); errAuth != nil {
		return errAuth
	}

	// Wait for Auth channel to return socket ID as confirmation of successful login, or timeout is hit
	connectTimeout := time.NewTimer(time.Second * 5)
	for {
		select {
		case authResponse := <-authChannel:
			// Check if Socket ID was set
			if len(authResponse.Sid) > 0 {
				c.socketID = authResponse.Sid
				log.Debugf("Assigned Socket ID: %v", c.socketID)
				return nil
			}
			// In any other case, error
			c.StopListeningToMessages()
			c.RemoveAllMessageHandlers()
			return fmt.Errorf("server returned invalid auth message result: %+v", authResponse)
		case <-connectTimeout.C:
			c.StopListeningToMessages()
			c.RemoveAllMessageHandlers()
			return errors.New("connection timeout")
		}
	}
}

func (c *KajiwotoWebSocketClient) IsConnected() bool {
	return c.wsConn != nil && len(c.socketID) > 0
}

// AddDefaultHandlers
// ensures all basic handlers required to operate the WebSocket Client long term are set up and added to the client.
func (c *KajiwotoWebSocketClient) AddDefaultHandlers() {
	// Ping Handler
	c.AddMessageHandler(NewKajiwotoWebSocketPingHandler(c), false)
}

func (c *KajiwotoWebSocketClient) StartListeningToMessages() {
	// Start goroutine to handle incoming messages if it's not active
	if c.listen.CompareAndSwap(false, true) {
		c.listenCtx, c.listenCtxStop = context.WithCancel(context.Background())
		go func(c *KajiwotoWebSocketClient) {
			log.Debugf("Listening to incoming messages...")
			for c.listen.Load() {
				message, errRead := c.ReadMessage(c.listenCtx)
				if errRead != nil {
					log.Errorf("error reading websocket messages. Error: %v", errRead.Error())
					continue
				}

				// Pass message to all handlers
				c.handlerMtx.RLock()
				for _, handler := range c.handlers {
					go func(c *KajiwotoWebSocketClient, h *MessageHandler) {
						// Execute the handler, remove in case it's set up to remove itself
						if h.handleFunc(message) == nil && h.removeOnSuccess {
							c.RemoveMessageHandler(h.handlerKey)
							log.Debugf("Removed Message Handler '%v' after successful execution", h.handlerKey)
						}
					}(c, handler)
				}
				c.handlerMtx.RUnlock()
			}
			log.Debugf("Stopped listening to incoming messages.")
		}(c)
	}
}

func (c *KajiwotoWebSocketClient) StopListeningToMessages() {
	if c.listen.CompareAndSwap(true, false) {
		// Finish listen context and unset it
		c.listenCtxStop()
		c.listenCtx = nil
	}
}

func (c *KajiwotoWebSocketClient) AddMessageHandler(handleFunc MessageHandlerFunc, removeOnSuccess bool) (handlerKey string) {
	handlerKey = uuid.New().String()
	c.handlerMtx.Lock()
	c.handlers[handlerKey] = &MessageHandler{
		handlerKey:      handlerKey,
		handleFunc:      handleFunc,
		removeOnSuccess: removeOnSuccess,
	}
	c.handlerMtx.Unlock()
	log.Debugf("Added Message Handler '%v'. Autoremove: %v", handlerKey, removeOnSuccess)
	return handlerKey
}

func (c *KajiwotoWebSocketClient) RemoveMessageHandler(handlerKey string) {
	c.handlerMtx.Lock()
	delete(c.handlers, handlerKey)
	c.handlerMtx.Unlock()
	log.Debugf("Removed Message Handler '%v'", handlerKey)
}

func (c *KajiwotoWebSocketClient) RemoveAllMessageHandlers() {
	c.handlerMtx.Lock()
	c.handlers = make(map[string]*MessageHandler)
	c.handlerMtx.Unlock()
}

func (c *KajiwotoWebSocketClient) SendMessage(message *KajiwotoWebSocketMessage) error {
	bytes, errMessage := message.ToBytes()
	if errMessage != nil {
		return errMessage
	}
	// TODO: better context handling here
	log.Debugf("Sending message: %v", string(bytes))
	if errWrite := c.wsConn.Write(context.Background(), 1, bytes); errWrite != nil {
		return errWrite
	}
	return nil
}

func (c *KajiwotoWebSocketClient) ReadMessage(ctx context.Context) (*KajiwotoWebSocketMessage, error) {
	// Ensure this is only called once
	if c.listenCtx != nil && c.listenCtx != ctx {
		return nil, fmt.Errorf("client is already listening for new messages. Stop listening to manually handle reads")
	}

	msgType, data, errAPIResponse := c.wsConn.Read(ctx)
	if errAPIResponse != nil {
		return nil, errAPIResponse
	}
	// Check if Server responded with valid message
	if strconv.Itoa(int(msgType)) != DataFrameText {
		return nil, fmt.Errorf("server did not respond with text frame. Message was: (%v)[%v]", msgType, string(data))
	}
	log.Debugf("Received message: %v", string(data))
	message := &KajiwotoWebSocketMessage{}
	if errMessage := message.FromBytes(data); errMessage != nil {
		return nil, errMessage
	}
	return message, nil
}

// BuildLocalUserTime is sent whenever the backend needs to know the current time at the location of the user
func (c *KajiwotoWebSocketClient) BuildLocalUserTime() int {
	// Build Time
	hours, minutes, _ := time.Now().Clock()
	if minutes < 30 {
		minutes = 0
	} else {
		minutes = 30
	}
	localUserTime := hours*100 + minutes
	return localUserTime
}
