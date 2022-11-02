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
	"context"
	"database/sql"
	"time"

	"github.com/golang/glog"

	"magma/orc8r/lib/go/service/config"
	"middlewareApp/magmanbi/orc8r/cloud/go/blobstore"
	"middlewareApp/magmanbi/orc8r/cloud/go/orc8r"
	"middlewareApp/magmanbi/orc8r/cloud/go/service"
	"middlewareApp/magmanbi/orc8r/cloud/go/services/state"
	state_config "middlewareApp/magmanbi/orc8r/cloud/go/services/state/config"
	"middlewareApp/magmanbi/orc8r/cloud/go/services/state/indexer/reindex"
	"middlewareApp/magmanbi/orc8r/cloud/go/services/state/metrics"
	indexer_protos "middlewareApp/magmanbi/orc8r/cloud/go/services/state/protos"
	protected_servicers "middlewareApp/magmanbi/orc8r/cloud/go/services/state/servicers/protected"
	servicers "middlewareApp/magmanbi/orc8r/cloud/go/services/state/servicers/southbound"
	"middlewareApp/magmanbi/orc8r/cloud/go/sqorc"
	"middlewareApp/magmanbi/orc8r/cloud/go/storage"
	"middlewareApp/magmanbi/orc8r/lib/go/protos"
)

// how often to report gateway status
const gatewayStatusReportInterval = time.Second * 60

const nonPostgresDriverMessage = `Configuration warning:

This deployment has automatic state reindexing enabled, but is targeting a
database driver other than Postgres. This will cause the state service
to log a (harmless) DB syntax error, due to its use of Postgres-specific
syntax for automatic reindexing.

(Option 1) Continue using non-Postgres driver. To clear this warning, update
the state.yml cloud config to set enable_automatic_reindexing to false.
Keep in mind that, for this option, you will have to perform manual state
reindexing on every Orc8r upgrade. We provide a CLI to manage this, and will
provide directions in the upgrade notes.

(Option 2) Switch to a Postgres driver.
`

func main() {
	srv, err := service.NewOrchestratorService(orc8r.ModuleName, state.ServiceName)
	if err != nil {
		glog.Fatalf("Error creating state service %v", err)
	}

	db, err := sqorc.Open(storage.GetSQLDriver(), storage.GetDatabaseSource())
	if err != nil {
		glog.Fatalf("Error connecting to database: %v", err)
	}
	store := blobstore.NewSQLStoreFactory(state.DBTableName, db, sqorc.GetSqlBuilder())
	err = store.InitializeFactory()
	if err != nil {
		glog.Fatalf("Error initializing state database: %v", err)
	}

	stateServicer := newStateServicer(store)
	protos.RegisterStateServiceServer(srv.GrpcServer, stateServicer)

	cloudStateServicer := newCloudStateServicer(store)
	protos.RegisterCloudStateServiceServer(srv.ProtectedGrpcServer, cloudStateServicer)

	singletonReindex := srv.Config.MustGetBool(state_config.EnableSingletonReindex)
	if !singletonReindex {
		glog.Info("Running reindexer")
		indexerManagerServer := newIndexerManagerServicer(srv.Config, db, store)
		indexer_protos.RegisterIndexerManagerServer(srv.ProtectedGrpcServer, indexerManagerServer)
	}

	go metrics.PeriodicallyReportGatewayStatus(gatewayStatusReportInterval)

	err = srv.Run()
	if err != nil {
		glog.Fatalf("Error running state service: %v", err)
	}
}

func newStateServicer(store blobstore.StoreFactory) protos.StateServiceServer {
	servicer, err := servicers.NewStateServicer(store)
	if err != nil {
		glog.Fatalf("Error creating state servicer: %v", err)
	}
	return servicer
}

func newCloudStateServicer(store blobstore.StoreFactory) protos.CloudStateServiceServer {
	servicer, err := protected_servicers.NewCloudStateServicer(store)
	if err != nil {
		glog.Fatalf("Error creating state servicer: %v", err)
	}
	return servicer
}

func newIndexerManagerServicer(cfg *config.Map, db *sql.DB, store blobstore.StoreFactory) indexer_protos.IndexerManagerServer {
	queue := reindex.NewSQLJobQueue(reindex.DefaultMaxAttempts, db, sqorc.GetSqlBuilder())
	err := queue.Initialize()
	if err != nil {
		glog.Fatal("Error initializing state reindex queue")
	}
	_, err = queue.PopulateJobs()
	if err != nil {
		glog.Fatalf("Unexpected error initializing reindex job queue: %s", err)
	}

	autoReindex := cfg.MustGetBool(state_config.EnableAutomaticReindexing)
	reindexer := reindex.NewReindexerQueue(queue, reindex.NewStore(store))
	servicer := protected_servicers.NewIndexerManagerServicer(reindexer, autoReindex)

	if autoReindex && storage.GetSQLDriver() != sqorc.PostgresDriver {
		glog.Warning(nonPostgresDriverMessage)
	}

	if autoReindex {
		glog.Info("Automatic reindexing enabled for state service")
		go reindexer.Run(context.Background())
	} else {
		glog.Info("Automatic reindexing disabled for state service")
	}

	return servicer
}
