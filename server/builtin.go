package server

import (
	"github.com/SECRYPT-2022/SECRYPT/consensus"
	consensusDev "github.com/SECRYPT-2022/SECRYPT/consensus/dev"
	consensusDummy "github.com/SECRYPT-2022/SECRYPT/consensus/dummy"
	consensusIBFT "github.com/SECRYPT-2022/SECRYPT/consensus/ibft"
	"github.com/SECRYPT-2022/SECRYPT/secrets"
	"github.com/SECRYPT-2022/SECRYPT/secrets/awsssm"
	"github.com/SECRYPT-2022/SECRYPT/secrets/gcpssm"
	"github.com/SECRYPT-2022/SECRYPT/secrets/hashicorpvault"
	"github.com/SECRYPT-2022/SECRYPT/secrets/local"
)

type ConsensusType string

const (
	DevConsensus   ConsensusType = "dev"
	IBFTConsensus  ConsensusType = "ibft"
	DummyConsensus ConsensusType = "dummy"
)

var consensusBackends = map[ConsensusType]consensus.Factory{
	DevConsensus:   consensusDev.Factory,
	IBFTConsensus:  consensusIBFT.Factory,
	DummyConsensus: consensusDummy.Factory,
}

// secretsManagerBackends defines the SecretManager factories for different
// secret management solutions
var secretsManagerBackends = map[secrets.SecretsManagerType]secrets.SecretsManagerFactory{
	secrets.Local:          local.SecretsManagerFactory,
	secrets.HashicorpVault: hashicorpvault.SecretsManagerFactory,
	secrets.AWSSSM:         awsssm.SecretsManagerFactory,
	secrets.GCPSSM:         gcpssm.SecretsManagerFactory,
}

func ConsensusSupported(value string) bool {
	_, ok := consensusBackends[ConsensusType(value)]

	return ok
}
