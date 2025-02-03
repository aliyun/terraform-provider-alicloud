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

func resourceAliCloudNasAccessRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudNasAccessRuleCreate,
		Read:   resourceAliCloudNasAccessRuleRead,
		Update: resourceAliCloudNasAccessRuleUpdate,
		Delete: resourceAliCloudNasAccessRuleDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"access_group_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"access_rule_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"file_system_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Computed:     true,
				ValidateFunc: StringInSlice([]string{"standard", "extreme"}, true),
			},
			"ipv6_source_cidr_ip": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"priority": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      1,
				ValidateFunc: IntBetween(0, 100),
			},
			"rw_access_type": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: StringInSlice([]string{"RDWR", "RDONLY"}, true),
			},
			"source_cidr_ip": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"user_access_type": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: StringInSlice([]string{"no_squash", "root_squash", "all_squash"}, true),
			},
		},
	}
}

func resourceAliCloudNasAccessRuleCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "CreateAccessRule"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	if v, ok := d.GetOk("file_system_type"); ok {
		query["FileSystemType"] = v
	} else {
		query["FileSystemType"] = "standard"
	}
	query["AccessGroupName"] = d.Get("access_group_name")

	if v, ok := d.GetOk("source_cidr_ip"); ok {
		request["SourceCidrIp"] = v
	}
	if v, ok := d.GetOk("ipv6_source_cidr_ip"); ok {
		request["Ipv6SourceCidrIp"] = v
	}
	if v, ok := d.GetOk("rw_access_type"); ok {
		request["RWAccessType"] = v
	}
	if v, ok := d.GetOk("user_access_type"); ok {
		request["UserAccessType"] = v
	}
	if v, ok := d.GetOk("priority"); ok {
		request["Priority"] = v
	}
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RpcPost("NAS", "2017-06-26", action, query, request, true)

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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_nas_access_rule", action, AlibabaCloudSdkGoERROR)
	}

	if query["FileSystemType"] == "standard" {
		d.SetId(fmt.Sprintf("%v:%v", query["AccessGroupName"], response["AccessRuleId"]))
	} else {
		d.SetId(fmt.Sprintf("%v:%v:%v", query["AccessGroupName"], response["AccessRuleId"], query["FileSystemType"]))
	}

	return resourceAliCloudNasAccessRuleRead(d, meta)
}

func resourceAliCloudNasAccessRuleRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	nasServiceV2 := NasServiceV2{client}

	objectRaw, err := nasServiceV2.DescribeNasAccessRule(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_nas_access_rule DescribeNasAccessRule Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("ipv6_source_cidr_ip", objectRaw["Ipv6SourceCidrIp"])
	d.Set("priority", formatInt(objectRaw["Priority"]))
	d.Set("rw_access_type", objectRaw["RWAccess"])
	d.Set("source_cidr_ip", objectRaw["SourceCidrIp"])
	d.Set("user_access_type", objectRaw["UserAccess"])
	d.Set("access_group_name", objectRaw["AccessGroupName"])
	d.Set("access_rule_id", objectRaw["AccessRuleId"])
	d.Set("file_system_type", objectRaw["FileSystemType"])

	parts := strings.Split(d.Id(), ":")
	d.Set("access_group_name", parts[0])
	d.Set("access_rule_id", parts[1])

	return nil
}

func resourceAliCloudNasAccessRuleUpdate(d *schema.ResourceData, meta interface{}) error {
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
	query["AccessGroupName"] = parts[0]
	query["AccessRuleId"] = parts[1]
	if len(parts) == 3 {
		query["FileSystemType"] = parts[2]
	}
	if d.HasChange("source_cidr_ip") {
		update = true
	}
	if v, ok := d.GetOk("source_cidr_ip"); ok {
		request["SourceCidrIp"] = v
	}

	if d.HasChange("priority") {
		update = true
	}
	request["Priority"] = d.Get("priority")

	if d.HasChange("rw_access_type") {
		update = true
	}
	request["RWAccessType"] = d.Get("rw_access_type")

	if d.HasChange("user_access_type") {
		update = true
	}
	request["UserAccessType"] = d.Get("user_access_type")

	if d.HasChange("ipv6_source_cidr_ip") {
		update = true
	}
	if v, ok := d.GetOk("ipv6_source_cidr_ip"); ok {
		request["Ipv6SourceCidrIp"] = v
	}

	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("NAS", "2017-06-26", action, query, request, true)

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

	return resourceAliCloudNasAccessRuleRead(d, meta)
}

func resourceAliCloudNasAccessRuleDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	parts := strings.Split(d.Id(), ":")
	action := "DeleteAccessRule"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	query["AccessGroupName"] = parts[0]
	query["AccessRuleId"] = parts[1]
	if len(parts) == 3 {
		query["FileSystemType"] = parts[2]
	}

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.RpcPost("NAS", "2017-06-26", action, query, request, true)

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

	return nil
}
