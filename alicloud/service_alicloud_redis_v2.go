package alicloud

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

type RedisServiceV2 struct {
	client *connectivity.AliyunClient
}

// DescribeRedisTairInstance <<< Encapsulated get interface for Redis TairInstance.

func (s *RedisServiceV2) DescribeRedisTairInstance(id string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	action := "DescribeInstanceAttribute"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["InstanceId"] = id
	request["RegionId"] = client.RegionId

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(1*time.Minute, func() *resource.RetryError {
		response, err = client.RpcPost("R-kvstore", "2015-01-01", action, query, request, true)

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
		if IsExpectedErrors(err, []string{"InvalidInstanceId.NotFound"}) {
			return object, WrapErrorf(NotFoundErr("TairInstance", id), NotFoundMsg, response)
		}
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}

	v, err := jsonpath.Get("$.Instances.DBInstanceAttribute[*]", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.Instances.DBInstanceAttribute[*]", response)
	}

	if len(v.([]interface{})) == 0 {
		return object, WrapErrorf(NotFoundErr("TairInstance", id), NotFoundMsg, response)
	}

	return v.([]interface{})[0].(map[string]interface{}), nil
}
func (s *RedisServiceV2) DescribeTairInstanceDescribeInstanceConfig(id string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	action := "DescribeInstanceConfig"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["InstanceId"] = id
	request["RegionId"] = client.RegionId

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(1*time.Minute, func() *resource.RetryError {
		response, err = client.RpcPost("R-kvstore", "2015-01-01", action, query, request, true)

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
		if IsExpectedErrors(err, []string{"InvalidInstanceId.NotFound"}) {
			return object, WrapErrorf(NotFoundErr("TairInstance", id), NotFoundMsg, response)
		}
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}

	return response, nil
}
func (s *RedisServiceV2) DescribeDescribeSecurityIps(id, securityIpGroupName string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	action := "DescribeSecurityIps"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["InstanceId"] = id
	request["RegionId"] = client.RegionId

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(1*time.Minute, func() *resource.RetryError {
		response, err = client.RpcPost("R-kvstore", "2015-01-01", action, query, request, true)

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
		if IsExpectedErrors(err, []string{"InvalidInstanceId.NotFound"}) {
			return object, WrapErrorf(NotFoundErr("TairInstance", id), NotFoundMsg, response)
		}
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}

	v, err := jsonpath.Get("$.SecurityIpGroups.SecurityIpGroup[*]", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.SecurityIpGroups.SecurityIpGroup[*]", response)
	}

	if securityIpGroupName == "" {
		securityIpGroupName = "default"
	}

	for _, item := range v.([]interface{}) {
		securityIpGroup := item.(map[string]interface{})
		if securityIpGroup["SecurityIpGroupName"] == securityIpGroupName {
			return securityIpGroup, nil
		}
	}

	if len(v.([]interface{})) == 0 {
		return object, WrapErrorf(NotFoundErr("TairInstance", id), NotFoundMsg, response)
	}

	return v.([]interface{})[0].(map[string]interface{}), nil
}
func (s *RedisServiceV2) DescribeTairInstanceDescribeSecurityGroupConfiguration(id string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	action := "DescribeSecurityGroupConfiguration"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["InstanceId"] = id
	request["RegionId"] = client.RegionId

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(1*time.Minute, func() *resource.RetryError {
		response, err = client.RpcPost("R-kvstore", "2015-01-01", action, query, request, true)

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
		if IsExpectedErrors(err, []string{"IncorrectDBInstanceState", "InvalidInstanceId.NotFound"}) {
			return object, WrapErrorf(NotFoundErr("TairInstance", id), NotFoundMsg, response)
		}
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}

	v, err := jsonpath.Get("$.Items.EcsSecurityGroupRelation[*]", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.Items.EcsSecurityGroupRelation[*]", response)
	}

	relations, ok := v.([]interface{})
	if !ok || len(relations) == 0 {
		return object, WrapErrorf(NotFoundErr("TairInstance", id), NotFoundMsg, response)
	}

	// Join every EcsSecurityGroupRelation[*].SecurityGroupId in sorted order so state is
	// stable across API return-order fluctuations. Build a fresh map (never alias into
	// relations[0]) so we do not mutate the RPC response; the caller only reads RegionId
	// and SecurityGroupId, so those are the only fields we carry over.
	securityGroupIds := make([]string, 0, len(relations))
	for _, relation := range relations {
		rel, ok := relation.(map[string]interface{})
		if !ok {
			continue
		}
		securityGroupIds = append(securityGroupIds, fmt.Sprint(rel["SecurityGroupId"]))
	}
	sort.Strings(securityGroupIds)

	object = map[string]interface{}{
		"SecurityGroupId": strings.Join(securityGroupIds, ","),
	}
	if first, ok := relations[0].(map[string]interface{}); ok {
		if regionId, ok := first["RegionId"]; ok {
			object["RegionId"] = regionId
		}
	}

	return object, nil
}
func (s *RedisServiceV2) DescribeTairInstanceDescribeInstanceTDEStatus(id string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["InstanceId"] = id
	request["RegionId"] = client.RegionId
	action := "DescribeInstanceTDEStatus"

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(1*time.Minute, func() *resource.RetryError {
		response, err = client.RpcPost("R-kvstore", "2015-01-01", action, query, request, true)

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
		if IsExpectedErrors(err, []string{"InvalidInstanceId.NotFound"}) {
			return object, WrapErrorf(NotFoundErr("TairInstance", id), NotFoundMsg, response)
		}
		// Defense in depth: instances that do not support TDE reject DescribeInstanceTDEStatus
		// with a 400 InstanceType.NotSupport. Treat it as "TDE not supported" by returning an
		// empty object without error, so the read is not aborted.
		if IsExpectedErrors(err, []string{"InstanceType.NotSupport"}) {
			return object, nil
		}
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}

	return response, nil
}
func (s *RedisServiceV2) DescribeTairInstanceDescribeInstanceSSL(id string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	action := "DescribeInstanceSSL"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["InstanceId"] = id
	request["RegionId"] = client.RegionId

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(1*time.Minute, func() *resource.RetryError {
		response, err = client.RpcPost("R-kvstore", "2015-01-01", action, query, request, true)

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
		if IsExpectedErrors(err, []string{"InvalidInstanceId.NotFound"}) {
			return object, WrapErrorf(NotFoundErr("TairInstance", id), NotFoundMsg, response)
		}
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}

	return response, nil
}
func (s *RedisServiceV2) DescribeTairInstanceDescribeIntranetAttribute(id string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	action := "DescribeIntranetAttribute"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["InstanceId"] = id
	request["RegionId"] = client.RegionId

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(1*time.Minute, func() *resource.RetryError {
		response, err = client.RpcPost("R-kvstore", "2015-01-01", action, query, request, true)

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
		if IsExpectedErrors(err, []string{"InvalidInstanceId.NotFound"}) {
			return object, WrapErrorf(NotFoundErr("TairInstance", id), NotFoundMsg, response)
		}
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}

	return response, nil
}

func (s *RedisServiceV2) RedisTairInstanceStateRefreshFunc(id string, field string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeRedisTairInstance(id)
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

// DescribeRedisTairInstance >>> Encapsulated.

// SetResourceTags <<< Encapsulated tag function for Redis.
func (s *RedisServiceV2) SetResourceTags(d *schema.ResourceData, resourceType string) error {
	if d.HasChange("tags") {
		var err error
		var action string
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
			request["RegionId"] = client.RegionId
			request["ResourceType"] = resourceType
			for i, key := range removedTagKeys {
				request[fmt.Sprintf("TagKey.%d", i+1)] = key
			}

			wait := incrementalWait(3*time.Second, 5*time.Second)
			err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
				response, err = client.RpcPost("R-kvstore", "2015-01-01", action, query, request, false)
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
				response, err = client.RpcPost("R-kvstore", "2015-01-01", action, query, request, false)
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

// DescribeRedisAccount <<< Encapsulated get interface for Redis Account.

func (s *RedisServiceV2) DescribeRedisAccount(id string) (object map[string]interface{}, err error) {
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
	request["AccountName"] = parts[1]
	request["InstanceId"] = parts[0]
	request["RegionId"] = client.RegionId
	action := "DescribeAccounts"

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(1*time.Minute, func() *resource.RetryError {
		response, err = client.RpcPost("R-kvstore", "2015-01-01", action, query, request, true)

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
		if IsExpectedErrors(err, []string{"InvalidInstanceId.NotFound"}) {
			return object, WrapErrorf(NotFoundErr("Account", id), NotFoundMsg, response)
		}
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}

	v, err := jsonpath.Get("$.Accounts.Account[*]", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.Accounts.Account[*]", response)
	}

	if len(v.([]interface{})) == 0 {
		return object, WrapErrorf(NotFoundErr("Account", id), NotFoundMsg, response)
	}

	return v.([]interface{})[0].(map[string]interface{}), nil
}

func (s *RedisServiceV2) RedisAccountStateRefreshFunc(id string, field string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeRedisAccount(id)
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

// DescribeRedisAccount >>> Encapsulated.
// DescribeRedisBackup <<< Encapsulated get interface for Redis Backup.

func (s *RedisServiceV2) DescribeRedisBackup(id string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	parts := strings.Split(id, ":")
	if len(parts) != 2 {
		err = WrapError(fmt.Errorf("invalid Resource Id %s. Expected parts' length %d, got %d", id, 2, len(parts)))
		return nil, err
	}
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["BackupId"] = parts[1]
	request["InstanceId"] = parts[0]
	request["RegionId"] = client.RegionId
	request["EndTime"] = "2050-01-01T10:00Z"
	request["StartTime"] = "2010-01-01T10:00Z"
	action := "DescribeBackups"

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(1*time.Minute, func() *resource.RetryError {
		response, err = client.RpcPost("R-kvstore", "2015-01-01", action, query, request, true)

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

	v, err := jsonpath.Get("$.Backups.Backup[*]", response)
	if err != nil {
		return object, WrapErrorf(NotFoundErr("Backup", id), NotFoundMsg, response)
	}

	if len(v.([]interface{})) == 0 {
		return object, WrapErrorf(NotFoundErr("Backup", id), NotFoundMsg, response)
	}

	backup := v.([]interface{})[0].(map[string]interface{})

	// Check if BackupId exists (it's the actual backup identifier)
	if backupID, ok := backup["BackupId"]; !ok || backupID == nil {
		return object, WrapErrorf(NotFoundErr("Backup", id), NotFoundMsg, response)
	}

	return backup, nil
}

func (s *RedisServiceV2) RedisBackupStateRefreshFunc(id string, field string, failStates []string) resource.StateRefreshFunc {
	return s.RedisBackupStateRefreshFuncWithApi(id, field, failStates, s.DescribeRedisBackup)
}

func (s *RedisServiceV2) RedisBackupStateRefreshFuncWithApi(id string, field string, failStates []string, call func(id string) (map[string]interface{}, error)) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := call(id)
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

func (s *RedisServiceV2) DescribeAsyncRedisBackupStateRefreshFunc(d *schema.ResourceData, res map[string]interface{}, field string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeAsyncDescribeBackups(d, res)
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

// DescribeRedisBackup >>> Encapsulated.
// DescribeAsyncDescribeBackups <<< Encapsulated for Redis.
func (s *RedisServiceV2) DescribeAsyncDescribeBackups(d *schema.ResourceData, res map[string]interface{}) (object map[string]interface{}, err error) {
	client := s.client
	id := d.Id()
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	parts := strings.Split(id, ":")
	if len(parts) != 2 {
		err = WrapError(fmt.Errorf("invalid Resource Id %s. Expected parts' length %d, got %d", id, 2, len(parts)))
		return nil, err
	}
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["InstanceId"] = parts[0]
	request["RegionId"] = client.RegionId
	request["StartTime"] = "2010-01-01T10:00Z"
	request["EndTime"] = "2050-01-01T10:00Z"
	request["PageSize"] = PageSizeLarge
	action := "DescribeBackups"

	// Get BackupJobID from create response
	backupJobID := ""
	if res != nil {
		if jobID, ok := res["BackupJobID"]; ok {
			backupJobID = fmt.Sprint(jobID)
		}
	}

	// If we don't have JobID from response, it might be in the ID (during waiting period)
	if backupJobID == "" {
		backupJobID = parts[1]
	}

	// Iterate through all pages to find the backup with matching JobID
	pageNumber := 1
	for {
		request["PageNumber"] = pageNumber

		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(1*time.Minute, func() *resource.RetryError {
			response, err = client.RpcPost("R-kvstore", "2015-01-01", action, query, request, true)

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

		// Filter backups by BackupJobID to find the one created by this job
		if backupJobID != "" {
			v, err := jsonpath.Get("$.Backups.Backup[*]", response)
			if err == nil && v != nil {
				if backups, ok := v.([]interface{}); ok {
					for _, item := range backups {
						if backup, ok := item.(map[string]interface{}); ok {
							if jobID, ok := backup["BackupJobID"]; ok {
								if fmt.Sprint(jobID) == backupJobID {
									// Found the backup with matching JobID
									// Check if BackupId is available (backup completed)
									if backupID, ok := backup["BackupId"]; ok && backupID != nil && fmt.Sprint(backupID) != "" {
										// Return a response structure that contains this backup
										filteredResponse := make(map[string]interface{})
										filteredResponse["Backups"] = map[string]interface{}{
											"Backup": []interface{}{backup},
										}
										return filteredResponse, nil
									}
									// Found the job but BackupId not ready yet, return not found error to continue waiting
									return object, WrapErrorf(NotFoundErr("Backup", id), NotFoundMsg, response)
								}
							}
						}
					}
				}
			}
		}

		// Check if there are more pages
		totalCount := 0
		if tc, ok := response["TotalCount"]; ok {
			if tcInt, ok := tc.(float64); ok {
				totalCount = int(tcInt)
			} else if tcInt, ok := tc.(int); ok {
				totalCount = tcInt
			}
		}

		// If we've checked all pages and haven't found it, break
		if pageNumber*30 >= totalCount {
			break
		}

		pageNumber++
	}

	// Backup job not found in any page, return not found error to continue waiting
	return object, WrapErrorf(NotFoundErr("Backup", id), NotFoundMsg, response)
}

// DescribeAsyncDescribeBackups >>> Encapsulated.

// DescribeRedisClusterBackupId <<< Resolve the cluster backup set id for a backup.

// DescribeRedisClusterBackupId returns the cluster backup set id (cb-*) that owns the
// given per-shard BackupId. Standalone instances legitimately have no cluster backup set;
// an empty result with a nil error means "not a cluster backup", not a failure.
func (s *RedisServiceV2) DescribeRedisClusterBackupId(instanceId string, backupId string) (string, error) {
	response, err := s.describeClusterBackupList(instanceId, "")
	if err != nil {
		return "", err
	}

	for _, item := range redisRpcList(response["ClusterBackups"], "ClusterBackup") {
		clusterBackup, ok := item.(map[string]interface{})
		if !ok {
			continue
		}
		clusterBackupId := fmt.Sprint(clusterBackup["ClusterBackupId"])
		if clusterBackupId == "" || clusterBackupId == "<nil>" {
			continue
		}
		for _, shard := range redisRpcList(clusterBackup["Backups"], "Backup") {
			shardBackup, ok := shard.(map[string]interface{})
			if !ok {
				continue
			}
			if fmt.Sprint(shardBackup["BackupId"]) == backupId {
				return clusterBackupId, nil
			}
		}
	}

	return "", nil
}

// ListClusterBackupShardIds returns every shard BackupId that belongs to the given
// cluster backup set. Used by resource_alicloud_redis_backup Delete to fan out shard-level
// DeleteBackup calls (R-kvstore exposes no cluster-level delete API). An empty slice with
// nil error means the cluster set is gone or has no shard backups — both are non-fatal.
func (s *RedisServiceV2) ListClusterBackupShardIds(instanceId string, clusterBackupId string) ([]string, error) {
	response, err := s.describeClusterBackupList(instanceId, clusterBackupId)
	if err != nil {
		return nil, err
	}

	shardIds := make([]string, 0)
	for _, item := range redisRpcList(response["ClusterBackups"], "ClusterBackup") {
		clusterBackup, ok := item.(map[string]interface{})
		if !ok {
			continue
		}
		if clusterBackupId != "" && fmt.Sprint(clusterBackup["ClusterBackupId"]) != clusterBackupId {
			continue
		}
		for _, shard := range redisRpcList(clusterBackup["Backups"], "Backup") {
			shardBackup, ok := shard.(map[string]interface{})
			if !ok {
				continue
			}
			if id := fmt.Sprint(shardBackup["BackupId"]); id != "" && id != "<nil>" {
				shardIds = append(shardIds, id)
			}
		}
	}
	return shardIds, nil
}

// describeClusterBackupList calls R-kvstore DescribeClusterBackupList with a wide UTC time
// window that covers any live backup set. StartTime / EndTime are Required by the API even
// when ClusterBackupId is supplied; a 30-day-back to 1-day-forward window absorbs clock
// skew and long-lived retention without risking a false miss.
func (s *RedisServiceV2) describeClusterBackupList(instanceId string, clusterBackupId string) (map[string]interface{}, error) {
	client := s.client
	var response map[string]interface{}
	var err error

	query := map[string]interface{}{
		"InstanceId": instanceId,
		"RegionId":   client.RegionId,
		"PageSize":   100,
	}
	now := time.Now().UTC()
	query["StartTime"] = now.Add(-30 * 24 * time.Hour).Format("2006-01-02T15:04Z")
	query["EndTime"] = now.Add(24 * time.Hour).Format("2006-01-02T15:04Z")
	if clusterBackupId != "" {
		query["ClusterBackupId"] = clusterBackupId
	}
	action := "DescribeClusterBackupList"

	// DescribeClusterBackupList only accepts GET (RpcPost is rejected with 403
	// UnsupportedHTTPMethod), so all parameters go in the query and the body is nil.
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(1*time.Minute, func() *resource.RetryError {
		response, err = client.RpcGet("R-kvstore", "2015-01-01", action, query, nil)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, query)
	if err != nil {
		return nil, WrapErrorf(err, DefaultErrorMsg, instanceId, action, AlibabaCloudSdkGoERROR)
	}
	return response, nil
}

// redisRpcList normalizes an R-kvstore RPC list field that may arrive either as a flat
// JSON array or wrapped in a single-key object ({"<elementKey>": [...]}).
func redisRpcList(v interface{}, elementKey string) []interface{} {
	switch t := v.(type) {
	case []interface{}:
		return t
	case map[string]interface{}:
		if inner, ok := t[elementKey]; ok {
			if list, ok := inner.([]interface{}); ok {
				return list
			}
		}
	}
	return nil
}

// DescribeRedisClusterBackupId >>> Encapsulated.
