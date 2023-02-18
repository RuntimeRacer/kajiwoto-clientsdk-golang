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
	"errors"
	"fmt"
	"github.com/mitchellh/mapstructure"
	"regexp"
)

const (
	// Data Frames - see: https://www.rfc-editor.org/rfc/rfc6455#section-5.6
	DataFrameText   = "1"
	DataFrameBinary = "2"

	// WS Codes - see: https://stackoverflow.com/questions/24564877/what-do-these-numbers-mean-in-socket-io-payload
	// Basic codes
	SocketCodeOpen  = "0"
	SocketCodeClose = "1"
	SocketCodePing  = "2"
	SocketCodePong  = "3"

	// Complex Codes
	SocketCodeMessageConnect    = "40"
	SocketCodeMessageDisconnect = "41"
	SocketCodeMessageEvent      = "42"
	SocketCodeMessageAck        = "43"
	SocketCodeMessageError      = "44"

	// RPC Message Types
	RPCMessageChatActivity = "chatActivity"
	RPCMessageChatEnter    = "chatEnter"
	RPCMessageChatLeave    = "chatLeave"
	RPCMessageChatSend     = "chatSend"
	RPCMessageChatSubmit   = "chatSubmit"
	RPCMessageLiveSub      = "liveSub"
	RPCMessageLogin        = "login"
	RPCMessageSubscribe    = "subscribe"
	RPCMessageUserStatus   = "userStatus"
	RPCMessageTyping       = "typing"

	// RPC Message Subtypes
	ChatActivitySubActivity = "activity"
	ChatActivityMessage     = "message"
	ChatActivityPetMessage  = "petMessage"
	ChatActivityJoinRoom    = "join-room"
)

var (
	ErrUnableToHandleMessage = errors.New("unable to handle message")
)

// Basic WebSocket Message Handling types
type KajiwotoWebSocketMessage struct {
	MessageCode    string
	MessageContent interface{}
}

func (k *KajiwotoWebSocketMessage) ToBytes() ([]byte, error) {
	messageBytes := []byte(k.MessageCode)
	if k.MessageContent != nil {
		messageContentBytes, errMarshal := json.Marshal(k.MessageContent)
		if errMarshal != nil {
			return nil, errMarshal
		}
		messageBytes = append(messageBytes, messageContentBytes...)
	}
	return messageBytes, nil
}

func (k *KajiwotoWebSocketMessage) FromBytes(bytes []byte) error {
	contentMatcher, errCompile := regexp.Compile(`(\d*)({.*}|\[.*\])`)
	if errCompile != nil {
		return errCompile
	}
	matches := contentMatcher.FindAllSubmatch(bytes, -1)
	if matches != nil {
		if len(matches[0]) != 3 {
			return fmt.Errorf("unable to parse message. message data: %v", string(bytes))
		}
		// Build from regex result
		k.MessageCode = string(matches[0][1])
		k.MessageContent = matches[0][2] // Unmarshal in response handler
	} else {
		// Assume message has no content, just a code, to be evaluated in handlers
		k.MessageCode = string(bytes)
	}

	return nil
}

// WebSocket Message Content types
type KaiwotoWebSocketAuthRequest struct {
	ApiKey string `json:"api_key"`
}
type KaiwotoWebSocketAuthResponse struct {
	Sid     string `json:"sid,omitempty"`
	Message string `json:"message,omitempty"`
}

// Basic RPC message handling types
type KajiwotoRPCMessage interface {
	ToRPCBaseMessage() *KaiwotoRPCBaseMessage
	FromRPCBaseMessage(message *KaiwotoRPCBaseMessage) bool
}

type KaiwotoRPCBaseMessage struct {
	Action  string
	Payload []interface{}
}

func (k *KaiwotoRPCBaseMessage) Serialize() []interface{} {
	messageData := make([]interface{}, 0)
	messageData = append(messageData, k.Action)
	messageData = append(messageData, k.Payload...)
	return messageData
}

