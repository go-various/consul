package consul

import (
	"testing"
	"github.com/hashicorp/go-hclog"

)

type exampleConfig struct {
	App interface{} `hcl:"mysql" yaml:"mysql"`
}

func TestClient_LoadConfig(t *testing.T) {
	cli, err := NewClient(mockConfig(),hclog.Default())
	if nil != err {
		t.Fatal(err)
	}

	//载入配置
	var cfg exampleConfig
	if err := cli.LoadConfig(&cfg); err != nil {
		t.Fatal(err)
	}
	t.Log(cfg)
}
