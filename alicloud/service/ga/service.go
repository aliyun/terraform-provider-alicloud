// Package ga holds the Framework-native Global Accelerator (GA) resources and
// data sources. The helpers here mirror the SDKv2 GaService methods that the
// original alicloud_ga_service_extension implementation relied on, so the
// state-convergence, retry, NotFound and tagging-tolerance semantics survive
// the move to terraform-plugin-framework. They live in this package because
// package alicloud (the monolith that still hosts the SDKv2 GA resources)
// cannot be imported from a service package without a cycle:
// alicloud -> provider/registry -> service/ga.
package ga

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/errors"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/denverdino/aliyungo/common"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
)

const (
	// Default timeouts match the SDKv2 resource's schema.DefaultTimeout(5m).
	defaultCreateTimeout = 5 * time.Minute
	defaultUpdateTimeout = 5 * time.Minute
	defaultDeleteTimeout = 5 * time.Minute

	// serviceExtensionResourceType is the ResourceType value the GA tagging
	// backend does not actually register for service extensions; TagResources
	// and ListTagResources reject it with Invalid.ResourceType. The Read/Update
	// paths tolerate that so refresh and destroy stay reachable, matching the
	// original SDKv2 behavior.
	serviceExtensionResourceType = "serviceextension"
)

// isExpectedErrors reports whether err carries one of the given codes. It is a
// trimmed port of package alicloud's IsExpectedErrors covering the error types
// GA RpcPost actually returns — *tea.SDKError, *errors.ServerError,
// *errs.ProviderError, *common.Error and the wrapped *errs.ComplexError — plus a
// string-contains fallback so wrapped or future error shapes still match.
func isExpectedErrors(err error, expectCodes []string) bool {
	if err == nil {
		return false
	}
	if e, ok := err.(*errs.ComplexError); ok {
		return isExpectedErrors(e.Cause, expectCodes)
	}
	if e, ok := err.(*tea.SDKError); ok {
		for _, code := range expectCodes {
			if e.Code != nil && (*e.Code == code || strings.HasPrefix(code, *e.Code)) {
				return true
			}
			if e.Data != nil && strings.Contains(*e.Data, code) {
				return true
			}
		}
		return false
	}
	if e, ok := err.(*errors.ServerError); ok {
		for _, code := range expectCodes {
			if e.ErrorCode() == code || strings.Contains(e.Message(), code) {
				return true
			}
		}
		return false
	}
	if e, ok := err.(*errs.ProviderError); ok {
		for _, code := range expectCodes {
			if e.ErrorCode() == code || strings.Contains(e.Message(), code) {
				return true
			}
		}
		return false
	}
	if e, ok := err.(*common.Error); ok {
		for _, code := range expectCodes {
			if e.Code == code || fmt.Sprint(e.StatusCode) == code || strings.Contains(e.Message, code) {
				return true
			}
		}
		return false
	}
	for _, code := range expectCodes {
		if strings.Contains(err.Error(), code) {
			return true
		}
	}
	return false
}

// buildStateConf mirrors package alicloud's BuildStateConf so the asynchronous
// GA operations (CreateServiceExtension, UpdateServiceExtension, the add/remove
// and associate/dissociate calls) wait for state convergence exactly as before.
func buildStateConf(pending, target []string, timeout, delay time.Duration, f retry.StateRefreshFunc) *retry.StateChangeConf {
	return &retry.StateChangeConf{
		Pending:    pending,
		Target:     target,
		Refresh:    f,
		Timeout:    timeout,
		Delay:      delay,
		MinTimeout: 3 * time.Second,
	}
}

