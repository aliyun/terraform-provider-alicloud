package alicloud

import (
	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
)

type Brain_industrialService struct {
	client *connectivity.AliyunClient
}

func (s *Brain_industrialService) DescribeBrainIndustrialPidOrganization(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	client := s.client
	action := "ListPidOrganizations"
	request := map[string]interface{}{
		"RegionId": s.client.RegionId,
	}
	idExist := false
	response, err = client.RpcPost("brain-industrial", "2020-09-20", action, nil, request, true)
	if err != nil {
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}
	addDebug(action, response, request)
	v, err := jsonpath.Get("$.OrganizationList", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.OrganizationList", response)
	}
	if len(v.([]interface{})) < 1 {
		return object, WrapErrorf(NotFoundErr("BrainIndustrial", id), NotFoundWithResponse, response)
	}
	for _, v := range v.([]interface{}) {
		if v.(map[string]interface{})["OrganizationId"].(string) == id {
			idExist = true
			return v.(map[string]interface{}), nil
		}
	}
	if !idExist {
		return object, WrapErrorf(NotFoundErr("BrainIndustrial", id), NotFoundWithResponse, response)
	}
	return object, nil
}

func (s *Brain_industrialService) DescribeBrainIndustrialPidProject(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	client := s.client
	action := "ListPidProjects"
	request := map[string]interface{}{
		"RegionId":    s.client.RegionId,
		"CurrentPage": 1,
		"PageSize":    20,
	}
	idExist := false
	for {
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		response, err = client.RpcPost("brain-industrial", "2020-09-20", action, nil, request, true)
		if err != nil {
			return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
		}
		addDebug(action, response, request)
		v, err := jsonpath.Get("$.PidProjectList", response)
		if err != nil {
			return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.PidProjectList", response)
		}
		if len(v.([]interface{})) < 1 {
			return object, WrapErrorf(NotFoundErr("BrainIndustrial", id), NotFoundWithResponse, response)
		}
		for _, v := range v.([]interface{}) {
			if v.(map[string]interface{})["PidProjectId"].(string) == id {
				idExist = true
				return v.(map[string]interface{}), nil
			}
		}
		if len(v.([]interface{})) < request["PageSize"].(int) {
			break
		}
		request["CurrentPage"] = request["CurrentPage"].(int) + 1
	}
	if !idExist {
		return object, WrapErrorf(NotFoundErr("BrainIndustrial", id), NotFoundWithResponse, response)
	}
	return
}

func (s *Brain_industrialService) DescribeBrainIndustrialPidLoop(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	client := s.client
	action := "GetLoop"
	request := map[string]interface{}{
		"RegionId":  s.client.RegionId,
		"PidLoopId": id,
	}
	response, err = client.RpcPost("brain-industrial", "2020-09-20", action, nil, request, true)
	if err != nil {
		if IsExpectedErrors(err, []string{"-106"}) {
			return object, WrapErrorf(err, NotFoundMsg, AlibabaCloudSdkGoERROR)
		}
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}
	addDebug(action, response, request)
	v, err := jsonpath.Get("$.Data", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.Data", response)
	}
	object = v.(map[string]interface{})
	return object, nil
}
