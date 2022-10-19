package staking

import (
	"fmt"
	"math/big"

	"github.com/SECRYPT-2022/SECRYPT/chain"
	"github.com/SECRYPT-2022/SECRYPT/helper/common"
	"github.com/SECRYPT-2022/SECRYPT/helper/hex"
	"github.com/SECRYPT-2022/SECRYPT/helper/keccak"
	"github.com/SECRYPT-2022/SECRYPT/types"
	"github.com/SECRYPT-2022/SECRYPT/validators"
)

var (
	MinValidatorCount = uint64(1)
	MaxValidatorCount = common.MaxSafeJSInt
)

// getAddressMapping returns the key for the SC storage mapping (address => something)
//
// More information:
// https://docs.soliditylang.org/en/latest/internals/layout_in_storage.html
func getAddressMapping(address types.Address, slot int64) []byte {
	bigSlot := big.NewInt(slot)

	finalSlice := append(
		common.PadLeftOrTrim(address.Bytes(), 32),
		common.PadLeftOrTrim(bigSlot.Bytes(), 32)...,
	)

	return keccak.Keccak256(nil, finalSlice)
}

// getIndexWithOffset is a helper method for adding an offset to the already found keccak hash
func getIndexWithOffset(keccakHash []byte, offset uint64) []byte {
	bigOffset := big.NewInt(int64(offset))
	bigKeccak := big.NewInt(0).SetBytes(keccakHash)

	bigKeccak.Add(bigKeccak, bigOffset)

	return bigKeccak.Bytes()
}

// getStorageIndexes is a helper function for getting the correct indexes
// of the storage slots which need to be modified during bootstrap.
//
// It is SC dependant, and based on the SC located at:
// https://github.com/SECRYPT-2022/staking-contracts/
func getStorageIndexes(validator validators.Validator, index int) *StorageIndexes {
	storageIndexes := &StorageIndexes{}
	address := validator.Addr()

	// Get the indexes for the mappings
	// The index for the mapping is retrieved with:
	// keccak(address . slot)
	// . stands for concatenation (basically appending the bytes)
	storageIndexes.AddressToIsValidatorIndex = getAddressMapping(
		address,
		addressToIsValidatorSlot,
	)

	storageIndexes.AddressToStakedAmountIndex = getAddressMapping(
		address,
		addressToStakedAmountSlot,
	)

	storageIndexes.AddressToValidatorIndexIndex = getAddressMapping(
		address,
		addressToValidatorIndexSlot,
	)

	storageIndexes.ValidatorBLSPublicKeyIndex = getAddressMapping(
		address,
		addressToBLSPublicKeySlot,
	)

	// Index for array types is calculated as keccak(slot) + index
	// The slot for the dynamic arrays that's put in the keccak needs to be in hex form (padded 64 chars)
	storageIndexes.ValidatorsIndex = getIndexWithOffset(
		keccak.Keccak256(nil, common.PadLeftOrTrim(big.NewInt(validatorsSlot).Bytes(), 32)),
		uint64(index),
	)

	return storageIndexes
}

// setBytesToStorage sets bytes data into storage map from specified base index
func setBytesToStorage(
	storageMap map[types.Hash]types.Hash,
	baseIndexBytes []byte,
	data []byte,
) {
	dataLen := len(data)
	baseIndex := types.BytesToHash(baseIndexBytes)

	if dataLen <= 31 {
		bytes := types.Hash{}

		copy(bytes[:len(data)], data)

		// Set 2*Size at the first byte
		bytes[len(bytes)-1] = byte(dataLen * 2)

		storageMap[baseIndex] = bytes

		return
	}

	// Set size at the base index
	baseSlot := types.Hash{}
	baseSlot[31] = byte(2*dataLen + 1)
	storageMap[baseIndex] = baseSlot

	zeroIndex := keccak.Keccak256(nil, baseIndexBytes)
	numBytesInSlot := 256 / 8

	for i := 0; i < dataLen; i++ {
		offset := i / numBytesInSlot

		slotIndex := types.BytesToHash(getIndexWithOffset(zeroIndex, uint64(offset)))
		byteIndex := i % numBytesInSlot

		slot := storageMap[slotIndex]
		slot[byteIndex] = data[i]

		storageMap[slotIndex] = slot
	}
}

// PredeployParams contains the values used to predeploy the PoS staking contract
type PredeployParams struct {
	MinValidatorCount uint64
	MaxValidatorCount uint64
}

