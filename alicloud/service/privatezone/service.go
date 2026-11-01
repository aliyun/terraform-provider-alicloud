package privatezone

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/tags"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Status mirrors the root-package sentinel type used by the legacy flat layout
// to drive resource lifecycle waits. It is kept package-local because the
// PrivateZone package cannot import the root alicloud package; only the
// Status values actually referenced by pvtz (Deleted) are reproduced here.
type Status string

const (
	// Deleted marks a resource as removed; the WaitFor* helpers return nil
	// once the describe call reports NotFound while waiting for Deleted.
	Deleted = Status("Deleted")

	// DefaultTimeout is the legacy flat-layout default wait timeout in
	// seconds, reproduced locally for WaitFor* call sites that previously
	// referenced the root-package constant.
	DefaultTimeout = 120
)

// Service wraps the PrivateZone (pvtz) API calls used by the resources and
// data sources of this package.
type Service struct {
	client *connectivity.AliyunClient
}

// NewService is the exported constructor: the client field stays unexported,
// so callers outside the package (including the flat layout's acceptance
// tests) build values with it.
func NewService(client *connectivity.AliyunClient) *Service {
	return &Service{client: client}
}

func (s *Service) DescribeZoneBasic(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	client := s.client
	action := "DescribeZoneInfo"

	request := map[string]interface{}{
		"RegionId": s.client.RegionId,
		"ZoneId":   id,
	}

	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(1*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("pvtz", "2018-01-01", action, nil, request, true)
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
		if errs.IsExpectedErrors(err, []string{"Zone.Invalid.Id", "Zone.Invalid.UserId", "Zone.NotExists", "ZoneVpc.NotExists.VpcId"}) {
			return object, errs.WrapErrorf(errs.NotFoundErr("PvtzZone", id), errs.NotFoundMsg, response)
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

func (s *Service) DescribeZoneAttachment(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	client := s.client
	action := "DescribeZoneInfo"

	request := map[string]interface{}{
		"RegionId": s.client.RegionId,
		"ZoneId":   id,
	}

	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(1*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("pvtz", "2018-01-01", action, nil, request, true)
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
		if errs.IsExpectedErrors(err, []string{"Zone.Invalid.Id", "Zone.Invalid.UserId", "Zone.NotExists", "ZoneVpc.NotExists.VpcId"}) {
			return object, errs.WrapErrorf(errs.NotFoundErr("PvtzZoneAttachment", id), errs.NotFoundMsg, response)
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

func (s *Service) WaitForZoneAttachment(id string, vpcIdMap map[string]string, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	var vpcId string
	for {
		object, err := s.DescribeZoneAttachment(id)
		if err != nil && !errs.NotFoundError(err) {
			return errs.WrapError(err)
		}

		equal := true
		vpcs := object["BindVpcs"].(map[string]interface{})["Vpc"].([]interface{})
		if len(vpcs) == len(vpcIdMap) {
			for _, vpc := range vpcs {
				vpc := vpc.(map[string]interface{})
				if _, ok := vpcIdMap[vpc["VpcId"].(string)]; !ok {
					equal = false
					vpcId = vpc["VpcId"].(string)
					break
				}
			}
		} else {
			equal = false
		}
		if equal {
			break
		}
		if time.Now().After(deadline) {
			return errs.WrapErrorf(err, errs.WaitTimeoutMsg, id, helper.GetFunc(1), timeout, "", vpcId, errs.ProviderERROR)
		}
	}
	return nil
}

func (s *Service) DescribeZoneRecord(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	client := s.client
	action := "DescribeZoneRecords"
	parts, err := helper.ParseResourceId(id, 2)
	if err != nil {
		err = errs.WrapError(err)
		return
	}
	request := map[string]interface{}{
		"RegionId":   s.client.RegionId,
		"ZoneId":     parts[1],
		"PageNumber": 1,
		"PageSize":   20,
	}
	for {
		wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
		err = retry.Retry(5*time.Minute, func() *retry.RetryError {
			response, err = client.RpcPost("pvtz", "2018-01-01", action, nil, request, true)
			if err != nil {
				if errs.IsExpectedErrors(err, []string{"System.Busy", "ServiceUnavailable", "Throttling.User"}) {
					wait()
					return retry.RetryableError(err)
				}
				if errs.IsExpectedErrors(err, []string{"Zone.Invalid.Id", "Zone.Invalid.UserId", "Zone.NotExists", "ZoneVpc.NotExists.VpcId"}) {
					err = errs.WrapErrorf(errs.NotFoundErr("PvtzZoneRecord", id), errs.NotFoundMsg, errs.ProviderERROR)
					return retry.NonRetryableError(err)
				}
				err = errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
				return retry.NonRetryableError(err)
			}
			return nil
		})
		if err != nil {
			return
		}
		logs.AddDebug(action, response, request)
		v, err := jsonpath.Get("$.Records.Record", response)
		if err != nil {
			return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.Records.Record", response)
		}
		if len(v.([]interface{})) < 1 {
			return object, errs.WrapErrorf(errs.NotFoundErr("PrivateZone", id), errs.NotFoundWithResponse, response)
		}
		for _, v := range v.([]interface{}) {
			if fmt.Sprint(v.(map[string]interface{})["RecordId"]) == parts[0] {
				return v.(map[string]interface{}), nil
			}
		}
		if len(v.([]interface{})) < request["PageSize"].(int) {
			break
		}
		request["PageNumber"] = request["PageNumber"].(int) + 1
	}
	return
}

func (s *Service) WaitForZone(id string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)

	for {
		object, err := s.DescribeZone(id)
		if err != nil {
			if errs.NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			} else {
				return errs.WrapError(err)
			}
		}
		if object["ZoneId"] == id {
			return nil
		}
		if time.Now().After(deadline) {
			return errs.WrapErrorf(err, errs.WaitTimeoutMsg, id, helper.GetFunc(1), timeout, object["ZoneId"], id, errs.ProviderERROR)
		}

	}
}

func (s *Service) WaitForZoneAttachmentStatus(id string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)

	for {
		object, err := s.DescribeZoneAttachment(id)
		if err != nil {
			if errs.NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			} else {
				return errs.WrapError(err)
			}
		}
		if object["ZoneId"] == id {
			return nil
		}
		if time.Now().After(deadline) {
			return errs.WrapErrorf(err, errs.WaitTimeoutMsg, id, helper.GetFunc(1), timeout, object["ZoneId"], id, errs.ProviderERROR)
		}

	}
}

func (s *Service) WaitForZoneRecord(id string, status Status, timeout int) error {
	parts, err := helper.ParseResourceId(id, 2)
	if err != nil {
		return errs.WrapError(err)
	}
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)

	for {
		object, err := s.DescribeZoneRecord(id)
		if err != nil {
			if errs.NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			} else {
				return errs.WrapError(err)
			}
		}
		if strconv.FormatInt(object["RecordId"].(int64), 10) == parts[0] {
			return nil
		}
		if time.Now().After(deadline) {
			return errs.WrapErrorf(err, errs.WaitTimeoutMsg, id, helper.GetFunc(1), timeout, strconv.FormatInt(object["RecordId"].(int64), 10), id, errs.ProviderERROR)
		}

	}
}

