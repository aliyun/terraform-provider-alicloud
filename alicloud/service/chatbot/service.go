package chatbot

import (
	"fmt"
	"strconv"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
)

// Service wraps the Chatbot API calls used by the resources and data
// sources of this package.
type Service struct {
	client *connectivity.AliyunClient
}

// NewService is the exported constructor: the client field stays
// unexported, so callers outside the package (including the flat
// layout's acceptance tests) build values with it.
func NewService(client *connectivity.AliyunClient) *Service {
	return &Service{client: client}
}

func (s *Service) DescribeChatbotPublishTask(id string) (object map[string]interface{}, err error) {
	client := s.client
	request := map[string]interface{}{
		"Id":       id,
		"RegionId": s.client.RegionId,
	}

	var response map[string]interface{}
	action := "GetPublishTaskState"
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		resp, err := client.RpcPost("Chatbot", "2022-04-08", action, nil, request, true)
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
		if errs.IsExpectedErrors(err, []string{"20021207"}) {
			return object, errs.WrapErrorf(err, errs.NotFoundMsg, errs.AlibabaCloudSdkGoERROR)
		}
		return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}
	v, err := jsonpath.Get("$", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$", response)
	}
	return v.(map[string]interface{}), nil
}

// ListFaqCategories queries the FAQ top-level category ids published on the
// caller's default business space. CreatePublishTask only accepts DataIdList
// when BizType is "faq", and the values are FAQ category ids; the acceptance
// test uses two real ids to exercise the data_id_list reorder coverage instead
// of a fabricated placeholder. AgentKey is optional on the ListCategory API, so
// omitting it lists the categories of the default business space.
func (s *Service) ListFaqCategories() ([]string, error) {
	client := s.client
	action := "ListCategory"
	request := map[string]interface{}{
		"RegionId":      client.RegionId,
		"KnowledgeType": 1, // 1 = faq knowledge base
	}
	var response map[string]interface{}
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err := retry.Retry(5*time.Minute, func() *retry.RetryError {
		resp, err := client.RpcPost("Chatbot", "2022-04-08", action, nil, request, true)
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
		return nil, errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_chatbot_publish_task", action, errs.AlibabaCloudSdkGoERROR)
	}
	// ListCategory returns the category set under $.Categories; some gateways
	// wrap it once more under $.Data.Categories. Accept either shape.
	var categories []interface{}
	if v, err := jsonpath.Get("$.Categories", response); err == nil {
		categories, _ = v.([]interface{})
	} else if v, err := jsonpath.Get("$.Data.Categories", response); err == nil {
		categories, _ = v.([]interface{})
	}
	ids := make([]string, 0)
	for _, v := range categories {
		item, ok := v.(map[string]interface{})
		if !ok {
			continue
		}
		id := fmt.Sprint(item["CategoryId"])
		if id == "" || id == "<nil>" {
			continue
		}
		ids = append(ids, id)
	}
	return ids, nil
}

// CreateFaqCategory creates a temporary FAQ category on the caller's default
// business space and returns its id. Used by the acceptance fixture to ensure
// the data_id_list reorder coverage has two real FAQ category ids when the
// account only exposes one; the caller deletes the created category via
// DeleteFaqCategory in a t.Cleanup step. The values are real FAQ category ids,
// not fabricated placeholders.
func (s *Service) CreateFaqCategory(name string) (string, error) {
	client := s.client
	action := "CreateCategory"
	request := map[string]interface{}{
		"RegionId":         client.RegionId,
		"Name":             name,
		"KnowledgeType":    1, // 1 = faq knowledge base
		"ParentCategoryId": -1,
	}
	response, err := client.RpcPost("Chatbot", "2022-04-08", action, nil, request, true)
	if err != nil {
		return "", errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_chatbot_publish_task", action, errs.AlibabaCloudSdkGoERROR)
	}
	v, err := jsonpath.Get("$.Category.CategoryId", response)
	if err != nil {
		return "", errs.WrapErrorf(err, errs.FailedGetAttributeMsg, "category", "$.Category.CategoryId", response)
	}
	return fmt.Sprint(v), nil
}

// DeleteFaqCategory removes a FAQ category by id. Used by the acceptance
// fixture cleanup to delete the temporary category created by CreateFaqCategory.
// Already-published categories move to a deleted-unpublished state; unpublished
// categories are removed directly.
func (s *Service) DeleteFaqCategory(id string) error {
	client := s.client
	action := "DeleteCategory"
	categoryIdInt, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, id, "strconv.ParseInt", errs.AlibabaCloudSdkGoERROR)
	}
	request := map[string]interface{}{
		"RegionId":   client.RegionId,
		"CategoryId": categoryIdInt,
	}
	_, err = client.RpcPost("Chatbot", "2022-04-08", action, nil, request, true)
	if err != nil {
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}
	return nil
}
