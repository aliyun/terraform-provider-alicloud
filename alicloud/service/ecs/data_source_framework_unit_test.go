package ecs

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUnitEcsFwInt64DecodesJSONNumber(t *testing.T) {
	var decoded map[string]interface{}
	decoder := json.NewDecoder(strings.NewReader(`{"InstancePoolsToUseCount":3,"InstanceAmount":2,"TotalCount":7,"Fractional":4.0}`))
	decoder.UseNumber()
	assert.NoError(t, decoder.Decode(&decoded))

	assert.Equal(t, types.Int64Value(3), ecsFwInt64(decoded["InstancePoolsToUseCount"]))
	assert.Equal(t, types.Int64Value(2), ecsFwInt64(decoded["InstanceAmount"]))
	assert.Equal(t, types.Int64Value(7), ecsFwInt64(decoded["TotalCount"]))
	assert.Equal(t, types.Int64Value(4), ecsFwInt64(decoded["Fractional"]))
}

func TestUnitEcsFwTagMapValidator(t *testing.T) {
	assertTagMap := func(v ecsFwTagMapValidator, tags map[string]attr.Value, wantError bool) {
		var response validator.MapResponse
		v.ValidateMap(context.Background(), validator.MapRequest{
			ConfigValue: types.MapValueMust(types.StringType, tags),
			Path:        path.Root("tags"),
		}, &response)
		assert.Equal(t, wantError, response.Diagnostics.HasError(), "tags: %v", tags)
	}

	apg := ecsFwTagMapValidator{allowEmptyValue: true}
	assertTagMap(apg, map[string]attr.Value{"env": types.StringValue("")}, false)
	assertTagMap(apg, map[string]attr.Value{"env": types.StringValue("test")}, false)
	assertTagMap(apg, map[string]attr.Value{"": types.StringValue("test")}, true)
	assertTagMap(apg, map[string]attr.Value{"aliyun-tag": types.StringValue("test")}, true)
	assertTagMap(apg, map[string]attr.Value{"acs:tag": types.StringValue("test")}, true)
	assertTagMap(apg, map[string]attr.Value{"env": types.StringValue("https://example.com")}, true)
	assertTagMap(apg, map[string]attr.Value{strings.Repeat("k", 129): types.StringValue("test")}, true)
	assertTagMap(apg, map[string]attr.Value{"env": types.StringValue(strings.Repeat("v", 129))}, true)

	reserved := ecsFwTagMapValidator{valueForbiddenPrefixes: []string{"acs:"}}
	assertTagMap(reserved, map[string]attr.Value{"env": types.StringValue("")}, true)
	assertTagMap(reserved, map[string]attr.Value{"env": types.StringValue("acs:reserved")}, true)
	assertTagMap(reserved, map[string]attr.Value{"env": types.StringValue("http://example.com")}, true)
	assertTagMap(reserved, map[string]attr.Value{"env": types.StringValue("test")}, false)
}

