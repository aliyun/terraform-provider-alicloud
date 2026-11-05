package ecs

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"testing"
	"time"

	utilV2 "github.com/alibabacloud-go/tea-utils/v2/service"
	"github.com/alibabacloud-go/tea/dara"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUnitFwInt64DecodesJSONNumber(t *testing.T) {
	var decoded map[string]interface{}
	decoder := json.NewDecoder(strings.NewReader(`{"InstancePoolsToUseCount":3,"InstanceAmount":2,"TotalCount":7,"Fractional":4.0}`))
	decoder.UseNumber()
	assert.NoError(t, decoder.Decode(&decoded))

	assert.Equal(t, types.Int64Value(3), fwInt64(decoded["InstancePoolsToUseCount"]))
	assert.Equal(t, types.Int64Value(2), fwInt64(decoded["InstanceAmount"]))
	assert.Equal(t, types.Int64Value(7), fwInt64(decoded["TotalCount"]))
	assert.Equal(t, types.Int64Value(4), fwInt64(decoded["Fractional"]))
}

func TestUnitFwTagMapValidator(t *testing.T) {
	assertTagMap := func(v fwTagMapValidator, tags map[string]attr.Value, wantError bool) {
		var response validator.MapResponse
		v.ValidateMap(context.Background(), validator.MapRequest{
			ConfigValue: types.MapValueMust(types.StringType, tags),
			Path:        path.Root("tags"),
		}, &response)
		assert.Equal(t, wantError, response.Diagnostics.HasError(), "tags: %v", tags)
	}

	apg := fwTagMapValidator{allowEmptyValue: true}
	assertTagMap(apg, map[string]attr.Value{"env": types.StringValue("")}, false)
	assertTagMap(apg, map[string]attr.Value{"env": types.StringValue("test")}, false)
	assertTagMap(apg, map[string]attr.Value{"": types.StringValue("test")}, true)
	assertTagMap(apg, map[string]attr.Value{"aliyun-tag": types.StringValue("test")}, true)
	assertTagMap(apg, map[string]attr.Value{"acs:tag": types.StringValue("test")}, true)
	assertTagMap(apg, map[string]attr.Value{"env": types.StringValue("https://example.com")}, true)
	assertTagMap(apg, map[string]attr.Value{strings.Repeat("k", 129): types.StringValue("test")}, true)
	assertTagMap(apg, map[string]attr.Value{"env": types.StringValue(strings.Repeat("v", 129))}, true)

	reserved := fwTagMapValidator{valueForbiddenPrefixes: []string{"acs:"}}
	assertTagMap(reserved, map[string]attr.Value{"env": types.StringValue("")}, true)
	assertTagMap(reserved, map[string]attr.Value{"env": types.StringValue("acs:reserved")}, true)
	assertTagMap(reserved, map[string]attr.Value{"env": types.StringValue("http://example.com")}, true)
	assertTagMap(reserved, map[string]attr.Value{"env": types.StringValue("test")}, false)
}

// A null value must not skip the key rules, and it must not be sent as an empty
// string: null, unknown and an explicit empty string are three different
// requests.
func TestUnitFwTagMapValidatorNullValues(t *testing.T) {
	validate := func(v fwTagMapValidator, tags map[string]attr.Value) (bool, string) {
		var response validator.MapResponse
		v.ValidateMap(context.Background(), validator.MapRequest{
			ConfigValue: types.MapValueMust(types.StringType, tags),
			Path:        path.Root("tags"),
		}, &response)
		if !response.Diagnostics.HasError() {
			return false, ""
		}
		return true, response.Diagnostics.Errors()[0].Summary()
	}

	apg := fwTagMapValidator{allowEmptyValue: true}
	hasError, _ := validate(apg, map[string]attr.Value{"env": types.StringNull()})
	assert.True(t, hasError, "a null value must be rejected")
	hasError, summary := validate(apg, map[string]attr.Value{"acs:env": types.StringNull()})
	assert.True(t, hasError, "the key rules apply even when the value is null")
	assert.Contains(t, summary, "Invalid Tag Key")
	hasError, _ = validate(apg, map[string]attr.Value{"env": types.StringUnknown()})
	assert.False(t, hasError, "an unknown value cannot be validated yet")
	hasError, _ = validate(apg, map[string]attr.Value{"env": types.StringValue("")})
	assert.False(t, hasError, "an explicit empty value stays allowed for this API")

	reserved := fwTagMapValidator{valueForbiddenPrefixes: []string{"acs:"}}
	hasError, _ = validate(reserved, map[string]attr.Value{"env": types.StringNull()})
	assert.True(t, hasError)
	hasError, _ = validate(reserved, map[string]attr.Value{"env": types.StringValue("")})
	assert.True(t, hasError)
}

func TestUnitFwApplyTagsSerialization(t *testing.T) {
	request := map[string]interface{}{}
	fwApplyTags(request, types.MapValueMust(types.StringType, map[string]attr.Value{
		"env":   types.StringValue("test"),
		"empty": types.StringValue(""),
		"null":  types.StringNull(),
		"known": types.StringUnknown(),
	}))

	serialized := map[string]string{}
	for key, value := range request {
		if !strings.HasPrefix(key, "Tag.") || !strings.HasSuffix(key, ".Key") {
			continue
		}
		index := strings.TrimSuffix(strings.TrimPrefix(key, "Tag."), ".Key")
		serialized[value.(string)] = request[fmt.Sprintf("Tag.%s.Value", index)].(string)
	}
	assert.Equal(t, map[string]string{"env": "test", "empty": ""}, serialized, "null and unknown values must not be serialized")
}

func TestUnitFwDecodeRecordList(t *testing.T) {
	items, err := fwDecodeRecordList(map[string]interface{}{}, "AutoProvisioningGroups", "AutoProvisioningGroup", "Describe")
	assert.NoError(t, err, "a missing member is a valid empty result")
	assert.Empty(t, items)

	items, err = fwDecodeRecordList(map[string]interface{}{
		"AutoProvisioningGroups": map[string]interface{}{
			"AutoProvisioningGroup": []interface{}{
				map[string]interface{}{"AutoProvisioningGroupId": "apg-1"},
			},
		},
	}, "AutoProvisioningGroups", "AutoProvisioningGroup", "Describe")
	require.NoError(t, err)
	assert.Len(t, items, 1)

	_, err = fwDecodeRecordList(map[string]interface{}{
		"AutoProvisioningGroups": "not-an-object",
	}, "AutoProvisioningGroups", "AutoProvisioningGroup", "Describe")
	assert.Error(t, err, "a malformed container must not be treated as empty")

	_, err = fwDecodeRecordList(map[string]interface{}{
		"AutoProvisioningGroups": map[string]interface{}{"AutoProvisioningGroup": "not-a-list"},
	}, "AutoProvisioningGroups", "AutoProvisioningGroup", "Describe")
	assert.Error(t, err, "a malformed list must not be treated as empty")

	_, err = fwDecodeRecordList(map[string]interface{}{
		"AutoProvisioningGroups": map[string]interface{}{
			"AutoProvisioningGroup": []interface{}{"not-an-object"},
		},
	}, "AutoProvisioningGroups", "AutoProvisioningGroup", "Describe")
	assert.Error(t, err, "a malformed record must not be dropped")
}

