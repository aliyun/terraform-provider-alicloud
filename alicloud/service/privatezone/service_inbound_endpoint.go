package privatezone

import (
	"context"
	"fmt"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	logs "github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
)

// inboundEndpointReadTimeout bounds a single DescribeInboundEndpoint retry loop when the caller's
// context carries no shorter deadline. It matches the read budget the SDKv2 DescribeEndpoint uses.
const inboundEndpointReadTimeout = 1 * time.Minute

// inboundEndpointFailStates are the terminal statuses DescribeInboundEndpoint documents as
// operation failures. Reaching one aborts a create or update wait with an error instead of
// waiting until the operation timeout.
var inboundEndpointFailStates = []string{"FAILED", "CHANGE_FAILED", "EXCEPTION"}

// DescribeInboundEndpointContext describes a single inbound endpoint by EndpointId. It returns a
// provider not-found error when the API answers with a typed 404 or a not-exists error code, so
// the Framework Read path can remove the resource from state. The not-exists code is not declared
// in the published DescribeInboundEndpoint contract, and the live API returns it with a spelling
// error ("InboundEndpoint.NotExsits" rather than "InboundEndpoint.NotExists"), so both spellings
// are matched; a malformed body is never turned into NotFound.
func (s *Service) DescribeInboundEndpointContext(ctx context.Context, id string) (object map[string]interface{}, err error) {
	client := s.client
	action := "DescribeInboundEndpoint"
	request := map[string]interface{}{
		"EndpointId": id,
		"Lang":       "en",
	}

	response, err := helper.RetryRpcPost(ctx, action, func() (map[string]interface{}, error) {
		return client.RpcPost("pvtz", "2018-01-01", action, nil, request, true)
	}, inboundEndpointReadTimeout, func(callErr error) bool {
		return helper.NeedRetry(callErr)
	})
	logs.AddDebug(action, response, request)
	if err != nil {
		if errs.IsTypedNotFound(err) || errs.IsExpectedErrors(err, []string{"InboundEndpoint.NotExists", "InboundEndpoint.NotExsits"}) {
			return object, errs.WrapErrorf(errs.NotFoundErr("PrivateZone:InboundEndpoint", id), errs.NotFoundMsg, response)
		}
		return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}

	v, err := jsonpath.Get("$", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$", response)
	}
	object = v.(map[string]interface{})
	return object, nil
}

// inboundEndpointStateRefreshFunc polls the inbound endpoint status for retry.StateChangeConf. A
// not-found object reports state "" so the caller's wait keeps polling while the resource is still
// being created; a fail-state status aborts the wait with an error.
func (s *Service) inboundEndpointStateRefreshFunc(ctx context.Context, id string, failStates []string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeInboundEndpointContext(ctx, id)
		if err != nil {
			if errs.IsProviderNotFound(err) {
				return nil, "", nil
			}
			return nil, "", errs.WrapError(err)
		}
		for _, failState := range failStates {
			if fmt.Sprint(object["Status"]) == failState {
				return object, fmt.Sprint(object["Status"]), errs.WrapError(errs.Error(errs.FailedToReachTargetStatus, fmt.Sprint(object["Status"])))
			}
		}
		return object, fmt.Sprint(object["Status"]), nil
	}
}

// inboundEndpointWaitForState polls the inbound endpoint until Status reaches one of target. A
// value in failStates aborts the wait with an error. ctx bounds the wait as well as the timeout.
func (s *Service) inboundEndpointWaitForState(ctx context.Context, id string, target []string, failStates []string, timeout time.Duration) (map[string]interface{}, error) {
	stateConf := &retry.StateChangeConf{
		Pending:    []string{},
		Target:     target,
		Refresh:    s.inboundEndpointStateRefreshFunc(ctx, id, failStates),
		Timeout:    timeout,
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
	}
	raw, err := stateConf.WaitForStateContext(ctx)
	if err != nil {
		return nil, err
	}
	if m, ok := raw.(map[string]interface{}); ok {
		return m, nil
	}
	return nil, nil
}
