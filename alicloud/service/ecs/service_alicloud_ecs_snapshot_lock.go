package ecs

import (
	"fmt"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/convert"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	logs "github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
)

// DescribeSnapshotLock <<< Encapsulated get interface for Ecs SnapshotLock.

func (s *Service) DescribeSnapshotLock(id string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["SnapshotIds"] = convert.ExpandSingletonToList(id)
	request["RegionId"] = client.RegionId
	action := "DescribeLockedSnapshots"

	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(1*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("Ecs", "2014-05-26", action, query, request, true)

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
		if errs.IsExpectedErrors(err, []string{"InvalidSnapshotLock.NotFound", "InvalidSnapshotId.NotFound"}) {
			return object, errs.WrapErrorf(errs.NotFoundErr("SnapshotLock", id), errs.NotFoundMsg, response)
		}
		return object, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}

	v, err := jsonpath.Get("$.LockedSnapshotsInfo[*]", response)
	if err != nil {
		return object, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.LockedSnapshotsInfo[*]", response)
	}

	for _, row := range v.([]interface{}) {
		if item, ok := row.(map[string]interface{}); ok && fmt.Sprint(item["SnapshotId"]) == id {
			return item, nil
		}
	}

	return object, errs.WrapErrorf(errs.NotFoundErr("SnapshotLock", id), errs.NotFoundMsg, response)
}

// DescribeSnapshotLock >>> Encapsulated.