func (k *KaiwotoRPCBaseMessage) Deserialize(messageData interface{}) error {
	var rpcMessageParts []interface{}
	// check if already array
	rpcMessageParts, okCastArray := messageData.([]interface{})
	if !okCastArray {
		// check if bytes
		messageBytes, okCastBytes := messageData.([]byte)
		if !okCastBytes {
			return fmt.Errorf("cannot deserialize data into rpc message")
		}
		var err error
		if rpcMessageParts, err = k.DeserializeFromBytes(messageBytes); err != nil {
			return err
		}
	}
	// Fill in fields
	if len(rpcMessageParts) > 0 {
		k.Action = rpcMessageParts[0].(string)
	}
	if len(rpcMessageParts) > 1 {
		k.Payload = rpcMessageParts[1:]
	}
	return nil
}

func (k *KaiwotoRPCBaseMessage) DeserializeFromBytes(messageBytes []byte) ([]interface{}, error) {
	// Try to decode message content into Array
	rpcMessageParts := make([]interface{}, 0)
	errUnmarshal := json.Unmarshal(messageBytes, &rpcMessageParts)
	if errUnmarshal != nil {
		return nil, errUnmarshal
	}
	return rpcMessageParts, nil
}

func (k *KaiwotoRPCBaseMessage) FetchDataFromPayload(output interface{}, ignoreUnset bool) bool {
	// Create a typesafe decoder
	decoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		Metadata:   nil,
		Result:     output,
		ErrorUnset: !ignoreUnset,
	})
	if err != nil {
		return false
	}
	// Decode into element
	for _, payloadElem := range k.Payload {
		if err = decoder.Decode(payloadElem); err == nil {
			return true
		}
	}
	return false
}

// RPC Message Types
type KajiwotoRPCTypingMessage struct {
	UserData   KajiwotoRPCUserData
	ChatRoomId KajiwotoRPCChatRoomId
	Secret     KajiwotoRPCSecret
}

func (k *KajiwotoRPCTypingMessage) ToRPCBaseMessage() *KaiwotoRPCBaseMessage {
	return &KaiwotoRPCBaseMessage{
		Action: RPCMessageTyping,
		Payload: []interface{}{
			k.UserData,
			k.ChatRoomId,
			k.Secret,
		},
	}
}
func (k *KajiwotoRPCTypingMessage) FromRPCBaseMessage(message *KaiwotoRPCBaseMessage) bool {
	if message.Action != RPCMessageTyping {
		return false
	}
	if fetch := message.FetchDataFromPayload(&k.UserData, false); !fetch {
		// do nothing
	}
	if fetch := message.FetchDataFromPayload(&k.ChatRoomId, false); !fetch {
		// do nothing
	}
	if fetch := message.FetchDataFromPayload(&k.Secret, false); !fetch {
		// do nothing
	}
	return true
}

type KajiwotoRPCLoginMessage struct {
	UserData   KajiwotoRPCUserData
	UserStatus KajiwotoRPCUserStatus
	Secret     KajiwotoRPCSecret
}

func (k *KajiwotoRPCLoginMessage) ToRPCBaseMessage() *KaiwotoRPCBaseMessage {
	return &KaiwotoRPCBaseMessage{
		Action: RPCMessageLogin,
		Payload: []interface{}{
			k.UserData,
			k.UserStatus,
			k.Secret,
		},
	}
}
func (k *KajiwotoRPCLoginMessage) FromRPCBaseMessage(message *KaiwotoRPCBaseMessage) bool {
	if message.Action != RPCMessageLogin {
		return false
	}
	if fetch := message.FetchDataFromPayload(&k.UserData, false); !fetch {
		// do nothing
	}
	if fetch := message.FetchDataFromPayload(&k.UserStatus, false); !fetch {
		// do nothing
	}
	if fetch := message.FetchDataFromPayload(&k.Secret, false); !fetch {
		// do nothing
	}
	return true
}

