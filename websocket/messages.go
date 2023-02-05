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
	"encoding/base64"
	"strconv"
	"time"
)

func CreateKajiwotoWebSocketEventMessage(rpcMessage KajiwotoRPCMessage) *KajiwotoWebSocketMessage {
	wsMessage := &KajiwotoWebSocketMessage{
		MessageCode:    SocketCodeMessageEvent,
		MessageContent: rpcMessage.ToRPCBaseMessage().Serialize(),
	}
	return wsMessage
}

func createMessageSecret() KajiwotoRPCSecret {
	// Build timestamp secret
	now := time.Now()
	milis := now.UnixMilli()

	milis = 1675606413185
	timestamp := strconv.FormatInt(milis, 10)

	// Get chars at pos 7 & 8 + their int representation
	ts7, _ := strconv.Atoi(string(timestamp[7]))
	ts8, _ := strconv.Atoi(string(timestamp[8]))

	var multiplierString string
	if ts8%2 == 1 {
		multiplierString = strconv.Itoa(ts7) + "1"
	} else {
		multiplierString = strconv.Itoa(ts8) + "1"
	}
	multiplier, _ := strconv.Atoi(multiplierString)

	secret := milis * int64(multiplier)
	secretValue := strconv.FormatInt(secret, 10)
	secretB64 := base64.StdEncoding.EncodeToString([]byte(secretValue))
	return KajiwotoRPCSecret{
		Timestamp: timestamp,
		Secret:    secretB64,
	}
}
