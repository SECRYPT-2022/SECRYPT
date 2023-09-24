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

## Build from Source (Ubuntu 20.04 + MacOS)
Requirements - `Go >=1.18.x`

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

### Install Go
```
wget https://go.dev/dl/go1.18.7.linux-amd64.tar.gz
sudo tar -xvf go1.18.7.linux-amd64.tar.gz
sudo mv go /usr/local && rm go1.18.7.linux-amd64.tar.gz

```
Check that it's installed
```
go version

```
You should see something like this:
```
go version go1.18.7 linux/amd64
```

### Prerequisites

```
sudo apt-get install build-essential
sudo apt-get update
sudo apt install build-essential

```

### Port Open

```
sudo ufw allow 30301
sudo ufw allow 30302
sudo ufw allow 30303
sudo ufw allow 26656
sudo ufw allow 26660
sudo ufw allow 8545
sudo ufw allow 8546
sudo ufw allow 1478
sudo ufw allow 7071
sudo ufw allow 1317
sudo ufw allow 9632

```

### Build SECRYPT
```
git clone https://github.com/SECRYPT-2022/SECRYPT.git
cd SECRYPT/
go build -o secrypt main.go
go mod tidy
go build -o secrypt main.go

```

### Create data directory
```
mkdir ~/.secrypt

```

## Running a Node
### Running a Full Validating Node
After you have [downloaded](https://github.com/SECRYPT-2022/SECRYPT/releases/latest) the binaries or [built from source](https://github.com/SECRYPT-2022/SECRYPT#build-from-source), go [here](ValidatorGuide.md) and follow the guide:

### Running a non-Validating node
```
./secrypt server --data-dir ~/.secrypt --chain mainnet-genesis.json --libp2p 0.0.0.0:1478 --nat <public_or_private_ip>
```

### Running a Full non-Validating node
```
./secrypt server --data-dir ~/.secrypt --chain mainnet-genesis.json --jsonrpc 0.0.0.0:8545 --libp2p 0.0.0.0:1478 --grpc 0.0.0.0:9632 --max-inbound-peers 128 --max-outbound-peers 16 --max-slots 40960 --nat <public_or_private_ip>
```

---
```
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
