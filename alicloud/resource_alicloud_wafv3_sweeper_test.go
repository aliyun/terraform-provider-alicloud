package alicloud

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func init() {
	resource.AddTestSweepers("alicloud_wafv3_domain", &resource.Sweeper{
		Name: "alicloud_wafv3_domain",
		F:    testSweepWafv3Domains,
	})
	resource.AddTestSweepers("alicloud_wafv3_defense_template", &resource.Sweeper{
		Name: "alicloud_wafv3_defense_template",
		F:    testSweepWafv3DefenseTemplates,
	})
	resource.AddTestSweepers("alicloud_wafv3_defense_resource_group", &resource.Sweeper{
		Name: "alicloud_wafv3_defense_resource_group",
		F:    testSweepWafv3DefenseResourceGroups,
	})
	resource.AddTestSweepers("alicloud_wafv3_defense_rule", &resource.Sweeper{
		Name: "alicloud_wafv3_defense_rule",
		F:    testSweepWafv3DefenseRules,
	})
	resource.AddTestSweepers("alicloud_wafv3_address_book", &resource.Sweeper{
		Name: "alicloud_wafv3_address_book",
		F:    testSweepWafv3AddressBooks,
	})
}

// wafv3SweepPrefixes lists the test-artifact name prefixes that wafv3 sweepers match
// (case-insensitive). These cover the naming conventions used across wafv3 *_test.go files:
//   - "tf-testacc" / "tf_testacc": canonical provider-wide test prefixes.
//   - "tfacc": used by wafv3_address_book, wafv3_defense_template, wafv3_defense_resource_group tests.
//   - "tftest": used by wafv3_domain tests (tftest<rand>.tftest.top).
//   - "tf-test" / "tf-tpl" / "tf-whitelist": used by wafv3_defense_rule tests.
var wafv3SweepPrefixes = []string{
	"tf-testacc",
	"tf_testacc",
	"tfacc",
	"tftest",
	"tf-test",
	"tf-tpl",
	"tf-whitelist",
}

// wafv3SweepPrefixMatch returns true if name matches any wafv3 sweep prefix (case-insensitive).
func wafv3SweepPrefixMatch(name string) bool {
	lower := strings.ToLower(name)
	for _, prefix := range wafv3SweepPrefixes {
		if strings.HasPrefix(lower, strings.ToLower(prefix)) {
			return true
		}
	}
	return false
}

// wafv3InstanceID fetches the single WAF v3 instance ID for the region via DescribeInstance.
// Returns "" when there is no WAF v3 instance, so the sweeper is a no-op.
func wafv3InstanceID(client *connectivity.AliyunClient) string {
	action := "DescribeInstance"
	request := map[string]interface{}{
		"RegionId": client.RegionId,
	}
	response, err := client.RpcPost("waf-openapi", "2021-10-01", action, nil, request, false)
	if err != nil {
		log.Printf("[ERROR] Failed to retrieve WAF v3 instance: %s", err)
		return ""
	}
	if response["InstanceId"] == nil {
		return ""
	}
	return fmt.Sprint(response["InstanceId"])
}

