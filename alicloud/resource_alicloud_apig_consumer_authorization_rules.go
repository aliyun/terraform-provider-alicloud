package alicloud

import (
	"fmt"
	"log"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudApigConsumerAuthorizationRules() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudApigConsumerAuthorizationRulesCreate,
		Read:   resourceAliCloudApigConsumerAuthorizationRulesRead,
		Delete: resourceAliCloudApigConsumerAuthorizationRulesDelete,
		Importer: &schema.ResourceImporter{
			State: resourceAliCloudApigConsumerAuthorizationRulesImport,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"consumer_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"environment_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"parent_resource_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"resource_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Default:      "HttpApiRoute",
				ValidateFunc: StringInSlice([]string{"HttpApiRoute"}, false),
			},
			"resource_ids": {
				Type:     schema.TypeSet,
				Required: true,
				ForceNew: true,
				MinItems: 1,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"principal_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Default:      "Consumer",
				ValidateFunc: StringInSlice([]string{"Consumer"}, false),
			},
			"expire_mode": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Default:      "LongTerm",
				ValidateFunc: StringInSlice([]string{"LongTerm"}, false),
			},
			"authorization_rule_ids": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func resourceAliCloudApigConsumerAuthorizationRulesImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	parts, err := apigParseCompositeID(d.Id(), 5)
	if err != nil {
		return nil, err
	}
	resourceIDs := strings.Split(parts[4], ",")
	seen := make(map[string]struct{}, len(resourceIDs))
	for _, resourceID := range resourceIDs {
		if resourceID == "" {
			return nil, fmt.Errorf("invalid import ID %q: resource IDs must not be empty", d.Id())
		}
		if _, exists := seen[resourceID]; exists {
			return nil, fmt.Errorf("invalid import ID %q: resource IDs must be unique", d.Id())
		}
		seen[resourceID] = struct{}{}
	}
	if err := d.Set("consumer_id", parts[0]); err != nil {
		return nil, err
	}
	if err := d.Set("environment_id", parts[1]); err != nil {
		return nil, err
	}
	if err := d.Set("parent_resource_id", parts[2]); err != nil {
		return nil, err
	}
	if err := d.Set("resource_type", parts[3]); err != nil {
		return nil, err
	}
	if err := d.Set("resource_ids", resourceIDs); err != nil {
		return nil, err
	}
	if err := d.Set("principal_type", "Consumer"); err != nil {
		return nil, err
	}
	if err := d.Set("expire_mode", "LongTerm"); err != nil {
		return nil, err
	}
	d.SetId(strings.Join(parts[:4], ":"))
	return []*schema.ResourceData{d}, nil
}

func apigAuthorizationRulesRequest(d *schema.ResourceData) map[string]interface{} {
	rules := make([]map[string]interface{}, 0, d.Get("resource_ids").(*schema.Set).Len())
	for _, resourceID := range apigSortedStringSet(d.Get("resource_ids")) {
		rules = append(rules, map[string]interface{}{
			"consumerId":    d.Get("consumer_id"),
			"expireMode":    d.Get("expire_mode"),
			"principalType": d.Get("principal_type"),
			"resourceType":  d.Get("resource_type"),
			"resourceIdentifier": map[string]interface{}{
				"resourceId":       resourceID,
				"parentResourceId": d.Get("parent_resource_id"),
				"environmentId":    d.Get("environment_id"),
			},
		})
	}
	return map[string]interface{}{"authorizationRules": rules}
}

func resourceAliCloudApigConsumerAuthorizationRulesCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "/v1/authorization-rules"
	body := apigAuthorizationRulesRequest(d)
	var response map[string]interface{}
	var err error
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RoaPost("APIG", "2024-03-27", action, nil, nil, body, true)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_apig_consumer_authorization_rules", action, AlibabaCloudSdkGoERROR)
	}
	data, err := apigResponseData(response)
	if err != nil {
		return err
	}
	createdIDs := apigStringSlice(data["consumerAuthorizationRuleIds"])
	resourceIDs := apigSortedStringSet(d.Get("resource_ids"))
	if len(createdIDs) != len(resourceIDs) {
		return fmt.Errorf("APIG created %d authorization rules, expected %d", len(createdIDs), len(resourceIDs))
	}
	d.SetId(fmt.Sprintf("%s:%s:%s:%s", d.Get("consumer_id"), d.Get("environment_id"), d.Get("parent_resource_id"), d.Get("resource_type")))

	createdIDSet := make(map[string]struct{}, len(createdIDs))
	for _, createdID := range createdIDs {
		createdIDSet[createdID] = struct{}{}
	}
	var observed map[string]string
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		observed, err = apigQueryAuthorizationRules(client, d)
		if err != nil {
			if NeedRetry(err) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		if len(observed) != len(resourceIDs) {
			return resource.RetryableError(fmt.Errorf("APIG authorization rules are not visible yet"))
		}
		for _, ruleID := range observed {
			if _, created := createdIDSet[ruleID]; !created {
				return resource.NonRetryableError(fmt.Errorf("APIG authorization read-back returned an unexpected rule ID"))
			}
		}
		return nil
	})
	if err != nil {
		return err
	}
	if err := d.Set("authorization_rule_ids", observed); err != nil {
		return err
	}
	return nil
}

