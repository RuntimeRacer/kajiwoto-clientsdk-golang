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

func (s *WebSocketTypesTestSuite) TestTypingRPCWebSocketMessage() {
	// IMPORTANT: Golang always serializes struct keys in ascending order; backend not necessarily does that!
	messageString := "42[\"typing\",{\"displayName\":\"RRacer2021\",\"guest\":false,\"profilePhotoUri\":\"2021_6/mwe1ntk2mj_zth3eg_1624667330826.jpg\",\"time\":2030,\"userId\":\"e8wz\",\"username\":\"RRacer2021\"},{\"chatRoomId\":\"vd8p\"},{\"timestamp\":\"1675538167859\",\"secret\":\"MTAyMjA3ODI4MjM5Mzk5\"}]"

	// Deserialize
	// Create WS message
	wsMessage := &KajiwotoWebSocketMessage{}
	errBytes := wsMessage.FromBytes([]byte(messageString))
	assert.Nil(s.T(), errBytes)
	// Deserialize Content
	rpcMessage := &KaiwotoRPCMessage{}
	errDeserialize := rpcMessage.Deserialize(wsMessage.MessageContent)
	assert.Nil(s.T(), errDeserialize)
	// Test data
	assert.Equal(s.T(), RPCMessageTyping, rpcMessage.Action)
	message := &KajiwotoRPCTypingMessage{}
	assert.True(s.T(), message.FromRPCMessage(rpcMessage))

	// Serialize
	// Create WebSocketMessage
	wsMessage = &KajiwotoWebSocketMessage{
		MessageCode:    SocketCodeMessageEvent,
		MessageContent: message.ToRPCMessage().Serialize(),
	}
	wsBytes, errBytes := wsMessage.ToBytes()
	assert.Nil(s.T(), errBytes)
	wsString := string(wsBytes)
	assert.Equal(s.T(), messageString, wsString)
}

func (s *WebSocketTypesTestSuite) TestSubscribeRPCWebSocketMessage() {
	// IMPORTANT: Golang always serializes struct keys in ascending order; backend not necessarily does that!
	messageString := "42[\"subscribe\",{\"displayName\":\"RRacer2021\",\"guest\":false,\"profilePhotoUri\":\"2021_6/mwe1ntk2mj_zth3eg_1624667330826.jpg\",\"time\":2030,\"userId\":\"e8wz\",\"username\":\"RRacer2021\"},{\"chatRoomIds\":[\"vd8p\"],\"kajiId\":null},{\"timestamp\":\"1675538034488\",\"secret\":\"MTY3NTUzODAzNDQ4OA==\"}]"

	// Deserialize
	// Create WS message
	wsMessage := &KajiwotoWebSocketMessage{}
	errBytes := wsMessage.FromBytes([]byte(messageString))
	assert.Nil(s.T(), errBytes)
	// Deserialize Content
	rpcMessage := &KaiwotoRPCMessage{}
	errDeserialize := rpcMessage.Deserialize(wsMessage.MessageContent)
	assert.Nil(s.T(), errDeserialize)
	// Test data
	assert.Equal(s.T(), RPCMessageSubscribe, rpcMessage.Action)
	message := &KajiwotoRPCSubscribeMessage{}
	assert.True(s.T(), message.FromRPCMessage(rpcMessage))

	// Create WebSocketMessage
	wsMessage = &KajiwotoWebSocketMessage{
		MessageCode:    SocketCodeMessageEvent,
		MessageContent: message.ToRPCMessage().Serialize(),
	}
	wsBytes, errBytes := wsMessage.ToBytes()
	assert.Nil(s.T(), errBytes)
	wsString := string(wsBytes)
	assert.Equal(s.T(), messageString, wsString)
}

func (s *WebSocketTypesTestSuite) TestChatEnterRPCWebSocketMessage() {
	// IMPORTANT: Golang always serializes struct keys in ascending order; backend not necessarily does that!
	messageString := "42[\"chatEnter\",{\"displayName\":\"RRacer2021\",\"guest\":false,\"profilePhotoUri\":\"2021_6/mwe1ntk2mj_zth3eg_1624667330826.jpg\",\"time\":2030,\"userId\":\"e8wz\",\"username\":\"RRacer2021\"},{\"chatRoomId\":\"vd8p\",\"lastMessages\":[{\"createdAt\":1675477983263,\"message\":\"/say good night my man\"},{\"createdAt\":1675477879022,\"message\":\"*whispers* sweet dreams my pretty mink\"}],\"isPreviewRoom\":false},{\"timestamp\":\"1675538039386\",\"secret\":\"MTY3NTUzODAzOTM4Ng==\"}]"

	// Deserialize
	// Create WS message
	wsMessage := &KajiwotoWebSocketMessage{}
	errBytes := wsMessage.FromBytes([]byte(messageString))
	assert.Nil(s.T(), errBytes)
	// Deserialize Content
	rpcMessage := &KaiwotoRPCMessage{}
	errDeserialize := rpcMessage.Deserialize(wsMessage.MessageContent)
	assert.Nil(s.T(), errDeserialize)
	// Test data
	assert.Equal(s.T(), RPCMessageChatEnter, rpcMessage.Action)
	message := &KajiwotoRPCChatEnterMessage{}
	assert.True(s.T(), message.FromRPCMessage(rpcMessage))

	// Create WebSocketMessage
	wsMessage = &KajiwotoWebSocketMessage{
		MessageCode:    SocketCodeMessageEvent,
		MessageContent: message.ToRPCMessage().Serialize(),
	}
	wsBytes, errBytes := wsMessage.ToBytes()
	assert.Nil(s.T(), errBytes)
	wsString := string(wsBytes)
	assert.Equal(s.T(), messageString, wsString)
}

