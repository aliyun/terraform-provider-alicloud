package cloudmonitorservice

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/cms"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/convert"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/tags"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// DefaultTimeout mirrors the root package constant used by CmsService wait
// helpers; ported verbatim from alicloud/common.go.
const DefaultTimeout = 120

type CmsService struct {
	client *connectivity.AliyunClient
}

type IspCities []map[string]string

func (s *CmsService) BuildCmsCommonRequest(region string) *requests.CommonRequest {
	request := requests.NewCommonRequest()
	return request
}

func (s *CmsService) BuildCmsAlarmRequest(id string) *requests.CommonRequest {

	request := s.BuildCmsCommonRequest(s.client.RegionId)
	request.QueryParams["Id"] = id

	return request
}

func (s *CmsService) DescribeAlarm(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	action := "DescribeMetricRuleList"
	client := s.client

	request := map[string]interface{}{
		"RuleIds": id,
	}

	idExist := false
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(6*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("Cms", "2019-01-01", action, nil, request, false)
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
		if errs.IsExpectedErrors(err, []string{"InternalError", "ResourceNotFound"}) {
			return object, errs.WrapErrorf(errs.NotFoundErr("Cms:Alarm", id), errs.NotFoundMsg, errs.ProviderERROR, fmt.Sprint(response["RequestId"]))
		}
		return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}

	resp, err := jsonpath.Get("$.Alarms.Alarm", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.Alarms.Alarm", response)
	}

	if v, ok := resp.([]interface{}); !ok || len(v) < 1 {
		return object, errs.WrapErrorf(errs.NotFoundErr("Cms:Alarm", id), errs.NotFoundWithResponse, response)
	}

	for _, v := range resp.([]interface{}) {
		if fmt.Sprint(v.(map[string]interface{})["RuleId"]) == id {
			idExist = true
			return v.(map[string]interface{}), nil
		}
	}

	if !idExist {
		return object, errs.WrapErrorf(errs.NotFoundErr("Cms:Alarm", id), errs.NotFoundWithResponse, response)
	}

	return object, nil
}

func (s *CmsService) WaitForCmsAlarm(id string, enabled bool, timeout int) error {
	if timeout <= 0 {
		timeout = DefaultTimeout
	}

	for {
		alarm, err := s.DescribeAlarm(id)
		if err != nil {
			return err
		}

		if alarm["EnableState"] == enabled {
			break
		}
		timeout = timeout - helper.DefaultIntervalShort
		if timeout <= 0 {
			return errs.GetTimeErrorFromString(errs.GetTimeoutMessage("Alarm", strconv.FormatBool(enabled)))
		}
		time.Sleep(helper.DefaultIntervalShort * time.Second)
	}
	return nil
}

func (s *CmsService) BuildJsonWebhook(webhook string) string {
	if webhook != "" {
		return fmt.Sprintf("{\"method\":\"post\",\"url\":\"%s\"}", webhook)
	}
	return ""
}

func (s *CmsService) ExtractWebhookFromJson(webhookJson string) (string, error) {
	byt := []byte(webhookJson)
	var dat map[string]interface{}
	if err := json.Unmarshal(byt, &dat); err != nil {
		return "", err
	}
	return dat["url"].(string), nil
}