func TestUnitFwDecodeCountAndToken(t *testing.T) {
	total, known, err := fwDecodeCount(map[string]interface{}{}, "TotalCount", "Describe")
	assert.NoError(t, err)
	assert.False(t, known)
	assert.Zero(t, total)

	total, known, err = fwDecodeCount(map[string]interface{}{"TotalCount": json.Number("7")}, "TotalCount", "Describe")
	assert.NoError(t, err)
	assert.True(t, known)
	assert.Equal(t, int64(7), total)

	_, _, err = fwDecodeCount(map[string]interface{}{"TotalCount": "many"}, "TotalCount", "Describe")
	assert.Error(t, err, "an invalid count must not become 0")

	token, err := fwDecodeToken(map[string]interface{}{}, "NextToken", "Describe")
	assert.NoError(t, err)
	assert.Empty(t, token)

	_, err = fwDecodeToken(map[string]interface{}{"NextToken": 42}, "NextToken", "Describe")
	assert.Error(t, err, "an invalid token must not be treated as the last page")

	_, err = fwDecodeToken(map[string]interface{}{"NextToken": "  "}, "NextToken", "Describe")
	assert.Error(t, err, "a blank token is not a cursor that can be followed")
}

func TestUnitFwDecodeTags(t *testing.T) {
	tags, err := fwDecodeTags(nil)
	require.NoError(t, err)
	assert.Empty(t, tags.Elements(), "an absent member is an empty map")

	tags, err = fwDecodeTags(map[string]interface{}{
		"Tag": []interface{}{
			map[string]interface{}{"TagKey": "env", "TagValue": "test"},
			map[string]interface{}{"TagKey": "aliyun-reserved", "TagValue": "skip"},
		},
	})
	require.NoError(t, err)
	assert.Equal(t, types.StringValue("test"), tags.Elements()["env"])
	assert.NotContains(t, tags.Elements(), "aliyun-reserved")

	_, err = fwDecodeTags(map[string]interface{}{"Tag": "not-a-list"})
	assert.Error(t, err)

	_, err = fwDecodeTags(map[string]interface{}{"Tag": []interface{}{"not-an-object"}})
	assert.Error(t, err, "a malformed tag entry must be reported instead of panicking")

	_, err = fwDecodeTags(map[string]interface{}{"Tag": []interface{}{map[string]interface{}{"TagValue": "test"}}})
	assert.Error(t, err, "a tag without a key is malformed")

	_, err = fwDecodeTags(map[string]interface{}{"Tag": []interface{}{map[string]interface{}{"TagKey": 1, "TagValue": "test"}}})
	assert.Error(t, err, "a non-string tag key is malformed")
}

func TestUnitFwCollectTokenPagesFollowsTokens(t *testing.T) {
	var visited []string
	objects, lastToken, err := fwCollectTokenPages(context.Background(), 10, "", "Discover", "ExecutionId", func(token string) ([]map[string]interface{}, string, error) {
		visited = append(visited, token)
		switch token {
		case "":
			return []map[string]interface{}{{"ExecutionId": "exec-1"}}, "token-1", nil
		case "token-1":
			return []map[string]interface{}{{"ExecutionId": "exec-2"}}, "", nil
		default:
			return nil, "", fmt.Errorf("unexpected token %q", token)
		}
	})

	assert.NoError(t, err)
	assert.Equal(t, []string{"", "token-1"}, visited)
	assert.Len(t, objects, 2)
	assert.Equal(t, "exec-2", objects[1]["ExecutionId"])
	assert.Empty(t, lastToken)
}

func TestUnitFwCollectTokenPagesRejectsRepeatedToken(t *testing.T) {
	calls := 0
	objects, lastToken, err := fwCollectTokenPages(context.Background(), 10, "", "Discover", "ExecutionId", func(token string) ([]map[string]interface{}, string, error) {
		calls++
		return []map[string]interface{}{{"ExecutionId": fmt.Sprintf("exec-%d", calls)}}, "same-token", nil
	})

	assert.Error(t, err)
	assert.Nil(t, objects, "an incomplete result must not be returned")
	assert.Empty(t, lastToken)
	assert.Equal(t, 2, calls, "the repeated token must not be followed again")
}

func TestUnitFwCollectTokenPagesFailsWhenBudgetExhausted(t *testing.T) {
	calls := 0
	objects, _, err := fwCollectTokenPages(context.Background(), 2, "", "Discover", "ExecutionId", func(token string) ([]map[string]interface{}, string, error) {
		calls++
		return []map[string]interface{}{{"ExecutionId": fmt.Sprintf("exec-%d", calls)}}, fmt.Sprintf("token-%d", calls), nil
	})

	assert.Error(t, err)
	assert.Nil(t, objects, "an incomplete result must not be returned")
	assert.Equal(t, 2, calls)
}

func TestUnitFwCollectTokenPagesStopsOnCancelledContext(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	calls := 0
	objects, _, err := fwCollectTokenPages(ctx, 10, "", "Discover", "ExecutionId", func(token string) ([]map[string]interface{}, string, error) {
		calls++
		return nil, "", nil
	})

	assert.Error(t, err)
	assert.Nil(t, objects)
	assert.Zero(t, calls, "a cancelled read must not issue another request")
}

func TestUnitEcsImagePipelineExecutionsStatusValidator(t *testing.T) {
	v := imagePipelineExecutionsStatusValidator{}

	validate := func(value string) bool {
		var response validator.StringResponse
		v.ValidateString(context.Background(), validator.StringRequest{
			ConfigValue: types.StringValue(value),
			Path:        path.Root("status"),
		}, &response)
		return response.Diagnostics.HasError()
	}

	assert.False(t, validate("SUCCESS, FAILED"))
	assert.False(t, validate("SUCCESS"))
	assert.False(t, validate(""), "an empty value means every status and is allowed")
	assert.True(t, validate(" "))
	assert.True(t, validate(","))
	assert.True(t, validate(" , "))
	assert.True(t, validate("SUCCESS,"))
	assert.True(t, validate(",SUCCESS"))
	assert.True(t, validate("SUCCESS,,FAILED"))
	assert.True(t, validate("SUCCESS,not-a-status"))
}

func TestUnitEcsImagePipelineExecutionsNormalizedStatus(t *testing.T) {
	assert.Equal(t, "SUCCESS,FAILED", imagePipelineExecutionsNormalizedStatus(types.StringValue(" SUCCESS , FAILED ")).ValueString())
	assert.True(t, imagePipelineExecutionsNormalizedStatus(types.StringNull()).IsNull())
}

