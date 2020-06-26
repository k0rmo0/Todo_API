package utils

import (
	"fmt"
	"path/filepath"
	"runtime"

	"github.com/sirupsen/logrus"
)

//Logger ...
var (
	Logger *logrus.Logger
	lvl    string
)

//LogError ...
func LogError(details string, err error) {
	if lvl != "prod" {
		fmt.Println(details, err)
	}

	_, filePath, line, _ := runtime.Caller(1)

	_, file := filepath.Split(filePath)

	Logger.WithField("file", file).WithField("line", line).Errorln(details, err.Error())
}
