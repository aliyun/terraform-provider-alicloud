// Package swasopen. This file is generated automatically. Please do not modify it manually, thank you!
package swasopen

import (
	"fmt"
	"strings"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/convert"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
)

type SimpleApplicationServerServiceV2 struct {
	client *connectivity.AliyunClient
}

// DescribeSimpleApplicationServerDisk <<< Encapsulated get interface for SimpleApplicationServer Disk.

func (s *SimpleApplicationServerServiceV2) DescribeSimpleApplicationServerDisk(id string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["DiskIds"] = convert.ExpandSingletonToList(id)
	request["RegionId"] = client.RegionId
	action := "ListDisks"

	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(1*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("SWAS-OPEN", "2020-06-01", action, query, request, true)

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

	v, err := jsonpath.Get("$.Disks[*]", response)
	if err != nil {
		return object, errs.WrapErrorf(errs.NotFoundErr("Disk", id), errs.NotFoundMsg, response)
	}

	if len(v.([]interface{})) == 0 {
		return object, errs.WrapErrorf(errs.NotFoundErr("Disk", id), errs.NotFoundMsg, response)
	}

	currentStatus := v.([]interface{})[0].(map[string]interface{})["DiskId"]
	if currentStatus == nil {
		return object, errs.WrapErrorf(errs.NotFoundErr("Disk", id), errs.NotFoundMsg, response)
	}

	return v.([]interface{})[0].(map[string]interface{}), nil
}

func (s *SimpleApplicationServerServiceV2) SimpleApplicationServerDiskStateRefreshFunc(id string, field string, failStates []string) retry.StateRefreshFunc {
	return s.SimpleApplicationServerDiskStateRefreshFuncWithApi(id, field, failStates, s.DescribeSimpleApplicationServerDisk)
}

func (s *SimpleApplicationServerServiceV2) SimpleApplicationServerDiskStateRefreshFuncWithApi(id string, field string, failStates []string, call func(id string) (map[string]interface{}, error)) retry.StateRefreshFunc {
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
				return object, currentStatus, errs.WrapError(errs.Error(errs.FailedToReachTargetStatus, currentStatus))
			}
		}
		return object, currentStatus, nil
	}
}

// DescribeSimpleApplicationServerDisk >>> Encapsulated.
