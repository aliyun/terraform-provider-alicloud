package ecs

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
)

func TestUnitEcsImagePipelineExecutionsCollectPagesFollowsTokens(t *testing.T) {
	var visited []string
	objects, lastToken, err := imagePipelineExecutionsCollectPages(10, "", func(token string) ([]map[string]interface{}, string, error) {
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

func TestUnitEcsImagePipelineExecutionsCollectPagesRejectsRepeatedToken(t *testing.T) {
	calls := 0
	objects, lastToken, err := imagePipelineExecutionsCollectPages(10, "", func(token string) ([]map[string]interface{}, string, error) {
		calls++
		return []map[string]interface{}{{"ExecutionId": fmt.Sprintf("exec-%d", calls)}}, "same-token", nil
	})

	assert.Error(t, err)
	assert.Nil(t, objects, "an incomplete result must not be returned")
	assert.Empty(t, lastToken)
	assert.Equal(t, 2, calls, "the repeated token must not be followed again")
}

func TestUnitEcsImagePipelineExecutionsCollectPagesFailsWhenBudgetExhausted(t *testing.T) {
	calls := 0
	objects, _, err := imagePipelineExecutionsCollectPages(2, "", func(token string) ([]map[string]interface{}, string, error) {
		calls++
		return []map[string]interface{}{{"ExecutionId": fmt.Sprintf("exec-%d", calls)}}, fmt.Sprintf("token-%d", calls), nil
	})

	assert.Error(t, err)
	assert.Nil(t, objects, "an incomplete result must not be returned")
	assert.Equal(t, 2, calls)
}

func TestUnitEcsImagePipelineExecutionsStatusValidator(t *testing.T) {
	v := imagePipelineExecutionsStatusValidator{}

	var ok validator.StringResponse
	v.ValidateString(context.Background(), validator.StringRequest{
		ConfigValue: types.StringValue("SUCCESS, FAILED"),
		Path:        path.Root("status"),
	}, &ok)
	assert.False(t, ok.Diagnostics.HasError())

	var bad validator.StringResponse
	v.ValidateString(context.Background(), validator.StringRequest{
		ConfigValue: types.StringValue("SUCCESS,not-a-status"),
		Path:        path.Root("status"),
	}, &bad)
	assert.True(t, bad.Diagnostics.HasError())
}

func TestUnitEcsImagePipelineExecutionsNormalizedStatus(t *testing.T) {
	assert.Equal(t, "SUCCESS,FAILED", imagePipelineExecutionsNormalizedStatus(types.StringValue(" SUCCESS , FAILED ")).ValueString())
	assert.True(t, imagePipelineExecutionsNormalizedStatus(types.StringNull()).IsNull())
}

func TestUnitEcsImagePipelineExecutionObject(t *testing.T) {
	object := imagePipelineExecutionObject(map[string]interface{}{
		"ExecutionId":     "exec-1",
		"ImagePipelineId": "pipeline-1",
		"Status":          "SUCCESS",
		"CreationTime":    "2026-11-05T10:00:00Z",
		"Tags": map[string]interface{}{
			"Tag": []interface{}{
				map[string]interface{}{"Key": "env", "Value": "test"},
			},
		},
	})

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
	object := autoProvisioningGroupObject(map[string]interface{}{
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
			"InstancePoolsToUseCount": 3,
		},
	})

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
}

func TestUnitEcsReservedInstanceObject(t *testing.T) {
	object := reservedInstanceObject(map[string]interface{}{
		"ReservedInstanceId":   "ri-1",
		"ReservedInstanceName": "example",
		"InstanceAmount":       2,
		"Status":               "Active",
		"Scope":                "Region",
	})

	attributes := object.Attributes()
	assert.Equal(t, types.StringValue("ri-1"), attributes["reserved_instance_id"])
	assert.Equal(t, types.Int64Value(2), attributes["instance_amount"])
	assert.Equal(t, types.StringValue("Active"), attributes["status"])
	assert.Equal(t, types.StringValue("Region"), attributes["scope"])
}

func TestUnitEcsAutoProvisioningGroupsCollectPagesStopsAtTotalCount(t *testing.T) {
	var pages []int64
	objects, totalCount, err := ecsFwCollectNumberedPages(10, 1, 2, autoProvisioningGroupsAction, func(pageNumber int64) ([]map[string]interface{}, int64, error) {
		pages = append(pages, pageNumber)
		switch pageNumber {
		case 1:
			return []map[string]interface{}{{}, {}}, 3, nil
		case 2:
			return []map[string]interface{}{{}}, 3, nil
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
	objects, totalCount, err := ecsFwCollectNumberedPages(2, 1, 1, autoProvisioningGroupsAction, func(pageNumber int64) ([]map[string]interface{}, int64, error) {
		return []map[string]interface{}{{}}, 0, nil
	})

	assert.Error(t, err)
	assert.Nil(t, objects)
	assert.Equal(t, int64(0), totalCount)
}

func TestUnitEcsReservedInstancesCollectPagesStopsOnShortPage(t *testing.T) {
	objects, totalCount, err := ecsFwCollectNumberedPages(10, 1, 2, reservedInstancesAction, func(pageNumber int64) ([]map[string]interface{}, int64, error) {
		if pageNumber == 1 {
			return []map[string]interface{}{{}, {}}, 0, nil
		}
		return []map[string]interface{}{{}}, 0, nil
	})

	assert.NoError(t, err)
	assert.Len(t, objects, 3)
	assert.Equal(t, int64(3), totalCount, "a missing server total falls back to the number of rows read")
}

func TestUnitEcsReservedInstancesCollectPagesKeepsServerTotal(t *testing.T) {
	_, totalCount, err := ecsFwCollectNumberedPages(10, 1, 10, reservedInstancesAction, func(pageNumber int64) ([]map[string]interface{}, int64, error) {
		return []map[string]interface{}{{}}, 7, nil
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
