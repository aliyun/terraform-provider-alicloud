package config

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/errors"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/location"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
)

type Endpoints struct {
	Endpoint []Endpoint `xml:"Endpoint"`
}

type Endpoint struct {
	Name      string    `xml:"name,attr"`
	RegionIds RegionIds `xml:"RegionIds"`
	Products  Products  `xml:"Products"`
}

type RegionIds struct {
	RegionId string `xml:"RegionId"`
}

type Products struct {
	Product []Product `xml:"Product"`
}

type Product struct {
	ProductName string `xml:"ProductName"`
	DomainName  string `xml:"DomainName"`
}

var localEndpointPath = "./endpoints.xml"
var localEndpointPathEnv = "TF_ENDPOINT_PATH"
var loadLocalEndpoint = false

func hasLocalEndpoint() bool {
	data, err := ioutil.ReadFile(localEndpointPath)
	if err != nil || len(data) <= 0 {
		d, e := ioutil.ReadFile(os.Getenv(localEndpointPathEnv))
		if e != nil {
			return false
		}
		data = d
	}
	return len(data) > 0
}

func LoadRegionalEndpoint(region string, serviceCode string) string {
	if region == "" || serviceCode == "" {
		return ""
	}
	return fmt.Sprintf("%s.%s.aliyuncs.com", serviceCode, region)
}

func LoadEndpoint(region string, serviceCode ServiceCode) string {
	endpoint := strings.TrimSpace(os.Getenv(fmt.Sprintf("%s_ENDPOINT", string(serviceCode))))
	if endpoint != "" {
		return endpoint
	}

	// Load current path endpoint file endpoints.xml, if failed, it will load from environment variables TF_ENDPOINT_PATH
	if !loadLocalEndpoint {
		return ""
	}
	data, err := ioutil.ReadFile(localEndpointPath)
	if err != nil || len(data) <= 0 {
		d, e := ioutil.ReadFile(os.Getenv(localEndpointPathEnv))
		if e != nil {
			return ""
		}
		data = d
	}
	var endpoints Endpoints
	err = xml.Unmarshal(data, &endpoints)
	if err != nil {
		return ""
	}
	for _, endpoint := range endpoints.Endpoint {
		if endpoint.RegionIds.RegionId == string(region) {
			for _, product := range endpoint.Products.Product {
				if strings.ToLower(product.ProductName) == strings.ToLower(string(serviceCode)) {
					return strings.TrimSpace(product.DomainName)
				}
			}
		}
	}

	return ""
}

// NOTE: The productCode must be lowed.
func (c *Config) LoadEndpoint(productCode string, regionID string) error {
	// Firstly, load endpoint from environment variables
	endpoint := strings.TrimSpace(os.Getenv(fmt.Sprintf("ALIBABA_CLOUD_ENDPOINT_%s", strings.ToUpper(productCode))))
	if endpoint == "" {
		// Compatible with the previous implementation method
		endpoint = strings.TrimSpace(os.Getenv(fmt.Sprintf("%s_ENDPOINT", strings.ToUpper(productCode))))
	}
	if endpoint != "" {
		c.Endpoints.Store(productCode, endpoint)
		return nil
	}

	// Secondly, load endpoint from known rules
	if endpointFmt, ok := irregularProductEndpoint[productCode]; ok {
		if v, ok := irregularProductEndpointForIntlRegion[productCode]; ok && c.isInternationalRegion() {
			endpointFmt = v
		}
		if v, ok := irregularProductEndpointForIntlAccount[productCode]; ok && c.isInternationalAccount() {
			endpointFmt = v
		}
		if v, ok := irregularProductEndpointForIntlAccountIntlRegion[productCode]; ok && c.isInternationalAccount() && c.isInternationalRegion() {
			endpointFmt = v
		}
		if strings.Contains(endpointFmt, "%s") {
			endpointFmt = fmt.Sprintf(endpointFmt, regionID)
		}
		c.Endpoints.Store(productCode, endpointFmt)
		return nil
	}

	// Thirdly, load endpoint from location
	endpoint, err := c.describeEndpointForService(productCode, regionID)
	if err == nil {
		if v, ok := regularProductEndpointForIntlAccount[productCode]; ok && c.isInternationalAccount() {
			endpoint = v
		}
		if v, ok := regularProductEndpointReplace[endpoint]; ok {
			endpoint = v
		}
		c.Endpoints.Store(strings.ToLower(productCode), endpoint)
	} else if endpointFmt, ok := regularProductEndpoint[productCode]; ok {
		if v, ok := regularProductEndpointForIntlRegion[productCode]; ok && c.isInternationalRegion() {
			endpointFmt = v
		}
		if v, ok := regularProductEndpointForIntlAccount[productCode]; ok && c.isInternationalAccount() {
			endpointFmt = v
		}
		if v, ok := regularProductEndpointForIntlAccountIntlRegion[productCode]; ok && c.isInternationalAccount() && c.isInternationalRegion() {
			endpointFmt = v
		}
		if strings.Contains(endpointFmt, "%s") {
			endpointFmt = fmt.Sprintf(endpointFmt, regionID)
		}
		if v, ok := regularProductEndpointReplace[endpointFmt]; ok {
			endpointFmt = v
		}
		c.Endpoints.Store(productCode, endpointFmt)
		log.Printf("[WARN] loading %s endpoint got an error: %#v. Using the endpoint %s instead.", productCode, err, endpointFmt)
		return nil
	}
	return err
}

