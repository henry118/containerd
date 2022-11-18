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
	"time"

	"github.com/containerd/containerd/errdefs"
	"github.com/containerd/containerd/metadata/boltutil"
	"github.com/containerd/containerd/namespaces"
	"github.com/containerd/typeurl"
	bolt "go.etcd.io/bbolt"
)

const (
	schemaVersion = "v0"
	dbVersion     = 0
)

var (
	bucketKeyVersion       = []byte(schemaVersion)
	bucketKeyObjectEnvs    = []byte("envs")
	bucketKeyObjectLabels  = []byte("labels")
	bucketKeyObjectOptions = []byte("options")
	bucketKeyObjectLeases  = []byte("leases")

	bucketKeyCreatedAt = []byte("createdat")
	bucketKeyUpdatedAt = []byte("updatedat")
)

type Resource struct {
	// ID uniquely identies the network resource in a network env
	ID string
	// Options contains the options used to set up the network resource
	Options map[string]typeurl.Any
	// Labels provides metdata extensions for a network resource
	Labels map[string]string
	// CreatedAt is the time at which the network resource was created
	CreatedAt time.Time
	// UpdatedAt is the time at which the network resource was updated
	UpdatedAt time.Time
}

type Store interface {
	// Get a network resource using the id for the specified env
	Get(ctx context.Context, env string, id string) (*Resource, error)
	// List returns resources for the specified env
	List(ctx context.Context, env string) ([]*Resource, error)
	// Create a network resource in the store
	Create(ctx context.Context, env string, resource *Resource) error
	// Update a network resource
	Update(ctx context.Context, env string, resource *Resource) error
	// Delete a network resource identified by id
	Delete(ctx context.Context, env string, id string) error
}

func NewStore(db *bolt.DB) (Store, error) {
	s := &store{
		db: db,
	}
	if err := s.Init(); err != nil {
		return nil, err
	}
	return s, nil
}

type store struct {
	db *bolt.DB
}

func (s *store) Init() error {
	return s.db.Update(func(tx *bolt.Tx) error {
		if _, err := tx.CreateBucketIfNotExists(bucketKeyVersion); err != nil {
			return err
		}
		return nil
	})
}

func (s *store) Get(ctx context.Context, env string, id string) (*Resource, error) {
	namespace, err := namespaces.NamespaceRequired(ctx)
	if err != nil {
		return nil, err
	}
	r := &Resource{ID: id}
	if err := s.db.View(func(tx *bolt.Tx) error {
		bkt := getEnvBucket(tx, namespace, env)
		if bkt == nil {
			return fmt.Errorf("", env, errdefs.ErrNotFound)
		}
		return readResource(r, bkt)
	}); err != nil {
		return nil, err
	}
	return r, nil
}

func (s *store) List(ctx context.Context, env string) ([]*Resource, error) {
	namespace, err := namespaces.NamespaceRequired(ctx)
	if err != nil {
		return nil, err
	}
	rs := []*Resource{}
	if err := s.db.View(func(tx *bolt.Tx) error {
		bkt := getEnvBucket(tx, namespace, env)
		if bkt == nil {
			return nil // empty
		}
		return bkt.ForEach(func(k, v []byte) error {
			rbkt := bkt.Bucket(k)
			if rbkt == nil {
				return nil
			}
			r := &Resource{ID: string(k)}
			if err := readResource(r, rbkt); err != nil {
				return fmt.Errorf("", err)
			}
			rs = append(rs, r)
			return nil
		})
	}); err != nil {
		return nil, err
	}
	return rs, nil
}

func (s *store) Create(ctx context.Context, env string, resource *Resource) error {
	namespace, err := namespaces.NamespaceRequired(ctx)
	if err != nil {
		return err
	}
	return s.db.Update(func(tx *bolt.Tx) error {
		nbkt, err := tx.CreateBucketIfNotExists([]byte(namespace))
		if err != nil {
			return err
		}
		bkt, err := nbkt.CreateBucketIfNotExists(bucketKeyObjectEnvs)
		if err != nil {
			return err
		}
		ebkt, err := bkt.CreateBucketIfNotExists([]byte(env))
		if err != nil {
			return err
		}
		rbkt, err := ebkt.CreateBucket([]byte(resource.ID))
		if err != nil {
			return err
		}
		resource.CreatedAt = time.Now().UTC()
		resource.UpdatedAt = resource.CreatedAt
		if err := writeResource(resource, rbkt); err != nil {
			return err
		}
		return nil
	})
}

func (s *store) Update(ctx context.Context, env string, resource *Resource) error {
	return fmt.Errorf("", errdefs.ErrNotImplemented)
	/*
		return s.db.Update(tx *bolt.Tx) error {
			bkt := getEnvBucket(tx, env)
			if bkt == nil {
				return fmt.Errorf("", errdefs.ErrNotFound)
			}

		})
	*/
}

func (s *store) Delete(ctx context.Context, env string, id string) error {
	namespace, err := namespaces.NamespaceRequired(ctx)
	if err != nil {
		return err
	}
	return s.db.Update(func(tx *bolt.Tx) error {
		bkt := getEnvBucket(tx, namespace, env)
		if bkt == nil {
			return fmt.Errorf("", errdefs.ErrNotFound)
		}
		if err := bkt.DeleteBucket([]byte(id)); err != nil {
			return err
		}
		return nil
	})
}

func getEnvBucket(tx *bolt.Tx, namespace string, id string) *bolt.Bucket {
	nbkt := tx.Bucket([]byte(namespace))
	if nbkt == nil {
		return nil
	}
	ebkt := nbkt.Bucket(bucketKeyObjectEnvs)
	if ebkt == nil {
		return nil
	}
	return ebkt.Bucket([]byte(id))
}

func readResource(r *Resource, bkt *bolt.Bucket) error {
	if options, err := readOptions(bkt); err != nil {
		return err
	} else {
		r.Options = options
	}
	if labels, err := boltutil.ReadLabels(bkt); err != nil {
		return err
	} else {
		r.Labels = labels
	}
	return boltutil.ReadTimestamps(bkt, &r.CreatedAt, &r.UpdatedAt)
}

func writeResource(r *Resource, bkt *bolt.Bucket) error {
	if err := writeOptions(r.Options, bkt); err != nil {
		return err
	}
	if err := boltutil.WriteLabels(bkt, r.Labels); err != nil {
		return err
	}
	if err := boltutil.WriteTimestamps(bkt, r.CreatedAt, r.UpdatedAt); err != nil {
		return err
	}
	return nil
}

func readOptions(bkt *bolt.Bucket) (map[string]typeurl.Any, error) {
	obkt := bkt.Bucket(bucketKeyObjectOptions)
	if obkt == nil {
		return nil, errdefs.ErrNotFound
	}
	opts := make(map[string]typeurl.Any)
	if err := obkt.ForEach(func(k, _ []byte) error {
		t, err := boltutil.ReadAny(obkt, k)
		if err != nil {
			return err
		}
		opts[string(k)] = t
		return nil
	}); err != nil {
		return nil, fmt.Errorf("")
	}
	return opts, nil
}

func writeOptions(options map[string]typeurl.Any, bkt *bolt.Bucket) error {
	for k, v := range options {
		if err := boltutil.WriteAny(bkt, []byte(k), v); err != nil {
			return err
		}
	}
	return nil
}
