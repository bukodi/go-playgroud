package playground

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/spf13/viper"
)

type cfgStruct struct {
	Field1 string `json:"field1"`
	Field2 int32  `json:"field2"`
}

var (
	cfgStr string
)

func TestViperCfg(t *testing.T) {
	x := viper.New()
	x.AddConfigPath(".")
	x.SetConfigName("testcfg")
	x.ReadInConfig()

	root := x.Get("root")

	root = cfgStruct{Field1: "Value1", Field2: 2222}

	x.Set("root", root)
	all := x.AllSettings()
	fmt.Println(all)
	str, _ := json.MarshalIndent(all, "", "  ")
	fmt.Println(string(str))
	x.WriteConfig()

}
