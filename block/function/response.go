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

type Response struct {
	Status int    `json:"s,omitempty"`
	Body   any    `json:"b,omitempty"`
	Error  string `json:"e,omitempty"`
}

func NewResponse() Response {
	return Response{
		Status: 204,
	}
}

// [T any]
// func (r *Response[T]) GetBody() (T, error) {
// 	var b T

// 	// json.Unmarshal(r.body, &b)

// 	return b, nil
// }
