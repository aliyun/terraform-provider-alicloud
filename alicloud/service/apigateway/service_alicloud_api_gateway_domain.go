package apigateway

import (
	"fmt"
	"strings"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
)

// DescribeApiGatewayDomain returns the custom domain bound to an API group.
// id is "<GroupId>:<DomainName>", matching the CloudSpec primary key.
func (s *ApiGatewayServiceV2) DescribeApiGatewayDomain(id string) (object map[string]interface{}, err error) {
	client := s.client
	groupId, domainName, err := splitApiGatewayDomainId(id)
	if err != nil {
		return object, err
	}
	action := "DescribeDomain"
	// CloudAPI RPC params go in the body slot (5th arg), matching the package
	// convention shared by every apigateway RpcPost caller.
	request := map[string]interface{}{
		"GroupId":    groupId,
		"DomainName": domainName,
	}

	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	var response map[string]interface{}
	err = retry.Retry(1*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("CloudAPI", "2016-07-14", action, nil, request, true)
		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(action, response, request)
	if err != nil {
		// CloudAPI returns NotFoundDomain when the domain is not bound to the group;
		// NotFoundError covers a 404 or a NotFound-bearing message as a fallback.
		if errs.IsExpectedErrors(err, []string{"NotFoundDomain"}) || errs.NotFoundError(err) {
			return object, errs.WrapErrorf(errs.NotFoundErr("ApiGatewayDomain", id), errs.NotFoundMsg, response)
		}
		return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}
	// DescribeDomain returns a flat object for the requested GroupId+DomainName.
	if g, _ := response["GroupId"].(string); g != groupId {
		return object, errs.WrapErrorf(errs.NotFoundErr("ApiGatewayDomain", id), errs.NotFoundMsg, response)
	}
	return response, nil
}

// splitApiGatewayDomainId splits "<GroupId>:<DomainName>" into its two parts.
func splitApiGatewayDomainId(id string) (groupId, domainName string, err error) {
	parts := strings.SplitN(id, COLON_SEPARATED, 2)
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return "", "", fmt.Errorf("invalid ApiGatewayDomain id %q: expected <GroupId>:<DomainName>", id)
	}
	return parts[0], parts[1], nil
}
