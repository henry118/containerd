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

/* Schema
v0
|-- <network-manager-id>
    |-- <network-id>
        |-- <attachment-id>
            |-- capargs
            |-- *key*: <bytes>
        |-- args
            |--- *key*: <string>
        |-- labels
            |-- lease: <string>
            |-- nspath: <string>
        |-- status: <int>
        |-- result: <bytes>
        |-- createdat: <datetime>
        |-- updatedat: <datetime>
*/

import (
	"context"
	"fmt"
	"time"

	"github.com/containernetworking/cni/pkg/types"
	bolt "go.etcd.io/bbolt"
)

const (
	schemaVersion = "v0"
)

var (
	bucketKeyVersion       = []byte(schemaVersion)
	bucketKeyObjectLabels  = []byte("labels")
	bucketKeyObjectCapArgs = []byte("capargs")
	bucketKeyObjectArgs    = []byte("args")

	bucketKeyStatus    = []byte("status")
	bucketKeyResult    = []byte("result")
	bucketKeyCreatedAt = []byte("createdat")
	bucketKeyUpdatedAt = []byte("updatedat")
)

type AttachmentRecord struct {
	ID        string
	Args      *AttachmentArgs
	Labels    map[string]string
	Status    int
	Result    types.Result
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Store interface {
	// Create a network resource in the store
	Create(ctx context.Context, manager, network string, r *AttachmentRecord) error
	// Update
	Update(ctx context.Context, manager, network string, r *AttachmentRecord) error
	// Get a network resource using the id for the specified env
	Get(ctx context.Context, manager, network, id string) (*AttachmentRecord, error)
	// Delete a network resource identified by id
	Delete(ctx context.Context, manager, network, id string) error
	// List returns resources for the specified env
	Walk(ctx context.Context, manager, network string, fn func(*AttachmentRecord) error) error
}

type DB interface {
	View(func(tx *bolt.Tx) error) error
	Update(func(tx *bolt.Tx) error) error
}

func NewStore(db DB) (Store, error) {
	s := &store{
		db: db,
	}
	if err := s.Init(); err != nil {
		return nil, err
	}
	return s, nil
}

type store struct {
	db DB
	t  map[string]*AttachmentRecord
}

var _ Store = (*store)(nil)

type transactionKey struct{}

func (s *store) Init() error {
	s.t = make(map[string]*AttachmentRecord)
	// TODO: also create network managers buckets
	return s.update(context.Background(), func(tx *bolt.Tx) error {
		if _, err := tx.CreateBucketIfNotExists(bucketKeyVersion); err != nil {
			return err
		}
		return nil
	})
}

func (s *store) Create(ctx context.Context, manager, network string, r *AttachmentRecord) error {
	return s.update(ctx, func(tx *bolt.Tx) error {
		s.t[r.ID] = r
		// find manager/network bucket
		// check record not exists
		// create new record
		// set status to PENDING
		return nil
	})
}

func (s *store) Update(ctx context.Context, manager, network string, r *AttachmentRecord) error {
	return s.update(ctx, func(tx *bolt.Tx) error {
		// find manager/network bucket
		// check record exists
		// only update the following:
		// - status
		// - result
		return nil
	})
}

func (s *store) Get(ctx context.Context, manager, network, id string) (*AttachmentRecord, error) {
	var rec *AttachmentRecord
	if err := s.view(ctx, func(tx *bolt.Tx) error {
		r, ok := s.t[id]
		if !ok {
			return fmt.Errorf("not found")
		}
		rec = r
		// find manager/network bucket
		// read record
		return nil
	}); err != nil {
		return nil, err
	}
	return rec, nil
}

func (s *store) Walk(ctx context.Context, manager, network string, fn func(*AttachmentRecord) error) error {
	return s.view(ctx, func(tx *bolt.Tx) error {
		for _, v := range s.t {
			fn(v)
		}
		// get manager/network bucket
		// iterator records and call fn
		return nil
	})
}

func (s *store) Delete(ctx context.Context, manager, network, id string) error {
	return s.update(ctx, func(tx *bolt.Tx) error {
		if _, ok := s.t[id]; !ok {
			return fmt.Errorf("not found")
		}
		delete(s.t, id)
		// find manager/network bucket
		// delete record
		return nil
	})
}

func (s *store) view(ctx context.Context, fn func(*bolt.Tx) error) error {
	tx, ok := ctx.Value(transactionKey{}).(*bolt.Tx)
	if !ok {
		return s.db.View(fn)
	}
	return fn(tx)
}

func (s *store) update(ctx context.Context, fn func(*bolt.Tx) error) error {
	tx, ok := ctx.Value(transactionKey{}).(*bolt.Tx)
	if !ok {
		return s.db.Update(fn)
	} else if !tx.Writable() {
		return fmt.Errorf("unable to use transaction from context: %w", bolt.ErrTxNotWritable)
	}
	return fn(tx)
}

/*
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

func readCapArgs(bkt *bolt.Bucket) (map[string]typeurl.Any, error) {
	obkt := bkt.Bucket(bucketKeyObjectCapArgs)
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

func writeCapArgs(options map[string]typeurl.Any, bkt *bolt.Bucket) error {
	for k, v := range options {
		if err := boltutil.WriteAny(bkt, []byte(k), v); err != nil {
			return err
		}
	}
	return nil
}

func getBucket(tx *bolt.Tx, names ...[]byte) *bolt.Bucket {
	bkt := tx.Bucket(bucketKeyVersion)
	if bkt == nil {
		return nil
	}
	for _, name := range names {
		bkt = bkt.Bucket(name)
		if bkt == nil {
			return nil
		}
	}
	return bkt
}

func getNetworkBucket(tx *bolt.Tx, manager string, id string) *bolt.Bucket {
	return getBucket(tx, bucketKeyObjectManagers, []byte(manager), []byte(id))
}

func getManagerBucket(tx *bolt.Tx, manager string) *bolt.Bucket {
	return getBucket(tx, bucketKeyObjectManagers, []byte(manager))
}

func getAttachmentBucket(tx *bolt.Tx, manager string, network string, namespace string, id string) *bolt.Bucket {
	return getBucket(tx, bucketKeyObjectManagers, []byte(manager), []byte(network), []byte(namespace), []byte(id))
}

func readNetworkCfg(c *NetworkCfg, bkt *bolt.Bucket) error {
	return nil
}

func writeNetworkCfg(c *NetworkCfg, bkt *bolt.Bucket) error {
	return nil
}

func readAttachment(a *Attachment, bkt *bolt.Bucket) error {
	return nil
}

func writeAttachment(a *Attachment, bkt *bolt.Bucket) error {
	return nil
}
*/
