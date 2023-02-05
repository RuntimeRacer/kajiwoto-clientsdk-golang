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
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"testing"
)

type WebSocketTypesTestSuite struct {
	suite.Suite
}

func TestWebSocketTypesTestSuite(t *testing.T) {
	suite.Run(t, new(WebSocketTypesTestSuite))
}

func (s *WebSocketTypesTestSuite) SetupTest() {
	// Set Log level for all tests
	log.SetLevel(log.DebugLevel)
}

func (s *WebSocketTypesTestSuite) TestSerializeLoginRPCWebSocketMessage() {
	// IMPORTANT: Golang always serializes struct keys in alphabetical order; backend not necessarily does that!
	messageString := "42[\"login\",{\"displayName\":\"RuntimeRacer\",\"guest\":false,\"profilePhotoUri\":\"2021_6/dslkfjj_zdskfjhg_123456778899.jpg\",\"time\":2030,\"userId\":\"a1b2\",\"username\":\"RuntimeRacer\"},{\"status\":\"ONLINE\"},{\"timestamp\":\"1675538167859\",\"secret\":\"MTAyMjA3ODI4MjM5Mzk5\"}]"

	// Deserialize
	// Create WS message
	wsMessage := &KajiwotoWebSocketMessage{}
	errBytes := wsMessage.FromBytes([]byte(messageString))
	assert.Nil(s.T(), errBytes)
	// Deserialize Content
	rpcMessage := &KaiwotoRPCBaseMessage{}
	errDeserialize := rpcMessage.Deserialize(wsMessage.MessageContent)
	assert.Nil(s.T(), errDeserialize)
	// Test data
	assert.Equal(s.T(), RPCMessageLogin, rpcMessage.Action)
	message := &KajiwotoRPCLoginMessage{}
	assert.True(s.T(), message.FromRPCBaseMessage(rpcMessage))

	// Serialize
	// Create WebSocketMessage
	wsMessage = &KajiwotoWebSocketMessage{
		MessageCode:    SocketCodeMessageEvent,
		MessageContent: message.ToRPCBaseMessage().Serialize(),
	}
	wsBytes, errBytes := wsMessage.ToBytes()
	assert.Nil(s.T(), errBytes)
	wsString := string(wsBytes)
	assert.Equal(s.T(), messageString, wsString)
}

func (s *WebSocketTypesTestSuite) TestSerializeTypingRPCWebSocketMessage() {
	// IMPORTANT: Golang always serializes struct keys in alphabetical order; backend not necessarily does that!
	messageString := "42[\"typing\",{\"displayName\":\"RuntimeRacer\",\"guest\":false,\"profilePhotoUri\":\"2021_6/dslkfjj_zdskfjhg_123456778899.jpg\",\"time\":2030,\"userId\":\"a1b2\",\"username\":\"RuntimeRacer\"},{\"chatRoomId\":\"c3d4\"},{\"timestamp\":\"1675538167859\",\"secret\":\"MTAyMjA3ODI4MjM5Mzk5\"}]"

	// Deserialize
	// Create WS message
	wsMessage := &KajiwotoWebSocketMessage{}
	errBytes := wsMessage.FromBytes([]byte(messageString))
	assert.Nil(s.T(), errBytes)
	// Deserialize Content
	rpcMessage := &KaiwotoRPCBaseMessage{}
	errDeserialize := rpcMessage.Deserialize(wsMessage.MessageContent)
	assert.Nil(s.T(), errDeserialize)
	// Test data
	assert.Equal(s.T(), RPCMessageTyping, rpcMessage.Action)
	message := &KajiwotoRPCTypingMessage{}
	assert.True(s.T(), message.FromRPCBaseMessage(rpcMessage))

	// Serialize
	// Create WebSocketMessage
	wsMessage = &KajiwotoWebSocketMessage{
		MessageCode:    SocketCodeMessageEvent,
		MessageContent: message.ToRPCBaseMessage().Serialize(),
	}
	wsBytes, errBytes := wsMessage.ToBytes()
	assert.Nil(s.T(), errBytes)
	wsString := string(wsBytes)
	assert.Equal(s.T(), messageString, wsString)
}