// Load current path endpoint file endpoints.xml, if failed, it will load from environment variables TF_ENDPOINT_PATH
func (c *Config) loadEndpointFromLocal() error {
	data, err := ioutil.ReadFile(localEndpointPath)
	if err != nil || len(data) <= 0 {
		d, e := ioutil.ReadFile(os.Getenv(localEndpointPathEnv))
		if e != nil {
			return e
		}
		data = d
	}
	var endpoints Endpoints
	err = xml.Unmarshal(data, &endpoints)
	if err != nil {
		return err
	}
	for _, endpoint := range endpoints.Endpoint {
		if endpoint.RegionIds.RegionId == string(c.RegionId) {
			for _, product := range endpoint.Products.Product {
				c.Endpoints.Store(strings.ToLower(product.ProductName), strings.TrimSpace(product.DomainName))
			}
		}
	}
	return nil
}

func (c *Config) LoadAPIEndpoint(productCode string, regionID string, accountID func() (string, error)) (string, error) {
	configEndpoints := productCodeToConfigEndpoints[productCode]
	if configEndpoints == "" {
		configEndpoints = productCode
	}

	if v, ok := c.Endpoints.Load(configEndpoints); !ok || v.(string) == "" {
		// Reuse a previously probed endpoint stored under the real product code.
		// productCodeToConfigEndpoints maps sls->log so the user-config key above
		// is "log", while loadEndpoint always stores under the product code (e.g.
		// "sls"). Without this pre-check the mapped read key never matches the
		// write key, so every SLS API call re-runs endpoint discovery. This
		// restores the symmetric cache path without affecting the user-config
		// priority above or non-mapped products whose keys already coincide.
		if v, ok := c.Endpoints.Load(productCode); ok && v.(string) != "" {
			return c.formatEndpointWithAccountID(productCode, v.(string), accountID)
		}
		if err := c.LoadEndpoint(productCode, regionID); err != nil {
			return "", fmt.Errorf("[ERROR] loading %s endpoint got an error: %#v.", productCode, err)
		}
	} else {
		return c.formatEndpointWithAccountID(productCode, v.(string), accountID)
	}
	if v, ok := c.Endpoints.Load(productCode); ok && v.(string) != "" {
		return c.formatEndpointWithAccountID(productCode, v.(string), accountID)
	}
	return "", fmt.Errorf("[ERROR] missing the product %s endpoint.", productCode)
}

func incrementalWait(firstDuration time.Duration, increaseDuration time.Duration) func() {
	retryCount := 1
	return func() {
		var waitTime time.Duration
		if retryCount == 1 {
			waitTime = firstDuration
		} else if retryCount > 1 {
			waitTime += increaseDuration
		}
		time.Sleep(waitTime)
		retryCount++
	}
}