// A null value must not skip the key rules, and it must not be sent as an empty
// string: null, unknown and an explicit empty string are three different
// requests.
func TestUnitEcsFwTagMapValidatorNullValues(t *testing.T) {
	validate := func(v ecsFwTagMapValidator, tags map[string]attr.Value) (bool, string) {
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

	apg := ecsFwTagMapValidator{allowEmptyValue: true}
	hasError, _ := validate(apg, map[string]attr.Value{"env": types.StringNull()})
	assert.True(t, hasError, "a null value must be rejected")
	hasError, summary := validate(apg, map[string]attr.Value{"acs:env": types.StringNull()})
	assert.True(t, hasError, "the key rules apply even when the value is null")
	assert.Contains(t, summary, "Invalid Tag Key")
	hasError, _ = validate(apg, map[string]attr.Value{"env": types.StringUnknown()})
	assert.False(t, hasError, "an unknown value cannot be validated yet")
	hasError, _ = validate(apg, map[string]attr.Value{"env": types.StringValue("")})
	assert.False(t, hasError, "an explicit empty value stays allowed for this API")

	reserved := ecsFwTagMapValidator{valueForbiddenPrefixes: []string{"acs:"}}
	hasError, _ = validate(reserved, map[string]attr.Value{"env": types.StringNull()})
	assert.True(t, hasError)
	hasError, _ = validate(reserved, map[string]attr.Value{"env": types.StringValue("")})
	assert.True(t, hasError)
}

func TestUnitEcsFwApplyTagsSerialization(t *testing.T) {
	request := map[string]interface{}{}
	ecsFwApplyTags(request, types.MapValueMust(types.StringType, map[string]attr.Value{
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

func TestUnitEcsFwDecodeRecordList(t *testing.T) {
	items, err := ecsFwDecodeRecordList(map[string]interface{}{}, "AutoProvisioningGroups", "AutoProvisioningGroup", "Describe")
	assert.NoError(t, err, "a missing member is a valid empty result")
	assert.Empty(t, items)

	items, err = ecsFwDecodeRecordList(map[string]interface{}{
		"AutoProvisioningGroups": map[string]interface{}{
			"AutoProvisioningGroup": []interface{}{
				map[string]interface{}{"AutoProvisioningGroupId": "apg-1"},
			},
		},
	}, "AutoProvisioningGroups", "AutoProvisioningGroup", "Describe")
	require.NoError(t, err)
	assert.Len(t, items, 1)

	_, err = ecsFwDecodeRecordList(map[string]interface{}{
		"AutoProvisioningGroups": "not-an-object",
	}, "AutoProvisioningGroups", "AutoProvisioningGroup", "Describe")
	assert.Error(t, err, "a malformed container must not be treated as empty")

	_, err = ecsFwDecodeRecordList(map[string]interface{}{
		"AutoProvisioningGroups": map[string]interface{}{"AutoProvisioningGroup": "not-a-list"},
	}, "AutoProvisioningGroups", "AutoProvisioningGroup", "Describe")
	assert.Error(t, err, "a malformed list must not be treated as empty")

	_, err = ecsFwDecodeRecordList(map[string]interface{}{
		"AutoProvisioningGroups": map[string]interface{}{
			"AutoProvisioningGroup": []interface{}{"not-an-object"},
		},
	}, "AutoProvisioningGroups", "AutoProvisioningGroup", "Describe")
	assert.Error(t, err, "a malformed record must not be dropped")
}

func TestUnitEcsFwDecodeCountAndToken(t *testing.T) {
	total, known, err := ecsFwDecodeCount(map[string]interface{}{}, "TotalCount", "Describe")
	assert.NoError(t, err)
	assert.False(t, known)
	assert.Zero(t, total)

	total, known, err = ecsFwDecodeCount(map[string]interface{}{"TotalCount": json.Number("7")}, "TotalCount", "Describe")
	assert.NoError(t, err)
	assert.True(t, known)
	assert.Equal(t, int64(7), total)

	_, _, err = ecsFwDecodeCount(map[string]interface{}{"TotalCount": "many"}, "TotalCount", "Describe")
	assert.Error(t, err, "an invalid count must not become 0")

	token, err := ecsFwDecodeToken(map[string]interface{}{}, "NextToken", "Describe")
	assert.NoError(t, err)
	assert.Empty(t, token)

	_, err = ecsFwDecodeToken(map[string]interface{}{"NextToken": 42}, "NextToken", "Describe")
	assert.Error(t, err, "an invalid token must not be treated as the last page")
}

func TestUnitEcsFwDecodeTags(t *testing.T) {
	tags, err := ecsFwDecodeTags(nil)
	require.NoError(t, err)
	assert.Empty(t, tags.Elements(), "an absent member is an empty map")

	tags, err = ecsFwDecodeTags(map[string]interface{}{
		"Tag": []interface{}{
			map[string]interface{}{"TagKey": "env", "TagValue": "test"},
			map[string]interface{}{"TagKey": "aliyun-reserved", "TagValue": "skip"},
		},
	})
	require.NoError(t, err)
	assert.Equal(t, types.StringValue("test"), tags.Elements()["env"])
	assert.NotContains(t, tags.Elements(), "aliyun-reserved")

	_, err = ecsFwDecodeTags(map[string]interface{}{"Tag": "not-a-list"})
	assert.Error(t, err)

	_, err = ecsFwDecodeTags(map[string]interface{}{"Tag": []interface{}{"not-an-object"}})
	assert.Error(t, err, "a malformed tag entry must be reported instead of panicking")

	_, err = ecsFwDecodeTags(map[string]interface{}{"Tag": []interface{}{map[string]interface{}{"TagValue": "test"}}})
	assert.Error(t, err, "a tag without a key is malformed")

	_, err = ecsFwDecodeTags(map[string]interface{}{"Tag": []interface{}{map[string]interface{}{"TagKey": 1, "TagValue": "test"}}})
	assert.Error(t, err, "a non-string tag key is malformed")
}

func TestUnitEcsFwCollectTokenPagesFollowsTokens(t *testing.T) {
	var visited []string
	objects, lastToken, err := ecsFwCollectTokenPages(context.Background(), 10, "", "Discover", func(token string) ([]map[string]interface{}, string, error) {
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

func TestUnitEcsFwCollectTokenPagesRejectsRepeatedToken(t *testing.T) {
	calls := 0
	objects, lastToken, err := ecsFwCollectTokenPages(context.Background(), 10, "", "Discover", func(token string) ([]map[string]interface{}, string, error) {
		calls++
		return []map[string]interface{}{{"ExecutionId": fmt.Sprintf("exec-%d", calls)}}, "same-token", nil
	})

	assert.Error(t, err)
	assert.Nil(t, objects, "an incomplete result must not be returned")
	assert.Empty(t, lastToken)
	assert.Equal(t, 2, calls, "the repeated token must not be followed again")
}

func TestUnitEcsFwCollectTokenPagesFailsWhenBudgetExhausted(t *testing.T) {
	calls := 0
	objects, _, err := ecsFwCollectTokenPages(context.Background(), 2, "", "Discover", func(token string) ([]map[string]interface{}, string, error) {
		calls++
		return []map[string]interface{}{{"ExecutionId": fmt.Sprintf("exec-%d", calls)}}, fmt.Sprintf("token-%d", calls), nil
	})

	assert.Error(t, err)
	assert.Nil(t, objects, "an incomplete result must not be returned")
	assert.Equal(t, 2, calls)
}

func TestUnitEcsFwCollectTokenPagesStopsOnCancelledContext(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	calls := 0
	objects, _, err := ecsFwCollectTokenPages(ctx, 10, "", "Discover", func(token string) ([]map[string]interface{}, string, error) {
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
	objects, totalCount, err := ecsFwCollectNumberedPages(context.Background(), 10, 1, 2, autoProvisioningGroupsAction, func(pageNumber int64) ([]map[string]interface{}, int64, error) {
		pages = append(pages, pageNumber)
		switch pageNumber {
		case 1:
			return []map[string]interface{}{{"AutoProvisioningGroupId": "apg-1"}, {"AutoProvisioningGroupId": "apg-2"}}, 3, nil
		case 2:
			return []map[string]interface{}{{"AutoProvisioningGroupId": "apg-3"}}, 3, nil
		default:
			return nil, 0, fmt.Errorf("unexpected page %d", pageNumber)
		}
	})

	assert.NoError(t, err)
	assert.Equal(t, []int64{1, 2}, pages)
	assert.Len(t, objects, 3)
	assert.Equal(t, int64(3), totalCount)
}

func TestUnitEcsAutoProvisioningGroupsCollectPagesFailsWhenBudgetExhausted(t *testing.T) {
	objects, totalCount, err := ecsFwCollectNumberedPages(context.Background(), 2, 1, 1, autoProvisioningGroupsAction, func(pageNumber int64) ([]map[string]interface{}, int64, error) {
		return []map[string]interface{}{{"AutoProvisioningGroupId": fmt.Sprintf("apg-%d", pageNumber)}}, 0, nil
	})

	assert.Error(t, err)
	assert.Nil(t, objects)
	assert.Equal(t, int64(0), totalCount)
}

func TestUnitEcsAutoProvisioningGroupsCollectPagesRejectsRepeatedPage(t *testing.T) {
	objects, _, err := ecsFwCollectNumberedPages(context.Background(), 10, 1, 2, autoProvisioningGroupsAction, func(pageNumber int64) ([]map[string]interface{}, int64, error) {
		return []map[string]interface{}{{"AutoProvisioningGroupId": "apg-1"}, {"AutoProvisioningGroupId": "apg-2"}}, 0, nil
	})

	assert.Error(t, err, "an endpoint that ignores PageNumber must not produce duplicated rows")
	assert.Nil(t, objects)
}

func TestUnitEcsAutoProvisioningGroupsCollectPagesStopsOnCancelledContext(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	calls := 0
	objects, _, err := ecsFwCollectNumberedPages(ctx, 10, 1, 2, autoProvisioningGroupsAction, func(pageNumber int64) ([]map[string]interface{}, int64, error) {
		calls++
		return nil, 0, nil
	})

	assert.Error(t, err)
	assert.Nil(t, objects)
	assert.Zero(t, calls, "a cancelled read must not issue another request")
}

func TestUnitEcsReservedInstancesCollectPagesStopsOnShortPage(t *testing.T) {
	objects, totalCount, err := ecsFwCollectNumberedPages(context.Background(), 10, 1, 2, reservedInstancesAction, func(pageNumber int64) ([]map[string]interface{}, int64, error) {
		if pageNumber == 1 {
			return []map[string]interface{}{{"ReservedInstanceId": "ri-1"}, {"ReservedInstanceId": "ri-2"}}, 0, nil
		}
		return []map[string]interface{}{{"ReservedInstanceId": "ri-3"}}, 0, nil
	})

	assert.NoError(t, err)
	assert.Len(t, objects, 3)
	assert.Equal(t, int64(3), totalCount, "a missing server total falls back to the number of rows read")
}

func TestUnitEcsReservedInstancesCollectPagesKeepsServerTotal(t *testing.T) {
	_, totalCount, err := ecsFwCollectNumberedPages(context.Background(), 10, 1, 10, reservedInstancesAction, func(pageNumber int64) ([]map[string]interface{}, int64, error) {
		return []map[string]interface{}{{"ReservedInstanceId": "ri-1"}}, 7, nil
	})

	assert.NoError(t, err)
	assert.Equal(t, int64(7), totalCount)
}

func TestUnitEcsDataSourceRequestHelpers(t *testing.T) {
	request := map[string]interface{}{}
	ecsFwApplyRepeatList(request, "AutoProvisioningGroupId", []string{"apg-1", "apg-2"})
	assert.Equal(t, "apg-1", request["AutoProvisioningGroupId.1"])
	assert.Equal(t, "apg-2", request["AutoProvisioningGroupId.2"])

	tagRequest := map[string]interface{}{}
	ecsFwApplyTags(tagRequest, types.MapValueMust(types.StringType, map[string]attr.Value{
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
