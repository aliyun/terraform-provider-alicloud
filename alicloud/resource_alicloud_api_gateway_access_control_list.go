// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"log"
	"regexp"
	"time"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudApiGatewayAccessControlList() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudApiGatewayAccessControlListCreate,
		Read:   resourceAliCloudApiGatewayAccessControlListRead,
		Update: resourceAliCloudApiGatewayAccessControlListUpdate,
		Delete: resourceAliCloudApiGatewayAccessControlListDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"access_control_list_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: StringMatch(regexp.MustCompile("^[\u4E00-\u9FA5A-Za-z0-9_-]+$"), "Access control list name"),
			},
			"acl_entrys": {
				Type:       schema.TypeSet,
				Optional:   true,
				Sensitive:  true,
				Computed:   true,
				Deprecated: "Field 'acl_entrys' has been deprecated from provider version v1.228.0, and it will be removed in the future version. Please use the new resource 'alicloud_api_gateway_acl_entry_attachment'.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"acl_entry_comment": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"acl_entry_ip": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"address_ip_version": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				ValidateFunc: StringInSlice([]string{"ipv4", "ipv6"}, false),
			},
		},
	}
}

func resourceAliCloudApiGatewayAccessControlListCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "CreateAccessControlList"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})

	request["AclName"] = d.Get("access_control_list_name")
	if v, ok := d.GetOk("address_ip_version"); ok {
		request["AddressIPVersion"] = v
	}
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RpcPost("CloudAPI", "2016-07-14", action, query, request, true)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_api_gateway_access_control_list", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["AclId"]))

	return resourceAliCloudApiGatewayAccessControlListUpdate(d, meta)
}

func resourceAliCloudApiGatewayAccessControlListRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	apiGatewayServiceV2 := ApiGatewayServiceV2{client}

	objectRaw, err := apiGatewayServiceV2.DescribeApiGatewayAccessControlList(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_api_gateway_access_control_list DescribeApiGatewayAccessControlList Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("access_control_list_name", objectRaw["AclName"])
	d.Set("address_ip_version", objectRaw["AddressIPVersion"])

	aclEntry1Raw, _ := jsonpath.Get("$.AclEntrys.AclEntry", objectRaw)
	aclEntrysMaps := make([]map[string]interface{}, 0)
	if aclEntry1Raw != nil {
		for _, aclEntryChild1Raw := range aclEntry1Raw.([]interface{}) {
			aclEntrysMap := make(map[string]interface{})
			aclEntryChild1Raw := aclEntryChild1Raw.(map[string]interface{})
			aclEntrysMap["acl_entry_comment"] = aclEntryChild1Raw["AclEntryComment"]
			aclEntrysMap["acl_entry_ip"] = aclEntryChild1Raw["AclEntryIp"]

			aclEntrysMaps = append(aclEntrysMaps, aclEntrysMap)
		}
	}
	d.Set("acl_entrys", aclEntrysMaps)

	return nil
}

func resourceAliCloudApiGatewayAccessControlListUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	var err error
	if d.HasChange("acl_entrys") {
		oldEntry, newEntry := d.GetChange("acl_entrys")
		oldEntrySet := oldEntry.(*schema.Set)
		newEntrySet := newEntry.(*schema.Set)
		removed := oldEntrySet.Difference(newEntrySet)
		added := newEntrySet.Difference(oldEntrySet)

		if removed.Len() > 0 {
			action := "RemoveAccessControlListEntry"
			request = make(map[string]interface{})
			query = make(map[string]interface{})
			query["AclId"] = d.Id()
			localData := removed.List()
			aclEntrysMaps := make([]interface{}, 0)
			for _, dataLoop := range localData {
				dataLoopTmp := dataLoop.(map[string]interface{})
				dataLoopMap := make(map[string]interface{})
				dataLoopMap["Entry"] = dataLoopTmp["acl_entry_ip"]
				dataLoopMap["Comment"] = dataLoopTmp["acl_entry_comment"]
				aclEntrysMaps = append(aclEntrysMaps, dataLoopMap)
			}
			request["AclEntrys"], _ = convertArrayObjectToJsonString(aclEntrysMaps)

			runtime := util.RuntimeOptions{}
			runtime.SetAutoretry(true)
			wait := incrementalWait(3*time.Second, 5*time.Second)
			err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
				response, err = client.RpcPost("CloudAPI", "2016-07-14", action, query, request, true)
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

		if added.Len() > 0 {
			action := "AddAccessControlListEntry"
			request = make(map[string]interface{})
			query = make(map[string]interface{})
			query["AclId"] = d.Id()
			localData := added.List()
			aclEntrysMaps := make([]interface{}, 0)
			for _, dataLoop := range localData {
				dataLoopTmp := dataLoop.(map[string]interface{})
				dataLoopMap := make(map[string]interface{})
				dataLoopMap["Entry"] = dataLoopTmp["acl_entry_ip"]
				dataLoopMap["Comment"] = dataLoopTmp["acl_entry_comment"]
				aclEntrysMaps = append(aclEntrysMaps, dataLoopMap)
			}
			request["AclEntrys"], _ = convertArrayObjectToJsonString(aclEntrysMaps)

			runtime := util.RuntimeOptions{}
			runtime.SetAutoretry(true)
			wait := incrementalWait(3*time.Second, 5*time.Second)
			err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
				response, err = client.RpcPost("CloudAPI", "2016-07-14", action, query, request, true)
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

	}
	return resourceAliCloudApiGatewayAccessControlListRead(d, meta)
}

func resourceAliCloudApiGatewayAccessControlListDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	action := "DeleteAccessControlList"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	query["AclId"] = d.Id()

	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.RpcPost("CloudAPI", "2016-07-14", action, query, request, true)

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
		if IsExpectedErrors(err, []string{"NotFoundAccessControlList"}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	return nil
}
