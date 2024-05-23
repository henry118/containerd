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

/*
   This file is copied and customized based on
   https://github.com/moby/moby/blob/master/pkg/idtools/idtools.go
*/

package userns

import (
	"fmt"
	"math"

	"github.com/opencontainers/runtime-spec/specs-go"
)

var (
	invalidUser = User{Uid: math.MaxUint32, Gid: math.MaxUint32}
)

// User is a Uid and Gid pair of a user
//
//nolint:revive
type User struct {
	Uid uint32
	Gid uint32
}

// IDMap contains the mappings of Uids and Gids.
//
//nolint:revive
type IDMap struct {
	UidMap []specs.LinuxIDMapping `json:"UidMap"`
	GidMap []specs.LinuxIDMapping `json:"GidMap"`
}

// RootID returns the ID pair for the root user
func (i IDMap) RootPair() (User, error) {
	uid, err := toHost(0, i.UidMap)
	if err != nil {
		return invalidUser, err
	}
	gid, err := toHost(0, i.GidMap)
	if err != nil {
		return invalidUser, err
	}
	return User{Uid: uid, Gid: gid}, nil
}

// ToHost returns the host User pair for the container uid, gid.
func (i IDMap) ToHost(pair User) (User, error) {
	var (
		target User
		err    error
	)
	target.Uid, err = toHost(pair.Uid, i.UidMap)
	if err != nil {
		return invalidUser, err
	}
	target.Gid, err = toHost(pair.Gid, i.GidMap)
	if err != nil {
		return invalidUser, err
	}
	return target, nil
}

// ToContainer returns the container Identify pair for the host uid and gid
func (i IDMap) ToContainer(pair User) (User, error) {
	var (
		target User
		err    error
	)
	target.Uid, err = toContainer(pair.Uid, i.UidMap)
	if err != nil {
		return invalidUser, err
	}
	target.Gid, err = toContainer(pair.Gid, i.GidMap)
	if err != nil {
		return invalidUser, err
	}
	return target, nil
}

// Empty returns true if there are no id mappings
func (i IDMap) Empty() bool {
	return len(i.UidMap) == 0 && len(i.GidMap) == 0
}

// toContainer takes an id mapping, and uses it to translate a
// host ID to the remapped ID. If no map is provided, then the translation
// assumes a 1-to-1 mapping and returns the passed in id
func toContainer(hostID uint32, idMap []specs.LinuxIDMapping) (uint32, error) {
	if idMap == nil {
		return hostID, nil
	}
	for _, m := range idMap {
		if (hostID >= m.HostID) && (hostID <= (m.HostID + m.Size - 1)) {
			contID := m.ContainerID + (hostID - m.HostID)
			return contID, nil
		}
	}
	return math.MaxUint32, fmt.Errorf("host ID %d cannot be mapped to a container ID", hostID)
}

// toHost takes an id mapping and a remapped ID, and translates the
// ID to the mapped host ID. If no map is provided, then the translation
// assumes a 1-to-1 mapping and returns the passed in id #
func toHost(contID uint32, idMap []specs.LinuxIDMapping) (uint32, error) {
	if idMap == nil {
		return contID, nil
	}
	for _, m := range idMap {
		if (contID >= m.ContainerID) && (contID <= (m.ContainerID + m.Size - 1)) {
			hostID := m.HostID + (contID - m.ContainerID)
			return hostID, nil
		}
	}
	return math.MaxUint32, fmt.Errorf("container ID %d cannot be mapped to a host ID", contID)
}