func (s *Service) DescribeUserVpcAuthorization(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	client := s.client
	action := "DescribeUserVpcAuthorizations"
	parts, err := helper.ParseResourceId(id, 2)
	if err != nil {
		err = errs.WrapError(err)
		return
	}
	request := map[string]interface{}{
		"AuthType":         parts[1],
		"AuthorizedUserId": parts[0],
	}
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("pvtz", "2018-01-01", action, nil, request, true)
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
		if errs.IsExpectedErrors(err, []string{"System.Busy"}) {
			return object, errs.WrapErrorf(errs.NotFoundErr("PrivateZone:UserVpcAuthorization", id), errs.NotFoundMsg, errs.ProviderERROR, fmt.Sprint(response["RequestId"]))
		}
		return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}
	v, err := jsonpath.Get("$.Users", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.Users", response)
	}
	if len(v.([]interface{})) < 1 {
		return object, errs.WrapErrorf(errs.NotFoundErr("PrivateZone", id), errs.NotFoundWithResponse, response)
	} else {
		if fmt.Sprint(v.([]interface{})[0].(map[string]interface{})["AuthType"]) != parts[1] {
			return object, errs.WrapErrorf(errs.NotFoundErr("PrivateZone", id), errs.NotFoundWithResponse, response)
		}
	}
	object = v.([]interface{})[0].(map[string]interface{})
	return object, nil
}

