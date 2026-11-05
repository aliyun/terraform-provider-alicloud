// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package wafopenapi

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/convert"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/validators"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceWafv3AddressBook() *schema.Resource {
	return &schema.Resource{
		Create: resourceWafv3AddressBookCreate,
		Read:   resourceWafv3AddressBookRead,
		Update: resourceWafv3AddressBookUpdate,
		Delete: resourceWafv3AddressBookDelete,
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
				ValidateFunc: validators.StringInSlice([]string{"ip"}, false),
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

func resourceWafv3AddressBookCreate(d *schema.ResourceData, meta interface{}) error {

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
		return errs.WrapError(err)
	}
	request["Rules"] = string(rulesJson)

	request["DefenseScene"] = "address_book"
	request["DefenseType"] = "global"
	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		response, err = client.RpcPost("waf-openapi", "2021-10-01", action, query, request, true)
		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(action, response, request)

	if err != nil {
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_wafv3_address_book", action, errs.AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprintf("%v:%v", request["InstanceId"], response["RuleIds"]))

	return resourceWafv3AddressBookUpdate(d, meta)
}

func resourceWafv3AddressBookRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	wafv3ServiceV2 := Wafv3ServiceV2{client}

	objectRaw, err := wafv3ServiceV2.DescribeWafv3AddressBook(d.Id())
	if err != nil {
		if !d.IsNewResource() && errs.NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_wafv3_address_book DescribeWafv3AddressBook Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return errs.WrapError(err)
	}

	d.Set("address_book_name", objectRaw["name"])
	d.Set("address_book_type", objectRaw["valueType"])
	d.Set("description", objectRaw["description"])

	addressesRaw, err := wafv3ServiceV2.DescribeWafv3AddressBookAddresses(d.Id())
	if err != nil && !errs.NotFoundError(err) {
		return errs.WrapError(err)
	}

	addressListRawObj, _ := jsonpath.Get("$.AddressList", addressesRaw)
	addressListRaw := make([]interface{}, 0)
	if addressListRawObj != nil {
		addressListRaw = convert.ConvertToInterfaceArray(addressListRawObj)
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

func resourceWafv3AddressBookUpdate(d *schema.ResourceData, meta interface{}) error {
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
			return errs.WrapError(err)
		}
		request["Rules"] = string(rulesJson)

		wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
		err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
			response, err = client.RpcPost("waf-openapi", "2021-10-01", action, query, request, true)
			if err != nil {
				if helper.NeedRetry(err) {
					wait()
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
			}
			return nil
		})
		logs.AddDebug(action, response, request)
		if err != nil {
			return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
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

			wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
			err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
				response, err = client.RpcPost("waf-openapi", "2021-10-01", action, query, request, true)
				if err != nil {
					if helper.NeedRetry(err) {
						wait()
						return retry.RetryableError(err)
					}
					return retry.NonRetryableError(err)
				}
				return nil
			})
			logs.AddDebug(action, response, request)
			if err != nil {
				return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
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

			wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
			err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
				response, err = client.RpcPost("waf-openapi", "2021-10-01", action, query, request, true)
				if err != nil {
					if helper.NeedRetry(err) {
						wait()
						return retry.RetryableError(err)
					}
					return retry.NonRetryableError(err)
				}
				return nil
			})
			logs.AddDebug(action, response, request)
			if err != nil {
				return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
			}

		}

	}
	return resourceWafv3AddressBookRead(d, meta)
}

func resourceWafv3AddressBookDelete(d *schema.ResourceData, meta interface{}) error {

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
	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutDelete), func() *retry.RetryError {
		response, err = client.RpcPost("waf-openapi", "2021-10-01", action, query, request, true)
		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(action, response, request)

	if err != nil {
		if errs.NotFoundError(err) {
			return nil
		}
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
	}

	return nil
}
