package alicloud

import (
	"fmt"
	"log"
	"strconv"

	"github.com/PaesslerAG/jsonpath"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"

	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
)

type PvtzService struct {
	client *connectivity.AliyunClient
}

func (s *PvtzService) DescribePvtzZoneBasic(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	client := s.client
	action := "DescribeZoneInfo"
	request := map[string]interface{}{
		"RegionId": s.client.RegionId,
		"ZoneId":   id,
	}
	response, err = client.RpcPost("pvtz", "2018-01-01", action, nil, request, true)
	if err != nil {
		if IsExpectedErrors(err, []string{"Zone.Invalid.Id", "Zone.Invalid.UserId", "Zone.NotExists", "ZoneVpc.NotExists.VpcId"}) {
			err = WrapErrorf(NotFoundErr("PvtzZone", id), NotFoundMsg, ProviderERROR)
			return object, err
		}
		err = WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
		return object, err
	}
	addDebug(action, response, request)
	v, err := jsonpath.Get("$", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$", response)
	}
	object = v.(map[string]interface{})
	return object, nil
}

func (s *PvtzService) DescribePvtzZoneAttachment(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	client := s.client
	action := "DescribeZoneInfo"
	request := map[string]interface{}{
		"RegionId": s.client.RegionId,
		"ZoneId":   id,
	}
	response, err = client.RpcPost("pvtz", "2018-01-01", action, nil, request, true)
	if err != nil {
		if IsExpectedErrors(err, []string{"Zone.Invalid.Id", "Zone.Invalid.UserId", "Zone.NotExists", "ZoneVpc.NotExists.VpcId"}) {
			err = WrapErrorf(NotFoundErr("PvtzZoneAttachment", id), NotFoundMsg, ProviderERROR)
			return object, err
		}
		err = WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
		return object, err
	}
	addDebug(action, response, request)
	v, err := jsonpath.Get("$", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$", response)
	}
	object = v.(map[string]interface{})
	return object, nil
}

func (s *PvtzService) WaitForZoneAttachment(id string, vpcIdMap map[string]string, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	var vpcId string
	for {
		object, err := s.DescribePvtzZoneAttachment(id)
		if err != nil && !NotFoundError(err) {
			return WrapError(err)
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
			return WrapErrorf(err, WaitTimeoutMsg, id, GetFunc(1), timeout, "", vpcId, ProviderERROR)
		}
	}
	return nil
}

func (s *PvtzService) DescribePvtzZoneRecord(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	client := s.client
	action := "DescribeZoneRecords"
	parts, err := ParseResourceId(id, 2)
	if err != nil {
		err = WrapError(err)
		return
	}
	request := map[string]interface{}{
		"RegionId":   s.client.RegionId,
		"ZoneId":     parts[1],
		"PageNumber": 1,
		"PageSize":   20,
	}
	for {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			response, err = client.RpcPost("pvtz", "2018-01-01", action, nil, request, true)
			if err != nil {
				if IsExpectedErrors(err, []string{"System.Busy", "ServiceUnavailable", "Throttling.User"}) {
					wait()
					return resource.RetryableError(err)
				}
				if IsExpectedErrors(err, []string{"Zone.Invalid.Id", "Zone.Invalid.UserId", "Zone.NotExists", "ZoneVpc.NotExists.VpcId"}) {
					err = WrapErrorf(NotFoundErr("PvtzZoneRecord", id), NotFoundMsg, ProviderERROR)
					return resource.NonRetryableError(err)
				}
				err = WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
				return resource.NonRetryableError(err)
			}
			return nil
		})
		if err != nil {
			return
		}
		addDebug(action, response, request)
		v, err := jsonpath.Get("$.Records.Record", response)
		if err != nil {
			return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.Records.Record", response)
		}
		if len(v.([]interface{})) < 1 {
			return object, WrapErrorf(NotFoundErr("PrivateZone", id), NotFoundWithResponse, response)
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

func (s *PvtzService) WaitForPvtzZone(id string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)

	for {
		object, err := s.DescribePvtzZone(id)
		if err != nil {
			if NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			} else {
				return WrapError(err)
			}
		}
		if object["ZoneId"] == id {
			return nil
		}
		if time.Now().After(deadline) {
			return WrapErrorf(err, WaitTimeoutMsg, id, GetFunc(1), timeout, object["ZoneId"], id, ProviderERROR)
		}

	}
}

