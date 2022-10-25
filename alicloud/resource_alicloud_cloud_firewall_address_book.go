package alicloud

import (
	"fmt"
	"log"
	"strings"
	"time"

	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceAlicloudCloudFirewallAddressBook() *schema.Resource {

	return &schema.Resource{
		Create: resourceAlicloudCloudFirewallAddressBookCreate,
		Read:   resourceAlicloudCloudFirewallAddressBookRead,
		Update: resourceAlicloudCloudFirewallAddressBookUpdate,
		Delete: resourceAlicloudCloudFirewallAddressBookDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"address_list": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"auto_add_tag_ecs": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntInSlice([]int{0, 1}),
			},
			"description": {
				Type:     schema.TypeString,
				Required: true,
			},
			"group_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"group_type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"ip", "tag"}, false),
			},
			"lang": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"ecs_tags": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"tag_key": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"tag_value": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"tag_relation": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"and", "or"}, false),
			},
		},
	}
}

func resourceAlicloudCloudFirewallAddressBookCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "AddAddressBook"
	request := make(map[string]interface{})
	conn, err := client.NewCloudfwClient()
	if err != nil {
		return WrapError(err)
	}
	if v, ok := d.GetOk("address_list"); ok {
		request["AddressList"] = strings.Join(expandStringList(v.([]interface{})), ",")
	}
	if v, ok := d.GetOk("auto_add_tag_ecs"); ok {
		request["AutoAddTagEcs"] = v
	}
	request["Description"] = d.Get("description")
	request["GroupName"] = d.Get("group_name")
	request["GroupType"] = d.Get("group_type")
	if v, ok := d.GetOk("lang"); ok {
		request["Lang"] = v
	}

	if v, ok := d.GetOk("tag_relation"); ok {
		request["TagRelation"] = v
	}
	if v, ok := d.GetOk("ecs_tags"); ok {
		for i, tagItem := range v.(*schema.Set).List() {
			tagItemArg := tagItem.(map[string]interface{})
			request[fmt.Sprintf("TagList.%d.TagValue", i+1)] = tagItemArg["tag_value"]
			request[fmt.Sprintf("TagList.%d.TagKey", i+1)] = tagItemArg["tag_key"]
		}
	}

	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2017-12-07"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_cloud_firewall_address_book", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["GroupUuid"]))

	return resourceAlicloudCloudFirewallAddressBookRead(d, meta)
}
func resourceAlicloudCloudFirewallAddressBookRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cloudfwService := CloudfwService{client}
	object, err := cloudfwService.DescribeCloudFirewallAddressBook(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_cloud_firewall_address_book cloudfwService.DescribeCloudFirewallAddressBook Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	addressListItems := make([]string, 0)
	for _, addressListArg := range object["AddressList"].([]interface{}) {
		addressListItems = append(addressListItems, fmt.Sprint(addressListArg))
	}
	d.Set("address_list", addressListItems)
	if v, ok := object["AutoAddTagEcs"]; ok {
		d.Set("auto_add_tag_ecs", formatInt(v))
	}
	d.Set("description", object["Description"])
	d.Set("group_name", object["GroupName"])
	d.Set("group_type", object["GroupType"])
	d.Set("tag_relation", object["TagRelation"])
	ecsTags := make([]map[string]interface{}, 0)

	for _, tagListItem := range object["TagList"].([]interface{}) {
		ecsTagItem := make(map[string]interface{})
		ecsTagItem["tag_value"] = tagListItem.(map[string]interface{})["TagValue"]
		ecsTagItem["tag_key"] = tagListItem.(map[string]interface{})["TagKey"]
		ecsTags = append(ecsTags, ecsTagItem)
	}

	d.Set("ecs_tags", ecsTags)
	return nil
}
func resourceAlicloudCloudFirewallAddressBookUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	update := false
	request := map[string]interface{}{
		"GroupUuid": d.Id(),
	}
	if d.HasChange("description") {
		update = true
	}
	request["Description"] = d.Get("description")
	if d.HasChange("group_name") {
		update = true
	}
	request["GroupName"] = d.Get("group_name")
	if d.HasChange("address_list") {
		update = true
		if v, ok := d.GetOk("address_list"); ok {
			request["AddressList"] = strings.Join(expandStringList(v.([]interface{})), ",")
		}
	}
	if d.HasChange("auto_add_tag_ecs") {
		update = true
		if v, ok := d.GetOkExists("auto_add_tag_ecs"); ok {
			request["AutoAddTagEcs"] = v
		}
	}
	if d.HasChange("tag_relation") {
		update = true
		if v, ok := d.GetOk("tag_relation"); ok {
			request["TagRelation"] = v
		}
	}
	if d.HasChange("ecs_tags") {
		update = true
		if v, ok := d.GetOk("ecs_tags"); ok {
			for i, tagItem := range v.(*schema.Set).List() {
				tagItemArg := tagItem.(map[string]interface{})
				request[fmt.Sprintf("TagList.%d.TagValue", i+1)] = tagItemArg["tag_value"]
				request[fmt.Sprintf("TagList.%d.TagKey", i+1)] = tagItemArg["tag_key"]
			}
		}
	}

	if update {
		if v, ok := d.GetOk("lang"); ok {
			request["Lang"] = v
		}
		if v, ok := d.GetOk("source_ip"); ok {
			request["SourceIp"] = v
		}

		action := "ModifyAddressBook"
		conn, err := client.NewCloudfwClient()
		if err != nil {
			return WrapError(err)
		}
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2017-12-07"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
	return resourceAlicloudCloudFirewallAddressBookRead(d, meta)
}
func resourceAlicloudCloudFirewallAddressBookDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "DeleteAddressBook"
	var response map[string]interface{}
	conn, err := client.NewCloudfwClient()
	if err != nil {
		return WrapError(err)
	}
	request := map[string]interface{}{
		"GroupUuid": d.Id(),
	}

	if v, ok := d.GetOk("lang"); ok {
		request["Lang"] = v
	}
	if v, ok := d.GetOk("source_ip"); ok {
		request["SourceIp"] = v
	}
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2017-12-07"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
	return nil
}
