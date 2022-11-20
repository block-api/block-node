// Copyright 2022 The block-node Authors
// This file is part of the block-node library.
//
// The block-node library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The block-node library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the block-node library. If not, see <http://www.gnu.org/licenses/>.

// Package function
package function

import (
	"encoding/json"
)

type Request struct {
	FromID string `json:"fid,omitempty"`
	Body   any    `json:"b,omitempty"`
}

func NewRequest(fromID string, body any) Request {
	return Request{
		FromID: fromID,
		Body:   body,
	}
}

func DecodeBody[T any](body any) T {
	var out T

	bodyBytes, _ := json.Marshal(body)

	_ = json.Unmarshal(bodyBytes, &out)
	// fmt.Println(e)
	// fmt.Println(out)
	return out
}