func resourceAliCloudApigConsumerAuthorizationRulesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	observed, err := apigQueryAuthorizationRules(client, d)
	if err != nil {
		return WrapError(err)
	}
	if len(observed) == 0 {
		if !d.IsNewResource() {
			log.Printf("[DEBUG] Resource alicloud_apig_consumer_authorization_rules %s was not found", d.Id())
		}
		d.SetId("")
		return nil
	}
	resourceIDs := make([]string, 0, len(observed))
	for resourceID := range observed {
		resourceIDs = append(resourceIDs, resourceID)
	}
	sort.Strings(resourceIDs)
	if err := d.Set("resource_ids", resourceIDs); err != nil {
		return err
	}
	if err := d.Set("authorization_rule_ids", observed); err != nil {
		return err
	}
	return nil
}

func apigQueryAuthorizationRules(client *connectivity.AliyunClient, d *schema.ResourceData) (map[string]string, error) {
	desired := make(map[string]struct{})
	for _, resourceID := range apigSortedStringSet(d.Get("resource_ids")) {
		desired[resourceID] = struct{}{}
	}
	result := make(map[string]string, len(desired))
	const pageSize = 100
	for page := 1; ; page++ {
		query := map[string]*string{
			"consumerId":       StringPointer(d.Get("consumer_id").(string)),
			"environmentId":    StringPointer(d.Get("environment_id").(string)),
			"parentResourceId": StringPointer(d.Get("parent_resource_id").(string)),
			"principalType":    StringPointer(d.Get("principal_type").(string)),
			"resourceType":     StringPointer(d.Get("resource_type").(string)),
			"pageNumber":       StringPointer(strconv.Itoa(page)),
			"pageSize":         StringPointer(strconv.Itoa(pageSize)),
		}
		response, err := client.RoaGet("APIG", "2024-03-27", "/v1/authorization-rules", query, nil, nil)
		if err != nil {
			return nil, err
		}
		data, err := apigResponseData(response)
		if err != nil {
			return nil, err
		}
		items := apigObjectSlice(data["items"])
		for _, item := range items {
			resourceID, resourceOK := item["resourceId"].(string)
			ruleID, ruleOK := item["consumerAuthorizationRuleId"].(string)
			if !resourceOK || !ruleOK {
				continue
			}
			if _, wanted := desired[resourceID]; !wanted || item["expireMode"] != d.Get("expire_mode") {
				continue
			}
			if _, duplicate := result[resourceID]; duplicate {
				return nil, fmt.Errorf("APIG returned multiple authorization rules for resource %s", resourceID)
			}
			result[resourceID] = ruleID
		}
		if len(items) < pageSize {
			return result, nil
		}
	}
}

func resourceAliCloudApigConsumerAuthorizationRulesDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	rules, err := apigQueryAuthorizationRules(client, d)
	if err != nil {
		return WrapError(err)
	}
	ruleIDs := make([]string, 0, len(rules))
	for _, ruleID := range rules {
		ruleIDs = append(ruleIDs, ruleID)
	}
	if len(ruleIDs) == 0 {
		return nil
	}
	sort.Strings(ruleIDs)
	query := map[string]*string{"consumerAuthorizationRuleIds": StringPointer(strings.Join(ruleIDs, ","))}
	action := "/v1/authorization-rules"
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		_, err = client.RoaDelete("APIG", "2024-03-27", action, query, nil, nil, true)
		if err != nil {
			if IsExpectedErrors(err, []string{"NotFound.ConsumerAuthorizationRuleNotFound"}) || NotFoundError(err) {
				return nil
			}
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		remaining, err := apigQueryAuthorizationRules(client, d)
		if err != nil {
			if NeedRetry(err) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		if len(remaining) != 0 {
			return resource.RetryableError(fmt.Errorf("APIG authorization rules still exist"))
		}
		return nil
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}
	return nil
}
