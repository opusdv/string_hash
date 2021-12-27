package handlers

import (
	"context"
	"fmt"
	"service1/pkg/hashservice"
	"service1/pkg/myhash"
	"sync"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/metadata"
)

var wg sync.WaitGroup

type Server struct {
	Hashs []*hashservice.Hash
	hashservice.UnimplementedStringHashServer
	Loger *logrus.Logger
}

func (s *Server) CreateHash(ctx context.Context, strings *hashservice.Strings) (*hashservice.Hashs, error) {
	hs := &hashservice.Hashs{}
	header, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		s.Loger.WithFields(logrus.Fields{
			"Message": "Reques id not found",
		}).Error()
	}
	hash := []*hashservice.Hash{}
	for _, ss := range strings.List {
		wg.Add(1)
		go func(ss string) {
			defer wg.Done()
			shash := myhash.Hash(ss)
			hash = append(hash, &hashservice.Hash{H: shash})
		}(ss)
	}
	wg.Wait()
	hs.List = hash
	hs.Success = true
	s.Loger.WithFields(logrus.Fields{
		"Message":    hs,
		"Request ID": fmt.Sprintf("Request ID %s", header["requestid"]),
	}).Info()
	return hs, nil
}
