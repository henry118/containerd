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
	"fmt"

	"github.com/containerd/containerd/gc"
	"github.com/containerd/containerd/log"
	"github.com/containerd/containerd/metadata"
	bolt "go.etcd.io/bbolt"
)

const (
	ResourceNetwork gc.ResourceType = 0x11
)

type ResourceCollector struct {
	db      *bolt.DB
	gcnodes []gc.Node
	envs    map[string]Env
}

func NewCollector(db *bolt.DB, envs map[string]Env) (*ResourceCollector, error) {
	c := &ResourceCollector{
		db:   db,
		envs: envs,
	}
	return c, nil
}

func (rc *ResourceCollector) StartCollection(context.Context) (metadata.CollectionContext, error) {
	return rc, nil
}

func (rc *ResourceCollector) ReferenceLabel() string {
	return "network."
}

func (rc *ResourceCollector) All(func(gc.Node)) {
	log.L.Debugf("All")
}

func (rc *ResourceCollector) Active(string, func(gc.Node)) {
	log.L.Debugf("Active")
	// noop, use the default gc that scans labels of existing containers
}

func (rc *ResourceCollector) Leased(namespace, lease string, fn func(gc.Node)) {
	log.L.Debugf("Leased")
	rc.db.View(func(tx *bolt.Tx) error {
		nbkt := tx.Bucket([]byte(namespace))
		if nbkt == nil {
			return nil
		}
		bkt := nbkt.Bucket(bucketKeyObjectLeases)
		if bkt == nil {
			return nil
		}
		lbkt := bkt.Bucket([]byte(lease))
		if lbkt == nil {
			return nil
		}
		lbkt.ForEach(func(k, _ []byte) error {
			fn(gcnode(ResourceNetwork, namespace, string(k)))
			return nil
		})
		return nil
	})
}

func (rc *ResourceCollector) Remove(n gc.Node) {
	log.L.Debugf("Remove")
	if n.Type == ResourceNetwork {
		rc.gcnodes = append(rc.gcnodes, n)
	}
}

func (rc *ResourceCollector) Cancel() error {
	log.L.Debugf("Cancel")
	rc.gcnodes = nil
	return nil
}

func (rc *ResourceCollector) Finish() error {
	log.L.Debugf("Finish")
	return rc.db.Update(func(tx *bolt.Tx) error {
		for _, n := range rc.gcnodes {
			if n.Type != ResourceNetwork {
				continue
			}
			var eid, rid string
			if _, err := fmt.Sscanf(n.Key, "%s/%s", &eid, &rid); err != nil {
				// log
				continue
			}
			if env, ok := rc.envs[eid]; !ok {
				// log
			} else {
				if err := env.Remove(context.Background(), rid); err != nil {
					// log
					continue
				}
				if bkt := getEnvBucket(tx, n.Namespace, eid); bkt == nil {
					// log
					continue
				} else if err := bkt.DeleteBucket([]byte(rid)); err != nil {
					// log
				}
			}
		}
		return nil
	})
}

func gcnode(t gc.ResourceType, ns, key string) gc.Node {
	return gc.Node{
		Type:      t,
		Namespace: ns,
		Key:       key,
	}
}
