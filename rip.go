/*
 * Copyright 2013 Murali Suriar
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License in the LICENSE file, or at:
 *
 *     https://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package garbler

import (
	"bytes"
	"encoding/binary"
	"log"
	"net"
	"syscall"
)

const (
	rip_cmd_response uint8 = 2
	rip_metric_healthy uint32 = 1
	rip_metric_unhealthy uint32 = 16
	rip_version uint8 = 2
)

type ripMsg struct {
	command uint8
	version uint8
	pad uint16
	afi uint16
	tag uint16
	network uint32
	mask uint32
	nexthop uint32
	metric uint32
}

func packPrefix(prefix string) (network, subnet_mask uint32) {
	ip, ipnet, err := net.ParseCIDR(prefix)
	if err != nil {
		log.Fatalf("Error while parsing prefix %s: %s", prefix, err)
	}
	ip = ip.To4()
	mask := net.IP(ipnet.Mask)
	return IPtoInt32(ip), IPtoInt32(mask)
}

func newRipMsg(prefix string, metric uint32) (rm *ripMsg) {

	network, mask := packPrefix(prefix)
	rm = &ripMsg{
		command: rip_cmd_response,
		version: rip_version,
		pad: 0,
		afi: syscall.AF_INET,
		tag: 0,
		network: network,
		mask: mask,
		nexthop: 0,
		metric: metric,
	}

	return rm
}

func newHealthyRipMsg(prefix string) (rm *ripMsg) {
	return newRipMsg(prefix, rip_metric_healthy)
}

func newUnhealthyRipMsg(prefix string) (rm *ripMsg) {
	return newRipMsg(prefix, rip_metric_unhealthy)
}

func sendRipMsg(rm *ripMsg) {
	serverAddr, err := net.ResolveUDPAddr("udp", "224.0.0.9:520")
	if err != nil {
		log.Fatal("Failure!")
	}
	conn, err := net.DialUDP("udp", nil, serverAddr)
	if err != nil {
		log.Fatal("Failure!", err)
	}

	buf := new(bytes.Buffer)

	err = binary.Write(buf, binary.BigEndian, rm)
	if err != nil {
		log.Fatal("Failure!", err)
	}

	_, err = conn.Write(buf.Bytes())
	if err != nil {
		log.Fatal("Failure!", err)
	}

	err = conn.Close()
	if err != nil {
		log.Fatalf("Failed to close socket: %s", err)
	}
}

func IPtoInt32(ip net.IP) (ipint uint32) {
	ipint = 0
	for i := 0; i < 4; i++ {
		shiftbits := 3-i
		sum := uint32(ip[i]) << uint32(8*shiftbits)
		ipint = ipint + sum
	}
	return ipint
}