func TestUnitEcsImagePipelineExecutionObject(t *testing.T) {
	object, err := imagePipelineExecutionObject(map[string]interface{}{
		"ExecutionId":     "exec-1",
		"ImagePipelineId": "pipeline-1",
		"Status":          "SUCCESS",
		"CreationTime":    "2026-11-05T10:00:00Z",
		"Tags": map[string]interface{}{
			"Tag": []interface{}{
				map[string]interface{}{"TagKey": "env", "TagValue": "test"},
			},
		},
	})
	require.NoError(t, err)

	attributes := object.Attributes()
	assert.Equal(t, types.StringValue("exec-1"), attributes["execution_id"])
	assert.Equal(t, types.StringValue("pipeline-1"), attributes["image_pipeline_id"])
	assert.Equal(t, types.StringValue("SUCCESS"), attributes["status"])
	assert.Equal(t, types.StringValue("2026-11-05T10:00:00Z"), attributes["creation_time"])
	assert.Equal(t, types.StringValue(""), attributes["error_code"])

	tags, ok := attributes["tags"].(types.Map)
	if assert.True(t, ok) {
		assert.Equal(t, types.StringValue("test"), tags.Elements()["env"])
	}
}

func TestUnitEcsAutoProvisioningGroupsObject(t *testing.T) {
	object, err := autoProvisioningGroupObject(map[string]interface{}{
		"AutoProvisioningGroupId":   "apg-1",
		"AutoProvisioningGroupName": "example",
		"TerminateInstances":        true,
		"LaunchTemplateConfigs": map[string]interface{}{
			"LaunchTemplateConfig": []interface{}{
				map[string]interface{}{
					"InstanceType":     "ecs.g6.large",
					"WeightedCapacity": 2,
					"VSwitchId":        "vsw-1",
				},
			},
		},
		"TargetCapacitySpecification": map[string]interface{}{
			"TotalTargetCapacity":       "3",
			"DefaultTargetCapacityType": "Spot",
			"SpotTargetCapacity":        "2",
			"PayAsYouGoTargetCapacity":  "1",
		},
		"SpotOptions": map[string]interface{}{
			"AllocationStrategy":      "lowest-price",
			"InstancePoolsToUseCount": json.Number("3"),
		},
		"CandidateOptions": map[string]interface{}{
			"TimeoutMinutes": json.Number("60"),
		},
		"CapacitySpecification": map[string]interface{}{
			"PayAsYouGoCapacity": json.Number("2"),
			"PrePaidCapacity":    json.Number("0"),
			"SpotCapacity":       json.Number("3"),
			"TotalCapacity":      json.Number("5"),
		},
		"SuspendedProcesses": map[string]interface{}{
			"SuspendedProcess": []interface{}{"ScaleOut"},
		},
	})
	require.NoError(t, err)

	attributes := object.Attributes()
	assert.Equal(t, types.StringValue("apg-1"), attributes["auto_provisioning_group_id"])
	assert.Equal(t, types.BoolValue(true), attributes["terminate_instances"])

	configs, ok := attributes["launch_template_configs"].(types.List)
	if assert.True(t, ok) && assert.Len(t, configs.Elements(), 1) {
		configAttributes := configs.Elements()[0].(types.Object).Attributes()
		assert.Equal(t, types.StringValue("ecs.g6.large"), configAttributes["instance_type"])
		assert.Equal(t, types.StringValue("2"), configAttributes["weighted_capacity"])
		assert.Equal(t, types.StringValue("vsw-1"), configAttributes["vswitch_id"])
	}

	capacity, ok := attributes["target_capacity_specification"].(types.List)
	if assert.True(t, ok) && assert.Len(t, capacity.Elements(), 1) {
		capacityAttributes := capacity.Elements()[0].(types.Object).Attributes()
		assert.Equal(t, types.StringValue("3"), capacityAttributes["total_target_capacity"])
		assert.Equal(t, types.StringValue("Spot"), capacityAttributes["default_target_capacity_type"])
	}

	spotOptions, ok := attributes["spot_options"].(types.List)
	if assert.True(t, ok) && assert.Len(t, spotOptions.Elements(), 1) {
		spotAttributes := spotOptions.Elements()[0].(types.Object).Attributes()
		assert.Equal(t, types.StringValue("lowest-price"), spotAttributes["allocation_strategy"])
		assert.Equal(t, types.Int64Value(3), spotAttributes["instance_pools_to_use_count"])
	}

	candidateOptions, ok := attributes["candidate_options"].(types.List)
	if assert.True(t, ok) && assert.Len(t, candidateOptions.Elements(), 1) {
		candidateAttributes := candidateOptions.Elements()[0].(types.Object).Attributes()
		assert.Equal(t, types.Int64Value(60), candidateAttributes["timeout_minutes"])
	}

	capacitySpecification, ok := attributes["capacity_specification"].(types.List)
	if assert.True(t, ok) && assert.Len(t, capacitySpecification.Elements(), 1) {
		capacityAttributes := capacitySpecification.Elements()[0].(types.Object).Attributes()
		assert.Equal(t, types.StringValue("2"), capacityAttributes["pay_as_you_go_capacity"])
		assert.Equal(t, types.StringValue("0"), capacityAttributes["pre_paid_capacity"])
		assert.Equal(t, types.StringValue("3"), capacityAttributes["spot_capacity"])
		assert.Equal(t, types.StringValue("5"), capacityAttributes["total_capacity"])
	}

	suspendedProcesses, ok := attributes["suspended_processes"].(types.List)
	if assert.True(t, ok) && assert.Len(t, suspendedProcesses.Elements(), 1) {
		assert.Equal(t, types.StringValue("ScaleOut"), suspendedProcesses.Elements()[0])
	}
}

func TestUnitEcsReservedInstanceObject(t *testing.T) {
	object, err := reservedInstanceObject(map[string]interface{}{
		"ReservedInstanceId":   "ri-1",
		"ReservedInstanceName": "example",
		"InstanceAmount":       json.Number("2"),
		"Status":               "Active",
		"Scope":                "region",
		"OperationLocks": map[string]interface{}{
			"OperationLock": []interface{}{
				map[string]interface{}{"LockReason": "security"},
			},
		},
	})
	require.NoError(t, err)

	attributes := object.Attributes()
	assert.Equal(t, types.StringValue("ri-1"), attributes["reserved_instance_id"])
	assert.Equal(t, types.Int64Value(2), attributes["instance_amount"])
	assert.Equal(t, types.StringValue("Active"), attributes["status"])
	assert.Equal(t, types.StringValue("region"), attributes["scope"])

	operationLocks, ok := attributes["operation_locks"].(types.List)
	if assert.True(t, ok) && assert.Len(t, operationLocks.Elements(), 1) {
		lockAttributes := operationLocks.Elements()[0].(types.Object).Attributes()
		assert.Equal(t, types.StringValue("security"), lockAttributes["lock_reason"])
	}
}

