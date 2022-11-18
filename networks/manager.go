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
	"fmt"
	"sync"

	"github.com/containerd/containerd/log"
	"github.com/containernetworking/cni/libcni"
)

type manager struct {
	name     string
	binDir   string
	confDir  string
	cni      libcni.CNI
	networks map[string]Network
	store    Store
	lock     sync.RWMutex
}

var _ Manager = (*manager)(nil)

func NewManager(name string, store Store, binDir string, confDir string) (Manager, error) {
	/*
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
	*/
	m := &manager{
		name:     name,
		store:    store,
		binDir:   binDir,
		confDir:  confDir,
		networks: make(map[string]Network),
	}
	if err := m.Init(context.Background()); err != nil {
		return nil, err
	}
	return m, nil
}

func (m *manager) Init(ctx context.Context) error {
	log.G(ctx).WithField("network manager", m.name).WithField("confdir", m.confDir).Debugf("init")
	return nil
}

func (m *manager) Create(ctx context.Context, name string, options ...NetworkOption) (Network, error) {
	log.G(ctx).WithField("network manager", m.name).WithField("network", name).Debugf("create")
	/*
		var conflist *libcni.NetworkConfigList

		if cfg == nil || cfg.NetworkManager != m.name {
			return errdefs.ErrInvalidArgument
		}
		if cfg.Type == containerd.NetworkConfigConf {
			if conf, err := libcni.ConfFromBytes(cfg.Bytes); err == nil {
				if conf.Network.Name != name {
					return errdefs.ErrInvalidArgument
				}
				if conflist, err = libcni.ConfListFromConf(conf); err != nil {
					return errdefs.ErrUnknown
				}
			} else {
				return errdefs.ErrInvalidArgument
			}
		} else if cfg.Type == containerd.NetworkConfigConfList {
			//if conflist, err := libcni.ConfListFromBytes(cfg.Bytes); err != nil {
			//	return errdefs.ErrInvalidArgument
			//}
		} else {
			return errdefs.ErrInvalidArgument
		}
	*/

	m.lock.Lock()
	defer m.lock.Unlock()
	n := &network{
		NetworkInfo: NetworkInfo{
			Manager: m.name,
			Name:    name,
		},
		cni:   m.cni,
		store: m.store,
	}
	m.networks[name] = n
	// TODO write network to dir
	return n, nil
}

func (m *manager) Delete(ctx context.Context, name string) error {
	log.G(ctx).WithField("network manager", m.name).WithField("network", name).Debugf("delete")
	m.lock.Lock()
	defer m.lock.Unlock()
	if _, ok := m.networks[name]; !ok {
		return fmt.Errorf("")
	}
	// TODO delete file
	delete(m.networks, name)
	return nil
}

func (m *manager) Network(ctx context.Context, name string) (Network, error) {
	log.G(ctx).WithField("network manager", m.name).WithField("network", name).Debugf("get")
	m.lock.RLock()
	defer m.lock.RUnlock()
	n, ok := m.networks[name]
	if !ok {
		return nil, fmt.Errorf("")
	}
	return n, nil
}

func (m *manager) Walk(ctx context.Context, fn func(context.Context, Network) error) error {
	log.G(ctx).WithField("network manager", m.name).Debugf("list")
	m.lock.RLock()
	defer m.lock.RUnlock()
	for _, n := range m.networks {
		fn(ctx, n)
	}
	return nil
}
