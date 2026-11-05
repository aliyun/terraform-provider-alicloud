// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package polardb

import (
	"fmt"
	"strings"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/blues/jsonata-go"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type PolarDbServiceV2 struct {
	client *connectivity.AliyunClient
}

func NewPolarDbServiceV2(client *connectivity.AliyunClient) *PolarDbServiceV2 {
	return &PolarDbServiceV2{client: client}
}

// DescribePolarDbParameterGroup <<< Encapsulated get interface for PolarDb ParameterGroup.

func (s *PolarDbServiceV2) DescribePolarDbParameterGroup(id string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["ParameterGroupId"] = id
	request["RegionId"] = client.RegionId
	action := "DescribeParameterGroup"

	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(1*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("polardb", "2017-08-01", action, query, request, true)

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
		if errs.IsExpectedErrors(err, []string{"ParamGroupsNotExist"}) {
			return object, errs.WrapErrorf(errs.NotFoundErr("ParameterGroup", id), NotFoundMsg, response)
		}
		return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}

	v, err := jsonpath.Get("$.ParameterGroup[*]", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.ParameterGroup[*]", response)
	}

	if len(v.([]interface{})) == 0 {
		return object, errs.WrapErrorf(errs.NotFoundErr("ParameterGroup", id), NotFoundMsg, response)
	}

	return v.([]interface{})[0].(map[string]interface{}), nil
}

func (s *PolarDbServiceV2) PolarDbParameterGroupStateRefreshFunc(id string, field string, failStates []string) retry.StateRefreshFunc {
	return s.PolarDbParameterGroupStateRefreshFuncWithApi(id, field, failStates, s.DescribePolarDbParameterGroup)
}

func (s *PolarDbServiceV2) PolarDbParameterGroupStateRefreshFuncWithApi(id string, field string, failStates []string, call func(id string) (map[string]interface{}, error)) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := call(id)
		if err != nil {
			if errs.NotFoundError(err) {
				return object, "", nil
			}
			return nil, "", errs.WrapError(err)
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
				return object, currentStatus, errs.WrapError(errs.Error(FailedToReachTargetStatus, currentStatus))
			}
		}
		return object, currentStatus, nil
	}
}

// DescribePolarDbParameterGroup >>> Encapsulated.

// DescribePolarDbExtension <<< Encapsulated get interface for PolarDb Extension.

func (s *PolarDbServiceV2) DescribePolarDbExtension(id string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	parts := strings.Split(id, ":")
	if len(parts) != 4 {
		err = errs.WrapError(fmt.Errorf("invalid Resource Id %s. Expected parts' length %d, got %d", id, 4, len(parts)))
		return nil, err
	}
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["DBClusterId"] = parts[0]
	request["DBName"] = parts[2]
	request["ExtensionName"] = parts[3]

	action := "DescribeExtensions"

	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(1*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("polardb", "2017-08-01", action, query, request, true)

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

	v, err := jsonpath.Get("$.InstalledExtensions[*]", response)
	if err != nil {
		return object, errs.WrapErrorf(errs.NotFoundErr("Extension", id), NotFoundMsg, response)
	}

	if len(v.([]interface{})) == 0 {
		return object, errs.WrapErrorf(errs.NotFoundErr("Extension", id), NotFoundMsg, response)
	}

	result, _ := v.([]interface{})
	for _, v := range result {
		item := v.(map[string]interface{})
		if fmt.Sprint(item["Owner"]) != parts[1] {
			continue
		}
		return item, nil
	}
	return object, errs.WrapErrorf(errs.NotFoundErr("Extension", id), NotFoundMsg, response)
}

func (s *PolarDbServiceV2) PolarDbExtensionStateRefreshFunc(id string, field string, failStates []string) retry.StateRefreshFunc {
	return s.PolarDbExtensionStateRefreshFuncWithApi(id, field, failStates, s.DescribePolarDbExtension)
}

func (s *PolarDbServiceV2) PolarDbExtensionStateRefreshFuncWithApi(id string, field string, failStates []string, call func(id string) (map[string]interface{}, error)) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := call(id)
		if err != nil {
			if errs.NotFoundError(err) {
				return object, "", nil
			}
			return nil, "", errs.WrapError(err)
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
				return object, currentStatus, errs.WrapError(errs.Error(FailedToReachTargetStatus, currentStatus))
			}
		}
		return object, currentStatus, nil
	}
}

// DescribePolarDbExtension >>> Encapsulated.

// DescribeAsyncDescribeExtensions <<< Encapsulated for PolarDb.
func (s *PolarDbServiceV2) DescribeAsyncDescribeExtensions(d *schema.ResourceData, res map[string]interface{}) (object map[string]interface{}, err error) {
	client := s.client
	id := d.Id()
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	parts := strings.Split(id, ":")
	if len(parts) != 4 {
		err = errs.WrapError(fmt.Errorf("invalid Resource Id %s. Expected parts' length %d, got %d", id, 4, len(parts)))
	}
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["ExtensionName"] = parts[3]
	request["DBClusterId"] = parts[0]
	request["DBName"] = parts[1]

	action := "DescribeExtensions"

	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(1*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("polardb", "2017-08-01", action, query, request, true)

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
		return response, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}

	return response, nil
}

// DescribeAsyncDescribeExtensions >>> Encapsulated.

// DescribePolarDbAccount <<< Encapsulated get interface for PolarDb Account.