type KajiwotoRPCSubscribeMessage struct {
	UserData      KajiwotoRPCUserData
	SubscribeArgs KajiwotoRPCSubscribeArgs
	Secret        KajiwotoRPCSecret
}

func (k *KajiwotoRPCSubscribeMessage) ToRPCBaseMessage() *KaiwotoRPCBaseMessage {
	return &KaiwotoRPCBaseMessage{
		Action: RPCMessageSubscribe,
		Payload: []interface{}{
			k.UserData,
			k.SubscribeArgs,
			k.Secret,
		},
	}
}
func (k *KajiwotoRPCSubscribeMessage) FromRPCBaseMessage(message *KaiwotoRPCBaseMessage) bool {
	if message.Action != RPCMessageSubscribe {
		return false
	}
	if fetch := message.FetchDataFromPayload(&k.UserData, false); !fetch {
		// do nothing
	}
	if fetch := message.FetchDataFromPayload(&k.SubscribeArgs, false); !fetch {
		// do nothing
	}
	if fetch := message.FetchDataFromPayload(&k.Secret, false); !fetch {
		// do nothing
	}
	return true
}

type KajiwotoRPCChatEnterMessage struct {
	UserData     KajiwotoRPCUserData
	ChatroomData KajiwotoRPCChatRoomData
	Secret       KajiwotoRPCSecret
}

func (k *KajiwotoRPCChatEnterMessage) ToRPCBaseMessage() *KaiwotoRPCBaseMessage {
	return &KaiwotoRPCBaseMessage{
		Action: RPCMessageChatEnter,
		Payload: []interface{}{
			k.UserData,
			k.ChatroomData,
			k.Secret,
		},
	}
}
func (k *KajiwotoRPCChatEnterMessage) FromRPCBaseMessage(message *KaiwotoRPCBaseMessage) bool {
	if message.Action != RPCMessageChatEnter {
		return false
	}
	if fetch := message.FetchDataFromPayload(&k.UserData, false); !fetch {
		// do nothing
	}
	if fetch := message.FetchDataFromPayload(&k.ChatroomData, false); !fetch {
		// do nothing
	}
	if fetch := message.FetchDataFromPayload(&k.Secret, false); !fetch {
		// do nothing
	}
	return true
}

type KajiwotoRPCChatSendMessage struct {
	UserData     KajiwotoRPCUserData
	ChatSendData KajiwotoRPCChatMessageCreate
	Secret       KajiwotoRPCSecret
}

func (k *KajiwotoRPCChatSendMessage) ToRPCBaseMessage() *KaiwotoRPCBaseMessage {
	return &KaiwotoRPCBaseMessage{
		Action: RPCMessageChatSend,
		Payload: []interface{}{
			k.UserData,
			k.ChatSendData,
			k.Secret,
		},
	}
}
func (k *KajiwotoRPCChatSendMessage) FromRPCBaseMessage(message *KaiwotoRPCBaseMessage) bool {
	if message.Action != RPCMessageChatSend {
		return false
	}
	if fetch := message.FetchDataFromPayload(&k.UserData, false); !fetch {
		// do nothing
	}
	if fetch := message.FetchDataFromPayload(&k.ChatSendData, false); !fetch {
		// do nothing
	}
	if fetch := message.FetchDataFromPayload(&k.Secret, false); !fetch {
		// do nothing
	}
	return true
}

type KajiwotoRPCChatSubmitMessage struct {
	UserData       KajiwotoRPCUserData
	ChatSubmitData KajiwotoRPCChatSubmitData
	Secret         KajiwotoRPCSecret
}

