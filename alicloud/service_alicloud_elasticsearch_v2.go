package alicloud

import (
	"fmt"
	"time"

	"github.com/PaesslerAG/jsonpath"
	roa "github.com/alibabacloud-go/tea-roa/client"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

type ElasticsearchServiceV2 struct {
	client *connectivity.AliyunClient
}

// DescribeElasticsearchLogstash <<< Encapsulated get interface for Elasticsearch Logstash.

func (s *ElasticsearchServiceV2) DescribeElasticsearchLogstash(id string) (object map[string]interface{}, err error) {

	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]*string
	var body map[string]interface{}
	InstanceId := id
	action := fmt.Sprintf("/openapi/logstashes/%s", InstanceId)
	conn, err := client.NewElasticsearchClient()
	if err != nil {
		return object, WrapError(err)
	}
	request = make(map[string]interface{})
	body = make(map[string]interface{})
	query = make(map[string]*string)
	request["InstanceId"] = id

	body["body"] = request
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(1*time.Minute, func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer("2017-06-13"), nil, StringPointer("GET"), StringPointer("AK"), StringPointer(action), query, nil, body, &util.RuntimeOptions{})

		if err != nil {
			if IsExpectedErrors(err, []string{"ServiceUnavailable"}) || NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, response, request)
		return nil
	})
	if err != nil {
		if IsExpectedErrors(err, []string{"InstanceNotFound"}) {
			return object, WrapErrorf(Error(GetNotFoundMessage("Logstash", id)), NotFoundMsg, ProviderERROR, fmt.Sprint(response["RequestId"]))
		}
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}

	v, err := jsonpath.Get("$.Result", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.Result", response)
	}

	return v.(map[string]interface{}), nil
}

func (s *ElasticsearchServiceV2) ElasticsearchLogstashStateRefreshFunc(id string, field string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeElasticsearchLogstash(id)
		if err != nil {
			if NotFoundError(err) {
				return object, "", nil
			}
			return nil, "", WrapError(err)
		}

		currentStatus := fmt.Sprint(object[field])
		for _, failState := range failStates {
			if currentStatus == failState {
				return object, currentStatus, WrapError(Error(FailedToReachTargetStatus, currentStatus))
			}
		}
		return object, currentStatus, nil
	}
}

// DescribeElasticsearchLogstash >>> Encapsulated.
// SetResourceTags <<< Encapsulated tag function for Elasticsearch.
func (s *ElasticsearchServiceV2) SetResourceTags(d *schema.ResourceData, resourceType string) error {
	if d.HasChange("tags") {
		var err error
		var action string
		var conn *roa.Client
		client := s.client
		var request map[string]interface{}
		var response map[string]interface{}
		query := make(map[string]*string)
		body := make(map[string]interface{})

		added, removed := parsingTags(d)
		removedTagKeys := make([]string, 0)
		for _, v := range removed {
			if !ignoredTags(v, "") {
				removedTagKeys = append(removedTagKeys, v)
			}
		}
		if len(removedTagKeys) > 0 {
			action = fmt.Sprintf("/openapi/tags")
			conn, err = client.NewElasticsearchClient()
			if err != nil {
				return WrapError(err)
			}
			request = make(map[string]interface{})
			query = make(map[string]*string)
			body = make(map[string]interface{})
			if v, ok := d.GetOk(""); ok {
				query[""] = StringPointer(v.(string))
			}
			if v, ok := d.GetOk(""); ok {
				query[""] = StringPointer(v.(string))
			}
			body["body"] = request
			wait := incrementalWait(3*time.Second, 5*time.Second)
			err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
				response, err = conn.DoRequest(StringPointer("2017-06-13"), nil, StringPointer("DELETE"), StringPointer("AK"), StringPointer(action), query, nil, body, &util.RuntimeOptions{})

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
			action = fmt.Sprintf("/openapi/tags")
			conn, err = client.NewElasticsearchClient()
			if err != nil {
				return WrapError(err)
			}
			request = make(map[string]interface{})
			query = make(map[string]*string)
			body = make(map[string]interface{})
			count := 1
			for key, value := range added {
				request[fmt.Sprintf("Tag.%d.Key", count)] = key
				request[fmt.Sprintf("Tag.%d.Value", count)] = value
				count++
			}

			request["ResourceType"] = resourceType
			body["body"] = request
			wait := incrementalWait(3*time.Second, 5*time.Second)
			err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
				response, err = conn.DoRequest(StringPointer("2017-06-13"), nil, StringPointer("POST"), StringPointer("AK"), StringPointer(action), query, nil, body, &util.RuntimeOptions{})

				if err != nil {
					if IsExpectedErrors(err, []string{"ServiceUnavailable"}) || NeedRetry(err) {
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

// SetResourceTags >>> tag function encapsulated.
