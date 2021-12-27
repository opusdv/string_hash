package mongoclient

import (
	"fmt"
	"service2/conf"
	"time"

	"github.com/globalsign/mgo"
	st "github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

func NewDbConnection(cfg *conf.Config, loger *logrus.Logger) *mgo.Session {
	db_string_con := newDbStringConnection(cfg.DB_HOST, cfg.DB_PORT, cfg.DB_NAME)
	var sess *mgo.Session
	sess, err := mgo.DialWithTimeout(db_string_con, time.Second*10)
	if err != nil {
		p := sess.Ping()
		fmt.Println(p)
		loger.WithFields(logrus.Fields{
			"Error":       err,
			"Stack Trace": st.WithStack(err),
		}).Errorf("Error connect db %s", db_string_con)
	}
	loger.WithFields(logrus.Fields{
		"Message": fmt.Sprintf("%s connection sucsses", db_string_con),
	}).Debug("connection sucsses")

	return sess
}

func newDbStringConnection(host string, port string, db_name string) string {
	db_string_con := fmt.Sprintf("mongodb://%s:%s/%s", host, port, db_name)
	return db_string_con
}