func (k *KajiwotoRPCChatSubmitMessage) ToRPCBaseMessage() *KaiwotoRPCBaseMessage {
	return &KaiwotoRPCBaseMessage{
		Action: RPCMessageChatSubmit,
		Payload: []interface{}{
			k.UserData,
			k.ChatSubmitData,
			k.Secret,
		},
	}
}
func (k *KajiwotoRPCChatSubmitMessage) FromRPCBaseMessage(message *KaiwotoRPCBaseMessage) bool {
	if message.Action != RPCMessageChatSubmit {
		return false
	}
	if fetch := message.FetchDataFromPayload(&k.UserData, false); !fetch {
		// do nothing
	}
	if fetch := message.FetchDataFromPayload(&k.ChatSubmitData, false); !fetch {
		// do nothing
	}
	if fetch := message.FetchDataFromPayload(&k.Secret, false); !fetch {
		// do nothing
	}
	return true
}

type KajiwotoRPCChatLeaveMessage struct {
	Field1   KajiwotoRPCEmptyObject
	ChatRoom KajiwotoRPCChatRoomId
	Secret   KajiwotoRPCSecret
}

func (k *KajiwotoRPCChatLeaveMessage) ToRPCBaseMessage() *KaiwotoRPCBaseMessage {
	return &KaiwotoRPCBaseMessage{
		Action: RPCMessageChatLeave,
		Payload: []interface{}{
			k.Field1,
			k.ChatRoom,
			k.Secret,
		},
	}
}
func (k *KajiwotoRPCChatLeaveMessage) FromRPCBaseMessage(message *KaiwotoRPCBaseMessage) bool {
	if message.Action != RPCMessageChatLeave {
		return false
	}
	if fetch := message.FetchDataFromPayload(&k.Field1, false); !fetch {
		// do nothing
	}
	if fetch := message.FetchDataFromPayload(&k.ChatRoom, false); !fetch {
		// do nothing
	}
	if fetch := message.FetchDataFromPayload(&k.Secret, false); !fetch {
		// do nothing
	}
	return true
}

type KajiwotoRPCChatActivityMessage struct {
	ActivityData KajiwotoRPCChatActivityData
}

func (k *KajiwotoRPCChatActivityMessage) ToRPCBaseMessage() *KaiwotoRPCBaseMessage {
	return &KaiwotoRPCBaseMessage{
		Action: RPCMessageChatActivity,
		Payload: []interface{}{
			k.ActivityData,
		},
	}
}
func (k *KajiwotoRPCChatActivityMessage) FromRPCBaseMessage(message *KaiwotoRPCBaseMessage) bool {
	if message.Action != RPCMessageChatActivity {
		return false
	}
	if fetch := message.FetchDataFromPayload(&k.ActivityData, true); !fetch {
		// do nothing
	}
	return true
}

type KajiwotoRPCUserStatusClientMessage struct {
	UserData   KajiwotoRPCUserData
	UserStatus KajiwotoRPCUserStatus
	Secret     KajiwotoRPCSecret
}

func (k *KajiwotoRPCUserStatusClientMessage) ToRPCBaseMessage() *KaiwotoRPCBaseMessage {
	return &KaiwotoRPCBaseMessage{
		Action: RPCMessageUserStatus,
		Payload: []interface{}{
			k.UserData,
			k.UserStatus,
			k.Secret,
		},
	}
}
func (k *KajiwotoRPCUserStatusClientMessage) FromRPCBaseMessage(message *KaiwotoRPCBaseMessage) bool {
	if message.Action != RPCMessageUserStatus {
		return false
	}
	if fetch := message.FetchDataFromPayload(&k.UserData, false); !fetch {
		// do nothing
	}
	if fetch := message.FetchDataFromPayload(&k.UserStatus, false); !fetch {
		// do nothing
	}
	if fetch := message.FetchDataFromPayload(&k.Secret, false); !fetch {
		// do nothing
	}
	return true
}

type KajiwotoRPCUserStatusServerMessage struct {
	StatusData KajiwotoRPCUserStatusData
}

func (k *KajiwotoRPCUserStatusServerMessage) ToRPCBaseMessage() *KaiwotoRPCBaseMessage {
	return &KaiwotoRPCBaseMessage{
		Action: RPCMessageUserStatus,
		Payload: []interface{}{
			k.StatusData,
		},
	}
}
func (k *KajiwotoRPCUserStatusServerMessage) FromRPCBaseMessage(message *KaiwotoRPCBaseMessage) bool {
	if message.Action != RPCMessageUserStatus {
		return false
	}
	if fetch := message.FetchDataFromPayload(&k.StatusData, true); !fetch {
		// do nothing
	}
	return true
}

