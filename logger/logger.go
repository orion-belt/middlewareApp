package logger

import (
	formatter "github.com/antonfisher/nested-logrus-formatter"
	"github.com/sirupsen/logrus"
	"time"
)

var (
	log *logrus.Logger

	AppLog        *logrus.Entry
	ConfigLog     *logrus.Entry
	GrpcLog       *logrus.Entry
	MagmaGrpcLog  *logrus.Entry
	MagmaGwRegLog *logrus.Entry
	OaiSbiLog     *logrus.Entry
	ApiConv     *logrus.Entry
	HttpLog       *logrus.Entry
)

func init() {
	log = logrus.New()
	log.SetReportCaller(false)
	log.Formatter = &formatter.Formatter{
		TimestampFormat: time.RFC3339,
		TrimMessages:    true,
		NoFieldsSpace:   true,
		HideKeys:        true,
	}

	AppLog = log.WithFields(logrus.Fields{"category": "APPS"})
	ConfigLog = log.WithFields(logrus.Fields{"category": "CONFIG"})
	MagmaGrpcLog = log.WithFields(logrus.Fields{"category": "MAGM", "component": "GRPC"})
	OaiSbiLog = log.WithFields(logrus.Fields{"category": "OAI", "component": "SBI"})
	ApiConv = log.WithFields(logrus.Fields{"category": "APIC", "component": "MO"})
	MagmaGwRegLog = log.WithFields(logrus.Fields{"category": "MAGM", "component": "HTTP"})
}

func SetLogLevel(level string) {
	lvl, err := logrus.ParseLevel(level)
	if err != nil {
		AppLog.Fatalln("Failed to parse log level:", err)
	}
	log.SetLevel(lvl)
}

func SetReportCaller(set bool) {
	log.SetReportCaller(set)
}
