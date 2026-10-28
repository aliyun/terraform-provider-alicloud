// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package elasticsearch

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/convert"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/tags"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/tidwall/sjson"
)

type ElasticsearchServiceV2 struct {
	client *connectivity.AliyunClient
}

// NewElasticsearchServiceV2 constructs an ElasticsearchServiceV2 bound to the
// shared connectivity client; acceptance tests use it from the external
// elasticsearch_test package.
func NewElasticsearchServiceV2(client *connectivity.AliyunClient) *ElasticsearchServiceV2 {
	return &ElasticsearchServiceV2{client: client}
}

// DescribeElasticsearchInstance <<< Encapsulated get interface for Elasticsearch Instance.

func (s *ElasticsearchServiceV2) DescribeElasticsearchInstance(id string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]*string
	var header map[string]*string
	InstanceId := id
	request = make(map[string]interface{})
	query = make(map[string]*string)
	header = make(map[string]*string)

	action := fmt.Sprintf("/openapi/instances/%s", InstanceId)

	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(1*time.Minute, func() *retry.RetryError {
		response, err = client.RoaGet("elasticsearch", "2017-06-13", action, query, header, nil)

		if err != nil {
			if errs.IsExpectedErrors(err, []string{"ServiceUnavailable"}) || helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(action, response, request)
	if err != nil {
		if errs.IsExpectedErrors(err, []string{"InstanceNotFound"}) {
			return object, errs.WrapErrorf(errs.NotFoundErr("Instance", id), errs.NotFoundMsg, response)
		}
		return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}

	v, err := jsonpath.Get("$.Result", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.Result", response)
	}

	return v.(map[string]interface{}), nil
}
func (s *ElasticsearchServiceV2) DescribeInstanceQueryAvailableInstances(d *schema.ResourceData) (object map[string]interface{}, err error) {
	client := s.client
	id := d.Id()
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	query["InstanceIDs"] = id

	query["ProductCode"] = "elasticsearch"
	var endpoint string
	request["ProductCode"] = "elasticsearch"
	request["ProductType"] = "elasticsearchpre"
	if client.IsInternationalAccount() {
		request["ProductType"] = "elasticsearchpre_intl"
	}
	action := "QueryAvailableInstances"

	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(1*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPostWithEndpoint("BssOpenApi", "2017-12-14", action, query, request, true, endpoint)

		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			if !client.IsInternationalAccount() && errs.IsExpectedErrors(err, []string{"NotApplicable"}) {
				request["ProductCode"] = "elasticsearch"
				request["ProductType"] = "elasticsearchpre_intl"
				endpoint = connectivity.BssOpenAPIEndpointInternational
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

	v, err := jsonpath.Get("$.Data.InstanceList[*]", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.Data.InstanceList[*]", response)
	}

	if len(v.([]interface{})) == 0 {
		return object, errs.WrapErrorf(errs.NotFoundErr("Instance", id), errs.NotFoundMsg, response)
	}

	return v.([]interface{})[0].(map[string]interface{}), nil
}
func (s *ElasticsearchServiceV2) DescribeInstanceListKibanaPvlNetwork(id string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]*string
	var header map[string]*string
	InstanceId := id
	request = make(map[string]interface{})
	query = make(map[string]*string)
	header = make(map[string]*string)

	action := fmt.Sprintf("/openapi/instances/%s/actions/get-kibana-private", InstanceId)

	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(1*time.Minute, func() *retry.RetryError {
		response, err = client.RoaGet("elasticsearch", "2017-06-13", action, query, header, nil)

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
		if errs.IsExpectedErrors(err, []string{"InstanceNotFound"}) {
			return object, errs.WrapErrorf(errs.NotFoundErr("Instance", id), errs.NotFoundMsg, response)
		}
		return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}

	v, err := jsonpath.Get("$.Result[*]", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.Result[*]", response)
	}

	kibanaPvlNetwork, err := parseKibanaPvlResult(v)
	if err != nil {
		return object, errs.WrapErrorf(err, "get kibana pvl network error %s", id)
	}

	if kibanaPvlNetwork == nil {
		return object, errs.WrapErrorf(errs.NotFoundErr("Instance", id), errs.NotFoundMsg, response)
	}

	return kibanaPvlNetwork, nil
}

func (s *ElasticsearchServiceV2) ElasticsearchInstanceStateRefreshFunc(id string, field string, failStates []string) retry.StateRefreshFunc {
	return s.ElasticsearchInstanceStateRefreshFuncWithApi(id, field, failStates, s.DescribeElasticsearchInstance)
}

func (s *ElasticsearchServiceV2) ElasticsearchInstanceStateRefreshFuncWithApi(id string, field string, failStates []string, call func(id string) (map[string]interface{}, error)) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := call(id)
		if err != nil {
			if errs.NotFoundError(err) {
				return object, "", nil
			}
			return nil, "", errs.WrapError(err)
		}
		object["paymentType"] = convertElasticsearchInstanceResultpaymentTypeResponse(object["paymentType"])
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
				return object, currentStatus, errs.WrapError(errs.Error(errs.FailedToReachTargetStatus, currentStatus))
			}
		}
		return object, currentStatus, nil
	}
}

// DescribeElasticsearchInstance >>> Encapsulated.

// SetResourceTags <<< Encapsulated tag function for Elasticsearch.
func (s *ElasticsearchServiceV2) SetResourceTags(d *schema.ResourceData, resourceType string) error {
	if d.HasChange("tags") {
		var action string
		var err error
		client := s.client
		var request map[string]interface{}
		var response map[string]interface{}
		header := make(map[string]*string)
		query := make(map[string]*string)
		body := make(map[string]interface{})

		added, removed := tags.ParsingTags(d)
		removedTagKeys := make([]string, 0)
		for _, v := range removed {
			if !tags.IgnoredTags(v, "") {
				removedTagKeys = append(removedTagKeys, v)
			}
		}
		if len(removedTagKeys) > 0 {
			action = fmt.Sprintf("/openapi/tags")
			request = make(map[string]interface{})
			query = make(map[string]*string)
			body = make(map[string]interface{})
			query["ResourceIds"] = helper.StringPointer(convert.ConvertObjectToJsonString(convert.ExpandSingletonToList(d.Id())))

			query["TagKeys"] = helper.StringPointer(convert.ConvertListToJsonString(convert.ConvertListStringToListInterface(removedTagKeys)))
			query["ResourceType"] = helper.StringPointer(resourceType)
			body = request
			wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
			err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
				response, err = client.RoaDelete("elasticsearch", "2017-06-13", action, query, header, nil, true)
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
				return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
			}

		}

		if len(added) > 0 {
			action = fmt.Sprintf("/openapi/tags")
			request = make(map[string]interface{})
			query = make(map[string]*string)
			body = make(map[string]interface{})

			count := 1
			tagsMaps := make([]map[string]interface{}, 0)
			for key, value := range added {
				tagsMap := make(map[string]interface{})
				tagsMap["value"] = value
				tagsMap["key"] = key
				tagsMaps = append(tagsMaps, tagsMap)
				count++
			}
			request["Tags"] = tagsMaps

			request["ResourceType"] = resourceType
			jsonString := convert.ConvertObjectToJsonString(request)
			jsonString, _ = sjson.Set(jsonString, "ResourceIds.0", d.Id())
			_ = json.Unmarshal([]byte(jsonString), &request)

			body = request
			wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
			err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
				response, err = client.RoaPost("elasticsearch", "2017-06-13", action, query, header, body, true)
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
				return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
			}

		}
	}

	return nil
}

// SetResourceTags >>> tag function encapsulated.
