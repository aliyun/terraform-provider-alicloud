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

func resourceAliCloudDfsAccessRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudDfsAccessRuleCreate,
		Read:   resourceAliCloudDfsAccessRuleRead,
		Update: resourceAliCloudDfsAccessRuleUpdate,
		Delete: resourceAliCloudDfsAccessRuleDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"access_group_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"access_rule_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"network_segment": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"priority": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: IntBetween(0, 100),
			},
			"rw_access_type": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: StringInSlice([]string{"RDWR", "RDONLY"}, true),
			},
		},
	}
}

func resourceAliCloudDfsAccessRuleCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "CreateAccessRule"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	query["AccessGroupId"] = d.Get("access_group_id")
	request["InputRegionId"] = client.RegionId

	request["NetworkSegment"] = d.Get("network_segment")
	request["Priority"] = d.Get("priority")
	if v, ok := d.GetOk("description"); ok {
		request["Description"] = v
	}
	request["RWAccessType"] = d.Get("rw_access_type")
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RpcPost("DFS", "2018-06-20", action, query, request, false)

		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, response, request)
		return nil
	})

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_dfs_access_rule", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprintf("%v:%v", query["AccessGroupId"], response["AccessRuleId"]))

	return resourceAliCloudDfsAccessRuleRead(d, meta)
}

func resourceAliCloudDfsAccessRuleRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	dfsServiceV2 := DfsServiceV2{client}

	objectRaw, err := dfsServiceV2.DescribeDfsAccessRule(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_dfs_access_rule DescribeDfsAccessRule Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("create_time", objectRaw["CreateTime"])
	d.Set("description", objectRaw["Description"])
	d.Set("network_segment", objectRaw["NetworkSegment"])
	d.Set("priority", objectRaw["Priority"])
	d.Set("rw_access_type", objectRaw["RWAccessType"])
	d.Set("access_group_id", objectRaw["AccessGroupId"])
	d.Set("access_rule_id", objectRaw["AccessRuleId"])

	parts := strings.Split(d.Id(), ":")
	d.Set("access_group_id", parts[0])
	d.Set("access_rule_id", parts[1])

	return nil
}

func resourceAliCloudDfsAccessRuleUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	update := false
	parts := strings.Split(d.Id(), ":")
	action := "ModifyAccessRule"
	var err error
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	query["AccessGroupId"] = parts[0]
	query["AccessRuleId"] = parts[1]
	request["InputRegionId"] = client.RegionId
	if d.HasChange("priority") {
		update = true
	}
	request["Priority"] = d.Get("priority")
	if d.HasChange("description") {
		update = true
		request["Description"] = d.Get("description")
	}

	if d.HasChange("rw_access_type") {
		update = true
	}
	request["RWAccessType"] = d.Get("rw_access_type")
	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("DFS", "2018-06-20", action, query, request, false)

			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			addDebug(action, response, request)
			return nil
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
	}

	return resourceAliCloudDfsAccessRuleRead(d, meta)
}

func resourceAliCloudDfsAccessRuleDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	parts := strings.Split(d.Id(), ":")
	action := "DeleteAccessRule"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	query["AccessGroupId"] = parts[0]
	query["AccessRuleId"] = parts[1]
	request["InputRegionId"] = client.RegionId

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.RpcPost("DFS", "2018-06-20", action, query, request, false)

		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, response, request)
		return nil
	})

	if err != nil {
		if IsExpectedErrors(err, []string{"InvalidParameter.AccessRuleNotFound"}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	return nil
}
