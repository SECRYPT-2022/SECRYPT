package command

import "github.com/SECRYPT-2022/SECRYPT/server"

const (
	DefaultGenesisFileName = "genesis.json"
	DefaultChainName       = "secrypt"
	DefaultChainID         = 1143
	DefaultPremineBalance  = "0x33B2E3C9FD0803CE8000000" // 1 billion units of native network currency
	DefaultConsensus       = server.IBFTConsensus
	DefaultGenesisGasUsed  = 458752  // 0x70000
	DefaultGenesisGasLimit = 30000000 // 0x500000
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
