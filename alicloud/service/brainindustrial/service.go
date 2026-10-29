package brainindustrial

import (
	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
)

// Service wraps the Brain Industrial API calls used by the resources and data
// sources of this package. The root alicloud package cannot be imported from a
// service package, so the encapsulated calls live here with the resources.
type Service struct {
	client *connectivity.AliyunClient
}

// NewService is the exported constructor: the client field stays unexported,
// so callers outside the package (including the acceptance tests) build values
// with it.
func NewService(client *connectivity.AliyunClient) *Service {
	return &Service{client: client}
}

func (s *Service) DescribePidOrganization(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	client := s.client
	action := "ListPidOrganizations"
	request := map[string]interface{}{
		"RegionId": s.client.RegionId,
	}
	idExist := false
	response, err = client.RpcPost("brain-industrial", "2020-09-20", action, nil, request, true)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}
	logs.AddDebug(action, response, request)
	v, err := jsonpath.Get("$.OrganizationList", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.OrganizationList", response)
	}
	if len(v.([]interface{})) < 1 {
		return object, errs.WrapErrorf(errs.NotFoundErr("BrainIndustrial", id), errs.NotFoundWithResponse, response)
	}
	for _, v := range v.([]interface{}) {
		if v.(map[string]interface{})["OrganizationId"].(string) == id {
			idExist = true
			return v.(map[string]interface{}), nil
		}
	}
	if !idExist {
		return object, errs.WrapErrorf(errs.NotFoundErr("BrainIndustrial", id), errs.NotFoundWithResponse, response)
	}
	return object, nil
}

func (s *Service) DescribePidProject(id string) (object map[string]interface{}, err error) {
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
			return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
		}
		logs.AddDebug(action, response, request)
		v, err := jsonpath.Get("$.PidProjectList", response)
		if err != nil {
			return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.PidProjectList", response)
		}
		if len(v.([]interface{})) < 1 {
			return object, errs.WrapErrorf(errs.NotFoundErr("BrainIndustrial", id), errs.NotFoundWithResponse, response)
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
		return object, errs.WrapErrorf(errs.NotFoundErr("BrainIndustrial", id), errs.NotFoundWithResponse, response)
	}
	return
}

func (s *Service) DescribePidLoop(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	client := s.client
	action := "GetLoop"
	request := map[string]interface{}{
		"RegionId":  s.client.RegionId,
		"PidLoopId": id,
	}
	response, err = client.RpcPost("brain-industrial", "2020-09-20", action, nil, request, true)
	if err != nil {
		if errs.IsExpectedErrors(err, []string{"-106"}) {
			return object, errs.WrapErrorf(err, errs.NotFoundMsg, errs.AlibabaCloudSdkGoERROR)
		}
		return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}
	logs.AddDebug(action, response, request)
	v, err := jsonpath.Get("$.Data", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.Data", response)
	}
	object = v.(map[string]interface{})
	return object, nil
}
