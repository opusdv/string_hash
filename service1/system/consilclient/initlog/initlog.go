package initlog

import (
	"os"

	formatter "github.com/fabienm/go-logrus-formatters"
	"github.com/sirupsen/logrus"
)

func InitLog(loglevel string) *logrus.Logger {
	gelfFmt := formatter.NewGelf("service1")
	l := logrus.New()
	ll, err := logrus.ParseLevel(loglevel)
	if err != nil {
		l.SetLevel(logrus.InfoLevel)
	}
	l.SetLevel(ll)
	l.SetFormatter(gelfFmt)
	l.SetOutput(os.Stdout)
	return l
}
