package connectivity

import (
	"fmt"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/auth"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/auth/credentials"
)

// Config of aliyun
type Config struct {
	AccessKey       string
	SecretKey       string
	Region          Region
	RegionId        string
	SecurityToken   string
	OtsInstanceName string
	AccountId       string

	EcsEndpoint        string
	RdsEndpoint        string
	SlbEndpoint        string
	VpcEndpoint        string
	CenEndpoint        string
	EssEndpoint        string
	OssEndpoint        string
	DnsEndpoint        string
	RamEndpoint        string
	CsEndpoint         string
	CdnEndpoint        string
	KmsEndpoint        string
	OtsEndpoint        string
	CmsEndpoint        string
	PvtzEndpoint       string
	StsEndpoint        string
	LogEndpoint        string
	DrdsEndpoint       string
	DdsEndpoint        string
	KVStoreEndpoint    string
	FcEndpoint         string
	ApigatewayEndpoint string
	DatahubEndpoint    string
	MnsEndpoint        string
	LocationEndpoint   string
	ElasticsearchEndpoint string
}

func (c *Config) loadAndValidate() error {
	err := c.validateRegion()
	if err != nil {
		return err
	}

	return nil
}

func (c *Config) validateRegion() error {

	for _, valid := range ValidRegions {
		if c.Region == valid {
			return nil
		}
	}

	return fmt.Errorf("Not a valid region: %s", c.Region)
}

func (c *Config) getAuthCredential(stsSupported bool) auth.Credential {
	if stsSupported {
		return credentials.NewStsTokenCredential(c.AccessKey, c.SecretKey, c.SecurityToken)
	}

	return credentials.NewAccessKeyCredential(c.AccessKey, c.SecretKey)
}
