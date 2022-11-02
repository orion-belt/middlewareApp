/*
Copyright 2020 The Magma Authors.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package servicers

import (
	"context"
	"fmt"

	"magma/orc8r/lib/go/definitions"
	streamer_protos "middlewareApp/magmanbi/orc8r/cloud/go/services/streamer/protos"
	stream_provider "middlewareApp/magmanbi/orc8r/cloud/go/services/streamer/providers/servicers/protected"
	"middlewareApp/magmanbi/orc8r/lib/go/protos"
)

type providerServicer struct{}

func NewProviderServicer() streamer_protos.StreamProviderServer {
	return &providerServicer{}
}

func (s *providerServicer) GetUpdates(ctx context.Context, req *protos.StreamRequest) (*protos.DataUpdateBatch, error) {
	var streamer stream_provider.StreamProvider
	switch req.GetStreamName() {
	case definitions.MconfigStreamName:
		streamer = &stream_provider.MconfigProvider{}
	default:
		return nil, fmt.Errorf("GetUpdates failed: unknown stream name provided: %s", req.GetStreamName())
	}

	update, err := streamer.GetUpdates(ctx, req.GetGatewayId(), req.GetExtraArgs())
	if err != nil {
		// Note: return blank err to properly receive EAGAIN from mconfig provider
		return nil, err
	}
	return &protos.DataUpdateBatch{Updates: update}, nil
}