func (c *Config) describeEndpointForService(productCode string, regionID string) (string, error) {
	locationCode := productCodeToLocationCode[productCode]
	if locationCode == "" {
		locationCode = productCode
	}
	args := location.CreateDescribeEndpointsRequest()
	args.ServiceCode = locationCode
	args.Id = c.RegionId
	args.Domain = c.LocationEndpoint
	if args.Domain == "" {
		args.Domain = LoadEndpoint(regionID, LOCATIONCode)
	}
	if args.Domain == "" {
		args.Domain = "location.aliyuncs.com"
	}

	locationClient, err := location.NewClientWithOptions(c.RegionId, c.SDKConfig(time.Duration(c.ClientReadTimeout)*time.Millisecond), c.AuthCredential(true))
	if err != nil {
		return "", fmt.Errorf("unable to initialize the location client: %#v", err)

	}
	locationClient.SetReadTimeout(time.Duration(c.ClientReadTimeout) * time.Millisecond)
	locationClient.SetConnectTimeout(time.Duration(c.ClientConnectTimeout) * time.Millisecond)
	locationClient.SourceIp = c.SourceIp
	locationClient.SecureTransport = c.SecureTransport
	defer locationClient.Shutdown()
	wait := incrementalWait(3*time.Second, 5*time.Second)
	var endpointResult string
	err = retry.Retry(2*time.Minute, func() *retry.RetryError {
		endpointsResponse, err := locationClient.DescribeEndpoints(args)
		if err != nil {
			if e, ok := err.(*errors.ClientError); ok && e.ErrorCode() == errors.TimeoutErrorCode {
				wait()
				args.Domain = "location-readonly.aliyuncs.com"
				return retry.RetryableError(err)
			}
			re := regexp.MustCompile("^Post [\"]*https://.*")
			if err.Error() != "" && re.MatchString(err.Error()) {
				wait()
				args.Domain = "location-readonly.aliyuncs.com"
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		if endpointsResponse != nil && len(endpointsResponse.Endpoints.Endpoint) > 0 {
			for _, e := range endpointsResponse.Endpoints.Endpoint {
				if e.Type == "openAPI" {
					endpointResult = e.Endpoint
					return nil
				}
			}
		}
		return nil
	})
	if err != nil {
		return "", fmt.Errorf("Describe %s endpoint using region: %#v got an error: %#v.", productCode, regionID, err)
	}
	if endpointResult == "" {
		return "", fmt.Errorf("There is no any available endpoint for %s in region %s.", productCode, regionID)
	}
	return endpointResult, nil
}

func (c *Config) isInternationalAccount() bool {
	return c.AccountType == "International"
}

func (c *Config) isInternationalRegion() bool {
	if c.RegionId == "cn-hongkong" {
		return true
	}
	return !strings.HasPrefix(c.RegionId, "cn-")
}

func (c *Config) formatEndpointWithAccountID(productCode string, endpoint string, accountID func() (string, error)) (string, error) {
	if endpoint == "" {
		return endpoint, nil
	}
	switch productCode {
	case "fc_open", "fc":
		accountId, err := accountID()
		if err != nil {
			return endpoint, err
		}
		if strings.HasPrefix(endpoint, accountId) {
			return endpoint, nil
		}
		// skip pop endpoint
		if !strings.HasSuffix(endpoint, ".fc.aliyuncs.com") {
			return endpoint, nil
		}
		return fmt.Sprintf("%s.%s", accountId, endpoint), nil
	}
	return endpoint, nil
}

const (
	OpenApiGatewayService          = "apigateway.cn-hangzhou.aliyuncs.com"
	OpenOtsService                 = "ots.cn-hangzhou.aliyuncs.com"
	OpenOssService                 = "oss-admin.aliyuncs.com"
	OpenNasService                 = "nas.cn-hangzhou.aliyuncs.com"
	OpenCdnService                 = "cdn.aliyuncs.com"
	OpenKmsService                 = "kms.cn-hangzhou.aliyuncs.com"
	OpenSaeService                 = "sae.cn-hangzhou.aliyuncs.com"
	OpenCmsService                 = "metrics.cn-hangzhou.aliyuncs.com"
	OpenDatahubService             = "datahub.aliyuncs.com"
	OpenOnsService                 = "ons.cn-hangzhou.aliyuncs.com"
	OpenDcdnService                = "dcdn.aliyuncs.com"
	OpenFcService                  = "fc-open.cn-hangzhou.aliyuncs.com"
	OpenAckService                 = "cs.aliyuncs.com"
	OpenPrivateLinkService         = "privatelink.cn-hangzhou.aliyuncs.com"
	OpenBrainIndustrialService     = "brain-industrial.cn-hangzhou.aliyuncs.com"
	OpenIotService                 = "iot.aliyuncs.com"
	OpenVsService                  = "vs.cn-shanghai.aliyuncs.com"
	OpenCrService                  = "cr.cn-hangzhou.aliyuncs.com"
	OpenMaxcomputeService          = "maxcompute.aliyuncs.com"
	OpenCloudStorageGatewayService = "sgw.cn-shanghai.aliyuncs.com"
	DataWorksService               = "dataworks.aliyuncs.com"
	OpenHbrService                 = "hbr.aliyuncs.com"
)

const (
	BssOpenAPIEndpointDomestic                = "business.aliyuncs.com"
	BssOpenAPIEndpointInternational           = "business.ap-southeast-1.aliyuncs.com"
	EcdOpenAPIEndpointUser                    = "eds-user.ap-southeast-1.aliyuncs.com"
	CloudFirewallOpenAPIEndpointControlPolicy = "cloudfw.ap-southeast-1.aliyuncs.com"
	SaSOpenAPIEndpointInternational           = "tds.ap-southeast-1.aliyuncs.com"
)
