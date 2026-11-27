package client

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	"github.com/alibabacloud-go/tea/tea"
	otsCommon "github.com/aliyun/aliyun-tablestore-go-sdk/common"
	"github.com/aliyun/aliyun-tablestore-go-sdk/tablestore"
	otsTunnel "github.com/aliyun/aliyun-tablestore-go-sdk/tunnel"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
)

// OTS reaches the service through four request paths: the control-plane POP
// API (product "ots", version "2016-06-20"), the REST data plane (product
// "tablestore", version "2020-12-09"), and the table and tunnel SDK clients.
// Only the table and tunnel SDKs can sign AKv4 themselves, so the other two
// paths are routed through the provider's common V4 layer. All four are
// selected by sign_version.default = v4 without an "ots" product override
// away from v4; every other configuration keeps the master signing behavior.

const (
	otsPopApiVersion  = "2016-06-20"
	otsV4Product      = "ots"
	otsV4RegionError  = "sign_version.default = v4 requires a signing region"
	otsTunnelV4Notice = "OTS tunnel workers are not supported with sign_version.default = v4"
)

// UseOtsV4 reports whether the OTS request paths sign with V4. The service
// package branches on this because the control-plane and REST call sites own
// legacy SDK request structs that cannot be signed by V4.
func (client *AliyunClient) UseOtsV4() bool {
	return client.useV4ForProduct(otsV4Product)
}

// otsV4SigningRegion returns the region used in OTS V4 credential scopes and
// the x-ots-signregion header. EO instances configure the physical region;
// the common V4 layer keeps the same requirement for the gateway paths.
func (client *AliyunClient) otsV4SigningRegion() (string, error) {
	if client.config.RegionId == "" {
		return "", fmt.Errorf("%s", otsV4RegionError)
	}
	return client.config.RegionId, nil
}

// otsV4CredentialsProvider adapts the provider credential chain to the
// tablestore SDK's CredentialsV4 mechanism. Both SDKs ask for credentials on
// every request, so GetRefreshCredential keeps temporary credentials as fresh
// as the v1 clients, and common.CreateByCredentials derives the daily OTS V4
// signing key for the configured region.
type otsV4CredentialsProvider struct {
	client *AliyunClient
	region string
}

func (provider *otsV4CredentialsProvider) GetCredentials() otsCommon.Credentials {
	accessKey, secretKey, token := provider.client.config.GetRefreshCredential()
	// V4Credentials.GetSnapshot drops V4SigningStsToken (an SDK v1.9.3 gap)
	// while both SDK request layers read the token from the snapshot, so the
	// token is restored on the snapshot the SDK signs with.
	snapshot := otsCommon.CreateByCredentials(&otsCommon.DefaultCredentials{
		AccessKeyID:     accessKey,
		AccessKeySecret: secretKey,
		SecurityToken:   token,
	}, provider.region).GetSnapshot()
	snapshot.V4SigningStsToken = token
	return snapshot
}

func (client *AliyunClient) newOtsV4CredentialsProvider() (otsCommon.CredentialsProvider, error) {
	region, err := client.otsV4SigningRegion()
	if err != nil {
		return nil, err
	}
	return &otsV4CredentialsProvider{client: client, region: region}, nil
}

// otsTunnelV4Client adapts a *otsTunnel.TunnelApi carrying both the V4
// credentials provider and the external headers to the TunnelClient interface
// the OTS service expects. The SDK has no constructor combining a credentials
// provider with external headers, and DefaultTunnelClient wraps an
// unexported api, so the adapter delegates the metadata API and reports
// tunnel workers explicitly: no provider code builds workers, and a silent
// fallback to v1 signing would hide the missing V4 support.
type otsTunnelV4Client struct {
	*otsTunnel.TunnelApi
}

func (c *otsTunnelV4Client) NewTunnelWorker(tunnelId string, workerConfig *otsTunnel.TunnelWorkerConfig) (otsTunnel.TunnelWorker, error) {
	return nil, fmt.Errorf("%s", otsTunnelV4Notice)
}

