package ecs

import (
	"fmt"
	"time"

	"github.com/PaesslerAG/jsonpath"
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
	request["SnapshotIds"] = expandSingletonToList(id)
	request["RegionId"] = client.RegionId
	action := "DescribeLockedSnapshots"

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(1*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("Ecs", "2014-05-26", action, query, request, true)

		if err != nil {
			if NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)
	if err != nil {
		if IsExpectedErrors(err, []string{"InvalidSnapshotLock.NotFound", "InvalidSnapshotId.NotFound"}) {
			return object, WrapErrorf(NotFoundErr("SnapshotLock", id), NotFoundMsg, response)
		}
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}

	v, err := jsonpath.Get("$.LockedSnapshotsInfo[*]", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, id, "$.LockedSnapshotsInfo[*]", response)
	}

	for _, row := range v.([]interface{}) {
		if item, ok := row.(map[string]interface{}); ok && fmt.Sprint(item["SnapshotId"]) == id {
			return item, nil
		}
	}

	return object, WrapErrorf(NotFoundErr("SnapshotLock", id), NotFoundMsg, response)
}

// DescribeSnapshotLock >>> Encapsulated.