func TestUnitEcsAutoProvisioningGroupsCollectPagesStopsAtTotalCount(t *testing.T) {
	var pages []int64
	objects, totalCount, err := fwCollectNumberedPages(context.Background(), 10, 1, 2, autoProvisioningGroupsAction, "AutoProvisioningGroupId", func(pageNumber int64) (fwNumberedPage, error) {
		pages = append(pages, pageNumber)
		switch pageNumber {
		case 1:
			return fwNumberedPage{items: []map[string]interface{}{{"AutoProvisioningGroupId": "apg-1"}, {"AutoProvisioningGroupId": "apg-2"}}, total: 3, totalKnown: true}, nil
		case 2:
			return fwNumberedPage{items: []map[string]interface{}{{"AutoProvisioningGroupId": "apg-3"}}, total: 3, totalKnown: true}, nil
		default:
			return fwNumberedPage{}, fmt.Errorf("unexpected page %d", pageNumber)
		}
	})

	assert.NoError(t, err)
	assert.Equal(t, []int64{1, 2}, pages)
	assert.Len(t, objects, 3)
	assert.Equal(t, int64(3), totalCount)
}

// A read that starts after the first page cannot see the earlier rows: a short
// final page is complete once it accounts for the readable suffix of the total.
func TestUnitEcsAutoProvisioningGroupsCollectPagesAcceptsReadableSuffix(t *testing.T) {
	objects, totalCount, err := fwCollectNumberedPages(context.Background(), 10, 2, 2, autoProvisioningGroupsAction, "AutoProvisioningGroupId", func(pageNumber int64) (fwNumberedPage, error) {
		if pageNumber != 2 {
			return fwNumberedPage{}, fmt.Errorf("unexpected page %d", pageNumber)
		}
		return fwNumberedPage{items: []map[string]interface{}{{"AutoProvisioningGroupId": "apg-3"}}, total: 3, totalKnown: true}, nil
	})

	require.NoError(t, err, "the last readable row must not be reported as a short page")
	assert.Len(t, objects, 1)
	assert.Equal(t, int64(3), totalCount)
}

// A start page beyond the reported total has no readable rows, so an empty page
// is a complete result rather than an error.
func TestUnitEcsAutoProvisioningGroupsCollectPagesAcceptsStartPageBeyondTotal(t *testing.T) {
	objects, totalCount, err := fwCollectNumberedPages(context.Background(), 10, 3, 2, autoProvisioningGroupsAction, "AutoProvisioningGroupId", func(pageNumber int64) (fwNumberedPage, error) {
		return fwNumberedPage{total: 3, totalKnown: true}, nil
	})

	require.NoError(t, err)
	assert.Empty(t, objects)
	assert.Equal(t, int64(3), totalCount)
}

// An explicit TotalCount of 0 is a real total: it must be reported as 0 and it
// must reject rows that contradict it, unlike a response that omits the total.
func TestUnitEcsReservedInstancesCollectPagesKeepsExplicitZeroTotal(t *testing.T) {
	objects, totalCount, err := fwCollectNumberedPages(context.Background(), 10, 1, 10, reservedInstancesAction, "ReservedInstanceId", func(pageNumber int64) (fwNumberedPage, error) {
		return fwNumberedPage{total: 0, totalKnown: true}, nil
	})

	require.NoError(t, err)
	assert.Empty(t, objects)
	assert.Equal(t, int64(0), totalCount)

	_, _, err = fwCollectNumberedPages(context.Background(), 10, 1, 10, reservedInstancesAction, "ReservedInstanceId", func(pageNumber int64) (fwNumberedPage, error) {
		return fwNumberedPage{items: []map[string]interface{}{{"ReservedInstanceId": "ri-1"}}, total: 0, totalKnown: true}, nil
	})
	assert.Error(t, err, "rows must not contradict an explicit TotalCount of 0")
}

func TestUnitEcsAutoProvisioningGroupsCollectPagesFailsWhenBudgetExhausted(t *testing.T) {
	objects, totalCount, err := fwCollectNumberedPages(context.Background(), 2, 1, 1, autoProvisioningGroupsAction, "AutoProvisioningGroupId", func(pageNumber int64) (fwNumberedPage, error) {
		return fwNumberedPage{items: []map[string]interface{}{{"AutoProvisioningGroupId": fmt.Sprintf("apg-%d", pageNumber)}}}, nil
	})

	assert.Error(t, err)
	assert.Nil(t, objects)
	assert.Equal(t, int64(0), totalCount)
}

func TestUnitEcsAutoProvisioningGroupsCollectPagesRejectsRepeatedPage(t *testing.T) {
	objects, _, err := fwCollectNumberedPages(context.Background(), 10, 1, 2, autoProvisioningGroupsAction, "AutoProvisioningGroupId", func(pageNumber int64) (fwNumberedPage, error) {
		return fwNumberedPage{items: []map[string]interface{}{{"AutoProvisioningGroupId": "apg-1"}, {"AutoProvisioningGroupId": "apg-2"}}}, nil
	})

	assert.Error(t, err, "an endpoint that ignores PageNumber must not produce duplicated rows")
	assert.Nil(t, objects)
}

func TestUnitEcsAutoProvisioningGroupsCollectPagesStopsOnCancelledContext(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	calls := 0
	objects, _, err := fwCollectNumberedPages(ctx, 10, 1, 2, autoProvisioningGroupsAction, "AutoProvisioningGroupId", func(pageNumber int64) (fwNumberedPage, error) {
		calls++
		return fwNumberedPage{}, nil
	})

	assert.Error(t, err)
	assert.Nil(t, objects)
	assert.Zero(t, calls, "a cancelled read must not issue another request")
}

func TestUnitEcsReservedInstancesCollectPagesStopsOnShortPage(t *testing.T) {
	objects, totalCount, err := fwCollectNumberedPages(context.Background(), 10, 1, 2, reservedInstancesAction, "ReservedInstanceId", func(pageNumber int64) (fwNumberedPage, error) {
		if pageNumber == 1 {
			return fwNumberedPage{items: []map[string]interface{}{{"ReservedInstanceId": "ri-1"}, {"ReservedInstanceId": "ri-2"}}}, nil
		}
		return fwNumberedPage{items: []map[string]interface{}{{"ReservedInstanceId": "ri-3"}}}, nil
	})

	assert.NoError(t, err)
	assert.Len(t, objects, 3)
	assert.Equal(t, int64(3), totalCount, "a missing server total falls back to the number of rows read")
}

