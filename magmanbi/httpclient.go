package magmanbi

import (
	"crypto/tls"
	"log"
	"net/http"
	"os"
	"middlewareApp/common"
)

func GetHttpClient() *http.Client {
	base_path, _ := os.Getwd()
	cert_path := base_path + "/magmanbi/.certs/"
	cert, err := tls.LoadX509KeyPair(cert_path+"admin_operator.pem", cert_path+"admin_operator.key.pem")
	if err != nil {
		log.Fatalf(
			"ERROR loading orchestrator certificate and key ('%s')",
			err,
		)
	}
	return &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
				Certificates:       []tls.Certificate{cert},
			},
		},
	}
}

func SendHttpRequest(method, url,
	payload string,
	// optional headers, each header is a slice of strings in the form:
	//   [<header name>, value1, value2...]
	headers ...[]string,
) (int, string, error) {
	var client = GetHttpClient()
	response, contents, error :=common.SendHttpRequest(method, url, payload, client, headers ...,)
	return response, string(contents), error
}
