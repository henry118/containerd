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
	"context"
	"os"
	"path/filepath"
	"time"

	api "github.com/containerd/containerd/api/services/networks/v1"
	"github.com/containerd/containerd/log"
	"github.com/containerd/containerd/metadata"
	"github.com/containerd/containerd/networks"
	"github.com/containerd/containerd/plugin"
	"github.com/containerd/containerd/services"
	bolt "go.etcd.io/bbolt"
)

func init() {
	defaultCfg := DefaultConfig()
	plugin.Register(&plugin.Registration{
		Type:   plugin.ServicePlugin,
		ID:     services.NetworkService,
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
			path := filepath.Join(ic.State, "networks.db")
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

			store, err := networks.NewStore(db)
			if err != nil {
				return nil, err
			}

			locals := make(map[string]networks.Service)
			managers := make(map[string]networks.Manager)
			pc := ic.Config.(*PluginConfig)
			for key, cfg := range pc.Managers {
				m, err := networks.NewManager(key, store, cfg.PluginBinDir, cfg.PluginConfDir)
				if err != nil {
					return nil, err
				}
				managers[key] = m
				locals[key] = &local{
					name: key,
					m:    m,
				}
			}

			collector, err := networks.NewResourceCollector(managers)
			if err != nil {
				return nil, err
			}

			m.(*metadata.DB).RegisterCollectibleResource(networks.ResourceNetwork, collector)

			return locals, nil
		},
	})
}

type local struct {
	name string
	m    networks.Manager
}

var _ networks.Service = (*local)(nil)

func (l *local) CreateNetwork(ctx context.Context, name string, options ...networks.NetworkOption) (networks.NetworkInfo, error) {
	n, err := l.m.Create(ctx, name, options...)
	if err != nil {
		return networks.NetworkInfo{}, err
	}
	return n.Info(ctx), nil
}

func (l *local) DeleteNetwork(ctx context.Context, name string) error {
	return l.m.Delete(ctx, name)
}

func (l *local) GetNetwork(ctx context.Context, name string) (networks.NetworkInfo, error) {
	n, err := l.m.Network(ctx, name)
	if err != nil {
		return networks.NetworkInfo{}, err
	}
	return n.Info(ctx), nil
}

func (l *local) ListNetworks(ctx context.Context) ([]networks.NetworkInfo, error) {
	var res []networks.NetworkInfo
	if err := l.m.Walk(ctx, func(ctx context.Context, n networks.Network) error {
		res = append(res, n.Info(ctx))
		return nil
	}); err != nil {
		return nil, err
	}
	return res, nil
}

func (l *local) AttachNetwork(ctx context.Context, network string, options ...networks.AttachmentOption) (networks.AttachmentInfo, error) {
	n, err := l.m.Network(ctx, network)
	if err != nil {
		return networks.AttachmentInfo{}, err
	}
	a, err := n.Attach(ctx, options...)
	if err != nil {
		return networks.AttachmentInfo{}, err
	}
	return a.Info(ctx), nil
}

func (l *local) DetachNetwork(ctx context.Context, network string, attachment string) error {
	n, err := l.m.Network(ctx, network)
	if err != nil {
		return err
	}
	a, err := n.Attachment(ctx, attachment)
	if err != nil {
		return err
	}
	return a.Remove(ctx)
}

func (l *local) GetAttachment(ctx context.Context, network string, attachment string) (networks.AttachmentInfo, error) {
	n, err := l.m.Network(ctx, network)
	if err != nil {
		return networks.AttachmentInfo{}, err
	}
	a, err := n.Attachment(ctx, attachment)
	if err != nil {
		return networks.AttachmentInfo{}, err
	}
	return a.Info(ctx), nil
}

func (l *local) CheckAttachment(ctx context.Context, network string, attachment string) error {
	n, err := l.m.Network(ctx, network)
	if err != nil {
		return err
	}
	a, err := n.Attachment(ctx, attachment)
	if err != nil {
		return err
	}
	return a.Check(ctx)
}

func (l *local) ListAttachments(ctx context.Context, network string, options ...networks.AttachmentFilterOption) ([]networks.AttachmentInfo, error) {
	n, err := l.m.Network(ctx, network)
	if err != nil {
		return nil, err
	}
	var res []networks.AttachmentInfo
	if err := n.Walk(ctx, func(ctx context.Context, a networks.Attachment) error {
		// TODO: process filter options here
		res = append(res, a.Info(ctx))
		return nil
	}); err != nil {
		return nil, err
	}
	return res, nil
}

func toNetworkConfig(p *api.NetworkConfig) networks.NetworkConfig {
	c := networks.NetworkConfig{
		Type:  networks.NetworkConfigInvalid,
		Bytes: p.Data,
	}
	switch p.Type {
	case api.NetworkConfigType_NC_CONF:
		c.Type = networks.NetworkConfigConf
	case api.NetworkConfigType_NC_CONFLIST:
		c.Type = networks.NetworkConfigConfList
	}
	return c
}

func fromNetworkConfig(c networks.NetworkConfig) *api.NetworkConfig {
	p := &api.NetworkConfig{
		Type: api.NetworkConfigType_NC_UNKNOWN,
		Data: c.Bytes,
	}
	switch c.Type {
	case networks.NetworkConfigConf:
		p.Type = api.NetworkConfigType_NC_CONF
	case networks.NetworkConfigConfList:
		p.Type = api.NetworkConfigType_NC_CONFLIST
	}
	return p
}

/*
func toAttachmentCreateOption(a *api.AttachmentArgument) (networks.AttachmentCreateOption, error) {
	v, err := typeurl.UnmarshalAny(a.Value)
	if err != nil {
		return nil, err
	}
	switch a.Type {
	case api.AttachmentArgumentType_ARG_CAPABILITY:
		return networks.WithCapabilityArg(a.Name, v), nil
	case api.AttachmentArgumentType_ARG_PLUGIN:
		return networks.WithPluginArg(a.Name, v.(*networks.StringValue).Val), nil
	}
	return nil, errdefs.ErrInvalidArgument
}
*/

// PluginConfig contains toml config related Network service
type PluginConfig struct {
	DefaultManager string                   `toml:"default_manager" json:"defaultManager"`
	Managers       map[string]ManagerConfig `toml:"managers" json:"managers"`
}

// EnvConfig contains toml config related network env
type ManagerConfig struct {
	// PluginBinDir is the directory in which the binaries for the plugin are kept.
	PluginBinDir string `toml:"bin_dir" json:"binDir"`
	// PluginConfDir is the directory in which the admin places a network(CNI) conf.
	PluginConfDir string `toml:"conf_dir" json:"confDir"`
}

func DefaultConfig() *PluginConfig {
	return &PluginConfig{
		DefaultManager: "default",
		Managers: map[string]ManagerConfig{
			"default": {
				PluginBinDir:  "/opt/cni/bin",
				PluginConfDir: "/etc/cni/net.d",
			},
		},
	}
}