// A short page that leaves rows the server total promised must never be
// accepted as a complete result.
func TestUnitEcsReservedInstancesCollectPagesRejectsShortPageWithOutstandingRows(t *testing.T) {
	objects, totalCount, err := fwCollectNumberedPages(context.Background(), 10, 1, 10, reservedInstancesAction, "ReservedInstanceId", func(pageNumber int64) (fwNumberedPage, error) {
		return fwNumberedPage{items: []map[string]interface{}{{"ReservedInstanceId": "ri-1"}}, total: 7, totalKnown: true}, nil
	})

	assert.Error(t, err)
	assert.Nil(t, objects)
	assert.Equal(t, int64(0), totalCount)
}

func TestUnitEcsReservedInstancesCollectPagesRejectsEmptyPageWithOutstandingRows(t *testing.T) {
	objects, _, err := fwCollectNumberedPages(context.Background(), 10, 1, 10, reservedInstancesAction, "ReservedInstanceId", func(pageNumber int64) (fwNumberedPage, error) {
		if pageNumber == 1 {
			return fwNumberedPage{items: []map[string]interface{}{{"ReservedInstanceId": "ri-1"}}, total: 2, totalKnown: true}, nil
		}
		return fwNumberedPage{total: 2, totalKnown: true}, nil
	})

	assert.Error(t, err)
	assert.Nil(t, objects)
}

func TestUnitEcsReservedInstancesCollectPagesRejectsInconsistentTotal(t *testing.T) {
	objects, _, err := fwCollectNumberedPages(context.Background(), 10, 1, 1, reservedInstancesAction, "ReservedInstanceId", func(pageNumber int64) (fwNumberedPage, error) {
		switch pageNumber {
		case 1:
			return fwNumberedPage{items: []map[string]interface{}{{"ReservedInstanceId": "ri-1"}}, total: 3, totalKnown: true}, nil
		default:
			return fwNumberedPage{items: []map[string]interface{}{{"ReservedInstanceId": "ri-2"}}, total: 4, totalKnown: true}, nil
		}
	})

	assert.Error(t, err, "a total that changes between pages must be rejected")
	assert.Nil(t, objects)
}

func TestUnitEcsReservedInstancesCollectPagesRejectsDuplicateIdentity(t *testing.T) {
	objects, _, err := fwCollectNumberedPages(context.Background(), 10, 1, 1, reservedInstancesAction, "ReservedInstanceId", func(pageNumber int64) (fwNumberedPage, error) {
		return fwNumberedPage{items: []map[string]interface{}{{"ReservedInstanceId": "ri-1"}}}, nil
	})

	assert.Error(t, err, "the same identity on two pages must not produce duplicate rows")
	assert.Nil(t, objects)
}

func TestUnitEcsReservedInstancesCollectPagesRejectsMissingIdentity(t *testing.T) {
	objects, _, err := fwCollectNumberedPages(context.Background(), 10, 1, 10, reservedInstancesAction, "ReservedInstanceId", func(pageNumber int64) (fwNumberedPage, error) {
		return fwNumberedPage{items: []map[string]interface{}{{"ReservedInstanceName": "no-id"}}, total: 1, totalKnown: true}, nil
	})

	assert.Error(t, err, "a record without an identity cannot be addressed")
	assert.Nil(t, objects)
}

// An identity of another type must not be formatted into a string: it cannot
// address the record, so it is a malformed response.
func TestUnitEcsReservedInstancesCollectPagesRejectsNonStringIdentity(t *testing.T) {
	objects, _, err := fwCollectNumberedPages(context.Background(), 10, 1, 10, reservedInstancesAction, "ReservedInstanceId", func(pageNumber int64) (fwNumberedPage, error) {
		return fwNumberedPage{items: []map[string]interface{}{{"ReservedInstanceId": json.Number("42")}}, total: 1, totalKnown: true}, nil
	})

	assert.Error(t, err)
	assert.Nil(t, objects)
}

func TestUnitEcsDataSourceRequestHelpers(t *testing.T) {
	request := map[string]interface{}{}
	fwApplyRepeatList(request, "AutoProvisioningGroupId", []string{"apg-1", "apg-2"})
	assert.Equal(t, "apg-1", request["AutoProvisioningGroupId.1"])
	assert.Equal(t, "apg-2", request["AutoProvisioningGroupId.2"])

	tagRequest := map[string]interface{}{}
	fwApplyTags(tagRequest, types.MapValueMust(types.StringType, map[string]attr.Value{
		"env": types.StringValue("test"),
	}))
	assert.Equal(t, "env", tagRequest["Tag.1.Key"])
	assert.Equal(t, "test", tagRequest["Tag.1.Value"])
}

// The request-to-state regressions below pin the exact request each data source
// builds from its configuration. They cover the cloud scenarios whose live test
// can only observe an empty result.
func TestUnitEcsAutoProvisioningGroupsRequest(t *testing.T) {
	model := autoProvisioningGroupsModel{
		Ids:                       types.ListValueMust(types.StringType, []attr.Value{types.StringValue("apg-1"), types.StringValue("apg-2")}),
		AutoProvisioningGroupName: types.StringValue("example"),
		ResourceGroupId:           types.StringValue("rg-1"),
		Status:                    types.ListValueMust(types.StringType, []attr.Value{types.StringValue("active")}),
		Types:                     types.ListValueMust(types.StringType, []attr.Value{types.StringValue("maintain")}),
		Tags: types.MapValueMust(types.StringType, map[string]attr.Value{
			"env": types.StringValue("test"),
		}),
	}

	request := (&autoProvisioningGroupsDataSource{}).autoProvisioningGroupsRequest(model, "cn-hangzhou", 2, 50)

	assert.Equal(t, "cn-hangzhou", request["RegionId"])
	assert.Equal(t, 2, request["PageNumber"])
	assert.Equal(t, 50, request["PageSize"])
	assert.Equal(t, "example", request["AutoProvisioningGroupName"])
	assert.Equal(t, "rg-1", request["ResourceGroupId"])
	assert.Equal(t, "apg-1", request["AutoProvisioningGroupId.1"])
	assert.Equal(t, "apg-2", request["AutoProvisioningGroupId.2"])
	assert.Equal(t, "active", request["AutoProvisioningGroupStatus.1"])
	assert.Equal(t, "maintain", request["AutoProvisioningGroupTypes.1"])
	assert.Equal(t, "env", request["Tag.1.Key"])
	assert.Equal(t, "test", request["Tag.1.Value"])
}

