package worker

type Config struct {
	AbiJSON         string `yaml:"abi_json"`
	BlockFrom       uint64 `yaml:"block_from"`
	ContractAddress string `yaml:"contract_address"`
	JSONRPCClient   string `yaml:"jsonrpc_client"`
}
