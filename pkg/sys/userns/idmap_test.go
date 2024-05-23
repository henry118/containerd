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

package userns

import (
	"os"
	"testing"

	"github.com/opencontainers/runtime-spec/specs-go"
	"github.com/stretchr/testify/assert"
)

func TestRootPair(t *testing.T) {
	idmap := IDMap{
		UidMap: []specs.LinuxIDMapping{
			{
				ContainerID: 0,
				HostID:      uint32(os.Getuid()),
				Size:        1,
			},
		},
		GidMap: []specs.LinuxIDMapping{
			{
				ContainerID: 0,
				HostID:      uint32(os.Getgid()),
				Size:        1,
			},
		},
	}

	identity, err := idmap.RootPair()
	assert.NoError(t, err)
	assert.Equal(t, uint32(os.Geteuid()), identity.Uid)
	assert.Equal(t, uint32(os.Getegid()), identity.Gid)

	idmapErr := IDMap{
		UidMap: []specs.LinuxIDMapping{
			{
				ContainerID: 1,
				HostID:      uint32(os.Getuid()),
				Size:        1,
			},
		},
	}
	_, err = idmapErr.RootPair()
	assert.EqualError(t, err, "container ID 0 cannot be mapped to a host ID")
}

func TestToContainer(t *testing.T) {
	idmap := IDMap{
		UidMap: []specs.LinuxIDMapping{
			{
				ContainerID: 0,
				HostID:      1,
				Size:        2,
			},
			{
				ContainerID: 2,
				HostID:      4,
				Size:        1000,
			},
		},
		GidMap: []specs.LinuxIDMapping{
			{
				ContainerID: 0,
				HostID:      2,
				Size:        4,
			},
			{
				ContainerID: 4,
				HostID:      8,
				Size:        1000,
			},
		},
	}
	for _, tt := range []struct {
		cid User
		hid User
	}{
		{
			hid: User{
				Uid: 1,
				Gid: 2,
			},
			cid: User{
				Uid: 0,
				Gid: 0,
			},
		},
		{
			hid: User{
				Uid: 2,
				Gid: 3,
			},
			cid: User{
				Uid: 1,
				Gid: 1,
			},
		},
		{
			hid: User{
				Uid: 4,
				Gid: 8,
			},
			cid: User{
				Uid: 2,
				Gid: 4,
			},
		},
		{
			hid: User{
				Uid: 102,
				Gid: 204,
			},
			cid: User{
				Uid: 100,
				Gid: 200,
			},
		},
		{
			hid: User{
				Uid: 1003,
				Gid: 1007,
			},
			cid: User{
				Uid: 1001,
				Gid: 1003,
			},
		},
		{
			hid: User{
				Uid: 1004,
				Gid: 1008,
			},
			cid: invalidUser,
		},
		{
			hid: User{
				Uid: 2000,
				Gid: 2000,
			},
			cid: invalidUser,
		},
	} {
		r, err := idmap.ToContainer(tt.hid)
		assert.Equal(t, tt.cid, r)
		if r == invalidUser {
			assert.Error(t, err)
		} else {
			assert.NoError(t, err)
		}
	}
}

func TestToHost(t *testing.T) {
	idmap := IDMap{
		UidMap: []specs.LinuxIDMapping{
			{
				ContainerID: 0,
				HostID:      1,
				Size:        2,
			},
			{
				ContainerID: 2,
				HostID:      4,
				Size:        1000,
			},
		},
		GidMap: []specs.LinuxIDMapping{
			{
				ContainerID: 0,
				HostID:      2,
				Size:        4,
			},
			{
				ContainerID: 4,
				HostID:      8,
				Size:        1000,
			},
		},
	}
	for _, tt := range []struct {
		cid User
		hid User
	}{
		{
			cid: User{
				Uid: 0,
				Gid: 0,
			},
			hid: User{
				Uid: 1,
				Gid: 2,
			},
		},
		{
			cid: User{
				Uid: 1,
				Gid: 1,
			},
			hid: User{
				Uid: 2,
				Gid: 3,
			},
		},
		{
			cid: User{
				Uid: 2,
				Gid: 4,
			},
			hid: User{
				Uid: 4,
				Gid: 8,
			},
		},
		{
			cid: User{
				Uid: 100,
				Gid: 200,
			},
			hid: User{
				Uid: 102,
				Gid: 204,
			},
		},
		{
			cid: User{
				Uid: 1001,
				Gid: 1003,
			},
			hid: User{
				Uid: 1003,
				Gid: 1007,
			},
		},
		{
			cid: User{
				Uid: 1004,
				Gid: 1008,
			},
			hid: invalidUser,
		},
		{
			cid: User{
				Uid: 2000,
				Gid: 2000,
			},
			hid: invalidUser,
		},
	} {
		r, err := idmap.ToHost(tt.cid)
		assert.Equal(t, tt.hid, r)
		if r == invalidUser {
			assert.Error(t, err)
		} else {
			assert.NoError(t, err)
		}
	}
}
