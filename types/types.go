package types

const (
	OneCoin          = 1000000 //1BCC = 1000000ubcc
	Coin             = "BCC"
	DefaultDenom     = "ubcc"
	DaysPerYear      = 360
	DaysPerQuarter   = 90
	BlockTime        = 10 //seconds
	BlocksPerHour    = 3600 / BlockTime
	BlocksPerDay     = (3600 / BlockTime) * 24
	BlocksPerQuarter = BlocksPerDay * DaysPerQuarter
	BlocksPerYear    = BlocksPerDay * DaysPerYear

	DefaultMintStartValue  = 4 * OneCoin         //minting 4 BCC per block
	DefaultMaxEmission     = 100000000 * OneCoin //100 000 000 BCC
	DefaultMintingInterval = 100                 //Add coins to fee pool every 100 block to decrease blockchain fast growing

	DefaultMaxGas      = -1
	DefaultMinGasPrice = "0.25ubcc"

	DefaultMinValidatorSelfDelegation = 4000000 * OneCoin //minimum tokens to be a validator
	DefaultValidatorDelegateEnabled   = false

	// AddrLen defines a valid address length
	AddrLen = 20
	// Bech32PrefixAccAddr defines the Bech32 prefix of an account's address
	Bech32MainPrefix = "boco"

	// Boco in https://github.com/satoshilabs/slips/blob/master/slip-0044.md
	CoinType = 3020 //0x80000BCC

	// BIP44Prefix is the parts of the BIP44 HD path that are fixed by
	// what we used during the fundraiser.
	FullFundraiserPath = "44'/3020'/0'/0/0"

	// PrefixAccount is the prefix for account keys
	PrefixAccount = "acc"
	// PrefixValidator is the prefix for validator keys
	PrefixValidator = "val"
	// PrefixConsensus is the prefix for consensus keys
	PrefixConsensus = "cons"
	// PrefixPublic is the prefix for public keys
	PrefixPublic = "pub"
	// PrefixOperator is the prefix for operator keys
	PrefixOperator = "oper"

	// PrefixAddress is the prefix for addresses
	PrefixAddress = "addr"

	// Bech32PrefixAccAddr defines the Bech32 prefix of an account's address
	Bech32PrefixAccAddr = Bech32MainPrefix
	// Bech32PrefixAccPub defines the Bech32 prefix of an account's public key
	Bech32PrefixAccPub = Bech32MainPrefix + PrefixPublic
	// Bech32PrefixValAddr defines the Bech32 prefix of a validator's operator address
	Bech32PrefixValAddr = Bech32MainPrefix + PrefixValidator + PrefixOperator
	// Bech32PrefixValPub defines the Bech32 prefix of a validator's operator public key
	Bech32PrefixValPub = Bech32MainPrefix + PrefixValidator + PrefixOperator + PrefixPublic
	// Bech32PrefixConsAddr defines the Bech32 prefix of a consensus node address
	Bech32PrefixConsAddr = Bech32MainPrefix + PrefixValidator + PrefixConsensus
	// Bech32PrefixConsPub defines the Bech32 prefix of a consensus node public key
	Bech32PrefixConsPub = Bech32MainPrefix + PrefixValidator + PrefixConsensus + PrefixPublic
)
