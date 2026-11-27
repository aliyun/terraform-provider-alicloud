package client

import (
	"context"
	"fmt"
	"strings"
	"sync"

	pop "github.com/alibabacloud-go/alibabacloud-gateway-pop/client"
	spi "github.com/alibabacloud-go/alibabacloud-gateway-spi/client"
	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	rpcutil "github.com/alibabacloud-go/tea-rpc-utils/service"
	"github.com/alibabacloud-go/tea/dara"
	"github.com/alibabacloud-go/tea/tea"
)

func effectiveSignVersion(versions *sync.Map, product string) string {
	if versions == nil {
		return ""
	}
	for _, key := range []string{product, "default"} {
		if value, ok := versions.Load(key); ok {
			if version, ok := value.(string); ok && version != "" {
				return version
			}
		}
	}
	return ""
}

func (client *AliyunClient) defaultV4() bool {
	return client.config != nil && effectiveSignVersion(client.config.SignVersion, "default") == "v4"
}

func (client *AliyunClient) useV4ForProduct(product string) bool {
	return client.defaultV4() && effectiveSignVersion(client.config.SignVersion, strings.ToLower(product)) == "v4"
}

// UseSignV4 reports whether requests for the given product sign with V4,
// honoring the product-level sign_version override. Callers whose request
// belongs to a different product than the resource they serve (for example
// the OTS attachment resolving a vswitch through VPC) must branch on the
// request's own product: default = v4 with ots = v1 keeps the VPC call on
// the V4 path while the OTS calls stay legacy.
func (client *AliyunClient) UseSignV4(product string) bool {
	return client.useV4ForProduct(product)
}

func (client *AliyunClient) configureV4(conn *openapi.Client, product string) error {
	return client.configureV4WithRegion(conn, product, client.config.RegionId)
}

// configureV4WithRegion signs conn with an explicit credential-scope region.
// The shared entry points sign with the client's physical region; a call
// pinned to a regionless central endpoint passes the gateway's central scope
// instead, because the gateway derives the scope region from the endpoint host
// and rejects a physical region there.
func (client *AliyunClient) configureV4WithRegion(conn *openapi.Client, product, region string) error {
	if region == "" {
		return fmt.Errorf("sign_version.default = v4 requires a signing region")
	}
	gateway, err := pop.NewClient()
	if err != nil {
		return err
	}
	conn.RegionId = tea.String(region)
	conn.ProductId = tea.String(strings.ToLower(product))
	conn.SignatureVersion = tea.String("v4")
	return conn.SetGatewayClient(gateway)
}

func (client *AliyunClient) rpcV4(ctx context.Context, product, endpoint, method, version, action string, query, body map[string]interface{}, retry bool) (map[string]interface{}, error) {
	return client.rpcV4WithApi(ctx, product, endpoint, "/", method, version, action, query, body, retry)
}

// rpcV4WithApi is the pathname-aware V4 entry for calls migrated off the old
// SDK clients whose request reached a fixed API path (the SLS POP client
// appended "/open-api" to its Domain); the pop gateway signs that pathname and
// appends it to the bare endpoint host.
func (client *AliyunClient) rpcV4WithApi(ctx context.Context, product, endpoint, pathname, method, version, action string, query, body map[string]interface{}, retry bool) (map[string]interface{}, error) {
	return client.rpcV4WithScopeRegion(ctx, product, endpoint, pathname, client.config.RegionId, method, version, action, query, body, retry)
}

// rpcV4WithScopeRegion is the credential-scope-region-explicit form of
// rpcV4WithApi, for calls pinned to an endpoint whose gateway associates a
// scope region other than the client's physical region: the regionless
// central hosts scope credentials to the gateway's central region.
func (client *AliyunClient) rpcV4WithScopeRegion(ctx context.Context, product, endpoint, pathname, scopeRegion, method, version, action string, query, body map[string]interface{}, retry bool) (map[string]interface{}, error) {
	flat, err := rpcutil.FlattenParams(query)
	if err != nil {
		return nil, err
	}
	params := &openapi.Params{Style: tea.String("RPC"), Method: tea.String(method), Version: tea.String(version),
		Action: tea.String(action), Pathname: tea.String(pathname), ReqBodyType: tea.String("formData")}
	req := &openapi.OpenApiRequest{Query: flat}
	if body != nil {
		req.Body = body
	}
	return client.commonV4WithScopeRegion(ctx, product, endpoint, scopeRegion, params, req, retry)
}

