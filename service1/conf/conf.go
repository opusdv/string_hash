package conf

import (
	"encoding/json"
	"log"
	"service1/system/consilclient"
	"sync"

	"github.com/hashicorp/consul/api"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	LOG_LEVEL string `json:"log_level" envconfig:"LOG_LEVEL" required:"true"`
	Wg        *sync.WaitGroup
}

func InitDefaultConfig() {
	var c Config
	var cfgConsul consilclient.ConsulConfig

	err := envconfig.Process("stgring_hash", &c)
	if err != nil {
		panic(err)
	}
	err = envconfig.Process("string_hash", &cfgConsul)
	if err != nil {
		panic(err)
	}

	j, err := json.MarshalIndent(c, " ", "\t")
	if err != nil {
		log.Fatal(err)
	}

	client := cfgConsul.NewConsulClient()
	kv := client.KV()
	p := &api.KVPair{
		Key:   "APP_CONFIG",
		Value: []byte(j),
	}
	_, err = kv.Put(p, nil)
	if err != nil {
		panic(err)
	}
}

func GetConfig() *Config {
	var cfg Config
	var cfgConsul consilclient.ConsulConfig
	err := envconfig.Process("string_hash", &cfgConsul)
	if err != nil {
		panic(err)
	}
	client := cfgConsul.NewConsulClient()
	kv := client.KV()
	pair, _, err := kv.Get("APP_CONFIG", nil)
	if err != nil {
		panic(err)
	}

	if err := json.Unmarshal([]byte(pair.Value), &cfg); err != nil {
		panic(err)
	}

	return &cfg
}
