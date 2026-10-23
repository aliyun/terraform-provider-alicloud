// Package polardb. This file is generated automatically. Please do not modify it manually, thank you!
package polardb

import (
	tfhelper "github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	tflogs "github.com/aliyun/terraform-provider-alicloud/alicloud/logs"

	tferrs "github.com/aliyun/terraform-provider-alicloud/alicloud/errs"

	"fmt"
	"strings"
	"time"

	"github.com/PaesslerAG/jsonpath"
	tfretry "github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	jsonata "github.com/blues/jsonata-go"
)

type PolardbServiceV2 struct {
	client *connectivity.AliyunClient
}

// DescribePolardbDatabase <<< Encapsulated get interface for Polardb Database.

func (s *PolardbServiceV2) DescribePolardbDatabase(id string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	parts := strings.Split(id, ":")
	if len(parts) != 2 {
		err = tferrs.WrapError(fmt.Errorf("invalid Resource Id %s. Expected parts' length %d, got %d", id, 2, len(parts)))
		return nil, err
	}
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["DBClusterId"] = parts[0]
	request["DBName"] = parts[1]

	action := "DescribeDatabases"

	wait := tfhelper.IncrementalWait(3*time.Second, 5*time.Second)
	err = tfretry.Retry(1*time.Minute, func() *tfretry.RetryError {
		response, err = client.RpcPost("polardb", "2017-08-01", action, query, request, true)

		if err != nil {
			if tferrs.IsExpectedErrors(err, []string{"InvalidDBClusterId.NotFound"}) || tfhelper.NeedRetry(err) {
				wait()
				return tfretry.RetryableError(err)
			}
			return tfretry.NonRetryableError(err)
		}
		return nil
	})
	tflogs.AddDebug(action, response, request)
	if err != nil {
		if tferrs.IsExpectedErrors(err, []string{"InvalidDBClusterId.NotFound"}) {
			return object, tferrs.WrapErrorf(tferrs.NotFoundErr("Database", id), tferrs.NotFoundMsg, response)
		}
		return object, tferrs.WrapErrorf(err, tferrs.DefaultErrorMsg, id, action, tferrs.AlibabaCloudSdkGoERROR)
	}

	v, err := jsonpath.Get("$.Databases.Database[*]", response)
	if err != nil {
		return object, tferrs.WrapErrorf(tferrs.NotFoundErr("Database", id), tferrs.NotFoundMsg, response)
	}

	if len(v.([]interface{})) == 0 {
		return object, tferrs.WrapErrorf(tferrs.NotFoundErr("Database", id), tferrs.NotFoundMsg, response)
	}

	currentStatus := v.([]interface{})[0].(map[string]interface{})["DBName"]
	if currentStatus == nil {
		return object, tferrs.WrapErrorf(tferrs.NotFoundErr("Database", id), tferrs.NotFoundMsg, response)
	}

	return v.([]interface{})[0].(map[string]interface{}), nil
}

func (s *PolardbServiceV2) PolardbDatabaseStateRefreshFunc(id string, field string, failStates []string) tfretry.StateRefreshFunc {
	return s.PolardbDatabaseStateRefreshFuncWithApi(id, field, failStates, s.DescribePolardbDatabase)
}

func (s *PolardbServiceV2) PolardbDatabaseStateRefreshFuncWithApi(id string, field string, failStates []string, call func(id string) (map[string]interface{}, error)) tfretry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := call(id)
		if err != nil {
			if tferrs.NotFoundError(err) {
				return object, "", nil
			}
			return nil, "", tferrs.WrapError(err)
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
				return object, currentStatus, tferrs.WrapError(tferrs.Error(tferrs.FailedToReachTargetStatus, currentStatus))
			}
		}
		return object, currentStatus, nil
	}
}

// DescribePolardbDatabase >>> Encapsulated.
