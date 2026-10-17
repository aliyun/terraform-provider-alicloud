package ecs

import (
	"fmt"
	"log"
	"regexp"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
)

// tagFilters are the Alibaba Cloud system-tag prefixes that the provider never
// owns; they are filtered out on both read (so the plan stays clean) and on
// the UntagResources side (so a delete never tries to remove them).
var tagFilters = []string{"^aliyun", "^acs:", "^http://", "^https://", "^sae.do.not.delete"}

// tagIgnored reports whether a tag key is a system tag the provider should
// ignore. It is reused verbatim from the root package's tagIgnored.
func tagIgnored(tagKey string) bool {
	for _, pattern := range tagFilters {
		if ok, _ := regexp.MatchString(pattern, tagKey); ok {
			log.Printf("[DEBUG] Found Alibaba Cloud specific tag with key: %s, ignoring.\n", tagKey)
			return true
		}
	}
	return false
}

// tagsToMap converts the Tags field of a Describe response into a
// map[string]string, dropping system tags. ECS returns tags as a list of
// {TagKey, TagValue} pairs, but the Key/Value shape is also accepted to stay
// aligned with the root helper.
func tagsToMap(tags interface{}) map[string]string {
	result := make(map[string]string)
	if tags == nil {
		return result
	}
	switch v := tags.(type) {
	case []interface{}:
		for _, raw := range v {
			t, ok := raw.(map[string]interface{})
			if !ok {
				continue
			}
			key, value := extractTagPair(t)
			if key == "" || tagIgnored(key) {
				continue
			}
			result[key] = value
		}
	case []map[string]interface{}:
		for _, t := range v {
			key, value := extractTagPair(t)
			if key == "" || tagIgnored(key) {
				continue
			}
			result[key] = value
		}
	case map[string]interface{}:
		for key, value := range v {
			if tagIgnored(key) {
				continue
			}
			if s, ok := value.(string); ok {
				result[key] = s
			}
		}
	default:
		log.Printf("[ERROR] Unknown tags type %T. The tags value is: %v.", v, tags)
	}
	return result
}

// extractTagPair reads a tag pair from either the {TagKey, TagValue} or the
// {Key, Value} shape into plain strings.
func extractTagPair(t map[string]interface{}) (string, string) {
	if v, ok := t["TagKey"].(string); ok {
		return v, fmt.Sprint(t["TagValue"])
	}
	if v, ok := t["Key"].(string); ok {
		return v, fmt.Sprint(t["Value"])
	}
	return "", ""
}

// setResourceTags reconciles the old and new tag maps for a PortRangeList by
// first UntagResources the removed/changed keys, then TagResources every key
// in the new map. Mirrors the SDKv2 EcsServiceV2.SetResourceTags flow and
// uses the portrangelist ResourceType for both calls.
func setResourceTags(client *connectivity.AliyunClient, id string, oldTags, newTags map[string]string) error {
	removed := make([]string, 0)
	for key, value := range oldTags {
		if nv, ok := newTags[key]; !ok || nv != value {
			if !tagIgnored(key) {
				removed = append(removed, key)
			}
		}
	}

	if len(removed) > 0 {
		request := map[string]interface{}{
			"ResourceId.1": id,
			"RegionId":     client.RegionId,
			"ResourceType": "portrangelist",
		}
		for i, key := range removed {
			request[fmt.Sprintf("TagKey.%d", i+1)] = key
		}
		if _, err := callRpcWithRetry(client, "UntagResources", request, false); err != nil {
			return err
		}
	}

	if len(newTags) > 0 {
		request := map[string]interface{}{
			"ResourceId.1": id,
			"RegionId":     client.RegionId,
			"ResourceType": "portrangelist",
		}
		count := 1
		for key, value := range newTags {
			request[fmt.Sprintf("Tag.%d.Key", count)] = key
			request[fmt.Sprintf("Tag.%d.Value", count)] = value
			count++
		}
		if _, err := callRpcWithRetry(client, "TagResources", request, false); err != nil {
			return err
		}
	}
	return nil
}
