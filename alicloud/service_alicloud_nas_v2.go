package alicloud

import (
	"fmt"
	"strings"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/blues/jsonata-go"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

type NasServiceV2 struct {
	client *connectivity.AliyunClient
}

// DescribeNasAccessRule <<< Encapsulated get interface for Nas AccessRule.

func (s *NasServiceV2) DescribeNasAccessRule(id string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	parts := strings.Split(id, ":")
	action := "DescribeAccessRules"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	query["AccessGroupName"] = parts[0]
	query["AccessRuleId"] = parts[1]
	if len(parts) == 3 {
		query["FileSystemType"] = parts[2]
	}

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(1*time.Minute, func() *resource.RetryError {
		response, err = client.RpcPost("NAS", "2017-06-26", action, query, request, true)

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
		if IsExpectedErrors(err, []string{"InvalidAccessRule.NotFound"}) {
			return object, WrapErrorf(NotFoundErr("AccessRule", id), NotFoundMsg, response)
		}
		addDebug(action, response, request)
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}

	v, err := jsonpath.Get("$.AccessRules.AccessRule[*]", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.AccessRules.AccessRule[*]", response)
	}

	if len(v.([]interface{})) == 0 {
		return object, WrapErrorf(NotFoundErr("AccessRule", id), NotFoundMsg, response)
	}

	return v.([]interface{})[0].(map[string]interface{}), nil
}

func (s *NasServiceV2) NasAccessRuleStateRefreshFunc(id string, field string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeNasAccessRule(id)
		if err != nil {
			if NotFoundError(err) {
				return object, "", nil
			}
			return nil, "", WrapError(err)
		}

		v, err := jsonpath.Get(field, object)
		currentStatus := fmt.Sprint(v)

		for _, failState := range failStates {
			if currentStatus == failState {
				return object, currentStatus, WrapError(Error(FailedToReachTargetStatus, currentStatus))
			}
		}
		return object, currentStatus, nil
	}
}

// DescribeNasAccessRule >>> Encapsulated.

// DescribeNasAccessGroup <<< Encapsulated get interface for Nas AccessGroup.

func (s *NasServiceV2) DescribeNasAccessGroup(id string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	parts := strings.Split(id, ":")
	if len(parts) != 2 {
		err = WrapError(fmt.Errorf("invalid Resource Id %s. Expected parts' length %d, got %d", id, 2, len(parts)))
	}
	action := "DescribeAccessGroups"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	query["AccessGroupName"] = parts[0]
	query["FileSystemType"] = parts[1]

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(1*time.Minute, func() *resource.RetryError {
		response, err = client.RpcPost("NAS", "2017-06-26", action, query, request, true)

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
		if IsExpectedErrors(err, []string{"InvalidAccessGroup.NotFound", "Resource.NotFound"}) {
			return object, WrapErrorf(NotFoundErr("AccessGroup", id), NotFoundMsg, response)
		}
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}

	v, err := jsonpath.Get("$.AccessGroups.AccessGroup[*]", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.AccessGroups.AccessGroup[*]", response)
	}

	if len(v.([]interface{})) == 0 {
		return object, WrapErrorf(NotFoundErr("AccessGroup", id), NotFoundMsg, response)
	}

	return v.([]interface{})[0].(map[string]interface{}), nil
}

func (s *NasServiceV2) NasAccessGroupStateRefreshFunc(id string, field string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeNasAccessGroup(id)
		if err != nil {
			if NotFoundError(err) {
				return object, "", nil
			}
			return nil, "", WrapError(err)
		}

		v, err := jsonpath.Get(field, object)
		currentStatus := fmt.Sprint(v)

		for _, failState := range failStates {
			if currentStatus == failState {
				return object, currentStatus, WrapError(Error(FailedToReachTargetStatus, currentStatus))
			}
		}
		return object, currentStatus, nil
	}
}

// DescribeNasAccessGroup >>> Encapsulated.

// DescribeNasFileSystem <<< Encapsulated get interface for Nas FileSystem.

func (s *NasServiceV2) DescribeNasFileSystem(id string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["FileSystemId"] = id

	action := "DescribeFileSystems"

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(1*time.Minute, func() *resource.RetryError {
		response, err = client.RpcPost("NAS", "2017-06-26", action, query, request, true)

		if err != nil {
			if IsExpectedErrors(err, []string{"InvalidFileSystemStatus.Ordering"}) || NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)
	if err != nil {
		if IsExpectedErrors(err, []string{"Resource.NotFound", "InvalidFileSystem.NotFound"}) {
			return object, WrapErrorf(NotFoundErr("FileSystem", id), NotFoundMsg, response)
		}
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}

	v, err := jsonpath.Get("$.FileSystems.FileSystem[*]", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.FileSystems.FileSystem[*]", response)
	}

	if len(v.([]interface{})) == 0 {
		return object, WrapErrorf(NotFoundErr("FileSystem", id), NotFoundMsg, response)
	}

	return v.([]interface{})[0].(map[string]interface{}), nil
}
func (s *NasServiceV2) DescribeFileSystemDescribeSmbAcl(id string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["FileSystemId"] = id

	action := "DescribeSmbAcl"

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(1*time.Minute, func() *resource.RetryError {
		response, err = client.RpcPost("NAS", "2017-06-26", action, query, request, true)

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

	v, err := jsonpath.Get("$.Acl", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.Acl", response)
	}

	return v.(map[string]interface{}), nil
}
func (s *NasServiceV2) DescribeFileSystemGetRecycleBinAttribute(id string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	query["FileSystemId"] = id

	action := "GetRecycleBinAttribute"

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(1*time.Minute, func() *resource.RetryError {
		response, err = client.RpcGet("NAS", "2017-06-26", action, query, request)

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
		if IsExpectedErrors(err, []string{"InvalidFileSystem.NotFound"}) {
			return object, WrapErrorf(NotFoundErr("FileSystem", id), NotFoundMsg, response)
		}
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}

	v, err := jsonpath.Get("$.RecycleBinAttribute", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.RecycleBinAttribute", response)
	}

	return v.(map[string]interface{}), nil
}
func (s *NasServiceV2) DescribeFileSystemListTagResources(id string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["ResourceId.1"] = id

	request["ResourceType"] = "filesystem"
	action := "ListTagResources"

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(1*time.Minute, func() *resource.RetryError {
		response, err = client.RpcPost("NAS", "2017-06-26", action, query, request, true)

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
		if IsExpectedErrors(err, []string{"InvalidResourceId.NotFound"}) {
			return object, WrapErrorf(NotFoundErr("FileSystem", id), NotFoundMsg, response)
		}
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}

	return response, nil
}
func (s *NasServiceV2) DescribeFileSystemDescribeNfsAcl(id string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["FileSystemId"] = id

	action := "DescribeNfsAcl"

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(1*time.Minute, func() *resource.RetryError {
		response, err = client.RpcPost("NAS", "2017-06-26", action, query, request, true)

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

	v, err := jsonpath.Get("$.Acl", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.Acl", response)
	}

	return v.(map[string]interface{}), nil
}

func (s *NasServiceV2) NasFileSystemStateRefreshFunc(id string, field string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeNasFileSystem(id)
		if err != nil {
			if NotFoundError(err) {
				return object, "", nil
			}
			return nil, "", WrapError(err)
		}

		v, err := jsonpath.Get(field, object)
		currentStatus := fmt.Sprint(v)

		if strings.HasPrefix(field, "#") {
			v, _ := jsonpath.Get(strings.TrimPrefix(field, "#"), object)
			if v != nil {
				currentStatus = "#CHECKSET"
			}
		}

		for _, failState := range failStates {
			if currentStatus == failState {
				return object, currentStatus, WrapError(Error(FailedToReachTargetStatus, currentStatus))
			}
		}
		return object, currentStatus, nil
	}
}

func (s *NasServiceV2) DescribeAsyncNasFileSystemStateRefreshFunc(d *schema.ResourceData, res map[string]interface{}, field string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeAsyncDescribeFileSystems(d, res)
		if err != nil {
			if NotFoundError(err) {
				return object, "", nil
			}
		}

		v, err := jsonpath.Get(field, object)
		currentStatus := fmt.Sprint(v)

		if strings.HasPrefix(field, "#") {
			v, _ := jsonpath.Get(strings.TrimPrefix(field, "#"), object)
			if v != nil {
				currentStatus = "#CHECKSET"
			}
		}

		for _, failState := range failStates {
			if currentStatus == failState {
				if _err, ok := object["error"]; ok {
					return _err, currentStatus, WrapError(Error(FailedToReachTargetStatus, currentStatus))
				}
				return object, currentStatus, WrapError(Error(FailedToReachTargetStatus, currentStatus))
			}
		}
		return object, currentStatus, nil
	}
}

// DescribeNasFileSystem >>> Encapsulated.

// SetResourceTags <<< Encapsulated tag function for Nas.
func (s *NasServiceV2) SetResourceTags(d *schema.ResourceData, resourceType string) error {
	if d.HasChange("tags") {
		var action string
		var err error
		client := s.client
		var request map[string]interface{}
		var response map[string]interface{}
		query := make(map[string]interface{})

		added, removed := parsingTags(d)
		removedTagKeys := make([]string, 0)
		for _, v := range removed {
			if !ignoredTags(v, "") {
				removedTagKeys = append(removedTagKeys, v)
			}
		}
		if len(removedTagKeys) > 0 {
			action = "UntagResources"
			request = make(map[string]interface{})
			query = make(map[string]interface{})
			request["ResourceId.1"] = d.Id()

			for i, key := range removedTagKeys {
				request[fmt.Sprintf("TagKey.%d", i+1)] = key
			}

			request["ResourceType"] = resourceType
			wait := incrementalWait(3*time.Second, 5*time.Second)
			err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
				response, err = client.RpcPost("NAS", "2017-06-26", action, query, request, true)
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
				return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
			}

		}

		if len(added) > 0 {
			action = "TagResources"
			request = make(map[string]interface{})
			query = make(map[string]interface{})
			request["ResourceId.1"] = d.Id()

			count := 1
			for key, value := range added {
				request[fmt.Sprintf("Tag.%d.Key", count)] = key
				request[fmt.Sprintf("Tag.%d.Value", count)] = value
				count++
			}

			request["ResourceType"] = resourceType
			wait := incrementalWait(3*time.Second, 5*time.Second)
			err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
				response, err = client.RpcPost("NAS", "2017-06-26", action, query, request, true)
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
				return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
			}

		}
	}

	return nil
}

// SetResourceTags >>> tag function encapsulated.

// DescribeNasAutoSnapshotPolicy <<< Encapsulated get interface for Nas AutoSnapshotPolicy.

func (s *NasServiceV2) DescribeNasAutoSnapshotPolicy(id string) (object map[string]interface{}, err error) {
	client := s.client
	var response map[string]interface{}
	var query map[string]interface{}
	action := "DescribeAutoSnapshotPolicies"
	query = make(map[string]interface{})
	query["AutoSnapshotPolicyId"] = id

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(1*time.Minute, func() *resource.RetryError {
		response, err = client.RpcPost("NAS", "2017-06-26", action, query, nil, true)

		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, response, query)
		return nil
	})
	if err != nil {
		if IsExpectedErrors(err, []string{"InvalidLifecyclePolicy.NotFound"}) {
			return object, WrapErrorf(NotFoundErr("AutoSnapshotPolicy", id), NotFoundMsg, response)
		}
		addDebug(action, response, query)
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}

	v, err := jsonpath.Get("$.AutoSnapshotPolicies.AutoSnapshotPolicy[*]", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.AutoSnapshotPolicies.AutoSnapshotPolicy[*]", response)
	}

	if len(v.([]interface{})) == 0 {
		return object, WrapErrorf(NotFoundErr("AutoSnapshotPolicy", id), NotFoundMsg, response)
	}

	return v.([]interface{})[0].(map[string]interface{}), nil
}

func (s *NasServiceV2) NasAutoSnapshotPolicyStateRefreshFunc(id string, field string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeNasAutoSnapshotPolicy(id)
		if err != nil {
			if NotFoundError(err) {
				return object, "", nil
			}
			return nil, "", WrapError(err)
		}

		v, err := jsonpath.Get(field, object)
		currentStatus := fmt.Sprint(v)
		if field == "$.RepeatWeekdays" {
			e := jsonata.MustCompile("$split($.RepeatWeekdays, \",\")")
			v, _ = e.Eval(object)
			currentStatus = fmt.Sprint(v)
		}
		if field == "$.TimePoints" {
			e := jsonata.MustCompile("$split($.TimePoints, \",\")")
			v, _ = e.Eval(object)
			currentStatus = fmt.Sprint(v)
		}

		for _, failState := range failStates {
			if currentStatus == failState {
				return object, currentStatus, WrapError(Error(FailedToReachTargetStatus, currentStatus))
			}
		}
		return object, currentStatus, nil
	}
}

// DescribeNasAutoSnapshotPolicy >>> Encapsulated.

// DescribeNasAccessPoint <<< Encapsulated get interface for Nas AccessPoint.

func (s *NasServiceV2) DescribeNasAccessPoint(id string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	parts := strings.Split(id, ":")
	if len(parts) != 2 {
		err = WrapError(fmt.Errorf("invalid Resource Id %s. Expected parts' length %d, got %d", id, 2, len(parts)))
	}
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["AccessPointId"] = parts[1]
	request["FileSystemId"] = parts[0]

	action := "DescribeAccessPoint"

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(1*time.Minute, func() *resource.RetryError {
		response, err = client.RpcPost("NAS", "2017-06-26", action, query, request, true)

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
		if IsExpectedErrors(err, []string{"NotFound"}) {
			return object, WrapErrorf(NotFoundErr("AccessPoint", id), NotFoundMsg, response)
		}
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}

	v, err := jsonpath.Get("$.AccessPoint", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.AccessPoint", response)
	}

	return v.(map[string]interface{}), nil
}

func (s *NasServiceV2) NasAccessPointStateRefreshFunc(id string, field string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeNasAccessPoint(id)
		if err != nil {
			if NotFoundError(err) {
				return nil, "", nil
			}
			return nil, "", WrapError(err)
		}

		v, err := jsonpath.Get(field, object)
		currentStatus := fmt.Sprint(v)

		for _, failState := range failStates {
			if currentStatus == failState {
				return object, currentStatus, WrapError(Error(FailedToReachTargetStatus, currentStatus))
			}
		}
		return object, currentStatus, nil
	}
}

// DescribeNasAccessPoint >>> Encapsulated.

// DescribeNasMountTarget <<< Encapsulated get interface for Nas MountTarget.

func (s *NasServiceV2) DescribeNasMountTarget(id string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	parts := strings.Split(id, ":")
	if len(parts) != 2 {
		err = WrapError(fmt.Errorf("invalid Resource Id %s. Expected parts' length %d, got %d", id, 2, len(parts)))
	}
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["FileSystemId"] = parts[0]
	request["MountTargetDomain"] = parts[1]

	action := "DescribeMountTargets"

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(1*time.Minute, func() *resource.RetryError {
		response, err = client.RpcPost("NAS", "2017-06-26", action, query, request, true)

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
		if IsExpectedErrors(err, []string{"InvalidMountTarget.NotFound", "InvalidFileSystem.NotFound", "Forbidden.NasNotFound", "InvalidLBid.NotFound", "VolumeUnavailable", "InvalidParam.MountTargetDomain"}) {
			return object, WrapErrorf(NotFoundErr("MountTarget", id), NotFoundMsg, response)
		}
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}

	v, err := jsonpath.Get("$.MountTargets.MountTarget[*]", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.MountTargets.MountTarget[*]", response)
	}

	if len(v.([]interface{})) == 0 {
		return object, WrapErrorf(NotFoundErr("MountTarget", id), NotFoundMsg, response)
	}

	return v.([]interface{})[0].(map[string]interface{}), nil
}

func (s *NasServiceV2) NasMountTargetStateRefreshFunc(id string, field string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeNasMountTarget(id)
		if err != nil {
			if NotFoundError(err) {
				return object, "", nil
			}
			return nil, "", WrapError(err)
		}

		v, err := jsonpath.Get(field, object)
		currentStatus := fmt.Sprint(v)

		if strings.HasPrefix(field, "#") {
			v, _ := jsonpath.Get(strings.TrimPrefix(field, "#"), object)
			if v != nil {
				currentStatus = "#CHECKSET"
			}
		}

		for _, failState := range failStates {
			if currentStatus == failState {
				return object, currentStatus, WrapError(Error(FailedToReachTargetStatus, currentStatus))
			}
		}
		return object, currentStatus, nil
	}
}

func (s *NasServiceV2) DescribeAsyncNasMountTargetStateRefreshFunc(d *schema.ResourceData, res map[string]interface{}, field string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeAsyncDescribeMountTargets(d, res)

		if err != nil {
			if NotFoundError(err) {
				return object, "", nil
			}
		}

		v, err := jsonpath.Get(field, object)
		currentStatus := fmt.Sprint(v)

		if strings.HasPrefix(field, "#") {
			v, _ := jsonpath.Get(strings.TrimPrefix(field, "#"), object)
			if v != nil {
				currentStatus = "#CHECKSET"
			}
		}

		for _, failState := range failStates {
			if currentStatus == failState {
				if _err, ok := object["error"]; ok {
					return _err, currentStatus, WrapError(Error(FailedToReachTargetStatus, currentStatus))
				}
				return object, currentStatus, WrapError(Error(FailedToReachTargetStatus, currentStatus))
			}
		}
		return object, currentStatus, nil
	}
}

// DescribeNasMountTarget >>> Encapsulated.

// DescribeAsyncDescribeMountTargets <<< Encapsulated for Nas.
func (s *NasServiceV2) DescribeAsyncDescribeMountTargets(d *schema.ResourceData, res map[string]interface{}) (object map[string]interface{}, err error) {
	client := s.client
	id := d.Id()
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	parts := strings.Split(id, ":")
	if len(parts) != 2 {
		err = WrapError(fmt.Errorf("invalid Resource Id %s. Expected parts' length %d, got %d", id, 2, len(parts)))
	}
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["FileSystemId"] = parts[0]

	action := "DescribeMountTargets"

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(1*time.Minute, func() *resource.RetryError {
		response, err = client.RpcPost("NAS", "2017-06-26", action, query, request, true)

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
		return response, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}

	return response, nil
}

// DescribeAsyncDescribeMountTargets >>> Encapsulated.

// DescribeAsyncDescribeFileSystems <<< Encapsulated for Nas.
func (s *NasServiceV2) DescribeAsyncDescribeFileSystems(d *schema.ResourceData, res map[string]interface{}) (object map[string]interface{}, err error) {
	client := s.client
	id := d.Id()
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["FileSystemId"] = d.Id()

	action := "DescribeFileSystems"

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(1*time.Minute, func() *resource.RetryError {
		response, err = client.RpcPost("NAS", "2017-06-26", action, query, request, true)

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
		return response, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}

	return response, nil
}

// DescribeAsyncDescribeFileSystems >>> Encapsulated.
