package reqid

import (
	"fmt"

	guid "github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
)

func GenerateReqId(loger *logrus.Logger) string {
	uuid := guid.Must(guid.NewV4(), nil)
	reqId := uuid.String()
	loger.WithFields(logrus.Fields{
		"Message": fmt.Sprintf("Generate reqid : %s", reqId),
	}).Debug()
	return reqId
}