func (s *PolarDbServiceV2) DescribePolarDbAccount(id string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	parts := strings.Split(id, ":")
	if len(parts) != 2 {
		err = errs.WrapError(fmt.Errorf("invalid Resource Id %s. Expected parts' length %d, got %d", id, 2, len(parts)))
		return nil, err
	}
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["AccountName"] = parts[1]
	request["DBClusterId"] = parts[0]

	action := "DescribeAccounts"

	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(1*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("polardb", "2017-08-01", action, query, request, true)

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
		if errs.IsExpectedErrors(err, []string{"InvalidDBClusterId.NotFound"}) {
			return object, errs.WrapErrorf(errs.NotFoundErr("Account", id), NotFoundMsg, response)
		}
		return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}

	v, err := jsonpath.Get("$.Accounts[*]", response)
	if err != nil {
		return object, errs.WrapErrorf(errs.NotFoundErr("Account", id), NotFoundMsg, response)
	}

	if len(v.([]interface{})) == 0 {
		return object, errs.WrapErrorf(errs.NotFoundErr("Account", id), NotFoundMsg, response)
	}

	currentStatus := v.([]interface{})[0].(map[string]interface{})["AccountName"]
	if currentStatus == nil {
		return object, errs.WrapErrorf(errs.NotFoundErr("Account", id), NotFoundMsg, response)
	}

	return v.([]interface{})[0].(map[string]interface{}), nil
}

func (s *PolarDbServiceV2) PolarDbAccountStateRefreshFunc(id string, field string, failStates []string) retry.StateRefreshFunc {
	return s.PolarDbAccountStateRefreshFuncWithApi(id, field, failStates, s.DescribePolarDbAccount)
}

func (s *PolarDbServiceV2) PolarDbAccountStateRefreshFuncWithApi(id string, field string, failStates []string, call func(id string) (map[string]interface{}, error)) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := call(id)
		if err != nil {
			if errs.NotFoundError(err) {
				return object, "", nil
			}
			return nil, "", errs.WrapError(err)
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
				return object, currentStatus, errs.WrapError(errs.Error(FailedToReachTargetStatus, currentStatus))
			}
		}
		return object, currentStatus, nil
	}
}

// DescribePolarDbAccount >>> Encapsulated.

// DescribePolarDbDatabase <<< Encapsulated get interface for PolarDb Database.

func (s *PolarDbServiceV2) DescribePolarDbDatabase(id string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	parts := strings.Split(id, ":")
	if len(parts) != 2 {
		err = errs.WrapError(fmt.Errorf("invalid Resource Id %s. Expected parts' length %d, got %d", id, 2, len(parts)))
		return nil, err
	}
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["DBClusterId"] = parts[0]
	request["DBName"] = parts[1]

	action := "DescribeDatabases"

	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(1*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("polardb", "2017-08-01", action, query, request, true)

		if err != nil {
			if errs.IsExpectedErrors(err, []string{"InvalidDBClusterId.NotFound"}) || helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(action, response, request)
	if err != nil {
		if errs.IsExpectedErrors(err, []string{"InvalidDBClusterId.NotFound"}) {
			return object, errs.WrapErrorf(errs.NotFoundErr("Database", id), NotFoundMsg, response)
		}
		return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}

	v, err := jsonpath.Get("$.Databases.Database[*]", response)
	if err != nil {
		return object, errs.WrapErrorf(errs.NotFoundErr("Database", id), NotFoundMsg, response)
	}

	if len(v.([]interface{})) == 0 {
		return object, errs.WrapErrorf(errs.NotFoundErr("Database", id), NotFoundMsg, response)
	}

	currentStatus := v.([]interface{})[0].(map[string]interface{})["DBName"]
	if currentStatus == nil {
		return object, errs.WrapErrorf(errs.NotFoundErr("Database", id), NotFoundMsg, response)
	}

	return v.([]interface{})[0].(map[string]interface{}), nil
}

func (s *PolarDbServiceV2) PolarDbDatabaseStateRefreshFunc(id string, field string, failStates []string) retry.StateRefreshFunc {
	return s.PolarDbDatabaseStateRefreshFuncWithApi(id, field, failStates, s.DescribePolarDbDatabase)
}

func (s *PolarDbServiceV2) PolarDbDatabaseStateRefreshFuncWithApi(id string, field string, failStates []string, call func(id string) (map[string]interface{}, error)) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := call(id)
		if err != nil {
			if errs.NotFoundError(err) {
				return object, "", nil
			}
			return nil, "", errs.WrapError(err)
		}
		v, err := jsonpath.Get(field, object)
		currentStatus := fmt.Sprint(v)
		if field == "$.Accounts.Account[*]" {
			e := jsonata.MustCompile("$.Accounts.Account[0].AccountName")
			v, _ = e.Eval(object)
			currentStatus = fmt.Sprint(v)
		}

		if strings.HasPrefix(field, "#") {
			v, _ := jsonpath.Get(strings.TrimPrefix(field, "#"), object)
			if v != nil {
				currentStatus = "#CHECKSET"
			}
		}

		for _, failState := range failStates {
			if currentStatus == failState {
				return object, currentStatus, errs.WrapError(errs.Error(FailedToReachTargetStatus, currentStatus))
			}
		}
		return object, currentStatus, nil
	}
}

// DescribePolarDbDatabase >>> Encapsulated.
