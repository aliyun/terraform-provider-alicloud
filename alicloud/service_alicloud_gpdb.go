package alicloud

import (
	"fmt"
	"log"
	"regexp"
	"strings"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/gpdb"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

type GpdbService struct {
	client *connectivity.AliyunClient
}

func (s *GpdbService) DescribeGpdbInstance(id string) (instanceAttribute gpdb.DBInstance, err error) {
	request := gpdb.CreateDescribeDBInstancesRequest()
	request.RegionId = s.client.RegionId
	request.DBInstanceIds = id
	var response *gpdb.DescribeDBInstancesResponse
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(3*time.Minute, func() *resource.RetryError {
		raw, err := s.client.WithGpdbClient(func(client *gpdb.Client) (interface{}, error) {
			return client.DescribeDBInstances(request)
		})
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(request.GetActionName(), response, request.RpcRequest, request)
		response, _ = raw.(*gpdb.DescribeDBInstancesResponse)
		return nil
	})

	if err != nil {
		// convert error code
		if IsExpectedErrors(err, []string{"InvalidDBInstanceId.NotFound"}) {
			err = WrapErrorf(err, NotFoundMsg, AlibabaCloudSdkGoERROR)
		} else {
			err = WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		return
	}

	if len(response.Items.DBInstance) == 0 {
		return instanceAttribute, WrapErrorf(NotFoundErr("Gpdb Instance", id), NotFoundMsg, ProviderERROR)
	}

	return response.Items.DBInstance[0], nil
}

func (s *GpdbService) DescribeDBInstanceAttribute(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	client := s.client
	action := "DescribeDBInstanceAttribute"
	request := map[string]interface{}{
		"DBInstanceId": id,
	}
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = client.RpcPost("gpdb", "2016-05-03", action, nil, request, true)
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
		if IsExpectedErrors(err, []string{"InvalidDBInstanceId.NotFound"}) {
			return object, WrapErrorf(NotFoundErr("GPDB:DBInstance", id), NotFoundMsg, ProviderERROR, fmt.Sprint(response["RequestId"]))
		}
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}
	v, err := jsonpath.Get("$.Items.DBInstanceAttribute", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.Items.DBInstanceAttribute", response)
	}
	if len(v.([]interface{})) < 1 {
		return object, WrapErrorf(NotFoundErr("GPDB", id), NotFoundWithResponse, response)
	} else {
		if fmt.Sprint(v.([]interface{})[0].(map[string]interface{})["DBInstanceId"]) != id {
			return object, WrapErrorf(NotFoundErr("GPDB", id), NotFoundWithResponse, response)
		}
	}
	object = v.([]interface{})[0].(map[string]interface{})
	return object, nil
}

func (s *GpdbService) DescribeGpdbElasticInstance(id string) (object map[string]interface{}, err error) {
	client := s.client
	action := "DescribeDBInstanceOnECSAttribute"
	request := map[string]interface{}{
		"RegionId":     s.client.RegionId,
		"DBInstanceId": id,
		"SourceIp":     s.client.SourceIp,
	}
	var response map[string]interface{}
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(3*time.Minute, func() *resource.RetryError {
		response, err = client.RpcPost("gpdb", "2016-05-03", action, nil, request, true)
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
		if IsExpectedErrors(err, []string{"InvalidDBInstanceId.NotFound", ServiceUnavailable}) {
			return nil, WrapErrorf(err, NotFoundMsg, AlibabaCloudSdkGoERROR)
		}
		return nil, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}

	v, err := jsonpath.Get("$.Items.DBInstanceAttribute", response)
	if err != nil {
		return nil, WrapErrorf(err, FailedGetAttributeMsg, id, "$.Items.DBInstanceAttribute", response)
	}
	if len(v.([]interface{})) < 1 {
		return nil, WrapErrorf(NotFoundErr("Gpdb elastic instance", id), NotFoundMsg, ProviderERROR)
	}
	return v.([]interface{})[0].(map[string]interface{}), nil
}

func (s *GpdbService) DescribeGpdbSecurityIps(id string) (ips []string, err error) {
	request := gpdb.CreateDescribeDBInstanceIPArrayListRequest()
	request.RegionId = s.client.RegionId
	request.DBInstanceId = id

	raw, err := s.client.WithGpdbClient(func(client *gpdb.Client) (interface{}, error) {
		return client.DescribeDBInstanceIPArrayList(request)
	})
	if err != nil {
		if IsExpectedErrors(err, []string{"InvalidDBInstanceId.NotFound"}) {
			err = WrapErrorf(err, NotFoundMsg, AlibabaCloudSdkGoERROR)
		} else {
			err = WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		return
	}
	response, _ := raw.(*gpdb.DescribeDBInstanceIPArrayListResponse)
	addDebug(request.GetActionName(), response, request.RpcRequest, request)
	var ipstr, separator string
	ipsMap := make(map[string]string)
	for _, ip := range response.Items.DBInstanceIPArray {
		if ip.DBInstanceIPArrayAttribute == "hidden" {
			continue
		}
		ipstr += separator + ip.SecurityIPList
		separator = COMMA_SEPARATED
	}
	for _, ip := range strings.Split(ipstr, COMMA_SEPARATED) {
		if ip == LOCAL_HOST_IP {
			continue
		}
		ipsMap[ip] = ip
	}

	var finalIps []string
	if len(ipsMap) > 0 {
		for key := range ipsMap {
			finalIps = append(finalIps, key)
		}
	}
	return finalIps, nil
}

func (s *GpdbService) ModifyGpdbSecurityIps(id, ips string) error {
	request := gpdb.CreateModifySecurityIpsRequest()
	request.RegionId = s.client.RegionId
	request.DBInstanceId = id
	request.SecurityIPList = ips
	raw, err := s.client.WithGpdbClient(func(client *gpdb.Client) (interface{}, error) {
		return client.ModifySecurityIps(request)
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	response := raw.(*gpdb.ModifySecurityIpsResponse)
	addDebug(request.GetActionName(), response, request.RpcRequest, request)

	return nil
}

func (s *GpdbService) DescribeGpdbConnection(id string) (*gpdb.DBInstanceNetInfo, error) {
	info := &gpdb.DBInstanceNetInfo{}
	parts, err := ParseResourceId(id, 2)
	if err != nil {
		return info, WrapError(err)
	}

	// Describe DB Instance Net Info
	request := gpdb.CreateDescribeDBInstanceNetInfoRequest()
	request.RegionId = s.client.RegionId
	request.DBInstanceId = parts[0]
	raw, err := s.client.WithGpdbClient(func(gpdbClient *gpdb.Client) (interface{}, error) {
		return gpdbClient.DescribeDBInstanceNetInfo(request)
	})
	if err != nil {
		if IsExpectedErrors(err, []string{"InvalidDBInstanceId.NotFound"}) {
			return info, WrapErrorf(err, NotFoundMsg, AlibabaCloudSdkGoERROR)
		}
		return info, WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response, _ := raw.(*gpdb.DescribeDBInstanceNetInfoResponse)
	if response.DBInstanceNetInfos.DBInstanceNetInfo != nil {
		for _, o := range response.DBInstanceNetInfos.DBInstanceNetInfo {
			if strings.HasPrefix(o.ConnectionString, parts[1]) {
				return &o, nil
			}
		}
	}

	return info, WrapErrorf(NotFoundErr("GpdbConnection", id), NotFoundMsg, ProviderERROR)
}

func (s *GpdbService) GpdbInstanceStateRefreshFunc(id string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeGpdbInstance(id)
		if err != nil {
			if NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", WrapError(err)
		}

		for _, failState := range failStates {
			if object.DBInstanceStatus == failState {
				return object, object.DBInstanceStatus, WrapError(Error(FailedToReachTargetStatus, object.DBInstanceStatus))
			}
		}
		return object, object.DBInstanceStatus, nil
	}
}

func (s *GpdbService) GpdbElasticInstanceStateRefreshFunc(id string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		instance, err := s.DescribeGpdbElasticInstance(id)
		if err != nil {
			if NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", WrapError(err)
		}

		for _, failState := range failStates {
			if instance["DBInstanceStatus"] == failState {
				return instance, instance["DBInstanceStatus"].(string), WrapError(Error(FailedToReachTargetStatus, instance["DBInstanceStatus"]))
			}
		}
		return instance, instance["DBInstanceStatus"].(string), nil
	}
}

func (s *GpdbService) WaitForGpdbConnection(id string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	for {
		object, err := s.DescribeGpdbConnection(id)
		if err != nil {
			if NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			} else {
				return WrapError(err)
			}
		}
		if object.ConnectionString != "" && status != Deleted {
			return nil
		}
		if time.Now().After(deadline) {
			return WrapErrorf(err, WaitTimeoutMsg, id, GetFunc(1), timeout, object.ConnectionString, id, ProviderERROR)
		}
	}
}

func (s *GpdbService) setInstanceTags(d *schema.ResourceData) error {
	oraw, nraw := d.GetChange("tags")
	o := oraw.(map[string]interface{})
	n := nraw.(map[string]interface{})
	create, remove := diffGpdbTags(gpdbTagsFromMap(o), gpdbTagsFromMap(n))

	if len(remove) > 0 {
		var tagKey []string
		for _, v := range remove {
			tagKey = append(tagKey, v.Key)
		}
		request := gpdb.CreateUntagResourcesRequest()
		request.ResourceId = &[]string{d.Id()}
		request.ResourceType = string(TagResourceInstance)
		request.TagKey = &tagKey
		request.RegionId = s.client.RegionId
		raw, err := s.client.WithGpdbClient(func(client *gpdb.Client) (interface{}, error) {
			return client.UntagResources(request)
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	}

	if len(create) > 0 {
		request := gpdb.CreateTagResourcesRequest()
		request.ResourceId = &[]string{d.Id()}
		request.Tag = &create
		request.ResourceType = string(TagResourceInstance)
		request.RegionId = s.client.RegionId
		raw, err := s.client.WithGpdbClient(func(client *gpdb.Client) (interface{}, error) {
			return client.TagResources(request)
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	}

	d.SetPartial("tags")
	return nil
}

func (s *GpdbService) tagsToMap(tags []gpdb.Tag) map[string]string {
	result := make(map[string]string)
	for _, t := range tags {
		if !s.ignoreTag(t) {
			result[t.Key] = t.Value
		}
	}
	return result
}

func (s *GpdbService) ignoreTag(t gpdb.Tag) bool {
	filter := []string{"^aliyun", "^acs:", "^http://", "^https://"}
	for _, v := range filter {
		log.Printf("[DEBUG] Matching prefix %v with %v\n", v, t.Key)
		ok, _ := regexp.MatchString(v, t.Key)
		if ok {
			log.Printf("[DEBUG] Found Alibaba Cloud specific t %s (val: %s), ignoring.\n", t.Key, t.Value)
			return true
		}
	}
	return false
}

func (s *GpdbService) DescribeGpdbAccount(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	client := s.client
	action := "DescribeAccounts"
	parts, err := ParseResourceId(id, 2)
	if err != nil {
		err = WrapError(err)
		return
	}
	request := map[string]interface{}{
		"AccountName":  parts[1],
		"DBInstanceId": parts[0],
	}
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {

		response, err = client.RpcPost("gpdb", "2016-05-03", action, nil, request, true)
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
		if IsExpectedErrors(err, []string{"InvalidDBInstanceId.NotFound"}) {
			return object, WrapErrorf(NotFoundErr("GPDB:Account", id), NotFoundMsg, ProviderERROR, fmt.Sprint(response["RequestId"]))
		}
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}
	v, err := jsonpath.Get("$.Accounts.DBInstanceAccount", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.Accounts.DBInstanceAccount", response)
	}
	if len(v.([]interface{})) < 1 {
		return object, WrapErrorf(NotFoundErr("GPDB", id), NotFoundWithResponse, response)
	} else {
		if fmt.Sprint(v.([]interface{})[0].(map[string]interface{})["AccountName"]) != parts[1] {
			return object, WrapErrorf(NotFoundErr("GPDB", id), NotFoundWithResponse, response)
		}
	}
	object = v.([]interface{})[0].(map[string]interface{})
	return object, nil
}

func (s *GpdbService) GpdbAccountStateRefreshFunc(id string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeGpdbAccount(id)
		if err != nil {
			if NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", WrapError(err)
		}

		for _, failState := range failStates {
			if fmt.Sprint(object["AccountStatus"]) == failState {
				return object, fmt.Sprint(object["AccountStatus"]), WrapError(Error(FailedToReachTargetStatus, fmt.Sprint(object["AccountStatus"])))
			}
		}
		return object, fmt.Sprint(object["AccountStatus"]), nil
	}
}

func (s *GpdbService) ListTagResources(id string, resourceType string) (object interface{}, err error) {
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
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			response, err := client.RpcPost("gpdb", "2016-05-03", action, nil, request, false)
			if err != nil {
				if IsExpectedErrors(err, []string{Throttling}) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			addDebug(action, response, request)
			v, err := jsonpath.Get("$.TagResources", response)
			if err != nil {
				return resource.NonRetryableError(WrapErrorf(err, FailedGetAttributeMsg, id, "$.TagResources", response))
			}
			if v != nil {
				tags = append(tags, v.([]interface{})...)
			}
			return nil
		})
		if err != nil {
			err = WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
			return
		}
		if response["NextToken"] == nil {
			break
		}
		request["NextToken"] = response["NextToken"]
	}

	return tags, nil
}

func (s *GpdbService) SetResourceTags(d *schema.ResourceData, resourceType string) error {

	if d.HasChange("tags") {
		added, removed := parsingTags(d)
		client := s.client

		removedTagKeys := make([]string, 0)
		for _, v := range removed {
			if !ignoredTags(v, "") {
				removedTagKeys = append(removedTagKeys, v)
			}
		}
		if len(removedTagKeys) > 0 {
			action := "UntagResources"
			request := map[string]interface{}{
				"RegionId":     s.client.RegionId,
				"ResourceType": resourceType,
				"ResourceId.1": d.Id(),
			}
			for i, key := range removedTagKeys {
				request[fmt.Sprintf("TagKey.%d", i+1)] = key
			}
			wait := incrementalWait(2*time.Second, 1*time.Second)
			err := resource.Retry(10*time.Minute, func() *resource.RetryError {
				response, err := client.RpcPost("gpdb", "2016-05-03", action, nil, request, false)
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

			wait := incrementalWait(2*time.Second, 1*time.Second)
			err := resource.Retry(10*time.Minute, func() *resource.RetryError {
				response, err := client.RpcPost("gpdb", "2016-05-03", action, nil, request, false)
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

func (s *GpdbService) DescribeDBInstanceIPArrayList(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	client := s.client
	action := "DescribeDBInstanceIPArrayList"
	request := map[string]interface{}{
		"DBInstanceId": id,
	}
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = client.RpcPost("gpdb", "2016-05-03", action, nil, request, true)
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
		if IsExpectedErrors(err, []string{"InvalidDBInstanceId.NotFound"}) {
			return object, WrapErrorf(NotFoundErr("GPDB:DBInstance", id), NotFoundMsg, ProviderERROR, fmt.Sprint(response["RequestId"]))
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

func (s *GpdbService) DescribeDBInstanceSSL(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	action := "DescribeDBInstanceSSL"
	request := map[string]interface{}{
		"RegionId":     s.client.RegionId,
		"DBInstanceId": id,
	}
	client := s.client
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = client.RpcPost("gpdb", "2016-05-03", action, nil, request, true)
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
		if IsExpectedErrors(err, []string{"InvalidDBInstanceId.NotFound"}) {
			return object, WrapErrorf(NotFoundErr("GPDB:DBInstance", id), NotFoundMsg, ProviderERROR, fmt.Sprint(response["RequestId"]))
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

func (s *GpdbService) DescribeDBResourceManagementMode(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	action := "DescribeDBResourceManagementMode"
	request := map[string]interface{}{
		"RegionId":     s.client.RegionId,
		"DBInstanceId": id,
	}
	client := s.client
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = client.RpcPost("gpdb", "2016-05-03", action, nil, request, true)
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
		if IsExpectedErrors(err, []string{"InvalidDBInstanceId.NotFound"}) {
			return object, WrapErrorf(NotFoundErr("GPDB:DBInstance", id), NotFoundMsg, ProviderERROR, fmt.Sprint(response["RequestId"]))
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

func (s *GpdbService) DescribeParameters(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	action := "DescribeParameters"
	request := map[string]interface{}{
		"DBInstanceId": id,
	}
	client := s.client
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = client.RpcPost("gpdb", "2016-05-03", action, nil, request, true)
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
		if IsExpectedErrors(err, []string{"InvalidDBInstanceId.NotFound"}) {
			return object, WrapErrorf(NotFoundErr("GPDB:DBInstance", id), NotFoundMsg, ProviderERROR, fmt.Sprint(response["RequestId"]))
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

func (s *GpdbService) DBInstanceSSLStateRefreshFunc(d *schema.ResourceData, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeDBInstanceSSL(d.Id())
		if err != nil {
			if NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", WrapError(err)
		}

		sslEnabled := convertGpdbDbInstanceSSLEnabledResponse(object["SSLEnabled"])
		for _, failState := range failStates {
			if fmt.Sprint(sslEnabled) == failState {
				return object, fmt.Sprint(sslEnabled), WrapError(Error(FailedToReachTargetStatus, fmt.Sprint(sslEnabled)))
			}
		}

		return object, fmt.Sprint(sslEnabled), nil
	}
}

func (s *GpdbService) DescribeGpdbDbInstance(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	client := s.client
	action := "DescribeDBInstanceAttribute"
	request := map[string]interface{}{
		"DBInstanceId": id,
	}
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = client.RpcPost("gpdb", "2016-05-03", action, nil, request, true)
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
		if IsExpectedErrors(err, []string{"InvalidDBInstanceId.NotFound"}) {
			return object, WrapErrorf(NotFoundErr("GPDB:DBInstance", id), NotFoundMsg, ProviderERROR, fmt.Sprint(response["RequestId"]))
		}
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}
	v, err := jsonpath.Get("$.Items.DBInstanceAttribute", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.Items.DBInstanceAttribute", response)
	}
	if len(v.([]interface{})) < 1 {
		return object, WrapErrorf(NotFoundErr("GPDB", id), NotFoundWithResponse, response)
	} else {
		if fmt.Sprint(v.([]interface{})[0].(map[string]interface{})["DBInstanceId"]) != id {
			return object, WrapErrorf(NotFoundErr("GPDB", id), NotFoundWithResponse, response)
		}
	}
	object = v.([]interface{})[0].(map[string]interface{})
	return object, nil
}

func (s *GpdbService) GpdbDbInstanceStateRefreshFunc(id string, field string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeGpdbDbInstance(id)
		if err != nil {
			if NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", WrapError(err)
		}

		for _, failState := range failStates {
			if fmt.Sprint(object[field]) == failState {
				return object, fmt.Sprint(object[field]), WrapError(Error(FailedToReachTargetStatus, fmt.Sprint(object[field])))
			}
		}
		return object, fmt.Sprint(object[field]), nil
	}
}

func (s *GpdbService) DescribeGpdbDbInstancePlan(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	action := "DescribeDBInstancePlans"
	client := s.client
	parts, err := ParseResourceId(id, 2)
	if err != nil {
		return object, WrapError(err)
	}

	request := map[string]interface{}{
		"DBInstanceId": parts[0],
		"PlanId":       parts[1],
	}

	idExist := false
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = client.RpcPost("gpdb", "2016-05-03", action, nil, request, true)
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
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}

	resp, err := jsonpath.Get("$.Items.PlanList", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.Items.PlanList", response)
	}

	if v, ok := resp.([]interface{}); !ok || len(v) < 1 {
		return object, WrapErrorf(NotFoundErr("Gpdb:DbInstancePlan", id), NotFoundWithResponse, response)
	}

	for _, v := range resp.([]interface{}) {
		if fmt.Sprint(v.(map[string]interface{})["DBInstanceId"]) == parts[0] && fmt.Sprint(v.(map[string]interface{})["PlanId"]) == parts[1] {
			idExist = true
			return v.(map[string]interface{}), nil
		}
	}

	if !idExist {
		return object, WrapErrorf(NotFoundErr("Gpdb:DbInstancePlan", id), NotFoundWithResponse, response)
	}

	return object, nil
}

func (s *GpdbService) GpdbDbInstancePlanStateRefreshFunc(id string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeGpdbDbInstancePlan(id)
		if err != nil {
			if NotFoundError(err) {
				return nil, "", nil
			}
			return nil, "", WrapError(err)
		}

		for _, failState := range failStates {
			if fmt.Sprint(object["PlanStatus"]) == failState {
				return object, fmt.Sprint(object["PlanStatus"]), WrapError(Error(FailedToReachTargetStatus, fmt.Sprint(object["PlanStatus"])))
			}
		}

		return object, fmt.Sprint(object["PlanStatus"]), nil
	}
}

func (s *GpdbService) DescribeGpdbDbInstanceDataShareStatus(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	action := "DescribeDataShareInstances"

	client := s.client

	request := map[string]interface{}{
		"RegionId":    s.client.RegionId,
		"SearchValue": id,
		"PageSize":    PageSizeLarge,
		"PageNumber":  1,
	}

	for {
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			response, err = client.RpcPost("gpdb", "2016-05-03", action, nil, request, true)
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
			return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
		}

		resp, err := jsonpath.Get("$.Items.DBInstance", response)
		if err != nil {
			return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.Items.DBInstance", response)
		}

		if v, ok := resp.([]interface{}); !ok || len(v) < 1 {
			return object, nil
		}

		for _, v := range resp.([]interface{}) {
			if fmt.Sprint(v.(map[string]interface{})["DBInstanceId"]) == id {
				return v.(map[string]interface{}), nil
			}
		}

		if len(resp.([]interface{})) < request["PageSize"].(int) {
			break
		}

		request["PageNumber"] = request["PageNumber"].(int) + 1
	}

	return object, nil
}

func (s *GpdbService) GpdbDbInstanceDataShareStatusStateRefreshFunc(id string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeGpdbDbInstanceDataShareStatus(id)
		if err != nil {
			if NotFoundError(err) {
				return nil, "", nil
			}
			return nil, "", WrapError(err)
		}

		for _, failState := range failStates {
			if fmt.Sprint(object["DataShareStatus"]) == failState {
				return object, fmt.Sprint(object["DataShareStatus"]), WrapError(Error(FailedToReachTargetStatus, fmt.Sprint(object["DataShareStatus"])))
			}
		}

		return object, fmt.Sprint(object["DataShareStatus"]), nil
	}
}
