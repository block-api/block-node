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
package log

import (
	"os"
	"time"

	"github.com/fatih/color"
)

type Log struct {
	Message   string
	Type      Type
	CreatedAt time.Time
}

type Type uint

const (
	TypeDefault Type = iota
	TypeWarning
	TypePanic
	TypeFatal
	TypeDebug
)

// Default is an alias for log.New("your text", TypeDefault)
func Default(data string) {
	New(data, TypeDefault)
}

// Debug is an alias for log.New("your text", TypeDebug)
// this message will be seen if environment variable DEBUG=true
func Debug(data string) {
	New(data, TypeDebug)
}

// Warning is an alias for log.New("your text", TypeWarning)
func Warning(data string) {
	New(data, TypeWarning)
}

// Fatal is an alias for log.New("your text", TypeFatal)
func Fatal(data string) {
	New(data, TypeFatal)
}

// Panic is an alias for log.New("your text", TypePanic)
func Panic(data string) {
	New(data, TypePanic)
}

func New(message string, logType Type) {
	log := &Log{
		Message:   message,
		CreatedAt: time.Now(),
	}

	ld := color.New(color.FgGreen)

	switch logType {
	case TypePanic:
		log.Type = TypePanic

		l := color.New(color.FgRed)
		_, _ = l.Println(log.String())
		panic(1)
	case TypeFatal:
		log.Type = TypeFatal

		l := color.New(color.FgRed)
		_, _ = l.Println(log.String())
		os.Exit(1)
	case TypeWarning:
		log.Type = TypeWarning

		c := color.New(color.FgYellow)
		_, _ = c.Println(log.String())
	case TypeDebug:
		log.Type = TypeDebug

		if os.Getenv("DEBUG") == "true" {
			c := color.New(color.FgBlue)
			_, _ = c.Println(log.String())
		}
	default:
		log.Type = TypeDefault
		_, _ = ld.Println(log.String())
	}
}

func (l *Log) String() string {
	output := "[" + l.CreatedAt.Format(time.RFC822) + "] > " + l.Message

	return output
}
