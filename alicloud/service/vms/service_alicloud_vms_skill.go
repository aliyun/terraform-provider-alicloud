package vms

import (
	"context"
	"fmt"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
)

// createSkill calls CloudCreateSkill and returns the response Data object, which
// carries the new skill Id (PK) alongside Name, Comment, EnterpriseId and
// CreateTime, so the resource state can be set without an extra read.
func (s *Service) createSkill(ctx context.Context, enterpriseId int64, name, comment string, timeout time.Duration) (map[string]interface{}, error) {
	action := "CloudCreateSkill"
	request := map[string]interface{}{
		"EnterpriseId": enterpriseId,
		"Name":         name,
	}
	if comment != "" {
		request["Comment"] = comment
	}
	response, err := retryRpcPost(ctx, s.client, action, request, timeout)
	if err != nil {
		return nil, errs.WrapErrorf(err, errs.DefaultErrorMsg, fmt.Sprintf("%d:%s", enterpriseId, name), action, errs.AlibabaCloudSdkGoERROR)
	}
	logs.AddDebug(action, response, request)
	data, ok := response["Data"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("creating alicloud_vms_skill: CloudCreateSkill response carries no Data: %#v", response)
	}
	if _, err := toInt64(data["Id"]); err != nil {
		return nil, fmt.Errorf("creating alicloud_vms_skill: CloudCreateSkill response Data.Id is not a valid integer: %#v", data["Id"])
	}
	return data, nil
}

// describeSkill calls CloudListSkill paginating with Limit/Offset until it finds
// the skill whose Id matches skillId. It returns the matching list item, or
// (nil, nil) when no skill with that Id remains in the enterprise, which the
// caller treats as a not-found resource.
func (s *Service) describeSkill(ctx context.Context, enterpriseId, skillId int64, timeout time.Duration) (map[string]interface{}, error) {
	action := "CloudListSkill"
	budget := timeout
	offset := 0
	for {
		// Recompute the per-call retry budget against the operation deadline so the
		// pagination does not receive a fresh full budget of its own each page.
		if deadline, ok := ctx.Deadline(); ok {
			if remaining := time.Until(deadline); remaining < budget {
				budget = remaining
			}
		}
		if budget <= 0 {
			if err := ctx.Err(); err != nil {
				return nil, err
			}
			return nil, context.DeadlineExceeded
		}
		request := map[string]interface{}{
			"EnterpriseId": enterpriseId,
			"Limit":        skillListPageSize,
			"Offset":       offset,
		}
		response, err := retryRpcPost(ctx, s.client, action, request, budget)
		if err != nil {
			return nil, errs.WrapErrorf(err, errs.DefaultErrorMsg, fmt.Sprintf("%d:%d", enterpriseId, skillId), action, errs.AlibabaCloudSdkGoERROR)
		}
		logs.AddDebug(action, response, request)
		data, ok := response["Data"].(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("reading alicloud_vms_skill: CloudListSkill response carries no Data: %#v", response)
		}
		list, _ := data["List"].([]interface{})
		for _, item := range list {
			obj, ok := item.(map[string]interface{})
			if !ok {
				continue
			}
			id, err := toInt64(obj["Id"])
			if err != nil {
				continue
			}
			if id == skillId {
				return obj, nil
			}
		}
		if !hasNextPage(data) {
			return nil, nil
		}
		offset += skillListPageSize
	}
}

// hasNextPage reports whether CloudListSkill's Data advertises a following page.
func hasNextPage(data map[string]interface{}) bool {
	if v, ok := data["HasNextPage"].(bool); ok {
		return v
	}
	return false
}

// deleteSkill calls CloudDeleteSkill by Id; EnterpriseId is sent as well so the
// call names the owning enterprise. CloudDeleteSkill documents that Id takes
// precedence when both Id and Name are passed. A not-found result is treated as
// success, because the skill is already gone and a re-run destroy must not fail.
func (s *Service) deleteSkill(ctx context.Context, enterpriseId, skillId int64, timeout time.Duration) error {
	action := "CloudDeleteSkill"
	request := map[string]interface{}{
		"EnterpriseId": enterpriseId,
		"Id":           skillId,
	}
	_, err := retryRpcPost(ctx, s.client, action, request, timeout)
	if err != nil {
		if errs.NotFoundError(err) {
			return nil
		}
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, fmt.Sprintf("%d:%d", enterpriseId, skillId), action, errs.AlibabaCloudSdkGoERROR)
	}
	logs.AddDebug(action, nil, request)
	return nil
}
