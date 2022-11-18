/*
   Copyright The containerd Authors.

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

package networks

import "net"

type NetworkConfigType int

const (
	NetworkConfigInvalid  NetworkConfigType = 0
	NetworkConfigConfList NetworkConfigType = 1
	NetworkConfigConf     NetworkConfigType = 2
)

type NetworkConfig struct {
	Type  NetworkConfigType
	Bytes []byte
}

type Interface struct {
	Name    string
	IPs     []IPConfig
	Mac     string
	Sandbox string
}

type IPConfig struct {
	IP      net.IP
	Gateway net.IP
}

type Route struct {
	Dst net.IPNet
	GW  net.IP
}

type DNS struct {
	Servers  []string
	Searches []string
	Options  []string
}

type PortMapping struct {
	HostPort      int32
	ContainerPort int32
	Protocol      string
	HostIP        string
}

type IPRanges struct {
	Subnet     string
	RangeStart string
	RangeEnd   string
	Gateway    string
}

type BandWidth struct {
	IngressRate  uint64
	IngressBurst uint64
	EgressRate   uint64
	EgressBurst  uint64
}

type StringValue struct {
	Val string
}
