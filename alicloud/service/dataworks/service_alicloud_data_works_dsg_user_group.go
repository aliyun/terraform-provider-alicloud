package dataworks

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
)

const (
	// dsgUserGroupTypeName is the public Terraform resource type name.
	dsgUserGroupTypeName = "alicloud_data_works_dsg_user_group"
	// dsgUserGroupTimeout replaces the SDKv2 resource-level create/update/delete
	// default timeouts (5 minutes each); the Framework does not support
	// user-configurable resource timeouts yet.
	dsgUserGroupTimeout = 5 * time.Minute
	// dsgUserGroupMaxPages bounds the QueryList pagination so an unbounded list
	// cannot stall Read/refresh indefinitely; reaching it is treated as a
	// search failure rather than NotFound.
	dsgUserGroupMaxPages = 100
)

// DsgUserGroupAddOrUpdate creates (no Id) or updates (with Id) a data
// desensitization user group via the DsgUserGroupAddOrUpdate API.
func (s *Service) DsgUserGroupAddOrUpdate(ctx context.Context, request map[string]interface{}) error {
	action := "DsgUserGroupAddOrUpdate"
	var response map[string]interface{}
	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err := retry.RetryContext(ctx, dsgUserGroupTimeout, func() *retry.RetryError {
		var callErr error
		response, callErr = s.client.RpcPost("dataworks-public", "2020-05-18", action, nil, request, true)
		if callErr != nil {
			if helper.NeedRetry(callErr) {
				wait()
				return retry.RetryableError(callErr)
			}
			return retry.NonRetryableError(callErr)
		}
		return nil
	})
	logs.AddDebug(action, response, request)
	if err != nil {
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, dsgUserGroupTypeName, action, errs.AlibabaCloudSdkGoERROR)
	}
	return nil
}

// DsgUserGroupDescribe reads a single user group by Id. The DsgUserGroupQueryList
// API has no Id filter, so when name is known it is used as a server-side filter
// and the matching Id is selected client-side; when name is empty (e.g. after a
// passthrough import) the full list is paginated until the Id is found. The API
// response carries no single Get-by-Id, so NotFound is only concluded after the
// configured pagination boundary is exhausted.
func (s *Service) DsgUserGroupDescribe(ctx context.Context, id string, name string) (map[string]interface{}, error) {
	action := "DsgUserGroupQueryList"
	for page := 1; page <= dsgUserGroupMaxPages; page++ {
		request := map[string]interface{}{
			"RegionId":   s.client.RegionId,
			"PageNumber": page,
			"PageSize":   100,
		}
		if name != "" {
			request["Name"] = name
		}
		var response map[string]interface{}
		wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
		err := retry.RetryContext(ctx, dsgUserGroupTimeout, func() *retry.RetryError {
			var callErr error
			response, callErr = s.client.RpcPost("dataworks-public", "2020-05-18", action, nil, request, true)
			if callErr != nil {
				if helper.NeedRetry(callErr) {
					wait()
					return retry.RetryableError(callErr)
				}
				return retry.NonRetryableError(callErr)
			}
			return nil
		})
		logs.AddDebug(action, response, request)
		if err != nil {
			return nil, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
		}
		pageData, ok := response["PageData"].(map[string]interface{})
		if !ok {
			break
		}
		data, _ := pageData["Data"].([]interface{})
		for _, raw := range data {
			item, ok := raw.(map[string]interface{})
			if !ok {
				continue
			}
			// When id is empty the caller just created the group and does not
			// know its Id yet, so the first entry matching the name filter is
			// returned; otherwise the entry whose Id matches is selected.
			if id == "" || fmt.Sprint(item["Id"]) == id {
				return item, nil
			}
		}
		totalCount := toInt(pageData["TotalCount"])
		pageSize := toInt(pageData["PageSize"])
		if pageSize <= 0 {
			pageSize = 100
		}
		if len(data) == 0 || page*pageSize >= totalCount {
			break
		}
	}
	return nil, errs.WrapErrorf(errs.NotFoundErr("dataworks", id), errs.NotFoundMsg, map[string]interface{}{"Id": id})
}

// DsgUserGroupDelete removes a user group by Id via the DsgUserGroupDelete API.
func (s *Service) DsgUserGroupDelete(ctx context.Context, id string) error {
	action := "DsgUserGroupDelete"
	idInt, perr := strconv.ParseInt(id, 10, 64)
	if perr != nil {
		return errs.WrapError(perr)
	}
	idsJSON, err := json.Marshal([]interface{}{idInt})
	if err != nil {
		return errs.WrapError(err)
	}
	request := map[string]interface{}{
		"RegionId": s.client.RegionId,
		"Ids":      string(idsJSON),
	}
	var response map[string]interface{}
	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.RetryContext(ctx, dsgUserGroupTimeout, func() *retry.RetryError {
		var callErr error
		response, callErr = s.client.RpcPost("dataworks-public", "2020-05-18", action, nil, request, true)
		if callErr != nil {
			if helper.NeedRetry(callErr) {
				wait()
				return retry.RetryableError(callErr)
			}
			return retry.NonRetryableError(callErr)
		}
		return nil
	})
	logs.AddDebug(action, response, request)
	if err != nil {
		// A user group that no longer exists leaves nothing to delete. The
		// published errorCodes table maps the not-exist case to
		// USERGROUP.ID.ERROR ("does not match the tenant or does not exist").
		if errs.IsExpectedErrors(err, []string{"USERGROUP.ID.ERROR"}) || errs.NotFoundError(err) {
			return nil
		}
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}
	return nil
}

// toInt coerces a JSON-decoded numeric value to int, accepting float64 and
// json.Number; non-numeric values yield 0.
func toInt(v interface{}) int {
	switch n := v.(type) {
	case float64:
		return int(n)
	case int:
		return n
	case int64:
		return int(n)
	}
	return 0
}
