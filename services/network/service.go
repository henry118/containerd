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
	"github.com/containerd/containerd/errdefs"
	"github.com/containerd/containerd/network"
	"github.com/containerd/containerd/plugin"
	"google.golang.org/grpc"
)

func init() {
	plugin.Register(&plugin.Registration{
		Type: plugin.GRPCPlugin,
		ID:   "networks",
		Requires: []plugin.Type{
			plugin.ServicePlugin,
		},
		InitFn: func(ic *plugin.InitContext) (interface{}, error) {
			return &service{}, nil
		},
	})
}

type service struct {
	envs map[string]network.Env
	// api
}

func (s *service) Register(gs *grpc.Server) error {
	return nil
}

func (s *service) getEnv(name string) (network.Env, error) {
	if name == "" {
		return nil, errdefs.ToGRPCf(errdefs.ErrInvalidArgument, "env name argument missing")
	}

	e := s.envs[name]
	if e == nil {
		return nil, errdefs.ToGRPCf(errdefs.ErrInvalidArgument, "network envs not configured: %s", name)
	}

	return e, nil
}

/*
func (s *service) Setup(ctx context.Context, id string, path string, opts ...Opts) (*Result, error) {
	return nil, nil
}

func (s *service) Remove(ctx context.Context, id string) error {
	return nil
}

func (s *service) Load(ctx context.Context) error {
	return nil
}

func (s *service) Check(ctx context.Context, id string) error {
	return nil
}

func (s *service) Status(ctx context.Context) error {
	return nil
}

func (s *service) GetConfig(ctx context.Context) (*ConfigResult, error) {
	return nil, nil
  }
*/
