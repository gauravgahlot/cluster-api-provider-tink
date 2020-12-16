/*
Copyright 2020 The Kubernetes Authors.

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

package client

import (
	"context"
	"fmt"

	"github.com/tinkerbell/tink/protos/hardware"
	"google.golang.org/grpc"
)

// Hardware client for Tinkerbell.
type Hardware struct {
	client hardware.HardwareServiceClient
}

// NewHardwareClient returns a Hardware client.
func NewHardwareClient(client hardware.HardwareServiceClient) Hardware {
	return Hardware{client: client}
}

// Get returns a Tinkerbell Hardware.
func (t *Hardware) Get(ctx context.Context, id, ip, mac string) (*hardware.Hardware, error) {
	var method func(context.Context, *hardware.GetRequest, ...grpc.CallOption) (*hardware.Hardware, error)

	req := &hardware.GetRequest{}

	switch {
	case id != "":
		req.Id = id
		method = t.client.ByID
	case mac != "":
		req.Mac = mac
		method = t.client.ByMAC
	case ip != "":
		req.Ip = ip
		method = t.client.ByIP
	}

	tinkHardware, err := method(ctx, req)
	if err != nil {
		if err.Error() == sqlErrorString {
			return nil, fmt.Errorf("hardware %w", ErrNotFound)
		}

		return nil, fmt.Errorf("failed to get hardware from Tinkerbell: %w", err)
	}

	return tinkHardware, nil
}

// Delete a Tinkerbell Hardware.
func (t *Hardware) Delete(ctx context.Context, id string) error {
	if _, err := t.client.Delete(ctx, &hardware.DeleteRequest{Id: id}); err != nil {
		if err.Error() == sqlErrorString {
			return fmt.Errorf("hardware %w", ErrNotFound)
		}

		return fmt.Errorf("failed to delete hardware from Tinkerbell: %w", err)
	}

	return nil
}
