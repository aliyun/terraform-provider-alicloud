package rocketmq

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Settle/retry tuning. The real ACC flow tolerates a 20-minute create settle
// (the DisasterRecovery backend 404s a freshly-referenced instance until its
// ledger settles); the no-cloud unit mock returns the steady state on the
// first call, so these timeouts never bind in tests.
const (
	createMaxAttempts   = 20
	createSettleTimeout = 20 * time.Minute
	deleteSettleTimeout = 5 * time.Minute
)

// writeExclusiveFields are sent on Create/Update but never echoed by
// GetDisasterRecoveryPlan. Read carries them over from prior state so the
// plan does not drift on every refresh.
var writeExclusiveFields = []string{"password", "endpoint_url", "network_type", "vpc_id", "vswitch_id"}

// incrementalWait returns a closure that sleeps a growing, bounded backoff,
// mirroring the SDKv2 helper used by the legacy RocketmqServiceV2 path.
func incrementalWait(min, max time.Duration) func() {
	step := min
	return func() {
		time.Sleep(step)
		step *= 2
		if step > max {
			step = max
		}
	}
}

// strFrom maps a JSON-decoded string field to a Framework string value. Empty
// and missing values read back as null, matching the SDKv2 d.Set behavior.
func strFrom(v interface{}) types.String {
	switch s := v.(type) {
	case string:
		if s == "" {
			return types.StringNull()
		}
		return types.StringValue(s)
	case nil:
		return types.StringNull()
	default:
		str := fmt.Sprint(v)
		if str == "" {
			return types.StringNull()
		}
		return types.StringValue(str)
	}
}

// boolFrom maps a JSON-decoded bool field to a Framework bool value.
func boolFrom(v interface{}) types.Bool {
	switch b := v.(type) {
	case bool:
		return types.BoolValue(b)
	case nil:
		return types.BoolNull()
	default:
		return types.BoolNull()
	}
}

// messagePropertyElementType is the Framework object type of a single
// message_property block, matching the resource schema.
func messagePropertyElementType() attr.Type {
	return types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"property_key":   types.StringType,
			"property_value": types.StringType,
		},
	}
}

// instancesElementType is the Framework object type of a single instances
// block element, matching the resource schema. types.ListValueFrom needs it
// to rebuild a list from a []instanceModel slice.
func instancesElementType() attr.Type {
	return types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"auth_type":         types.StringType,
			"consumer_group_id": types.StringType,
			"endpoint_url":      types.StringType,
			"instance_id":       types.StringType,
			"instance_role":     types.StringType,
			"instance_type":     types.StringType,
			"message_property":  types.ListType{ElemType: messagePropertyElementType()},
			"network_type":      types.StringType,
			"password":          types.StringType,
			"region_id":         types.StringType,
			"security_group_id": types.StringType,
			"username":          types.StringType,
			"vpc_id":            types.StringType,
			"vswitch_id":        types.StringType,
		},
	}
}

// flattenInstances maps the API instances array onto resource models, carrying
// over write-exclusive fields from the prior state by index. The index match
// mirrors the SDKv2 behavior: positions correspond across read and the
// previously-stored instances list.
func flattenInstances(raw interface{}, prior []instanceModel) []instanceModel {
	out := make([]instanceModel, 0)
	rawList, ok := raw.([]interface{})
	if !ok {
		return out
	}
	for i, item := range rawList {
		obj, ok := item.(map[string]interface{})
		if !ok {
			continue
		}
		m := instanceModel{
			AuthType:        strFrom(obj["authType"]),
			ConsumerGroupId: strFrom(obj["consumerGroupId"]),
			EndpointUrl:     strFrom(obj["endpointUrl"]),
			InstanceId:      strFrom(obj["instanceId"]),
			InstanceRole:    strFrom(obj["instanceRole"]),
			InstanceType:    strFrom(obj["instanceType"]),
			NetworkType:     strFrom(obj["networkType"]),
			Password:        strFrom(obj["password"]),
			RegionId:        strFrom(obj["regionId"]),
			SecurityGroupId: strFrom(obj["securityGroupId"]),
			Username:        strFrom(obj["username"]),
			VpcId:           strFrom(obj["vpcId"]),
			VswitchId:       strFrom(obj["vSwitchId"]),
			MessageProperty: types.ListNull(messagePropertyElementType()),
		}
		if mp := flattenMessageProperty(obj["messageProperty"]); mp != nil {
			if list, diags := types.ListValueFrom(context.Background(), messagePropertyElementType(), mp); !diags.HasError() {
				m.MessageProperty = list
			}
		}
		if i < len(prior) {
			carryOverInto(&m, prior[i])
		}
		out = append(out, m)
	}
	return out
}

func flattenMessageProperty(raw interface{}) []messagePropertyModel {
	obj, ok := raw.(map[string]interface{})
	if !ok || len(obj) == 0 {
		return nil
	}
	return []messagePropertyModel{{
		PropertyKey:   strFrom(obj["propertyKey"]),
		PropertyValue: strFrom(obj["propertyValue"]),
	}}
}

// carryOverInto restores write-exclusive fields the API omitted from the prior
// state element, so a refresh does not drop password/endpoint/network/vpc.
func carryOverInto(dst *instanceModel, src instanceModel) {
	for _, field := range writeExclusiveFields {
		switch field {
		case "password":
			if dst.Password.IsNull() {
				dst.Password = src.Password
			}
		case "endpoint_url":
			if dst.EndpointUrl.IsNull() {
				dst.EndpointUrl = src.EndpointUrl
			}
		case "network_type":
			if dst.NetworkType.IsNull() {
				dst.NetworkType = src.NetworkType
			}
		case "vpc_id":
			if dst.VpcId.IsNull() {
				dst.VpcId = src.VpcId
			}
		case "vswitch_id":
			if dst.VswitchId.IsNull() {
				dst.VswitchId = src.VswitchId
			}
		}
	}
}
