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

// Package network
package network

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"time"

	"github.com/block-api/block-node/block/function"
	"github.com/block-api/block-node/event"

	"github.com/block-api/block-node/log"
	"github.com/block-api/block-node/network/delivery"
	"github.com/block-api/block-node/network/packet"
	"github.com/block-api/block-node/network/router"
	"github.com/block-api/block-node/network/transport"
	"github.com/block-api/block-node/params"
)

type TCPTransport struct {
	ITransport
	networkManager *Manager
	nodeID         string
	cReceive       chan Packet
}

func (tr TCPTransport) Start(cReceive chan<- Packet) error {
	log.Debug("-- tcp_transport::start --")
	if tr.networkManager.config.Transport != transport.TCP {
		return errors.New("invalid transport type")
	}

	tcpSettings, err := params.GetNetworkSettings[params.NetworkTCPSettings](tr.networkManager.config)
	if err != nil {
		return err
	}

	listener, err := startTCPListener(&tcpSettings)
	if err != nil {
		return err
	}

	defer func(listener net.Listener) {
		err := listener.Close()
		if err != nil {
			log.Debug(err.Error())
		}
	}(listener)

	log.Default("network::tcp_transport::listening -> " + tcpSettings.BindHost + ":" + tcpSettings.BindPort)
	for {
		conn, _ := listener.Accept()

		go listenerConnWorker(tr.networkManager, tr.nodeID, conn, manager.config, tr.cReceive)
	}
}

func startTCPListener(tcpSettings *params.NetworkTCPSettings) (net.Listener, error) {
	listenAt := tcpSettings.BindHost + ":" + tcpSettings.BindPort
	listener, err := net.Listen("tcp", listenAt)

	if err != nil {
		return nil, err
	}

	return listener, err
}

func listenerConnWorker(networkManager *Manager, nodeID string, conn net.Conn, config *params.NetworkConfig, cReceive chan<- Packet) {
	log.Debug("network::tcp_transport::connection_worker::start")

	defer func(conn net.Conn) {
		err := conn.Close()
		if err != nil {
			// log.Warning(err.Error())
		}
		log.Debug("network::tcp_transport::connection_worker::stop")
	}(conn)

	var mb = 1000000

	tmp := make([]byte, mb*1)
	data := make([]byte, 0)
	length := 0

	for {
		n, err := conn.Read(tmp)
		if err != nil {
			if err != io.EOF {
				fmt.Printf("Read error - %s\n", err)
			}
			break
		}
		data = append(data, tmp[:n]...)
		length += n
	}

	packetHeaderIndex := bytes.Index(data, Header)
	if packetHeaderIndex < 0 {
		_, _ = conn.Write(EncodePacket(ErrInvalidHeader))
		return
	}
	conn.Write(make([]byte, 0))
	_ = conn.Close()

	data = data[len(Header):]
	decodedPacket := DecodePacket(data)

	// fmt.Println(data)
	err := decodedPacket.Validate()
	if err != nil {
		log.Debug(err.Error())
		return
	}

	log.Debug("network::tcp_transport::connection_worker::received_packet")

	// log.Default("IS RESPONSE HASH")
	// fmt.Println(decodedPacket.Type)
	// fmt.Println(decodedPacket.TargetID)
	// fmt.Println(decodedPacket.ResponseHash)
	// log.Default("====")

	if decodedPacket.ResponseHash != nil {
		// log.Warning("==== PROCESS RESPONSE HASH ===")
		// fmt.Println(decodedPacket)
		return
	}

	if decodedPacket.Type == packet.Heartbeat {
		var heartbeatBody packet.HeartbeatBody

		bodyBytes, _ := json.Marshal(decodedPacket.Body)
		json.Unmarshal(bodyBytes, &heartbeatBody)

		networkManager.router.UpdateLastSeen(decodedPacket.FromID, time.Now().UnixMilli())
		// add known nodes
		for bodyNodeID, bodyNode := range heartbeatBody.KnownNodes {
			if bodyNodeID == networkManager.nodeID {
				continue
			}

			_ = networkManager.router.Add(bodyNodeID, &router.Node{
				Transport:  bodyNode.Transport,
				NodeID:     bodyNode.NodeID,
				PublicHost: bodyNode.PublicHost,
				PublicPort: bodyNode.PublicPort,
				Functions:  bodyNode.Functions,
			})
			// log.Default(err.Error())
		}

		return
	}

	if decodedPacket.Type == packet.Function {
		fn, err := networkManager.functionManager.Get(decodedPacket.TargetNodeFunction)
		if err != nil {
			log.Debug(err.Error())
			return
		}

		reqFn := function.Request{}
		resFn := function.Response{}

		// TODO: dd
		fnResponse, fnErr := fn(&reqFn, &resFn)
		if fnErr != nil {
			log.Debug(fnErr.Error())
		}

		// send it back
		// log.Default("SEND IT BACK >>>>")
		// fmt.Println(fnResponse.Body)

		resPacket := NewPacket(delivery.All, packet.Function, networkManager.nodeID, decodedPacket.FromID, "", fnResponse, &decodedPacket.Hash)
		networkManager.Send(resPacket)
		// log.Default("<<<< SENT")
		// eventManager := event.GetManager()
		// eventManager.Emit(event.Event{
		// 	Name:    decodedPacket.Hash,
		// 	Payload: fnResponse,
		// })
	}
}

func (tr TCPTransport) Stop() error {
	log.Debug("network::transport:tcp::stop")

	return nil
}

func (tr TCPTransport) Send(sendPacket Packet) {
	log.Debug("network::transport::tcp::send")

	targetNodes := tr.networkManager.router.GetTarget(sendPacket.Delivery, sendPacket.TargetNodeFunction, sendPacket.TargetID)

	for _, targetNode := range targetNodes {

		servAddr := targetNode.PublicHost + ":" + targetNode.PublicPort
		tcpAddr, err := net.ResolveTCPAddr("tcp", servAddr)
		if err != nil {
			log.Debug(err.Error())
			// println("ResolveTCPAddr failed:", err.Error())
			continue
		}
		// log.Default(servAddr)
		conn, err := net.DialTCP("tcp", nil, tcpAddr)
		if err != nil {
			log.Debug(err.Error())
			continue
		}
		// fmt.Println(sendPacket)
		packetBytes, _ := sendPacket.Bytes()
		_, err = conn.Write(packetBytes)
		if err != nil {
			_ = conn.Close()
			log.Debug(err.Error())
			continue
		}

		_ = conn.Close()
		log.Debug("packet sent to node: " + string(sendPacket.Delivery))
		// eventManager := event.GetManager()
		// eventManager.On(string(sendPacket.Hash), waitForResponse)

		// eventManager.Emit(event.Event{
		// 	Name:    decodedPacket.Hash,
		// 	Payload: fnResponse,
		// })
	}
}

func waitForResponse(e event.Event) {
	// var timer = time.NewTimer(time.Second)

	// config.Get
}