type KajiwotoRPCLiveSubMessage struct {
	Field1 KajiwotoRPCEmptyObject // No idea what these do yet, but they're always empty on initial call
	Field2 KajiwotoRPCEmptyObject // No idea what these do yet, but they're always empty on initial call
	Secret KajiwotoRPCSecret
}

func (k *KajiwotoRPCLiveSubMessage) ToRPCBaseMessage() *KaiwotoRPCBaseMessage {
	return &KaiwotoRPCBaseMessage{
		Action: RPCMessageLiveSub,
		Payload: []interface{}{
			k.Field1,
			k.Field2,
			k.Secret,
		},
	}
}
func (k *KajiwotoRPCLiveSubMessage) FromRPCBaseMessage(message *KaiwotoRPCBaseMessage) bool {
	if message.Action != RPCMessageLiveSub {
		return false
	}
	if fetch := message.FetchDataFromPayload(&k.Field1, false); !fetch {
		// do nothing
	}
	if fetch := message.FetchDataFromPayload(&k.Field2, false); !fetch {
		// do nothing
	}
	if fetch := message.FetchDataFromPayload(&k.Secret, false); !fetch {
		// do nothing
	}
	return true
}

// RPC Message Content types
type KajiwotoRPCUserData struct {
	DisplayName     string  `json:"displayName"`
	Guest           bool    `json:"guest"`
	ProfilePhotoUri *string `json:"profilePhotoUri"`
	Time            int     `json:"time"`
	UserID          string  `json:"userId"`
	Username        string  `json:"username"`
}

type KajiwotoRPCStatusUserData struct {
	DisplayName     string  `json:"displayName"`
	Guest           bool    `json:"guest"`
	ProfilePhotoUri *string `json:"profilePhotoUri"`
	UserID          string  `json:"userId"`
	Username        string  `json:"username"`
	Status          string  `json:"status"`
}

type KajiwotoRPCUserStatus struct {
	FriendIDs []string `json:"friendIds,omitempty"`
	Status    string   `json:"status"`
}

type KajiwotoRPCChatRoomId struct {
	ChatRoomId string `json:"chatRoomId"`
}

type KajiwotoRPCUserStatusData struct {
	Data KajiwotoRPCStatusUserData `json:"data"`
}

type KajiwotoRPCChatActivityData struct {
	Data KajiwotoRPCChatActivity `json:"data"`
}

type KajiwotoRPCChatActivity struct {
	Action      string                              `json:"action"`
	ChatRoomId  string                              `json:"chatRoomId"`
	EventType   *string                             `json:"eventType,omitempty"`
	Message     *KajiwotoRPCChatActivitySubMessage  `json:"message,omitempty"`
	Activity    *KajiwotoRPCChatActivitySubActivity `json:"activity,omitempty"`
	PetData     *KajiwotoRPCChatActivityPetData     `json:"petData,omitempty"`
	Channel     *KajiwotoRPCChatActivityChannel     `json:"channel,omitempty"`
	Interaction *KajiwotoRPCChatActivityInteraction `json:"interaction,omitempty"`
	SocketIds   []string                            `json:"socketIds,omitempty"`
}

type KajiwotoRPCChatActivityInteraction struct {
	ShowScene bool   `json:"showScene"`
	Type      string `json:"type"`
}

type KajiwotoRPCChatActivitySubMessage struct {
	ClientId        string  `json:"clientId,omitempty"`
	ChatRoomId      string  `json:"chatRoomId"`
	KajiwotoPetId   string  `json:"kajiwotoPetId,omitempty"`
	Message         string  `json:"message"`
	AttachmentUri   *string `json:"attachmentUri"`
	Id              string  `json:"id"`
	UserId          string  `json:"userId,omitempty"`
	UserName        string  `json:"username,omitempty"`
	DisplayName     string  `json:"displayName,omitempty"`
	ProfilePhotoUri *string `json:"profilePhotoUri,omitempty"`
	CreatedAt       uint64  `json:"createdAt"`
}

