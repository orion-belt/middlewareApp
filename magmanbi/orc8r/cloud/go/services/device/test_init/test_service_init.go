/*
 * Copyright 2020 The Magma Authors.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package test_init

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"middlewareApp/magmanbi/orc8r/cloud/go/orc8r"
	"middlewareApp/magmanbi/orc8r/cloud/go/services/device"
	"middlewareApp/magmanbi/orc8r/cloud/go/services/device/protos"
	servicers "middlewareApp/magmanbi/orc8r/cloud/go/services/device/servicers/protected"
	"middlewareApp/magmanbi/orc8r/cloud/go/test_utils"
)

// StartTestService instantiates a service backed by an in-memory storage
func StartTestService(t *testing.T) {
	factory := test_utils.NewSQLBlobstore(t, "device_test_service_blobstore")
	srv, lis, plis := test_utils.NewTestService(t, orc8r.ModuleName, device.ServiceName)
	server, err := servicers.NewDeviceServicer(factory)
	assert.NoError(t, err)
	protos.RegisterDeviceServer(srv.ProtectedGrpcServer, server)
	go srv.RunTest(lis, plis)
}
