package alicloud

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"strconv"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/cms"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

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
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(6*time.Minute, func() *resource.RetryError {
		response, err = client.RpcPost("Cms", "2019-01-01", action, nil, request, false)
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
		if IsExpectedErrors(err, []string{"InternalError", "ResourceNotFound"}) {
			return object, WrapErrorf(Error(GetNotFoundMessage("Cms:Alarm", id)), NotFoundMsg, ProviderERROR, fmt.Sprint(response["RequestId"]))
		}
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}

	resp, err := jsonpath.Get("$.Alarms.Alarm", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.Alarms.Alarm", response)
	}

	if v, ok := resp.([]interface{}); !ok || len(v) < 1 {
		return object, WrapErrorf(Error(GetNotFoundMessage("Cms:Alarm", id)), NotFoundWithResponse, response)
	}

	for _, v := range resp.([]interface{}) {
		if fmt.Sprint(v.(map[string]interface{})["RuleId"]) == id {
			idExist = true
			return v.(map[string]interface{}), nil
		}
	}

	if !idExist {
		return object, WrapErrorf(Error(GetNotFoundMessage("Cms:Alarm", id)), NotFoundWithResponse, response)
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
		timeout = timeout - DefaultIntervalShort
		if timeout <= 0 {
			return GetTimeErrorFromString(GetTimeoutMessage("Alarm", strconv.FormatBool(enabled)))
		}
		time.Sleep(DefaultIntervalShort * time.Second)
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
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		raw, err = s.client.WithCmsClient(func(cmsClient *cms.Client) (interface{}, error) {
			return cmsClient.DescribeSiteMonitorList(listRequest)
		})
		if err != nil {
			if IsExpectedErrors(err, []string{ThrottlingUser, "ExceedingQuota"}) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	if err != nil {
		return siteMonitor, WrapError(err)
	}
	list := raw.(*cms.DescribeSiteMonitorListResponse)
	if len(list.SiteMonitors.SiteMonitor) < 1 {
		return siteMonitor, GetNotFoundErrorFromString(GetNotFoundMessage("Site Monitor", id))

	}
	for _, v := range list.SiteMonitors.SiteMonitor {
		if v.TaskName == keyword || v.TaskId == id {
			return v, nil
		}
	}
	return siteMonitor, GetNotFoundErrorFromString(GetNotFoundMessage("Site Monitor", id))
}

func (s *CmsService) GetIspCities(id string) (ispCities IspCities, err error) {
	request := cms.CreateDescribeSiteMonitorAttributeRequest()
	request.TaskId = id

	var raw interface{}
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		raw, err = s.client.WithCmsClient(func(cmsClient *cms.Client) (interface{}, error) {
			return cmsClient.DescribeSiteMonitorAttribute(request)
		})
		if err != nil {
			if IsExpectedErrors(err, []string{ThrottlingUser, "ExceedingQuota"}) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	if err != nil {
		return nil, WrapError(err)
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
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		raw, err := s.client.WithCmsClient(func(cmsClient *cms.Client) (interface{}, error) {
			return cmsClient.DescribeContactList(request)
		})
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		response, _ = raw.(*cms.DescribeContactListResponse)
		return nil
	})
	if err != nil {
		if IsExpectedErrors(err, []string{"ContactNotExists", "ResourceNotFound"}) {
			err = WrapErrorf(Error(GetNotFoundMessage("CmsAlarmContact", id)), NotFoundMsg, ProviderERROR)
			return
		}
		err = WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR)
		return
	}

	if response.Code != "200" {
		err = Error("DescribeContactList failed for " + response.Message)
		return
	}

	if len(response.Contacts.Contact) < 1 {
		err = WrapErrorf(Error(GetNotFoundMessage("CmsAlarmContact", id)), NotFoundMsg, ProviderERROR, response.RequestId)
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
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			raw, err = s.client.WithCmsClient(func(cmsClient *cms.Client) (interface{}, error) {
				return cmsClient.DescribeContactGroupList(request)
			})
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})

		if err != nil {
			if IsExpectedErrors(err, []string{"ContactGroupNotExists", "ResourceNotFound"}) {
				err = WrapErrorf(Error(GetNotFoundMessage("CmsAlarmContactGroup", id)), NotFoundMsg, ProviderERROR)
				return object, err
			}
			err = WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR)
			return object, err
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		response, _ := raw.(*cms.DescribeContactGroupListResponse)
		if response.Code != "200" {
			err = Error("DescribeContactGroupList failed for " + response.Message)
			return object, err
		}

		if len(response.ContactGroupList.ContactGroup) < 1 {
			err = WrapErrorf(Error(GetNotFoundMessage("CmsAlarmContactGroup", id)), NotFoundMsg, ProviderERROR, response.RequestId)
			return object, err
		}
		for _, object := range response.ContactGroupList.ContactGroup {
			if object.Name == id {
				return object, nil
			}
		}
		if len(response.ContactGroupList.ContactGroup) < PageSizeMedium {
			break
		}
		if page, err := getNextpageNumber(request.PageNumber); err != nil {
			return object, WrapError(err)
		} else {
			request.PageNumber = page
		}
	}
	err = WrapErrorf(Error(GetNotFoundMessage("CmsAlarmContactGroup", id)), NotFoundMsg, ProviderERROR)
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
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(6*time.Minute, func() *resource.RetryError {
		response, err = client.RpcPost("Cms", "2019-01-01", action, nil, request, false)
		if err != nil {
			if IsExpectedErrors(err, []string{"ExceedingQuota"}) || NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)

	if err != nil {
		if IsExpectedErrors(err, []string{"GroupMetricRuleNotExists", "ResourceNotFound", "ResourceNotFoundError"}) {
			return object, WrapErrorf(Error(GetNotFoundMessage("Cms:GroupMetricRule", id)), NotFoundMsg, ProviderERROR, fmt.Sprint(response["RequestId"]))
		}
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}

	resp, err := jsonpath.Get("$.Alarms.Alarm", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.Alarms.Alarm", response)
	}

	if v, ok := resp.([]interface{}); !ok || len(v) < 1 {
		return object, WrapErrorf(Error(GetNotFoundMessage("Cms:GroupMetricRule", id)), NotFoundWithResponse, response)
	}

	for _, v := range resp.([]interface{}) {
		if fmt.Sprint(v.(map[string]interface{})["RuleId"]) == id {
			idExist = true
			return v.(map[string]interface{}), nil
		}
	}

	if !idExist {
		return object, WrapErrorf(Error(GetNotFoundMessage("Cms:GroupMetricRule", id)), NotFoundWithResponse, response)
	}

	return object, nil
}

func (s *CmsService) SetResourceTags(d *schema.ResourceData, resourceType string) error {
	client := s.client
	if d.HasChange("tags") {
		added, removed := parsingTags(d)
		removedTagKeys := make([]string, 0)
		for _, v := range removed {
			if !ignoredTags(v, "") {
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
			wait := incrementalWait(2*time.Second, 1*time.Second)
			err := resource.Retry(10*time.Minute, func() *resource.RetryError {
				response, err := client.RpcPost("Cms", "2019-01-01", action, nil, request, false)
				if err != nil {
					if NeedRetry(err) {
						wait()
						return resource.RetryableError(err)

					}
					return resource.NonRetryableError(err)
				}
				addDebug(action, response, request)
				return nil
			})
			if err != nil {
				return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
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

			wait := incrementalWait(2*time.Second, 1*time.Second)
			err := resource.Retry(10*time.Minute, func() *resource.RetryError {
				response, err := client.RpcPost("Cms", "2019-01-01", action, nil, request, false)
				if err != nil {
					if NeedRetry(err) {
						wait()
						return resource.RetryableError(err)

					}
					return resource.NonRetryableError(err)
				}
				addDebug(action, response, request)
				return nil
			})
			if err != nil {
				return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
			}
		}
		d.SetPartial("tags")
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
		err = WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
		return
	}
	addDebug(action, response, request)
	if IsExpectedErrorCodes(fmt.Sprintf("%v", response["Code"]), []string{"GroupsNotExists", "ResourceNotFound"}) {
		err = WrapErrorf(Error(GetNotFoundMessage("CmsMonitorGroup", id)), NotFoundMsg, ProviderERROR)
		return object, err
	}
	if fmt.Sprintf(`%v`, response["Code"]) != "200" {
		err = Error("DescribeMonitorGroups failed for " + response["Message"].(string))
		return object, err
	}
	v, err := jsonpath.Get("$.Resources.Resource", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.Resources.Resource", response)
	}
	if len(v.([]interface{})) < 1 {
		return object, WrapErrorf(Error(GetNotFoundMessage("CloudMonitorService", id)), NotFoundWithResponse, response)
	} else {
		if fmt.Sprint(formatInt(v.([]interface{})[0].(map[string]interface{})["GroupId"])) != id {
			return object, WrapErrorf(Error(GetNotFoundMessage("CloudMonitorService", id)), NotFoundWithResponse, response)
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
	request["PageSize"] = PageSizeMedium
	request["PageNumber"] = 1
	for {
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			response, err = client.RpcPost("Cms", "2019-01-01", action, nil, request, false)
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
			if IsExpectedErrors(err, []string{"ResourceNotFound", "ResourceNotFoundError"}) {
				err = WrapErrorf(Error(GetNotFoundMessage("CmsMonitorGroupInstances", id)), NotFoundMsg, ProviderERROR)
				return object, err
			}
			return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
		}
		v, err := jsonpath.Get("$.Resources.Resource", response)
		if err != nil {
			return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.Resources.Resource", response)
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
		return object, WrapErrorf(Error(GetNotFoundMessage("CmsMonitorGroupInstances", id)), NotFoundWithResponse, response)
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

	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = client.RpcPost("Cms", "2019-01-01", action, nil, request, false)
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
		if IsExpectedErrors(err, []string{"ResourceNotFound"}) {
			return object, WrapErrorf(Error(GetNotFoundMessage("Cms:MetricRuleTemplate", id)), NotFoundMsg, ProviderERROR, fmt.Sprint(response["RequestId"]))
		}
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}

	v, err := jsonpath.Get("$.Resource", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.Resource", response)
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
		"PageSize":         PageSizeLarge,
		"PageNumber":       1,
	}

	idExist := false
	for {
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			response, err = client.RpcPost("Cms", "2019-01-01", action, nil, request, false)
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
			return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
		}

		resp, err := jsonpath.Get("$.TagGroupList.TagGroup", response)
		if err != nil {
			return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.TagGroupList.TagGroup", response)
		}

		if v, ok := resp.([]interface{}); !ok || len(v) < 1 {
			return object, WrapErrorf(Error(GetNotFoundMessage("Cms:DynamicTagGroup", id)), NotFoundWithResponse, response)
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
		return object, WrapErrorf(Error(GetNotFoundMessage("Cms:DynamicTagGroup", id)), NotFoundWithResponse, response)
	}

	return object, nil
}

func (s *CmsService) DescribeCmsNamespace(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	action := "DescribeHybridMonitorNamespaceList"
	client := s.client

	request := map[string]interface{}{
		"Namespace":  id,
		"PageSize":   PageSizeLarge,
		"PageNumber": 1,
	}

	idExist := false
	for {
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			response, err = client.RpcPost("Cms", "2019-01-01", action, nil, request, false)
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
			return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
		}

		resp, err := jsonpath.Get("$.DescribeHybridMonitorNamespace", response)
		if err != nil {
			return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.DescribeHybridMonitorNamespace", response)
		}

		if v, ok := resp.([]interface{}); !ok || len(v) < 1 {
			return object, WrapErrorf(Error(GetNotFoundMessage("CloudMonitorService:Namespace", id)), NotFoundWithResponse, response)
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
		return object, WrapErrorf(Error(GetNotFoundMessage("CloudMonitorService:Namespace", id)), NotFoundWithResponse, response)
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
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = client.RpcPost("Cms", "2019-01-01", action, nil, request, false)
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
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}
	if IsExpectedErrorCodes(fmt.Sprint(response["Code"]), []string{"400"}) {
		return object, WrapErrorf(Error(GetNotFoundMessage("CloudMonitorService:SlsGroup", id)), NotFoundMsg, ProviderERROR)
	}
	v, err := jsonpath.Get("$.List", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.List", response)
	}
	if len(v.([]interface{})) < 1 {
		return object, WrapErrorf(Error(GetNotFoundMessage("CloudMonitorService", id)), NotFoundWithResponse, response)
	} else {
		if fmt.Sprint(v.([]interface{})[0].(map[string]interface{})["SLSGroupName"]) != id {
			return object, WrapErrorf(Error(GetNotFoundMessage("CloudMonitorService", id)), NotFoundWithResponse, response)
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
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = client.RpcPost("Cms", "2019-01-01", action, nil, request, false)
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
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}
	if IsExpectedErrorCodes(fmt.Sprint(response["Code"]), []string{"ResourceNotFound"}) {
		return object, WrapErrorf(Error(GetNotFoundMessage("CloudMonitorService:HybridMonitorSlsTask", id)), NotFoundMsg, ProviderERROR)
	}
	v, err := jsonpath.Get("$.TaskList", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.TaskList", response)
	}
	if len(v.([]interface{})) < 1 {
		return object, WrapErrorf(Error(GetNotFoundMessage("CloudMonitorService", id)), NotFoundWithResponse, response)
	} else {
		if fmt.Sprint(v.([]interface{})[0].(map[string]interface{})["TaskId"]) != id {
			return object, WrapErrorf(Error(GetNotFoundMessage("CloudMonitorService", id)), NotFoundWithResponse, response)
		}
	}
	object = v.([]interface{})[0].(map[string]interface{})
	return object, nil
}

func (s *CmsService) DescribeCmsHybridMonitorFcTask(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	client := s.client
	action := "DescribeHybridMonitorTaskList"
	parts, err := ParseResourceId(id, 2)
	if err != nil {
		err = WrapError(err)
		return
	}
	request := map[string]interface{}{
		"TargetUserId": parts[0],
		"Namespace":    parts[1],
		"TaskType":     "aliyun_fc",
	}
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = client.RpcPost("Cms", "2019-01-01", action, nil, request, false)
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
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}
	if IsExpectedErrorCodes(fmt.Sprint(response["Code"]), []string{"ResourceNotFound"}) {
		return object, WrapErrorf(Error(GetNotFoundMessage("CloudMonitorService:HybridMonitorFcTask", id)), NotFoundMsg, ProviderERROR)
	}
	v, err := jsonpath.Get("$.TaskList", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.TaskList", response)
	}
	if len(v.([]interface{})) < 1 {
		return object, WrapErrorf(Error(GetNotFoundMessage("CloudMonitorService", id)), NotFoundWithResponse, response)
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
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = client.RpcPost("Cms", "2019-01-01", action, nil, request, false)
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
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}
	resp, err := jsonpath.Get("$.EventRules.EventRule", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.EventRules.EventRule", response)
	}

	if v, ok := resp.([]interface{}); !ok || len(v) < 1 {
		return object, WrapErrorf(Error(GetNotFoundMessage("Cms:EventRule", id)), NotFoundWithResponse, response)
	}

	for _, v := range resp.([]interface{}) {
		if fmt.Sprint(v.(map[string]interface{})["Name"]) == id {
			idExist = true
			return v.(map[string]interface{}), nil
		}
	}

	if !idExist {
		return object, WrapErrorf(Error(GetNotFoundMessage("Cms:EventRule", id)), NotFoundWithResponse, response)
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

	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(6*time.Minute, func() *resource.RetryError {
		response, err = client.RpcPost("Cms", "2019-01-01", action, nil, request, false)
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
		return objects, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}

	resp, err := jsonpath.Get("$.Targets.Target", response)
	if err != nil {
		return objects, WrapErrorf(err, FailedGetAttributeMsg, id, "$.Targets.Target", response)
	}

	if v, ok := resp.([]interface{}); !ok || len(v) < 1 {
		return objects, WrapErrorf(Error(GetNotFoundMessage("Cms:GroupMetricRule", id)), NotFoundWithResponse, response)
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
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		resp, err := client.RpcPost("Cms", "2019-01-01", action, nil, request, false)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		response = resp
		addDebug(action, response, request)
		return nil
	})
	if err != nil {
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}
	v, err := jsonpath.Get("$.DescribeMetricRuleBlackList[0]", response)
	totalCount, _ := jsonpath.Get("$.Total", response)
	total, _ := totalCount.(json.Number).Int64()
	if err != nil && total == 0 {
		err = WrapErrorf(Error(GetNotFoundMessage("Cms", id)), NotFoundMsg, ProviderERROR)
		return object, err
	}
	return v.(map[string]interface{}), nil
}

func (s *CmsService) CmsMetricRuleBlackListStateRefreshFunc(d *schema.ResourceData, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeCmsMetricRuleBlackList(d.Id())
		if err != nil {
			if NotFoundError(err) {
				return nil, "", nil
			}
			return nil, "", WrapError(err)
		}
		for _, failState := range failStates {
			if fmt.Sprint(object[""]) == failState {
				return object, fmt.Sprint(object[""]), WrapError(Error(FailedToReachTargetStatus, fmt.Sprint(object[""])))
			}
		}
		return object, fmt.Sprint(object[""]), nil
	}
}

func (s *CmsService) CmsDynamicTagGroupStateRefreshFunc(id string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeCmsDynamicTagGroup(id)
		if err != nil {
			if NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", WrapError(err)
		}

		for _, failState := range failStates {
			if fmt.Sprint(object["Status"]) == failState {
				return object, fmt.Sprint(object["Status"]), WrapError(Error(FailedToReachTargetStatus, fmt.Sprint(object["Status"])))
			}
		}

		return object, fmt.Sprint(object["Status"]), nil
	}
}
