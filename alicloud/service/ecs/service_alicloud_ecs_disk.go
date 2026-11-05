package ecs

import (
	"fmt"
	"strings"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// SetResourceTags <<< Encapsulated tag function for Ecs.
func (s *Service) SetResourceTags(d *schema.ResourceData, resourceType string) error {
	if d.HasChange("tags") {
		var err error
		var action string
		client := s.client
		var request map[string]interface{}
		var response map[string]interface{}
		query := make(map[string]interface{})

		added, removed := parsingTags(d)
		removedTagKeys := make([]string, 0)
		for _, v := range removed {
			if !ignoredTags(v, "") {
				removedTagKeys = append(removedTagKeys, v)
			}
		}
		if len(removedTagKeys) > 0 {
			action = "UntagResources"
			request = make(map[string]interface{})
			query = make(map[string]interface{})
			request["ResourceId.1"] = d.Id()
			request["RegionId"] = client.RegionId
			for i, key := range removedTagKeys {
				request[fmt.Sprintf("TagKey.%d", i+1)] = key
			}

			request["ResourceType"] = resourceType
			wait := incrementalWait(3*time.Second, 5*time.Second)
			err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
				response, err = client.RpcPost("Ecs", "2014-05-26", action, query, request, false)
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
				return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
			}

		}

		if len(added) > 0 {
			action = "TagResources"
			request = make(map[string]interface{})
			query = make(map[string]interface{})
			request["ResourceId.1"] = d.Id()
			request["RegionId"] = client.RegionId
			count := 1
			for key, value := range added {
				request[fmt.Sprintf("Tag.%d.Key", count)] = key
				request[fmt.Sprintf("Tag.%d.Value", count)] = value
				count++
			}

			request["ResourceType"] = resourceType
			wait := incrementalWait(3*time.Second, 5*time.Second)
			err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
				response, err = client.RpcPost("Ecs", "2014-05-26", action, query, request, false)
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
				return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
			}

		}
	}

	return nil
}

// SetResourceTags >>> tag function encapsulated.

// DescribeDisk <<< Encapsulated get interface for Ecs Disk.

func (s *Service) DescribeDisk(id string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	action := "DescribeDisks"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["DiskIds"] = convertListToJsonString([]interface{}{id})
	request["RegionId"] = client.RegionId

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(1*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("Ecs", "2014-05-26", action, query, request, true)

		if err != nil {
			if IsExpectedErrors(err, []string{"InternalError"}) || NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)
	if err != nil {
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}

	v, err := jsonpath.Get("$.Disks.Disk[*]", response)
	if err != nil {
		return object, WrapErrorf(NotFoundErr("Disk", id), NotFoundMsg, response)
	}

	if len(v.([]interface{})) == 0 {
		return object, WrapErrorf(NotFoundErr("Disk", id), NotFoundMsg, response)
	}

	currentStatus := v.([]interface{})[0].(map[string]interface{})["Status"]
	if currentStatus == nil {
		return object, WrapErrorf(NotFoundErr("Disk", id), NotFoundMsg, response)
	}

	return v.([]interface{})[0].(map[string]interface{}), nil
}

func (s *Service) DiskStateRefreshFunc(id string, field string, failStates []string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeDisk(id)
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

// DescribeDisk >>> Encapsulated.

// diskSweepTimeout bounds each step of the forced-release cleanup route.
const diskSweepTimeout = 5 * time.Minute

// DiskSweepPrefixes are the cloud-side disk name prefixes generated by this
// provider's disk acceptance tests and documentation examples. Matching is
// case-sensitive so only names this provider generates are considered.
var DiskSweepPrefixes = []string{"tf-testAcc", "tf-testacc", "terraform-example", "tf-example", "tf_example"}

// IsDiskSweepCandidateName reports whether name was generated by a disk
// acceptance test or documentation example of this provider.
func IsDiskSweepCandidateName(name string) bool {
	for _, prefix := range DiskSweepPrefixes {
		if strings.HasPrefix(name, prefix) {
			return true
		}
	}
	return false
}

// SweepDisk releases one disk left behind by an interrupted acceptance test.
// DeleteDisk only accepts PostPaid disks, so a PrePaid disk is converted,
// detached with DeleteWithInstance=false and waited for until it is available
// before it is released. Every step is bounded by diskSweepTimeout and an error
// is returned instead of being logged as a successful delete.
func (s *Service) SweepDisk(id string) error {
	object, err := s.DescribeDisk(id)
	if err != nil {
		if NotFoundError(err) {
			return nil
		}
		return WrapError(err)
	}
	return forceDeleteDisk(id, diskSweepTimeout, s.client, s, object, "")
}
