package ens

import (
	"fmt"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
)

const (
	ensProduct  = "Ens"
	ensVersion  = "2017-11-10"
	ensPageSize = 50 // matches the provider's PageSizeLarge used by the SDKv2 ENS data sources
)

// describeNetworkInterface queries DescribeNetworkInterfaces for a single ENI
// by id and returns the matching element, or (nil, nil) when it no longer
// exists. The ENS Describe API returns a list even for a single id, so an
// empty list is the not-found signal rather than an API error.
func describeNetworkInterface(client *connectivity.AliyunClient, id string) (map[string]interface{}, error) {
	response, err := client.RpcPost(ensProduct, ensVersion, "DescribeNetworkInterfaces", nil, map[string]interface{}{
		"NetworkInterfaceId": id,
	}, true)
	if err != nil {
		if errs.NotFoundError(err) {
			return nil, nil
		}
		return nil, errs.WrapErrorf(err, "DescribeNetworkInterfaces: %s", id)
	}
	item, err := firstNetworkInterface(response)
	if err != nil {
		if err == errEmptyInterfaceSet {
			return nil, nil
		}
		return nil, err
	}
	return item, nil
}

var errEmptyInterfaceSet = fmt.Errorf("the response carried no NetworkInterfaceSets.NetworkInterfaceSet")

// firstNetworkInterface returns the first element of
// NetworkInterfaceSets.NetworkInterfaceSet, or errEmptyInterfaceSet when the
// list is absent or empty. ENS never omits the wrapper object for a successful
// call, so a missing list is treated as not-found by the caller.
func firstNetworkInterface(response map[string]interface{}) (map[string]interface{}, error) {
	sets, ok := response["NetworkInterfaceSets"].(map[string]interface{})
	if !ok {
		return nil, errEmptyInterfaceSet
	}
	list, ok := sets["NetworkInterfaceSet"].([]interface{})
	if !ok || len(list) == 0 {
		return nil, errEmptyInterfaceSet
	}
	item, ok := list[0].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("the NetworkInterfaceSet element is not an object: %v", list[0])
	}
	return item, nil
}

// networkInterfaceList returns the full NetworkInterfaceSet list from a
// DescribeNetworkInterfaces response, or an empty slice when absent.
func networkInterfaceList(response map[string]interface{}) []interface{} {
	sets, ok := response["NetworkInterfaceSets"].(map[string]interface{})
	if !ok {
		return nil
	}
	list, ok := sets["NetworkInterfaceSet"].([]interface{})
	if !ok {
		return nil
	}
	return list
}

// securityGroupIDs reads SecurityGroupIds.SecurityGroup from an ENI object and
// returns the IDs as a slice of strings. Missing or non-list values yield an
// empty slice, matching how the SDKv2 resource read the field.
func securityGroupIDs(object map[string]interface{}) []string {
	wrapper, ok := object["SecurityGroupIds"].(map[string]interface{})
	if !ok {
		return nil
	}
	raw, ok := wrapper["SecurityGroup"].([]interface{})
	if !ok {
		return nil
	}
	ids := make([]string, 0, len(raw))
	for _, v := range raw {
		if s, ok := v.(string); ok && s != "" {
			ids = append(ids, s)
		}
	}
	return ids
}