func TestUnitEcsReservedInstancesRequest(t *testing.T) {
	model := reservedInstancesModel{
		Ids:                  types.ListValueMust(types.StringType, []attr.Value{types.StringValue("ri-1")}),
		ReservedInstanceName: types.StringValue("example"),
		ZoneId:               types.StringValue("cn-hangzhou-k"),
		LockReason:           types.StringValue("financial"),
		InstanceType:         types.StringValue("ecs.g6.large"),
		InstanceTypeFamily:   types.StringValue("ecs.g6"),
		OfferingType:         types.StringValue("No Upfront"),
		AllocationType:       types.StringValue("Normal"),
		Status:               types.ListValueMust(types.StringType, []attr.Value{types.StringValue("Active")}),
		Tags: types.MapValueMust(types.StringType, map[string]attr.Value{
			"env": types.StringValue("test"),
		}),
	}

	request := reservedInstancesRequest(model, "cn-hangzhou", "Zone", 1, 100)

	assert.Equal(t, "cn-hangzhou", request["RegionId"])
	assert.Equal(t, 1, request["PageNumber"])
	assert.Equal(t, 100, request["PageSize"])
	assert.Equal(t, "cn-hangzhou-k", request["ZoneId"])
	assert.Equal(t, "Zone", request["Scope"], "zone_id without scope requests Scope = Zone")
	assert.Equal(t, "ri-1", request["ReservedInstanceId.1"])
	assert.Equal(t, "Active", request["Status.1"])
	assert.Equal(t, "env", request["Tag.1.Key"])
}

func TestUnitEcsImagePipelineExecutionsRequest(t *testing.T) {
	model := imagePipelineExecutionsModel{
		ImagePipelineId: types.StringValue("pipeline-1"),
		Status:          types.StringValue(" SUCCESS , FAILED "),
		MaxResults:      types.Int64Value(10),
		NextToken:       types.StringValue("token-1"),
		Tags: types.MapValueMust(types.StringType, map[string]attr.Value{
			"env": types.StringValue("test"),
		}),
	}

	request := imagePipelineExecutionsRequest(model, "cn-hangzhou")

	assert.Equal(t, "cn-hangzhou", request["RegionId"])
	assert.Equal(t, "pipeline-1", request["ImagePipelineId"])
	assert.Equal(t, "SUCCESS,FAILED", request["Status"])
	assert.Equal(t, 10, request["MaxResults"])
	assert.Equal(t, "token-1", request["NextToken"])
	assert.Equal(t, "env", request["Tag.1.Key"])

	defaults := imagePipelineExecutionsRequest(imagePipelineExecutionsModel{
		ImagePipelineId: types.StringValue("pipeline-1"),
		Status:          types.StringNull(),
		MaxResults:      types.Int64Null(),
		NextToken:       types.StringNull(),
		Tags:            types.MapNull(types.StringType),
	}, "cn-hangzhou")
	assert.Equal(t, imagePipelineExecutionsDefaultMaxResults, defaults["MaxResults"])
	assert.NotContains(t, defaults, "Status")
	assert.NotContains(t, defaults, "NextToken")
}

func TestUnitEcsAutoProvisioningGroupsIdsLimit(t *testing.T) {
	validate := func(count int) bool {
		elements := make([]attr.Value, 0, count)
		for i := 0; i < count; i++ {
			elements = append(elements, types.StringValue(fmt.Sprintf("apg-%d", i)))
		}
		var response validator.ListResponse
		listvalidator.SizeAtMost(autoProvisioningGroupsMaxIds).ValidateList(context.Background(), validator.ListRequest{
			ConfigValue: types.ListValueMust(types.StringType, elements),
			Path:        path.Root("ids"),
		}, &response)
		return response.Diagnostics.HasError()
	}

	assert.False(t, validate(autoProvisioningGroupsMaxIds), "20 IDs are within the documented limit")
	assert.True(t, validate(autoProvisioningGroupsMaxIds+1), "21 IDs exceed the documented limit")
}

func TestUnitFwDecodeCountRejectsInvalidValues(t *testing.T) {
	for _, value := range []interface{}{"1.5", json.Number("1.5"), -1, json.Number("-1"), 2.5, float32(2.5)} {
		_, _, err := fwDecodeCount(map[string]interface{}{"TotalCount": value}, "TotalCount", "Describe")
		assert.Error(t, err, "value %v must not be a valid count", value)
	}
	assert.Equal(t, 0, fwRemainingBudgetMillis(context.Background()), "a context without deadline leaves the SDK defaults in place")
}

// The transport receives the remaining read budget and the SDK retry loop is
// off; every retry stays under the caller's own deadline.
func TestUnitFwRpcPostBoundsTransportWithReadBudget(t *testing.T) {
	original := fwRpcCall
	defer func() { fwRpcCall = original }()

	var captured *utilV2.RuntimeOptions
	fwRpcCall = func(ctx context.Context, client *connectivity.AliyunClient, action string, request map[string]interface{}, runtime *utilV2.RuntimeOptions) (map[string]interface{}, error) {
		captured = runtime
		return map[string]interface{}{}, nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 90*time.Second)
	defer cancel()
	_, err := fwRpcPost(ctx, nil, "type", "Action", map[string]interface{}{})
	require.NoError(t, err)
	require.NotNil(t, captured)
	require.NotNil(t, captured.Autoretry)
	assert.False(t, *captured.Autoretry, "the SDK internal retries must be off")
	require.NotNil(t, captured.ReadTimeout)
	require.NotNil(t, captured.ConnectTimeout)
	assert.Greater(t, *captured.ReadTimeout, 0)
	assert.LessOrEqual(t, *captured.ReadTimeout, 90*1000)
	assert.Greater(t, *captured.ConnectTimeout, 0)
	assert.LessOrEqual(t, *captured.ConnectTimeout, 90*1000)
}

// The request must carry the read context, so cancelling the read aborts the
// call that is already in flight, and the cancelled response is never retried or
// returned to the caller.
func TestUnitFwRpcPostStopsWhenTheReadIsCancelledInFlight(t *testing.T) {
	original := fwRpcCall
	defer func() { fwRpcCall = original }()

	ctx, cancel := context.WithCancel(context.Background())
	calls := 0
	var requestCtx context.Context
	fwRpcCall = func(callCtx context.Context, client *connectivity.AliyunClient, action string, request map[string]interface{}, runtime *utilV2.RuntimeOptions) (map[string]interface{}, error) {
		calls++
		requestCtx = callCtx
		cancel()
		<-callCtx.Done()
		return nil, callCtx.Err()
	}

	_, err := fwRpcPost(ctx, nil, "type", "Action", map[string]interface{}{})
	assert.Error(t, err)
	assert.Equal(t, 1, calls, "a request cancelled in flight must not be retried")
	require.NotNil(t, requestCtx)
	assert.ErrorIs(t, requestCtx.Err(), context.Canceled, "the transport must receive the read context")
}

func TestUnitFwRpcPostRetriesThrottlingThenStops(t *testing.T) {
	originalCall, originalBackoff := fwRpcCall, fwRetryBackoff
	defer func() { fwRpcCall, fwRetryBackoff = originalCall, originalBackoff }()
	fwRetryBackoff = time.Millisecond

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	calls := 0
	fwRpcCall = func(callCtx context.Context, client *connectivity.AliyunClient, action string, request map[string]interface{}, runtime *utilV2.RuntimeOptions) (map[string]interface{}, error) {
		calls++
		if calls == 1 {
			return nil, &dara.SDKError{Code: dara.String("Rejected.Throttling"), Message: dara.String("throttled")}
		}
		return map[string]interface{}{"ok": true}, nil
	}

	response, err := fwRpcPost(ctx, nil, "type", "Action", map[string]interface{}{})
	require.NoError(t, err)
	assert.Equal(t, map[string]interface{}{"ok": true}, response)
	assert.Equal(t, 2, calls, "a throttled read is retried inside the read budget")
}