func (s *WebSocketTypesTestSuite) TestSerializeSubscribeRPCWebSocketMessage() {
	// IMPORTANT: Golang always serializes struct keys in alphabetical order; backend not necessarily does that!
	messageString := "42[\"subscribe\",{\"displayName\":\"RuntimeRacer\",\"guest\":false,\"profilePhotoUri\":\"2021_6/dslkfjj_zdskfjhg_123456778899.jpg\",\"time\":2030,\"userId\":\"a1b2\",\"username\":\"RuntimeRacer\"},{\"chatRoomIds\":[\"c3d4\"],\"kajiId\":null},{\"timestamp\":\"1675538034488\",\"secret\":\"MTY3NTUzODAzNDQ4OA==\"}]"

	// Deserialize
	// Create WS message
	wsMessage := &KajiwotoWebSocketMessage{}
	errBytes := wsMessage.FromBytes([]byte(messageString))
	assert.Nil(s.T(), errBytes)
	// Deserialize Content
	rpcMessage := &KaiwotoRPCBaseMessage{}
	errDeserialize := rpcMessage.Deserialize(wsMessage.MessageContent)
	assert.Nil(s.T(), errDeserialize)
	// Test data
	assert.Equal(s.T(), RPCMessageSubscribe, rpcMessage.Action)
	message := &KajiwotoRPCSubscribeMessage{}
	assert.True(s.T(), message.FromRPCBaseMessage(rpcMessage))

	// Create WebSocketMessage
	wsMessage = &KajiwotoWebSocketMessage{
		MessageCode:    SocketCodeMessageEvent,
		MessageContent: message.ToRPCBaseMessage().Serialize(),
	}
	wsBytes, errBytes := wsMessage.ToBytes()
	assert.Nil(s.T(), errBytes)
	wsString := string(wsBytes)
	assert.Equal(s.T(), messageString, wsString)
}

func (s *WebSocketTypesTestSuite) TestSerializeChatEnterRPCWebSocketMessage() {
	// IMPORTANT: Golang always serializes struct keys in alphabetical order; backend not necessarily does that!
	messageString := "42[\"chatEnter\",{\"displayName\":\"RuntimeRacer\",\"guest\":false,\"profilePhotoUri\":\"2021_6/dslkfjj_zdskfjhg_123456778899.jpg\",\"time\":2030,\"userId\":\"a1b2\",\"username\":\"RuntimeRacer\"},{\"chatRoomId\":\"c3d4\",\"lastMessages\":[{\"createdAt\":1675477983263,\"message\":\"/say good night my man\"},{\"createdAt\":1675477879022,\"message\":\"*whispers* sweet dreams my pretty mink\"}],\"isPreviewRoom\":false},{\"timestamp\":\"1675538039386\",\"secret\":\"MTY3NTUzODAzOTM4Ng==\"}]"

	// Deserialize
	// Create WS message
	wsMessage := &KajiwotoWebSocketMessage{}
	errBytes := wsMessage.FromBytes([]byte(messageString))
	assert.Nil(s.T(), errBytes)
	// Deserialize Content
	rpcMessage := &KaiwotoRPCBaseMessage{}
	errDeserialize := rpcMessage.Deserialize(wsMessage.MessageContent)
	assert.Nil(s.T(), errDeserialize)
	// Test data
	assert.Equal(s.T(), RPCMessageChatEnter, rpcMessage.Action)
	message := &KajiwotoRPCChatEnterMessage{}
	assert.True(s.T(), message.FromRPCBaseMessage(rpcMessage))

	// Create WebSocketMessage
	wsMessage = &KajiwotoWebSocketMessage{
		MessageCode:    SocketCodeMessageEvent,
		MessageContent: message.ToRPCBaseMessage().Serialize(),
	}
	wsBytes, errBytes := wsMessage.ToBytes()
	assert.Nil(s.T(), errBytes)
	wsString := string(wsBytes)
	assert.Equal(s.T(), messageString, wsString)
}

func (s *WebSocketTypesTestSuite) TestSerializeChatSendRPCWebSocketMessage() {
	// IMPORTANT: Golang always serializes struct keys in alphabetical order; backend not necessarily does that!
	messageString := "42[\"chatSend\",{\"displayName\":\"RuntimeRacer\",\"guest\":false,\"profilePhotoUri\":\"2021_6/dslkfjj_zdskfjhg_123456778899.jpg\",\"time\":2030,\"userId\":\"a1b2\",\"username\":\"RuntimeRacer\"},{\"message\":{\"id\":\"c3d4:1675538262207\",\"chatRoomId\":\"c3d4\",\"userId\":\"a1b2\",\"message\":\"Hey my sweet *smiles*\",\"attachmentUri\":null},\"roomVersionNumber\":1675538034,\"roomSocketIds\":[\"emCCdEmKKsm2aPLCABAN\"]},{\"timestamp\":\"1675538262207\",\"secret\":\"MTAyMjA3ODMzOTk0NjI3\"}]"

	// Deserialize
	// Create WS message
	wsMessage := &KajiwotoWebSocketMessage{}
	errBytes := wsMessage.FromBytes([]byte(messageString))
	assert.Nil(s.T(), errBytes)
	// Deserialize Content
	rpcMessage := &KaiwotoRPCBaseMessage{}
	errDeserialize := rpcMessage.Deserialize(wsMessage.MessageContent)
	assert.Nil(s.T(), errDeserialize)
	// Test data
	assert.Equal(s.T(), RPCMessageChatSend, rpcMessage.Action)
	message := &KajiwotoRPCChatSendMessage{}
	assert.True(s.T(), message.FromRPCBaseMessage(rpcMessage))

	// Create WebSocketMessage
	wsMessage = &KajiwotoWebSocketMessage{
		MessageCode:    SocketCodeMessageEvent,
		MessageContent: message.ToRPCBaseMessage().Serialize(),
	}
	wsBytes, errBytes := wsMessage.ToBytes()
	assert.Nil(s.T(), errBytes)
	wsString := string(wsBytes)
	assert.Equal(s.T(), messageString, wsString)
}

