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
	client := GetKajiwotoWebSocketClient("wss://socket.chiefhappiness.co/socket.io/?EIO=4&transport=websocket", brokenKey)
	// Connect to Websocket Server
	errConnect := client.Connect()
	assert.Nil(s.T(), errConnect)
	// Check for Socket ID not assigned
	time.Sleep(time.Second * 2)
	assert.Empty(s.T(), client.socketID)
}

func (s *WebSocketClientTestSuite) TestWebSocketLoginCorrectAPIKey() {
	// Init client correct key
	client := GetKajiwotoWebSocketClient("wss://socket.chiefhappiness.co/socket.io/?EIO=4&transport=websocket", os.Getenv("WEBSOCKET_CLIENT_KEY"))
	// Connect to Websocket Server
	errConnect := client.Connect()
	assert.Nil(s.T(), errConnect)
	time.Sleep(time.Second)
	assert.NotEmpty(s.T(), client.socketID)
}

func (s *WebSocketClientTestSuite) TestWebSocketStopListening() {
	// Init client wrong key
	client := GetKajiwotoWebSocketClient("wss://socket.chiefhappiness.co/socket.io/?EIO=4&transport=websocket", os.Getenv("WEBSOCKET_CLIENT_KEY"))
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

func (s *WebSocketClientTestSuite) helperUserStatusChannelWithHandler() (chan *KajiwotoRPCUserStatusServerMessage, MessageHandlerFunc) {
	userStatusChannel := make(chan *KajiwotoRPCUserStatusServerMessage, 1)
	handleFunc := func(message *KajiwotoWebSocketMessage) error {
		if message.MessageCode == SocketCodeMessageEvent {
			// Deserialize Content
			rpcMessage := &KaiwotoRPCBaseMessage{}
			errDeserialize := rpcMessage.Deserialize(message.MessageContent)
			assert.Nil(s.T(), errDeserialize)

			// Handle Activity Message
			if rpcMessage.Action == RPCMessageUserStatus {
				statusMessage := &KajiwotoRPCUserStatusServerMessage{}
				assert.True(s.T(), statusMessage.FromRPCBaseMessage(rpcMessage))
				userStatusChannel <- statusMessage
			}

			return nil
		}
		return ErrUnableToHandleMessage
	}
	return userStatusChannel, handleFunc
}

func (s *WebSocketClientTestSuite) helperBuildDefaultRequestData(client *KajiwotoWebSocketClient) (userData KajiwotoRPCUserData) {
	// UserData
	photoUri := os.Getenv("WEBSOCKET_USER_PHOTO_URI")
	userData = KajiwotoRPCUserData{
		Guest:           false,
		UserID:          os.Getenv("WEBSOCKET_USER_ID"),
		DisplayName:     os.Getenv("WEBSOCKET_USER_DISPLAYNAME"),
		Username:        os.Getenv("WEBSOCKET_USER_USERNAME"),
		ProfilePhotoUri: &photoUri,
		Time:            client.BuildLocalUserTime(),
	}
	return userData
}

func (s *WebSocketClientTestSuite) TestWebSocketKajiRoomFlow() {
	// Init client correct key
	client := GetKajiwotoWebSocketClient("wss://socket.chiefhappiness.co/socket.io/?EIO=4&transport=websocket", os.Getenv("WEBSOCKET_CLIENT_KEY"))

	// Define channels used to wait for responses
	finishTestChannel := make(chan bool, 1)
	authChannel, authHandler := s.helperAuthChannelWithHandler()
	chatActivityChannel, chatActivityHandler := s.helperChatActivityChannelWithHandler()
	userStatusChannel, userStatusHandler := s.helperUserStatusChannelWithHandler()

	// Add Handlers
	_ = client.AddMessageHandler(authHandler, true)
	chatActivityHandlerKey := client.AddMessageHandler(chatActivityHandler, false)
	userStatusHandlerKey := client.AddMessageHandler(userStatusHandler, false)

	// Connect to Websocket Server
	errConnect := client.Connect()
	assert.Nil(s.T(), errConnect)

	// Build comman data required for handling
	userData := s.helperBuildDefaultRequestData(client)

	// Wait until timeout or all target messages received
	waitTimeout := time.NewTimer(time.Second * 5)
	done := false

	// Handle all events
	for !done {
		select {
		case socketId := <-authChannel:
			assert.NotEmpty(s.T(), socketId)

			// Create Login message & send it
			loginMessage := &KajiwotoRPCLoginMessage{
				UserData: userData,
				UserStatus: KajiwotoRPCUserStatus{
					Status: "ONLINE",
				},
				Secret: createMessageSecret(),
			}
			wsMessage := CreateKajiwotoWebSocketEventMessage(loginMessage)
			errSend := client.SendMessage(wsMessage)
			assert.Nil(s.T(), errSend)

		case userStatusUpdate := <-userStatusChannel:
			assert.NotNil(s.T(), userStatusUpdate)

			// Check if OK and then trigger subscribe Message
			equalUser := userStatusUpdate.StatusData.Data.UserID == userData.UserID
			userOnline := userStatusUpdate.StatusData.Data.Status == "ONLINE"
			assert.True(s.T(), equalUser)
			assert.True(s.T(), userOnline)

			if !equalUser || !userOnline {
				finishTestChannel <- true
			} else {
				// Send subscribe Message
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
			}

		case chatActivityUpdate := <-chatActivityChannel:
			assert.NotNil(s.T(), chatActivityUpdate)

			// Handle according to subtype
			switch chatActivityUpdate.ActivityData.Data.Action {
			case ChatActivityJoinRoom: // <- First Message to be received; after subscribe
				equalChatRoom := chatActivityUpdate.ActivityData.Data.ChatRoomId == os.Getenv("WEBSOCKET_CHATROOM_ID")
				assert.True(s.T(), equalChatRoom)
				if !equalChatRoom {
					finishTestChannel <- true
				} else {
					// Send enter Chat Message
					subscribeMessage := &KajiwotoRPCChatEnterMessage{
						UserData: userData,
						ChatroomData: KajiwotoRPCChatRoomData{
							ChatRoomId:    os.Getenv("WEBSOCKET_CHATROOM_ID"),
							IsPreviewRoom: false,
							LastMessages:  []KajiwotoRPCChatMessage{}, // TODO: Not sure if or how the AI is affected if these are omitted.
						},
						Secret: createMessageSecret(),
					}
					wsMessage := CreateKajiwotoWebSocketEventMessage(subscribeMessage)
					errSend := client.SendMessage(wsMessage)
					assert.Nil(s.T(), errSend)

					// This gives no further feeback from the backend, just finish test after sending it
					finishTestChannel <- true
				}
			//case ChatActivityPetMessage: // <- First Message to be received; after "chat enter"
			//	equalChatRoom := chatActivityUpdate.ActivityData.Data.ChatRoomId == os.Getenv("WEBSOCKET_CHATROOM_ID")
			//	assert.True(s.T(), equalChatRoom)
			//	if !equalChatRoom {
			//		finishTestChannel <- true
			//	} else {
			//		// Evaluate Pet Data
			//		petNotEmpty := chatActivityUpdate.ActivityData.Data.PetData != nil
			//		equalOwner := chatActivityUpdate.ActivityData.Data.PetData.OwnerId == userData.UserID
			//		assert.True(s.T(), petNotEmpty)
			//		assert.True(s.T(), equalOwner)
			//		if !petNotEmpty || !equalOwner {
			//			finishTestChannel <- true
			//		} else {
			//			// Send enter Chat Leave
			//			subscribeMessage := &KajiwotoRPCChatLeaveMessage{
			//				ChatRoom: KajiwotoRPCChatRoomId{
			//					ChatRoomId: os.Getenv("WEBSOCKET_CHATROOM_ID"),
			//				},
			//				Secret: createMessageSecret(),
			//			}
			//			wsMessage := CreateKajiwotoWebSocketEventMessage(subscribeMessage)
			//			errSend := client.SendMessage(wsMessage)
			//			assert.Nil(s.T(), errSend)
			//
			//			// This gives no further feeback from the backend, just finish test after sending it
			//			finishTestChannel <- true
			//		}
			//	}
			default:
				assert.True(s.T(), false, "unexpected Chat activity message received")
				finishTestChannel <- true
			}
		case <-waitTimeout.C:
			assert.True(s.T(), false, "timeout reached")
			// Shutdown
			finishTestChannel <- true
		case <-finishTestChannel:
			// Remove Handlers & Shutdown the client
			client.RemoveMessageHandler(chatActivityHandlerKey)
			client.RemoveMessageHandler(userStatusHandlerKey)
			client.StopListeningToMessages()
			done = true
			break
		}
	}

}
