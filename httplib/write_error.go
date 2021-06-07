package httplib

import (
    "net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"github.com/artisanhe/tools/courier/status_error"
)

var serviceName = ""

func SetServiceName(name string) {
	serviceName = name
}

func getServiceName() string {
	if serviceName == "" {
		SetServiceName(os.Getenv("PROJECT_NAME"))
	}
	return serviceName
}

func errorWithSource(err *status_error.StatusError) *status_error.StatusError {
	return err.WithSource(getServiceName())
}

func WriteError(c *gin.Context, err error) {
	wrapErrorForEmqx(c, err)
	//statusError := status_error.FromError(err)
	//if statusError.Code == int64(status_error.UnknownError) {
	//	logrus.Warnf("got UnknownError %s", err.Error())
	//}
	//statusError = errorWithSource(statusError)
	//c.Error(statusError)
	//c.JSON(statusError.Status(), statusError)
}

type EmqxError struct {
	Code    int32  `json:"code"`
	Message string `json:"message"`
	Data struct{} `json:"data"`
}

func wrapErrorForEmqx(c *gin.Context, err error) {
	statusError := status_error.FromError(err)
	if statusError.Code == int64(status_error.UnknownError) {
		logrus.Warnf("got UnknownError %s", err.Error())
	}

	errResp := EmqxError{
		Code : int32(statusError.Code),
		Message: statusError.Desc + statusError.Msg + statusError.ErrorFields.String(),
	}

	c.JSON(http.StatusOK, errResp)
}
