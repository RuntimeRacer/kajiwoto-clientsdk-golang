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
)

// MessageHandlerFunc describes a function interface to be bound to the websocket client
type MessageHandlerFunc func(message *KajiwotoWebSocketMessage) error

// MessageHandler is used to bind handler functions to the websocket client
type MessageHandler struct {
	handlerKey      string // uuid to identify handler, internal
	handleFunc      MessageHandlerFunc
	removeOnSuccess bool
}

// KajiwotoClient is a custom websocket client for kajiwoto reqeusts using the websocket API
type KajiwotoClient struct {
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

func GetKajiwotoClient(endpoint, apiKey string) (*KajiwotoClient, error) {
	// Init Websocket Client Config
	//origin := "http://localhost/"
	options := &websocket.DialOptions{}

	// Set Handshake Headers

	//if errMarshal != nil {
	//	return nil, errMarshal
	//}
	//config.Header.Add("auth", string(authHeader))

	return &KajiwotoClient{
		endpoint: endpoint,
		apiKey:   apiKey,
		options:  options,
		handlers: make(map[string]*MessageHandler),
	}, nil
}

func (c *KajiwotoClient) Connect() error {
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

	// Get Welcome message - TODO: Replace with Handler
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
	return nil
}

// AddDefaultHandlers
// ensures all basic handlers required to operate the WebSocket Client long term are set up and added to the client.
func (c *KajiwotoClient) AddDefaultHandlers() {
	// Ping Handler
	c.AddMessageHandler(NewKajiwotoWebSocketPingHandler(c), false)
	// Auth Response handler - Updates Socket ID
	c.AddMessageHandler(NewKajiwotoWebSocketAuthResponseHandler(c), false)
}

func (c *KajiwotoClient) StartListeningToMessages() {
	// Start goroutine to handle incoming messages if it's not active
	if c.listen.CompareAndSwap(false, true) {
		c.listenCtx, c.listenCtxStop = context.WithCancel(context.Background())
		go func(c *KajiwotoClient) {
			for c.listen.Load() {
				message, errRead := c.ReadMessage(c.listenCtx)
				if errRead != nil {
					log.Errorf("error reading websocket messages. Error: %v", errRead.Error())
					continue
				}

				// Pass message to all handlers
				c.handlerMtx.RLock()
				for _, handler := range c.handlers {
					go func(c *KajiwotoClient, h *MessageHandler) {
						// Execute the handler, remove in case it's set up to remove itself
						if h.handleFunc(message) == nil && h.removeOnSuccess {
							c.RemoveMessageHandler(h.handlerKey)
							log.Debugf("Removed Message Handler '%v' after successful execution", h.handlerKey)
						}
					}(c, handler)
				}
				c.handlerMtx.RUnlock()
			}
		}(c)
	}
}

func (c *KajiwotoClient) StopListeningToMessages() {
	if c.listen.CompareAndSwap(true, false) {
		// Finish listen context and unset it
		c.listenCtxStop()
		c.listenCtx = nil
	}
}

func (c *KajiwotoClient) AddMessageHandler(handleFunc MessageHandlerFunc, removeOnSuccess bool) {
	c.handlerMtx.Lock()
	handlerKey := uuid.New()
	c.handlers[handlerKey.String()] = &MessageHandler{
		handlerKey:      handlerKey.String(),
		handleFunc:      handleFunc,
		removeOnSuccess: removeOnSuccess,
	}
	c.handlerMtx.Unlock()
	log.Debugf("Added Message Handler '%v'. Autoremove: %v", handlerKey, removeOnSuccess)
}

func (c *KajiwotoClient) RemoveMessageHandler(handlerKey string) {
	c.handlerMtx.Lock()
	delete(c.handlers, handlerKey)
	c.handlerMtx.Unlock()
	log.Debugf("Removed Message Handler '%v'", handlerKey)
}

func (c *KajiwotoClient) SendMessage(message *KajiwotoWebSocketMessage) error {
	bytes, errMessage := message.ToBytes()
	if errMessage != nil {
		return errMessage
	}
	// TODO: better context handling here
	if errWrite := c.wsConn.Write(context.Background(), 1, bytes); errWrite != nil {
		return errWrite
	}
	return nil
}

func (c *KajiwotoClient) ReadMessage(ctx context.Context) (*KajiwotoWebSocketMessage, error) {
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
	message := &KajiwotoWebSocketMessage{}
	if errMessage := message.FromBytes(data); errMessage != nil {
		return nil, errMessage
	}
	return message, nil
}