func (s *PvtzService) WaitForPvtzZoneAttachment(id string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)

	for {
		object, err := s.DescribePvtzZoneAttachment(id)
		if err != nil {
			if NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			} else {
				return WrapError(err)
			}
		}
		if object["ZoneId"] == id {
			return nil
		}
		if time.Now().After(deadline) {
			return WrapErrorf(err, WaitTimeoutMsg, id, GetFunc(1), timeout, object["ZoneId"], id, ProviderERROR)
		}

	}
}

func (s *PvtzService) WaitForPvtzZoneRecord(id string, status Status, timeout int) error {
	parts, err := ParseResourceId(id, 2)
	if err != nil {
		return WrapError(err)
	}
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)

	for {
		object, err := s.DescribePvtzZoneRecord(id)
		if err != nil {
			if NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			} else {
				return WrapError(err)
			}
		}
		if strconv.FormatInt(object["RecordId"].(int64), 10) == parts[0] {
			return nil
		}
		if time.Now().After(deadline) {
			return WrapErrorf(err, WaitTimeoutMsg, id, GetFunc(1), timeout, strconv.FormatInt(object["RecordId"].(int64), 10), id, ProviderERROR)
		}

	}
}

func (s *PvtzService) DescribePvtzUserVpcAuthorization(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	client := s.client
	action := "DescribeUserVpcAuthorizations"
	parts, err := ParseResourceId(id, 2)
	if err != nil {
		err = WrapError(err)
		return
	}
	request := map[string]interface{}{
		"AuthType":         parts[1],
		"AuthorizedUserId": parts[0],
	}
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = client.RpcPost("pvtz", "2018-01-01", action, nil, request, true)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)
	if err != nil {
		if IsExpectedErrors(err, []string{"System.Busy"}) {
			return object, WrapErrorf(NotFoundErr("PrivateZone:UserVpcAuthorization", id), NotFoundMsg, ProviderERROR, fmt.Sprint(response["RequestId"]))
		}
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}
	v, err := jsonpath.Get("$.Users", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.Users", response)
	}
	if len(v.([]interface{})) < 1 {
		return object, WrapErrorf(NotFoundErr("PrivateZone", id), NotFoundWithResponse, response)
	} else {
		if fmt.Sprint(v.([]interface{})[0].(map[string]interface{})["AuthType"]) != parts[1] {
			return object, WrapErrorf(NotFoundErr("PrivateZone", id), NotFoundWithResponse, response)
		}
	}
	object = v.([]interface{})[0].(map[string]interface{})
	return object, nil
}

func (s *PvtzService) DescribePvtzEndpoint(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	client := s.client
	action := "DescribeResolverEndpoint"
	request := map[string]interface{}{
		"EndpointId": id,
	}
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {

		response, err = client.RpcPost("pvtz", "2018-01-01", action, nil, request, true)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)
	if err != nil {
		if IsExpectedErrors(err, []string{"ResolverEndpoint.NotExists"}) {
			return object, WrapErrorf(NotFoundErr("PrivateZone:Endpoint", id), NotFoundMsg, ProviderERROR, fmt.Sprint(response["RequestId"]))
		}
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}
	v, err := jsonpath.Get("$", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$", response)
	}
	object = v.(map[string]interface{})
	return object, nil
}

func (s *PvtzService) PvtzEndpointStateRefreshFunc(id string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribePvtzEndpoint(id)
		if err != nil {
			if NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", WrapError(err)
		}

		for _, failState := range failStates {
			if fmt.Sprint(object["Status"]) == failState {
				return object, fmt.Sprint(object["Status"]), WrapError(Error(FailedToReachTargetStatus, fmt.Sprint(object["Status"])))
			}
		}
		return object, fmt.Sprint(object["Status"]), nil
	}
}

