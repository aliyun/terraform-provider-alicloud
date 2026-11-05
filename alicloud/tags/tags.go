package tags

import (
	"context"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

// Generic tag helpers moved from the root package (alicloud/tags.go, the
// Tag struct from alicloud/extension_tags.go, the tag converters from
// alicloud/common.go and the TagResourceType family from
// alicloud/common.go); package alicloud re-exports them as var/type/const
// aliases so the flat layout keeps compiling unchanged.
// Product-SDK-typed tag helpers (setCdnTags, updateCdnTags, setVolumeTags,
// diffGpdbTags, gpdbTagsFromMap, the *TagsToMap/*TagIgnored families) stay in
// the root package and migrate with their products.

// Tag is the generic key/value pair used by DiffTags and TagsFromMap.
type Tag struct {
	Key   string
	Value string
}

func TagsSchemaWithElements() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeMap,
		Optional: true,
		Elem: &schema.Schema{
			Type: schema.TypeString,
		},
	}
}

func TagsSchemaForceNewWithElements() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeMap,
		Optional: true,
		ForceNew: true,
		Elem: &schema.Schema{
			Type: schema.TypeString,
		},
	}
}

func TagsSchemaComputedWithElements() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeMap,
		Optional: true,
		Computed: true,
		Elem: &schema.Schema{
			Type: schema.TypeString,
		},
	}
}

// lintignore: S006
// Deprecated: Use TagsSchemaWithElements instead, which declares the string
// element type explicitly.
func TagsSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeMap,
		Optional: true,
	}
}

// lintignore: S006
// Deprecated: Use TagsSchemaForceNewWithElements instead, which declares the
// string element type explicitly.
func TagsSchemaForceNew() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeMap,
		Optional: true,
		ForceNew: true,
	}
}

// lintignore: S006
// Deprecated: Use TagsSchemaComputedWithElements instead, which declares the
// string element type explicitly.
func TagsSchemaComputed() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeMap,
		Optional: true,
		Computed: true,
	}
}

func TagsSchemaWithIgnore() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeMap,
		Optional: true,
		Elem: &schema.Schema{
			Type:         schema.TypeString,
			ValidateFunc: validation.StringDoesNotMatch(regexp.MustCompile(`(^acs:.*)|(^aliyun.*)|(/.*http://.*\.\w+/gm)|(/.*https://.*\.\w+/gm)`), "It cannot begin with \"aliyun\", \"acs:\"; without \"http://\", and \"https://\"."),
		},
	}
}

// parse template_tags
func ParsingTemplateTags(d *schema.ResourceData) (map[string]interface{}, []string) {
	oraw, nraw := d.GetChange("template_tags")
	removedTags := oraw.(map[string]interface{})
	addedTags := nraw.(map[string]interface{})
	removed := make([]string, 0)
	for key, value := range removedTags {
		old, ok := addedTags[key]
		if !ok || old != value {
			removed = append(removed, key)
		}
	}

	return addedTags, removed
}

func ParsingTags(d *schema.ResourceData) (map[string]interface{}, []string) {
	oraw, nraw := d.GetChange("tags")
	removedTags := oraw.(map[string]interface{})
	addedTags := nraw.(map[string]interface{})
	// Build the list of what to remove
	removed := make([]string, 0)
	for key, value := range removedTags {
		old, ok := addedTags[key]
		if !ok || old != value {
			// Delete it!
			removed = append(removed, key)
		}
	}

	return addedTags, removed
}

func TagsToMap(tags interface{}) map[string]interface{} {
	result := make(map[string]interface{})
	if tags == nil {
		return result
	}
	switch v := tags.(type) {
	case map[string]interface{}:
		for key, value := range tags.(map[string]interface{}) {
			if !TagIgnored(key, value) {
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
			} else if v, ok := t["tag_key"]; ok {
				tagKey = v.(string)
				tagValue = t["tag_value"]
			} else if v, ok := t["tagKey"]; ok {
				tagKey = v.(string)
				tagValue = t["tagValue"]
			} else if v, ok := t["key"]; ok {
				tagKey = v.(string)
				tagValue = t["value"]
			}
			if !TagIgnored(tagKey, tagValue) {
				result[tagKey] = tagValue
			}
		}
	case []map[string]interface{}:
		for _, tag := range tags.([]map[string]interface{}) {
			var tagKey string
			var tagValue interface{}
			if v, ok := tag["tag_key"]; ok {
				tagKey = v.(string)
				tagValue = tag["tag_value"]
			}
			if !TagIgnored(tagKey, tagValue) {
				result[tagKey] = tagValue
			}
		}
	default:
		log.Printf("\u001B[31m[ERROR]\u001B[0m Unknown tags type %s. The tags value is: %v.", v, tags)
	}
	return result
}