func (s *CmsService) DescribeSiteMonitor(id, keyword string) (siteMonitor cms.SiteMonitor, err error) {
	listRequest := cms.CreateDescribeSiteMonitorListRequest()
	listRequest.Keyword = keyword
	listRequest.TaskId = id
	var raw interface{}
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		raw, err = s.client.WithCmsClient(func(cmsClient *cms.Client) (interface{}, error) {
			return cmsClient.DescribeSiteMonitorList(listRequest)
		})
		if err != nil {
			if errs.IsExpectedErrors(err, []string{errs.ThrottlingUser, "ExceedingQuota"}) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	if err != nil {
		return siteMonitor, errs.WrapError(err)
	}
	list := raw.(*cms.DescribeSiteMonitorListResponse)
	if len(list.SiteMonitors.SiteMonitor) < 1 {
		return siteMonitor, errs.GetNotFoundErrorFromString(errs.GetNotFoundMessage("Site Monitor", id))

	}
	for _, v := range list.SiteMonitors.SiteMonitor {
		if v.TaskName == keyword || v.TaskId == id {
			return v, nil
		}
	}
	return siteMonitor, errs.GetNotFoundErrorFromString(errs.GetNotFoundMessage("Site Monitor", id))
}

func (s *CmsService) GetIspCities(id string) (ispCities IspCities, err error) {
	request := cms.CreateDescribeSiteMonitorAttributeRequest()
	request.TaskId = id

	var raw interface{}
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		raw, err = s.client.WithCmsClient(func(cmsClient *cms.Client) (interface{}, error) {
			return cmsClient.DescribeSiteMonitorAttribute(request)
		})
		if err != nil {
			if errs.IsExpectedErrors(err, []string{errs.ThrottlingUser, "ExceedingQuota"}) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	if err != nil {
		return nil, errs.WrapError(err)
	}

	response := raw.(*cms.DescribeSiteMonitorAttributeResponse)
	ispCity := response.SiteMonitors.IspCities.IspCity

	var list []map[string]string
	for _, element := range ispCity {
		list = append(list, map[string]string{"city": element.City, "isp": element.Isp})
	}

	return list, nil
}

func (s *CmsService) DescribeCmsAlarmContact(id string) (object cms.Contact, err error) {
	request := cms.CreateDescribeContactListRequest()
	request.RegionId = s.client.RegionId

	request.ContactName = id
	var response *cms.DescribeContactListResponse
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		raw, err := s.client.WithCmsClient(func(cmsClient *cms.Client) (interface{}, error) {
			return cmsClient.DescribeContactList(request)
		})
		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		logs.AddDebug(request.GetActionName(), raw, request.RpcRequest, request)
		response, _ = raw.(*cms.DescribeContactListResponse)
		return nil
	})
	if err != nil {
		if errs.IsExpectedErrors(err, []string{"ContactNotExists", "ResourceNotFound"}) {
			err = errs.WrapErrorf(errs.NotFoundErr("CmsAlarmContact", id), errs.NotFoundMsg, errs.ProviderERROR)
			return
		}
		err = errs.WrapErrorf(err, errs.DefaultErrorMsg, id, request.GetActionName(), errs.AlibabaCloudSdkGoERROR)
		return
	}

	if response.Code != "200" {
		err = errs.Error("DescribeContactList failed for %s", response.Message)
		return
	}

	if len(response.Contacts.Contact) < 1 {
		err = errs.WrapErrorf(errs.NotFoundErr("CmsAlarmContact", id), errs.NotFoundMsg, errs.ProviderERROR, response.RequestId)
		return
	}
	return response.Contacts.Contact[0], nil
}