// StorageIndexes is a wrapper for different storage indexes that
// need to be modified
type StorageIndexes struct {
	ValidatorsIndex              []byte // []address
	ValidatorBLSPublicKeyIndex   []byte // mapping(address => byte[])
	AddressToIsValidatorIndex    []byte // mapping(address => bool)
	AddressToStakedAmountIndex   []byte // mapping(address => uint256)
	AddressToValidatorIndexIndex []byte // mapping(address => uint256)
}

// Slot definitions for SC storage
var (
	validatorsSlot              = int64(0) // Slot 0
	addressToIsValidatorSlot    = int64(1) // Slot 1
	addressToStakedAmountSlot   = int64(2) // Slot 2
	addressToValidatorIndexSlot = int64(3) // Slot 3
	stakedAmountSlot            = int64(4) // Slot 4
	minNumValidatorSlot         = int64(5) // Slot 5
	maxNumValidatorSlot         = int64(6) // Slot 6
	addressToBLSPublicKeySlot   = int64(7) // Slot 7
)

const (
	DefaultStakedBalance = "0x6F05B59D3B200000" // 8 ETH
	//nolint: lll
	StakingSCBytecode = "0x6080604052600436106101185760003560e01c80637a6eea37116100a0578063d94c111b11610064578063d94c111b1461040a578063e387a7ed14610433578063e804fbf61461045e578063f90ecacc14610489578063facd743b146104c657610186565b80637a6eea37146103215780637dceceb81461034c578063af6da36e14610389578063c795c077146103b4578063ca1e7819146103df57610186565b8063373d6132116100e7578063373d6132146102595780633a4b66f1146102845780633c561f041461028e57806351a9ab32146102b9578063714ff425146102f657610186565b806302b751991461018b578063065ae171146101c85780632367f6b5146102055780632def66201461024257610186565b366101865761013c3373ffffffffffffffffffffffffffffffffffffffff16610503565b1561017c576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610173906117d8565b60405180910390fd5b610184610516565b005b600080fd5b34801561019757600080fd5b506101b260048036038101906101ad91906113ce565b6105ed565b6040516101bf9190611833565b60405180910390f35b3480156101d457600080fd5b506101ef60048036038101906101ea91906113ce565b610605565b6040516101fc919061173b565b60405180910390f35b34801561021157600080fd5b5061022c600480360381019061022791906113ce565b610625565b6040516102399190611833565b60405180910390f35b34801561024e57600080fd5b5061025761066e565b005b34801561026557600080fd5b5061026e610759565b60405161027b9190611833565b60405180910390f35b61028c610763565b005b34801561029a57600080fd5b506102a36107cc565b6040516102b09190611719565b60405180910390f35b3480156102c557600080fd5b506102e060048036038101906102db91906113ce565b610972565b6040516102ed9190611756565b60405180910390f35b34801561030257600080fd5b5061030b610a12565b6040516103189190611833565b60405180910390f35b34801561032d57600080fd5b50610336610a1c565b6040516103439190611818565b60405180910390f35b34801561035857600080fd5b50610373600480360381019061036e91906113ce565b610a28565b6040516103809190611833565b60405180910390f35b34801561039557600080fd5b5061039e610a40565b6040516103ab9190611833565b60405180910390f35b3480156103c057600080fd5b506103c9610a46565b6040516103d69190611833565b60405180910390f35b3480156103eb57600080fd5b506103f4610a4c565b60405161040191906116f7565b60405180910390f35b34801561041657600080fd5b50610431600480360381019061042c91906113fb565b610ada565b005b34801561043f57600080fd5b50610448610b7f565b6040516104559190611833565b60405180910390f35b34801561046a57600080fd5b50610473610b85565b6040516104809190611833565b60405180910390f35b34801561049557600080fd5b506104b060048036038101906104ab9190611444565b610b8f565b6040516104bd91906116dc565b60405180910390f35b3480156104d257600080fd5b506104ed60048036038101906104e891906113ce565b610bce565b6040516104fa919061173b565b60405180910390f35b600080823b905060008111915050919050565b34600460008282546105289190611954565b9250508190555034600260003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020600082825461057e9190611954565b9250508190555061058e33610c24565b1561059d5761059c33610c9c565b5b3373ffffffffffffffffffffffffffffffffffffffff167f9e71bc8eea02a63969f509818f2dafb9254532904319f9dbda79b67bd34a5f3d346040516105e39190611833565b60405180910390a2565b60036020528060005260406000206000915090505481565b60016020528060005260406000206000915054906101000a900460ff1681565b6000600260008373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020549050919050565b61068d3373ffffffffffffffffffffffffffffffffffffffff16610503565b156106cd576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016106c4906117d8565b60405180910390fd5b6000600260003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020541161074f576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161074690611778565b60405180910390fd5b610757610deb565b565b6000600454905090565b6107823373ffffffffffffffffffffffffffffffffffffffff16610503565b156107c2576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016107b9906117d8565b60405180910390fd5b6107ca610516565b565b60606000808054905067ffffffffffffffff8111156107ee576107ed611bec565b5b60405190808252806020026020018201604052801561082157816020015b606081526020019060019003908161080c5790505b50905060005b60008054905081101561096a576007600080838154811061084b5761084a611bbd565b5b9060005260206000200160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002080546108bb90611a84565b80601f01602080910402602001604051908101604052809291908181526020018280546108e790611a84565b80156109345780601f1061090957610100808354040283529160200191610934565b820191906000526020600020905b81548152906001019060200180831161091757829003601f168201915b505050505082828151811061094c5761094b611bbd565b5b6020026020010181905250808061096290611ae7565b915050610827565b508091505090565b6007602052806000526040600020600091509050805461099190611a84565b80601f01602080910402602001604051908101604052809291908181526020018280546109bd90611a84565b8015610a0a5780601f106109df57610100808354040283529160200191610a0a565b820191906000526020600020905b8154815290600101906020018083116109ed57829003601f168201915b505050505081565b6000600554905090565b676f05b59d3b20000081565b60026020528060005260406000206000915090505481565b60065481565b60055481565b60606000805480602002602001604051908101604052809291908181526020018280548015610ad057602002820191906000526020600020905b8160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019060010190808311610a86575b5050505050905090565b80600760003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020016000209080519060200190610b2d929190611291565b503373ffffffffffffffffffffffffffffffffffffffff167f472da4d064218fa97032725fbcff922201fa643fed0765b5ffe0ceef63d7b3dc82604051610b749190611756565b60405180910390a250565b60045481565b6000600654905090565b60008181548110610b9f57600080fd5b906000526020600020016000915054906101000a900473ffffffffffffffffffffffffffffffffffffffff1681565b6000600160008373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060009054906101000a900460ff169050919050565b6000610c2f82610f3d565b158015610c955750676f05b59d3b2000006fffffffffffffffffffffffffffffffff16600260008473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020016000205410155b9050919050565b60065460008054905010610ce5576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610cdc90611798565b60405180910390fd5b60018060008373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060006101000a81548160ff021916908315150217905550600080549050600360008373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020819055506000819080600181540180825580915050600190039060005260206000200160009091909190916101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555050565b6000600260003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020016000205490506000600260003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020819055508060046000828254610e8691906119aa565b92505081905550610e9633610f3d565b15610ea557610ea433610f93565b5b3373ffffffffffffffffffffffffffffffffffffffff166108fc829081150290604051600060405180830381858888f19350505050158015610eeb573d6000803e3d6000fd5b503373ffffffffffffffffffffffffffffffffffffffff167f0f5bb82176feb1b5e747e28471aa92156a04d9f3ab9f45f28e2d704232b93f7582604051610f329190611833565b60405180910390a250565b6000600160008373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060009054906101000a900460ff169050919050565b60055460008054905011610fdc576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610fd3906117f8565b60405180910390fd5b600080549050600360008373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020016000205410611062576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401611059906117b8565b60405180910390fd5b6000600360008373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020549050600060016000805490506110ba91906119aa565b90508082146111a85760008082815481106110d8576110d7611bbd565b5b9060005260206000200160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff169050806000848154811061111a57611119611bbd565b5b9060005260206000200160006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555082600360008373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002081905550505b6000600160008573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060006101000a81548160ff0219169083151502179055506000600360008573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002081905550600080548061125757611256611b8e565b5b6001900381819060005260206000200160006101000a81549073ffffffffffffffffffffffffffffffffffffffff02191690559055505050565b82805461129d90611a84565b90600052602060002090601f0160209004810192826112bf5760008555611306565b82601f106112d857805160ff1916838001178555611306565b82800160010185558215611306579182015b828111156113055782518255916020019190600101906112ea565b5b5090506113139190611317565b5090565b5b80821115611330576000816000905550600101611318565b5090565b600061134761134284611873565b61184e565b90508281526020810184848401111561136357611362611c20565b5b61136e848285611a42565b509392505050565b60008135905061138581611d59565b92915050565b600082601f8301126113a05761139f611c1b565b5b81356113b0848260208601611334565b91505092915050565b6000813590506113c881611d70565b92915050565b6000602082840312156113e4576113e3611c2a565b5b60006113f284828501611376565b91505092915050565b60006020828403121561141157611410611c2a565b5b600082013567ffffffffffffffff81111561142f5761142e611c25565b5b61143b8482850161138b565b91505092915050565b60006020828403121561145a57611459611c2a565b5b6000611468848285016113b9565b91505092915050565b600061147d838361149d565b60208301905092915050565b6000611495838361159d565b905092915050565b6114a6816119de565b82525050565b6114b5816119de565b82525050565b60006114c6826118c4565b6114d081856118ff565b93506114db836118a4565b8060005b8381101561150c5781516114f38882611471565b97506114fe836118e5565b9250506001810190506114df565b5085935050505092915050565b6000611524826118cf565b61152e8185611910565b935083602082028501611540856118b4565b8060005b8581101561157c578484038952815161155d8582611489565b9450611568836118f2565b925060208a01995050600181019050611544565b50829750879550505050505092915050565b611597816119f0565b82525050565b60006115a8826118da565b6115b28185611921565b93506115c2818560208601611a51565b6115cb81611c2f565b840191505092915050565b60006115e1826118da565b6115eb8185611932565b93506115fb818560208601611a51565b61160481611c2f565b840191505092915050565b600061161c601d83611943565b915061162782611c40565b602082019050919050565b600061163f602783611943565b915061164a82611c69565b604082019050919050565b6000611662601283611943565b915061166d82611cb8565b602082019050919050565b6000611685601a83611943565b915061169082611ce1565b602082019050919050565b60006116a8604083611943565b91506116b382611d0a565b604082019050919050565b6116c7816119fc565b82525050565b6116d681611a38565b82525050565b60006020820190506116f160008301846114ac565b92915050565b6000602082019050818103600083015261171181846114bb565b905092915050565b600060208201905081810360008301526117338184611519565b905092915050565b6000602082019050611750600083018461158e565b92915050565b6000602082019050818103600083015261177081846115d6565b905092915050565b600060208201905081810360008301526117918161160f565b9050919050565b600060208201905081810360008301526117b181611632565b9050919050565b600060208201905081810360008301526117d181611655565b9050919050565b600060208201905081810360008301526117f181611678565b9050919050565b600060208201905081810360008301526118118161169b565b9050919050565b600060208201905061182d60008301846116be565b92915050565b600060208201905061184860008301846116cd565b92915050565b6000611858611869565b90506118648282611ab6565b919050565b6000604051905090565b600067ffffffffffffffff82111561188e5761188d611bec565b5b61189782611c2f565b9050602081019050919050565b6000819050602082019050919050565b6000819050602082019050919050565b600081519050919050565b600081519050919050565b600081519050919050565b6000602082019050919050565b6000602082019050919050565b600082825260208201905092915050565b600082825260208201905092915050565b600082825260208201905092915050565b600082825260208201905092915050565b600082825260208201905092915050565b600061195f82611a38565b915061196a83611a38565b9250827fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0382111561199f5761199e611b30565b5b828201905092915050565b60006119b582611a38565b91506119c083611a38565b9250828210156119d3576119d2611b30565b5b828203905092915050565b60006119e982611a18565b9050919050565b60008115159050919050565b60006fffffffffffffffffffffffffffffffff82169050919050565b600073ffffffffffffffffffffffffffffffffffffffff82169050919050565b6000819050919050565b82818337600083830152505050565b60005b83811015611a6f578082015181840152602081019050611a54565b83811115611a7e576000848401525b50505050565b60006002820490506001821680611a9c57607f821691505b60208210811415611ab057611aaf611b5f565b5b50919050565b611abf82611c2f565b810181811067ffffffffffffffff82111715611ade57611add611bec565b5b80604052505050565b6000611af282611a38565b91507fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff821415611b2557611b24611b30565b5b600182019050919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052602260045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603160045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b600080fd5b600080fd5b600080fd5b600080fd5b6000601f19601f8301169050919050565b7f4f6e6c79207374616b65722063616e2063616c6c2066756e6374696f6e000000600082015250565b7f56616c696461746f72207365742068617320726561636865642066756c6c206360008201527f6170616369747900000000000000000000000000000000000000000000000000602082015250565b7f696e646578206f7574206f662072616e67650000000000000000000000000000600082015250565b7f4f6e6c7920454f412063616e2063616c6c2066756e6374696f6e000000000000600082015250565b7f56616c696461746f72732063616e2774206265206c657373207468616e20746860008201527f65206d696e696d756d2072657175697265642076616c696461746f72206e756d602082015250565b611d62816119de565b8114611d6d57600080fd5b50565b611d7981611a38565b8114611d8457600080fd5b5056fea264697066735822122091251243e38d4ab79bc9bb795d9a55198a975905e2aca4f3efebad77e3cf0e7b64736f6c63430008070033"
)