func TagIgnored(tagKey string, tagValue interface{}) bool {
	filter := []string{"^aliyun", "^acs:", "^http://", "^https://", "^sae.do.not.delete"}
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

// SetTags is a helper to set the tags for a resource. It expects the
// tags field to be named "tags"
func SetTags(client *connectivity.AliyunClient, resourceType TagResourceType, d *schema.ResourceData) error {
	if d.HasChange("tags") {
		oraw, nraw := d.GetChange("tags")
		return UpdateTags(client, []string{d.Id()}, resourceType, oraw, nraw)
	}

	return nil
}

func UpdateTags(client *connectivity.AliyunClient, ids []string, resourceType TagResourceType, oraw, nraw interface{}) error {
	o := oraw.(map[string]interface{})
	n := nraw.(map[string]interface{})
	create, remove := DiffTags(TagsFromMap(o), TagsFromMap(n))

	// Set tags
	if len(remove) > 0 {
		request := ecs.CreateUntagResourcesRequest()
		request.ResourceType = string(resourceType)
		request.ResourceId = &ids
		var tagsKey []string

		for i := 0; i < len(remove); i = i + 20 {
			tagsKey = tagsKey[:0]

			if len(remove[i:]) <= 20 {
				for _, t := range remove[i:] {
					tagsKey = append(tagsKey, t.Key)
				}

				request.TagKey = &tagsKey

				wait := helper.IncrementalWait(1*time.Second, 1*time.Second)
				err := retry.Retry(10*time.Minute, func() *retry.RetryError {
					raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
						return ecsClient.UntagResources(request)
					})
					if err != nil {
						if helper.NeedRetry(err) {
							wait()
							return retry.RetryableError(err)

						}
						return retry.NonRetryableError(err)
					}
					logs.AddDebug(request.GetActionName(), raw, request.RpcRequest, request)
					return nil
				})

				if err != nil {
					return errs.WrapErrorf(err, errs.DefaultErrorMsg, ids, request.GetActionName(), errs.AlibabaCloudSdkGoERROR)
				}
			} else {
				for _, t := range remove[i : i+20] {
					tagsKey = append(tagsKey, t.Key)
				}

				request.TagKey = &tagsKey

				wait := helper.IncrementalWait(1*time.Second, 1*time.Second)
				err := retry.Retry(10*time.Minute, func() *retry.RetryError {
					raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
						return ecsClient.UntagResources(request)
					})
					if err != nil {
						if helper.NeedRetry(err) {
							wait()
							return retry.RetryableError(err)

						}
						return retry.NonRetryableError(err)
					}
					logs.AddDebug(request.GetActionName(), raw, request.RpcRequest, request)
					return nil
				})

				if err != nil {
					return errs.WrapErrorf(err, errs.DefaultErrorMsg, ids, request.GetActionName(), errs.AlibabaCloudSdkGoERROR)
				}
			}
		}
	}

	if len(create) > 0 {
		request := ecs.CreateTagResourcesRequest()
		request.ResourceType = string(resourceType)
		request.ResourceId = &ids
		var tags []ecs.TagResourcesTag

		for i := 0; i < len(create); i = i + 20 {
			tags = tags[:0]

			if len(create[i:]) <= 20 {
				for _, t := range create[i:] {
					tags = append(tags, ecs.TagResourcesTag{
						Key:   t.Key,
						Value: t.Value,
					})
				}

				request.Tag = &tags

				wait := helper.IncrementalWait(1*time.Second, 1*time.Second)
				err := retry.Retry(10*time.Minute, func() *retry.RetryError {
					raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
						return ecsClient.TagResources(request)
					})
					if err != nil {
						if helper.NeedRetry(err) {
							wait()
							return retry.RetryableError(err)

						}
						return retry.NonRetryableError(err)
					}
					logs.AddDebug(request.GetActionName(), raw, request.RpcRequest, request)
					return nil
				})

				if err != nil {
					return errs.WrapErrorf(err, errs.DefaultErrorMsg, ids, request.GetActionName(), errs.AlibabaCloudSdkGoERROR)
				}
			} else {
				for _, t := range create[i : i+20] {
					tags = append(tags, ecs.TagResourcesTag{
						Key:   t.Key,
						Value: t.Value,
					})
				}

				request.Tag = &tags

				wait := helper.IncrementalWait(1*time.Second, 1*time.Second)
				err := retry.Retry(10*time.Minute, func() *retry.RetryError {
					raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
						return ecsClient.TagResources(request)
					})
					if err != nil {
						if helper.NeedRetry(err) {
							wait()
							return retry.RetryableError(err)

						}
						return retry.NonRetryableError(err)
					}
					logs.AddDebug(request.GetActionName(), raw, request.RpcRequest, request)
					return nil
				})

				if err != nil {
					return errs.WrapErrorf(err, errs.DefaultErrorMsg, ids, request.GetActionName(), errs.AlibabaCloudSdkGoERROR)
				}
			}
		}
	}

	return nil
}

