package config

import "fmt"

func Validate(cfg *Config) error {
	if len(cfg.Chains) == 0 {
		return fmt.Errorf("no chains configured")
	}

	for _, chain := range cfg.Chains {
		if chain.ChainID == 0 {
			return fmt.Errorf("chain %s has invalid chain_id", chain.Name)
		}

		switch chain.Type {
		case "ethereum":
			if chain.Confirmations <= 0 {
				return fmt.Errorf(
					"ethereum chain %s requires confirmations > 0",
					chain.Name,
				)
			}
		case "opstack":
			if chain.ReorgWindow <= 0 {
				return fmt.Errorf(
					"opstack chain %s requires reorg_window > 0",
					chain.Name,
				)
			}
		default:
			return fmt.Errorf(
				"chain %s has unknown type %s",
				chain.Name, chain.Type,
			)
		}

		if len(chain.Contracts) == 0 {
			return fmt.Errorf(
				"chain %s has no contracts configured",
				chain.Name,
			)
		}

		for _, c := range chain.Contracts {
			if c.Address == "" {
				return fmt.Errorf("chain %s has empty contract address", chain.Name)
			}
			if c.TokenDecimals != 18 {
				return fmt.Errorf(
					"contract %s on chain %s must have 18 decimals",
					c.Address, chain.Name,
				)
			}
		}
	}

	return nil
}
