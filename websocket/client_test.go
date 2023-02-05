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
	"encoding/json"
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
	// Init client correct key
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

func (s *WebSocketClientTestSuite) helperAuthChannelWithHandler() (chan string, MessageHandlerFunc) {
	authChannel := make(chan string, 1)
	handleFunc := func(message *KajiwotoWebSocketMessage) error {
		if message.MessageCode == SocketCodeMessageConnect {
			// Try to umarshall into required response
			// If this won't work, message is not of expected type
			response := &KaiwotoWebSocketAuthResponse{}
			if errUnmarshall := json.Unmarshal(message.MessageContent.([]byte), response); errUnmarshall != nil {
				return errUnmarshall
			}
			authChannel <- response.Sid
			return nil
		}
		return ErrUnableToHandleMessage
	}
	return authChannel, handleFunc
}

func (s *WebSocketClientTestSuite) helperChatActivityChannelWithHandler() (chan *KajiwotoRPCChatActivityMessage, MessageHandlerFunc) {
	chatActivityChannel := make(chan *KajiwotoRPCChatActivityMessage, 1)
	handleFunc := func(message *KajiwotoWebSocketMessage) error {
		if message.MessageCode == SocketCodeMessageEvent {
			// Deserialize Content
			rpcMessage := &KaiwotoRPCBaseMessage{}
			errDeserialize := rpcMessage.Deserialize(message.MessageContent)
			assert.Nil(s.T(), errDeserialize)

			// Handle Activity Message
			if rpcMessage.Action == RPCMessageChatActivity {
				activityMessage := &KajiwotoRPCChatActivityMessage{}
				assert.True(s.T(), activityMessage.FromRPCBaseMessage(rpcMessage))
				chatActivityChannel <- activityMessage
			}

			return nil
		}
		return ErrUnableToHandleMessage
	}
	return chatActivityChannel, handleFunc
}

func (s *WebSocketClientTestSuite) TestWebSocketSubscribeToRoom() {
	// Init client correct key
	client, errClient := GetKajiwotoClient("wss://socket.chiefhappiness.co/socket.io/?EIO=4&transport=websocket", os.Getenv("WEBSOCKET_CLIENT_KEY"))
	assert.Nil(s.T(), errClient)

	// Define channels used to wait for responses
	closeChannel := make(chan bool, 1)
	authChannel, authHandler := s.helperAuthChannelWithHandler()
	chatActivityChannel, chatActivityHandler := s.helperChatActivityChannelWithHandler()

	// Add Handlers
	client.AddMessageHandler(authHandler, true)
	client.AddMessageHandler(chatActivityHandler, true)

	// Connect to Websocket Server
	errConnect := client.Connect()
	assert.Nil(s.T(), errConnect)

	// Handle all events
	// Wait until timeout or target message received
	waitTimeout := time.NewTimer(time.Second * 5)
	done := false
	for !done {
		select {
		case socketId := <-authChannel:
			assert.NotEmpty(s.T(), socketId)

			// Build User Data
			photoUri := os.Getenv("WEBSOCKET_USER_PHOTO_URI")
			userData := KajiwotoRPCUserData{
				Guest:           false,
				UserID:          os.Getenv("WEBSOCKET_USER_ID"),
				DisplayName:     os.Getenv("WEBSOCKET_USER_DISPLAYNAME"),
				Username:        os.Getenv("WEBSOCKET_USER_USERNAME"),
				ProfilePhotoUri: &photoUri,
				Time:            client.BuildLocalUserTime(),
			}

			// Create Login message & send it
			//loginMessage := &KajiwotoRPCLoginMessage{
			//	UserData: userData,
			//	UserStatus: KajiwotoRPCUserStatus{
			//		Status: "ONLINE",
			//	},
			//	Secret: client.createSecret(),
			//}
			//wsMessage := &KajiwotoWebSocketMessage{
			//	MessageCode:    SocketCodeMessageEvent,
			//	MessageContent: loginMessage.ToRPCBaseMessage().Serialize(),
			//}
			//errSend := client.SendMessage(wsMessage)
			//assert.Nil(s.T(), errSend)
			//time.Sleep(time.Millisecond * 500)

			//// Create User Status message & send it
			//statusMessage := &KajiwotoRPCUserStatusClientMessage{
			//	UserData: userData,
			//	UserStatus: KajiwotoRPCUserStatus{
			//		Status: "ONLINE",
			//	},
			//	Secret: client.createSecret(),
			//}
			//wsMessage = &KajiwotoWebSocketMessage{
			//	MessageCode:    SocketCodeMessageEvent,
			//	MessageContent: statusMessage.ToRPCBaseMessage().Serialize(),
			//}
			//errSend = client.SendMessage(wsMessage)
			//assert.Nil(s.T(), errSend)
			//time.Sleep(time.Millisecond * 100)

			// Create Live Sub message & send it
			//liveSubMessage := &KajiwotoRPCLiveSubMessage{
			//	Secret: client.createSecret(),
			//}
			//wsMessage := &KajiwotoWebSocketMessage{
			//	MessageCode:    SocketCodeMessageEvent,
			//	MessageContent: liveSubMessage.ToRPCBaseMessage().Serialize(),
			//}
			//errSend := client.SendMessage(wsMessage)
			//assert.Nil(s.T(), errSend)
			//time.Sleep(time.Millisecond * 100)

			// Create suscribe message & send it
			subscribeMessage := &KajiwotoRPCSubscribeMessage{
				UserData: userData,
				SubscribeArgs: KajiwotoRPCSubscribeArgs{
					ChatRoomIds: []string{os.Getenv("WEBSOCKET_CHATROOM_ID")},
				},
				Secret: createMessageSecret(),
			}
			wsMessage := CreateKajiwotoWebSocketEventMessage(subscribeMessage)
			errSend := client.SendMessage(wsMessage)
			assert.Nil(s.T(), errSend)
		case activityUpdate := <-chatActivityChannel:
			assert.NotNil(s.T(), activityUpdate)
			// TODO: Some more checks on the result
			// Shutdown
			closeChannel <- true
		case <-waitTimeout.C:
			assert.True(s.T(), false, "timeout reached")
			// Shutdown
			closeChannel <- true
		case <-closeChannel:
			// Create leave message & send it
			//subscribeMessage := &KajiwotoRPCChatLeaveMessage{
			//	ChatRoom: KajiwotoRPCChatRoomId{
			//		ChatRoomId: "vd8p",
			//	},
			//	Secret: client.createSecret(),
			//}
			//wsMessage := &KajiwotoWebSocketMessage{
			//	MessageCode:    SocketCodeMessageEvent,
			//	MessageContent: subscribeMessage.ToRPCBaseMessage().Serialize(),
			//}
			//errSend := client.SendMessage(wsMessage)
			//assert.Nil(s.T(), errSend)
			//time.Sleep(time.Millisecond * 500)

			client.StopListeningToMessages()
			done = true
			break
		}
	}

}
