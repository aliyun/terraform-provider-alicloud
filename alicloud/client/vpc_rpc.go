package client

import (
	"strings"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/responses"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/vpc"
)

// VpcSdkRpc runs one typed VPC SDK request through the shared
// sign-version-aware RPC entry (LegacySdkRpc). Without
// sign_version.default = "v4" the generated vpc client owns serialization,
// signing, retries and error decoding exactly as the per-call-site
// WithVpcClient closures did, so the legacy signature, endpoint chain,
// request semantics and error mapping are preserved. With it (and no "vpc"
// override away from v4) the typed request is translated onto the common V4
// gateway, which honors the product endpoint, the physical signing region,
// temporary credentials and the sign_version override precedence.
//
// The typed response passed by the caller is the same object the legacy
// branch decodes into; the V4 branch rehydrates it from the shared response
// map with LegacyRpcResponse. Either way the populated response is returned,
// so callers keep their existing type assertions and field access:
//
//	raw, err := client.VpcSdkRpc(request, vpc.CreateDescribeVSwitchAttributesResponse())
//	response, _ := raw.(*vpc.DescribeVSwitchAttributesResponse)
func (client *AliyunClient) VpcSdkRpc(request requests.AcsRequest, response responses.AcsResponse) (interface{}, error) {
	// The legacy branch decodes straight into the typed response via DoAction,
	// so the map LegacySdkRpc returns would only be marshalled and unmarshalled
	// back into the object the SDK already populated. Route the legacy branch
	// directly and keep the LegacyRpcResponse rehydration for the V4 branch,
	// whose shared response map is the only source of the typed fields.
	if !client.useV4ForProduct(request.GetProduct()) {
		return client.WithVpcClient(func(c *vpc.Client) (interface{}, error) { return response, c.DoAction(request, response) })
	}
	raw, err := client.legacySdkRpc(request, "/", response, strings.HasPrefix(request.GetActionName(), "Describe") || strings.HasPrefix(request.GetActionName(), "List"))
	if err != nil {
		return nil, err
	}
	if raw == nil {
		return response, nil
	}
	if err := LegacyRpcResponse(raw, response); err != nil {
		return nil, err
	}
	return response, nil
}
