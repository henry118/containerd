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

	api "github.com/containerd/containerd/api/services/networks/v1"
	"github.com/containerd/containerd/errdefs"
)

type proxy struct {
	manager string
	client  api.NetworkClient
}

var _ Service = (*proxy)(nil)

func NewProxy(manager string, client api.NetworkClient) Service {
	return &proxy{
		manager: manager,
		client:  client,
	}
}

func (p *proxy) CreateNetwork(ctx context.Context, name string, options ...NetworkOption) (NetworkInfo, error) {
	_, err := p.client.CreateNetwork(ctx,
		&api.CreateNetworkRequest{
			NetworkManager: p.manager,
			NetworkName:    name,
		})
	if err != nil {
		return NetworkInfo{}, errdefs.FromGRPC(err)
	}
	return NetworkInfo{
		Manager: p.manager,
		Name:    name,
		//Config:  config,
	}, nil
}

func (p *proxy) DeleteNetwork(ctx context.Context, name string) error {
	_, err := p.client.DeleteNetwork(ctx,
		&api.DeleteNetworkRequest{
			NetworkManager: p.manager,
			NetworkName:    name,
		})
	if err != nil {
		return errdefs.FromGRPC(err)
	}
	return nil
}

func (p *proxy) GetNetwork(ctx context.Context, name string) (NetworkInfo, error) {
	resp, err := p.client.GetNetwork(ctx,
		&api.GetNetworkRequest{
			NetworkManager: p.manager,
			NetworkName:    name,
		})
	if err != nil {
		return NetworkInfo{}, errdefs.FromGRPC(err)
	}
	return NetworkInfo{
		Manager: resp.NetworkManager,
		Name:    resp.NetworkConfig.Name,
	}, nil
}

func (p *proxy) ListNetworks(ctx context.Context) ([]NetworkInfo, error) {
	_, err := p.client.ListNetworks(ctx,
		&api.ListNetworksRequest{
			NetworkManager: p.manager,
		})
	if err != nil {
		return nil, errdefs.FromGRPC(err)
	}
	return nil, nil
}

func (p *proxy) AttachNetwork(ctx context.Context, network string, options ...AttachmentOption) (AttachmentInfo, error) {
	_, err := p.client.AttachNetwork(ctx,
		&api.AttachNetworkRequest{
			NetworkManager: p.manager,
			NetworkName:    network,
		})
	if err != nil {
		return AttachmentInfo{}, errdefs.FromGRPC(err)
	}
	return AttachmentInfo{
		//	Interfaces []Interface
		//Routes     []Route
		//DNS        DNS
	}, nil
}

func (p *proxy) DetachNetwork(ctx context.Context, network string, attachment string) error {
	_, err := p.client.DetachNetwork(ctx,
		&api.DetachNetworkRequest{
			NetworkManager: p.manager,
			NetworkName:    network,
			AttachmentId:   attachment,
		})
	if err != nil {
		return errdefs.FromGRPC(err)
	}
	return nil
}

func (p *proxy) GetAttachment(ctx context.Context, network string, attachment string) (AttachmentInfo, error) {
	_, err := p.client.GetAttachment(ctx,
		&api.GetAttachmentRequest{
			NetworkManager: p.manager,
			NetworkName:    network,
			AttachmentId:   attachment,
		})
	if err != nil {
		return AttachmentInfo{}, errdefs.FromGRPC(err)
	}
	return AttachmentInfo{
		Manager: p.manager,
		Network: network,
	}, nil
}

func (p *proxy) CheckAttachment(ctx context.Context, network string, attachment string) error {
	_, err := p.client.CheckAttachment(ctx,
		&api.CheckAttachmentRequest{
			NetworkManager: p.manager,
			NetworkName:    network,
			AttachmentId:   attachment,
		})
	if err != nil {
		return errdefs.FromGRPC(err)
	}
	return nil
}

func (p *proxy) ListAttachments(ctx context.Context, network string, options ...AttachmentFilterOption) ([]AttachmentInfo, error) {
	_, err := p.client.ListAttachments(ctx,
		&api.ListAttachmentsRequest{
			NetworkManager: p.manager,
			NetworkName:    network,
		})
	if err != nil {
		return nil, errdefs.FromGRPC(err)
	}
	return nil, nil
}