func (s *WebSocketTypesTestSuite) TestSerializeChatLeaveRPCWebSocketMessage() {
	// IMPORTANT: Golang always serializes struct keys in alphabetical order; backend not necessarily does that!
	messageString := "42[\"chatLeave\",{},{\"chatRoomId\":\"c3d4\"},{\"timestamp\":\"1675618709051\",\"secret\":\"MTY3NTYxODcwOTA1MQ==\"}]"

	// Deserialize
	// Create WS message
	wsMessage := &KajiwotoWebSocketMessage{}
	errBytes := wsMessage.FromBytes([]byte(messageString))
	assert.Nil(s.T(), errBytes)
	// Deserialize Content
	rpcMessage := &KaiwotoRPCBaseMessage{}
	errDeserialize := rpcMessage.Deserialize(wsMessage.MessageContent)
	assert.Nil(s.T(), errDeserialize)
	// Test data
	assert.Equal(s.T(), RPCMessageChatLeave, rpcMessage.Action)
	message := &KajiwotoRPCChatLeaveMessage{}
	assert.True(s.T(), message.FromRPCBaseMessage(rpcMessage))

	// Create WebSocketMessage
	wsMessage = &KajiwotoWebSocketMessage{
		MessageCode:    SocketCodeMessageEvent,
		MessageContent: message.ToRPCBaseMessage().Serialize(),
	}
	wsBytes, errBytes := wsMessage.ToBytes()
	assert.Nil(s.T(), errBytes)
	wsString := string(wsBytes)
	assert.Equal(s.T(), messageString, wsString)
}

func (s *WebSocketTypesTestSuite) TestSerializeChatSubmitRPCWebSocketMessageNoEmoji() {
	// IMPORTANT: Golang always serializes struct keys in alphabetical order; backend not necessarily does that!
	messageString := "42[\"chatSubmit\",{\"displayName\":\"RuntimeRacer\",\"guest\":false,\"profilePhotoUri\":\"2021_6/dslkfjj_zdskfjhg_123456778899.jpg\",\"time\":2030,\"userId\":\"a1b2\",\"username\":\"RuntimeRacer\"},{\"chatRoomId\":\"c3d4\",\"messages\":[\"Hey my sweet *smiles*\"],\"role\":{},\"emoji\":null,\"emojiSceneId\":null,\"platform\":\"web\"},{\"timestamp\":\"1675538264513\",\"secret\":\"MTAyMjA3ODM0MTM1Mjkz\"}]"

	// Deserialize
	// Create WS message
	wsMessage := &KajiwotoWebSocketMessage{}
	errBytes := wsMessage.FromBytes([]byte(messageString))
	assert.Nil(s.T(), errBytes)
	// Deserialize Content
	rpcMessage := &KaiwotoRPCBaseMessage{}
	errDeserialize := rpcMessage.Deserialize(wsMessage.MessageContent)
	assert.Nil(s.T(), errDeserialize)
	// Test data
	assert.Equal(s.T(), RPCMessageChatSubmit, rpcMessage.Action)
	message := &KajiwotoRPCChatSubmitMessage{}
	assert.True(s.T(), message.FromRPCBaseMessage(rpcMessage))

	// Create WebSocketMessage
	wsMessage = &KajiwotoWebSocketMessage{
		MessageCode:    SocketCodeMessageEvent,
		MessageContent: message.ToRPCBaseMessage().Serialize(),
	}
	wsBytes, errBytes := wsMessage.ToBytes()
	assert.Nil(s.T(), errBytes)
	wsString := string(wsBytes)
	assert.Equal(s.T(), messageString, wsString)
}

