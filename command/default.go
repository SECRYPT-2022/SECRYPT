package command

import "github.com/SECRYPT-2022/SECRYPT/server"

const (
	DefaultGenesisFileName = "genesis.json"
	DefaultChainName       = "secrypt"
	DefaultChainID         = 1143
	DefaultPremineBalance  = "0xD3C21BCECCEDA1000000" // 1 million units of native network currency
	DefaultConsensus       = server.IBFTConsensus
	DefaultGenesisGasUsed  = 458752  // 0x70000
	DefaultGenesisGasLimit = 5242880 // 0x500000
)

const (
	JSONOutputFlag  = "json"
	GRPCAddressFlag = "grpc-address"
	JSONRPCFlag     = "jsonrpc"
)

// GRPCAddressFlagLEGACY Legacy flag that needs to be present to preserve backwards
// compatibility with running clients
const (
	GRPCAddressFlagLEGACY = "grpc"
)