func (s *WebSocketTypesTestSuite) TestChatSendRPCWebSocketMessage() {
	// IMPORTANT: Golang always serializes struct keys in ascending order; backend not necessarily does that!
	messageString := "42[\"chatSend\",{\"displayName\":\"RRacer2021\",\"guest\":false,\"profilePhotoUri\":\"2021_6/mwe1ntk2mj_zth3eg_1624667330826.jpg\",\"time\":2030,\"userId\":\"e8wz\",\"username\":\"RRacer2021\"},{\"message\":{\"id\":\"vd8p:1675538262207\",\"chatRoomId\":\"vd8p\",\"userId\":\"e8wz\",\"message\":\"Hey my sweet *smiles*\",\"attachmentUri\":null},\"roomVersionNumber\":1675538034,\"roomSocketIds\":[\"emCCdEmKKsm2aPLCABAN\"]},{\"timestamp\":\"1675538262207\",\"secret\":\"MTAyMjA3ODMzOTk0NjI3\"}]"

	// Deserialize
	// Create WS message
	wsMessage := &KajiwotoWebSocketMessage{}
	errBytes := wsMessage.FromBytes([]byte(messageString))
	assert.Nil(s.T(), errBytes)
	// Deserialize Content
	rpcMessage := &KaiwotoRPCMessage{}
	errDeserialize := rpcMessage.Deserialize(wsMessage.MessageContent)
	assert.Nil(s.T(), errDeserialize)
	// Test data
	assert.Equal(s.T(), RPCMessageChatSend, rpcMessage.Action)
	message := &KajiwotoRPCChatSendMessage{}
	assert.True(s.T(), message.FromRPCMessage(rpcMessage))

	// Create WebSocketMessage
	wsMessage = &KajiwotoWebSocketMessage{
		MessageCode:    SocketCodeMessageEvent,
		MessageContent: message.ToRPCMessage().Serialize(),
	}
	wsBytes, errBytes := wsMessage.ToBytes()
	assert.Nil(s.T(), errBytes)
	wsString := string(wsBytes)
	assert.Equal(s.T(), messageString, wsString)
}

func (s *WebSocketTypesTestSuite) TestChatActivityRPCWebSocketMessageVariant1() {
	// IMPORTANT: Golang always serializes struct keys in ascending order; backend not necessarily does that!
	messageString := "42[\"chatActivity\",{\"data\":{\"action\":\"join-room\",\"chatRoomId\":\"vd8p\",\"petData\":{\"id\":\"RxWJ\",\"chatRoomId\":\"vd8p\",\"petSpeciesId\":\"EDPW\",\"kajiId\":\"EDPW\",\"ownerId\":\"e8wz\",\"ownerDisplayName\":\"RRacer2021\",\"ownerProfilePhotoUri\":\"2021_6/mwe1ntk2mj_zth3eg_1624667330826.jpg\",\"name\":\"Wanda\",\"kajiName\":\"Wanda (WIP)\",\"gender\":\"F\",\"persona\":\"canine musketeer significant other\",\"stage\":null,\"state\":\"DEFAULT\",\"mood\":\"DEFAULT\",\"statusPhotoUri\":\"2021_6/tm9ybwfsxz_zth3eg_1622766488811.png\",\"dominantColors\":[\"#dc9744\",\"#fcd49c\"],\"statusMessage\":\"..\"},\"channel\":{\"v\":1675538034,\"list\":[{\"id\":\"e8wz\",\"guestId\":\"OTMuMTk5LjEyOS4yMTg=*\",\"socketIds\":[\"emCCdEmKKsm2aPLCABAN\"],\"guest\":false,\"displayName\":\"RRacer2021\",\"username\":\"RRacer2021\",\"profilePhotoUri\":\"2021_6/mwe1ntk2mj_zth3eg_1624667330826.jpg\"}]}}}]"

	// Deserialize
	// Create WS message
	wsMessage := &KajiwotoWebSocketMessage{}
	errBytes := wsMessage.FromBytes([]byte(messageString))
	assert.Nil(s.T(), errBytes)
	// Deserialize Content
	rpcMessage := &KaiwotoRPCMessage{}
	errDeserialize := rpcMessage.Deserialize(wsMessage.MessageContent)
	assert.Nil(s.T(), errDeserialize)
	// Test data
	assert.Equal(s.T(), RPCMessageChatActivity, rpcMessage.Action)
	message := &KajiwotoRPCChatActivityMessage{}
	assert.True(s.T(), message.FromRPCMessage(rpcMessage))

	// Create WebSocketMessage
	wsMessage = &KajiwotoWebSocketMessage{
		MessageCode:    SocketCodeMessageEvent,
		MessageContent: message.ToRPCMessage().Serialize(),
	}
	wsBytes, errBytes := wsMessage.ToBytes()
	assert.Nil(s.T(), errBytes)
	wsString := string(wsBytes)
	assert.Equal(s.T(), messageString, wsString)
}

