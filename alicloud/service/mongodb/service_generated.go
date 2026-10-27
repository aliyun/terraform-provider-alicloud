// Package mongodb. This file is generated automatically. Please do not modify it manually, thank you!
package mongodb

import (
	tfhelper "github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	tflogs "github.com/aliyun/terraform-provider-alicloud/alicloud/logs"

	tferrs "github.com/aliyun/terraform-provider-alicloud/alicloud/errs"

	"fmt"
	"strings"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	tfretry "github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
)

type MongodbServiceV2 struct {
	client *connectivity.AliyunClient
}

// DescribeMongodbAccount <<< Encapsulated get interface for Mongodb Account.

func (s *MongodbServiceV2) DescribeMongodbAccount(id string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	parts := strings.Split(id, ":")
	if len(parts) != 2 && len(parts) != 3 {
		err = tferrs.WrapError(fmt.Errorf("invalid Resource Id %s. Expected parts' length %d or %d, got %d", id, 2, 3, len(parts)))
		return nil, err
	}
	accountName := parts[len(parts)-1]
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["AccountName"] = accountName
	request["DBInstanceId"] = parts[0]
	request["RegionId"] = client.RegionId
	if len(parts) == 3 {
		request["DBName"] = parts[1]
	}
	action := "DescribeAccounts"

	wait := tfhelper.IncrementalWait(3*time.Second, 5*time.Second)
	err = tfretry.Retry(1*time.Minute, func() *tfretry.RetryError {
		response, err = client.RpcPost("Dds", "2015-12-01", action, query, request, true)

		if err != nil {
			if tfhelper.NeedRetry(err) {
				wait()
				return tfretry.RetryableError(err)
			}
			return tfretry.NonRetryableError(err)
		}
		return nil
	})
	tflogs.AddDebug(action, response, request)
	if err != nil {
		return object, tferrs.WrapErrorf(err, tferrs.DefaultErrorMsg, id, action, tferrs.AlibabaCloudSdkGoERROR)
	}

	v, err := jsonpath.Get("$.Accounts.Account[*]", response)
	if err != nil {
		return object, tferrs.WrapErrorf(err, tferrs.FailedGetAttributeMsg, id, "$.Accounts.Account[*]", response)
	}

	accounts, ok := v.([]interface{})
	if !ok {
		return object, tferrs.WrapErrorf(fmt.Errorf("unexpected DescribeAccounts response format for %s", id), tferrs.FailedGetAttributeMsg, id, "$.Accounts.Account[*]", response)
	}
	matched := make([]map[string]interface{}, 0, 1)
	for _, item := range accounts {
		account, ok := item.(map[string]interface{})
		if !ok {
			continue
		}
		if fmt.Sprint(account["AccountName"]) != accountName {
			continue
		}
		// Account names are only unique inside a database, so filter by the
		// database scope as well whenever the caller supplied one.
		if len(parts) == 3 && fmt.Sprint(account["DBName"]) != parts[1] {
			continue
		}
		matched = append(matched, account)
	}

	if len(matched) == 0 {
		return object, tferrs.WrapErrorf(tferrs.NotFoundErr("Account", id), tferrs.NotFoundMsg, response)
	}
	if len(matched) > 1 {
		return object, tferrs.WrapError(fmt.Errorf("found %d accounts named %s in instance %s, use id <instance_id>:<db_name>:<account_name> to identify the account", len(matched), accountName, parts[0]))
	}

	return matched[0], nil
}

func (s *MongodbServiceV2) MongodbAccountStateRefreshFunc(id string, field string, failStates []string) tfretry.StateRefreshFunc {
	return s.MongodbAccountStateRefreshFuncWithApi(id, field, failStates, s.DescribeMongodbAccount)
}

func (s *MongodbServiceV2) MongodbAccountStateRefreshFuncWithApi(id string, field string, failStates []string, call func(id string) (map[string]interface{}, error)) tfretry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := call(id)
		if err != nil {
			if tferrs.NotFoundError(err) {
				return object, "", nil
			}
			return nil, "", tferrs.WrapError(err)
		}
		var currentStatus string
		if strings.HasPrefix(field, "#") {
			v, _ := jsonpath.Get(strings.TrimPrefix(field, "#"), object)
			if v != nil {
				currentStatus = "#CHECKSET"
			}
		} else {
			v, err := jsonpath.Get(field, object)
			if err != nil {
				return nil, "", tferrs.WrapError(err)
			}
			if v == nil {
				return nil, "", tferrs.WrapError(fmt.Errorf("field %s is missing in the DescribeAccounts response of %s", field, id))
			}
			currentStatus = fmt.Sprint(v)
		}

		for _, failState := range failStates {
			if currentStatus == failState {
				return object, currentStatus, tferrs.WrapError(tferrs.Error(tferrs.FailedToReachTargetStatus, currentStatus))
			}
		}
		return object, currentStatus, nil
	}
}

// DescribeMongodbAccount >>> Encapsulated.
