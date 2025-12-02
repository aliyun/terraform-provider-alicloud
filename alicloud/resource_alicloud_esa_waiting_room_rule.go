// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudEsaWaitingRoomRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudEsaWaitingRoomRuleCreate,
		Read:   resourceAliCloudEsaWaitingRoomRuleRead,
		Update: resourceAliCloudEsaWaitingRoomRuleUpdate,
		Delete: resourceAliCloudEsaWaitingRoomRuleDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"rule": {
				Type:     schema.TypeString,
				Required: true,
			},
			"rule_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"site_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"status": {
				Type:     schema.TypeString,
				Required: true,
			},
			"waiting_room_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"waiting_room_rule_id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func resourceAliCloudEsaWaitingRoomRuleCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "CreateWaitingRoomRule"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	if v, ok := d.GetOk("site_id"); ok {
		request["SiteId"] = v
	}
	if v, ok := d.GetOk("waiting_room_id"); ok {
		request["WaitingRoomId"] = v
	}

	request["RuleEnable"] = d.Get("status")
	request["RuleName"] = d.Get("rule_name")
	request["Rule"] = d.Get("rule")
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RpcPost("ESA", "2024-09-10", action, query, request, true)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_esa_waiting_room_rule", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprintf("%v:%v:%v", request["SiteId"], request["WaitingRoomId"], response["WaitingRoomRuleId"]))

	return resourceAliCloudEsaWaitingRoomRuleRead(d, meta)
}

func resourceAliCloudEsaWaitingRoomRuleRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	esaServiceV2 := EsaServiceV2{client}

	objectRaw, err := esaServiceV2.DescribeEsaWaitingRoomRule(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_esa_waiting_room_rule DescribeEsaWaitingRoomRule Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("rule", objectRaw["Rule"])
	d.Set("rule_name", objectRaw["RuleName"])
	d.Set("status", objectRaw["RuleEnable"])
	d.Set("waiting_room_rule_id", objectRaw["WaitingRoomRuleId"])

	parts := strings.Split(d.Id(), ":")
	d.Set("site_id", fmt.Sprint(parts[0]))
	d.Set("waiting_room_id", parts[1])

	return nil
}

func resourceAliCloudEsaWaitingRoomRuleUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	update := false

	var err error
	parts := strings.Split(d.Id(), ":")
	action := "UpdateWaitingRoomRule"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["SiteId"] = parts[0]
	request["WaitingRoomRuleId"] = parts[2]
	request["RegionId"] = client.RegionId
	if d.HasChange("status") {
		update = true
	}
	request["RuleEnable"] = d.Get("status")
	if d.HasChange("rule_name") {
		update = true
	}
	request["RuleName"] = d.Get("rule_name")
	if d.HasChange("rule") {
		update = true
	}
	request["Rule"] = d.Get("rule")
	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("ESA", "2024-09-10", action, query, request, true)
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
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
	}

	return resourceAliCloudEsaWaitingRoomRuleRead(d, meta)
}

func resourceAliCloudEsaWaitingRoomRuleDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	parts := strings.Split(d.Id(), ":")
	action := "DeleteWaitingRoomRule"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["SiteId"] = parts[0]
	request["WaitingRoomRuleId"] = parts[2]

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.RpcPost("ESA", "2024-09-10", action, query, request, true)
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
		if NotFoundError(err) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	return nil
}