// rpcPostRetry calls a GA RpcPost action with the same throttle retry policy the
// SDKv2 resources used: NeedRetry-driven exponential backoff up to timeout.
func rpcPostRetry(client *connectivity.AliyunClient, action string, request map[string]interface{}, timeout time.Duration, autoRetry bool) (map[string]interface{}, error) {
	var response map[string]interface{}
	var err error
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(timeout, func() *retry.RetryError {
		response, err = client.RpcPost("Ga", "2019-11-20", action, nil, request, autoRetry)
		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(action, response, request)
	return response, err
}

// describeServiceExtension queries GetServiceExtension and returns the response
// object. NotExist.ServiceExtension is mapped to a NotFound error so Read and
// Delete can remove the resource from state, matching GaService.
func describeServiceExtension(client *connectivity.AliyunClient, id string) (map[string]interface{}, error) {
	action := "GetServiceExtension"
	request := map[string]interface{}{
		"RegionId":           client.RegionId,
		"ServiceExtensionId": id,
	}
	response, err := rpcPostRetry(client, action, request, 5*time.Minute, true)
	if err != nil {
		if isExpectedErrors(err, []string{"NotExist.ServiceExtension"}) {
			return nil, errs.WrapErrorf(errs.NotFoundErr("Ga:ServiceExtension", id), errs.NotFoundMsg, errs.ProviderERROR, fmt.Sprint(response["RequestId"]))
		}
		return nil, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}
	v, err := jsonpath.Get("$", response)
	if err != nil {
		return nil, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$", response)
	}
	object, ok := v.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("the %s response carried no object: %v", action, response)
	}
	return object, nil
}

// serviceExtensionStateRefreshFunc returns the refresh function used by
// buildStateConf to wait for a service extension to leave the pending/updating
// state. NotFound collapses to an empty state so delete convergence terminates.
func serviceExtensionStateRefreshFunc(client *connectivity.AliyunClient, id string, failStates []string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := describeServiceExtension(client, id)
		if err != nil {
			if errs.NotFoundError(err) {
				return object, "", nil
			}
			return nil, "", errs.WrapError(err)
		}
		state := fmt.Sprint(object["State"])
		for _, fs := range failStates {
			if state == fs {
				return nil, "", fmt.Errorf("service extension %s entered fail state %q", id, state)
			}
		}
		return object, state, nil
	}
}

// tagsToMapGA maps the TagResources.TagResource list returned by
// ListTagResources onto a string-keyed map. GA returns each tag as
// {TagKey, TagValue}. Provider-managed prefixes (aliyun/acs:) are filtered.
func tagsToMapGA(tags interface{}) map[string]string {
	result := make(map[string]string)
	if tags == nil {
		return result
	}
	list, ok := tags.([]interface{})
	if !ok {
		return result
	}
	for _, t := range list {
		item, ok := t.(map[string]interface{})
		if !ok {
			continue
		}
		key, _ := item["TagKey"].(string)
		if key == "" {
			continue
		}
		if isIgnoredTag(key) {
			continue
		}
		result[key] = fmt.Sprint(item["TagValue"])
	}
	return result
}

// isIgnoredTag mirrors ignoredTags: drop Alibaba Cloud provider-managed tags.
func isIgnoredTag(tagKey string) bool {
	filter := []string{"^aliyun", "^acs:", "^http://", "^https://"}
	for _, v := range filter {
		if ok, _ := regexp.MatchString(v, tagKey); ok {
			return true
		}
	}
	return false
}

// listServiceExtensionTags calls ListTagResources for the service extension,
// returning the tag map. GA tagging rejects serviceextension with
// Invalid.ResourceType; that is tolerated as an empty tag set so Read/Update
// do not block refresh or destroy.
func listServiceExtensionTags(client *connectivity.AliyunClient, id string) (map[string]string, error) {
	action := "ListTagResources"
	request := map[string]interface{}{
		"RegionId":     client.RegionId,
		"ResourceType": serviceExtensionResourceType,
		"ResourceId.1": id,
	}
	var collected []interface{}
	for {
		response, err := rpcPostRetry(client, action, request, 5*time.Minute, true)
		if err != nil {
			return nil, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
		}
		v, err := jsonpath.Get("$.TagResources.TagResource", response)
		if err != nil {
			return nil, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, id, "$.TagResources.TagResource", response)
		}
		if list, ok := v.([]interface{}); ok {
			collected = append(collected, list...)
		}
		if nextToken, ok := response["NextToken"].(string); ok && nextToken != "" {
			request["NextToken"] = nextToken
			continue
		}
		break
	}
	return tagsToMapGA(collected), nil
}

