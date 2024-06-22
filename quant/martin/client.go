package martin

import (
	"github.com/iaping/go-okx/common"
	"github.com/iaping/go-okx/rest"
)

type OkxApiConf struct {
	ApiHost     string
	ApiKey      string
	SecretKey   string
	Passphrase  string
	IsSimulated bool
}

func NewOkxClient(c OkxApiConf) (client *rest.Client) {
	host := c.ApiHost
	apiKey := c.ApiKey
	secretKey := c.SecretKey
	passphrase := c.Passphrase
	isSimulated := c.IsSimulated

	var Auth = common.NewAuth(apiKey, secretKey, passphrase, isSimulated)

	client = rest.New(host, Auth, nil)
	return
}
