package sign

import (
	"github.com/artisanhe/tools/httplib"
	"github.com/gin-gonic/gin"
)

func getAccessKeyIdentifier() string {
	return AccessKey
}

type AccessKeyer interface {
	GetMerchantBy(string) (interface{}, error)
	GetSecret(interface{}) string
}

type AccessKeyParam struct {
	// AccessKey
	AccessKey string `json:"AccessKey" in:"header"`
}

func fetchMerchantAndStore(accessKeyer AccessKeyer, c *gin.Context) (interface{}, error) {
	req := AccessKeyParam{}
	err := httplib.TransToReq(c, &req)
	if err != nil {
		return nil, err
	}

	resp, err := accessKeyer.GetMerchantBy(req.AccessKey)
	if err != nil {
		return nil, err
	}

	c.Set(getAccessKeyIdentifier(), resp)
	return resp, nil
}

func FetchMerchant(accessKeyer AccessKeyer) gin.HandlerFunc {
	return func(c *gin.Context) {
		if _, err := fetchMerchantAndStore(accessKeyer, c); err != nil {
			httplib.WriteError(c, err)
			c.Abort()
			return
		}

		return
	}
}

func GetMerchantFromContext(c *gin.Context) (interface{}, bool) {
	return c.Get(getAccessKeyIdentifier())
}

func VerifySign(accessKeyer AccessKeyer) gin.HandlerFunc {
	return func(c *gin.Context) {
		var err error
		merchant, exist := GetMerchantFromContext(c)
		if !exist {
			merchant, err = fetchMerchantAndStore(accessKeyer, c)
			if err != nil {
				httplib.WriteError(c, err)
				c.Abort()
				return
			}
		}

		keyFunc := func(key string) (string, error) {
			return accessKeyer.GetSecret(merchant), nil
		}

		WithSignBy(keyFunc)(c)
	}
}
