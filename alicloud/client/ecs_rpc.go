package client

import (
	"strings"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/responses"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
)

// EcsSdkRpc takes the generated request and response for the same ECS action.
func (client *AliyunClient) EcsSdkRpc(req requests.AcsRequest, resp responses.AcsResponse) (interface{}, error) {
	if !client.useV4ForProduct("ecs") {
		return client.WithEcsClient(func(c *ecs.Client) (interface{}, error) { return resp, c.DoAction(req, resp) })
	}
	raw, err := client.legacySdkRpc(req, "/", resp, strings.HasPrefix(req.GetActionName(), "Describe") || strings.HasPrefix(req.GetActionName(), "List"))
	if err == nil {
		err = LegacyRpcResponse(raw, resp)
	}
	return resp, err
}

// EcsCommonRpc keeps the original SDK response for security group rule actions.
func (client *AliyunClient) EcsCommonRpc(req *requests.CommonRequest) (interface{}, error) {
	if !client.useV4ForProduct("ecs") {
		return client.WithEcsClient(func(c *ecs.Client) (interface{}, error) { return c.ProcessCommonRequest(req) })
	}
	req.TransToAcsRequest()
	return client.LegacySdkRpc(req, "/", responses.NewCommonResponse())
}