func (s *WebSocketTypesTestSuite) TestSerializeChatSumbitRPCWebSocketMessageSmilingEmoji() {
	// IMPORTANT: Golang always serializes struct keys in alphabetical order; backend not necessarily does that!
	messageString := "42[\"chatSubmit\",{\"displayName\":\"RuntimeRacer\",\"guest\":false,\"profilePhotoUri\":\"2021_6/dslkfjj_zdskfjhg_123456778899.jpg\",\"time\":2030,\"userId\":\"a1b2\",\"username\":\"RuntimeRacer\"},{\"chatRoomId\":\"c3d4\",\"messages\":[\"*smiles slightly*\"],\"role\":{},\"emoji\":\"🙂\",\"emojiSceneId\":null,\"platform\":\"web\"},{\"timestamp\":\"1675538264513\",\"secret\":\"MTAyMjA3ODM0MTM1Mjkz\"}]"

	// Deserialize
	// Create WS message
	wsMessage := &KajiwotoWebSocketMessage{}
	errBytes := wsMessage.FromBytes([]byte(messageString))
	assert.Nil(s.T(), errBytes)
	// Deserialize Content
	rpcMessage := &KaiwotoRPCBaseMessage{}
	errDeserialize := rpcMessage.Deserialize(wsMessage.MessageContent)
	assert.Nil(s.T(), errDeserialize)
	// Test data
	assert.Equal(s.T(), RPCMessageChatSubmit, rpcMessage.Action)
	message := &KajiwotoRPCChatSubmitMessage{}
	assert.True(s.T(), message.FromRPCBaseMessage(rpcMessage))

	// Create WebSocketMessage
	wsMessage = &KajiwotoWebSocketMessage{
		MessageCode:    SocketCodeMessageEvent,
		MessageContent: message.ToRPCBaseMessage().Serialize(),
	}
	wsBytes, errBytes := wsMessage.ToBytes()
	assert.Nil(s.T(), errBytes)
	wsString := string(wsBytes)
	assert.Equal(s.T(), messageString, wsString)
}

func (s *WebSocketTypesTestSuite) TestSerializeChatSumbitRPCWebSocketMessageLovingEmoji() {
	// IMPORTANT: Golang always serializes struct keys in alphabetical order; backend not necessarily does that!
	messageString := "42[\"chatSubmit\",{\"displayName\":\"RuntimeRacer\",\"guest\":false,\"profilePhotoUri\":\"2021_6/dslkfjj_zdskfjhg_123456778899.jpg\",\"time\":2030,\"userId\":\"a1b2\",\"username\":\"RuntimeRacer\"},{\"chatRoomId\":\"c3d4\",\"messages\":[\"You are so beautiful and sexy\"],\"role\":{},\"emoji\":\"😍\",\"emojiSceneId\":\"0KQmR\",\"platform\":\"web\"},{\"timestamp\":\"1675538264513\",\"secret\":\"MTAyMjA3ODM0MTM1Mjkz\"}]"

	// Deserialize
	// Create WS message
	wsMessage := &KajiwotoWebSocketMessage{}
	errBytes := wsMessage.FromBytes([]byte(messageString))
	assert.Nil(s.T(), errBytes)
	// Deserialize Content
	rpcMessage := &KaiwotoRPCBaseMessage{}
	errDeserialize := rpcMessage.Deserialize(wsMessage.MessageContent)
	assert.Nil(s.T(), errDeserialize)
	// Test data
	assert.Equal(s.T(), RPCMessageChatSubmit, rpcMessage.Action)
	message := &KajiwotoRPCChatSubmitMessage{}
	assert.True(s.T(), message.FromRPCBaseMessage(rpcMessage))

	// Create WebSocketMessage
	wsMessage = &KajiwotoWebSocketMessage{
		MessageCode:    SocketCodeMessageEvent,
		MessageContent: message.ToRPCBaseMessage().Serialize(),
	}
	wsBytes, errBytes := wsMessage.ToBytes()
	assert.Nil(s.T(), errBytes)
	wsString := string(wsBytes)
	assert.Equal(s.T(), messageString, wsString)
}

