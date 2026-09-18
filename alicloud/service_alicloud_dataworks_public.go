package alicloud

import (
	"strconv"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

type DataworksPublicService struct {
	client *connectivity.AliyunClient
}

func (s *DataworksPublicService) DescribeDataWorksFolder(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	client := s.client
	action := "GetFolder"
	parts, err := ParseResourceId(id, 2)
	if err != nil {
		err = WrapError(err)
		return
	}
	request := map[string]interface{}{
		"FolderId":  parts[0],
		"ProjectId": parts[1],
	}
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = client.RpcPost("dataworks-public", "2020-05-18", action, nil, request, true)
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
	v, err := jsonpath.Get("$.Data", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.Data", response)
	}
	object = v.(map[string]interface{})
	if len(object) < 1 {
		return object, WrapErrorf(NotFoundErr("dataworks", id), NotFoundWithResponse, response)
	}
	return object, nil
}

func (s *DataworksPublicService) GetFolder(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	client := s.client
	action := "GetFolder"
	parts, err := ParseResourceId(id, 2)
	if err != nil {
		err = WrapError(err)
		return
	}
	request := map[string]interface{}{
		"FolderId":  parts[0],
		"ProjectId": parts[1],
	}
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = client.RpcPost("dataworks-public", "2020-05-18", action, nil, request, true)
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
	v, err := jsonpath.Get("$.Data", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.Data", response)
	}
	object = v.(map[string]interface{})
	return object, nil
}

// DescribeDataWorksMetaCategory reads a single MetaCategory resource by its composite id
// "<CategoryId>:<ParentCategoryId>". Because the upstream GetMetaCategory API only supports
// querying by parent (it returns the children of the given ParentCategoryId), the resource id
// carries the parent along with the category id; Read queries the parent and then searches the
// returned DataEntityList for the entry whose CategoryId matches.
func (s *DataworksPublicService) DescribeDataWorksMetaCategory(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	client := s.client
	action := "GetMetaCategory"
	parts, err := ParseResourceId(id, 2)
	if err != nil {
		err = WrapError(err)
		return
	}
	categoryId, _ := strconv.Atoi(parts[0])
	parentCategoryId, _ := strconv.Atoi(parts[1])
	request := map[string]interface{}{
		"ParentCategoryId": parentCategoryId,
		"PageNum":          1,
		"PageSize":         100,
	}
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = client.RpcPost("dataworks-public", "2020-05-18", action, nil, request, true)
		if err != nil {
			if NotFoundError(err) {
				return nil
			}
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
	v, err := jsonpath.Get("$.Data.DataEntityList", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.Data.DataEntityList", response)
	}
	list, ok := v.([]interface{})
	if !ok || len(list) == 0 {
		return object, WrapErrorf(NotFoundErr("dataworks meta category", id), NotFoundWithResponse, response)
	}
	for _, item := range list {
		obj, ok := item.(map[string]interface{})
		if !ok {
			continue
		}
		cid, ok := obj["CategoryId"]
		if !ok || cid == nil {
			continue
		}
		// RpcPost uses UseNumber decoding, so CategoryId arrives as json.Number;
		// a .(float64) assertion returns ok=false, the loop never matches, and
		// Read falls through to NotFoundErr -> d.SetId("") -> Apply post-check
		// reports the resource "was present, but now absent". toInt tolerates
		// json.Number (and int/float/string) so the match succeeds.
		if cidVal, err := toInt(cid); err == nil && cidVal == categoryId {
			return obj, nil
		}
	}
	return object, WrapErrorf(NotFoundErr("dataworks meta category", id), NotFoundWithResponse, response)
}
