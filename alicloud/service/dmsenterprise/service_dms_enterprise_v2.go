package dmsenterprise

import (
	"fmt"
	"strings"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	tferrs "github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	tfhelper "github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	tflogs "github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
)

type DMSEnterpriseServiceV2 struct {
	client *connectivity.AliyunClient
}

// DescribeDMSEnterpriseAuthorityTemplate <<< Encapsulated get interface for DMSEnterprise AuthorityTemplate.

func (s *DMSEnterpriseServiceV2) DescribeDMSEnterpriseAuthorityTemplate(id string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	parts := strings.Split(id, ":")
	if len(parts) != 2 {
		err = tferrs.WrapError(fmt.Errorf("invalid Resource Id %s. Expected parts' length %d, got %d", id, 2, len(parts)))
	}
	action := "GetAuthorityTemplate"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	query["TemplateId"] = parts[1]
	query["Tid"] = parts[0]

	wait := tfhelper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(1*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("dms-enterprise", "2018-11-01", action, query, request, true)

		if err != nil {
			if tfhelper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		tflogs.AddDebug(action, response, request)
		return nil
	})

	if err != nil {
		return object, tferrs.WrapErrorf(err, tferrs.DefaultErrorMsg, id, action, tferrs.AlibabaCloudSdkGoERROR)
	}

	if _, ok := response["AuthorityTemplateView"]; !ok {
		return object, tferrs.WrapErrorf(tferrs.NotFoundErr("AuthorityTemplate", id), tferrs.NotFoundMsg, response)
	}

	return response, nil
}

func (s *DMSEnterpriseServiceV2) DMSEnterpriseAuthorityTemplateStateRefreshFunc(id string, field string, failStates []string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeDMSEnterpriseAuthorityTemplate(id)
		if err != nil {
			if tferrs.NotFoundError(err) {
				return object, "", nil
			}
			return nil, "", tferrs.WrapError(err)
		}

		v, err := jsonpath.Get(field, object)
		currentStatus := fmt.Sprint(v)
		for _, failState := range failStates {
			if currentStatus == failState {
				return object, currentStatus, tferrs.WrapError(tferrs.Error(tferrs.FailedToReachTargetStatus, currentStatus))
			}
		}
		return object, currentStatus, nil
	}
}

// DescribeDMSEnterpriseAuthorityTemplate >>> Encapsulated.

// DescribeDmsEnterpriseWorkspace <<< Encapsulated get interface for DmsEnterprise Workspace.

func (s *DMSEnterpriseServiceV2) DescribeDmsEnterpriseWorkspace(id string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["WorkspaceId"] = id
	request["RegionId"] = client.RegionId
	action := "GetWorkspace"

	wait := tfhelper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(1*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("dms-enterprise", "2018-11-01", action, query, request, true)

		if err != nil {
			if tfhelper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	tflogs.AddDebug(action, response, request)
	if err != nil {
		if tferrs.IsExpectedErrors(err, []string{"WorkspaceNotExist"}) {
			return object, tferrs.WrapErrorf(tferrs.NotFoundErr("Workspace", id), tferrs.NotFoundMsg, response)
		}
		return object, tferrs.WrapErrorf(err, tferrs.DefaultErrorMsg, id, action, tferrs.AlibabaCloudSdkGoERROR)
	}

	v, err := jsonpath.Get("$.Data", response)
	if err != nil {
		return object, tferrs.WrapErrorf(err, tferrs.FailedGetAttributeMsg, id, "$.Data", response)
	}

	return v.(map[string]interface{}), nil
}

func (s *DMSEnterpriseServiceV2) DmsEnterpriseWorkspaceStateRefreshFunc(id string, field string, failStates []string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeDmsEnterpriseWorkspace(id)
		if err != nil {
			if tferrs.NotFoundError(err) {
				return object, "", nil
			}
			return nil, "", tferrs.WrapError(err)
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
				return object, currentStatus, tferrs.WrapError(tferrs.Error(tferrs.FailedToReachTargetStatus, currentStatus))
			}
		}
		return object, currentStatus, nil
	}
}

// DescribeDmsEnterpriseWorkspace >>> Encapsulated.
