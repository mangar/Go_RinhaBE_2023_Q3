package helpers

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

func SetupLog() os.File {

	flagDebug, _ := strconv.ParseBool(os.Getenv("FLAG_DEBUG"))
	if flagDebug {
		logrus.SetLevel(logrus.DebugLevel)
	} else {
		logrus.SetLevel(logrus.InfoLevel)
	}

	agora := time.Now()
	dataFormatada := agora.Format("20060102")
	fileName := os.Getenv("LOG_OUTPUT_DIR") + "/" + dataFormatada + "_" + os.Getenv("SERVER_NAME") + ".log"
	logFile, err := os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		logrus.Fatalf("Erro ao abrir arquivo de log: %v", err)
	}

	mw := io.MultiWriter(os.Stdout, logFile)
	logrus.SetOutput(mw)

    logrus.SetFormatter(new(CustomFormatter)) // Configura o logger para usar o formatter personalizado

	return *logFile

}

type CustomFormatter struct{}

func (f *CustomFormatter) Format(entry *logrus.Entry) ([]byte, error) {
    timestamp := entry.Time.Format("2006-01-02 15:04:05.000000")
    logMessage := fmt.Sprintf("[%s] [%s] %s\n", timestamp, strings.ToUpper(entry.Level.String()), entry.Message)
    return []byte(logMessage), nil
}

func LogOnError(err error, msg string) error {
	if err != nil {
		logrus.Error(msg + " .. " + err.Error())
		return errors.New(msg + " .. " + err.Error())
	} else {
		return nil
	}
}

func ExitOnError(err error, msg string) {
	if err != nil {
		logrus.Error(msg + " .. " + err.Error())
		os.Exit(1)
	}
}
