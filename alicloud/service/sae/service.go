package sae

import (
	"strconv"

	"github.com/alibabacloud-go/tea/tea"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
)

const (
	saeAPIVersion       = "2019-05-06"
	saePageSize         = 100
	saeMaxPageLoopCount = 10000
)

type Service struct {
	client *connectivity.AliyunClient
}

func NewService(client *connectivity.AliyunClient) *Service {
	return &Service{client: client}
}

// DescribeApplicationGroups returns every application group of the given
// application, paging until a short page arrives (the API exposes no
// TotalSize to terminate on).
func (s *Service) DescribeApplicationGroups(appId string) ([]map[string]interface{}, error) {
	const action = "DescribeApplicationGroups"
	groups := make([]map[string]interface{}, 0)
	for page := 1; page <= saeMaxPageLoopCount; page++ {
		request := map[string]*string{
			"AppId":       tea.String(appId),
			"CurrentPage": tea.String(strconv.Itoa(page)),
			"PageSize":    tea.String(strconv.Itoa(saePageSize)),
		}
		response, err := s.client.RoaGet("sae", saeAPIVersion, "/pop/v1/sam/app/describeApplicationGroups", request, nil, nil)
		if err != nil {
			return nil, errs.WrapErrorf(err, errs.DataDefaultErrorMsg, "alicloud_sae_application_instances", action, errs.AlibabaCloudSdkGoERROR)
		}
		data, _ := response["Data"].([]interface{})
		for _, item := range data {
			if group, ok := item.(map[string]interface{}); ok {
				groups = append(groups, group)
			}
		}
		if len(data) < saePageSize {
			break
		}
	}
	return groups, nil
}

// DescribeApplicationInstances returns every instance of the given
// application group, paging until the reported TotalSize is reached or a
// short page arrives. A response carrying no Data yields an empty result
// instead of an error.
func (s *Service) DescribeApplicationInstances(appId, groupId string) ([]map[string]interface{}, error) {
	const action = "DescribeApplicationInstances"
	instances := make([]map[string]interface{}, 0)
	for page := 1; page <= saeMaxPageLoopCount; page++ {
		request := map[string]*string{
			"AppId":       tea.String(appId),
			"GroupId":     tea.String(groupId),
			"CurrentPage": tea.String(strconv.Itoa(page)),
			"PageSize":    tea.String(strconv.Itoa(saePageSize)),
		}
		response, err := s.client.RoaGet("sae", saeAPIVersion, "/pop/v1/sam/app/describeApplicationInstances", request, nil, nil)
		if err != nil {
			return nil, errs.WrapErrorf(err, errs.DataDefaultErrorMsg, "alicloud_sae_application_instances", action, errs.AlibabaCloudSdkGoERROR)
		}
		data, _ := response["Data"].(map[string]interface{})
		items, _ := data["Instances"].([]interface{})
		for _, item := range items {
			if instance, ok := item.(map[string]interface{}); ok {
				instances = append(instances, instance)
			}
		}
		if len(items) < saePageSize {
			break
		}
		if totalSize, ok := data["TotalSize"].(float64); ok && len(instances) >= int(totalSize) {
			break
		}
	}
	return instances, nil
}
