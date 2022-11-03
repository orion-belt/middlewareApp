# middlewareApp
Middleware application for Orchestration OAI 5G Core Network


* Prerequisite
  * [Build and run Magma Orchestrator](https://docs.magmacore.org/docs/next/basics/quick_start_guide#terminal-tab-2-build-orchestrator)
  * [Build and run Magma Magma NMS](https://docs.magmacore.org/docs/next/basics/quick_start_guide#using-the-nms-ui)
  * [Deploy OAI 5G Core Network](https://gitlab.eurecom.fr/oai/cn5g/oai-cn5g-fed/-/blob/master/docs/DEPLOY_HOME.md)
  * Install Go and other dependencies for middlewareApp1

 ```bash
 wget https://artifactory.magmacore.org/artifactory/generic/go1.18.3.linux-amd64.tar.gz
 sudo rm -rf /usr/local/go && sudo tar -C /usr/local -xzf go1.18.3.linux-amd64.tar.gz
 export PATH=$PATH:/usr/local/go/bin
 go install github.com/fullstorydev/grpcurl/cmd/grpcurl@latest
 sudo apt install jq
 ```

* Build
```bash
git clone https://github.com/orion-belt/middlewareApp.git
cd middlewareApp
go mod tidy
go build .
```

* Copy Magma Orchestrator certificates to middlewareApp
```bash
cp $MAGMA_ROOT/.cache/test_certs/admin_operator.* middlewareApp/magmanbi/.certs/
```

* Run
```bash
./middlewareApp
```


* Build docker image and run middlewareApp as docker container [OPTIONAL]
```bash
docker build --target middlewareApp --tag middlewareapp:latest .
```