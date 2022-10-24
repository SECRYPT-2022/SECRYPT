# How to setup a Validator on SECRYPT Testnet

## Part 1 - Build from Source (Ubuntu 20.04)
Follow [this guide](https://github.com/SECRYPT-2022/SECRYPT/blob/main/README.md#build-from-source-ubuntu-2004) to download and compile the binaries. 

## Part 2 - Generate Validator Keys
Create a directory to store your keys and testnet blockchain data
```
mkdir ~/.testsxc
cd ~/SECRYPT
```

Run this to generate your keys
```
./secrypt secrets init  --data-dir ~/.testsxc
```

You should see something like this:
```
[SECRETS INIT]
Public key (address) = 0xXXXXXXXXXX...
BLS Public key       = 0xXXXXXXXXXX...
Node ID              = 16UXXXXXXXXX...
```

Save that somewhere and continue to the next step. 

## Part 3 - Stake Your Coins
### Setup staking contracts
Download staking contracts:
```
git clone https://github.com/SECRYPT-2022/staking-contracts.git
cd staking-contracts
npm i
cp .env.example .env
```

Now open up the .env file and fill out the missing parts. Fill out the variables like so:
```
JSONRPC_URL=http://ttv-wkqbnkghpr.dynamic-m.com:8544
PRIVATE_KEYS=YOUR_VALIDATOR_PRIVATE_KEY_HERE
STAKING_CONTRACT_ADDRESS=0x0000000000000000000000000000000000001001
```

The `YOUR_VALIDATOR_PRIVATE_KEY_HERE` variable can be found in this file: 
```
~/.testsxc/consensus/validator.key
```

Save the `.env` file and continue on

### Fund your wallet
Send at least 8.1 SXC to the wallet address that you generated in Part 2: 
```
Public key (address) = 0xXXXXX...
```

### Stake to become a validator
Once you have completed the steps above, run the following:
```
npm run stake
```

It run run a few seconds until the transaction confirms. Once the script is done running, you can run this script to confirm that your validator was added as well as the total staked amount and list of validators. 
```
npm run info
```

When you are done staking and want to un-stake, run this command:
```
npm run info
```

## Part 4 - Start Your Validator
```
./secrypt server --data-dir ~/.testsxc --chain testnet-genesis.json --seal --max-slots 40960 --grpc 0.0.0.0:9632 --libp2p 0.0.0.0:1478 --jsonrpc 0.0.0.0:8545 --max-inbound-peers 128 --max-outbound-peers 16
```

Once you run that, your validator will start syncing with the testnet blockchain and begin validating. 


