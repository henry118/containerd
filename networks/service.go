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

import "context"

type Service interface {
	CreateNetwork(ctx context.Context, name string, options ...NetworkOption) (NetworkInfo, error)
	DeleteNetwork(ctx context.Context, name string) error
	GetNetwork(ctx context.Context, name string) (NetworkInfo, error)
	ListNetworks(ctx context.Context) ([]NetworkInfo, error)
	AttachNetwork(ctx context.Context, network string, options ...AttachmentOption) (AttachmentInfo, error)
	DetachNetwork(ctx context.Context, network string, attachment string) error
	GetAttachment(ctx context.Context, network string, attachment string) (AttachmentInfo, error)
	CheckAttachment(ctx context.Context, network string, attachment string) error
	ListAttachments(ctx context.Context, network string, options ...AttachmentFilterOption) ([]AttachmentInfo, error)
}

type Network interface {
	Name(ctx context.Context) string
	Delete(ctx context.Context) error
	Info(ctx context.Context) NetworkInfo
	Attach(ctx context.Context, options ...AttachmentOption) (Attachment, error)
	Attachment(ctx context.Context, id string) (Attachment, error)
	Walk(ctx context.Context, fn func(context.Context, Attachment) error) error
}

type Attachment interface {
	ID(ctx context.Context) string
	Remove(ctx context.Context) error
	Info(ctx context.Context) AttachmentInfo
	Check(ctx context.Context) error
}

type Manager interface {
	Create(ctx context.Context, name string, options ...NetworkOption) (Network, error)
	Delete(ctx context.Context, name string) error
	Network(ctx context.Context, name string) (Network, error)
	Walk(ctx context.Context, fn func(context.Context, Network) error) error
}

type NetworkInfo struct {
	Manager string
	Name    string
	Config  NetworkConfig
}

type AttachmentInfo struct {
	ID         string
	Manager    string
	Network    string
	Lease      string
	Interfaces []Interface
	Routes     []Route
	DNS        DNS
}
