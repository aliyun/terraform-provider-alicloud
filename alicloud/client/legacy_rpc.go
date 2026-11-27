package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"strings"

	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/responses"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/cms"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/sls"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/client/config"
)

// LegacyRpcResponse rehydrates a shared-path RPC response map into the typed
// alibaba-cloud-sdk-go response, so a caller migrated off its per-product SDK
// client keeps its downstream typed field access unchanged. It is the
// response-side counterpart of LegacySdkRpc and stays free of product-specific
// envelope assumptions; callers that can rely on a fixed top-level field (for
// example the CMS Code) should validate it around this helper.
func LegacyRpcResponse(raw map[string]interface{}, response interface{}) error {
	data, err := json.Marshal(raw)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, response)
}

// LegacySdkRpc keeps the original product client for legacy signing. Only V4
// needs the typed request translated to the shared gateway's query and body.
func (client *AliyunClient) LegacySdkRpc(request requests.AcsRequest, apiPathname string, response responses.AcsResponse) (map[string]interface{}, error) {
	return client.legacySdkRpc(request, apiPathname, response, false)
}

// retryV4 applies only to callers that enable gateway retries.
func (client *AliyunClient) legacySdkRpc(request requests.AcsRequest, apiPathname string, response responses.AcsResponse, retryV4 bool) (map[string]interface{}, error) {
	product := request.GetProduct()
	if !client.useV4ForProduct(product) {
		return client.callLegacyProductRPC(request, response)
	}
	// Every generated request factory pre-seeds the headers map with the
	// SDK's own tracking headers (x-sdk-client, x-sdk-invoke-type and the
	// x-acs-version/action pair from InitWithApiInfo), so a plain len()
	// check would reject every real call. The snapshot must therefore be
	// taken before InitParams — the first serialization of this request —
	// and compared after: InitParams turns position:"header" business
	// fields into additional entries (it never clears or rewrites existing
	// ones), and only those would be silently dropped by the shared entry
	// points, so fail loudly when one appears. Raw bodies land in the
	// content the same way. Path and domain parameters live only on
	// ROA-style requests, which the AcsRequest interface does not expose
	// (there is no GetPathParams/GetDomainParams to consult): every migrated
	// call is a pure query/form RPC request, and a ROA-style caller would
	// fail on the missing path contract server-side rather than succeed
	// silently.
	baselineHeaders := make(map[string]string, len(request.GetHeaders()))
	for key, value := range request.GetHeaders() {
		baselineHeaders[key] = value
	}
	if err := requests.InitParams(request); err != nil {
		return nil, err
	}
	for key, value := range request.GetHeaders() {
		if baselineHeaders[key] != value {
			return nil, fmt.Errorf("LegacySdkRpc does not support requests with header parameters: %s", request.GetActionName())
		}
	}
	if len(request.GetContent()) > 0 {
		return nil, fmt.Errorf("LegacySdkRpc does not support requests with raw content: %s", request.GetActionName())
	}
	query := make(map[string]interface{}, len(request.GetQueryParams()))
	for key, value := range request.GetQueryParams() {
		query[key] = value
	}
	// buildRequestWithSigner resolves the signing region from the request and
	// falls back to the client region, and signRpcRequest then fills RegionId
	// into the query when the request leaves it unset. Keep both contracts.
	regionId := client.RegionId
	if request.GetRegionId() != "" {
		regionId = request.GetRegionId()
	}
	if _, ok := query["RegionId"]; !ok && regionId != "" {
		query["RegionId"] = regionId
	}
	var body map[string]interface{}
	if len(request.GetFormParams()) > 0 {
		body = make(map[string]interface{}, len(request.GetFormParams()))
		for key, value := range request.GetFormParams() {
			body[key] = value
		}
	}
	if apiPathname == "" {
		apiPathname = "/"
	}
	if strings.EqualFold(product, "vpc") && regionId != "" && regionId != client.RegionId {
		endpoint := client.config.VpcEndpoint
		if endpoint == "" {
			var err error
			endpoint, err = client.loadApiEndpointForRegion(context.Background(), "vpc", regionId)
			if err != nil {
				return nil, err
			}
		}
		return client.rpcV4WithScopeRegion(context.Background(), product, stripEndpointScheme(endpoint), apiPathname, regionId, request.GetMethod(), request.GetVersion(), request.GetActionName(), query, body, retryV4)
	}
	endpoint, err := client.loadApiEndpoint(strings.ToLower(ConvertKebabToSnake(product)))
	if err != nil {
		return nil, err
	}
	// The pop gateway appends the pathname to the endpoint host, so the
	// signing endpoint must be a bare host even when the configured
	// endpoint carries a scheme prefix.
	return client.rpcV4WithApi(context.Background(), product, stripEndpointScheme(endpoint), apiPathname,
		request.GetMethod(), request.GetVersion(), request.GetActionName(), query, body, retryV4)
}

