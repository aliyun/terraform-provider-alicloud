package alicloud

import (
	"fmt"
	"strings"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type RedisServiceV2 struct {
	client *connectivity.AliyunClient
}

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
	err = retry.Retry(1*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("R-kvstore", "2015-01-01", action, query, request, true)

		if err != nil {
			if NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
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

func (s *RedisServiceV2) RedisAccountStateRefreshFunc(id string, field string, failStates []string) retry.StateRefreshFunc {
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
	err = retry.Retry(1*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("R-kvstore", "2015-01-01", action, query, request, true)

		if err != nil {
			if NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
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

func (s *RedisServiceV2) RedisBackupStateRefreshFunc(id string, field string, failStates []string) retry.StateRefreshFunc {
	return s.RedisBackupStateRefreshFuncWithApi(id, field, failStates, s.DescribeRedisBackup)
}

func (s *RedisServiceV2) RedisBackupStateRefreshFuncWithApi(id string, field string, failStates []string, call func(id string) (map[string]interface{}, error)) retry.StateRefreshFunc {
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

func (s *RedisServiceV2) DescribeAsyncRedisBackupStateRefreshFunc(d *schema.ResourceData, res map[string]interface{}, field string, failStates []string) retry.StateRefreshFunc {
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
		err = retry.Retry(1*time.Minute, func() *retry.RetryError {
			response, err = client.RpcPost("R-kvstore", "2015-01-01", action, query, request, true)

			if err != nil {
				if NeedRetry(err) {
					wait()
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
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
// DescribeRedisGlobalSecurityIpGroup <<< Encapsulated get interface for Redis GlobalSecurityIpGroup.

func (s *RedisServiceV2) DescribeRedisGlobalSecurityIpGroup(id string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["GlobalSecurityGroupId"] = id
	request["RegionId"] = client.RegionId
	action := "DescribeGlobalSecurityIPGroup"

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(1*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("R-kvstore", "2015-01-01", action, query, request, true)

		if err != nil {
			if NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)
	if err != nil {
		if IsExpectedErrors(err, []string{"RequiredParam.NotFound"}) {
			return object, WrapErrorf(NotFoundErr("GlobalSecurityIpGroup", id), NotFoundMsg, response)
		}
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}

	v, err := jsonpath.Get("$.GlobalSecurityIPGroup[*]", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.GlobalSecurityIPGroup[*]", response)
	}

	if len(v.([]interface{})) == 0 {
		return object, WrapErrorf(NotFoundErr("GlobalSecurityIpGroup", id), NotFoundMsg, response)
	}

	// DescribeGlobalSecurityIPGroup is a list query that returns all global
	// security IP groups under the account, so the matching item must be
	// located by GlobalSecurityGroupId. Blindly taking the first element causes
	// a persistent diff when multiple groups exist (e.g. count > 1).
	result, _ := v.([]interface{})
	for _, item := range result {
		group := item.(map[string]interface{})
		if fmt.Sprint(group["GlobalSecurityGroupId"]) != id {
			continue
		}
		return group, nil
	}
	return object, WrapErrorf(NotFoundErr("GlobalSecurityIpGroup", id), NotFoundMsg, response)
}

func (s *RedisServiceV2) RedisGlobalSecurityIpGroupStateRefreshFunc(id string, field string, failStates []string) retry.StateRefreshFunc {
	return s.RedisGlobalSecurityIpGroupStateRefreshFuncWithApi(id, field, failStates, s.DescribeRedisGlobalSecurityIpGroup)
}

func (s *RedisServiceV2) RedisGlobalSecurityIpGroupStateRefreshFuncWithApi(id string, field string, failStates []string, call func(id string) (map[string]interface{}, error)) retry.StateRefreshFunc {
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

// DescribeRedisGlobalSecurityIpGroup >>> Encapsulated.
