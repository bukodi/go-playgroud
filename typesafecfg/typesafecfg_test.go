package typesafecfg

import (
	"encoding/json"
	"fmt"
	"testing"

	"gopkg.in/yaml.v2"
)

type CfgRoot struct {
	API   CfgListener `json:"api"`
	Admin CfgListener `json:"admin"`
}

type CfgListener struct {
	Address string `json:"address"`
	Port    int16  `json:"port"`
}

func TestJson(t *testing.T) {
	cfg := CfgRoot{
		API: CfgListener{
			Address: "localhost",
			Port:    8080,
		},
		Admin: CfgListener{
			Address: "localhost",
			Port:    8088,
		},
	}

	fmt.Printf("%#v\n", cfg)
	cfgTxt, _ := json.MarshalIndent(cfg, "", "  ")
	fmt.Printf("%s\n", cfgTxt)

	var cfg2 CfgRoot
	json.Unmarshal(cfgTxt, &cfg2)
	fmt.Printf("%#v\n", cfg2)
}

func TestYaml(t *testing.T) {
	cfg := CfgRoot{
		API: CfgListener{
			Address: "localhost",
			Port:    8080,
		},
		Admin: CfgListener{
			Address: "localhost",
			Port:    8088,
		},
	}
aa
	fmt.Printf("%#v\n", cfg)
	cfgTxt, _ := yaml.Marshal(cfg)
	fmt.Printf("%s\n", cfgTxt)

	var cfg2 CfgRoot
	yaml.UnmarshalStrict(cfgTxt, &cfg2)
	fmt.Printf("%#v\n", cfg2)
}
