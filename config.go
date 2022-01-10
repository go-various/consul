package consul

import (
	"encoding/json"
	"fmt"
	"github.com/hashicorp/hcl"
	"gopkg.in/yaml.v2"
	"strings"
)

//配置key格式规范: /config/{APP_NAME},{PROFILE}/{DATA_KEY}
//配置支持继承
// /config/application/{DATA_KEY}
// /config/application,dev/{DATA_KEY}
// /config/application,{PROFILE}/{DATA_KEY}
// /config/{APP_NAME}/{DATA_KEY}
// /config/{APP_NAME},{PROFILE}/{DATA_KEY}
func (c *client) LoadConfig(out interface{}) error {
	if c.config.Config.DataKey == "" {
		c.config.Config.DataKey = "1.0.0"
	}
	cfg := c.config.Config
	app := c.config.Application

	keyPaths := []string{ fmt.Sprintf("/config/application/%s", cfg.DataKey)}
	if app.Profile != ""{
		val := fmt.Sprintf("/config/application,%s/%s", app.Profile, cfg.DataKey)
		keyPaths = append(keyPaths, val)
	}

	val := fmt.Sprintf("/config/%s/%s", app.Name, cfg.DataKey)
	keyPaths = append(keyPaths, val)

	if app.Profile != ""{
		val = fmt.Sprintf("/config/%s,%s/%s", app.Name, app.Profile, cfg.DataKey)
		keyPaths = append(keyPaths, val)
	}

	options := c.queryOptions(nil)
	var succCount = 0
	for _, path := range keyPaths {
		kvp, _, err := c.client.KV().Get(path, options)
		if nil != err {
			return fmt.Errorf("load consul config: %s", err.Error())
		}
		if kvp == nil {
			continue
		}
		if err := c.decodeConfig(path, kvp.Value, out); err != nil {
			return err
		}
		succCount++
	}
	if 0 == succCount{
		return fmt.Errorf("load consul config %v is empty", keyPaths)
	}
	return nil
}

//default format using hcl
func (c *client)decodeConfig(key string, value []byte, out interface{})(err error)  {
	var format = strings.ToLower(c.config.Config.Format)
	switch format {
	case "json":
		err = json.Unmarshal(value, out)
	case "yaml":
		fallthrough
	case "yml":
		err = yaml.Unmarshal(value, out)
	default:
		err = hcl.Unmarshal(value, out)
	}
	if err != nil {
		return fmt.Errorf("unmarshal:%s => format:%s err:%s", key, format, err.Error())
	}
	return nil
}