// applyServiceExtensionTags reconciles the old and new tag sets via
// UntagResources/TagResources, mirroring GaService.SetResourceTags. The GA
// tagging backend rejects serviceextension with Invalid.ResourceType; that is
// tolerated (logged) so a tag change never fatal-blocks an otherwise-valid
// update, exactly as the SDKv2 resource did.
func applyServiceExtensionTags(client *connectivity.AliyunClient, id string, oldTags, newTags map[string]string) error {
	removed := make([]string, 0)
	for key := range oldTags {
		newVal, ok := newTags[key]
		if !ok || newVal != oldTags[key] {
			if !isIgnoredTag(key) {
				removed = append(removed, key)
			}
		}
	}

	if len(removed) > 0 {
		request := map[string]interface{}{
			"RegionId":     client.RegionId,
			"ResourceType": serviceExtensionResourceType,
			"ResourceId.1": id,
		}
		for i, key := range removed {
			request[fmt.Sprintf("TagKey.%d", i+1)] = key
		}
		_, err := rpcPostRetry(client, "UntagResources", request, 10*time.Minute, false)
		if err != nil {
			if isExpectedErrors(err, []string{"Invalid.ResourceType"}) {
				log.Printf("[WARN] alicloud_ga_service_extension tagging is not supported by the backend (Invalid.ResourceType), skipping tag removal for %s: %s", id, err)
			} else {
				return errs.WrapErrorf(err, errs.DefaultErrorMsg, id, "UntagResources", errs.AlibabaCloudSdkGoERROR)
			}
		}
	}

	if len(newTags) > 0 {
		request := map[string]interface{}{
			"RegionId":     client.RegionId,
			"ResourceType": serviceExtensionResourceType,
			"ResourceId.1": id,
		}
		count := 1
		for key, value := range newTags {
			if isIgnoredTag(key) {
				continue
			}
			request[fmt.Sprintf("Tag.%d.Key", count)] = key
			request[fmt.Sprintf("Tag.%d.Value", count)] = value
			count++
		}
		if count > 1 {
			_, err := rpcPostRetry(client, "TagResources", request, 10*time.Minute, false)
			if err != nil {
				if isExpectedErrors(err, []string{"Invalid.ResourceType"}) {
					log.Printf("[WARN] alicloud_ga_service_extension tagging is not supported by the backend (Invalid.ResourceType), skipping tag application for %s: %s", id, err)
				} else {
					return errs.WrapErrorf(err, errs.DefaultErrorMsg, id, "TagResources", errs.AlibabaCloudSdkGoERROR)
				}
			}
		}
	}
	return nil
}

// writeToFile writes the data source results to the requested output file,
// matching package alicloud's writeToFile (JSON, ~ expansion, replace).
func writeToFile(filePath string, data interface{}) error {
	if filePath == "" {
		return nil
	}
	var out string
	switch v := data.(type) {
	case string:
		out = v
	case nil:
		return nil
	default:
		bs, err := json.MarshalIndent(v, "", "\t")
		if err != nil {
			return fmt.Errorf("MarshalIndent data %#v got an error: %#v", data, err)
		}
		out = string(bs)
	}
	if strings.HasPrefix(filePath, "~") {
		home, err := userHomeDir()
		if err != nil {
			return err
		}
		if home != "" {
			filePath = strings.Replace(filePath, "~", home, 1)
		}
	}
	if _, err := os.Stat(filePath); err == nil {
		if err := os.Remove(filePath); err != nil {
			return err
		}
	}
	dir := filepath.Dir(filePath)
	if dir != "" && dir != "." {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return err
		}
	}
	return os.WriteFile(filePath, []byte(out), 0644)
}

func userHomeDir() (string, error) {
	if home := os.Getenv("HOME"); home != "" {
		return home, nil
	}
	return os.UserHomeDir()
}