type KajiwotoRPCChatActivitySubActivity struct {
	Type        string `json:"type"`
	UserId      string `json:"userId"`
	DisplayName string `json:"displayName"`
	ActivityAt  uint64 `json:"activityAt"`
}

type KajiwotoRPCChatActivityPetData struct {
	Id                   string   `json:"id"`
	ChatRoomId           string   `json:"chatRoomId"`
	PetSpeciesId         string   `json:"petSpeciesId"`
	KajiId               string   `json:"kajiId"`
	OwnerId              string   `json:"ownerId"`
	OwnerDisplayName     string   `json:"ownerDisplayName"`
	OwnerProfilePhotoUri *string  `json:"ownerProfilePhotoUri"`
	Name                 string   `json:"name"`
	KajiName             string   `json:"kajiName"`
	Gender               string   `json:"gender"`
	Persona              string   `json:"persona"`
	Stage                *string  `json:"stage"`
	State                string   `json:"state"`
	Mood                 string   `json:"mood"`
	StatusPhotoUri       *string  `json:"statusPhotoUri"`
	DominantColors       []string `json:"dominantColors"`
	StatusMessage        string   `json:"statusMessage"`
}

type KajiwotoRPCChatActivityChannel struct {
	V    uint64                               `json:"v"`              // Channel version
	List []KajiwotoRPCChatActivityChannelUser `json:"list,omitempty"` // Channel user list
}

type KajiwotoRPCChatActivityChannelUser struct {
	Id              string   `json:"id"`
	GuestId         string   `json:"guestId"`
	SocketIds       []string `json:"socketIds"`
	Guest           bool     `json:"guest"`
	DisplayName     string   `json:"displayName"`
	Username        string   `json:"username"`
	ProfilePhotoUri *string  `json:"profilePhotoUri"`
}

type KajiwotoRPCChatRoomData struct {
	ChatRoomId    string                   `json:"chatRoomId"`
	LastMessages  []KajiwotoRPCChatMessage `json:"lastMessages"`
	IsPreviewRoom bool                     `json:"isPreviewRoom"`
}

type KajiwotoRPCChatSubmitData struct {
	ChatRoomId   string                        `json:"chatRoomId"`
	Messages     []string                      `json:"messages"`
	Role         KajiwotoRPCChatSubmitDataRole `json:"role"`
	Emoji        *string                       `json:"emoji"`
	EmojiSceneId *string                       `json:"emojiSceneId"`
	Platform     string                        `json:"platform"`
}

type KajiwotoRPCChatSubmitDataRole struct {
}

type KajiwotoRPCEmptyObject struct {
}

type KajiwotoRPCChatMessage struct {
	CreatedAt uint64 `json:"createdAt"`
	Message   string `json:"message"`
}

type KajiwotoRPCChatMessageCreate struct {
	Message           KajiwotoRPCChatMessageCreateData `json:"message"`
	RoomVersionNumber int64                            `json:"roomVersionNumber"`
	RoomSocketIds     []string                         `json:"roomSocketIds"`
}

type KajiwotoRPCChatMessageCreateData struct {
	Id            string  `json:"id"`
	ChatRoomId    string  `json:"chatRoomId"`
	UserID        string  `json:"userId"`
	Message       string  `json:"message"`
	AttachmentUri *string `json:"attachmentUri"`
}

type KajiwotoRPCSecret struct {
	Timestamp string `json:"timestamp"`
	Secret    string `json:"secret"`
}

type KajiwotoRPCSubscribeArgs struct {
	ChatRoomIds []string `json:"chatRoomIds"`
	KajiId      *string  `json:"kajiId"`
}