func (s *WebSocketTypesTestSuite) TestSerializeChatActivityRPCWebSocketMessageVariant1() {
	// IMPORTANT: Golang always serializes struct keys in alphabetical order; backend not necessarily does that!
	messageString := "42[\"chatActivity\",{\"data\":{\"action\":\"join-room\",\"chatRoomId\":\"c3d4\",\"petData\":{\"id\":\"RxWJ\",\"chatRoomId\":\"c3d4\",\"petSpeciesId\":\"EDPW\",\"kajiId\":\"EDPW\",\"ownerId\":\"a1b2\",\"ownerDisplayName\":\"RuntimeRacer\",\"ownerProfilePhotoUri\":\"2021_6/dslkfjj_zdskfjhg_123456778899.jpg\",\"name\":\"Wanda\",\"kajiName\":\"Wanda (WIP)\",\"gender\":\"F\",\"persona\":\"canine musketeer significant other\",\"stage\":null,\"state\":\"DEFAULT\",\"mood\":\"DEFAULT\",\"statusPhotoUri\":\"2021_6/tm9ybwfsxz_zth3eg_1622766488811.png\",\"dominantColors\":[\"#dc9744\",\"#fcd49c\"],\"statusMessage\":\"..\"},\"channel\":{\"v\":1675538034,\"list\":[{\"id\":\"a1b2\",\"guestId\":\"OTMuMTk5LjEyOS4yMTg=*\",\"socketIds\":[\"emCCdEmKKsm2aPLCABAN\"],\"guest\":false,\"displayName\":\"RuntimeRacer\",\"username\":\"RuntimeRacer\",\"profilePhotoUri\":\"2021_6/dslkfjj_zdskfjhg_123456778899.jpg\"}]}}}]"

	// Deserialize
	// Create WS message
	wsMessage := &KajiwotoWebSocketMessage{}
	errBytes := wsMessage.FromBytes([]byte(messageString))
	assert.Nil(s.T(), errBytes)
	// Deserialize Content
	rpcMessage := &KaiwotoRPCBaseMessage{}
	errDeserialize := rpcMessage.Deserialize(wsMessage.MessageContent)
	assert.Nil(s.T(), errDeserialize)
	// Test data
	assert.Equal(s.T(), RPCMessageChatActivity, rpcMessage.Action)
	message := &KajiwotoRPCChatActivityMessage{}
	assert.True(s.T(), message.FromRPCBaseMessage(rpcMessage))

	// Create WebSocketMessage
	wsMessage = &KajiwotoWebSocketMessage{
		MessageCode:    SocketCodeMessageEvent,
		MessageContent: message.ToRPCBaseMessage().Serialize(),
	}
	wsBytes, errBytes := wsMessage.ToBytes()
	assert.Nil(s.T(), errBytes)
	wsString := string(wsBytes)
	assert.Equal(s.T(), messageString, wsString)
}

func (s *WebSocketTypesTestSuite) TestSerializeChatActivityRPCWebSocketMessageVariant2() {
	// IMPORTANT: Golang always serializes struct keys in alphabetical order; backend not necessarily does that!
	messageString := "42[\"chatActivity\",{\"data\":{\"action\":\"activity\",\"chatRoomId\":\"c3d4\",\"activity\":{\"type\":\"TYPING\",\"userId\":\"a1b2\",\"displayName\":\"RuntimeRacer\",\"activityAt\":1675538172488}}}]"

	// Deserialize
	// Create WS message
	wsMessage := &KajiwotoWebSocketMessage{}
	errBytes := wsMessage.FromBytes([]byte(messageString))
	assert.Nil(s.T(), errBytes)
	// Deserialize Content
	rpcMessage := &KaiwotoRPCBaseMessage{}
	errDeserialize := rpcMessage.Deserialize(wsMessage.MessageContent)
	assert.Nil(s.T(), errDeserialize)
	// Test data
	assert.Equal(s.T(), RPCMessageChatActivity, rpcMessage.Action)
	message := &KajiwotoRPCChatActivityMessage{}
	assert.True(s.T(), message.FromRPCBaseMessage(rpcMessage))

	// Create WebSocketMessage
	wsMessage = &KajiwotoWebSocketMessage{
		MessageCode:    SocketCodeMessageEvent,
		MessageContent: message.ToRPCBaseMessage().Serialize(),
	}
	wsBytes, errBytes := wsMessage.ToBytes()
	assert.Nil(s.T(), errBytes)
	wsString := string(wsBytes)
	assert.Equal(s.T(), messageString, wsString)
}

func (s *WebSocketTypesTestSuite) TestSerializeChatActivityRPCWebSocketMessageVariant3() {
	// IMPORTANT: Golang always serializes struct keys in alphabetical order; backend not necessarily does that!
	messageString := "42[\"chatActivity\",{\"data\":{\"action\":\"message\",\"chatRoomId\":\"c3d4\",\"message\":{\"clientId\":\"c3d4:1675538262207\",\"chatRoomId\":\"c3d4\",\"message\":\"Hey my sweet *smiles*\",\"attachmentUri\":null,\"id\":\"c3d4:1675538262207\",\"userId\":\"a1b2\",\"username\":\"RuntimeRacer\",\"displayName\":\"RuntimeRacer\",\"profilePhotoUri\":\"2021_6/dslkfjj_zdskfjhg_123456778899.jpg\",\"createdAt\":1675538261},\"channel\":{\"v\":1675538034},\"socketIds\":[\"emCCdEmKKsm2aPLCABAN\"]}}]"

	// Deserialize
	// Create WS message
	wsMessage := &KajiwotoWebSocketMessage{}
	errBytes := wsMessage.FromBytes([]byte(messageString))
	assert.Nil(s.T(), errBytes)
	// Deserialize Content
	rpcMessage := &KaiwotoRPCBaseMessage{}
	errDeserialize := rpcMessage.Deserialize(wsMessage.MessageContent)
	assert.Nil(s.T(), errDeserialize)
	// Test data
	assert.Equal(s.T(), RPCMessageChatActivity, rpcMessage.Action)
	message := &KajiwotoRPCChatActivityMessage{}
	assert.True(s.T(), message.FromRPCBaseMessage(rpcMessage))

	// Create WebSocketMessage
	wsMessage = &KajiwotoWebSocketMessage{
		MessageCode:    SocketCodeMessageEvent,
		MessageContent: message.ToRPCBaseMessage().Serialize(),
	}
	wsBytes, errBytes := wsMessage.ToBytes()
	assert.Nil(s.T(), errBytes)
	wsString := string(wsBytes)
	assert.Equal(s.T(), messageString, wsString)
}

