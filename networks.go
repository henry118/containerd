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

package containerd

import (
	"context"
	"fmt"

	"github.com/containerd/containerd/errdefs"
	"github.com/containerd/containerd/networks"
)

var (
	_ networks.Network    = (*network)(nil)
	_ networks.Manager    = (*manager)(nil)
	_ networks.Attachment = (*attachment)(nil)
)

type network struct {
	networks.NetworkInfo
	client *Client
}

type manager struct {
	name   string
	client *Client
}

type attachment struct {
	networks.AttachmentInfo
	client *Client
}

func NewNetworkManager(name string, client *Client) networks.Manager {
	return &manager{
		name:   name,
		client: client,
	}
}

func (n *network) Name(ctx context.Context) string {
	return n.NetworkInfo.Name
}

func (n *network) Info(ctx context.Context) networks.NetworkInfo {
	return n.NetworkInfo
}

func (n *network) Delete(ctx context.Context) error {
	s := n.client.NetworkService(n.Manager)
	if s == nil {
		return fmt.Errorf("")
	}
	return s.DeleteNetwork(ctx, n.NetworkInfo.Name)
}

func (n *network) Attach(ctx context.Context, options ...networks.AttachmentOption) (networks.Attachment, error) {
	s := n.client.NetworkService(n.Manager)
	if s == nil {
		return nil, fmt.Errorf("")
	}
	a, err := s.AttachNetwork(ctx, n.NetworkInfo.Name)
	if err != nil {
		return nil, fmt.Errorf("")
	}
	return &attachment{
		AttachmentInfo: a,
		client:         n.client,
	}, nil
}

func (n *network) Attachment(ctx context.Context, id string) (networks.Attachment, error) {
	s := n.client.NetworkService(n.Manager)
	if s == nil {
		return nil, fmt.Errorf("")
	}
	a, err := s.GetAttachment(ctx, n.NetworkInfo.Name, id)
	if err != nil {
		return nil, errdefs.FromGRPC(err)
	}
	return &attachment{
		AttachmentInfo: a,
		client:         n.client,
	}, nil
}

func (n *network) Walk(ctx context.Context, fn func(context.Context, networks.Attachment) error) error {
	s := n.client.NetworkService(n.Manager)
	if s == nil {
		return fmt.Errorf("")
	}
	al, err := s.ListAttachments(ctx, n.NetworkInfo.Name)
	if err != nil {
		return errdefs.FromGRPC(err)
	}
	for _, a := range al {
		fn(ctx, &attachment{AttachmentInfo: a, client: n.client})
	}
	return nil
}

func (m *manager) Create(ctx context.Context, name string, options ...networks.NetworkOption) (networks.Network, error) {
	s := m.client.NetworkService(m.name)
	if s == nil {
		return nil, fmt.Errorf("")
	}
	n, err := s.CreateNetwork(ctx, name, options...)
	if err != nil {
		return nil, errdefs.FromGRPC(err)
	}
	return &network{
		NetworkInfo: n,
		client:      m.client,
	}, nil
}

func (m *manager) Delete(ctx context.Context, name string) error {
	s := m.client.NetworkService(m.name)
	if s == nil {
		return fmt.Errorf("")
	}
	err := s.DeleteNetwork(ctx, m.name)
	if err != nil {
		return errdefs.FromGRPC(err)
	}
	return nil

}

func (m *manager) Network(ctx context.Context, name string) (networks.Network, error) {
	s := m.client.NetworkService(m.name)
	if s == nil {
		return nil, fmt.Errorf("")
	}
	n, err := s.GetNetwork(ctx, name)
	if err != nil {
		return nil, errdefs.FromGRPC(err)
	}
	return &network{
		NetworkInfo: n,
		client:      m.client,
	}, nil
}

func (m *manager) Walk(ctx context.Context, fn func(context.Context, networks.Network) error) error {
	s := m.client.NetworkService(m.name)
	if s == nil {
		return fmt.Errorf("")
	}
	nl, err := s.ListNetworks(ctx)
	if err != nil {
		return errdefs.FromGRPC(err)
	}
	for _, i := range nl {
		fn(ctx, &network{
			NetworkInfo: i,
			client:      m.client,
		})
	}
	return nil
}

func (a *attachment) ID(ctx context.Context) string {
	return a.AttachmentInfo.ID
}

func (a *attachment) Check(ctx context.Context) error {
	s := a.client.NetworkService(a.AttachmentInfo.Manager)
	if s == nil {
		return fmt.Errorf("")
	}
	err := s.CheckAttachment(ctx, a.AttachmentInfo.Network, a.AttachmentInfo.ID)
	if err != nil {
		return errdefs.FromGRPC(err)
	}
	return nil
}

func (a *attachment) Remove(ctx context.Context) error {
	s := a.client.NetworkService(a.AttachmentInfo.Manager)
	if s == nil {
		return fmt.Errorf("")
	}
	err := s.DetachNetwork(ctx, a.AttachmentInfo.Network, a.AttachmentInfo.ID)
	if err != nil {
		return errdefs.FromGRPC(err)
	}
	return nil
}

func (a *attachment) Info(ctx context.Context) networks.AttachmentInfo {
	return a.AttachmentInfo
}