// testSweepWafv3Domains sweeps alicloud_wafv3_domain resources.
// Domain names may include real存/ICP-backed domains, so the prefix filter is ALWAYS
// applied (even in SWEEP_ALL) to preserve them.
func testSweepWafv3Domains(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return WrapError(err)
	}
	client := rawClient.(*connectivity.AliyunClient)

	instanceId := wafv3InstanceID(client)
	if instanceId == "" {
		return nil
	}

	action := "DescribeDomains"
	pageNumber := 1
	pageSize := PageSizeLarge
	domains := make([]string, 0)
	for {
		request := map[string]interface{}{
			"RegionId":   client.RegionId,
			"InstanceId": instanceId,
			"PageNumber":  pageNumber,
			"PageSize":    pageSize,
		}
		response, err := client.RpcPost("waf-openapi", "2021-10-01", action, nil, request, false)
		if err != nil {
			log.Printf("[ERROR] Failed to retrieve WAF v3 domains: %s", err)
			return nil
		}
		// Handle missing list key as empty (empty-list defect pattern).
		resp, err := jsonpath.Get("$.Domains", response)
		if err != nil {
			return nil
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			item, ok := v.(map[string]interface{})
			if !ok {
				continue
			}
			domains = append(domains, fmt.Sprint(item["Domain"]))
		}
		if len(result) < pageSize {
			break
		}
		pageNumber++
	}

	for _, domain := range domains {
		// NODELETING guard runs in BOTH sweepAll + normal mode.
		if protectedSweepVpcName(domain) {
			log.Printf("[INFO] Skipping WAF v3 domain (NODELETING): %s", domain)
			continue
		}
		// ALWAYS apply prefix filter to protect real存/ICP-backed domains,
		// EVEN in SWEEP_ALL mode.
		if !wafv3SweepPrefixMatch(domain) {
			log.Printf("[INFO] Skipping WAF v3 domain (no test prefix): %s", domain)
			continue
		}
		log.Printf("[INFO] Deleting WAF v3 domain: %s", domain)
		deleteRequest := map[string]interface{}{
			"RegionId":   client.RegionId,
			"InstanceId": instanceId,
			"Domain":     domain,
		}
		if _, err := client.RpcPost("waf-openapi", "2021-10-01", "DeleteDomain", nil, deleteRequest, false); err != nil {
			log.Printf("[ERROR] Failed to delete WAF v3 domain (%s): %s", domain, err)
		}
	}
	return nil
}

// testSweepWafv3DefenseTemplates sweeps alicloud_wafv3_defense_template resources.
func testSweepWafv3DefenseTemplates(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return WrapError(err)
	}
	client := rawClient.(*connectivity.AliyunClient)

	instanceId := wafv3InstanceID(client)
	if instanceId == "" {
		return nil
	}

	type templateItem struct {
		TemplateId string
		Name       string
	}
	items := make([]templateItem, 0)
	action := "DescribeDefenseTemplates"
	pageNumber := 1
	pageSize := PageSizeLarge
	for {
		request := map[string]interface{}{
			"RegionId":   client.RegionId,
			"InstanceId": instanceId,
			"PageNumber":  pageNumber,
			"PageSize":    pageSize,
		}
		response, err := client.RpcPost("waf-openapi", "2021-10-01", action, nil, request, false)
		if err != nil {
			log.Printf("[ERROR] Failed to retrieve WAF v3 defense templates: %s", err)
			return nil
		}
		resp, err := jsonpath.Get("$.Templates", response)
		if err != nil {
			return nil
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			item, ok := v.(map[string]interface{})
			if !ok {
				continue
			}
			items = append(items, templateItem{
				TemplateId: fmt.Sprint(item["TemplateId"]),
				Name:       fmt.Sprint(item["TemplateName"]),
			})
		}
		if len(result) < pageSize {
			break
		}
		pageNumber++
	}

	for _, item := range items {
		// NODELETING guard runs in BOTH sweepAll + normal mode.
		if protectedSweepVpcName(item.Name) {
			log.Printf("[INFO] Skipping WAF v3 defense template (NODELETING): %s", item.Name)
			continue
		}
		if !sweepAll() && !wafv3SweepPrefixMatch(item.Name) {
			log.Printf("[INFO] Skipping WAF v3 defense template (no test prefix): %s", item.Name)
			continue
		}
		log.Printf("[INFO] Deleting WAF v3 defense template: %s (%s)", item.Name, item.TemplateId)
		deleteRequest := map[string]interface{}{
			"RegionId":   client.RegionId,
			"InstanceId": instanceId,
			"TemplateId": item.TemplateId,
		}
		if _, err := client.RpcPost("waf-openapi", "2021-10-01", "DeleteDefenseTemplate", nil, deleteRequest, false); err != nil {
			log.Printf("[ERROR] Failed to delete WAF v3 defense template (%s): %s", item.TemplateId, err)
		}
	}
	return nil
}

