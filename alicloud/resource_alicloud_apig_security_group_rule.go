// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudApigSecurityGroupRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudApigSecurityGroupRuleCreate,
		Read:   resourceAliCloudApigSecurityGroupRuleRead,
		Delete: resourceAliCloudApigSecurityGroupRuleDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"gateway_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"security_group_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"port_ranges": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"security_group_rule_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"port_range": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ip_protocol": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"auth_cidrs": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"source_security_group_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"security_group_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"region_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAliCloudApigSecurityGroupRuleCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	gatewayId := d.Get("gateway_id").(string)
	action := fmt.Sprintf("/v1/gateways/%s/security-group-rules", gatewayId)
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]*string)
	body := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})

	request["securityGroupId"] = d.Get("security_group_id")
	if v, ok := d.GetOk("port_ranges"); ok {
		portRanges := make([]string, 0, len(v.([]interface{})))
		for _, vv := range v.([]interface{}) {
			if vv == nil {
				continue
			}
			portRanges = append(portRanges, vv.(string))
		}
		request["portRanges"] = portRanges
	}
	if v, ok := d.GetOk("description"); ok {
		request["description"] = v
	}
	body = request
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RoaPost("APIG", "2024-03-27", action, query, nil, body, true)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_apig_security_group_rule", action, AlibabaCloudSdkGoERROR)
	}

	// The Create API response only returns code/requestId/success and does not include
	// the created security group rule id, so query the list endpoint to locate the
	// newly created rule (matched by description, or the first item when description is
	// absent) and extract its securityGroupRuleId.
	descriptionVal := d.Get("description").(string)
	listAction := fmt.Sprintf("/v1/gateways/%s/authorized-security-groups-rules", gatewayId)
	var listResponse map[string]interface{}
	listQuery := make(map[string]*string)
	var ruleId string
	listWait := incrementalWait(3*time.Second, 5*time.Second)
	var listErr error
	listErr = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		listResponse, listErr = client.RoaGet("APIG", "2024-03-27", listAction, listQuery, nil, nil)
		if listErr != nil {
			if NeedRetry(listErr) {
				listWait()
				return resource.RetryableError(listErr)
			}
			return resource.NonRetryableError(listErr)
		}
		itemsRaw, gerr := jsonpath.Get("$.data.items", listResponse)
		if gerr != nil {
			return resource.NonRetryableError(gerr)
		}
		items, ok := itemsRaw.([]interface{})
		if !ok || len(items) == 0 {
			listWait()
			return resource.RetryableError(fmt.Errorf("the created security group rule has not appeared in the list yet"))
		}
		for _, raw := range items {
			item, okItem := raw.(map[string]interface{})
			if !okItem {
				continue
			}
			if descriptionVal != "" {
				if fmt.Sprint(item["description"]) == descriptionVal {
					ruleId = fmt.Sprint(item["securityGroupRuleId"])
					return nil
				}
				continue
			}
			ruleId = fmt.Sprint(item["securityGroupRuleId"])
			return nil
		}
		listWait()
		return resource.RetryableError(fmt.Errorf("the created security group rule matched by description has not appeared in the list yet"))
	})
	addDebug(listAction, listResponse, request)
	if listErr != nil {
		return WrapErrorf(listErr, DefaultErrorMsg, "alicloud_apig_security_group_rule", listAction, AlibabaCloudSdkGoERROR)
	}
	if ruleId == "" {
		return WrapError(fmt.Errorf("failed to locate the created security group rule id for gateway %s", gatewayId))
	}
	d.SetId(fmt.Sprintf("%v:%v", gatewayId, ruleId))

	return resourceAliCloudApigSecurityGroupRuleRead(d, meta)
}

func resourceAliCloudApigSecurityGroupRuleRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	apigServiceV2 := ApigServiceV2{client}

	objectRaw, err := apigServiceV2.DescribeApigSecurityGroupRule(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_apig_security_group_rule DescribeApigSecurityGroupRule Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("security_group_rule_id", objectRaw["securityGroupRuleId"])
	d.Set("security_group_id", objectRaw["securityGroupId"])
	d.Set("port_range", objectRaw["portRange"])
	d.Set("ip_protocol", objectRaw["ipProtocol"])
	d.Set("description", objectRaw["description"])
	d.Set("source_security_group_id", objectRaw["sourceSecurityGroupId"])
	d.Set("security_group_name", objectRaw["securityGroupName"])
	d.Set("vpc_id", objectRaw["vpcId"])
	// auth_cidrs is a Computed TypeList; the list endpoint omits authCidrs for
	// type:SecurityGroup rules, so always write a list (empty when absent) to keep
	// state populated and avoid a spurious "" => "<computed>" diff.
	if authCidrs, ok := objectRaw["authCidrs"].([]interface{}); ok {
		d.Set("auth_cidrs", authCidrs)
	} else {
		d.Set("auth_cidrs", []interface{}{})
	}
	if parts := strings.Split(d.Id(), ":"); len(parts) == 2 {
		d.Set("gateway_id", parts[0])
	}
	d.Set("region_id", client.RegionId)

	return nil
}

func resourceAliCloudApigSecurityGroupRuleDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	parts := strings.Split(d.Id(), ":")
	if len(parts) != 2 {
		return WrapError(fmt.Errorf("invalid Resource Id %s. Expected parts' length %d, got %d", d.Id(), 2, len(parts)))
	}
	gatewayId := parts[0]
	securityGroupRuleId := parts[1]
	action := fmt.Sprintf("/v1/gateways/%s/security-group-rules/%s", gatewayId, securityGroupRuleId)
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]*string)
	var err error
	request = make(map[string]interface{})

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.RoaDelete("APIG", "2024-03-27", action, query, nil, nil, true)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)

	if err != nil {
		if IsExpectedErrors(err, []string{"NotFound.SecurityGroupRule"}) || NotFoundError(err) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	return nil
}