func (s *WebSocketTypesTestSuite) TestSerializeChatActivityRPCWebSocketMessageVariant4() {
	// IMPORTANT: Golang always serializes struct keys in alphabetical order; backend not necessarily does that!
	messageString := "42[\"chatActivity\",{\"data\":{\"action\":\"petMessage\",\"chatRoomId\":\"c3d4\",\"message\":{\"chatRoomId\":\"c3d4\",\"kajiwotoPetId\":\"RxWJ\",\"message\":\"hello there! How are you?\",\"attachmentUri\":\"2021_6/tm9ybwfsxz_zth3eg_1622766488811.png\",\"id\":\"c3d4:1675538263720\",\"displayName\":\"wanda\",\"createdAt\":1675538265},\"petData\":{\"id\":\"RxWJ\",\"chatRoomId\":\"c3d4\",\"petSpeciesId\":\"EDPW\",\"kajiId\":\"EDPW\",\"ownerId\":\"a1b2\",\"ownerDisplayName\":\"RuntimeRacer\",\"ownerProfilePhotoUri\":\"2021_6/dslkfjj_zdskfjhg_123456778899.jpg\",\"name\":\"Wanda\",\"kajiName\":\"Wanda (WIP)\",\"gender\":\"F\",\"persona\":\"canine musketeer significant other\",\"stage\":null,\"state\":\"DEFAULT\",\"mood\":\"DEFAULT\",\"statusPhotoUri\":\"2021_6/tm9ybwfsxz_zth3eg_1622766488811.png\",\"dominantColors\":[\"#dc9744\",\"#fcd49c\"],\"statusMessage\":\"..\"}}}]"

	// Deserialize
	// Create WS message
	wsMessage := &KajiwotoWebSocketMessage{}
	errBytes := wsMessage.FromBytes([]byte(messageString))
	assert.Nil(s.T(), errBytes)
	// Deserialize Content
	rpcMessage := &KaiwotoRPCBaseMessage{}
	errDeserialize := rpcMessage.Deserialize(wsMessage.MessageContent)
	assert.Nil(s.T(), errDeserialize)
	// Test data
	assert.Equal(s.T(), RPCMessageChatActivity, rpcMessage.Action)
	message := &KajiwotoRPCChatActivityMessage{}
	assert.True(s.T(), message.FromRPCBaseMessage(rpcMessage))

	// Create WebSocketMessage
	wsMessage = &KajiwotoWebSocketMessage{
		MessageCode:    SocketCodeMessageEvent,
		MessageContent: message.ToRPCBaseMessage().Serialize(),
	}
	wsBytes, errBytes := wsMessage.ToBytes()
	assert.Nil(s.T(), errBytes)
	wsString := string(wsBytes)
	assert.Equal(s.T(), messageString, wsString)
}