// DiffTags takes our tags locally and the ones remotely and returns
// the set of tags that must be created, and the set of tags that must
// be destroyed.
func DiffTags(oldTags, newTags []Tag) ([]Tag, []Tag) {
	// First, we're creating everything we have
	create := make(map[string]interface{})
	for _, t := range newTags {
		create[t.Key] = t.Value
	}

	// Build the list of what to remove
	var remove []Tag
	for _, t := range oldTags {
		old, ok := create[t.Key]
		if !ok || old != t.Value {
			// Delete it!
			remove = append(remove, t)
		}
	}

	return TagsFromMap(create), remove
}

// TagsFromMap returns the tags for the given map of data.
func TagsFromMap(m map[string]interface{}) []Tag {
	result := make([]Tag, 0, len(m))
	for k, v := range m {
		result = append(result, Tag{
			Key:   k,
			Value: v.(string),
		})
	}

	return result
}

func TagsMapEqual(expectMap map[string]interface{}, compareMap map[string]string) bool {
	if len(expectMap) != len(compareMap) {
		return false
	} else {
		for key, eVal := range expectMap {
			if eStr, ok := eVal.(string); !ok {
				// type is mismatch.
				return false
			} else {
				if cStr, ok := compareMap[key]; ok {
					if eStr != cStr {
						return false
					}
				} else {
					return false
				}
			}
		}
	}
	return true
}

func IgnoredTags(tagKey string, tagValue interface{}) bool {
	filter := []string{"^aliyun", "^acs:", "^http://", "^https://"}
	for _, v := range filter {
		ok, _ := regexp.MatchString(v, tagKey)
		if ok {
			log.Printf("[DEBUG] Found Alibaba Cloud specific tag with key: %s and value: %s, ignoring.\n", tagKey, tagValue)
			return true
		}
	}
	return false
}

// TagResourceType enumerates the ECS-style resource types accepted by the
// generic tag helpers (moved from alicloud/common.go).
type TagResourceType string

const (
	TagResourceImage         = TagResourceType("image")
	TagResourceInstance      = TagResourceType("instance")
	TagResourceAcl           = TagResourceType("acl")
	TagResourceCertificate   = TagResourceType("certificate")
	TagResourceDisk          = TagResourceType("disk")
	TagResourceSecurityGroup = TagResourceType("securitygroup")
	TagResourceCdn           = TagResourceType("DOMAIN")
	TagResourceApp           = TagResourceType("app")
	TagResourceTopic         = TagResourceType("topic")
	TagResourceCluster       = TagResourceType("cluster")
)

func ConvertTags(tagsMap map[string]interface{}) []map[string]interface{} {
	tags := make([]map[string]interface{}, 0)
	for key, value := range tagsMap {
		if value != nil {
			if v, ok := value.(string); ok {
				tags = append(tags, map[string]interface{}{
					"Key":   key,
					"Value": v,
				})
			}
		}
	}

	return tags
}

func ConvertLowercaseTags(tagsMap map[string]interface{}) []map[string]interface{} {
	tags := make([]map[string]interface{}, 0)
	for key, value := range tagsMap {
		if value != nil {
			if v, ok := value.(string); ok {
				tags = append(tags, map[string]interface{}{
					"key":   key,
					"value": v,
				})
			}
		}
	}

	return tags
}

func ConvertTagsForKms(tagsMap map[string]interface{}) []map[string]interface{} {
	tags := make([]map[string]interface{}, 0)
	for key, value := range tagsMap {
		if value != nil {
			if v, ok := value.(string); ok {
				tags = append(tags, map[string]interface{}{
					"TagKey":   key,
					"TagValue": v,
				})
			}
		}
	}

	return tags
}

func ExpandTagsToMap(originMap map[string]interface{}, tags []map[string]interface{}) map[string]interface{} {
	for i, tag := range tags {
		for key, value := range tag {
			if key == "Key" || key == "Value" {
				newKey := "Tag" + "." + strconv.Itoa(i+1) + "." + key
				originMap[newKey] = fmt.Sprintf("%v", value)
			}
		}
	}
	return originMap
}

func ExpandTagsToMapWithTags(originMap map[string]interface{}, tags []map[string]interface{}) map[string]interface{} {
	for i, tag := range tags {
		for key, value := range tag {
			if key == "Key" || key == "Value" {
				newKey := "Tags" + "." + strconv.Itoa(i+1) + "." + key
				originMap[newKey] = fmt.Sprintf("%v", value)
			}
		}
	}

	return originMap
}

