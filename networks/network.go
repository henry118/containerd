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

	"github.com/containerd/containerd/log"
	"github.com/containerd/containerd/namespaces"
	"github.com/containernetworking/cni/libcni"
)

type network struct {
	NetworkInfo
	manager Manager
	cni     libcni.CNI
	store   Store
}

var _ Network = (*network)(nil)

func (n *network) Name(ctx context.Context) string {
	return n.NetworkInfo.Name
}

func (n *network) Info(ctx context.Context) NetworkInfo {
	return n.NetworkInfo
}

func (n *network) Delete(ctx context.Context) error {
	return n.manager.Delete(ctx, n.NetworkInfo.Name)
}

func (n *network) Attach(ctx context.Context, options ...AttachmentOption) (Attachment, error) {
	log.G(ctx).WithField("network manager", n.NetworkInfo.Manager).WithField("network", n.NetworkInfo.Name).Debugf("attach")

	ns, err := namespaces.NamespaceRequired(ctx)
	if err != nil {
		return nil, err
	}

	/*
		// get namespace from context
		// get lease from context

		// apply options
		// generate attachment id

		// lock
		// get network config
		// unlock

		// create attachment in store
		// call libcni to add network list
	*/

	args := AttachmentArgs{}
	for _, o := range options {
		if err := o(&args); err != nil {
			return nil, err
		}
	}
	rec := AttachmentRecord{
		ID:     generateAttachmentID(n.NetworkInfo.Manager, n.NetworkInfo.Name, ns, args.Container, args.IFName),
		Args:   &args,
		Status: 0,
	}
	err = n.store.Create(ctx, n.NetworkInfo.Manager, n.NetworkInfo.Name, &rec)
	if err != nil {
		return nil, err
	}
	return &attachment{
		AttachmentInfo: AttachmentInfo{
			Manager: n.NetworkInfo.Manager,
			Network: n.NetworkInfo.Name,
		},
		cni:   n.cni,
		store: n.store,
	}, nil
}

func (n *network) Attachment(ctx context.Context, id string) (Attachment, error) {
	log.G(ctx).WithField("network manager", n.NetworkInfo.Manager).WithField("network", n.NetworkInfo.Name).WithField("attachment", id).Debugf("get")
	// find attachment from store
	return &attachment{
		AttachmentInfo: AttachmentInfo{
			Manager: n.NetworkInfo.Manager,
			Network: n.NetworkInfo.Name,
		},
		cni:   n.cni,
		store: n.store,
	}, nil
}

func (n *network) Walk(ctx context.Context, fn func(context.Context, Attachment) error) error {
	log.G(ctx).WithField("network manager", n.NetworkInfo.Manager).WithField("network", n.NetworkInfo.Name).Debugf("list")
	if err := n.store.Walk(ctx, n.NetworkInfo.Manager, n.NetworkInfo.Name,
		func(r *AttachmentRecord) error {
			fn(ctx, &attachment{
				AttachmentInfo: AttachmentInfo{
					Manager: n.NetworkInfo.Manager,
					Network: n.NetworkInfo.Name,
				},
				cni:   n.cni,
				store: n.store,
			})
			return nil
		}); err != nil {
		return err
	}
	return nil
}

func generateAttachmentID(manager, network, namespace, container, ifname string) string {
	return fmt.Sprintf("%s/%s/%s/%s/%s", manager, network, namespace, container, ifname)
}
