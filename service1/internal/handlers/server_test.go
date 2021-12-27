package handlers

import (
	"context"
	"net"
	"service1/pkg/hashservice"
	"testing"

	formatter "github.com/fabienm/go-logrus-formatters"
	graylog "github.com/gemnasium/logrus-graylog-hook"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
)

const bufSize = 1024 * 1024

var l *logrus.Logger
var lis *bufconn.Listener

func init() {
	gelfFmt := formatter.NewGelf("service1_testing")
	l = logrus.New()
	hook := graylog.NewGraylogHook("localhost:12201", map[string]interface{}{})
	l.SetLevel(logrus.InfoLevel)
	l.SetFormatter(gelfFmt)
	l.AddHook(hook)
	lis = bufconn.Listen(bufSize)
	s := grpc.NewServer()
	server := Server{}
	server.Loger = l
	hashservice.RegisterStringHashServer(s, &server)
	go func() {
		if err := s.Serve(lis); err != nil {
			l.WithFields(logrus.Fields{
				"Error": err,
			})
		}
	}()
}
func bufDialer(context.Context, string) (net.Conn, error) {
	return lis.Dial()
}

func TestCreateHash(t *testing.T) {
	req := require.New(t)
	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithInsecure())
	if err != nil {
		l.WithFields(logrus.Fields{
			"Error": err,
		})
	}
	defer conn.Close()
	client := hashservice.NewStringHashClient(conn)
	s := &hashservice.Strings{
		List: []string{"teststring", "teststring2"},
	}
	want_resp, err := client.CreateHash(ctx, s)
	tests := map[string]struct {
		firstVal  context.Context
		secondVal *hashservice.Strings
		want      *hashservice.Hashs
		err       error
	}{
		"simple": {
			firstVal:  ctx,
			secondVal: s,
			want:      want_resp,
			err:       nil,
		},
	}

	for name, testCase := range tests {
		t.Run(name, func(t *testing.T) {
			resp, err := client.CreateHash(testCase.firstVal, testCase.secondVal)
			req.Equal(testCase.err, err)
			req.Equal(testCase.want.List, resp.List)
		})
	}
}
