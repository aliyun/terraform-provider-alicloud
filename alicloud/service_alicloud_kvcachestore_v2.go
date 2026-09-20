// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"strings"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

type KvcachestoreServiceV2 struct {
	client *connectivity.AliyunClient
}

// DescribeKvcachestoreKVCacheStore <<< Encapsulated get interface for Kvcachestore KVCacheStore.

func (s *KvcachestoreServiceV2) DescribeKvcachestoreKVCacheStore(id string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["KvcsId"] = id
	request["RegionId"] = client.RegionId
	action := "GetKVCacheStore"

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(1*time.Minute, func() *resource.RetryError {
		response, err = client.RpcPost("Kvcachestore", "2026-06-17", action, query, request, true)
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
		if IsExpectedErrors(err, []string{"NoSuchResource"}) {
			return object, WrapErrorf(NotFoundErr("KVCacheStore", id), NotFoundMsg, response)
		}
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}

	v, err := jsonpath.Get("$.KvCacheStore", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.KvCacheStore", response)
	}
	return v.(map[string]interface{}), nil
}

func (s *KvcachestoreServiceV2) KvcachestoreKVCacheStoreStateRefreshFunc(id string, field string, failStates []string) resource.StateRefreshFunc {
	return s.KvcachestoreKVCacheStoreStateRefreshFuncWithApi(id, field, failStates, s.DescribeKvcachestoreKVCacheStore)
}

func (s *KvcachestoreServiceV2) KvcachestoreKVCacheStoreStateRefreshFuncWithApi(id string, field string, failStates []string, call func(id string) (map[string]interface{}, error)) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := call(id)
		if err != nil {
			if NotFoundError(err) {
				return object, "", nil
			}
			return nil, "", WrapError(err)
		}
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
				return object, currentStatus, WrapError(Error(FailedToReachTargetStatus, currentStatus))
			}
		}
		return object, currentStatus, nil
	}
}

// DescribeKvcachestoreKVCacheStore >>> Encapsulated.

// ListKvcachestoreKVCacheStores <<< Encapsulated list interface for Kvcachestore KVCacheStore (data source).

func (s *KvcachestoreServiceV2) ListKvcachestoreKVCacheStores(query map[string]interface{}) (response map[string]interface{}, err error) {
	client := s.client
	action := "ListKVCacheStores"
	request := map[string]interface{}{}
	for k, v := range query {
		request[k] = v
	}
	request["RegionId"] = client.RegionId

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = client.RpcPost("Kvcachestore", "2026-06-17", action, nil, request, true)
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
		return response, WrapErrorf(err, DataDefaultErrorMsg, "alicloud_kv_cache_stores", action, AlibabaCloudSdkGoERROR)
	}
	return response, nil
}

// ListKvcachestoreKVCacheStores >>> Encapsulated.
