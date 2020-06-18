package alicloud

import (
	"encoding/json"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ddoscoo"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

type DdoscooService struct {
	client *connectivity.AliyunClient
}

func (s *DdoscooService) DescribeDdoscooInstance(id string) (object ddoscoo.InstanceSpec, err error) {
	request := ddoscoo.CreateDescribeInstanceSpecsRequest()
	request.RegionId = s.client.RegionId

	request.InstanceIds = &[]string{id}

	raw, err := s.client.WithDdoscooClient(func(ddoscooClient *ddoscoo.Client) (interface{}, error) {
		return ddoscooClient.DescribeInstanceSpecs(request)
	})
	if err != nil {
		if IsExpectedErrors(err, []string{"InstanceNotFound", "ddos_coop3301"}) {
			err = WrapErrorf(Error(GetNotFoundMessage("DdoscooInstance", id)), NotFoundMsg, ProviderERROR)
			return
		}
		err = WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR)
		return
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response, _ := raw.(*ddoscoo.DescribeInstanceSpecsResponse)

	if len(response.InstanceSpecs) < 1 {
		err = WrapErrorf(Error(GetNotFoundMessage("DdoscooInstance", id)), NotFoundMsg, ProviderERROR)
		return
	}
	return response.InstanceSpecs[0], nil
}

func (s *DdoscooService) convertRulesToString(v []interface{}) (string, error) {
	arrayMaps := make([]ddoscoo.Rule, len(v))
	for i, vv := range v {
		item := vv.(map[string]interface{})
		arrayMaps[i] = ddoscoo.Rule{
			Priority:  item["priority"].(int),
			RegionId:  item["region_id"].(string),
			Status:    item["status"].(int),
			Type:      item["type"].(string),
			Value:     item["value"].(string),
			ValueType: item["value_type"].(int),
		}
	}
	maps, err := json.Marshal(arrayMaps)
	if err != nil {
		return "", WrapError(err)
	}
	return string(maps), nil
}

func (s *DdoscooService) DescribeDdoscooSchedulerRule(id string) (object ddoscoo.SchedulerRule, err error) {
	request := ddoscoo.CreateDescribeSchedulerRulesRequest()
	request.RegionId = s.client.RegionId

	request.RuleName = id
	request.PageSize = requests.NewInteger(10)

	raw, err := s.client.WithDdoscooClient(func(ddoscooClient *ddoscoo.Client) (interface{}, error) {
		return ddoscooClient.DescribeSchedulerRules(request)
	})
	if err != nil {
		err = WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR)
		return
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response, _ := raw.(*ddoscoo.DescribeSchedulerRulesResponse)

	if len(response.SchedulerRules) < 1 {
		err = WrapErrorf(Error(GetNotFoundMessage("DdoscooSchedulerRule", id)), NotFoundMsg, ProviderERROR)
		return
	}
	return response.SchedulerRules[0], nil
}

func (s *DdoscooService) DescribeInstances(id string) (object ddoscoo.Instance, err error) {
	request := ddoscoo.CreateDescribeInstancesRequest()
	request.RegionId = s.client.RegionId

	request.PageSize = "20"
	request.PageNumber = "1"
	for {

		raw, err := s.client.WithDdoscooClient(func(ddoscooClient *ddoscoo.Client) (interface{}, error) {
			return ddoscooClient.DescribeInstances(request)
		})
		if err != nil {
			if IsExpectedErrors(err, []string{"InstanceNotFound"}) {
				err = WrapErrorf(Error(GetNotFoundMessage("ddoscooinstance", id)), NotFoundMsg, ProviderERROR)
				return object, err
			}
			err = WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR)
			return object, err
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		response, _ := raw.(*ddoscoo.DescribeInstancesResponse)

		if len(response.Instances) < 1 {
			err = WrapErrorf(Error(GetNotFoundMessage("ddoscooinstance", id)), NotFoundMsg, ProviderERROR)
			return object, err
		}
		for _, object := range response.Instances {
			if object.InstanceId == id {
				return object, nil
			}
		}
		if len(response.Instances) < PageSizeMedium {
			break
		}
		if page, err := getNextpageNumber(requests.Integer(request.PageNumber)); err != nil {
			return object, WrapError(err)
		} else {
			request.PageNumber = string(page)
		}
	}
	err = WrapErrorf(Error(GetNotFoundMessage("ddoscooinstance", id)), NotFoundMsg, ProviderERROR)
	return
}

func (s *DdoscooService) DescribeTagResources(id string) (object ddoscoo.DescribeTagResourcesResponse, err error) {
	request := ddoscoo.CreateDescribeTagResourcesRequest()
	request.RegionId = s.client.RegionId

	request.ResourceIds = &[]string{id}
	request.ResourceType = "INSTANCE"

	raw, err := s.client.WithDdoscooClient(func(ddoscooClient *ddoscoo.Client) (interface{}, error) {
		return ddoscooClient.DescribeTagResources(request)
	})
	if err != nil {
		if IsExpectedErrors(err, []string{"InstanceNotFound"}) {
			err = WrapErrorf(Error(GetNotFoundMessage("ddoscooinstance", id)), NotFoundMsg, ProviderERROR)
			return
		}
		err = WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR)
		return
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response, _ := raw.(*ddoscoo.DescribeTagResourcesResponse)
	return *response, nil
}

func (s *DdoscooService) SetResourceTags(d *schema.ResourceData, resourceType string) error {
	oldItems, newItems := d.GetChange("tags")
	added := make([]ddoscoo.CreateTagResourcesTags, 0)
	for key, value := range newItems.(map[string]interface{}) {
		added = append(added, ddoscoo.CreateTagResourcesTags{
			Key:   key,
			Value: value.(string),
		})
	}
	removed := make([]string, 0)
	for key, _ := range oldItems.(map[string]interface{}) {
		removed = append(removed, key)
	}
	if len(removed) > 0 {
		request := ddoscoo.CreateDeleteTagResourcesRequest()
		request.RegionId = s.client.RegionId
		request.ResourceIds = &[]string{d.Id()}
		request.ResourceType = resourceType
		request.TagKey = &removed
		raw, err := s.client.WithDdoscooClient(func(ddoscooClient *ddoscoo.Client) (interface{}, error) {
			return ddoscooClient.DeleteTagResources(request)
		})
		addDebug(request.GetActionName(), raw)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
	}
	if len(added) > 0 {
		request := ddoscoo.CreateCreateTagResourcesRequest()
		request.RegionId = s.client.RegionId
		request.ResourceIds = &[]string{d.Id()}
		request.ResourceType = resourceType
		request.Tags = &added
		raw, err := s.client.WithDdoscooClient(func(ddoscooClient *ddoscoo.Client) (interface{}, error) {
			return ddoscooClient.CreateTagResources(request)
		})
		addDebug(request.GetActionName(), raw)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
	}
	return nil
}
