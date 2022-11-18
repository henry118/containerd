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
	"os"
	"path/filepath"
	"time"

	"github.com/containerd/containerd/log"
	"github.com/containerd/containerd/metadata"
	"github.com/containerd/containerd/network"
	"github.com/containerd/containerd/plugin"
	"github.com/containerd/containerd/services"
	bolt "go.etcd.io/bbolt"
)

func init() {
	defaultCfg := network.DefaultConfig()
	plugin.Register(&plugin.Registration{
		Type:   plugin.ServicePlugin,
		ID:     services.NetworksService,
		Config: defaultCfg,
		Requires: []plugin.Type{
			plugin.MetadataPlugin,
		},
		InitFn: func(ic *plugin.InitContext) (interface{}, error) {
			m, err := ic.Get(plugin.MetadataPlugin)
			if err != nil {
				return nil, err
			}

			if err := os.MkdirAll(ic.State, 0711); err != nil {
				return nil, err
			}
			path := filepath.Join(ic.State, "network.db")
			options := *bolt.DefaultOptions
			// Reading bbolt's freelist sometimes fails when the file has a data corruption.
			// Disabling freelist sync reduces the chance of the breakage.
			// https://github.com/etcd-io/bbolt/pull/1
			// https://github.com/etcd-io/bbolt/pull/6
			options.NoFreelistSync = true
			// Without the timeout, bbolt.Open would block indefinitely due to flock(2).
			//options.Timeout = timeout.Get(0)

			doneCh := make(chan struct{})
			go func() {
				t := time.NewTimer(10 * time.Second)
				defer t.Stop()
				select {
				case <-t.C:
					log.G(ic.Context).WithField("plugin", "bolt").Warn("waiting for response from boltdb open")
				case <-doneCh:
					return
				}
			}()

			db, err := bolt.Open(path, 0644, &options)
			close(doneCh)
			if err != nil {
				return nil, err
			}

			store, err := network.NewStore(db)
			if err != nil {
				return nil, err
			}

			var (
				envs   = make(map[string]network.Env)
				locals []network.Env
			)

			pcfg := ic.Config.(*network.PluginConfig)
			for key, cfg := range pcfg.Envs {
				l, err := network.NewEnv(&cfg)
				if err != nil {
					return nil, err
				}
				envs[key] = &env{
					name:  key,
					l:     l,
					store: store,
				}
				locals = append(locals, l)
			}

			mdb := m.(*metadata.DB)
			collector, err := network.NewCollector(db, nil)
			if err != nil {
				return nil, err
			}
			mdb.RegisterCollectibleResource(network.ResourceNetwork, collector)

			return envs, nil
		},
	})
}

type env struct {
	name  string
	l     network.Env
	store network.Store
}

func (e *env) Setup(ctx context.Context, id string, path string, opts ...network.Opts) (*network.Result, error) {
	return e.l.Setup(ctx, id, path, opts...)
}

func (e *env) Remove(ctx context.Context, id string) error {
	return e.l.Remove(ctx, id)
}

func (e *env) Load(ctx context.Context) error {
	return e.l.Load(ctx)
}

func (e *env) Check(ctx context.Context, id string) error {
	return e.l.Check(ctx, id)
}

func (e *env) Status(ctx context.Context) error {
	return e.l.Status(ctx)
}

func (e *env) GetConfig(ctx context.Context) (*network.ConfigResult, error) {
	return e.l.GetConfig(ctx)
}
