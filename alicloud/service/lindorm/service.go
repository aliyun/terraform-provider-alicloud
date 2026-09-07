package lindorm

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/tags"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Service wraps the Lindorm (hitsdb 2020-06-15) API calls used by the
// resources and data sources of this package. The root alicloud package
// cannot be imported from a service package, so the encapsulated calls that
// previously lived on the root HitsdbService (lindorm-only methods) and
// LindormServiceV2 live here, merged onto a single *Service.
type Service struct {
	client *connectivity.AliyunClient
}

// NewService is the exported constructor: the client field stays unexported,
// so callers outside the package (including the acceptance tests) build
// values with it.
func NewService(client *connectivity.AliyunClient) *Service {
	return &Service{client: client}
}

func (s *Service) GetInstanceIpWhiteList(id string) (object []interface{}, err error) {
	var response map[string]interface{}
	client := s.client
	action := "GetInstanceIpWhiteList"
	request := map[string]interface{}{
		"InstanceId": id,
	}
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("hitsdb", "2020-06-15", action, nil, request, true)
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
		return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}
	v, err := jsonpath.Get("$.IpList", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.IpList", response)
	}
	if len(v.([]interface{})) < 1 {
		return object, nil
	} else {
		if fmt.Sprint(response["InstanceId"]) != id {
			return object, errs.WrapErrorf(errs.NotFoundErr("Lindorm", id), errs.NotFoundWithResponse, response)
		}
	}
	return v.([]interface{}), nil
}

func (s *Service) GetLindormInstanceEngineInfo(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	client := s.client
	action := "GetLindormInstanceEngineInfo"
	request := map[string]interface{}{
		"InstanceId": id,
	}
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("hitsdb", "2020-06-15", action, nil, request, true)
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
		return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}
	v, err := jsonpath.Get("$", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$", response)
	}
	object = v.(map[string]interface{})
	return object, nil
}

