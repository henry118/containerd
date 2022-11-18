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

	"github.com/containerd/containerd/log"
	"github.com/containernetworking/cni/libcni"
)

type attachment struct {
	AttachmentInfo
	cni   libcni.CNI
	store Store
}

var _ Attachment = (*attachment)(nil)

func (a *attachment) ID(ctx context.Context) string {
	return a.AttachmentInfo.ID
}

func (a *attachment) Remove(ctx context.Context) error {
	log.G(ctx).WithField("network manager", a.AttachmentInfo.Manager).WithField("network", a.AttachmentInfo.Network).WithField("attachment", a.AttachmentInfo.ID).Debugf("remove")
	// find attachment from store

	// find networkcfg

	// call libcni delete attachment

	// delete attachment from store

	/*
		var args attachArgs
		var netns, ifname string
		// TODO: load args from store
		// TODO: delete instance from store
		return env.cni.DelNetworkList(ctx, cfg, args.config(network.Container, netns, ifname))
	*/
	return nil
}

func (a *attachment) Info(ctx context.Context) AttachmentInfo {
	return a.AttachmentInfo
}

func (a *attachment) Check(ctx context.Context) error {
	log.G(ctx).WithField("network manager", a.AttachmentInfo.Manager).WithField("network", a.AttachmentInfo.Network).WithField("attachment", a.AttachmentInfo.ID).Debugf("check")
	// find attachment from store

	// find network config

	// call libcni

	//return env.cni.CheckNetworkList(ctx, cfg, args.config(network.Container, netns, ifname))
	return nil
}
