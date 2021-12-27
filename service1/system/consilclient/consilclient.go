package consilclient

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/consul/api"
)

type ConsulConfig struct {
	ConsulHost string `envconfig:"CONSUL_HOST" required:"true"`
	ConsulPort string `envconfig:"CONSUL_PORT" required:"true"`
}

func (cc *ConsulConfig) NewConsulClient() *api.Client {
	port := strings.TrimSpace(cc.ConsulPort)
	port = strings.Trim(port, "\"")
	host := strings.TrimSpace(cc.ConsulHost)
	host = strings.Trim(host, "\"")
	string_con := newConsulStringConnection(host, port)
	client, err := api.NewClient(&api.Config{
		Address: string_con,
	})
	if err != nil {
		log.Fatal("err")
	}
	for {
		sta := client.Status()
		l, err := sta.Leader()
		fmt.Println(l)
		if err != nil {
			log.Printf("Consul reconect %s", string_con)
			time.Sleep(time.Second * 3)
			continue
		}
		break
	}
	return client
}

func newConsulStringConnection(host string, port string) string {
	return fmt.Sprintf("%s:%s", host, port)
}