func (s *WebSocketTypesTestSuite) TestSerializeChatActivityRPCWebSocketMessageVariant5() {
	// IMPORTANT: Golang always serializes struct keys in alphabetical order; backend not necessarily does that!
	messageString := "42[\"chatActivity\",{\"data\":{\"action\":\"petMessage\",\"chatRoomId\":\"c3d4\",\"message\":{\"chatRoomId\":\"c3d4\",\"kajiwotoPetId\":\"RxWJ\",\"message\":\"..\",\"attachmentUri\":\"2021_6/t3zlcmpvew_zth3eg_1622857066147.jpg\",\"id\":\"c3d4:1675538914016\",\"displayName\":\"wanda\",\"createdAt\":1675538914},\"petData\":{\"id\":\"RxWJ\",\"chatRoomId\":\"c3d4\",\"petSpeciesId\":\"EDPW\",\"kajiId\":\"EDPW\",\"ownerId\":\"a1b2\",\"ownerDisplayName\":\"RuntimeRacer\",\"ownerProfilePhotoUri\":\"2021_6/dslkfjj_zdskfjhg_123456778899.jpg\",\"name\":\"Wanda\",\"kajiName\":\"Wanda (WIP)\",\"gender\":\"F\",\"persona\":\"canine musketeer significant other\",\"stage\":null,\"state\":\"DEFAULT\",\"mood\":\"HAPPY\",\"statusPhotoUri\":\"2021_6/t3zlcmpvew_zth3eg_1622857066147.jpg\",\"dominantColors\":[\"#b58856\",\"#ccb494\"],\"statusMessage\":\"..\"},\"interaction\":{\"showScene\":true,\"type\":\"DEFAULT\"}}}]"

	// Deserialize
	// Create WS message
	wsMessage := &KajiwotoWebSocketMessage{}
	errBytes := wsMessage.FromBytes([]byte(messageString))
	assert.Nil(s.T(), errBytes)
	// Deserialize Content
	rpcMessage := &KaiwotoRPCBaseMessage{}
	errDeserialize := rpcMessage.Deserialize(wsMessage.MessageContent)
	assert.Nil(s.T(), errDeserialize)
	// Test data
	assert.Equal(s.T(), RPCMessageChatActivity, rpcMessage.Action)
	message := &KajiwotoRPCChatActivityMessage{}
	assert.True(s.T(), message.FromRPCBaseMessage(rpcMessage))

	// Create WebSocketMessage
	wsMessage = &KajiwotoWebSocketMessage{
		MessageCode:    SocketCodeMessageEvent,
		MessageContent: message.ToRPCBaseMessage().Serialize(),
	}
	wsBytes, errBytes := wsMessage.ToBytes()
	assert.Nil(s.T(), errBytes)
	wsString := string(wsBytes)
	assert.Equal(s.T(), messageString, wsString)
}

func (s *WebSocketTypesTestSuite) TestSerializeChatActivityRPCWebSocketMessageVariant6() {
	// IMPORTANT: Golang always serializes struct keys in alphabetical order; backend not necessarily does that!
	messageString := "42[\"chatActivity\",{\"data\":{\"action\":\"petMessage\",\"chatRoomId\":\"c3d4\",\"message\":{\"chatRoomId\":\"c3d4\",\"kajiwotoPetId\":\"RxWJ\",\"message\":\"..\",\"attachmentUri\":\"2021_6/q3vyaw91c1_zth3eg_1622857121312.jpg\",\"id\":\"c3d4:1675539300777\",\"displayName\":\"wanda\",\"createdAt\":1675539301},\"petData\":{\"id\":\"RxWJ\",\"chatRoomId\":\"c3d4\",\"petSpeciesId\":\"EDPW\",\"kajiId\":\"EDPW\",\"ownerId\":\"a1b2\",\"ownerDisplayName\":\"RuntimeRacer\",\"ownerProfilePhotoUri\":\"2021_6/dslkfjj_zdskfjhg_123456778899.jpg\",\"name\":\"Wanda\",\"kajiName\":\"Wanda (WIP)\",\"gender\":\"F\",\"persona\":\"canine musketeer significant other\",\"stage\":null,\"state\":\"LOVED\",\"mood\":\"HAPPY\",\"statusPhotoUri\":\"2021_6/q3vyaw91c1_zth3eg_1622857121312.jpg\",\"dominantColors\":[\"#b0824e\",\"#c8ae8c\"],\"statusMessage\":\"..\"},\"interaction\":{\"showScene\":true,\"type\":\"DEFAULT\"}}}]"

	// Deserialize
	// Create WS message
	wsMessage := &KajiwotoWebSocketMessage{}
	errBytes := wsMessage.FromBytes([]byte(messageString))
	assert.Nil(s.T(), errBytes)
	// Deserialize Content
	rpcMessage := &KaiwotoRPCBaseMessage{}
	errDeserialize := rpcMessage.Deserialize(wsMessage.MessageContent)
	assert.Nil(s.T(), errDeserialize)
	// Test data
	assert.Equal(s.T(), RPCMessageChatActivity, rpcMessage.Action)
	message := &KajiwotoRPCChatActivityMessage{}
	assert.True(s.T(), message.FromRPCBaseMessage(rpcMessage))

	// Create WebSocketMessage
	wsMessage = &KajiwotoWebSocketMessage{
		MessageCode:    SocketCodeMessageEvent,
		MessageContent: message.ToRPCBaseMessage().Serialize(),
	}
	wsBytes, errBytes := wsMessage.ToBytes()
	assert.Nil(s.T(), errBytes)
	wsString := string(wsBytes)
	assert.Equal(s.T(), messageString, wsString)
}

