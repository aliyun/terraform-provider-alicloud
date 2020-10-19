package alicloud

import (
	"encoding/json"
	"fmt"
	"time"

	"strconv"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/cms"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
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

func (s *CmsService) DescribeAlarm(id string) (alarm cms.Alarm, err error) {
	request := cms.CreateDescribeMetricRuleListRequest()
	request.RuleIds = id

	wait := incrementalWait(3*time.Second, 5*time.Second)
	var response *cms.DescribeMetricRuleListResponse
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		raw, err := s.client.WithCmsClient(func(cmsClient *cms.Client) (interface{}, error) {
			return cmsClient.DescribeMetricRuleList(request)
		})
		if err != nil {
			if IsExpectedErrors(err, []string{ThrottlingUser}) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		response, _ = raw.(*cms.DescribeMetricRuleListResponse)
		return nil
	})
	if err != nil {
		return alarm, err
	}

	if len(response.Alarms.Alarm) < 1 {
		return alarm, GetNotFoundErrorFromString(GetNotFoundMessage("Alarm Rule", id))
	}

	return response.Alarms.Alarm[0], nil
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

		if alarm.EnableState == enabled {
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
	raw, err := s.client.WithCmsClient(func(cmsClient *cms.Client) (interface{}, error) {
		return cmsClient.DescribeSiteMonitorList(listRequest)
	})
	if err != nil {
		return siteMonitor, err
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

	raw, err := s.client.WithCmsClient(func(cmsClient *cms.Client) (interface{}, error) {
		return cmsClient.DescribeSiteMonitorAttribute(request)
	})

	if err != nil {
		return nil, err
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

	raw, err := s.client.WithCmsClient(func(cmsClient *cms.Client) (interface{}, error) {
		return cmsClient.DescribeContactList(request)
	})
	if err != nil {
		if IsExpectedErrors(err, []string{"ContactNotExists", "ResourceNotFound"}) {
			err = WrapErrorf(Error(GetNotFoundMessage("CmsAlarmContact", id)), NotFoundMsg, ProviderERROR)
			return
		}
		err = WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR)
		return
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response, _ := raw.(*cms.DescribeContactListResponse)
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

		raw, err := s.client.WithCmsClient(func(cmsClient *cms.Client) (interface{}, error) {
			return cmsClient.DescribeContactGroupList(request)
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
