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

func resourceAliCloudApigGatewaySecurityGroupRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudApigGatewaySecurityGroupRuleCreate,
		Read:   resourceAliCloudApigGatewaySecurityGroupRuleRead,
		Delete: resourceAliCloudApigGatewaySecurityGroupRuleDelete,
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
			"port_range": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
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
			"source_security_group_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ip_protocol": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAliCloudApigGatewaySecurityGroupRuleCreate(d *schema.ResourceData, meta interface{}) error {
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
	request["portRanges"] = []string{d.Get("port_range").(string)}
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_apig_gateway_security_group_rule", action, AlibabaCloudSdkGoERROR)
	}

	// AddGatewaySecurityGroupRule does not document a securityGroupRuleId in its response
	// (only requestId/code/message). Try common response paths first; if absent, fall back
	// to listing authorized rules and matching by securityGroupId + description.
	securityGroupRuleId := ""
	if id, e := jsonpath.Get("$.data.securityGroupRuleId", response); e == nil && id != nil {
		securityGroupRuleId = fmt.Sprint(id)
	}
	if securityGroupRuleId == "" {
		if id, e := jsonpath.Get("$.securityGroupRuleId", response); e == nil && id != nil {
			securityGroupRuleId = fmt.Sprint(id)
		}
	}
	if securityGroupRuleId == "" {
		securityGroupId := d.Get("security_group_id").(string)
		description := d.Get("description").(string)
		portRangeHint := d.Get("port_range").(string)
		listWait := incrementalWait(2*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
			items, listErr := (&ApigServiceV2{client}).ListApigGatewaySecurityGroupRules(gatewayId)
			if listErr != nil {
				return resource.NonRetryableError(listErr)
			}
			for _, item := range items {
				if fmt.Sprint(item["securityGroupId"]) != securityGroupId {
					continue
				}
				if description != "" && fmt.Sprint(item["description"]) != description {
					continue
				}
				if description == "" && portRangeHint != "" {
					pr := fmt.Sprint(item["portRange"])
					if pr != "" && !strings.Contains(pr, portRangeHint) {
						continue
					}
				}
				securityGroupRuleId = fmt.Sprint(item["securityGroupRuleId"])
				return nil
			}
			listWait()
			return resource.RetryableError(fmt.Errorf("the security group rule is not visible in the authorized list yet"))
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, "alicloud_apig_gateway_security_group_rule", action, AlibabaCloudSdkGoERROR)
		}
	}

	d.SetId(fmt.Sprintf("%s:%s", gatewayId, securityGroupRuleId))

	return resourceAliCloudApigGatewaySecurityGroupRuleRead(d, meta)
}

func resourceAliCloudApigGatewaySecurityGroupRuleRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	apigServiceV2 := ApigServiceV2{client}

	object, err := apigServiceV2.DescribeApigGatewaySecurityGroupRule(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_apig_gateway_security_group_rule DescribeApigGatewaySecurityGroupRule Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("security_group_rule_id", object["securityGroupRuleId"])
	d.Set("security_group_id", object["securityGroupId"])
	d.Set("port_range", object["portRange"])
	d.Set("description", object["description"])
	d.Set("source_security_group_id", object["sourceSecurityGroupId"])
	d.Set("ip_protocol", object["ipProtocol"])

	parts := strings.Split(d.Id(), ":")
	if len(parts) == 2 {
		d.Set("gateway_id", parts[0])
	}

	return nil
}

func resourceAliCloudApigGatewaySecurityGroupRuleDelete(d *schema.ResourceData, meta interface{}) error {
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
		// When the parent gateway has already been deleted out-of-band, the API
		// returns Conflict.GatewayIsDeleted (409); the rule is gone with it, so
		// treat it as a successful delete.
		if NotFoundError(err) || IsExpectedErrors(err, []string{"Conflict.GatewayIsDeleted"}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	return nil
}
