package alicloud

import (
	"fmt"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

type DbsService struct {
	client *connectivity.AliyunClient
}

func (s *DbsService) DescribeBackupPlanBilling(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	client := s.client
	action := "DescribeBackupPlanBilling"
	request := map[string]interface{}{
		"BackupPlanId": id,
	}
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = client.RpcPost("Dbs", "2019-03-06", action, nil, request, true)
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
		if IsExpectedErrors(err, []string{"InvalidJobId"}) {
			return object, WrapErrorf(NotFoundErr("DBS:BackupPlan", id), NotFoundMsg, ProviderERROR, fmt.Sprint(response["RequestId"]))
		}
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}
	v, err := jsonpath.Get("$.Item", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.Item", response)
	}
	object = v.(map[string]interface{})
	return object, nil
}

func (s *DbsService) DescribeDbsBackupPlan(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	client := s.client
	action := "DescribeBackupPlanList"
	request := map[string]interface{}{
		"Region":       s.client.RegionId,
		"BackupPlanId": id,
	}
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = client.RpcPost("Dbs", "2019-03-06", action, nil, request, true)
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
		if IsExpectedErrors(err, []string{"InvalidParameterValid"}) {
			return object, WrapErrorf(NotFoundErr("DBS:BackupPlan", id), NotFoundMsg, ProviderERROR, fmt.Sprint(response["RequestId"]))
		}
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}
	v, err := jsonpath.Get("$.Items.BackupPlanDetail", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.Items.BackupPlanDetail", response)
	}
	if len(v.([]interface{})) < 1 {
		return object, WrapErrorf(NotFoundErr("DBS", id), NotFoundWithResponse, response)
	} else {
		if fmt.Sprint(v.([]interface{})[0].(map[string]interface{})["BackupPlanId"]) != id {
			return object, WrapErrorf(NotFoundErr("DBS", id), NotFoundWithResponse, response)
		}
	}
	object = v.([]interface{})[0].(map[string]interface{})
	return object, nil
}

func (s *DbsService) DbsBackupPlanStateRefreshFunc(id string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeDbsBackupPlan(id)
		if err != nil {
			if NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", WrapError(err)
		}

		for _, failState := range failStates {
			if fmt.Sprint(object["BackupPlanStatus"]) == failState {
				return object, fmt.Sprint(object["BackupPlanStatus"]), WrapError(Error(FailedToReachTargetStatus, fmt.Sprint(object["BackupPlanStatus"])))
			}
		}
		return object, fmt.Sprint(object["BackupPlanStatus"]), nil
	}
}
