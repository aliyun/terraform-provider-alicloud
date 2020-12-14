package alicloud

import (
	"encoding/json"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

type DcdnService struct {
	client *connectivity.AliyunClient
}

func (s *DcdnService) convertSourcesToString(v []interface{}) (string, error) {
	arrayMaps := make([]interface{}, len(v))
	for i, vv := range v {
		item := vv.(map[string]interface{})
		arrayMaps[i] = map[string]interface{}{
			"Content":  item["content"],
			"Port":     item["port"],
			"Priority": item["priority"],
			"Type":     item["type"],
			"Weight":   item["weight"],
		}
	}
	maps, err := json.Marshal(arrayMaps)
	if err != nil {
		return "", WrapError(err)
	}
	return string(maps), nil
}

func (s *DcdnService) DescribeDcdnDomainCertificateInfo(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	conn, err := s.client.NewDcdnClient()
	if err != nil {
		return nil, WrapError(err)
	}
	action := "DescribeDcdnDomainCertificateInfo"
	request := map[string]interface{}{
		"RegionId":   s.client.RegionId,
		"DomainName": id,
	}
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2018-01-15"), StringPointer("AK"), nil, request, &runtime)
	if err != nil {
		err = WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
		return
	}
	addDebug(action, response, request)
	v, err := jsonpath.Get("$.CertInfos.CertInfo", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.CertInfos.CertInfo", response)
	}
	if len(v.([]interface{})) < 1 {
		return object, WrapErrorf(Error(GetNotFoundMessage("DCDN", id)), NotFoundWithResponse, response)
	} else {
		if v.([]interface{})[0].(map[string]interface{})["DomainName"].(string) != id {
			return object, WrapErrorf(Error(GetNotFoundMessage("DCDN", id)), NotFoundWithResponse, response)
		}
	}
	object = v.([]interface{})[0].(map[string]interface{})
	return object, nil
}

func (s *DcdnService) DescribeDcdnDomain(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	conn, err := s.client.NewDcdnClient()
	if err != nil {
		return nil, WrapError(err)
	}
	action := "DescribeDcdnDomainDetail"
	request := map[string]interface{}{
		"RegionId":   s.client.RegionId,
		"DomainName": id,
	}
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2018-01-15"), StringPointer("AK"), nil, request, &runtime)
	if err != nil {
		if IsExpectedErrors(err, []string{"InvalidDomain.NotFound"}) {
			err = WrapErrorf(Error(GetNotFoundMessage("DcdnDomain", id)), NotFoundMsg, ProviderERROR)
			return object, err
		}
		err = WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
		return object, err
	}
	addDebug(action, response, request)
	v, err := jsonpath.Get("$.DomainDetail", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.DomainDetail", response)
	}
	object = v.(map[string]interface{})
	return object, nil
}

func (s *DcdnService) DcdnDomainStateRefreshFunc(id string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeDcdnDomain(id)
		if err != nil {
			if NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", WrapError(err)
		}

		for _, failState := range failStates {
			if object["DomainStatus"].(string) == failState {
				return object, object["DomainStatus"].(string), WrapError(Error(FailedToReachTargetStatus, object["DomainStatus"].(string)))
			}
		}
		return object, object["DomainStatus"].(string), nil
	}
}
