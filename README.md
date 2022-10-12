## SECRYPT
* Name - SECRYPT
* Symbol - SXC
* Supply - 1 billion
* Blocktime - 2 minutes
* Consensus - PoS
* P2P Port - 1478
* JSON-RPC Port - 8545 
* ChainID Main - 1143
* ChainID Test - 1144
* EVM Compatible

## Official Links
* Website - https://secrypt.tech
* Mainnet Explorer - https://explorer-mainnet.secrypt.tech
* Testnet Explorer - https://explorer-testnet.secrypt.tech
* Mainnet RPC - https://mainnet.secrypt.tech
* Testnet RPC - https://testnet.secrypt.tech

### Build from Source (Ubuntu 20.04)
Requirements - `Go >=1.17`

#### Setup Go Path
```
sudo nano ~/.profile
```
Paste this into the bottom of the file
```
export GOPATH=$HOME/work
export PATH=$PATH:/usr/local/go/bin:$GOPATH/bin
```
```
source ~/.profile
```

#### Install Go
```
wget https://go.dev/dl/go1.17.13.linux-amd64.tar.gz
sudo tar -xvf go1.17.13.linux-amd64.tar.gz
sudo mv go /usr/local && rm go1.17.13.linux-amd64.tar.gz
```
Check that it's installed
```
go version
```
You should see something like this:
```
go version go1.17.13 linux/amd64
```

#### Build SECRYPT
```
git clone https://github.com/SECRYPT-2022/SECRYPT.git
cd SECRYPT/
go build -o secrypt main.go
```

#### Running a Full Validating Node
After you have [downloaded](https://github.com/SECRYPT-2022/SECRYPT/releases/latest) the binaries or [built from source](https://github.com/SECRYPT-2022/SECRYPT#build-from-source), enter the `SECRYPT` directory and run the following:
```
./secrypt server --data-dir ~/.secrypt --chain mainnet-genesis.json --seal --max-slots 40960 --grpc 0.0.0.0:9632 --libp2p 0.0.0.0:1478 --jsonrpc 0.0.0.0:8545 --max-inbound-peers 128 --max-outbound-peers 16
```

#### Running a non-Validating node
```
./secrypt server --data-dir ~/.secrypt --chain mainnet-genesis.json --libp2p 0.0.0.0:1478 --nat <public_or_private_ip>
```


---
```
Copyright 2022 Polygon Technology
Copyright 2022 SECRYPT

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
```