// OtsPopV4 invokes an OTS control-plane POP API (product "ots", version
// "2016-06-20") through the provider's RPC path, which selects the common V4
// signing layer for this product. The legacy POP SDK client can only sign
// HMAC-SHA1, so the V4 branch bypasses it entirely; the endpoint chain
// (config.Endpoints "ots" key, then env rules and location) is the one
// rpcRequest already resolves. SourceIp and SecureTransport are carried as
// query parameters the same way the legacy client injected them. The response
// is the raw map; callers convert it into the legacy POP response struct.
func (client *AliyunClient) OtsPopV4(method string, action string, query map[string]interface{}) (map[string]interface{}, error) {
	if query == nil {
		query = make(map[string]interface{})
	}
	if _, ok := query["RegionId"]; !ok {
		query["RegionId"] = client.RegionId
	}
	if _, ok := query["SourceIp"]; !ok && client.config.SourceIp != "" {
		query["SourceIp"] = client.config.SourceIp
	}
	if _, ok := query["SecureTransport"]; !ok && client.config.SecureTransport != "" {
		query["SecureTransport"] = client.config.SecureTransport
	}
	return otsV4Response(client.rpcRequest(method, otsV4Product, otsPopApiVersion, action, query, nil, false, ""))
}

// OtsRestV4 sends an OTS REST data-plane request through the common V4 layer
// with the same ROA style the roaRequest V4 branch uses. The endpoint chain
// mirrors NewOtsRoaClient (config.Endpoints under the product code, then the
// env/rule/location loader), with the explicit ots endpoint consulted for EO
// instances whose station endpoint is not derivable from a public cloud rule.
// The call is a single attempt: the OTS callers own their
// temporarily-unavailable retry loop, matching the v1 path's plain runtime.
func (client *AliyunClient) OtsRestV4(product string, version string, actionPath string, method string, urlQuery map[string]*string, headers map[string]*string, requestBody map[string]interface{}) (map[string]interface{}, error) {
	endpoint := ""
	if v, ok := client.config.Endpoints.Load(product); ok && v.(string) != "" {
		endpoint = v.(string)
	}
	if endpoint == "" && client.config.OtsEndpoint != "" {
		endpoint = client.config.OtsEndpoint
	}
	if endpoint == "" {
		if err := client.loadEndpoint(context.Background(), product); err != nil {
			return nil, err
		}
		if v, ok := client.config.Endpoints.Load(product); ok && v.(string) != "" {
			endpoint = v.(string)
		}
	}
	if endpoint == "" {
		return nil, fmt.Errorf("[ERROR] missing the product %s endpoint", product)
	}
	// The common V4 layer composes the URL from the configured protocol, so
	// an endpoint that conventionally carries a scheme (the explicit ots
	// endpoint) must be reduced to its host, the same trim the FC and MNS
	// clients apply before handing endpoints to Tea SDK clients.
	endpoint = strings.TrimPrefix(strings.TrimPrefix(endpoint, "https://"), "http://")
	params := &openapi.Params{Style: tea.String("ROA"), Method: tea.String(method), Version: tea.String(version),
		Pathname: tea.String(actionPath), ReqBodyType: tea.String("json")}
	request := &openapi.OpenApiRequest{Query: urlQuery, Headers: headers, Body: requestBody}
	return otsV4Response(client.commonV4(context.Background(), product, endpoint, params, request, false))
}

// The legacy OTS client reports missing instances as ServerError. Preserve
// that lifecycle signal when the V4 gateway returns the same service error.
func otsV4Response(result map[string]interface{}, err error) (map[string]interface{}, error) {
	var sdkErr *tea.SDKError
	if errors.As(err, &sdkErr) {
		code := tea.StringValue(sdkErr.Code)
		switch code {
		case errs.InstanceNotFound, errs.RamInstanceNotFound, errs.NotFound:
			return nil, errs.WrapErrorf(err, errs.NotFoundMsg, errs.AlibabaCloudSdkGoERROR)
		}
		status := tea.IntValue(sdkErr.StatusCode)
		if (status == http.StatusBadRequest || status == http.StatusOK) && strings.Contains(strings.ToLower(tea.StringValue(sdkErr.Message)), errs.MessageInstanceNotFound) {
			return nil, errs.WrapErrorf(err, errs.NotFoundMsg, errs.AlibabaCloudSdkGoERROR)
		}
	}
	return result, err
}

// applyTableStoreV4 injects the V4 credentials provider into a table client
// built with NewClientWithExternalHeader, preserving the external headers and
// the table store config the v1 constructor already applied. The credentials
// handed to the constructor only back the default provider that is replaced
// here, so the values themselves become unused.
func applyTableStoreV4(tableStoreClient *tablestore.TableStoreClient, provider otsCommon.CredentialsProvider) {
	tablestore.SetCredentialsProvider(provider)(tableStoreClient)
}
