package common

import (
	"log"
	"os"

	"github.com/hhkbp2/go-logging"
)

var logger logging.Logger

func InitLogger(fileName string) {
	if logger != nil {
		return
	}
	format := "%(asctime)s %(levelname)s (%(filename)s:%(lineno)d) " +
		"%(name)s %(message)s"
	dateFormat := "%Y-%m-%d %H:%M:%S.%3n"
	formatter := logging.NewStandardFormatter(format, dateFormat)

	handler, err := logging.NewFileHandler(fileName, os.O_APPEND, 4096)
	if err != nil {
		log.Fatal("[InitLogger]: " + err.Error())
	}
	handler.SetFormatter(formatter)
	logger = logging.GetLogger("[anyway2pocket]")
	logger.SetLevel(logging.LevelInfo)
	logger.AddHandler(handler)
}

func GetLogger() logging.Logger {
	if logger == nil {
		log.Fatal("[GetLogger]: call InitLogger first")
	}
	return logger
}
