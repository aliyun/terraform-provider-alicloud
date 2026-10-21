package oss

import (
	"fmt"
	"strings"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
)

const accessPointAction = "/?accessPoint"

func (s *Service) DescribeAccessPoint(id string) (object map[string]interface{}, err error) {
	var request map[string]interface{}
	var response map[string]interface{}
	parts := strings.Split(id, ":")
	if len(parts) != 2 {
		return object, errs.WrapError(fmt.Errorf("invalid Resource Id %s. Expected parts' length %d, got %d", id, 2, len(parts)))
	}
	request = make(map[string]interface{})
	query := make(map[string]*string)
	hostMap := make(map[string]*string)
	query["x-oss-access-point-name"] = helper.StringPointer(parts[1])
	hostMap["bucket"] = helper.StringPointer(parts[0])
	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(1*time.Minute, func() *retry.RetryError {
		response, err = s.client.Do("Oss", helper.XmlParam("GET", "2019-05-17", "GetAccessPoint", accessPointAction), query, nil, nil, hostMap, true)
		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(accessPointAction, response, request)
	if err != nil {
		if errs.IsExpectedErrors(err, []string{"NoSuchAccessPoint"}) {
			return object, errs.WrapErrorf(errs.NotFoundErr("AccessPoint", id), errs.NotFoundMsg, response)
		}
		return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, accessPointAction, errs.AlibabaCloudSdkGoERROR)
	}
	if response == nil {
		return object, errs.WrapErrorf(errs.NotFoundErr("AccessPoint", id), errs.NotFoundMsg, response)
	}

	v, err := jsonpath.Get("$.GetAccessPointResult", response)
	if err != nil {
		return object, errs.WrapErrorf(errs.NotFoundErr("AccessPoint", id), errs.NotFoundMsg, response)
	}

	return v.(map[string]interface{}), nil
}

func (s *Service) AccessPointStateRefreshFunc(id string, field string, failStates []string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeAccessPoint(id)
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
