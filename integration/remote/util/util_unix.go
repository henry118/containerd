//go:build unix

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
Copyright 2017 The Kubernetes Authors.

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

package util

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/url"
	"os"
	"path/filepath"

	"golang.org/x/sys/unix"
)

const (
	// unixProtocol is the network protocol of unix socket.
	unixProtocol = "unix"
)

// CreateListener creates a listener on the specified endpoint.
func CreateListener(endpoint string) (net.Listener, error) {
	protocol, addr, err := parseEndpointWithFallbackProtocol(endpoint, unixProtocol)
	if err != nil {
		return nil, err
	}
	if protocol != unixProtocol {
		return nil, fmt.Errorf("only support unix socket endpoint")
	}

	// Unlink to cleanup the previous socket file.
	err = unix.Unlink(addr)
	if err != nil && !os.IsNotExist(err) {
		return nil, fmt.Errorf("failed to unlink socket file %q: %v", addr, err)
	}

	if err := os.MkdirAll(filepath.Dir(addr), 0750); err != nil {
		return nil, fmt.Errorf("error creating socket directory %q: %v", filepath.Dir(addr), err)
	}

	// Create the socket on a tempfile and move it to the destination socket to handle improprer cleanup
	file, err := os.CreateTemp(filepath.Dir(addr), "")
	if err != nil {
		return nil, fmt.Errorf("failed to create temporary file: %v", err)
	}

	if err := os.Remove(file.Name()); err != nil {
		return nil, fmt.Errorf("failed to remove temporary file: %v", err)
	}

	l, err := net.Listen(protocol, file.Name())
	if err != nil {
		return nil, err
	}

	if err = os.Rename(file.Name(), addr); err != nil {
		return nil, fmt.Errorf("failed to move temporary file to addr %q: %v", addr, err)
	}

	return l, nil
}

// GetAddressAndDialer returns the address parsed from the given endpoint and a context dialer.
func GetAddressAndDialer(endpoint string) (string, func(ctx context.Context, addr string) (net.Conn, error), error) {
	protocol, addr, err := parseEndpointWithFallbackProtocol(endpoint, unixProtocol)
	if err != nil {
		return "", nil, err
	}
	if protocol != unixProtocol {
		return "", nil, errors.New("only support unix socket endpoint")
	}

	// Use passthrough as the scheme so it allows us to use our custom dialer:
	//
	// "grpc.Dial uses "passthrough" as the default name resolver for backward compatibility while grpc.NewClient
	// uses "dns" as its default name resolver. This subtle difference is important to legacy systems that also
	// specified a custom dialer and expected it to receive the target string directly."
	// https://github.com/grpc/grpc-go/blob/master/Documentation/anti-patterns.md#the-wrong-way-grpcdial
	return fmt.Sprintf("passthrough:///%s", addr), dial, nil
}

func dial(ctx context.Context, addr string) (net.Conn, error) {
	return (&net.Dialer{}).DialContext(ctx, unixProtocol, addr)
}

func parseEndpointWithFallbackProtocol(endpoint string, fallbackProtocol string) (protocol string, addr string, err error) {
	if protocol, addr, err = parseEndpoint(endpoint); err != nil && protocol == "" {
		fallbackEndpoint := fallbackProtocol + "://" + endpoint
		protocol, addr, err = parseEndpoint(fallbackEndpoint)
	}
	return
}

func parseEndpoint(endpoint string) (string, string, error) {
	u, err := url.Parse(endpoint)
	if err != nil {
		return "", "", err
	}

	switch u.Scheme {
	case "tcp":
		return "tcp", u.Host, nil

	case "unix":
		return "unix", u.Path, nil

	case "":
		return "", "", fmt.Errorf("using %q as endpoint is deprecated, please consider using full url format", endpoint)

	default:
		return u.Scheme, "", fmt.Errorf("protocol %q not supported", u.Scheme)
	}
}

// IsUnixDomainSocket returns whether a given file is a AF_UNIX socket file
func IsUnixDomainSocket(filePath string) (bool, error) {
	fi, err := os.Stat(filePath)
	if err != nil {
		return false, fmt.Errorf("stat file %s failed: %v", filePath, err)
	}
	if fi.Mode()&os.ModeSocket == 0 {
		return false, nil
	}
	return true, nil
}

// NormalizePath is a no-op for Linux for now
func NormalizePath(path string) string {
	return path
}