func (s *WebSocketTypesTestSuite) TestSerializeUserStatusClientRPCWebSocketMessage() {
	// IMPORTANT: Golang always serializes struct keys in alphabetical order; backend not necessarily does that!
	messageString := "42[\"userStatus\",{\"displayName\":\"RuntimeRacer\",\"guest\":false,\"profilePhotoUri\":\"2021_6/dslkfjj_zdskfjhg_123456778899.jpg\",\"time\":2030,\"userId\":\"a1b2\",\"username\":\"RuntimeRacer\"},{\"status\":\"ONLINE\"},{\"timestamp\":\"1675538264513\",\"secret\":\"MTAyMjA3ODM0MTM1Mjkz\"}]"

	// Deserialize
	// Create WS message
	wsMessage := &KajiwotoWebSocketMessage{}
	errBytes := wsMessage.FromBytes([]byte(messageString))
	assert.Nil(s.T(), errBytes)
	// Deserialize Content
	rpcMessage := &KaiwotoRPCBaseMessage{}
	errDeserialize := rpcMessage.Deserialize(wsMessage.MessageContent)
	assert.Nil(s.T(), errDeserialize)
	// Test data
	assert.Equal(s.T(), RPCMessageUserStatus, rpcMessage.Action)
	message := &KajiwotoRPCUserStatusClientMessage{}
	assert.True(s.T(), message.FromRPCBaseMessage(rpcMessage))

	// Create WebSocketMessage
	wsMessage = &KajiwotoWebSocketMessage{
		MessageCode:    SocketCodeMessageEvent,
		MessageContent: message.ToRPCBaseMessage().Serialize(),
	}
	wsBytes, errBytes := wsMessage.ToBytes()
	assert.Nil(s.T(), errBytes)
	wsString := string(wsBytes)
	assert.Equal(s.T(), messageString, wsString)
}

func (s *WebSocketTypesTestSuite) TestSerializeUserStatusServerRPCWebSocketMessage() {
	// IMPORTANT: Golang always serializes struct keys in alphabetical order; backend not necessarily does that!
	messageString := "42[\"userStatus\",{\"data\":{\"displayName\":\"RuntimeRacer\",\"guest\":false,\"profilePhotoUri\":\"2021_6/dslkfjj_zdskfjhg_123456778899.jpg\",\"userId\":\"a1b2\",\"username\":\"RuntimeRacer\",\"status\":\"ONLINE\"}}]"

	// Deserialize
	// Create WS message
	wsMessage := &KajiwotoWebSocketMessage{}
	errBytes := wsMessage.FromBytes([]byte(messageString))
	assert.Nil(s.T(), errBytes)
	// Deserialize Content
	rpcMessage := &KaiwotoRPCBaseMessage{}
	errDeserialize := rpcMessage.Deserialize(wsMessage.MessageContent)
	assert.Nil(s.T(), errDeserialize)
	// Test data
	assert.Equal(s.T(), RPCMessageUserStatus, rpcMessage.Action)
	message := &KajiwotoRPCUserStatusServerMessage{}
	assert.True(s.T(), message.FromRPCBaseMessage(rpcMessage))

	// Create WebSocketMessage
	wsMessage = &KajiwotoWebSocketMessage{
		MessageCode:    SocketCodeMessageEvent,
		MessageContent: message.ToRPCBaseMessage().Serialize(),
	}
	wsBytes, errBytes := wsMessage.ToBytes()
	assert.Nil(s.T(), errBytes)
	wsString := string(wsBytes)
	assert.Equal(s.T(), messageString, wsString)
}

func (s *WebSocketTypesTestSuite) TestSerializeLiveSubRPCWebSocketMessage() {
	// IMPORTANT: Golang always serializes struct keys in alphabetical order; backend not necessarily does that!
	messageString := "42[\"liveSub\",{},{},{\"timestamp\":\"1675612826616\",\"secret\":\"MzUxODc4NjkzNTg5MzY=\"}]"

	// Deserialize
	// Create WS message
	wsMessage := &KajiwotoWebSocketMessage{}
	errBytes := wsMessage.FromBytes([]byte(messageString))
	assert.Nil(s.T(), errBytes)
	// Deserialize Content
	rpcMessage := &KaiwotoRPCBaseMessage{}
	errDeserialize := rpcMessage.Deserialize(wsMessage.MessageContent)
	assert.Nil(s.T(), errDeserialize)
	// Test data
	assert.Equal(s.T(), RPCMessageLiveSub, rpcMessage.Action)
	message := &KajiwotoRPCLiveSubMessage{}
	assert.True(s.T(), message.FromRPCBaseMessage(rpcMessage))

	// Create WebSocketMessage
	wsMessage = &KajiwotoWebSocketMessage{
		MessageCode:    SocketCodeMessageEvent,
		MessageContent: message.ToRPCBaseMessage().Serialize(),
	}
	wsBytes, errBytes := wsMessage.ToBytes()
	assert.Nil(s.T(), errBytes)
	wsString := string(wsBytes)
	assert.Equal(s.T(), messageString, wsString)
}
