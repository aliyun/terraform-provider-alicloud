package config

import (
	"context"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk"
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
// LoadEndpoint resolves a product endpoint without a caller context. It keeps
// the historical entry points working and delegates to the context-aware
// resolver.
func (c *Config) LoadEndpoint(productCode string, regionID string) error {
	return c.loadEndpoint(context.Background(), productCode, regionID, nil)
}

// LoadEndpointWithContext bounds Location discovery with the caller's context.
func (c *Config) LoadEndpointWithContext(ctx context.Context, productCode string, regionID string, discoverFns ...func(context.Context, string, string) (string, error)) error {
	var discover func(context.Context, string, string) (string, error)
	if len(discoverFns) > 0 {
		discover = discoverFns[0]
	}
	return c.loadEndpoint(ctx, productCode, regionID, discover)
}

func (c *Config) loadEndpoint(ctx context.Context, productCode string, regionID string, discover func(context.Context, string, string) (string, error)) error {
	endpoint, err := c.resolveEndpointForRegion(ctx, productCode, regionID, discover)
	if err != nil {
		return err
	}
	// Every caller of this chain passes a lower-case product code (see the
	// NOTE on LoadEndpoint), so this single key keeps the historical
	// per-branch write keys equivalent.
	c.Endpoints.Store(productCode, endpoint)
	return nil
}

// resolveEndpointForRegion walks the endpoint rules for one region —
// environment overrides first, then the irregular product rules, then a
// Location probe bounded by ctx, then the regular rules as a fallback — and
// returns the resolved host. Caching stays with the callers: the shared
// per-product entry for the provider's own region, and the region-scoped
// composite entry for a request targeting another region.
func (c *Config) resolveEndpointForRegion(ctx context.Context, productCode string, regionID string, discover func(context.Context, string, string) (string, error)) (string, error) {
	if err := ctx.Err(); err != nil {
		return "", err
	}
	// Firstly, load endpoint from environment variables
	endpoint := strings.TrimSpace(os.Getenv(fmt.Sprintf("ALIBABA_CLOUD_ENDPOINT_%s", strings.ToUpper(productCode))))
	if endpoint == "" {
		// Compatible with the previous implementation method
		endpoint = strings.TrimSpace(os.Getenv(fmt.Sprintf("%s_ENDPOINT", strings.ToUpper(productCode))))
	}
	if endpoint != "" {
		return endpoint, nil
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
		return endpointFmt, nil
	}

	// Thirdly, load endpoint from location
	if discover == nil {
		discover = c.describeEndpointForServiceWithContext
	}
	endpoint, err := discover(ctx, productCode, regionID)
	if ctxErr := ctx.Err(); ctxErr != nil {
		return "", ctxErr
	}
	if err == nil {
		if v, ok := regularProductEndpointForIntlAccount[productCode]; ok && c.isInternationalAccount() {
			endpoint = v
		}
		if v, ok := regularProductEndpointReplace[endpoint]; ok {
			endpoint = v
		}
		return endpoint, nil
	}
	if endpointFmt, ok := regularProductEndpoint[productCode]; ok {
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
		log.Printf("[WARN] loading %s endpoint got an error: %#v. Using the endpoint %s instead.", productCode, err, endpointFmt)
		return endpointFmt, nil
	}
	return "", err
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

// LoadAPIEndpoint resolves a product endpoint without a caller context.
func (c *Config) LoadAPIEndpoint(productCode string, regionID string, accountID func() (string, error)) (string, error) {
	return c.loadAPIEndpoint(context.Background(), productCode, regionID, accountID, nil)
}

// LoadAPIEndpointWithContext applies the configured endpoint rules with bounded discovery.
func (c *Config) LoadAPIEndpointWithContext(ctx context.Context, productCode string, regionID string, accountID func() (string, error), discoverFns ...func(context.Context, string, string) (string, error)) (string, error) {
	var discover func(context.Context, string, string) (string, error)
	if len(discoverFns) > 0 {
		discover = discoverFns[0]
	}
	return c.loadAPIEndpoint(ctx, productCode, regionID, accountID, discover)
}

// RegionalEndpointCacheKey is the cache key of a product endpoint resolved
// for one specific region. The "@" separator keeps it from ever colliding
// with the plain product-code keys every other endpoint reader and writer
// uses, so region-scoped entries can live in the same shared map without
// being observed by the per-product lookups.
func RegionalEndpointCacheKey(productCode, regionID string) string {
	return productCode + "@" + regionID
}

// LoadAPIEndpointForRegion resolves the endpoint of a product for a region
// other than the provider's own. Unlike LoadAPIEndpointWithContext it
// neither reads nor writes the shared per-product cache entry, because that
// entry holds the endpoint of whichever region resolved first — a
// same-region call usually warms it during apply — while a cross-region
// request signs with its own region, so a warm shared cache would sign one
// region against another region's host. The environment override and the
// per-region rule chain still apply, and resolved endpoints are cached under
// RegionalEndpointCacheKey.
func (c *Config) LoadAPIEndpointForRegion(ctx context.Context, productCode string, regionID string, accountID func() (string, error), discoverFns ...func(context.Context, string, string) (string, error)) (string, error) {
	var discover func(context.Context, string, string) (string, error)
	if len(discoverFns) > 0 {
		discover = discoverFns[0]
	}
	return c.loadAPIEndpointForRegion(ctx, productCode, regionID, accountID, discover)
}

func (c *Config) loadAPIEndpointForRegion(ctx context.Context, productCode string, regionID string, accountID func() (string, error), discover func(context.Context, string, string) (string, error)) (string, error) {
	if err := ctx.Err(); err != nil {
		return "", err
	}
	key := RegionalEndpointCacheKey(productCode, regionID)
	if v, ok := c.Endpoints.Load(key); ok && v.(string) != "" {
		return c.formatEndpointWithAccountID(productCode, v.(string), accountID)
	}
	endpoint, err := c.resolveEndpointForRegion(ctx, productCode, regionID, discover)
	if err != nil {
		if ctxErr := ctx.Err(); ctxErr != nil {
			return "", ctxErr
		}
		return "", fmt.Errorf("[ERROR] loading %s endpoint got an error: %w.", productCode, err)
	}
	if endpoint == "" {
		return "", fmt.Errorf("[ERROR] missing the product %s endpoint.", productCode)
	}
	c.Endpoints.Store(key, endpoint)
	return c.formatEndpointWithAccountID(productCode, endpoint, accountID)
}

func (c *Config) loadAPIEndpoint(ctx context.Context, productCode string, regionID string, accountID func() (string, error), discover func(context.Context, string, string) (string, error)) (string, error) {
	if err := ctx.Err(); err != nil {
		return "", err
	}
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
		if err := c.loadEndpoint(ctx, productCode, regionID, discover); err != nil {
			if ctxErr := ctx.Err(); ctxErr != nil {
				return "", ctxErr
			}
			return "", fmt.Errorf("[ERROR] loading %s endpoint got an error: %w.", productCode, err)
		}
	} else {
		return c.formatEndpointWithAccountID(productCode, v.(string), accountID)
	}
	if err := ctx.Err(); err != nil {
		return "", err
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

type locationContextTransport struct {
	ctx  context.Context
	base http.RoundTripper
}

func (t locationContextTransport) RoundTrip(request *http.Request) (*http.Response, error) {
	if err := t.ctx.Err(); err != nil {
		return nil, err
	}
	return t.base.RoundTrip(request.WithContext(t.ctx))
}

// LocationCodeForProduct returns the Location service code for a product.
func LocationCodeForProduct(productCode string) string {
	locationCode := productCodeToLocationCode[productCode]
	if locationCode == "" {
		return productCode
	}
	return locationCode
}

// describeEndpointForService keeps the historical probe entry point and
// delegates to the context-aware resolver.
func (c *Config) describeEndpointForService(productCode string, regionID string) (string, error) {
	return c.describeEndpointForServiceWithContext(context.Background(), productCode, regionID)
}

func (c *Config) describeEndpointForServiceWithContext(ctx context.Context, productCode string, regionID string) (string, error) {
	locationCode := productCodeToLocationCode[productCode]
	if locationCode == "" {
		locationCode = productCode
	}
	args := location.CreateDescribeEndpointsRequest()
	args.ServiceCode = locationCode
	// Location must describe the endpoints of the region the request targets,
	// not the provider's client region, or a cross-region probe returns the
	// client region's host.
	args.Id = regionID
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
	defer locationClient.Shutdown()
	readTimeout := time.Duration(c.ClientReadTimeout) * time.Millisecond
	connectTimeout := time.Duration(c.ClientConnectTimeout) * time.Millisecond
	locationClient.SetReadTimeout(readTimeout)
	locationClient.SetConnectTimeout(connectTimeout)
	if deadline, ok := ctx.Deadline(); ok {
		remaining := time.Until(deadline)
		if remaining <= 0 {
			return "", ctx.Err()
		}
		if readTimeout <= 0 || remaining < readTimeout {
			readTimeout = remaining
			locationClient.SetReadTimeout(remaining)
		}
		if connectTimeout <= 0 || remaining < connectTimeout {
			connectTimeout = remaining
			locationClient.SetConnectTimeout(remaining)
		}
	}
	locationClient.SourceIp = c.SourceIp
	locationClient.SecureTransport = c.SecureTransport
	if ctx.Done() != nil {
		transport := c.getTransport()
		transport.DialContext = sdk.Timeout(connectTimeout)
		transport.Proxy = http.ProxyFromEnvironment
		locationClient.SetTransport(locationContextTransport{ctx: ctx, base: transport})
	}
	wait := incrementalWait(3*time.Second, 5*time.Second)
	waitRetry := func() error {
		done := make(chan struct{})
		go func() {
			wait()
			close(done)
		}()
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-done:
			return nil
		}
	}
	var endpointResult string
	err = retry.RetryContext(ctx, 2*time.Minute, func() *retry.RetryError {
		if ctxErr := ctx.Err(); ctxErr != nil {
			return retry.NonRetryableError(ctxErr)
		}
		endpointsResponse, err := locationClient.DescribeEndpoints(args)
		if err != nil {
			if e, ok := err.(*errors.ClientError); ok && e.ErrorCode() == errors.TimeoutErrorCode {
				if waitErr := waitRetry(); waitErr != nil {
					return retry.NonRetryableError(waitErr)
				}
				args.Domain = "location-readonly.aliyuncs.com"
				return retry.RetryableError(err)
			}
			re := regexp.MustCompile("^Post [\"]*https://.*")
			if err.Error() != "" && re.MatchString(err.Error()) {
				if waitErr := waitRetry(); waitErr != nil {
					return retry.NonRetryableError(waitErr)
				}
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
