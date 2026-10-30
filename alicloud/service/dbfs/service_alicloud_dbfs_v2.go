package dbfs

import (
	"fmt"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
)

type ServiceV2 struct {
	client *connectivity.AliyunClient
}

// NewServiceV2 is the exported constructor for the v2 service variant.
func NewServiceV2(client *connectivity.AliyunClient) *ServiceV2 {
	return &ServiceV2{client: client}
}

// DescribeDbfsDbfsInstance <<< Encapsulated get interface for Dbfs DbfsInstance.

func (s *ServiceV2) DescribeDbfsDbfsInstance(id string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	action := "GetDbfs"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	query["FsId"] = id
	request["RegionId"] = client.RegionId
	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(1*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("DBFS", "2020-04-18", action, query, request, true)

		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		logs.AddDebug(action, response, request)
		return nil
	})

	if err != nil {
		if errs.IsExpectedErrors(err, []string{"EntityNotExist.DBFS"}) {
			return object, errs.WrapErrorf(errs.NotFoundErr("DbfsInstance", id), errs.NotFoundMsg, response)
		}
		return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}

	v, err := jsonpath.Get("$.DBFSInfo", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.DBFSInfo", response)
	}

	return v.(map[string]interface{}), nil
}

func (s *ServiceV2) DbfsDbfsInstanceStateRefreshFunc(id string, field string, failStates []string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeDbfsDbfsInstance(id)
		if err != nil {
			if errs.NotFoundError(err) {
				return nil, "", nil
			}
			return nil, "", errs.WrapError(err)
		}

		v, err := jsonpath.Get(field, object)
		currentStatus := fmt.Sprint(v)

		for _, failState := range failStates {
			if currentStatus == failState {
				return object, currentStatus, errs.WrapError(errs.Error(errs.FailedToReachTargetStatus, currentStatus))
			}
		}
		return object, currentStatus, nil
	}
}

// DescribeDbfsDbfsInstance >>> Encapsulated.