// PredeployStakingSC is a helper method for setting up the staking smart contract account,
// using the passed in validators as pre-staked validators
func PredeployStakingSC(
	vals validators.Validators,
	params PredeployParams,
) (*chain.GenesisAccount, error) {
	// Set the code for the staking smart contract
	// Code retrieved from https://github.com/SECRYPT-2022/staking-contracts
	scHex, _ := hex.DecodeHex(StakingSCBytecode)
	stakingAccount := &chain.GenesisAccount{
		Code: scHex,
	}

	// Parse the default staked balance value into *big.Int
	val := DefaultStakedBalance
	bigDefaultStakedBalance, err := types.ParseUint256orHex(&val)

	if err != nil {
		return nil, fmt.Errorf("unable to generate DefaultStatkedBalance, %w", err)
	}

	// Generate the empty account storage map
	storageMap := make(map[types.Hash]types.Hash)
	bigTrueValue := big.NewInt(1)
	stakedAmount := big.NewInt(0)
	bigMinNumValidators := big.NewInt(int64(params.MinValidatorCount))
	bigMaxNumValidators := big.NewInt(int64(params.MaxValidatorCount))
	valsLen := big.NewInt(0)

	if vals != nil {
		valsLen = big.NewInt(int64(vals.Len()))

		for idx := 0; idx < vals.Len(); idx++ {
			validator := vals.At(uint64(idx))

			// Update the total staked amount
			stakedAmount = stakedAmount.Add(stakedAmount, bigDefaultStakedBalance)

			// Get the storage indexes
			storageIndexes := getStorageIndexes(validator, idx)

			// Set the value for the validators array
			storageMap[types.BytesToHash(storageIndexes.ValidatorsIndex)] =
				types.BytesToHash(
					validator.Addr().Bytes(),
				)

			if blsValidator, ok := validator.(*validators.BLSValidator); ok {
				setBytesToStorage(
					storageMap,
					storageIndexes.ValidatorBLSPublicKeyIndex,
					blsValidator.BLSPublicKey,
				)
			}

			// Set the value for the address -> validator array index mapping
			storageMap[types.BytesToHash(storageIndexes.AddressToIsValidatorIndex)] =
				types.BytesToHash(bigTrueValue.Bytes())

			// Set the value for the address -> staked amount mapping
			storageMap[types.BytesToHash(storageIndexes.AddressToStakedAmountIndex)] =
				types.StringToHash(hex.EncodeBig(bigDefaultStakedBalance))

			// Set the value for the address -> validator index mapping
			storageMap[types.BytesToHash(storageIndexes.AddressToValidatorIndexIndex)] =
				types.StringToHash(hex.EncodeUint64(uint64(idx)))
		}
	}

	// Set the value for the total staked amount
	storageMap[types.BytesToHash(big.NewInt(stakedAmountSlot).Bytes())] =
		types.BytesToHash(stakedAmount.Bytes())

	// Set the value for the size of the validators array
	storageMap[types.BytesToHash(big.NewInt(validatorsSlot).Bytes())] =
		types.BytesToHash(valsLen.Bytes())

	// Set the value for the minimum number of validators
	storageMap[types.BytesToHash(big.NewInt(minNumValidatorSlot).Bytes())] =
		types.BytesToHash(bigMinNumValidators.Bytes())

	// Set the value for the maximum number of validators
	storageMap[types.BytesToHash(big.NewInt(maxNumValidatorSlot).Bytes())] =
		types.BytesToHash(bigMaxNumValidators.Bytes())

	// Save the storage map
	stakingAccount.Storage = storageMap

	// Set the Staking SC balance to numValidators * defaultStakedBalance
	stakingAccount.Balance = stakedAmount

	return stakingAccount, nil
}
