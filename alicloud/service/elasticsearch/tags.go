package elasticsearch

import (
	"log"
	"regexp"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/elasticsearch"
)

// elasticsearchTagsToMap converts the SDK tag slice into the Terraform map
// representation, dropping Alibaba Cloud system tags.
func elasticsearchTagsToMap(tags []elasticsearch.Tag) map[string]string {
	result := make(map[string]string)
	for _, t := range tags {
		if !elasticsearchTagIgnored(t.TagKey, t.TagValue) {
			result[t.TagKey] = t.TagValue
		}
	}

	return result
}

// elasticsearchTagIgnored reports whether the tag is an Alibaba Cloud system
// tag that must be ignored during read/diff.
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