func (s *WebSocketTypesTestSuite) TestChatActivityRPCWebSocketMessageVariant2() {
	// IMPORTANT: Golang always serializes struct keys in ascending order; backend not necessarily does that!
	messageString := "42[\"chatActivity\",{\"data\":{\"action\":\"activity\",\"chatRoomId\":\"vd8p\",\"activity\":{\"type\":\"TYPING\",\"userId\":\"e8wz\",\"displayName\":\"RRacer2021\",\"activityAt\":1675538172488}}}]"

	// Deserialize
	// Create WS message
	wsMessage := &KajiwotoWebSocketMessage{}
	errBytes := wsMessage.FromBytes([]byte(messageString))
	assert.Nil(s.T(), errBytes)
	// Deserialize Content
	rpcMessage := &KaiwotoRPCMessage{}
	errDeserialize := rpcMessage.Deserialize(wsMessage.MessageContent)
	assert.Nil(s.T(), errDeserialize)
	// Test data
	assert.Equal(s.T(), RPCMessageChatActivity, rpcMessage.Action)
	message := &KajiwotoRPCChatActivityMessage{}
	assert.True(s.T(), message.FromRPCMessage(rpcMessage))

	// Create WebSocketMessage
	wsMessage = &KajiwotoWebSocketMessage{
		MessageCode:    SocketCodeMessageEvent,
		MessageContent: message.ToRPCMessage().Serialize(),
	}
	wsBytes, errBytes := wsMessage.ToBytes()
	assert.Nil(s.T(), errBytes)
	wsString := string(wsBytes)
	assert.Equal(s.T(), messageString, wsString)
}

func (s *WebSocketTypesTestSuite) TestChatActivityRPCWebSocketMessageVariant3() {
	// IMPORTANT: Golang always serializes struct keys in ascending order; backend not necessarily does that!
	messageString := "42[\"chatActivity\",{\"data\":{\"action\":\"message\",\"chatRoomId\":\"vd8p\",\"message\":{\"clientId\":\"vd8p:1675538262207\",\"chatRoomId\":\"vd8p\",\"message\":\"Hey my sweet *smiles*\",\"attachmentUri\":null,\"id\":\"vd8p:1675538262207\",\"userId\":\"e8wz\",\"username\":\"RRacer2021\",\"displayName\":\"RRacer2021\",\"profilePhotoUri\":\"2021_6/mwe1ntk2mj_zth3eg_1624667330826.jpg\",\"createdAt\":1675538261},\"channel\":{\"v\":1675538034,\"list\":null},\"socketIds\":[\"emCCdEmKKsm2aPLCABAN\"]}}]"

	// Deserialize
	// Create WS message
	wsMessage := &KajiwotoWebSocketMessage{}
	errBytes := wsMessage.FromBytes([]byte(messageString))
	assert.Nil(s.T(), errBytes)
	// Deserialize Content
	rpcMessage := &KaiwotoRPCMessage{}
	errDeserialize := rpcMessage.Deserialize(wsMessage.MessageContent)
	assert.Nil(s.T(), errDeserialize)
	// Test data
	assert.Equal(s.T(), RPCMessageChatActivity, rpcMessage.Action)
	message := &KajiwotoRPCChatActivityMessage{}
	assert.True(s.T(), message.FromRPCMessage(rpcMessage))

	// Create WebSocketMessage
	wsMessage = &KajiwotoWebSocketMessage{
		MessageCode:    SocketCodeMessageEvent,
		MessageContent: message.ToRPCMessage().Serialize(),
	}
	wsBytes, errBytes := wsMessage.ToBytes()
	assert.Nil(s.T(), errBytes)
	wsString := string(wsBytes)
	assert.Equal(s.T(), messageString, wsString)
}
