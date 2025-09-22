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

package unpack

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/containerd/containerd/v2/core/diff"
	"github.com/containerd/containerd/v2/core/images"
	"github.com/containerd/containerd/v2/core/mount"
	"github.com/containerd/errdefs"
	"github.com/containerd/log"
	"github.com/opencontainers/go-digest"
	ocispec "github.com/opencontainers/image-spec/specs-go/v1"
)

func (u *Unpacker) supportParrallelUnpack(unpack *Platform) bool {
	return u.parallel && slices.Contains(unpack.SnapshotterCapabilities, "parallel-unpack")
}

func tempUnpackRoot(unpack *Platform) string {
	root := unpack.SnapshotterExports["unpack"]
	if len(root) > 0 {
		return root
	}
	root = unpack.SnapshotterExports["root"]
	if len(root) > 0 {
		return filepath.Join(root, "unpack")
	}
	return filepath.Join(getTempDir(), "unpack")
}

func tempUnpackHandler(h images.Handler, applier diff.Applier, unpackRoot string) images.HandlerFunc {
	return images.HandlerFunc(func(ctx context.Context, desc ocispec.Descriptor) ([]ocispec.Descriptor, error) {
		children, err := h.Handle(ctx, desc)
		if err != nil {
			return children, err
		}

		if images.IsLayerType(desc.MediaType) {
			unpackDir := filepath.Join(unpackRoot, uniquePart())
			fsDir := filepath.Join(unpackDir, "fs")
			mounts := []mount.Mount{
				{
					Type:    "overlay",
					Source:  "overlay",
					Options: []string{fmt.Sprintf("upperdir=%s", fsDir)},
				},
			}
			if err := os.MkdirAll(fsDir, 0711); err != nil {
				return nil, err
			}
			diff, err := applier.Apply(ctx, desc, mounts)
			if err != nil {
				return nil, err
			}
			if err := os.WriteFile(filepath.Join(unpackDir, "digest"), []byte(desc.Digest.String()), 0600); err != nil {
				return nil, err
			}
			if err := os.WriteFile(filepath.Join(unpackDir, "diff"), []byte(diff.Digest.String()), 0600); err != nil {
				return nil, err
			}
		}

		return children, err
	})
}

func tempRebasePath(mounts []mount.Mount) (string, error) {
	if len(mounts) == 1 {
		switch mounts[0].Type {
		case "overlay":
			path, _, err := tempOverlayPath(mounts[0].Options)
			if err != nil {
				return "", err
			}
			return path, nil
		case "bind":
			return mounts[0].Source, nil
		}
	}
	return "", fmt.Errorf("unable to determine rebase path from mounts: %+v", mounts)
}

func tempOverlayPath(options []string) (upper string, lower []string, err error) {
	const upperdirPrefix = "upperdir="
	const lowerdirPrefix = "lowerdir="

	for _, o := range options {
		if strings.HasPrefix(o, upperdirPrefix) {
			upper = strings.TrimPrefix(o, upperdirPrefix)
		} else if strings.HasPrefix(o, lowerdirPrefix) {
			lower = strings.Split(strings.TrimPrefix(o, lowerdirPrefix), ":")
		}
	}
	if upper == "" {
		return "", nil, fmt.Errorf("upperdir not found: %w", errdefs.ErrInvalidArgument)
	}

	return
}

func tempRebaseSnapshot(ctx context.Context, mounts []mount.Mount, desc ocispec.Descriptor, unpackRoot string) (ocispec.Descriptor, error) {
	dest, err := tempRebasePath(mounts)
	if err != nil {
		return ocispec.Descriptor{}, err
	}
	unpacks, err := os.ReadDir(unpackRoot)
	if err != nil {
		return ocispec.Descriptor{}, err
	}
	rebase := func(source, dest string) (ocispec.Descriptor, error) {
		var (
			rdesc ocispec.Descriptor
			rerr  error
			rdgst digest.Digest
			rb    []byte
		)
		if rb, rerr = os.ReadFile(filepath.Join(source, "digest")); rerr != nil {
			return rdesc, rerr
		}
		if rdgst, rerr = digest.Parse(string(rb)); rerr != nil {
			return rdesc, rerr
		}
		if rdgst != desc.Digest {
			return rdesc, fmt.Errorf("digest does not match")
		}
		if rerr = os.Rename(filepath.Join(source, "fs"), dest); rerr != nil {
			return rdesc, rerr
		}
		if rb, rerr = os.ReadFile(filepath.Join(source, "diff")); rerr != nil {
			return rdesc, rerr
		}
		if rdesc.Digest, rerr = digest.Parse(string(rb)); rerr != nil {
			return rdesc, rerr
		}
		os.RemoveAll(source)
		log.G(ctx).Debugf("rebased layer %s to %s with diff %s", rdgst, dest, rdesc.Digest)
		return rdesc, nil
	}

	if err := os.RemoveAll(dest); err != nil {
		return ocispec.Descriptor{}, err
	}
	for _, u := range unpacks {
		if u.IsDir() {
			diff, err := rebase(filepath.Join(unpackRoot, u.Name()), dest)
			if err == nil {
				return diff, nil
			}
		}
	}

	return ocispec.Descriptor{}, fmt.Errorf("unable to find unpacked layer for %s: %w", desc.Digest, errdefs.ErrNotFound)
}

func getTempDir() string {
	if xdg := os.Getenv("XDG_RUNTIME_DIR"); xdg != "" {
		return xdg
	}
	return os.TempDir()
}