func TestUnitEcsAutoProvisioningGroupsObjectRejectsMalformedNestedMembers(t *testing.T) {
	_, err := autoProvisioningGroupObject(map[string]interface{}{
		"AutoProvisioningGroupId": "apg-1",
		"LaunchTemplateConfigs":   map[string]interface{}{"LaunchTemplateConfig": []interface{}{"not-an-object"}},
	})
	assert.Error(t, err, "a malformed launch template config must be reported")

	_, err = autoProvisioningGroupObject(map[string]interface{}{
		"AutoProvisioningGroupId": "apg-1",
		"SuspendedProcesses":      map[string]interface{}{"SuspendedProcess": []interface{}{1}},
	})
	assert.Error(t, err, "a malformed suspended process must be reported")

	_, err = autoProvisioningGroupObject(map[string]interface{}{
		"AutoProvisioningGroupName": "no-id",
	})
	assert.Error(t, err, "a record without an identity must be reported")
}

func TestUnitEcsReservedInstanceObjectRejectsMalformedNestedMembers(t *testing.T) {
	_, err := reservedInstanceObject(map[string]interface{}{
		"ReservedInstanceId": "ri-1",
		"OperationLocks":     map[string]interface{}{"OperationLock": []interface{}{"not-an-object"}},
	})
	assert.Error(t, err, "a malformed operation lock must be reported")

	_, err = reservedInstanceObject(map[string]interface{}{
		"ReservedInstanceName": "no-id",
	})
	assert.Error(t, err, "a record without an identity must be reported")
}

// The request-to-state tests below run the same load() the data source Read
// uses: request construction, every page, record decoding and the final state.
func TestUnitEcsAutoProvisioningGroupsLoadRequestToState(t *testing.T) {
	original := fwRpcCall
	defer func() { fwRpcCall = original }()

	requests := make([]map[string]interface{}, 0)
	fwRpcCall = func(ctx context.Context, client *connectivity.AliyunClient, action string, request map[string]interface{}, runtime *utilV2.RuntimeOptions) (map[string]interface{}, error) {
		snapshot := map[string]interface{}{}
		for key, value := range request {
			snapshot[key] = value
		}
		requests = append(requests, snapshot)
		page, _ := request["PageNumber"].(int)
		items := map[int]map[string]interface{}{
			1: {"AutoProvisioningGroupId": "apg-a", "AutoProvisioningGroupName": "group-a", "Status": "active"},
			2: {"AutoProvisioningGroupId": "apg-b", "AutoProvisioningGroupName": "group-b", "Status": "active"},
		}
		item, ok := items[page]
		if !ok {
			return nil, fmt.Errorf("unexpected page %d", page)
		}
		return map[string]interface{}{
			"TotalCount": json.Number("2"),
			"AutoProvisioningGroups": map[string]interface{}{
				"AutoProvisioningGroup": []interface{}{item},
			},
		}, nil
	}

	model := autoProvisioningGroupsModel{
		Ids:                       types.ListValueMust(types.StringType, []attr.Value{types.StringValue("apg-b"), types.StringValue("apg-b"), types.StringValue(" ")}),
		AutoProvisioningGroupName: types.StringValue("group-b"),
		Status:                    types.ListValueMust(types.StringType, []attr.Value{types.StringValue("active")}),
		PageSize:                  types.Int64Value(1),
	}

	state, err := (&autoProvisioningGroupsDataSource{}).load(context.Background(), nil, "cn-hangzhou", model)
	require.NoError(t, err)

	require.Len(t, requests, 2, "the target on the second page must be read")
	assert.Equal(t, "group-b", requests[0]["AutoProvisioningGroupName"])
	assert.Equal(t, "apg-b", requests[0]["AutoProvisioningGroupId.1"], "duplicate and blank IDs are sent once")
	assert.NotContains(t, requests[0], "AutoProvisioningGroupId.2")
	assert.Equal(t, "active", requests[0]["AutoProvisioningGroupStatus.1"])
	assert.Equal(t, 2, requests[1]["PageNumber"])

	require.Len(t, state.Ids.Elements(), 2)
	assert.Equal(t, types.StringValue("apg-a"), state.Ids.Elements()[0])
	assert.Equal(t, types.StringValue("apg-b"), state.Ids.Elements()[1])
	assert.Equal(t, int64(2), state.TotalCount.ValueInt64())
	assert.Len(t, state.Groups.Elements(), 2)
	assert.Equal(t, types.StringValue("group-b"), state.Names.Elements()[1])
}

func TestUnitEcsReservedInstancesLoadRequestToState(t *testing.T) {
	original := fwRpcCall
	defer func() { fwRpcCall = original }()

	var request map[string]interface{}
	fwRpcCall = func(ctx context.Context, client *connectivity.AliyunClient, action string, value map[string]interface{}, runtime *utilV2.RuntimeOptions) (map[string]interface{}, error) {
		request = map[string]interface{}{}
		for key, item := range value {
			request[key] = item
		}
		return map[string]interface{}{
			"TotalCount": json.Number("1"),
			"ReservedInstances": map[string]interface{}{
				"ReservedInstance": []interface{}{
					map[string]interface{}{
						"ReservedInstanceId":   "ri-target",
						"ReservedInstanceName": "reserved-b",
						"Status":               "Active",
						"InstanceAmount":       json.Number("2"),
					},
				},
			},
		}, nil
	}

	model := reservedInstancesModel{
		ZoneId: types.StringValue("cn-hangzhou-k"),
		Status: types.ListValueMust(types.StringType, []attr.Value{types.StringValue("Active")}),
	}

	state, err := (&reservedInstancesDataSource{}).load(context.Background(), nil, "cn-hangzhou", "Zone", 1, 100, model)
	require.NoError(t, err)
	require.NotNil(t, request)
	assert.Equal(t, "cn-hangzhou-k", request["ZoneId"])
	assert.Equal(t, "Zone", request["Scope"], "zone_id without scope requests Scope = Zone")
	assert.Equal(t, "Active", request["Status.1"])

	require.Len(t, state.Ids.Elements(), 1)
	assert.Equal(t, types.StringValue("ri-target"), state.Ids.Elements()[0])
	assert.Equal(t, int64(1), state.TotalCount.ValueInt64())
	require.Len(t, state.Instances.Elements(), 1)
	attributes := state.Instances.Elements()[0].(types.Object).Attributes()
	assert.Equal(t, types.Int64Value(2), attributes["instance_amount"])
}