// testSweepWafv3DefenseResourceGroups sweeps alicloud_wafv3_defense_resource_group resources.
func testSweepWafv3DefenseResourceGroups(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return WrapError(err)
	}
	client := rawClient.(*connectivity.AliyunClient)

	instanceId := wafv3InstanceID(client)
	if instanceId == "" {
		return nil
	}

	groups := make([]string, 0)
	action := "DescribeDefenseResourceGroups"
	pageNumber := 1
	pageSize := PageSizeLarge
	for {
		request := map[string]interface{}{
			"RegionId":   client.RegionId,
			"InstanceId": instanceId,
			"PageNumber":  pageNumber,
			"PageSize":    pageSize,
		}
		response, err := client.RpcPost("waf-openapi", "2021-10-01", action, nil, request, false)
		if err != nil {
			log.Printf("[ERROR] Failed to retrieve WAF v3 defense resource groups: %s", err)
			return nil
		}
		resp, err := jsonpath.Get("$.Groups", response)
		if err != nil {
			return nil
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			item, ok := v.(map[string]interface{})
			if !ok {
				continue
			}
			groups = append(groups, fmt.Sprint(item["GroupName"]))
		}
		if len(result) < pageSize {
			break
		}
		pageNumber++
	}

	for _, groupName := range groups {
		// NODELETING guard runs in BOTH sweepAll + normal mode.
		if protectedSweepVpcName(groupName) {
			log.Printf("[INFO] Skipping WAF v3 defense resource group (NODELETING): %s", groupName)
			continue
		}
		if !sweepAll() && !wafv3SweepPrefixMatch(groupName) {
			log.Printf("[INFO] Skipping WAF v3 defense resource group (no test prefix): %s", groupName)
			continue
		}
		log.Printf("[INFO] Deleting WAF v3 defense resource group: %s", groupName)
		deleteRequest := map[string]interface{}{
			"RegionId":   client.RegionId,
			"InstanceId": instanceId,
			"GroupName":  groupName,
		}
		if _, err := client.RpcPost("waf-openapi", "2021-10-01", "DeleteDefenseResourceGroup", nil, deleteRequest, false); err != nil {
			log.Printf("[ERROR] Failed to delete WAF v3 defense resource group (%s): %s", groupName, err)
		}
	}
	return nil
}

// testSweepWafv3DefenseRules sweeps alicloud_wafv3_defense_rule resources.
// Address-book rules (DefenseScene="address_book") are skipped here — they are
// handled by testSweepWafv3AddressBooks to avoid duplicate deletion.
func testSweepWafv3DefenseRules(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return WrapError(err)
	}
	client := rawClient.(*connectivity.AliyunClient)

	instanceId := wafv3InstanceID(client)
	if instanceId == "" {
		return nil
	}

	type ruleItem struct {
		RuleId      string
		DefenseType string
		Name        string
	}
	items := make([]ruleItem, 0)
	action := "DescribeDefenseRules"
	pageNumber := 1
	pageSize := PageSizeLarge
	for {
		request := map[string]interface{}{
			"RegionId":   client.RegionId,
			"InstanceId": instanceId,
			"PageNumber":  pageNumber,
			"PageSize":    pageSize,
		}
		response, err := client.RpcPost("waf-openapi", "2021-10-01", action, nil, request, false)
		if err != nil {
			log.Printf("[ERROR] Failed to retrieve WAF v3 defense rules: %s", err)
			return nil
		}
		resp, err := jsonpath.Get("$.Rules", response)
		if err != nil {
			return nil
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			item, ok := v.(map[string]interface{})
			if !ok {
				continue
			}
			// Skip address-book rules — handled by testSweepWafv3AddressBooks.
			if fmt.Sprint(item["DefenseScene"]) == "address_book" {
				continue
			}
			items = append(items, ruleItem{
				RuleId:      fmt.Sprint(item["RuleId"]),
				DefenseType: fmt.Sprint(item["DefenseType"]),
				Name:        fmt.Sprint(item["RuleName"]),
			})
		}
		if len(result) < pageSize {
			break
		}
		pageNumber++
	}

	for _, item := range items {
		// NODELETING guard runs in BOTH sweepAll + normal mode.
		if protectedSweepVpcName(item.Name) {
			log.Printf("[INFO] Skipping WAF v3 defense rule (NODELETING): %s", item.Name)
			continue
		}
		if !sweepAll() && !wafv3SweepPrefixMatch(item.Name) {
			log.Printf("[INFO] Skipping WAF v3 defense rule (no test prefix): %s", item.Name)
			continue
		}
		log.Printf("[INFO] Deleting WAF v3 defense rule: %s (%s)", item.Name, item.RuleId)
		deleteRequest := map[string]interface{}{
			"RegionId":    client.RegionId,
			"InstanceId":  instanceId,
			"RuleIds":      item.RuleId,
			"DefenseType":  item.DefenseType,
		}
		if _, err := client.RpcPost("waf-openapi", "2021-10-01", "DeleteDefenseRule", nil, deleteRequest, false); err != nil {
			log.Printf("[ERROR] Failed to delete WAF v3 defense rule (%s): %s", item.RuleId, err)
		}
	}
	return nil
}

