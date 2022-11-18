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
	"context"

	gocni "github.com/containerd/go-cni"
)

type Env interface {
	// Setup sets up a networks for a namespace
	Setup(ctx context.Context, id string, path string, opts ...Opts) (*Result, error)

	// Remove tears down a networks from a namespace
	Remove(ctx context.Context, id string) error

	// Load load/reload networks configurations
	Load(ctx context.Context) error

	// Check if the networks is in desired state
	Check(ctx context.Context, id string) error

	// Status checks the status of the cni initialization
	Status(ctx context.Context) error

	// GetConfig returns a copy of networks configurations
	GetConfig(ctx context.Context) (*ConfigResult, error)
}

type env struct {
	cni gocni.CNI
}

func NewEnv(cfg *EnvConfig) (Env, error) {
	i, err := gocni.New(
		gocni.WithPluginDir([]string{cfg.PluginBinDir}),
		gocni.WithPluginConfDir(cfg.PluginConfDir),
		gocni.WithPluginMaxConfNum(cfg.MaxConfNum),
		gocni.WithMinNetworkCount(cfg.MinConfNum),
		//gocni.WithLoNetwork(cfg.LoNetwork),
		gocni.WithInterfacePrefix(cfg.IFPrefix),
	)
	if err != nil {
		return nil, err
	}
	return &env{
		cni: i,
	}, nil
}

func (e *env) Setup(ctx context.Context, id string, path string, opts ...Opts) (*Result, error) {
	return nil, nil
}

func (e *env) Remove(ctx context.Context, id string) error {
	return nil
}

func (e *env) Load(ctx context.Context) error {
	return nil
}

func (e *env) Check(ctx context.Context, id string) error {
	return nil
}

func (e *env) Status(ctx context.Context) error {
	return nil
}

func (e *env) GetConfig(ctx context.Context) (*ConfigResult, error) {
	return nil, nil
}