func (s *CmsService) DescribeCmsAlarmContactGroup(id string) (object cms.ContactGroup, err error) {
	request := cms.CreateDescribeContactGroupListRequest()
	request.RegionId = s.client.RegionId

	request.PageNumber = requests.NewInteger(1)
	request.PageSize = requests.NewInteger(20)
	for {

		var raw interface{}
		wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
		err = retry.Retry(5*time.Minute, func() *retry.RetryError {
			raw, err = s.client.WithCmsClient(func(cmsClient *cms.Client) (interface{}, error) {
				return cmsClient.DescribeContactGroupList(request)
			})
			if err != nil {
				if helper.NeedRetry(err) {
					wait()
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
			}
			return nil
		})

		if err != nil {
			if errs.IsExpectedErrors(err, []string{"ContactGroupNotExists", "ResourceNotFound"}) {
				err = errs.WrapErrorf(errs.NotFoundErr("CmsAlarmContactGroup", id), errs.NotFoundMsg, errs.ProviderERROR)
				return object, err
			}
			err = errs.WrapErrorf(err, errs.DefaultErrorMsg, id, request.GetActionName(), errs.AlibabaCloudSdkGoERROR)
			return object, err
		}
		logs.AddDebug(request.GetActionName(), raw, request.RpcRequest, request)
		response, _ := raw.(*cms.DescribeContactGroupListResponse)
		if response.Code != "200" {
			err = errs.Error("DescribeContactGroupList failed for %s", response.Message)
			return object, err
		}

		if len(response.ContactGroupList.ContactGroup) < 1 {
			err = errs.WrapErrorf(errs.NotFoundErr("CmsAlarmContactGroup", id), errs.NotFoundMsg, errs.ProviderERROR, response.RequestId)
			return object, err
		}
		for _, object := range response.ContactGroupList.ContactGroup {
			if object.Name == id {
				return object, nil
			}
		}
		if len(response.ContactGroupList.ContactGroup) < helper.PageSizeMedium {
			break
		}
		if page, err := helper.GetNextpageNumber(request.PageNumber); err != nil {
			return object, errs.WrapError(err)
		} else {
			request.PageNumber = page
		}
	}
	err = errs.WrapErrorf(errs.NotFoundErr("CmsAlarmContactGroup", id), errs.NotFoundMsg, errs.ProviderERROR)
	return
}

func (s *CmsService) DescribeCmsGroupMetricRule(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	action := "DescribeMetricRuleList"
	client := s.client

	request := map[string]interface{}{
		"RuleIds": id,
	}

	idExist := false
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(6*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("Cms", "2019-01-01", action, nil, request, false)
		if err != nil {
			if errs.IsExpectedErrors(err, []string{"ExceedingQuota"}) || helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(action, response, request)

	if err != nil {
		if errs.IsExpectedErrors(err, []string{"GroupMetricRuleNotExists", "ResourceNotFound", "ResourceNotFoundError"}) {
			return object, errs.WrapErrorf(errs.NotFoundErr("Cms:GroupMetricRule", id), errs.NotFoundMsg, errs.ProviderERROR, fmt.Sprint(response["RequestId"]))
		}
		return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}

	resp, err := jsonpath.Get("$.Alarms.Alarm", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.Alarms.Alarm", response)
	}

	if v, ok := resp.([]interface{}); !ok || len(v) < 1 {
		return object, errs.WrapErrorf(errs.NotFoundErr("Cms:GroupMetricRule", id), errs.NotFoundWithResponse, response)
	}

	for _, v := range resp.([]interface{}) {
		if fmt.Sprint(v.(map[string]interface{})["RuleId"]) == id {
			idExist = true
			return v.(map[string]interface{}), nil
		}
	}

	if !idExist {
		return object, errs.WrapErrorf(errs.NotFoundErr("Cms:GroupMetricRule", id), errs.NotFoundWithResponse, response)
	}

	return object, nil
}

func (s *CmsService) SetResourceTags(d *schema.ResourceData, resourceType string) error {
	client := s.client
	if d.HasChange("tags") {
		added, removed := tags.ParsingTags(d)
		removedTagKeys := make([]string, 0)
		for _, v := range removed {
			if !tags.IgnoredTags(v, "") {
				removedTagKeys = append(removedTagKeys, v)
			}
		}
		if len(removedTagKeys) > 0 {
			action := "RemoveTags"
			request := map[string]interface{}{
				"RegionId":   s.client.RegionId,
				"GroupIds.1": d.Id(),
			}
			oraw, _ := d.GetChange("tags")
			removedTags := oraw.(map[string]interface{})
			count := 1
			for _, key := range removedTagKeys {
				request[fmt.Sprintf("Tag.%d.Key", count)] = key
				request[fmt.Sprintf("Tag.%d.Value", count)] = removedTags[key]
			}
			wait := helper.IncrementalWait(2*time.Second, 1*time.Second)
			err := retry.Retry(10*time.Minute, func() *retry.RetryError {
				response, err := client.RpcPost("Cms", "2019-01-01", action, nil, request, false)
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
				return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
			}
		}
		if len(added) > 0 {
			action := "AddTags"
			request := map[string]interface{}{
				"RegionId":   s.client.RegionId,
				"GroupIds.1": d.Id(),
			}
			count := 1
			for key, value := range added {
				request[fmt.Sprintf("Tag.%d.Key", count)] = key
				request[fmt.Sprintf("Tag.%d.Value", count)] = value
				count++
			}

			wait := helper.IncrementalWait(2*time.Second, 1*time.Second)
			err := retry.Retry(10*time.Minute, func() *retry.RetryError {
				response, err := client.RpcPost("Cms", "2019-01-01", action, nil, request, false)
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
				return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
			}
		}
	}
	return nil
}

func (s *CmsService) DescribeCmsMonitorGroup(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	client := s.client
	action := "DescribeMonitorGroups"
	request := map[string]interface{}{
		"RegionId": s.client.RegionId,
		"GroupId":  id,
	}
	response, err = client.RpcPost("Cms", "2019-01-01", action, nil, request, false)
	if err != nil {
		err = errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
		return
	}
	logs.AddDebug(action, response, request)
	if errs.IsExpectedErrorCodes(fmt.Sprintf("%v", response["Code"]), []string{"GroupsNotExists", "ResourceNotFound"}) {
		err = errs.WrapErrorf(errs.NotFoundErr("CmsMonitorGroup", id), errs.NotFoundMsg, errs.ProviderERROR)
		return object, err
	}
	if fmt.Sprintf(`%v`, response["Code"]) != "200" {
		err = errs.Error("DescribeMonitorGroups failed for %s", response["Message"].(string))
		return object, err
	}
	v, err := jsonpath.Get("$.Resources.Resource", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.Resources.Resource", response)
	}
	if len(v.([]interface{})) < 1 {
		return object, errs.WrapErrorf(errs.NotFoundErr("CloudMonitorService", id), errs.NotFoundWithResponse, response)
	} else {
		groupId, err := convert.FormatInt(v.([]interface{})[0].(map[string]interface{})["GroupId"])
		if err != nil {
			return nil, err
		}
		if fmt.Sprint(groupId) != id {
			return object, errs.WrapErrorf(errs.NotFoundErr("CloudMonitorService", id), errs.NotFoundWithResponse, response)
		}
	}
	object = v.([]interface{})[0].(map[string]interface{})
	return object, nil
}

func (s *CmsService) DescribeCmsMonitorGroupInstances(id string) (object []map[string]interface{}, err error) {
	var response map[string]interface{}
	client := s.client
	action := "DescribeMonitorGroupInstances"
	request := map[string]interface{}{
		"RegionId": s.client.RegionId,
		"GroupId":  id,
	}
	request["PageSize"] = helper.PageSizeMedium
	request["PageNumber"] = 1
	for {
		wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
		err = retry.Retry(5*time.Minute, func() *retry.RetryError {
			response, err = client.RpcPost("Cms", "2019-01-01", action, nil, request, false)
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
			if errs.IsExpectedErrors(err, []string{"ResourceNotFound", "ResourceNotFoundError"}) {
				err = errs.WrapErrorf(errs.NotFoundErr("CmsMonitorGroupInstances", id), errs.NotFoundMsg, errs.ProviderERROR)
				return object, err
			}
			return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
		}
		v, err := jsonpath.Get("$.Resources.Resource", response)
		if err != nil {
			return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.Resources.Resource", response)
		}
		if len(v.([]interface{})) < 1 {
			break
		}

		for _, v := range v.([]interface{}) {
			object = append(object, v.(map[string]interface{}))
		}
		if len(v.([]interface{})) < request["PageSize"].(int) {
			break
		}
		request["PageNumber"] = request["PageNumber"].(int) + 1
	}

	if len(object) < 1 {
		return object, errs.WrapErrorf(errs.NotFoundErr("CmsMonitorGroupInstances", id), errs.NotFoundWithResponse, response)
	}

	return object, nil
}

func (s *CmsService) DescribeCmsMetricRuleTemplate(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	action := "DescribeMetricRuleTemplateAttribute"
	client := s.client

	request := map[string]interface{}{
		"TemplateId": id,
	}

	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("Cms", "2019-01-01", action, nil, request, false)
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
		if errs.IsExpectedErrors(err, []string{"ResourceNotFound"}) {
			return object, errs.WrapErrorf(errs.NotFoundErr("Cms:MetricRuleTemplate", id), errs.NotFoundMsg, errs.ProviderERROR, fmt.Sprint(response["RequestId"]))
		}
		return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}

	v, err := jsonpath.Get("$.Resource", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.Resource", response)
	}

	object = v.(map[string]interface{})

	return object, nil
}

func (s *CmsService) DescribeCmsDynamicTagGroup(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	action := "DescribeDynamicTagRuleList"
	client := s.client

	request := map[string]interface{}{
		"TagRegionId":      s.client.RegionId,
		"DynamicTagRuleId": id,
		"PageSize":         helper.PageSizeLarge,
		"PageNumber":       1,
	}

	idExist := false
	for {
		wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
		err = retry.Retry(5*time.Minute, func() *retry.RetryError {
			response, err = client.RpcPost("Cms", "2019-01-01", action, nil, request, false)
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

		resp, err := jsonpath.Get("$.TagGroupList.TagGroup", response)
		if err != nil {
			return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.TagGroupList.TagGroup", response)
		}

		if v, ok := resp.([]interface{}); !ok || len(v) < 1 {
			return object, errs.WrapErrorf(errs.NotFoundErr("Cms:DynamicTagGroup", id), errs.NotFoundWithResponse, response)
		}

		for _, v := range resp.([]interface{}) {
			if fmt.Sprint(v.(map[string]interface{})["DynamicTagRuleId"]) == id {
				idExist = true
				return v.(map[string]interface{}), nil
			}
		}

		if len(resp.([]interface{})) < request["PageSize"].(int) {
			break
		}

		request["PageNumber"] = request["PageNumber"].(int) + 1
	}

	if !idExist {
		return object, errs.WrapErrorf(errs.NotFoundErr("Cms:DynamicTagGroup", id), errs.NotFoundWithResponse, response)
	}

	return object, nil
}

func (s *CmsService) DescribeCmsNamespace(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	action := "DescribeHybridMonitorNamespaceList"
	client := s.client

	request := map[string]interface{}{
		"Namespace":  id,
		"PageSize":   helper.PageSizeLarge,
		"PageNumber": 1,
	}

	idExist := false
	for {
		wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
		err = retry.Retry(5*time.Minute, func() *retry.RetryError {
			response, err = client.RpcPost("Cms", "2019-01-01", action, nil, request, false)
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

		resp, err := jsonpath.Get("$.DescribeHybridMonitorNamespace", response)
		if err != nil {
			return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.DescribeHybridMonitorNamespace", response)
		}

		if v, ok := resp.([]interface{}); !ok || len(v) < 1 {
			return object, errs.WrapErrorf(errs.NotFoundErr("CloudMonitorService:Namespace", id), errs.NotFoundWithResponse, response)
		}

		for _, v := range resp.([]interface{}) {
			if fmt.Sprint(v.(map[string]interface{})["Namespace"]) == id {
				idExist = true
				return v.(map[string]interface{}), nil
			}
		}

		if len(resp.([]interface{})) < request["PageSize"].(int) {
			break
		}

		request["PageNumber"] = request["PageNumber"].(int) + 1
	}

	if !idExist {
		return object, errs.WrapErrorf(errs.NotFoundErr("CloudMonitorService:Namespace", id), errs.NotFoundWithResponse, response)
	}

	return object, nil
}

func (s *CmsService) DescribeCmsSlsGroup(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	client := s.client
	action := "DescribeHybridMonitorSLSGroup"
	request := map[string]interface{}{
		"SLSGroupName": id,
	}
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("Cms", "2019-01-01", action, nil, request, false)
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
	if errs.IsExpectedErrorCodes(fmt.Sprint(response["Code"]), []string{"400"}) {
		return object, errs.WrapErrorf(errs.NotFoundErr("CloudMonitorService:SlsGroup", id), errs.NotFoundMsg, errs.ProviderERROR)
	}
	v, err := jsonpath.Get("$.List", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.List", response)
	}
	if len(v.([]interface{})) < 1 {
		return object, errs.WrapErrorf(errs.NotFoundErr("CloudMonitorService", id), errs.NotFoundWithResponse, response)
	} else {
		if fmt.Sprint(v.([]interface{})[0].(map[string]interface{})["SLSGroupName"]) != id {
			return object, errs.WrapErrorf(errs.NotFoundErr("CloudMonitorService", id), errs.NotFoundWithResponse, response)
		}
	}
	object = v.([]interface{})[0].(map[string]interface{})
	return object, nil
}

func (s *CmsService) DescribeCmsHybridMonitorSlsTask(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	client := s.client
	action := "DescribeHybridMonitorTaskList"
	request := map[string]interface{}{
		"TaskId": id,
	}
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("Cms", "2019-01-01", action, nil, request, false)
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
	if errs.IsExpectedErrorCodes(fmt.Sprint(response["Code"]), []string{"ResourceNotFound"}) {
		return object, errs.WrapErrorf(errs.NotFoundErr("CloudMonitorService:HybridMonitorSlsTask", id), errs.NotFoundMsg, errs.ProviderERROR)
	}
	v, err := jsonpath.Get("$.TaskList", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.TaskList", response)
	}
	if len(v.([]interface{})) < 1 {
		return object, errs.WrapErrorf(errs.NotFoundErr("CloudMonitorService", id), errs.NotFoundWithResponse, response)
	} else {
		if fmt.Sprint(v.([]interface{})[0].(map[string]interface{})["TaskId"]) != id {
			return object, errs.WrapErrorf(errs.NotFoundErr("CloudMonitorService", id), errs.NotFoundWithResponse, response)
		}
	}
	object = v.([]interface{})[0].(map[string]interface{})
	return object, nil
}

func (s *CmsService) DescribeCmsHybridMonitorFcTask(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	client := s.client
	action := "DescribeHybridMonitorTaskList"
	parts, err := helper.ParseResourceId(id, 2)
	if err != nil {
		err = errs.WrapError(err)
		return
	}
	request := map[string]interface{}{
		"TargetUserId": parts[0],
		"Namespace":    parts[1],
		"TaskType":     "aliyun_fc",
	}
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("Cms", "2019-01-01", action, nil, request, false)
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
	if errs.IsExpectedErrorCodes(fmt.Sprint(response["Code"]), []string{"ResourceNotFound"}) {
		return object, errs.WrapErrorf(errs.NotFoundErr("CloudMonitorService:HybridMonitorFcTask", id), errs.NotFoundMsg, errs.ProviderERROR)
	}
	v, err := jsonpath.Get("$.TaskList", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.TaskList", response)
	}
	if len(v.([]interface{})) < 1 {
		return object, errs.WrapErrorf(errs.NotFoundErr("CloudMonitorService", id), errs.NotFoundWithResponse, response)
	}
	object = v.([]interface{})[0].(map[string]interface{})
	return object, nil
}

func (s *CmsService) DescribeCmsEventRule(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	action := "DescribeEventRuleList"
	client := s.client

	request := map[string]interface{}{
		"NamePrefix": id,
	}

	idExist := false
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("Cms", "2019-01-01", action, nil, request, false)
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
	resp, err := jsonpath.Get("$.EventRules.EventRule", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.EventRules.EventRule", response)
	}

	if v, ok := resp.([]interface{}); !ok || len(v) < 1 {
		return object, errs.WrapErrorf(errs.NotFoundErr("Cms:EventRule", id), errs.NotFoundWithResponse, response)
	}

	for _, v := range resp.([]interface{}) {
		if fmt.Sprint(v.(map[string]interface{})["Name"]) == id {
			idExist = true
			return v.(map[string]interface{}), nil
		}
	}

	if !idExist {
		return object, errs.WrapErrorf(errs.NotFoundErr("Cms:EventRule", id), errs.NotFoundWithResponse, response)
	}

	return object, nil
}

func (s *CmsService) DescribeMetricRuleTargets(id string) (objects []interface{}, err error) {
	var response map[string]interface{}
	action := "DescribeMetricRuleTargets"
	client := s.client

	request := map[string]interface{}{
		"RuleId": id,
	}

	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(6*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("Cms", "2019-01-01", action, nil, request, false)
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
		return objects, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}

	resp, err := jsonpath.Get("$.Targets.Target", response)
	if err != nil {
		return objects, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.Targets.Target", response)
	}

	if v, ok := resp.([]interface{}); !ok || len(v) < 1 {
		return objects, errs.WrapErrorf(errs.NotFoundErr("Cms:GroupMetricRule", id), errs.NotFoundWithResponse, response)
	}

	objects = resp.([]interface{})

	return objects, nil
}

func (s *CmsService) DescribeCmsMetricRuleBlackList(id string) (object map[string]interface{}, err error) {
	client := s.client

	request := map[string]interface{}{
		"Ids.1": id,
	}

	var response map[string]interface{}
	action := "DescribeMetricRuleBlackList"
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		resp, err := client.RpcPost("Cms", "2019-01-01", action, nil, request, false)
		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		response = resp
		logs.AddDebug(action, response, request)
		return nil
	})
	if err != nil {
		return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}
	v, err := jsonpath.Get("$.DescribeMetricRuleBlackList[0]", response)
	totalCount, _ := jsonpath.Get("$.Total", response)
	total, _ := totalCount.(json.Number).Int64()
	if err != nil && total == 0 {
		err = errs.WrapErrorf(errs.NotFoundErr("Cms", id), errs.NotFoundMsg, errs.ProviderERROR)
		return object, err
	}
	return v.(map[string]interface{}), nil
}

func (s *CmsService) CmsMetricRuleBlackListStateRefreshFunc(d *schema.ResourceData, failStates []string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeCmsMetricRuleBlackList(d.Id())
		if err != nil {
			if errs.NotFoundError(err) {
				return nil, "", nil
			}
			return nil, "", errs.WrapError(err)
		}
		for _, failState := range failStates {
			if fmt.Sprint(object[""]) == failState {
				return object, fmt.Sprint(object[""]), errs.WrapError(errs.Error(errs.FailedToReachTargetStatus, fmt.Sprint(object[""])))
			}
		}
		return object, fmt.Sprint(object[""]), nil
	}
}

func (s *CmsService) CmsDynamicTagGroupStateRefreshFunc(id string, failStates []string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeCmsDynamicTagGroup(id)
		if err != nil {
			if errs.NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", errs.WrapError(err)
		}

		for _, failState := range failStates {
			if fmt.Sprint(object["Status"]) == failState {
				return object, fmt.Sprint(object["Status"]), errs.WrapError(errs.Error(errs.FailedToReachTargetStatus, fmt.Sprint(object["Status"])))
			}
		}

		return object, fmt.Sprint(object["Status"]), nil
	}
}

// NewCmsService constructs a CmsService bound to the provider client.
func NewCmsService(client *connectivity.AliyunClient) *CmsService {
	return &CmsService{client: client}
}
