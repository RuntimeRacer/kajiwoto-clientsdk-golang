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
	"fmt"
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
)

// Basic Message Handling type
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
	contentMatcher, errCompile := regexp.Compile(`(\d*)({.*})`)
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

// Message Content types
type KaiwotoWebSocketAuthRequest struct {
	ApiKey string `json:"api_key"`
}
type KaiwotoWebSocketAuthResponse struct {
	Sid string `json:"sid"`
}
