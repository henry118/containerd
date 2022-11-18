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

	"github.com/containerd/containerd/gc"
	"github.com/containerd/containerd/log"
	"github.com/containerd/containerd/metadata"
)

const (
	ResourceNetwork gc.ResourceType = 0x11
)

type ResourceCollector struct {
	managers map[string]Manager
}

func NewResourceCollector(managers map[string]Manager) (*ResourceCollector, error) {
	c := &ResourceCollector{
		managers: make(map[string]Manager),
	}
	for k, m := range managers {
		c.managers[k] = m
	}
	return c, nil
}

func (rc *ResourceCollector) StartCollection(ctx context.Context) (metadata.CollectionContext, error) {
	c := &gccontext{
		managers: rc.managers,
		leases:   make(map[string][]string),
	}
	for _, m := range rc.managers {
		if err := m.Walk(ctx, func(ctx context.Context, n Network) error {

			return nil
		}); err != nil {
			// log
		}
		/*
			if err != nil {
				// log
				continue
			}
			for _, n := range nl {
				al, err := n.ListAttachments(ctx)
				if err != nil {
					// log
					continue
				}
				for _, a := range al {
					c.all = append(c.all, a.ID(ctx))
					lease := a.Info(ctx).Lease
					if len(lease) > 0 {
						c.leases[lease] = append(c.leases[lease], a.ID(ctx))
					}
				}
			}
		*/
	}
	return c, nil
}

func (rc *ResourceCollector) ReferenceLabel() string {
	return "network."
}

type gccontext struct {
	managers map[string]Manager
	all      []string
	leases   map[string][]string
	removes  []string
}

func (c *gccontext) All(fn func(gc.Node)) {
	log.L.Debugf("All")
	for _, id := range c.all {
		// get ns from id
		fn(gcnode(ResourceNetwork, "", string(id)))
	}
}

func (c *gccontext) Active(string, func(gc.Node)) {
	// noop, use the default gc that scans labels of existing containers
}

func (c *gccontext) Leased(namespace, lease string, fn func(gc.Node)) {
	log.L.Debugf("Leased")
	ls, ok := c.leases[lease]
	if !ok {
		return
	}
	for _, l := range ls {
		// split l into namespace,
		gcnode(ResourceNetwork, "", l)
	}
}

func (c *gccontext) Remove(n gc.Node) {
	log.L.Debugf("Remove")
	if n.Type == ResourceNetwork {
		c.removes = append(c.removes, n.Key)
	}
}

func (c *gccontext) Cancel() error {
	log.L.Debugf("Cancel")
	c.removes = nil
	return nil
}

func (c *gccontext) Finish() error {
	log.L.Debugf("Finish")
	for _, id := range c.removes {
		// split id into network
		net := "default"
		if m, ok := c.managers[net]; ok {
			if err := m.Delete(context.Background(), id); err != nil {
				// log
			}
		} else {
			// log
		}
	}
	c.removes = nil
	return nil
}

func gcnode(t gc.ResourceType, ns, key string) gc.Node {
	return gc.Node{
		Type:      t,
		Namespace: ns,
		Key:       key,
	}
}
