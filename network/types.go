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

package network

import (
	"net"

	cni "github.com/containernetworking/cni/pkg/types"
)

type IFConfig struct {
	IPConfigs []*IPConfig
	Mac       string
	Sandbox   string
}

type IPConfig struct {
	IP      net.IP
	Gateway net.IP
}

type Route struct {
	Dst net.IPNet
	GW  net.IP
}

type Namespace struct {
	Id             string
	Path           string
	CapabilityArgs map[string]interface{}
	Args           map[string]string
}

// typeurl registered
type DNS struct {
	Servers  []string
	Searches []string
	Options  []string
}

// typeurl registered
type PortMapping struct {
	HostPort      int32
	ContainerPort int32
	Protocol      string
	HostIP        string
}

// typeurl registered
type IPRanges struct {
	Subnet     string
	RangeStart string
	RangeEnd   string
	Gateway    string
}

// typeurl registered
type BandWidth struct {
	IngressRate  uint64
	IngressBurst uint64
	EgressRate   uint64
	EgressBurst  uint64
}

type Opts func(ns *Namespace) error

type Result struct {
	Interfaces map[string]*IFConfig
	DNS        []DNS
	Routes     []Route
}

// typeurl registered any type
type ConfigResult struct {
	PluginDirs       []string
	PluginConfDir    string
	PluginMaxConfNum int
	Prefix           string
	Networks         []*ConfNetwork
}

type ConfNetwork struct {
	Config *NetworkConfList
	IFName string
}

type NetworkConfList struct {
	Name       string
	CNIVersion string
	Plugins    []*NetworkConf
	Source     string
}

type NetworkConf struct {
	Network *cni.NetConf
	Source  string
}