// Let the SDK own request serialization, signing, retries and error decoding,
// just as the generated ECS, CMS and SLS methods do.
func (client *AliyunClient) callLegacyProductRPC(request requests.AcsRequest, response responses.AcsResponse) (map[string]interface{}, error) {
	call := func(conn *sdk.Client) (interface{}, error) { return nil, conn.DoAction(request, response) }
	var err error
	switch strings.ToLower(request.GetProduct()) {
	case "ecs":
		_, err = client.WithEcsClient(func(c *ecs.Client) (interface{}, error) { return call(&c.Client) })
	case "cms":
		_, err = client.WithCmsClient(func(c *cms.Client) (interface{}, error) { return call(&c.Client) })
	case "sls":
		_, err = client.WithLogPopClient(func(c *sls.Client) (interface{}, error) { return call(&c.Client) })
	default:
		return nil, fmt.Errorf("unsupported legacy RPC product %q", request.GetProduct())
	}
	if err != nil {
		return nil, err
	}
	// The SLS audit reader consumes the full DescribeApp body, including fields
	// missing from the SDK model. Other callers consume the generated fields and
	// rely on the SDK's scalar coercion and zero values for omitted fields.
	var body []byte
	if _, raw := response.(*sls.DescribeAppResponse); raw {
		body = response.GetHttpContentBytes()
	} else {
		body, err = json.Marshal(response)
		if err != nil {
			return nil, err
		}
	}
	var result map[string]interface{}
	if len(body) > 0 {
		decoder := json.NewDecoder(bytes.NewReader(body))
		decoder.UseNumber()
		if err := decoder.Decode(&result); err != nil {
			return nil, fmt.Errorf("decode %s response: %w", request.GetActionName(), err)
		}
	}
	return result, nil
}

// OpenNasService uses the original account activation endpoint unless NAS
// itself selects V4. A regional NAS data endpoint is independent of activation.
func (client *AliyunClient) OpenNasService() (map[string]interface{}, error) {
	if !client.useV4ForProduct("nas") {
		conn, err := client.NewTeaCommonClient(config.OpenNasService)
		if err != nil {
			return nil, err
		}
		return conn.DoRequest(tea.String("OpenNASService"), nil, tea.String("POST"), tea.String("2017-06-26"), tea.String("AK"), nil, nil, &util.RuntimeOptions{})
	}
	endpoint := stripEndpointScheme(client.config.NasEndpoint)
	if endpoint == "" {
		endpoint = config.OpenNasService
	}
	return client.RpcPostWithEndpoint("Nas", "2017-06-26", "OpenNASService", nil, nil, false, endpoint)
}

func stripEndpointScheme(endpoint string) string {
	endpoint = strings.TrimSpace(endpoint)
	for _, prefix := range []string{"https://", "http://"} {
		if len(endpoint) >= len(prefix) && strings.EqualFold(endpoint[:len(prefix)], prefix) {
			endpoint = endpoint[len(prefix):]
			break
		}
	}
	return endpoint
}
