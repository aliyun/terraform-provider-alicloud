package alicloud

import (
	"fmt"
	"log"
	"regexp"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/cdn"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/elasticsearch"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ess"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/gpdb"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ots"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/slb"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/vpc"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/tags"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Generic tag helpers moved to alicloud/tags (tags/tags.go); the aliases
// below keep every flat file compiling unchanged. Product-SDK-typed tag
// helpers stay in this file and migrate with their products.

var String = helper.String

var (
	tagsSchemaWithElements         = tags.TagsSchemaWithElements
	tagsSchemaForceNewWithElements = tags.TagsSchemaForceNewWithElements
	tagsSchemaComputedWithElements = tags.TagsSchemaComputedWithElements
	tagsSchema                     = tags.TagsSchema
	tagsSchemaForceNew             = tags.TagsSchemaForceNew
	tagsSchemaComputed             = tags.TagsSchemaComputed
	tagsSchemaWithIgnore           = tags.TagsSchemaWithIgnore
	parsingTemplateTags            = tags.ParsingTemplateTags
	parsingTags                    = tags.ParsingTags
	tagsToMap                      = tags.TagsToMap
	tagIgnored                     = tags.TagIgnored
	setTags                        = tags.SetTags
	updateTags                     = tags.UpdateTags
	diffTags                       = tags.DiffTags
	tagsFromMap                    = tags.TagsFromMap
	tagsMapEqual                   = tags.TagsMapEqual
	ignoredTags                    = tags.IgnoredTags
)

func setCdnTags(client *connectivity.AliyunClient, resourceType TagResourceType, d *schema.ResourceData) error {
	if d.HasChange("tags") {
		oraw, nraw := d.GetChange("tags")
		return updateCdnTags(client, []string{d.Id()}, resourceType, oraw, nraw)
	}

	return nil
}

func setVolumeTags(client *connectivity.AliyunClient, resourceType TagResourceType, d *schema.ResourceData) error {
	ecsService := EcsService{client}
	if d.HasChange("volume_tags") {
		request := ecs.CreateDescribeDisksRequest()
		request.InstanceId = d.Id()
		var response *ecs.DescribeDisksResponse
		wait := incrementalWait(1*time.Second, 1*time.Second)
		err := retry.Retry(10*time.Minute, func() *retry.RetryError {
			raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
				return ecsClient.DescribeDisks(request)
			})
			if err != nil {
				if NeedRetry(err) {
					wait()
					return retry.RetryableError(err)

				}
				return retry.NonRetryableError(err)
			}
			addDebug(request.GetActionName(), raw, request.RpcRequest, request)
			response, _ = raw.(*ecs.DescribeDisksResponse)
			return nil
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
		}

		if len(response.Disks.Disk) == 0 {
			return WrapError(fmt.Errorf("the system disk cannot be queried in this instance %s. Please check whether you have permission to access the API DescribeDisks. Last response is: %v", d.Id(), response))
		}

		var ids []string
		systemDiskTag := make(map[string]interface{})
		for _, disk := range response.Disks.Disk {
			ids = append(ids, disk.DiskId)
			if disk.Type == "system" {
				for _, t := range disk.Tags.Tag {
					if !ecsService.ecsTagIgnored(t) {
						systemDiskTag[t.TagKey] = t.TagValue
					}
				}
			}
		}

		oraw, nraw := d.GetChange("volume_tags")
		if d.IsNewResource() {
			oraw = systemDiskTag
		}
		return updateTags(client, ids, resourceType, oraw, nraw)
	}

	return nil
}

