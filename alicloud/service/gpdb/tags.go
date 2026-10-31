package gpdb

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/services/gpdb"
)

// diffTags computes the tag create/remove sets for gpdb resources, whose tags
// use the SDK-specific gpdb.TagResourcesTag type. Migrated verbatim from the
// flat alicloud package's alicloud/tags.go (diffGpdbTags), renamed to drop the
// product prefix per the subpackage convention.
func diffTags(oldTags, newTags []gpdb.TagResourcesTag) ([]gpdb.TagResourcesTag, []gpdb.TagResourcesTag) {
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