func (s *Service) DescribeLindormInstance(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	client := s.client
	action := "GetLindormInstance"
	request := map[string]interface{}{
		"InstanceId": id,
	}
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("hitsdb", "2020-06-15", action, nil, request, true)
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
		if errs.IsExpectedErrors(err, []string{"Lindorm.Errorcode.InstanceNotFound", "Instance.IsDeleted"}) {
			return object, errs.WrapErrorf(errs.NotFoundErr("Lindorm:Instance", id), errs.NotFoundMsg, errs.ProviderERROR, fmt.Sprint(response["RequestId"]))
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

func (s *Service) LindormInstanceStateRefreshFunc(id string, failStates []string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeLindormInstance(id)
		if err != nil {
			if errs.NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", errs.WrapError(err)
		}

		for _, failState := range failStates {
			if fmt.Sprint(object["InstanceStatus"]) == failState {
				return object, fmt.Sprint(object["InstanceStatus"]), errs.WrapError(errs.Error(errs.FailedToReachTargetStatus, fmt.Sprint(object["InstanceStatus"])))
			}
		}
		return object, fmt.Sprint(object["InstanceStatus"]), nil
	}
}

// LindormPostCreateReadTolerance keeps the last tolerated post-create read
// error. The SDK runs the refresh function in its own goroutine, so access is
// guarded.
type LindormPostCreateReadTolerance struct {
	mu      sync.Mutex
	lastErr error
}

func (t *LindormPostCreateReadTolerance) record(err error) {
	if t == nil {
		return
	}
	t.mu.Lock()
	t.lastErr = err
	t.mu.Unlock()
}

func (t *LindormPostCreateReadTolerance) last() error {
	if t == nil {
		return nil
	}
	t.mu.Lock()
	defer t.mu.Unlock()
	return t.lastErr
}

// LindormInstanceStateRefreshFuncWithPostCreateReadTolerance wraps
// LindormInstanceStateRefreshFunc so that a freshly created instance which is
// briefly unreadable (HTTP 451 OperationFailed from GetLindormInstance) is
// treated as still provisioning until deadline. Any other error is returned
// unchanged, and the original error is returned once the deadline passes, so a
// permanent failure is not masked and the wait stays bounded. The last tolerated
// error is recorded in tolerance so the caller can surface it when the wait
// itself times out. Only the Create wait uses this wrapper; Read and Delete keep
// the strict behaviour.
func (s *Service) LindormInstanceStateRefreshFuncWithPostCreateReadTolerance(id string, failStates []string, deadline time.Time, tolerance *LindormPostCreateReadTolerance) retry.StateRefreshFunc {
	refresh := s.LindormInstanceStateRefreshFunc(id, failStates)
	return func() (interface{}, string, error) {
		object, state, err := refresh()
		if err != nil && time.Now().Before(deadline) && isLindormTransientPostCreateReadError(err) {
			tolerance.record(err)
			return map[string]interface{}{"InstanceId": id}, "CREATING", nil
		}
		return object, state, err
	}
}

func isLindormTransientPostCreateReadError(err error) bool {
	if err == nil {
		return false
	}
	if e, ok := err.(*errs.ComplexError); ok {
		return isLindormTransientPostCreateReadError(e.Cause)
	}
	if e, ok := err.(*tea.SDKError); ok {
		code := ""
		if e.Code != nil {
			code = *e.Code
		}
		return tea.IntValue(e.StatusCode) == http.StatusUnavailableForLegalReasons && code == "OperationFailed"
	}
	// The shared client can also surface the legacy SDK error shape
	// (alibaba-cloud-sdk-go ServerError) which exposes HttpStatus/ErrorCode
	// instead of the tea fields.
	if e, ok := err.(interface {
		HttpStatus() int
		ErrorCode() string
	}); ok {
		return e.HttpStatus() == http.StatusUnavailableForLegalReasons && e.ErrorCode() == "OperationFailed"
	}
	return false
}

func (s *Service) ListTagResources(id string, resourceType string) (object interface{}, err error) {
	client := s.client
	action := "ListTagResources"
	request := map[string]interface{}{
		"RegionId":     s.client.RegionId,
		"ResourceType": resourceType,
		"ResourceId.1": id,
	}
	tags := make([]interface{}, 0)
	var response map[string]interface{}

	for {
		wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
		err = retry.Retry(5*time.Minute, func() *retry.RetryError {
			response, err := client.RpcPost("hitsdb", "2020-06-15", action, nil, request, true)
			if err != nil {
				if helper.NeedRetry(err) {
					wait()
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
			}
			logs.AddDebug(action, response, request)
			v, err := jsonpath.Get("$.TagResources", response)
			if err != nil {
				return retry.NonRetryableError(errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.TagResources", response))
			}
			if v != nil {
				tags = append(tags, v.([]interface{})...)
			}
			return nil
		})
		if err != nil {
			err = errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
			return
		}
		if response["NextToken"] == nil {
			break
		}
		request["NextToken"] = response["NextToken"]
	}

	return tags, nil
}

func (s *Service) SetResourceTags(d *schema.ResourceData, resourceType string) error {

	if d.HasChange("tags") {
		added, removed := tags.ParsingTags(d)
		client := s.client

		removedTagKeys := make([]string, 0)
		for _, v := range removed {
			if !tags.IgnoredTags(v, "") {
				removedTagKeys = append(removedTagKeys, v)
			}
		}
		if len(removedTagKeys) > 0 {
			action := "UnTagResources"
			request := map[string]interface{}{
				"RegionId":     s.client.RegionId,
				"ResourceType": resourceType,
				"ResourceId.1": d.Id(),
			}
			for i, key := range removedTagKeys {
				request[fmt.Sprintf("TagKey.%d", i+1)] = key
			}
			wait := helper.IncrementalWait(2*time.Second, 1*time.Second)
			err := retry.Retry(10*time.Minute, func() *retry.RetryError {
				response, err := client.RpcPost("hitsdb", "2020-06-15", action, nil, request, false)
				if err != nil {
					if helper.NeedRetry(err) {
						wait()
						return retry.RetryableError(err)

					}
					return retry.NonRetryableError(err)
				}
				logs.AddDebug(action, response, request)
				return nil
			})
			if err != nil {
				return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
			}
		}
		if len(added) > 0 {
			action := "TagResources"
			request := map[string]interface{}{
				"RegionId":     s.client.RegionId,
				"ResourceType": resourceType,
				"ResourceId.1": d.Id(),
			}
			count := 1
			for key, value := range added {
				request[fmt.Sprintf("Tag.%d.Key", count)] = key
				request[fmt.Sprintf("Tag.%d.Value", count)] = value
				count++
			}

			wait := helper.IncrementalWait(2*time.Second, 1*time.Second)
			err := retry.Retry(10*time.Minute, func() *retry.RetryError {
				response, err := client.RpcPost("hitsdb", "2020-06-15", action, nil, request, false)
				if err != nil {
					if helper.NeedRetry(err) {
						wait()
						return retry.RetryableError(err)

					}
					return retry.NonRetryableError(err)
				}
				logs.AddDebug(action, response, request)
				return nil
			})
			if err != nil {
				return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
			}
		}
	}
	return nil
}

// describeVSwitchAttributes is a local copy of VpcService.DescribeVSwitchWithTeadsl.
// The Lindorm instance create path falls back to it to resolve VpcId/ZoneId from a
// VSwitchId, and the root VpcService is not importable from this service package.
// It targets Vpc 2016-04-28 DescribeVSwitchAttributes.
func (s *Service) describeVSwitchAttributes(id string) (object map[string]interface{}, err error) {
	client := s.client
	action := "DescribeVSwitchAttributes"
	request := map[string]interface{}{
		"RegionId":  s.client.RegionId,
		"VSwitchId": id,
	}

	response, err := client.RpcPost("Vpc", "2016-04-28", action, nil, request, true)
	if err != nil {
		if errs.IsExpectedErrors(err, []string{"InvalidVswitchID.NotFound"}) {
			return nil, errs.WrapErrorf(err, errs.NotFoundMsg, errs.AlibabaCloudSdkGoERROR)
		}
		return nil, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}
	logs.AddDebug(action, response, request)
	if v, ok := response["VSwitchId"].(string); ok && v != id {
		return nil, errs.WrapErrorf(errs.NotFoundErr("vswitch", id), errs.NotFoundMsg, errs.ProviderERROR)
	}
	return response, nil
}