// A read that starts on a later page returns the readable suffix of the total
// through the same load() path Read uses, including the zone inference.
func TestUnitEcsReservedInstancesLoadLaterPage(t *testing.T) {
	original := fwRpcCall
	defer func() { fwRpcCall = original }()

	requests := make([]map[string]interface{}, 0)
	fwRpcCall = func(ctx context.Context, client *connectivity.AliyunClient, action string, request map[string]interface{}, runtime *utilV2.RuntimeOptions) (map[string]interface{}, error) {
		snapshot := map[string]interface{}{}
		for key, value := range request {
			snapshot[key] = value
		}
		requests = append(requests, snapshot)
		if page, _ := request["PageNumber"].(int); page != 2 {
			return nil, fmt.Errorf("unexpected page %v", request["PageNumber"])
		}
		return map[string]interface{}{
			"TotalCount": json.Number("3"),
			"ReservedInstances": map[string]interface{}{
				"ReservedInstance": []interface{}{
					map[string]interface{}{"ReservedInstanceId": "ri-3", "Status": "Active"},
				},
			},
		}, nil
	}

	model := reservedInstancesModel{ZoneId: types.StringValue("cn-hangzhou-k")}
	state, err := (&reservedInstancesDataSource{}).load(context.Background(), nil, "cn-hangzhou", "Zone", 2, 2, model)
	require.NoError(t, err, "the last readable row of page 2 must not be reported as a short page")

	require.Len(t, requests, 1)
	assert.Equal(t, 2, requests[0]["PageNumber"])
	assert.Equal(t, "cn-hangzhou-k", requests[0]["ZoneId"])
	assert.Equal(t, "Zone", requests[0]["Scope"], "zone_id without scope requests Scope = Zone")
	assert.Len(t, state.Instances.Elements(), 1)
	assert.Equal(t, types.StringValue("ri-3"), state.Ids.Elements()[0])
	assert.Equal(t, int64(3), state.TotalCount.ValueInt64(), "the total still counts the rows of the skipped page")
}

// A configured next_token is the starting point of the first request, not a
// value that only appears afterwards.
func TestUnitEcsImagePipelineExecutionsLoadStartsFromConfiguredToken(t *testing.T) {
	original := fwRpcCall
	defer func() { fwRpcCall = original }()

	requests := make([]map[string]interface{}, 0)
	fwRpcCall = func(ctx context.Context, client *connectivity.AliyunClient, action string, request map[string]interface{}, runtime *utilV2.RuntimeOptions) (map[string]interface{}, error) {
		snapshot := map[string]interface{}{}
		for key, value := range request {
			snapshot[key] = value
		}
		requests = append(requests, snapshot)
		return map[string]interface{}{
			"ImagePipelineExecution": map[string]interface{}{
				"ImagePipelineExecutionSet": []interface{}{
					map[string]interface{}{"ExecutionId": "exec-from-token", "Status": "SUCCESS"},
				},
			},
		}, nil
	}

	model := imagePipelineExecutionsModel{
		ImagePipelineId: types.StringValue("pipeline-1"),
		NextToken:       types.StringValue("token-start"),
	}
	state, err := (&imagePipelineExecutionsDataSource{}).load(context.Background(), nil, "cn-hangzhou", model)
	require.NoError(t, err)

	require.Len(t, requests, 1)
	assert.Equal(t, "token-start", requests[0]["NextToken"], "the configured token must be sent with the first request")
	assert.Equal(t, types.StringValue("exec-from-token"), state.Ids.Elements()[0])
	assert.Equal(t, types.StringValue(""), state.NextToken, "the read ends when the server returns no token")
}

func TestUnitEcsImagePipelineExecutionsLoadRequestToState(t *testing.T) {
	original := fwRpcCall
	defer func() { fwRpcCall = original }()

	requests := make([]map[string]interface{}, 0)
	fwRpcCall = func(ctx context.Context, client *connectivity.AliyunClient, action string, request map[string]interface{}, runtime *utilV2.RuntimeOptions) (map[string]interface{}, error) {
		snapshot := map[string]interface{}{}
		for key, value := range request {
			snapshot[key] = value
		}
		requests = append(requests, snapshot)
		if _, ok := request["NextToken"]; !ok {
			return map[string]interface{}{
				"NextToken": "token-1",
				"ImagePipelineExecution": map[string]interface{}{
					"ImagePipelineExecutionSet": []interface{}{
						map[string]interface{}{"ExecutionId": "exec-first", "Status": "SUCCESS"},
					},
				},
			}, nil
		}
		return map[string]interface{}{
			"ImagePipelineExecution": map[string]interface{}{
				"ImagePipelineExecutionSet": []interface{}{
					map[string]interface{}{"ExecutionId": "exec-target", "Status": "SUCCESS"},
				},
			},
		}, nil
	}

	model := imagePipelineExecutionsModel{
		ImagePipelineId: types.StringValue("pipeline-1"),
		Status:          types.StringValue("SUCCESS, FAILED"),
	}

	state, err := (&imagePipelineExecutionsDataSource{}).load(context.Background(), nil, "cn-hangzhou", model)
	require.NoError(t, err)

	require.Len(t, requests, 2, "the continuation token must be followed")
	assert.Equal(t, "pipeline-1", requests[0]["ImagePipelineId"])
	assert.Equal(t, "SUCCESS,FAILED", requests[0]["Status"])
	assert.Equal(t, "token-1", requests[1]["NextToken"])

	require.Len(t, state.Ids.Elements(), 2)
	assert.Equal(t, types.StringValue("exec-first"), state.Ids.Elements()[0])
	assert.Equal(t, types.StringValue("exec-target"), state.Ids.Elements()[1])
	assert.Equal(t, int64(2), state.TotalCount.ValueInt64())
	assert.Equal(t, types.StringValue(""), state.NextToken)
	assert.Len(t, state.Executions.Elements(), 2)
}

func TestUnitEcsImagePipelineExecutionsTagsRejectNull(t *testing.T) {
	response := &datasource.SchemaResponse{}
	(&imagePipelineExecutionsDataSource{}).Schema(context.Background(), datasource.SchemaRequest{}, response)
	attribute, ok := response.Schema.Attributes["tags"].(schema.MapAttribute)
	require.True(t, ok)
	require.NotEmpty(t, attribute.Validators, "the executions tags attribute must validate its entries")

	var mapResponse validator.MapResponse
	for _, entry := range attribute.Validators {
		entry.ValidateMap(context.Background(), validator.MapRequest{
			ConfigValue: types.MapValueMust(types.StringType, map[string]attr.Value{"env": types.StringNull()}),
			Path:        path.Root("tags"),
		}, &mapResponse)
	}
	assert.True(t, mapResponse.Diagnostics.HasError(), "a null tag value must be rejected at the query entry")
}
