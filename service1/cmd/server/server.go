package main

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	"service1/conf"
	"service1/internal/handlers"
	"service1/pkg/hashservice"
	"service1/system/consilclient/initlog"

	st "github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

func main() {
	conf.InitDefaultConfig()
	cfg := conf.GetConfig()
	l := initlog.InitLog(cfg.LOG_LEVEL)
	l.WithFields(logrus.Fields{
		"Message": "Starting server",
	}).Info()
	errChan := make(chan error)

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	lis, err := net.Listen("tcp", ":8090")
	if err != nil {
		l.Panic("")
		panic(err)
	}
	s := grpc.NewServer()
	server := &handlers.Server{}
	server.Loger = l
	hashservice.RegisterStringHashServer(s, server)
	go func() {
		if err := s.Serve(lis); err != nil {
			errChan <- err
		}
	}()

	defer func() {
		s.GracefulStop()
	}()

	select {
	case err := <-errChan:
		l.WithFields(logrus.Fields{
			"Error": err,
		}).Error()
		log.Printf("Error : %v\n", err)
	case <-ctx.Done():
		l.WithFields(logrus.Fields{
			"Message":     "Server stop",
			"Stack Trace": st.WithStack(err),
		}).Info()
		cancel()
	}
}