func updateCdnTags(client *connectivity.AliyunClient, ids []string, resourceType TagResourceType, oraw, nraw interface{}) error {
	o := oraw.(map[string]interface{})
	n := nraw.(map[string]interface{})
	create, remove := diffTags(tagsFromMap(o), tagsFromMap(n))

	// Set tags
	if len(remove) > 0 {
		request := cdn.CreateUntagResourcesRequest()
		request.ResourceType = string(resourceType)
		request.ResourceId = &ids

		var tagsKey []string
		for _, t := range remove {
			tagsKey = append(tagsKey, t.Key)
		}
		request.TagKey = &tagsKey

		raw, err := client.WithCdnClient_new(func(cdnClient *cdn.Client) (interface{}, error) {
			return cdnClient.UntagResources(request)
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, ids, request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	}

	if len(create) > 0 {
		request := cdn.CreateTagResourcesRequest()
		request.ResourceType = string(resourceType)
		request.ResourceId = &ids

		var tags []cdn.TagResourcesTag
		for _, t := range create {
			tags = append(tags, cdn.TagResourcesTag{
				Key:   t.Key,
				Value: t.Value,
			})
		}
		request.Tag = &tags

		raw, err := client.WithCdnClient_new(func(cdnClient *cdn.Client) (interface{}, error) {
			return cdnClient.TagResources(request)
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, ids, request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	}

	return nil
}

func diffGpdbTags(oldTags, newTags []gpdb.TagResourcesTag) ([]gpdb.TagResourcesTag, []gpdb.TagResourcesTag) {
	// First, we're creating everything we have
	create := make(map[string]interface{})
	for _, t := range newTags {
		create[t.Key] = t.Value
	}

	// Build the list of what to remove
	var remove []gpdb.TagResourcesTag
	for _, t := range oldTags {
		old, ok := create[t.Key]
		if !ok || old != t.Value {
			// Delete it!
			remove = append(remove, t)
		}
	}

	return gpdbTagsFromMap(create), remove
}

func gpdbTagsFromMap(m map[string]interface{}) []gpdb.TagResourcesTag {
	result := make([]gpdb.TagResourcesTag, 0, len(m))
	for k, v := range m {
		result = append(result, gpdb.TagResourcesTag{
			Key:   k,
			Value: v.(string),
		})
	}

	return result
}

func ecsTagsToMap(tags []ecs.Tag) map[string]string {
	result := make(map[string]string)
	for _, t := range tags {
		if !ecsTagIgnored(t) {
			result[t.TagKey] = t.TagValue
		}
	}

	return result
}

func elasticsearchTagsToMap(tags []elasticsearch.Tag) map[string]string {
	result := make(map[string]string)
	for _, t := range tags {
		if !elasticsearchTagIgnored(t.TagKey, t.TagValue) {
			result[t.TagKey] = t.TagValue
		}
	}

	return result
}

func vpcTagsToMap(tags []vpc.Tag) map[string]string {
	result := make(map[string]string)
	for _, t := range tags {
		if !vpcTagIgnored(t) {
			result[t.Key] = t.Value
		}
	}
	return result
}

func cdnTagsToMap(tags []cdn.TagItem) map[string]string {
	result := make(map[string]string)
	for _, t := range tags {
		if !cdnTagIgnored(t) {
			result[t.Key] = t.Value
		}
	}

	return result
}

func slbTagsToMap(tags []slb.TagResource) map[string]string {
	result := make(map[string]string)
	for _, t := range tags {
		if !slbTagIgnored(t) {
			result[t.TagKey] = t.TagValue
		}
	}

	return result
}

func essTagsToMap(tags []ess.Tag) map[string]string {
	result := make(map[string]string)
	for _, t := range tags {
		if !essTagIgnored(t) {
			result[t.Key] = t.Value
		}
	}

	return result
}

func otsTagsToMap(tags []ots.TagInfo) map[string]string {
	result := make(map[string]string)
	for _, t := range tags {
		result[t.TagKey] = t.TagValue
	}

	return result
}

func albTagsToMap(tags interface{}) map[string]interface{} {
	result := make(map[string]interface{})
	if tags == nil {
		return result
	}
	switch v := tags.(type) {
	case map[string]interface{}:
		for key, value := range tags.(map[string]interface{}) {
			if !albTagIgnored(key, value) {
				result[key] = value
			}
		}
	case []interface{}:
		if len(tags.([]interface{})) < 1 {
			return result
		}
		for _, tag := range tags.([]interface{}) {
			t := tag.(map[string]interface{})
			var tagKey string
			var tagValue interface{}
			if v, ok := t["TagKey"]; ok {
				tagKey = v.(string)
				tagValue = t["TagValue"]
			} else if v, ok := t["Key"]; ok {
				tagKey = v.(string)
				tagValue = t["Value"]
			}
			if !albTagIgnored(tagKey, tagValue) {
				result[tagKey] = tagValue
			}
		}
	default:
		log.Printf("\u001B[31m[ERROR]\u001B[0m Unknown tags type %s. The tags value is: %v.", v, tags)
	}
	return result
}

// tagIgnored compares a tag against a list of strings and checks if it should be ignored or not
func ecsTagIgnored(t ecs.Tag) bool {
	filter := []string{"^aliyun", "^acs:", "^http://", "^https://"}
	for _, v := range filter {
		log.Printf("[DEBUG] Matching prefix %v with %v\n", v, t.TagKey)
		ok, _ := regexp.MatchString(v, t.TagKey)
		if ok {
			log.Printf("[DEBUG] Found Alibaba Cloud specific tag %s (val: %s), ignoring.\n", t.TagKey, t.TagValue)
			return true
		}
	}
	return false
}

func vpcTagIgnored(t vpc.Tag) bool {
	filter := []string{"^aliyun", "^acs:", "^http://", "^https://"}
	for _, v := range filter {
		log.Printf("[DEBUG] Matching prefix %v with %v\n", v, t.Key)
		ok, _ := regexp.MatchString(v, t.Key)
		if ok {
			log.Printf("[DEBUG] Found Alibaba Cloud specific tag %s (val: %s), ignoring.\n", t.Key, t.Value)
			return true
		}
	}
	return false
}

// tagIgnored compares a tag against a list of strings and checks if it should be ignored or not
func essTagIgnored(t ess.Tag) bool {
	filter := []string{"^aliyun", "^http://", "^https://"}
	for _, v := range filter {
		log.Printf("[DEBUG] Matching prefix %v with %v\n", v, t.Key)
		ok, _ := regexp.MatchString(v, t.Key)
		if ok {
			log.Printf("[DEBUG] Found Alibaba Cloud specific tag %s (val: %s), ignoring.\n", t.Key, t.Value)
			return true
		}
	}
	return false
}

func cdnTagIgnored(t cdn.TagItem) bool {
	filter := []string{"^aliyun", "^acs:", "^http://", "^https://"}
	for _, v := range filter {
		log.Printf("[DEBUG] Matching prefix %v with %v\n", v, t.Key)
		ok, _ := regexp.MatchString(v, t.Key)
		if ok {
			log.Printf("[DEBUG] Found Alibaba Cloud specific tag %s (val: %s), ignoring.\n", t.Key, t.Value)
			return true
		}
	}
	return false
}

func slbTagIgnored(t slb.TagResource) bool {
	filter := []string{"^aliyun", "^acs:", "^http://", "^https://"}
	for _, v := range filter {
		log.Printf("[DEBUG] Matching prefix %v with %v\n", v, t.TagKey)
		ok, _ := regexp.MatchString(v, t.TagKey)
		if ok {
			log.Printf("[DEBUG] Found Alibaba Cloud specific tag %s (val: %s), ignoring.\n", t.TagKey, t.TagValue)
			return true
		}
	}
	return false
}

func albTagIgnored(tagKey string, tagValue interface{}) bool {
	filter := []string{"^aliyun", "^acs:", "^http://", "^https://", "^ack", "^ingress"}
	for _, v := range filter {
		log.Printf("[DEBUG] Matching prefix %v with %v\n", v, tagKey)
		ok, _ := regexp.MatchString(v, tagKey)
		if ok {
			log.Printf("[DEBUG] Found Alibaba Cloud specific tag %s (val: %s), ignoring.\n", tagKey, tagValue)
			return true
		}
	}
	return false
}

func elasticsearchTagIgnored(tagKey, tagValue string) bool {
	filter := []string{"^aliyun", "^acs:", "^http://", "^https://"}
	for _, v := range filter {
		log.Printf("[DEBUG] Matching prefix %v with %v\n", v, tagKey)
		ok, _ := regexp.MatchString(v, tagKey)
		if ok {
			log.Printf("[DEBUG] Found Alibaba Cloud specific tag %s (val: %s), ignoring.\n", tagKey, tagValue)
			return true
		}
	}
	return false
}
