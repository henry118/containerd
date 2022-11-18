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
	"errors"

	api "github.com/containerd/containerd/api/services/networks/v1"
	"github.com/containerd/containerd/errdefs"
	"github.com/containerd/containerd/log"
	"github.com/containerd/containerd/networks"
	"github.com/containerd/containerd/plugin"
	ptypes "github.com/containerd/containerd/protobuf/types"
	"github.com/containerd/containerd/services"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

var (
	_     api.NetworkServer = (*service)(nil)
	empty *ptypes.Empty     = &ptypes.Empty{}
)

func init() {
	plugin.Register(&plugin.Registration{
		Type: plugin.GRPCPlugin,
		ID:   "network",
		Requires: []plugin.Type{
			plugin.ServicePlugin,
		},
		InitFn: func(ic *plugin.InitContext) (interface{}, error) {
			plugins, err := ic.GetByType(plugin.ServicePlugin)
			if err != nil {
				return nil, err
			}
			p, ok := plugins[services.NetworkService]
			if !ok {
				return nil, errors.New("networks manager not found")
			}
			i, err := p.Instance()
			if err != nil {
				return nil, err
			}
			return &service{
				locals: i.(map[string]networks.Service),
			}, nil
		},
	})
}

type service struct {
	locals map[string]networks.Service
	api.UnimplementedNetworkServer
}

func (s *service) Register(gs *grpc.Server) error {
	api.RegisterNetworkServer(gs, s)
	return nil
}

func (s *service) getService(name string) (networks.Service, error) {
	if name == "" {
		return nil, errdefs.ToGRPCf(errdefs.ErrInvalidArgument, "network manager name missing")
	}

	m := s.locals[name]
	if m == nil {
		return nil, errdefs.ToGRPCf(errdefs.ErrInvalidArgument, "network manager not configured: %s", name)
	}

	return m, nil
}

func (s *service) CreateNetwork(ctx context.Context, r *api.CreateNetworkRequest) (*emptypb.Empty, error) {
	log.G(ctx).WithField("manager", r.NetworkManager).WithField("network", r.NetworkName).Debug("create network")
	ns, err := s.getService(r.NetworkManager)
	if err != nil {
		return nil, err
	}
	_, err = ns.CreateNetwork(ctx, r.NetworkName)
	if err != nil {
		return nil, errdefs.ToGRPC(err)
	}
	return empty, nil
}

func (s *service) DeleteNetwork(ctx context.Context, r *api.DeleteNetworkRequest) (*emptypb.Empty, error) {
	log.G(ctx).WithField("manager", r.NetworkManager).WithField("network", r.NetworkName).Debug("delete network")
	ns, err := s.getService(r.NetworkManager)
	if err != nil {
		return nil, err
	}
	err = ns.DeleteNetwork(ctx, r.NetworkName)
	if err != nil {
		return nil, errdefs.ToGRPC(err)
	}
	return empty, nil
}

func (s *service) GetNetwork(ctx context.Context, r *api.GetNetworkRequest) (*api.GetNetworkResponse, error) {
	log.G(ctx).WithField("manager", r.NetworkManager).WithField("network", r.NetworkName).Debug("get network")
	ns, err := s.getService(r.NetworkManager)
	if err != nil {
		return nil, err
	}
	n, err := ns.GetNetwork(ctx, r.NetworkName)
	if err != nil {
		return nil, errdefs.ToGRPC(err)
	}
	return &api.GetNetworkResponse{
		NetworkManager: n.Manager,
		//Config:         n.Info(ctx).Config,
	}, nil
}

func (s *service) ListNetworks(ctx context.Context, r *api.ListNetworksRequest) (*api.ListNetworksResponse, error) {
	log.G(ctx).WithField("manager", r.NetworkManager).Debug("list networks")
	ns, err := s.getService(r.NetworkManager)
	if err != nil {
		return nil, err
	}
	_, err = ns.ListNetworks(ctx)
	if err != nil {
		return nil, errdefs.ToGRPC(err)
	}
	return &api.ListNetworksResponse{}, nil
}

func (s *service) AttachNetwork(ctx context.Context, r *api.AttachNetworkRequest) (*api.AttachNetworkResponse, error) {
	log.G(ctx).WithField("manager", r.NetworkManager).WithField("network", r.NetworkName).Debugf("attach network")
	ns, err := s.getService(r.NetworkManager)
	if err != nil {
		return nil, err
	}
	_, err = ns.AttachNetwork(ctx, r.NetworkName)
	if err != nil {
		return nil, errdefs.ToGRPC(err)
	}
	return &api.AttachNetworkResponse{
		NetworkManager: r.NetworkManager,
	}, nil
}

func (s *service) DetachNetwork(ctx context.Context, r *api.DetachNetworkRequest) (*emptypb.Empty, error) {
	log.G(ctx).WithField("manager", r.NetworkManager).WithField("network", r.NetworkName).WithField("attachment", r.AttachmentId).Debugf("detach network")
	ns, err := s.getService(r.NetworkManager)
	if err != nil {
		return nil, err
	}
	err = ns.DetachNetwork(ctx, r.NetworkName, r.AttachmentId)
	if err != nil {
		return nil, errdefs.ToGRPC(err)
	}
	return empty, nil
}

func (s *service) GetAttachment(ctx context.Context, r *api.GetAttachmentRequest) (*api.GetAttachmentResponse, error) {
	log.G(ctx).WithField("manager", r.NetworkManager).WithField("network", r.NetworkName).WithField("attachment", r.AttachmentId).Debugf("get attachment")
	ns, err := s.getService(r.NetworkManager)
	_, err = ns.GetAttachment(ctx, r.NetworkName, r.AttachmentId)
	if err != nil {
		return nil, errdefs.ToGRPC(err)
	}
	return &api.GetAttachmentResponse{
		NetworkManager: r.NetworkManager,
	}, nil

}
func (s *service) CheckAttachment(ctx context.Context, r *api.CheckAttachmentRequest) (*api.CheckAttachmentResponse, error) {
	log.G(ctx).WithField("manager", r.NetworkManager).WithField("network", r.NetworkName).WithField("attachment", r.AttachmentId).Debugf("detach network")
	ns, err := s.getService(r.NetworkManager)
	if err != nil {
		return nil, err
	}
	err = ns.CheckAttachment(ctx, r.NetworkName, r.AttachmentId)
	return &api.CheckAttachmentResponse{
		NetworkManager: r.NetworkManager,
		NetworkName:    r.NetworkName,
		AttachmentId:   r.AttachmentId,
		Ok:             err != nil,
	}, nil

}
func (s *service) ListAttachments(ctx context.Context, r *api.ListAttachmentsRequest) (*api.ListAttachmentsResponse, error) {
	log.G(ctx).WithField("manager", r.NetworkManager).WithField("network", r.NetworkName).Debugf("list attachments")
	ns, err := s.getService(r.NetworkManager)
	_, err = ns.ListAttachments(ctx, r.NetworkName)
	if err != nil {
		return nil, errdefs.ToGRPC(err)
	}
	return &api.ListAttachmentsResponse{
		NetworkManager: r.NetworkManager,
		//Network:        r.NetworkName,
	}, nil
}
