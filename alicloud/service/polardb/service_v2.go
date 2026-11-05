// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package polardb

import (
	"errors"
	"fmt"
	"strconv"
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

// DescribeAccounts is queried page by page because the API contract does not guarantee that
// AccountName filters the result set; an account beyond the first page must not be reported as
// missing, which would make Read drop a live account from state.
const (
	polarDbAccountPageSize = 100
	polarDbAccountMaxPages = 100
)

// errPolarDbAccountPageLimitExceeded reports that pagination stopped before a short page; the account
// may still exist, so callers must not treat it as missing and clear the resource from state.
var errPolarDbAccountPageLimitExceeded = errors.New("DescribeAccounts pagination exceeded the page limit before reaching the last page")

// DescribePolarDbAccount <<< Encapsulated get interface for PolarDb Account.

func (s *PolarDbServiceV2) DescribePolarDbAccount(id string) (object map[string]interface{}, err error) {
	var request map[string]interface{}
	var query map[string]interface{}
	parts := strings.Split(id, ":")
	if len(parts) != 2 && len(parts) != 3 {
		err = errs.WrapError(fmt.Errorf("invalid Resource Id %s. Expected parts' length %d or %d, got %d", id, 2, 3, len(parts)))
		return nil, err
	}
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["AccountName"] = parts[1]
	request["DBClusterId"] = parts[0]
	request["PageSize"] = polarDbAccountPageSize
	if len(parts) == 3 && parts[2] != "" {
		request["NodeType"] = parts[2]
	}

	action := "DescribeAccounts"
	account, lastResponse, found, err := findPolarDbAccountByPages(parts[1], polarDbAccountPageSize, func(pageNumber int) (map[string]interface{}, error) {
		request["PageNumber"] = pageNumber
		return s.describePolarDbAccountsPage(action, query, request)
	})
	if err != nil {
		if errs.IsExpectedErrors(err, []string{"InvalidDBClusterId.NotFound"}) {
			return object, errs.WrapErrorf(errs.NotFoundErr("Account", id), NotFoundMsg, lastResponse)
		}
		return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}
	if found {
		return account, nil
	}
	return object, errs.WrapErrorf(errs.NotFoundErr("Account", id), NotFoundMsg, lastResponse)
}

func (s *PolarDbServiceV2) describePolarDbAccountsPage(action string, query, request map[string]interface{}) (map[string]interface{}, error) {
	var response map[string]interface{}
	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err := retry.Retry(1*time.Minute, func() *retry.RetryError {
		var err error
		response, err = s.client.RpcPost("polardb", "2017-08-01", action, query, request, true)

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
	return response, err
}

// findPolarDbAccountByPages walks the DescribeAccounts pages until the account is found or the
// last page is reached. fetch returns the response for a 1-based page number.
func findPolarDbAccountByPages(name string, pageSize int, fetch func(pageNumber int) (map[string]interface{}, error)) (map[string]interface{}, map[string]interface{}, bool, error) {
	var lastResponse map[string]interface{}
	for pageNumber := 1; pageNumber <= polarDbAccountMaxPages; pageNumber++ {
		response, err := fetch(pageNumber)
		if err != nil {
			return nil, lastResponse, false, err
		}
		lastResponse = response
		accounts, err := jsonpath.Get("$.Accounts[*]", response)
		if err != nil {
			return nil, lastResponse, false, nil
		}
		list, _ := accounts.([]interface{})
		if account, ok := findPolarDbAccountByName(list, name); ok {
			return account, lastResponse, true, nil
		}
		if polarDbAccountPageRecordCount(response, len(list)) < pageSize {
			return nil, lastResponse, false, nil
		}
	}
	return nil, lastResponse, false, errPolarDbAccountPageLimitExceeded
}

func findPolarDbAccountByName(accounts []interface{}, name string) (map[string]interface{}, bool) {
	for _, item := range accounts {
		account, ok := item.(map[string]interface{})
		if !ok {
			continue
		}
		if accountName, ok := account["AccountName"].(string); ok && accountName == name {
			return account, true
		}
	}
	return nil, false
}

func polarDbAccountPageRecordCount(response map[string]interface{}, fallback int) int {
	if v, err := jsonpath.Get("$.PageRecordCount", response); err == nil {
		switch n := v.(type) {
		case float64:
			return int(n)
		case int:
			return n
		case string:
			if count, err := strconv.Atoi(n); err == nil {
				return count
			}
		}
	}
	return fallback
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