// testSweepWafv3AddressBooks sweeps alicloud_wafv3_address_book resources.
// Address books are defense rules with DefenseScene="address_book" and DefenseType="global";
// the address-book name is stored inside the rule's Config JSON (config.name).
func testSweepWafv3AddressBooks(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return WrapError(err)
	}
	client := rawClient.(*connectivity.AliyunClient)

	instanceId := wafv3InstanceID(client)
	if instanceId == "" {
		return nil
	}

	type bookItem struct {
		RuleId string
		Name   string
	}
	items := make([]bookItem, 0)
	queryFilter, _ := json.Marshal(map[string]interface{}{"scene": "address_book"})
	action := "DescribeDefenseRules"
	pageNumber := 1
	pageSize := PageSizeLarge
	for {
		request := map[string]interface{}{
			"RegionId":    client.RegionId,
			"InstanceId":  instanceId,
			"DefenseType": "global",
			"Query":       string(queryFilter),
			"PageNumber":  pageNumber,
			"PageSize":    pageSize,
		}
		response, err := client.RpcPost("waf-openapi", "2021-10-01", action, nil, request, false)
		if err != nil {
			log.Printf("[ERROR] Failed to retrieve WAF v3 address books: %s", err)
			return nil
		}
		resp, err := jsonpath.Get("$.Rules", response)
		if err != nil {
			return nil
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			item, ok := v.(map[string]interface{})
			if !ok {
				continue
			}
			bookName := ""
			if config := decodeWafv3AddressBookConfig(item["Config"]); config != nil {
				if n, ok := config["name"].(string); ok {
					bookName = n
				}
			}
			items = append(items, bookItem{
				RuleId: fmt.Sprint(item["RuleId"]),
				Name:   bookName,
			})
		}
		if len(result) < pageSize {
			break
		}
		pageNumber++
	}

	for _, item := range items {
		// NODELETING guard runs in BOTH sweepAll + normal mode.
		if protectedSweepVpcName(item.Name) {
			log.Printf("[INFO] Skipping WAF v3 address book (NODELETING): %s", item.Name)
			continue
		}
		if !sweepAll() && !wafv3SweepPrefixMatch(item.Name) {
			log.Printf("[INFO] Skipping WAF v3 address book (no test prefix): %s", item.Name)
			continue
		}
		log.Printf("[INFO] Deleting WAF v3 address book: %s (%s)", item.Name, item.RuleId)
		deleteRequest := map[string]interface{}{
			"RegionId":    client.RegionId,
			"InstanceId":  instanceId,
			"RuleIds":      item.RuleId,
			"DefenseType": "global",
		}
		if _, err := client.RpcPost("waf-openapi", "2021-10-01", "DeleteDefenseRule", nil, deleteRequest, false); err != nil {
			log.Printf("[ERROR] Failed to delete WAF v3 address book (%s): %s", item.RuleId, err)
		}
	}
	return nil
}