// MergeTags overlays resourceTags on top of defaultTags: a key present in the
// resource's own tags overrides the provider default of the same key, matching
// the AWS default_tags precedence. A nil defaultTags is treated as empty so
// callers can pass client.DefaultTags unconditionally.
func MergeTags(defaultTags, resourceTags map[string]interface{}) map[string]interface{} {
	result := make(map[string]interface{}, len(defaultTags)+len(resourceTags))
	for k, v := range defaultTags {
		result[k] = v
	}
	for k, v := range resourceTags {
		result[k] = v
	}
	return result
}

// TagsAllSchema is the Computed mirror of a resource's tags field that carries
// the provider default_tags merged with the resource's own tags. Resources that
// opt into provider-level default_tags declare it alongside their tags field;
// CustomizeDiffDefaultTags keeps it in sync so a change to provider default_tags
// alone plans an update. With no default_tags configured tags_all equals tags,
// so behaviour is unchanged for providers that do not configure default_tags.
func TagsAllSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeMap,
		Computed: true,
		Elem: &schema.Schema{
			Type: schema.TypeString,
		},
	}
}

// tagsMapsEqual reports whether two tag maps hold the same key/value pairs.
// Values are compared by their string form so a default tag value read back as
// a Go string matches the interface{} the provider configuration held.
func tagsMapsEqual(a, b map[string]interface{}) bool {
	if len(a) != len(b) {
		return false
	}
	for k, v := range a {
		bv, ok := b[k]
		if !ok || fmt.Sprintf("%v", v) != fmt.Sprintf("%v", bv) {
			return false
		}
	}
	return true
}

// StripDefaultTags returns cloudTags without the entries that match a provider
// default tag (same key and value). Use it when reading a resource's tags back
// into its user-managed tags field so that default tags are not reported as
// user-configured drift; the merged set is reported via tags_all instead. A
// resource that overrides a default key with a different value keeps that
// override in the returned tags.
func StripDefaultTags(cloudTags map[string]string, defaultTags map[string]interface{}) map[string]string {
	result := make(map[string]string, len(cloudTags))
	for k, v := range cloudTags {
		if dv, ok := defaultTags[k]; ok && fmt.Sprintf("%v", dv) == v {
			continue
		}
		result[k] = v
	}
	return result
}

// SetTagsWithDefault is the default_tags-aware counterpart of SetTags. When
// provider default_tags are configured the effective tag set is the merge of
// default_tags and the resource's tags, so a change to provider default_tags
// (with the resource's own tags unchanged) is still propagated. The prior
// effective set is read from tags_all, which the Read path keeps equal to the
// cloud's full tag set; with no default_tags configured it falls back to
// SetTags, preserving the original tags-only behaviour.
func SetTagsWithDefault(client *connectivity.AliyunClient, resourceType TagResourceType, d *schema.ResourceData) error {
	if len(client.DefaultTags) == 0 {
		return SetTags(client, resourceType, d)
	}
	oldEffective, _ := d.Get("tags_all").(map[string]interface{})
	newEffective := MergeTags(client.DefaultTags, d.Get("tags").(map[string]interface{}))
	if !tagsMapsEqual(oldEffective, newEffective) {
		return UpdateTags(client, []string{d.Id()}, resourceType, oldEffective, newEffective)
	}
	return nil
}

// CustomizeDiffDefaultTags keeps a resource's tags_all equal to the merge of
// the provider default_tags and the resource's own tags. It is the only place
// that can turn a change to provider default_tags (which the resource's tags
// field does not reflect) into a planned update: when the merged set differs
// from the tags_all held in state it marks tags_all for recompute, which drives
// the resource's Update so SetTagsWithDefault can apply the delta. With no
// default_tags configured it is a no-op and the tags field drives updates as
// before. A nil client (meta unavailable at plan time) is treated as no
// default_tags, so the resource still updates on its own tags field.
func CustomizeDiffDefaultTags(_ context.Context, d *schema.ResourceDiff, meta interface{}) error {
	client, ok := meta.(*connectivity.AliyunClient)
	if !ok || client == nil || len(client.DefaultTags) == 0 {
		return nil
	}
	// On a brand-new resource d.Id() is empty: tags_all has no prior state to
	// reconcile against, and Create sets it after apply, so skip the drift
	// guard. (*schema.ResourceDiff has no IsNewResource; Id()=="" is the SDK
	// v2 CustomizeDiff convention, e.g. cloudfirewall instance v2.)
	if d.Id() == "" {
		return nil
	}
	oldEffective, _ := d.Get("tags_all").(map[string]interface{})
	newEffective := MergeTags(client.DefaultTags, d.Get("tags").(map[string]interface{}))
	if !tagsMapsEqual(oldEffective, newEffective) {
		return d.SetNewComputed("tags_all")
	}
	return nil
}
