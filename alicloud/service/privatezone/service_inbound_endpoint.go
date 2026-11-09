package privatezone

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/alibabacloud-go/tea/tea"
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

// inboundEndpointNotExistsCodes are the not-exists responses for the inbound endpoint. The code
// is not declared in the published contracts, and the live API returns it with a spelling error
// ("InboundEndpoint.NotExsits" rather than "InboundEndpoint.NotExists"), so both are matched.
var inboundEndpointNotExistsCodes = []string{"InboundEndpoint.NotExists", "InboundEndpoint.NotExsits"}

// IsInboundEndpointNotFound reports whether err proves the endpoint no longer exists.
// helper.RetryRpcPost wraps the API error in errs.ComplexError, which has no Unwrap method, so
// errs.IsTypedNotFound (errors.As based) cannot see the tea SDK error through the wrap; walk the
// wrap cause explicitly. The sweeper uses the same decision.
func IsInboundEndpointNotFound(err error) bool {
	for err != nil {
		if errs.IsTypedNotFound(err) {
			return true
		}
		switch e := err.(type) {
		case *errs.ComplexError:
			err = e.Cause
			continue
		case *tea.SDKError:
			return inboundEndpointTeaCodeMatches(e, inboundEndpointNotExistsCodes)
		}
		// Other error types are safe with the shared matcher; tea errors are handled above
		// because errs.IsExpectedErrors dereferences SDKError.Data, which may be nil.
		if errs.IsExpectedErrors(err, inboundEndpointNotExistsCodes) {
			return true
		}
		if wrapped, ok := err.(interface{ Unwrap() error }); ok {
			err = wrapped.Unwrap()
			continue
		}
		return false
	}
	return false
}

func inboundEndpointTeaCodeMatches(e *tea.SDKError, codes []string) bool {
	if e == nil {
		return false
	}
	code := tea.StringValue(e.Code)
	data := tea.StringValue(e.Data)
	for _, expected := range codes {
		// Exact match only: the service also answers endpoint-scoped failures with a short code
		// such as "InboundEndpoint", which is a prefix of the not-exists code but says nothing
		// about the endpoint being gone. A prefix match would erase live state and make a failed
		// delete look successful.
		if code != "" && code == expected {
			return true
		}
		if data != "" && strings.Contains(data, expected) {
			return true
		}
	}
	return false
}

// DescribeInboundEndpointContext describes a single inbound endpoint by EndpointId. It returns a
// provider not-found error when the API answers with a typed 404 or a not-exists error code, so
// the Framework Read path can remove the resource from state; a malformed body is never turned
// into NotFound. The retry backoff is bound to ctx so a cancelled operation stops retrying.
func (s *Service) DescribeInboundEndpointContext(ctx context.Context, id string) (object map[string]interface{}, err error) {
	client := s.client
	action := "DescribeInboundEndpoint"
	request := map[string]interface{}{
		"EndpointId": id,
		"Lang":       "en",
	}

	response, err := helper.RetryRpcPost(ctx, action, func() (map[string]interface{}, error) {
		return client.RpcPost("pvtz", "2018-01-01", action, nil, request, true)
	}, inboundEndpointReadTimeout, helper.NeedRetry)
	logs.AddDebug(action, response, request)
	if err != nil {
		if IsInboundEndpointNotFound(err) {
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

// CreateInboundEndpoint calls AddInboundEndpoint. The request map, including its single
// ClientToken, is built by the caller before this call, so every retry of the same logical create
// reuses the same token instead of creating a second endpoint.
func (s *Service) CreateInboundEndpoint(ctx context.Context, request map[string]interface{}) (map[string]interface{}, error) {
	return s.writeInboundEndpoint(ctx, "AddInboundEndpoint", request)
}

// UpdateInboundEndpoint calls UpdateInboundEndpoint. Only the fields present in request are
// changed; the caller decides which values differ from state.
func (s *Service) UpdateInboundEndpoint(ctx context.Context, request map[string]interface{}) (map[string]interface{}, error) {
	return s.writeInboundEndpoint(ctx, "UpdateInboundEndpoint", request)
}

// DeleteInboundEndpoint calls DeleteInboundEndpoint. A successful call only means the service
// accepted the delete; call inboundEndpointWaitForDeleted when the endpoint must be observably
// gone before dependents are released.
func (s *Service) DeleteInboundEndpoint(ctx context.Context, request map[string]interface{}) (map[string]interface{}, error) {
	return s.writeInboundEndpoint(ctx, "DeleteInboundEndpoint", request)
}

func (s *Service) writeInboundEndpoint(ctx context.Context, action string, request map[string]interface{}) (map[string]interface{}, error) {
	response, err := helper.RetryRpcPost(ctx, action, func() (map[string]interface{}, error) {
		return s.client.RpcPost("pvtz", "2018-01-01", action, nil, request, false)
	}, inboundEndpointTimeout, helper.NeedRetry)
	logs.AddDebug(action, response, request)
	return response, err
}

// inboundEndpointWaitForDeleted polls DescribeInboundEndpoint until it reports NotFound. Delete
// is asynchronous, so returning right after DeleteInboundEndpoint lets Terraform release
// dependents while the endpoint still exists. The same provider NotFound chain used by Read and
// CheckDestroy decides when the object is gone; any other read error aborts the wait.
func (s *Service) inboundEndpointWaitForDeleted(ctx context.Context, id string, timeout time.Duration) error {
	stateConf := &retry.StateChangeConf{
		Pending: []string{"Deleting"},
		Target:  []string{"Deleted"},
		Refresh: func() (interface{}, string, error) {
			object, err := s.DescribeInboundEndpointContext(ctx, id)
			if err != nil {
				if errs.IsProviderNotFound(err) {
					return object, "Deleted", nil
				}
				return nil, "", errs.WrapError(err)
			}
			return object, "Deleting", nil
		},
		Timeout:    timeout,
		Delay:      2 * time.Second,
		MinTimeout: 2 * time.Second,
	}
	_, err := stateConf.WaitForStateContext(ctx)
	return err
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
