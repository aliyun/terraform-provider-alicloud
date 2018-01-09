package alicloud

import (
	"fmt"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/denverdino/aliyungo/cdn"
	"github.com/denverdino/aliyungo/common"
	"github.com/denverdino/aliyungo/dns"
	"github.com/denverdino/aliyungo/ecs"
	"github.com/denverdino/aliyungo/ess"
	"github.com/denverdino/aliyungo/location"
	"github.com/denverdino/aliyungo/ram"
	"github.com/denverdino/aliyungo/rds"
	"github.com/denverdino/aliyungo/slb"

	"github.com/denverdino/aliyungo/cs"
	"github.com/hashicorp/terraform/terraform"
	"log"
	"strings"
)

// Config of aliyun
type Config struct {
	AccessKey     string
	SecretKey     string
	Region        common.Region
	SecurityToken string
}

// AliyunClient of aliyun
type AliyunClient struct {
	Region  common.Region
	ecsconn *ecs.Client
	essconn *ess.Client
	rdsconn *rds.Client
	// use new version
	ecsNewconn *ecs.Client
	vpcconn    *ecs.Client
	slbconn    *slb.Client
	ossconn    *oss.Client
	dnsconn    *dns.Client
	ramconn    ram.RamClientInterface
	csconn     *cs.Client
	cdnconn    *cdn.CdnClient
}

// Client for AliyunClient
func (c *Config) Client() (*AliyunClient, error) {
	err := c.loadAndValidate()
	if err != nil {
		return nil, err
	}

	ecsconn, err := c.ecsConn()
	if err != nil {
		return nil, err
	}

	ecsNewconn, err := c.ecsConn()
	if err != nil {
		return nil, err
	}
	ecsNewconn.SetVersion(EcsApiVersion20160314)

	rdsconn, err := c.rdsConn()
	if err != nil {
		return nil, err
	}

	slbconn, err := c.slbConn()
	if err != nil {
		return nil, err
	}

	vpcconn, err := c.vpcConn()
	if err != nil {
		return nil, err
	}

	essconn, err := c.essConn()
	if err != nil {
		return nil, err
	}
	ossconn, err := c.ossConn()
	if err != nil {
		return nil, err
	}
	dnsconn, err := c.dnsConn()
	if err != nil {
		return nil, err
	}
	ramconn, err := c.ramConn()
	if err != nil {
		return nil, err
	}
	csconn, err := c.csConn()
	if err != nil {
		return nil, err
	}
	cdnconn, err := c.cdnConn()
	if err != nil {
		return nil, err
	}

	return &AliyunClient{
		Region:     c.Region,
		ecsconn:    ecsconn,
		ecsNewconn: ecsNewconn,
		vpcconn:    vpcconn,
		slbconn:    slbconn,
		rdsconn:    rdsconn,
		essconn:    essconn,
		ossconn:    ossconn,
		dnsconn:    dnsconn,
		ramconn:    ramconn,
		csconn:     csconn,
		cdnconn:    cdnconn,
	}, nil
}

const BusinessInfoKey = "Terraform"

func (c *Config) loadAndValidate() error {
	err := c.validateRegion()
	if err != nil {
		return err
	}

	return nil
}

func (c *Config) validateRegion() error {

	for _, valid := range common.ValidRegions {
		if c.Region == valid {
			return nil
		}
	}

	return fmt.Errorf("Not a valid region: %s", c.Region)
}

func (c *Config) ecsConn() (*ecs.Client, error) {
	client := ecs.NewECSClientWithSecurityToken(c.AccessKey, c.SecretKey, c.SecurityToken, c.Region)
	client.SetBusinessInfo(BusinessInfoKey)
	client.SetUserAgent(getUserAgent())

	if _, err := client.DescribeRegions(); err != nil {
		return nil, err
	}

	return client, nil
}

func (c *Config) rdsConn() (*rds.Client, error) {
	client := rds.NewRDSClient(c.AccessKey, c.SecretKey, c.Region)
	client.SetBusinessInfo(BusinessInfoKey)
	client.SetUserAgent(getUserAgent())
	return client, nil
}

func (c *Config) slbConn() (*slb.Client, error) {
	client := slb.NewSLBClient(c.AccessKey, c.SecretKey, c.Region)
	client.SetBusinessInfo(BusinessInfoKey)
	client.SetUserAgent(getUserAgent())
	return client, nil
}

func (c *Config) vpcConn() (*ecs.Client, error) {
	client := ecs.NewVPCClientWithSecurityToken(c.AccessKey, c.SecretKey, c.SecurityToken, c.Region)
	client.SetBusinessInfo(BusinessInfoKey)
	client.SetUserAgent(getUserAgent())
	return client, nil

}
func (c *Config) essConn() (*ess.Client, error) {
	client := ess.NewESSClient(c.AccessKey, c.SecretKey, c.Region)
	client.SetBusinessInfo(BusinessInfoKey)
	client.SetUserAgent(getUserAgent())
	return client, nil
}
func (c *Config) ossConn() (*oss.Client, error) {

	endpointClient := location.NewClient(c.AccessKey, c.SecretKey)
	endpointClient.SetSecurityToken(c.SecurityToken)
	args := &location.DescribeEndpointsArgs{
		Id:          c.Region,
		ServiceCode: "oss",
		Type:        "openAPI",
	}

	endpoints, err := endpointClient.DescribeEndpoints(args)
	if err != nil {
		return nil, fmt.Errorf("Describe endpoint using region: %#v got an error: %#v.", c.Region, err)
	}
	endpointItem := endpoints.Endpoints.Endpoint
	var endpoint string
	if endpointItem == nil || len(endpointItem) <= 0 {
		log.Printf("Cannot find endpoint in the region: %#v", c.Region)
		endpoint = ""
	} else {
		endpoint = strings.ToLower(endpointItem[0].Protocols.Protocols[0]) + "://" + endpointItem[0].Endpoint
	}

	if endpoint == "" {
		endpoint = fmt.Sprintf("http://oss-%s.aliyuncs.com", c.Region)
	}

	log.Printf("[DEBUG] Instantiate OSS client using endpoint: %#v", endpoint)
	client, err := oss.New(endpoint, c.AccessKey, c.SecretKey, oss.UserAgent(getUserAgent()))

	return client, err
}

func (c *Config) dnsConn() (*dns.Client, error) {
	client := dns.NewClientNew(c.AccessKey, c.SecretKey)
	client.SetBusinessInfo(BusinessInfoKey)
	client.SetUserAgent(getUserAgent())
	return client, nil
}

func (c *Config) ramConn() (ram.RamClientInterface, error) {
	client := ram.NewClient(c.AccessKey, c.SecretKey)
	return client, nil
}

func (c *Config) csConn() (*cs.Client, error) {
	client := cs.NewClientForAussumeRole(c.AccessKey, c.SecretKey, c.SecurityToken)
	client.SetUserAgent(getUserAgent())
	return client, nil
}

func (c *Config) cdnConn() (*cdn.CdnClient, error) {
	client := cdn.NewClient(c.AccessKey, c.SecretKey)
	client.SetBusinessInfo(BusinessInfoKey)
	client.SetUserAgent(getUserAgent())
	return client, nil
}

func getUserAgent() string {
	return fmt.Sprintf("HashiCorp-Terraform-v%s", terraform.VersionString())
}
