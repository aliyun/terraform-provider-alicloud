// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudWafv3AddressBook() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudWafv3AddressBookCreate,
		Read:   resourceAliCloudWafv3AddressBookRead,
		Update: resourceAliCloudWafv3AddressBookUpdate,
		Delete: resourceAliCloudWafv3AddressBookDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"address_book_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"address_book_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"address_book_type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: StringInSlice([]string{"ip"}, false),
			},
			"address_list": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceAliCloudWafv3AddressBookCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "CreateDefenseRule"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	if v, ok := d.GetOk("instance_id"); ok {
		request["InstanceId"] = v
	}
	request["RegionId"] = client.RegionId

	rulesDataList := make(map[string]interface{})

	if v, ok := d.GetOk("address_book_type"); ok {
		rulesDataList["valueType"] = v
	}

	if v, ok := d.GetOk("description"); ok {
		rulesDataList["description"] = v
	}

	if v, ok := d.GetOk("address_book_name"); ok {
		rulesDataList["name"] = v
	}

	RulesMap := make([]interface{}, 0)
	RulesMap = append(RulesMap, rulesDataList)
	rulesJson, err := json.Marshal(RulesMap)
	if err != nil {
		return WrapError(err)
	}
	request["Rules"] = string(rulesJson)

	request["DefenseScene"] = "address_book"
	request["DefenseType"] = "global"
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RpcPost("waf-openapi", "2021-10-01", action, query, request, true)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_wafv3_address_book", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprintf("%v:%v", request["InstanceId"], response["RuleIds"]))

	return resourceAliCloudWafv3AddressBookUpdate(d, meta)
}

func resourceAliCloudWafv3AddressBookRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	wafv3ServiceV2 := Wafv3ServiceV2{client}

	objectRaw, err := wafv3ServiceV2.DescribeWafv3AddressBook(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_wafv3_address_book DescribeWafv3AddressBook Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("address_book_name", objectRaw["name"])
	d.Set("address_book_type", objectRaw["valueType"])
	d.Set("description", objectRaw["description"])

	addressesRaw, err := wafv3ServiceV2.DescribeWafv3AddressBookAddresses(d.Id())
	if err != nil && !NotFoundError(err) {
		return WrapError(err)
	}

	addressListRawObj, _ := jsonpath.Get("$.AddressList", addressesRaw)
	addressListRaw := make([]interface{}, 0)
	if addressListRawObj != nil {
		addressListRaw = convertToInterfaceArray(addressListRawObj)
	}

	addresses := make([]interface{}, 0)
	for _, item := range addressListRaw {
		entry, ok := item.(map[string]interface{})
		if !ok {
			continue
		}
		if addr, ok := entry["Address"]; ok && addr != nil {
			addresses = append(addresses, addr)
		}
	}

	d.Set("address_list", addresses)

	parts := strings.Split(d.Id(), ":")
	d.Set("instance_id", parts[0])
	d.Set("address_book_id", parts[1])

	return nil
}

func resourceAliCloudWafv3AddressBookUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	parts := strings.Split(d.Id(), ":")

	if d.HasChange("description") || d.HasChange("address_book_name") {
		action := "ModifyDefenseRule"
		request = make(map[string]interface{})
		query = make(map[string]interface{})
		request["InstanceId"] = parts[0]
		request["RegionId"] = client.RegionId
		request["DefenseScene"] = "address_book"
		request["DefenseType"] = "global"

		rulesDataList := map[string]interface{}{
			"id":        parts[1],
			"valueType": d.Get("address_book_type"),
		}
		if v, ok := d.GetOk("description"); ok {
			rulesDataList["description"] = v
		}
		if v, ok := d.GetOk("address_book_name"); ok {
			rulesDataList["name"] = v
		}

		rulesJson, err := json.Marshal([]interface{}{rulesDataList})
		if err != nil {
			return WrapError(err)
		}
		request["Rules"] = string(rulesJson)

		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("waf-openapi", "2021-10-01", action, query, request, true)
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

	if d.HasChange("address_list") {
		var err error
		oldEntry, newEntry := d.GetChange("address_list")
		oldEntrySet := oldEntry.(*schema.Set)
		newEntrySet := newEntry.(*schema.Set)
		removed := oldEntrySet.Difference(newEntrySet)
		added := newEntrySet.Difference(oldEntrySet)

		if removed.Len() > 0 {
			parts := strings.Split(d.Id(), ":")
			action := "DeleteAddress"
			request = make(map[string]interface{})
			query = make(map[string]interface{})
			request["RuleId"] = parts[1]
			request["InstanceId"] = parts[0]

			localData := removed.List()
			addressListMapsArray := localData
			request["AddressList"] = addressListMapsArray

			wait := incrementalWait(3*time.Second, 5*time.Second)
			err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
				response, err = client.RpcPost("waf-openapi", "2021-10-01", action, query, request, true)
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

		if added.Len() > 0 {
			parts := strings.Split(d.Id(), ":")
			action := "AddAddress"
			request = make(map[string]interface{})
			query = make(map[string]interface{})
			request["RuleId"] = parts[1]
			request["InstanceId"] = parts[0]

			localData := added.List()
			addressListMapsArray := localData
			request["AddressList"] = addressListMapsArray

			wait := incrementalWait(3*time.Second, 5*time.Second)
			err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
				response, err = client.RpcPost("waf-openapi", "2021-10-01", action, query, request, true)
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

	}
	return resourceAliCloudWafv3AddressBookRead(d, meta)
}

func resourceAliCloudWafv3AddressBookDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	parts := strings.Split(d.Id(), ":")
	action := "DeleteDefenseRule"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["RuleIds"] = parts[1]
	request["InstanceId"] = parts[0]
	request["RegionId"] = client.RegionId

	request["DefenseType"] = "global"
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.RpcPost("waf-openapi", "2021-10-01", action, query, request, true)
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