func (s *Service) DescribeEndpoint(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	client := s.client
	action := "DescribeResolverEndpoint"
	request := map[string]interface{}{
		"EndpointId": id,
	}
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {

		response, err = client.RpcPost("pvtz", "2018-01-01", action, nil, request, true)
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
		if errs.IsExpectedErrors(err, []string{"ResolverEndpoint.NotExists"}) {
			return object, errs.WrapErrorf(errs.NotFoundErr("PrivateZone:Endpoint", id), errs.NotFoundMsg, errs.ProviderERROR, fmt.Sprint(response["RequestId"]))
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

func (s *Service) EndpointStateRefreshFunc(id string, failStates []string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeEndpoint(id)
		if err != nil {
			if errs.NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
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

func (s *Service) DescribeRule(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	client := s.client
	action := "DescribeResolverRule"
	request := map[string]interface{}{
		"RuleId": id,
	}
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("pvtz", "2018-01-01", action, nil, request, true)
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
		if errs.IsExpectedErrors(err, []string{"ResolverRule.NotExists"}) {
			return object, errs.WrapErrorf(errs.NotFoundErr("PrivateZone:Rule", id), errs.NotFoundMsg, errs.ProviderERROR)
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

func (s *Service) DescribeRuleAttachment(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	client := s.client
	action := "DescribeResolverRule"
	request := map[string]interface{}{
		"RuleId": id,
	}
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("pvtz", "2018-01-01", action, nil, request, true)
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
		if errs.IsExpectedErrors(err, []string{"ResolverRule.NotExists"}) {
			return object, errs.WrapErrorf(errs.NotFoundErr("PrivateZone:RuleAttachment", id), errs.NotFoundMsg, errs.ProviderERROR)
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

func (s *Service) DescribeSyncEcsHostTask(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	client := s.client
	action := "DescribeSyncEcsHostTask"
	request := map[string]interface{}{
		"ZoneId": id,
	}
	response, err = client.RpcPost("pvtz", "2018-01-01", action, nil, request, true)
	if err != nil {
		if errs.IsExpectedErrors(err, []string{"System.Busy", "Ecs.SyncTask.NotExists", "ServiceUnavailable", "Throttling.User"}) {
			err = errs.WrapErrorf(errs.NotFoundErr("PvtzZone", id), errs.NotFoundMsg, errs.ProviderERROR)
			return object, err
		}
		err = errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
		return object, err
	}
	logs.AddDebug(action, response, request)
	v, err := jsonpath.Get("$", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$", response)
	}
	object = v.(map[string]interface{})
	return object, nil
}

func (s *Service) DescribeZone(id string) (object map[string]interface{}, err error) {
	object, err = s.DescribeZoneBasic(id)
	if err != nil {
		return nil, errs.WrapError(err)
	}
	syncObj, err := s.DescribeSyncEcsHostTask(id)
	if err != nil {
		if errs.NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_pvtz_zone pvtzService.DescribeSyncEcsHostTask Failed!!! %s", err)
			return object, nil
		}
		return nil, errs.WrapError(err)
	}
	object["SyncHostTask"] = syncObj
	return object, nil
}

func (s *Service) SetResourceTags(d *schema.ResourceData, resourceType string) error {
	if d.HasChange("tags") {
		var err error
		var action string
		client := s.client
		var request map[string]interface{}
		var response map[string]interface{}

		added, removed := tags.ParsingTags(d)
		removedTagKeys := make([]string, 0)
		for _, v := range removed {
			if !tags.IgnoredTags(v, "") {
				removedTagKeys = append(removedTagKeys, v)
			}
		}
		if len(removedTagKeys) > 0 {
			action = "UnTagResources"
			request = make(map[string]interface{})
			request["ResourceId.1"] = d.Id()
			request["RegionId"] = client.RegionId
			for i, key := range removedTagKeys {
				request[fmt.Sprintf("TagKey.%d", i+1)] = key
			}

			request["ResourceType"] = resourceType
			wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
			err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
				response, err = client.RpcPost("pvtz", "2018-01-01", action, nil, request, true)

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
			action = "TagResources"
			request = make(map[string]interface{})
			request["ResourceId.1"] = d.Id()
			request["RegionId"] = client.RegionId
			count := 1
			for key, value := range added {
				request[fmt.Sprintf("Tag.%d.Key", count)] = key
				request[fmt.Sprintf("Tag.%d.Value", count)] = value
				count++
			}

			request["ResourceType"] = resourceType
			wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
			err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
				response, err = client.RpcPost("pvtz", "2018-01-01", action, nil, request, true)

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

func (s *Service) DescribeListTagResources(id string) (object map[string]interface{}, err error) {

	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	action := "ListTagResources"
	request = make(map[string]interface{})
	query = make(map[string]interface{})

	request["ResourceId.1"] = id
	request["RegionId"] = client.RegionId

	request["ResourceType"] = "ZONE"
	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(1*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("pvtz", "2018-01-01", action, query, request, false)

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
		if errs.IsExpectedErrors(err, []string{}) {
			return object, errs.WrapErrorf(errs.NotFoundErr("ZONE", id), errs.NotFoundMsg, errs.ProviderERROR, fmt.Sprint(response["RequestId"]))
		}
		return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}

	return response, nil
}

func (s *Service) DescribeServiceStatus() (object map[string]interface{}, err error) {
	var response map[string]interface{}
	client := s.client
	action := "DescribeUserServiceStatus"
	request := make(map[string]interface{})

	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("pvtz", "2018-01-01", action, nil, request, true)
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
		return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, "PvtzService", action, errs.AlibabaCloudSdkGoERROR)
	}

	v, err := jsonpath.Get("$", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, "PvtzService", "$", response)
	}

	object = v.(map[string]interface{})

	return object, nil
}

func (s *Service) ServiceStatusStateRefreshFunc(failStates []string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeServiceStatus()
		if err != nil {
			if errs.NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
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
