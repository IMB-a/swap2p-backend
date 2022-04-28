package worker

type Config struct {
	E20E20         string `yaml:"e20_e20"`
	E20E20Contract string `yaml:"e20_e20_contract"`

	E20E721         string `yaml:"e20_e721"`
	E20E721Contract string `yaml:"e20_e721_contract"`

	E721E20         string `yaml:"e721_e20"`
	E721E20Contract string `yaml:"e721_e20_contract"`

	E721E721         string `yaml:"e721_e721"`
	E721E721Contract string `yaml:"e721_e721_contract"`

	BlockFrom     uint64 `yaml:"block_from"`
	JSONRPCClient string `yaml:"jsonrpc_client"`
}
