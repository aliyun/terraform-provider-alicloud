// Package cms hosts the Framework-native Cms EndpointConnector resource and
// data source, hand-written from the Cms 2024-03-30 CloudSpec definition.
package cms

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	logs "github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
)

// Service serves the Cms EndpointConnector resources in this package. The root
// alicloud package cannot be imported from a service package, so the Describe
// and List calls live here alongside the resource.
type Service struct {
	client *connectivity.AliyunClient
}

func NewService(client *connectivity.AliyunClient) *Service {
	return &Service{client: client}
}

// endpointConnectorID splits the composite id "<workspace>:<connectorId>".
func endpointConnectorID(id string) (workspace, connectorId string, err error) {
	parts := strings.Split(id, ":")
	if len(parts) != 2 {
		return "", "", fmt.Errorf("invalid Resource Id %s. Expected parts' length %d, got %d", id, 2, len(parts))
	}
	return parts[0], parts[1], nil
}

// DescribeEndpointConnector calls GetEndpointConnector (ROA GET).
func (s *Service) DescribeEndpointConnector(id string) (object map[string]interface{}, err error) {
	client := s.client
	workspace, connectorId, err := endpointConnectorID(id)
	if err != nil {
		return object, errs.WrapError(err)
	}
	action := fmt.Sprintf("/api/v1/endpoint-connectors/%s/%s", workspace, connectorId)
	query := make(map[string]*string)
	var response map[string]interface{}
	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(1*time.Minute, func() *retry.RetryError {
		response, err = client.RoaGet("Cms", "2024-03-30", action, query, nil, nil)
		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(action, response)
	if err != nil {
		if errs.IsExpectedErrors(err, []string{"EndpointConnectorNotExist", "WorkspaceNotExist"}) {
			return object, errs.WrapErrorf(errs.NotFoundErr("Cms:EndpointConnector", id), errs.NotFoundMsg, response)
		}
		return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}
	return response, nil
}

// mapInt64 coerces an API response number (json.Number, float64, int, int64 or
// numeric string) into a types.Int64. The tea-rpc client may decode numbers as
// json.Number or int rather than float64, so every numeric encoding the
// transport produces must map back to Int64 or planned values would collapse to
// null ("inconsistent result after apply").
func mapInt64(object map[string]interface{}, key string) types.Int64 {
	switch v := object[key].(type) {
	case float64:
		return types.Int64Value(int64(v))
	case json.Number:
		if n, err := v.Int64(); err == nil {
			return types.Int64Value(n)
		}
	case int:
		return types.Int64Value(int64(v))
	case int64:
		return types.Int64Value(v)
	case string:
		if n, err := strconv.ParseInt(v, 10, 64); err == nil {
			return types.Int64Value(n)
		}
	}
	return types.Int64Null()
}

func (s *Service) ListEndpointConnectors(workspace, connectorType, name string) (objects []map[string]interface{}, err error) {
	client := s.client
	action := "/api/v1/endpoint-connectors"
	query := make(map[string]*string)
	query["workspace"] = helper.StringPointer(workspace)
	if connectorType != "" {
		query["type"] = helper.StringPointer(connectorType)
	}
	if name != "" {
		query["name"] = helper.StringPointer(name)
	}
	query["maxResults"] = helper.StringPointer(fmt.Sprintf("%d", helper.PageSizeLarge))
	for {
		var response map[string]interface{}
		wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
		err = retry.Retry(1*time.Minute, func() *retry.RetryError {
			response, err = client.RoaGet("Cms", "2024-03-30", action, query, nil, nil)
			if err != nil {
				if helper.NeedRetry(err) {
					wait()
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
			}
			return nil
		})
		logs.AddDebug(action, response, query)
		if err != nil {
			return objects, errs.WrapErrorf(err, errs.DataDefaultErrorMsg, "alicloud_cms_endpoint_connectors", action, errs.AlibabaCloudSdkGoERROR)
		}
		resp, _ := jsonpath.Get("$.endpointConnectors[*]", response)
		result, _ := resp.([]interface{})
		for _, v := range result {
			if item, ok := v.(map[string]interface{}); ok {
				objects = append(objects, item)
			}
		}
		if nextToken, ok := response["nextToken"].(string); ok && nextToken != "" {
			query["nextToken"] = helper.StringPointer(nextToken)
		} else {
			break
		}
	}
	return objects, nil
}
