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

package main

import (
	"github.com/golang/glog"

	"middlewareApp/magmanbi/orc8r/cloud/go/blobstore"
	"middlewareApp/magmanbi/orc8r/cloud/go/orc8r"
	"middlewareApp/magmanbi/orc8r/cloud/go/service"
	"middlewareApp/magmanbi/orc8r/cloud/go/services/obsidian"
	swagger_protos "middlewareApp/magmanbi/orc8r/cloud/go/services/obsidian/swagger/protos"
	swagger_servicers "middlewareApp/magmanbi/orc8r/cloud/go/services/obsidian/swagger/servicers/protected"
	"middlewareApp/magmanbi/orc8r/cloud/go/services/tenants"
	"middlewareApp/magmanbi/orc8r/cloud/go/services/tenants/obsidian/handlers"
	"middlewareApp/magmanbi/orc8r/cloud/go/services/tenants/protos"
	servicers "middlewareApp/magmanbi/orc8r/cloud/go/services/tenants/servicers/protected"
	"middlewareApp/magmanbi/orc8r/cloud/go/services/tenants/servicers/storage"
	"middlewareApp/magmanbi/orc8r/cloud/go/sqorc"
	storage2 "middlewareApp/magmanbi/orc8r/cloud/go/storage"
)

func main() {
	srv, err := service.NewOrchestratorService(orc8r.ModuleName, tenants.ServiceName)
	if err != nil {
		glog.Fatalf("Error creating tenants service %s", err)
	}
	db, err := sqorc.Open(storage2.GetSQLDriver(), storage2.GetDatabaseSource())
	if err != nil {
		glog.Fatalf("Failed to connect to database: %s", err)
	}
	factory := blobstore.NewSQLStoreFactory(tenants.DBTableName, db, sqorc.GetSqlBuilder())
	err = factory.InitializeFactory()
	if err != nil {
		glog.Fatalf("Error initializing tenant database: %s", err)
	}
	store := storage.NewBlobstoreStore(factory)

	server, err := servicers.NewTenantsServicer(store)
	if err != nil {
		glog.Fatalf("Error creating tenants server: %s", err)
	}
	protos.RegisterTenantsServiceServer(srv.ProtectedGrpcServer, server)

	swagger_protos.RegisterSwaggerSpecServer(srv.ProtectedGrpcServer, swagger_servicers.NewSpecServicerFromFile(tenants.ServiceName))

	obsidian.AttachHandlers(srv.EchoServer, handlers.GetObsidianHandlers())

	err = srv.Run()
	if err != nil {
		glog.Fatalf("Error running service: %s", err)
	}
}
