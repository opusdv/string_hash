package rpcclient

import (
	"context"
	"fmt"
	"service2/conf"
	"service2/pkg/hashservice"
	"time"

	st "github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

func NewRpcClient(cfg *conf.Config, loger *logrus.Logger) hashservice.StringHashClient {
	cwt, _ := context.WithTimeout(context.Background(), time.Second*10)
	rpc_string_con := newRpcStringConnection(cfg.RPC_SERVVER, cfg.RPC_PORT)
	conn, err := grpc.DialContext(cwt, rpc_string_con, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		loger.WithFields(logrus.Fields{
			"Error":       err,
			"Stack Trace": st.WithStack(err),
		}).Errorf("Error connect rpc %s", rpc_string_con)
		panic(err)
	}
	uc := hashservice.NewStringHashClient(conn)
	loger.WithFields(logrus.Fields{
		"Message": fmt.Sprintf("%s connection sucsses", rpc_string_con),
	}).Debug("connection sucsses")
	return uc
}

func newRpcStringConnection(host string, port string) string {
	return fmt.Sprintf("%s:%s", host, port)
}