func (client *AliyunClient) commonV4(ctx context.Context, product, endpoint string, params *openapi.Params, req *openapi.OpenApiRequest, retry bool) (map[string]interface{}, error) {
	return client.commonV4WithScopeRegion(ctx, product, endpoint, client.config.RegionId, params, req, retry)
}

func (client *AliyunClient) commonV4WithScopeRegion(ctx context.Context, product, endpoint, scopeRegion string, params *openapi.Params, req *openapi.OpenApiRequest, retry bool) (map[string]interface{}, error) {
	cfg := client.teaRoaOpenapiConfig
	if tea.StringValue(params.Style) == "RPC" {
		cfg = client.teaRpcOpenapiConfig
	}
	cfg.SetEndpoint(endpoint).SetCredential(client.config.Credential).SetProtocol(client.config.Protocol)
	useProxy, proxyURL, bypass, err := client.resolveTeaProxy(endpoint)
	if err != nil {
		return nil, err
	}
	if useProxy {
		cfg.SetHttpProxy(proxyURL).SetHttpsProxy(proxyURL)
	} else if bypass != "" {
		cfg.SetNoProxy(bypass)
	}
	conn, err := openapi.NewClient(&cfg)
	if err != nil {
		return nil, err
	}
	if err := client.configureV4WithRegion(conn, product, scopeRegion); err != nil {
		return nil, err
	}
	params.AuthType = tea.String("AK")
	params.BodyType = tea.String("json")
	params.Protocol = cfg.Protocol
	// ExecuteWithCtx reads RetryOptions rather than RuntimeOptions.Autoretry.
	// Retain the Tea common clients' three retries for retryable errors only.
	var response map[string]interface{}
	for attempt := 0; ; attempt++ {
		response, err = conn.ExecuteWithCtx(ctx, params, req, &dara.RuntimeOptions{})
		if !retry || attempt >= 3 || !tea.BoolValue(tea.Retryable(err)) || ctx.Err() != nil {
			break
		}
	}
	if err != nil {
		return nil, err
	}
	if body, ok := response["body"].(map[string]interface{}); ok {
		if tea.StringValue(params.Style) == "RPC" {
			body["_headers"] = response["headers"]
		}
		response = body
	}
	return response, formatError(response, nil)
}

// The context-aware RPC path owns its HTTP transport and error decoding.
// Sign with the SDK without handing a signed request to SDK debug logging.
func (client *AliyunClient) signContextRPCV4(product, endpoint, version, action string, request *tea.Request, query, body map[string]interface{}) error {
	gateway, err := pop.NewClient()
	if err != nil {
		return err
	}
	flat, err := rpcutil.FlattenParams(query)
	if err != nil {
		return err
	}
	if flat == nil {
		flat = make(map[string]*string)
	}
	if client.config.SourceIp != "" {
		flat["SourceIp"] = tea.String(client.config.SourceIp)
	}
	if client.config.SecureTransport != "" {
		flat["SecureTransport"] = tea.String(client.config.SecureTransport)
	}
	signed := &spi.InterceptorContextRequest{
		ProductId: tea.String(strings.ToLower(product)), Action: tea.String(action), Version: tea.String(version),
		Credential: client.config.Credential, AuthType: tea.String("AK"), UserAgent: tea.String(client.config.UserAgent()),
		Pathname: request.Pathname, Method: request.Method, Query: flat, ReqBodyType: tea.String("formData"),
	}
	if body != nil {
		signed.Body = body
	}
	if client.config.RegionId == "" {
		return fmt.Errorf("sign_version.default = v4 requires a signing region")
	}
	err = gateway.ModifyRequest(&spi.InterceptorContext{Request: signed, Configuration: &spi.InterceptorContextConfiguration{
		RegionId: tea.String(client.config.RegionId), Endpoint: tea.String(endpoint),
	}}, nil)
	if err == nil {
		request.Query, request.Headers, request.Body = signed.Query, signed.Headers, signed.Stream
	}
	return err
}
