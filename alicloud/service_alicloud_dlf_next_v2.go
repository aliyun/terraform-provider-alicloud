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

type DlfNextServiceV2 struct {
	client *connectivity.AliyunClient
}

// DescribeDlfNextCatalog <<< Encapsulated get interface for DlfNext Catalog.
func (s *DlfNextServiceV2) DescribeDlfNextCatalog(id string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]*string

	catalog := id
	request = make(map[string]interface{})
	query = make(map[string]*string)

	action := fmt.Sprintf("/dlf/v1/catalogs/%s", catalog)

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(1*time.Minute, func() *resource.RetryError {
		response, err = client.RoaGet("DlfNext", "2025-03-10", action, query, nil, nil)
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
		if IsExpectedErrors(err, []string{"CatalogNotFound", "NotFound"}) {
			return object, WrapErrorf(NotFoundErr("DlfCatalog", id), NotFoundMsg, response)
		}
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}

	return response, nil
}

// ListDlfNextCatalogs <<< Encapsulated list interface for DlfNext Catalog.
func (s *DlfNextServiceV2) ListDlfNextCatalogs(catalogNamePattern string, pageToken string, maxResults int) (objects []interface{}, nextToken string, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	action := fmt.Sprintf("/dlf/v1/catalogs")
	request = make(map[string]interface{})
	query := make(map[string]*string)

	if catalogNamePattern != "" {
		v := catalogNamePattern
		query["catalogNamePattern"] = &v
	}
	if pageToken != "" {
		v := pageToken
		query["pageToken"] = &v
	}
	if maxResults > 0 {
		v := fmt.Sprint(maxResults)
		query["maxResults"] = &v
	}

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(1*time.Minute, func() *resource.RetryError {
		response, err = client.RoaGet("DlfNext", "2025-03-10", action, query, nil, nil)
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
		return objects, nextToken, WrapErrorf(err, DefaultErrorMsg, "list", action, AlibabaCloudSdkGoERROR)
	}

	nextToken = ""
	if v, ok := response["nextPageToken"]; ok && v != nil {
		nextToken = fmt.Sprint(v)
	}

	if v, ok := response["catalogs"].([]interface{}); ok {
		objects = v
	}

	return objects, nextToken, nil
}

func (s *DlfNextServiceV2) DlfNextCatalogStateRefreshFunc(id string, field string, failStates []string) resource.StateRefreshFunc {
	return s.DlfNextCatalogStateRefreshFuncWithApi(id, field, failStates, s.DescribeDlfNextCatalog)
}

func (s *DlfNextServiceV2) DlfNextCatalogStateRefreshFuncWithApi(id string, field string, failStates []string, call func(id string) (map[string]interface{}, error)) resource.StateRefreshFunc {
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

// DescribeDlfNextCatalog >>> Encapsulated.
