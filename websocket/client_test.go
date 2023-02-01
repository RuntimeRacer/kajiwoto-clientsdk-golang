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
	log "github.com/sirupsen/logrus"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type WebSocketClientTestSuite struct {
	suite.Suite
}

func TestWebSocketClientTestSuite(t *testing.T) {
	suite.Run(t, new(WebSocketClientTestSuite))
}

func (s *WebSocketClientTestSuite) SetupTest() {
	// Set Log level for all tests
	log.SetLevel(log.DebugLevel)
}

func (s *WebSocketClientTestSuite) TestWebSocketLoginWrongAPIKey() {
	// Init client wrong key
	brokenKey := os.Getenv("WEBSOCKET_CLIENT_KEY")[1:4]
	client, errClient := GetKajiwotoClient("wss://socket.chiefhappiness.co/socket.io/?EIO=4&transport=websocket", brokenKey)
	assert.Nil(s.T(), errClient)
	// Connect to Websocket Server
	errConnect := client.Connect()
	assert.Nil(s.T(), errConnect)
	// Check for Socket ID not assigned
	time.Sleep(time.Second * 2)
	assert.Empty(s.T(), client.socketID)

}

func (s *WebSocketClientTestSuite) TestWebSocketLoginCorrectAPIKey() {
	// Init client wrong key
	client, errClient := GetKajiwotoClient("wss://socket.chiefhappiness.co/socket.io/?EIO=4&transport=websocket", os.Getenv("WEBSOCKET_CLIENT_KEY"))
	assert.Nil(s.T(), errClient)
	// Connect to Websocket Server
	errConnect := client.Connect()
	assert.Nil(s.T(), errConnect)
	time.Sleep(time.Second)
	assert.NotEmpty(s.T(), client.socketID)
}

func (s *WebSocketClientTestSuite) TestWebSocketStopListening() {
	// Init client wrong key
	client, errClient := GetKajiwotoClient("wss://socket.chiefhappiness.co/socket.io/?EIO=4&transport=websocket", os.Getenv("WEBSOCKET_CLIENT_KEY"))
	assert.Nil(s.T(), errClient)
	// Connect to Websocket Server
	errConnect := client.Connect()
	assert.Nil(s.T(), errConnect)
	// Stop listening before backend sends back auth message
	client.StopListeningToMessages()
	time.Sleep(time.Second)
	assert.Empty(s.T(), client.socketID)
}
