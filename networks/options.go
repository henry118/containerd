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

import (
	"github.com/containernetworking/cni/libcni"
)

type AttachmentArgs struct {
	Container      string
	NetNS          string
	IFName         string
	CapabilityArgs map[string]interface{}
	PluginArgs     map[string]string
}

type AttachmentFilter struct {
	Container string
	IFname    string
}

type NetworkOption func(args *NetworkInfo) error
type AttachmentOption func(args *AttachmentArgs) error
type AttachmentFilterOption func(filter *AttachmentFilter) error

func WithCapabilityArg(name string, value interface{}) AttachmentOption {
	return func(args *AttachmentArgs) error {
		args.CapabilityArgs[name] = value
		return nil
	}
}

func WithPluginArg(name string, value string) AttachmentOption {
	return func(args *AttachmentArgs) error {
		args.PluginArgs[name] = value
		return nil
	}
}

func (a *AttachmentArgs) Validate() error {
	return nil
}

func (a *AttachmentArgs) Config() *libcni.RuntimeConf {
	c := &libcni.RuntimeConf{}
	for k, v := range a.PluginArgs {
		c.Args = append(c.Args, [2]string{k, v})
	}
	c.CapabilityArgs = a.CapabilityArgs
	return c
}

func WithNetworkConfig(config NetworkConfig) NetworkOption {
	return func(args *NetworkInfo) error {
		args.Config = config
		return nil
	}
}
