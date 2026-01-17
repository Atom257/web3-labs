package config

type Config struct {
	App      AppConfig      `toml:"app"`
	Database DatabaseConfig `toml:"database"`
	Redis    RedisConfig    `toml:"redis"`
	Chains   []ChainConfig  `toml:"chains"`
}

type AppConfig struct {
	Name     string `toml:"name"`
	Timezone string `toml:"timezone"`
}

type DatabaseConfig struct {
	MaxOpenConns int `toml:"max_open_conns"`
	MaxIdleConns int `toml:"max_idle_conns"`
}

type RedisConfig struct {
	KeyPrefix string `toml:"key_prefix"`
}

type ChainConfig struct {
	Name           string `toml:"name"`
	ChainID        int64  `toml:"chain_id"`
	Type           string `toml:"type"` // ethereum | opstack
	RPCEnvKey      string `toml:"rpc_env_key"`
	Confirmations  int64  `toml:"confirmations"`
	ReorgWindow    int64  `toml:"reorg_window"`
	ChunkSize      uint64 `toml:"chunk_size"`       // 每次同步的区块数量，默认 10
	RequestDelayMs int64  `toml:"request_delay_ms"` // 每次请求之间的延迟（毫秒），默认 100

	Contracts []ContractConfig `toml:"contracts"`

	// 派生字段（不来自 toml）
	RPCURL string `toml:"-"`
}

type ContractConfig struct {
	Address       string `toml:"address"`
	StartBlock    int64  `toml:"start_block"`
	TokenDecimals int64  `toml:"token_decimals"`
}