func (s *PvtzService) DescribePvtzRule(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	client := s.client
	action := "DescribeResolverRule"
	request := map[string]interface{}{
		"RuleId": id,
	}
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = client.RpcPost("pvtz", "2018-01-01", action, nil, request, true)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)
	if err != nil {
		if IsExpectedErrors(err, []string{"ResolverRule.NotExists"}) {
			return object, WrapErrorf(NotFoundErr("PrivateZone:Rule", id), NotFoundMsg, ProviderERROR)
		}
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}
	v, err := jsonpath.Get("$", response)
	if err != nil {

		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$", response)
	}
	object = v.(map[string]interface{})
	return object, nil
}

func (s *PvtzService) DescribePvtzRuleAttachment(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	client := s.client
	action := "DescribeResolverRule"
	request := map[string]interface{}{
		"RuleId": id,
	}
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = client.RpcPost("pvtz", "2018-01-01", action, nil, request, true)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)
	if err != nil {
		if IsExpectedErrors(err, []string{"ResolverRule.NotExists"}) {
			return object, WrapErrorf(NotFoundErr("PrivateZone:RuleAttachment", id), NotFoundMsg, ProviderERROR)
		}
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}
	v, err := jsonpath.Get("$", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$", response)
	}
	object = v.(map[string]interface{})
	return object, nil
}

func (s *PvtzService) DescribeSyncEcsHostTask(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	client := s.client
	action := "DescribeSyncEcsHostTask"
	request := map[string]interface{}{
		"ZoneId": id,
	}
	response, err = client.RpcPost("pvtz", "2018-01-01", action, nil, request, true)
	if err != nil {
		if IsExpectedErrors(err, []string{"System.Busy", "Ecs.SyncTask.NotExists", "ServiceUnavailable", "Throttling.User"}) {
			err = WrapErrorf(NotFoundErr("PvtzZone", id), NotFoundMsg, ProviderERROR)
			return object, err
		}
		err = WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
		return object, err
	}
	addDebug(action, response, request)
	v, err := jsonpath.Get("$", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$", response)
	}
	object = v.(map[string]interface{})
	return object, nil
}

func (s *PvtzService) DescribePvtzZone(id string) (object map[string]interface{}, err error) {
	object, err = s.DescribePvtzZoneBasic(id)
	if err != nil {
		return nil, WrapError(err)
	}
	syncObj, err := s.DescribeSyncEcsHostTask(id)
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_pvtz_zone pvtzService.DescribeSyncEcsHostTask Failed!!! %s", err)
			return object, nil
		}
		return nil, WrapError(err)
	}
	object["SyncHostTask"] = syncObj
	return object, nil
}

func (s *PvtzService) SetResourceTags(d *schema.ResourceData, resourceType string) error {
	if d.HasChange("tags") {
		var err error
		var action string
		client := s.client
		var request map[string]interface{}
		var response map[string]interface{}

		added, removed := parsingTags(d)
		removedTagKeys := make([]string, 0)
		for _, v := range removed {
			if !ignoredTags(v, "") {
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
			wait := incrementalWait(3*time.Second, 5*time.Second)
			err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
				response, err = client.RpcPost("pvtz", "2018-01-01", action, nil, request, true)

				if err != nil {
					if NeedRetry(err) {
						wait()
						return resource.RetryableError(err)
					}
					return resource.NonRetryableError(err)
				}
				addDebug(action, response, request)
				return nil
			})
			if err != nil {
				return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
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
			wait := incrementalWait(3*time.Second, 5*time.Second)
			err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
				response, err = client.RpcPost("pvtz", "2018-01-01", action, nil, request, true)

				if err != nil {
					if NeedRetry(err) {
						wait()
						return resource.RetryableError(err)
					}
					return resource.NonRetryableError(err)
				}
				addDebug(action, response, request)
				return nil
			})
			if err != nil {
				return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
			}

		}
		d.SetPartial("tags")
	}

	return nil
}

func (s *PvtzService) DescribeListTagResources(id string) (object map[string]interface{}, err error) {

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
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(1*time.Minute, func() *resource.RetryError {
		response, err = client.RpcPost("pvtz", "2018-01-01", action, query, request, false)

		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, response, request)
		return nil
	})
	if err != nil {
		if IsExpectedErrors(err, []string{}) {
			return object, WrapErrorf(NotFoundErr("ZONE", id), NotFoundMsg, ProviderERROR, fmt.Sprint(response["RequestId"]))
		}
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}

	return response, nil
}
