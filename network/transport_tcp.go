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
	"github.com/block-api/block-node/block/function"
	"io"
	"net"
	"time"

	"github.com/block-api/block-node/log"
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
			log.Warning(err.Error())
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
	_ = conn.Close()

	data = data[len(Header):]
	decodedPacket := DecodePacket(data)

	err := decodedPacket.Validate()
	if err != nil {
		log.Debug(err.Error())
		return
	}

	log.Debug("network::tcp_transport::connection_worker::received_packet")

	if decodedPacket.Type == packet.Heartbeat {
		var heartbeatBody packet.HeartbeatBody

		bodyBytes, _ := json.Marshal(decodedPacket.Body)
		json.Unmarshal(bodyBytes, &heartbeatBody)

		networkManager.router.UpdateLastSeen(decodedPacket.FromID, time.Now().UnixMilli())

		// add known nodes
		for bodyNodeID, bodyNode := range heartbeatBody.KnownNodes {
			if bodyNodeID == networkManager.nodeID {
				// log.Default("self id, skip known node")
				continue
			}
			networkManager.router.Add(bodyNodeID, &router.Node{
				Transport:  bodyNode.Transport,
				NodeID:     bodyNode.NodeID,
				PublicHost: bodyNode.PublicHost,
				PublicPort: bodyNode.PublicPort,
				Functions:  bodyNode.Functions,
			})
		}

		return
	}

	if decodedPacket.Type == packet.Function {
		fn, err := networkManager.functionManager.Get(decodedPacket.TargetNodeFunction)
		if err != nil {
			log.Warning(err.Error())
			return
		}

		reqFn := function.Request{}
		resFn := function.Response{}
		log.Default("XxXxXxXxXxXxXxXxx")
		fnResponse, fnErr := fn(&reqFn, &resFn)
		if fnErr != nil {
			log.Warning(fnErr.Error())
		}

		// send it back
		fmt.Println(fnResponse.Body)
	}
}

func (tr TCPTransport) Stop() error {
	log.Debug("-- tcp_transport::stop --")

	return nil
}

func (tr TCPTransport) Send(sendPacket Packet) {
	log.Debug("-- tcp_transport::send --")

	targetNodes := tr.networkManager.router.GetTarget(sendPacket.Delivery, sendPacket.TargetNodeFunction, sendPacket.TargetID)

	for _, targetNode := range targetNodes {

		servAddr := targetNode.PublicHost + ":" + targetNode.PublicPort
		tcpAddr, err := net.ResolveTCPAddr("tcp", servAddr)
		if err != nil {
			println("ResolveTCPAddr failed:", err.Error())
			continue
		}

		conn, err := net.DialTCP("tcp", nil, tcpAddr)
		if err != nil {
			log.Warning(err.Error())
			continue
		}

		packetBytes, _ := sendPacket.Bytes()
		_, err = conn.Write(packetBytes)
		if err != nil {
			_ = conn.Close()
			log.Warning(err.Error())
			continue
		}

		_ = conn.Close()
		log.Default("packet sent to node: " + string(sendPacket.Delivery))
	}
}
