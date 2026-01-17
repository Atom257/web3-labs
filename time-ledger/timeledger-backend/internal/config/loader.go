package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/pelletier/go-toml/v2"
)

func Load(configPath string) (*Config, error) {
	//加载 .env（允许不存在）
	_ = godotenv.Load()

	//读取 toml
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("read config file failed: %w", err)
	}

	var cfg Config
	if err := toml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("parse toml failed: %w", err)
	}

	//注入 RPC URL
	for i := range cfg.Chains {
		chain := &cfg.Chains[i]

		if chain.RPCEnvKey == "" {
			return nil, fmt.Errorf("chain %s missing rpc_env_key", chain.Name)
		}

		rpcURL := os.Getenv(chain.RPCEnvKey)
		if rpcURL == "" {
			return nil, fmt.Errorf(
				"env %s not set for chain %s",
				chain.RPCEnvKey, chain.Name,
			)
		}

		chain.RPCURL = rpcURL
	}

	//启动期校验
	if err := Validate(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
