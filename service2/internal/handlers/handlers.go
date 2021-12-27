package handlers

import (
	"context"
	"math/rand"
	"service2/conf"
	"service2/internal/handler/operations"
	"service2/models"
	"service2/pkg/hashservice"
	"service2/system/database/mongoclient"
	"service2/system/rpcclient"
	"strconv"
	"sync"
	"time"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"github.com/go-openapi/runtime/middleware"
	st "github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/metadata"
)

type HashStruct struct {
	rpcClient hashservice.StringHashClient
	dbClient  *mgo.Session
	Wg        *sync.WaitGroup
	loger     *logrus.Logger
}

func NewHashStruct(cfg *conf.Config, loger *logrus.Logger) *HashStruct {
	r := rpcclient.NewRpcClient(cfg, loger)
	d := mongoclient.NewDbConnection(cfg, loger)
	return &HashStruct{
		rpcClient: r,
		dbClient:  d,
		Wg:        cfg.Wg,
		loger:     loger,
	}
}

func (hs *HashStruct) UpdateConfig() error {
	cfg := conf.GetConfig()
	r := rpcclient.NewRpcClient(cfg, hs.loger)
	d := mongoclient.NewDbConnection(cfg, hs.loger)
	hs.rpcClient = r
	hs.dbClient = d

	return nil
}

func (hs *HashStruct) SendHashHandler(psp operations.PostSendParams) middleware.Responder {
	hs.Wg.Add(1)
	defer hs.Wg.Done()
	md := metadata.New(map[string]string{"requestid": psp.HTTPRequest.Context().Value("requestID").(string)})
	ctx := metadata.NewOutgoingContext(context.Background(), md)
	s := []string{}

	for _, v := range psp.Params {
		s = append(s, v)
	}

	us, err := hs.rpcClient.CreateHash(ctx, &hashservice.Strings{List: s})
	if err != nil {
		hs.loger.WithFields(
			logrus.Fields{
				"Error":       err,
				"Stack trace": st.WithStack(err),
				"Request ID":  psp.HTTPRequest.Context().Value("requestID").(string),
			}).Error("Error rpc create hash ")
		return operations.NewPostSendInternalServerError()
	}
	ah := models.ArrayOfHash{}

	for _, v := range us.List {
		r := randId(1000, 500000)
		ah = append(ah, &models.Hash{Hash: &v.H, ID: r})
		err = hs.dbClient.DB("string_hash").C("hash").Insert(models.Hash{Hash: &v.H, ID: r})
		if err != nil {
			hs.loger.WithFields(
				logrus.Fields{
					"Error":       err,
					"Stack trace": st.WithStack(err),
					"Request ID":  psp.HTTPRequest.Context().Value("requestID").(string),
				}).Error("Errir db insert")
			return operations.NewPostSendInternalServerError()
		}
	}

	hs.loger.WithFields(
		logrus.Fields{
			"Message":    ah,
			"Request ID": psp.HTTPRequest.Context().Value("requestID").(string),
		}).Info("Info hash")
	return operations.NewPostSendOK().WithPayload(ah)
}

func (hs *HashStruct) CheckHashHandler(gcp operations.GetCheckParams) middleware.Responder {
	hs.Wg.Add(1)
	defer hs.Wg.Done()
	var h models.Hash
	ah := models.ArrayOfHash{}

	for _, v := range gcp.Ids {
		id, err := strconv.Atoi(v)
		if err != nil {
			hs.loger.WithFields(
				logrus.Fields{
					"Error":       err,
					"Stack trace": st.WithStack(err),
					"Request ID":  gcp.HTTPRequest.Context().Value("requestID").(string),
				}).Error("Error ids convertation")
			return operations.NewGetCheckBadRequest()
		}
		q := bson.M{
			"id": id,
		}
		err = hs.dbClient.DB("string_hash").C("hash").Find(q).One(&h)
		if err != nil {
			hs.loger.WithFields(
				logrus.Fields{
					"Error":       err,
					"Stack trace": st.WithStack(err),
					"Request ID":  gcp.HTTPRequest.Context().Value("requestID").(string),
				}).Error("Error db find")
		}
		ah = append(ah, &h)
	}
	hs.loger.WithFields(logrus.Fields{
		"Message":    ah,
		"Request ID": gcp.HTTPRequest.Context().Value("requestID").(string),
	}).Info()
	return operations.NewGetCheckOK().WithPayload(ah)
}

func randId(min int64, max int64) *int64 {
	rand.Seed(time.Now().UnixNano())
	r := rand.Int63n(max-min+1) + min
	return &r
}
