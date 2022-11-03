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

// Package registry for Magma microservices

package magmanbi

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"flag"
	"fmt"
	"github.com/golang/glog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/backoff"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/keepalive"
	"middlewareApp/config"
	"time"
)

const (
	ControlProxyServiceName = "CONTROL_PROXY"
	grpcMaxTimeoutSec       = 60
	grpcMaxDelaySec         = 20
)

var (
	defaultTimeoutDuration = grpcMaxTimeoutSec * time.Second
	keepaliveParams        = keepalive.ClientParameters{
		Time:                59 * time.Second,
		Timeout:             20 * time.Second,
		PermitWithoutStream: true,
	}
	proxiedKeepaliveParams = keepalive.ClientParameters{
		Time:                47 * time.Second,
		Timeout:             10 * time.Second,
		PermitWithoutStream: true,
	}
	grpcKeepAliveFlag = flag.Bool("grpc_keepalive", false, "Use keepalive option for all GRPC connections")
)

func GetCloudConnection(authority string, addr string) (*grpc.ClientConn, error) {
	ctx, cancel := context.WithTimeout(context.Background(), grpcMaxTimeoutSec*time.Second)
	defer cancel()

	opts, err := getDialOptions(authority)
	if err != nil {
		return nil, err
	}

	glog.V(2).Infof("connecting to: %s, authority: %s", addr, authority)
	conn, err := grpc.DialContext(ctx, addr, opts...)
	// conn, err := grpc.Dial(addr, opts...)

	if err != nil {
		return nil, fmt.Errorf("Address: %s GRPC Dial error: %s", addr, err)
	} else if ctx.Err() != nil {
		return nil, ctx.Err()
	}

	return conn, nil
}
func getDialOptions(authority string) ([]grpc.DialOption, error) {
	bckoff := backoff.DefaultConfig
	bckoff.MaxDelay = grpcMaxDelaySec * time.Second
	var opts = []grpc.DialOption{
		grpc.WithConnectParams(grpc.ConnectParams{
			Backoff:           bckoff,
			MinConnectTimeout: grpcMaxTimeoutSec * time.Second,
		}),
		grpc.WithBlock(),
		grpc.WithUnaryInterceptor(TimeoutInterceptor),
	}

	// always try to add OS certs
	certPool, err := x509.SystemCertPool()
	if err != nil {
		glog.Warningf("OS Cert Pool initialization error: %v", err)
		certPool = x509.NewCertPool()
	}

	tlsCfg := &tls.Config{ServerName: authority}
	if len(certPool.Subjects()) > 0 {
		tlsCfg.RootCAs = certPool
	} else {
		glog.Warning("Empty server certificate pool, using TLS InsecureSkipVerify")

		tlsCfg.InsecureSkipVerify = true
	}
	clientCaFile, clientKeyFile := config.GetGatewayCerds()

	clientCert, err := tls.LoadX509KeyPair(clientCaFile, clientKeyFile)
	if err == nil {
		tlsCfg.Certificates = []tls.Certificate{clientCert}
	} else {
		glog.Errorf("failed to load Client Certificate & Key from '%s', '%s': %v",
			clientCaFile, clientKeyFile, err)
	}

	tlsCfg.InsecureSkipVerify = true

	opts = append(opts, grpc.WithTransportCredentials(credentials.NewTLS(tlsCfg)))
	if *grpcKeepAliveFlag {
		opts = append(opts, grpc.WithKeepaliveParams(keepaliveParams))
	}

	return opts, nil
}

// TimeoutInterceptor is a generic client connection interceptor which sets default timeout option for RPC if the
// currently used CTX does not already specify it's own deadline option
func TimeoutInterceptor(ctx context.Context, method string, req, resp interface{},
	cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {

	// check if given CTX already has a deadline & only add default deadline if not
	if _, deadlineIsSet := ctx.Deadline(); !deadlineIsSet {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, defaultTimeoutDuration)
		// cleanup timer after invoke call chain completion
		defer cancel()
	}
	return invoker(ctx, method, req, resp, cc, opts...)
}
