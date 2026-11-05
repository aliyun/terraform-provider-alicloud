package hologram

import (
	"fmt"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	logs "github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
)

// holoWebLoginSettingTimeout replaces the SDKv2 resource-level default timeouts
// (Create/Update/Delete each 5 minutes); the Framework resource has no
// user-configurable timeouts yet.
const holoWebLoginSettingTimeout = 5 * time.Minute

// DescribeHoloWebLoginSetting queries whether the HoloWeb login setting allows
// external accounts to log in to the given Hologres instance. The setting is a
// single boolean that always exists for an existing instance, so NotFound here
// means the underlying instance is gone.
func (s *Service) DescribeHoloWebLoginSetting(instanceId string) (map[string]interface{}, error) {
	action := fmt.Sprintf("/api/v1/instances/%s/getHoloWebLoginSetting", instanceId)
	query := make(map[string]*string)
	body := make(map[string]interface{})
	var response map[string]interface{}
	var err error
	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(1*time.Minute, func() *retry.RetryError {
		response, err = s.client.RoaGet("Hologram", "2022-06-01", action, query, nil, body)
		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(action, response, query)
	if err != nil {
		if errs.IsExpectedErrors(err, []string{"resource not exists failed"}) {
			return nil, errs.WrapErrorf(errs.NotFoundErr("HoloWebLoginSetting", instanceId), errs.NotFoundMsg, response)
		}
		return nil, errs.WrapErrorf(err, errs.DefaultErrorMsg, instanceId, action, errs.AlibabaCloudSdkGoERROR)
	}
	data, ok := response["Data"].(map[string]interface{})
	if !ok {
		return nil, errs.WrapErrorf(fmt.Errorf("the response Data of %s is not an object: %v", action, response["Data"]), errs.DefaultErrorMsg, instanceId, action, errs.AlibabaCloudSdkGoERROR)
	}
	return data, nil
}

// UpdateHoloWebLoginSetting toggles the external accounts login setting of the
// given Hologres instance. Create and Update pass the configured value; Delete
// reverts it to false, which is the service default.
func (s *Service) UpdateHoloWebLoginSetting(instanceId string, allowExternalAccountsLogin bool) error {
	action := fmt.Sprintf("/api/v1/instances/%s/updateHoloWebLoginSetting", instanceId)
	query := make(map[string]*string)
	body := map[string]interface{}{
		"allowExternalAccountsLogin": allowExternalAccountsLogin,
	}
	var response map[string]interface{}
	var err error
	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(holoWebLoginSettingTimeout, func() *retry.RetryError {
		response, err = s.client.RoaPost("Hologram", "2022-06-01", action, query, nil, body, true)
		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(action, response, body)
	if err != nil {
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, instanceId, action, errs.AlibabaCloudSdkGoERROR)
	}
	return nil
}
