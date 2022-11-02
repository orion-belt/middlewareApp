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

package main

import (
	"github.com/golang/glog"

	"middlewareApp/magmanbi/orc8r/cloud/go/blobstore"
	"middlewareApp/magmanbi/orc8r/cloud/go/orc8r"
	"middlewareApp/magmanbi/orc8r/cloud/go/service"
	"middlewareApp/magmanbi/orc8r/cloud/go/services/ctraced"
	"middlewareApp/magmanbi/orc8r/cloud/go/services/ctraced/obsidian/handlers"
	servicers "middlewareApp/magmanbi/orc8r/cloud/go/services/ctraced/servicers/southbound"
	ctraced_storage "middlewareApp/magmanbi/orc8r/cloud/go/services/ctraced/storage"
	"middlewareApp/magmanbi/orc8r/cloud/go/services/obsidian"
	swagger_protos "middlewareApp/magmanbi/orc8r/cloud/go/services/obsidian/swagger/protos"
	swagger_servicers "middlewareApp/magmanbi/orc8r/cloud/go/services/obsidian/swagger/servicers/protected"
	"middlewareApp/magmanbi/orc8r/cloud/go/sqorc"
	"middlewareApp/magmanbi/orc8r/cloud/go/storage"
	"middlewareApp/magmanbi/orc8r/lib/go/protos"
)

func main() {
	// Create service
	srv, err := service.NewOrchestratorService(orc8r.ModuleName, ctraced.ServiceName)
	if err != nil {
		glog.Fatalf("Error creating ctraced service: %+v", err)
	}

	// Init storage
	db, err := sqorc.Open(storage.GetSQLDriver(), storage.GetDatabaseSource())
	if err != nil {
		glog.Fatalf("Error opening db connection: %+v", err)
	}
	fact := blobstore.NewSQLStoreFactory(ctraced.LookupTableBlobstore, db, sqorc.GetSqlBuilder())
	err = fact.InitializeFactory()
	if err != nil {
		glog.Fatalf("Error initializing ctraced table: %+v", err)
	}
	ctracedBlobstore := ctraced_storage.NewCtracedBlobstore(fact)

	// Init gRPC servicer
	protos.RegisterCallTraceControllerServer(srv.GrpcServer, servicers.NewCallTraceServicer(ctracedBlobstore))
	swagger_protos.RegisterSwaggerSpecServer(srv.ProtectedGrpcServer, swagger_servicers.NewSpecServicerFromFile(ctraced.ServiceName))

	gwClient := handlers.NewGwCtracedClient()
	obsidian.AttachHandlers(srv.EchoServer, handlers.GetObsidianHandlers(gwClient, ctracedBlobstore))

	// Run service
	err = srv.Run()
	if err != nil {
		glog.Fatalf("Error running ctraced service: %+v", err)
	}
}
