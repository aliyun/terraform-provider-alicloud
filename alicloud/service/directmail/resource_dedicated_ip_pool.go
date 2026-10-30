// Package directmail. This file is generated automatically. Please do not modify it manually, thank you!
package directmail

import (
	"fmt"
	"log"
	"sort"
	"strings"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceDedicatedIpPool() *schema.Resource {
	return &schema.Resource{
		Create: ResourceDedicatedIpPoolCreate,
		Read:   ResourceDedicatedIpPoolRead,
		Update: ResourceDedicatedIpPoolUpdate,
		Delete: ResourceDedicatedIpPoolDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"buy_resource_ids": {
				Type:     schema.TypeString,
				Optional: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					// Member order is insignificant to the API; compare the
					// normalized (trimmed, sorted) sets to avoid a permanent
					// update loop when the configuration lists IDs unordered.
					return fmt.Sprint(DirectMailBuyResourceIdSet(old)) == fmt.Sprint(DirectMailBuyResourceIdSet(new))
				},
			},
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ip_count": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"ips": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ip": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"zone_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func ResourceDedicatedIpPoolCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "DedicatedIpPoolCreate"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["RegionId"] = client.RegionId

	if v, ok := d.GetOk("buy_resource_ids"); ok {
		request["BuyResourceIds"] = strings.Join(DirectMailBuyResourceIdSet(v.(string)), ",")
	}
	request["Name"] = d.Get("name")
	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		response, err = client.RpcPost("Dm", "2015-11-23", action, query, request, true)
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
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_direct_mail_dedicated_ip_pool", action, errs.AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["Id"]))

	return ResourceDedicatedIpPoolRead(d, meta)
}

func ResourceDedicatedIpPoolRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	directMailServiceV2 := NewDirectMailServiceV2(client)

	objectRaw, err := directMailServiceV2.DescribeDirectMailDedicatedIpPool(d.Id())
	if err != nil {
		if !d.IsNewResource() && errs.NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_direct_mail_dedicated_ip_pool DescribeDirectMailDedicatedIpPool Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return errs.WrapError(err)
	}

	d.Set("create_time", objectRaw["CreateTime"])
	d.Set("ip_count", objectRaw["IpCount"])
	d.Set("name", objectRaw["Name"])

	ipsMaps := make([]map[string]interface{}, 0)
	buyResourceIds := make([]string, 0)
	if ipsRaw, ok := objectRaw["Ips"].([]interface{}); ok {
		for _, item := range ipsRaw {
			if ipItem, ok := item.(map[string]interface{}); ok {
				ipsMaps = append(ipsMaps, map[string]interface{}{
					"id":      ipItem["Id"],
					"ip":      ipItem["Ip"],
					"zone_id": ipItem["ZoneId"],
				})
				if id, ok := ipItem["Id"].(string); ok && id != "" {
					buyResourceIds = append(buyResourceIds, id)
				}
			}
		}
	}
	if err := d.Set("ips", ipsMaps); err != nil {
		return err
	}
	// BuyResourceIds is accepted as a comma-separated string while the Get API
	// returns Ips[*].Id; record the normalized actual membership, including an
	// empty one, so drift is reconciled by a regular update.
	sort.Strings(buyResourceIds)
	if err := d.Set("buy_resource_ids", strings.Join(buyResourceIds, ",")); err != nil {
		return errs.WrapError(err)
	}

	return nil
}

func ResourceDedicatedIpPoolUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "DedicatedIpPoolUpdate"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["Id"] = d.Id()
	request["RegionId"] = client.RegionId

	if d.HasChange("buy_resource_ids") {
		request["BuyResourceIds"] = strings.Join(DirectMailBuyResourceIdSet(d.Get("buy_resource_ids").(string)), ",")
		request["UpdateResource"] = true
	} else {
		return ResourceDedicatedIpPoolRead(d, meta)
	}

	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
		response, err = client.RpcPost("Dm", "2015-11-23", action, query, request, true)
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

	directMailServiceV2 := NewDirectMailServiceV2(client)
	err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
		objectRaw, err := directMailServiceV2.DescribeDirectMailDedicatedIpPool(d.Id())
		if err != nil {
			return retry.NonRetryableError(err)
		}
		ipsRaw, _ := objectRaw["Ips"].([]interface{})
		want := DirectMailBuyResourceIdSet(d.Get("buy_resource_ids").(string))
		got := make([]string, 0)
		for _, item := range ipsRaw {
			if ipItem, ok := item.(map[string]interface{}); ok {
				if id, ok := ipItem["Id"].(string); ok && id != "" {
					got = append(got, id)
				}
			}
		}
		sort.Strings(got)
		if fmt.Sprint(want) != fmt.Sprint(got) {
			return retry.RetryableError(fmt.Errorf("waiting for IP pool %s members to converge", d.Id()))
		}
		return nil
	})
	if err != nil {
		return errs.WrapError(err)
	}

	// Refresh Computed fields (create_time, ip_count, ips) from the converged
	// pool so state reflect the update result.
	return ResourceDedicatedIpPoolRead(d, meta)
}

func ResourceDedicatedIpPoolDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	action := "DedicatedIpPoolDelete"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["Id"] = d.Id()
	request["RegionId"] = client.RegionId

	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutDelete), func() *retry.RetryError {
		response, err = client.RpcPost("Dm", "2015-11-23", action, query, request, true)
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

// DirectMailBuyResourceIdSet splits a comma-separated BuyResourceIds value,
// trims surrounding spaces from each member, drops empty members and sorts
// the result so configuration and state comparisons are order-insensitive.
func DirectMailBuyResourceIdSet(value string) []string {
	items := make([]string, 0)
	for _, item := range strings.Split(value, ",") {
		if trimmed := strings.TrimSpace(item); trimmed != "" {
			items = append(items, trimmed)
		}
	}
	sort.Strings(items)
	return items
}
