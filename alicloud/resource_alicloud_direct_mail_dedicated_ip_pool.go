// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"log"
	"sort"
	"strings"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudDirectMailDedicatedIpPool() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudDirectMailDedicatedIpPoolCreate,
		Read:   resourceAliCloudDirectMailDedicatedIpPoolRead,
		Update: resourceAliCloudDirectMailDedicatedIpPoolUpdate,
		Delete: resourceAliCloudDirectMailDedicatedIpPoolDelete,
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

func resourceAliCloudDirectMailDedicatedIpPoolCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "DedicatedIpPoolCreate"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["RegionId"] = client.RegionId

	if v, ok := d.GetOk("buy_resource_ids"); ok {
		request["BuyResourceIds"] = v
	}
	request["Name"] = d.Get("name")
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RpcPost("Dm", "2015-11-23", action, query, request, true)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_direct_mail_dedicated_ip_pool", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["Id"]))

	return resourceAliCloudDirectMailDedicatedIpPoolRead(d, meta)
}

func resourceAliCloudDirectMailDedicatedIpPoolRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	directMailServiceV2 := DirectMailServiceV2{client}

	objectRaw, err := directMailServiceV2.DescribeDirectMailDedicatedIpPool(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_direct_mail_dedicated_ip_pool DescribeDirectMailDedicatedIpPool Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
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
	// returns Ips[*].Id; normalize the read-back form so state matches config.
	if len(buyResourceIds) > 0 {
		sort.Strings(buyResourceIds)
		d.Set("buy_resource_ids", strings.Join(buyResourceIds, ","))
	} else if _, ok := d.GetOk("buy_resource_ids"); ok {
		return WrapError(fmt.Errorf("IP pool %s no longer contains any purchased IP instances; remove buy_resource_ids from the configuration or re-create the pool", d.Id()))
	}

	return nil
}

func resourceAliCloudDirectMailDedicatedIpPoolUpdate(d *schema.ResourceData, meta interface{}) error {
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
		request["BuyResourceIds"] = d.Get("buy_resource_ids")
		request["UpdateResource"] = true
	} else {
		return resourceAliCloudDirectMailDedicatedIpPoolRead(d, meta)
	}

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
		response, err = client.RpcPost("Dm", "2015-11-23", action, query, request, true)
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

	directMailServiceV2 := DirectMailServiceV2{client}
	return resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
		objectRaw, err := directMailServiceV2.DescribeDirectMailDedicatedIpPool(d.Id())
		if err != nil {
			return resource.NonRetryableError(err)
		}
		ipsRaw, _ := objectRaw["Ips"].([]interface{})
		want := strings.Split(d.Get("buy_resource_ids").(string), ",")
		got := make([]string, 0)
		for _, item := range ipsRaw {
			if ipItem, ok := item.(map[string]interface{}); ok {
				if id, ok := ipItem["Id"].(string); ok && id != "" {
					got = append(got, id)
				}
			}
		}
		sort.Strings(want)
		sort.Strings(got)
		if fmt.Sprint(want) != fmt.Sprint(got) {
			return resource.RetryableError(fmt.Errorf("waiting for IP pool %s members to converge", d.Id()))
		}
		return nil
	})
}

func resourceAliCloudDirectMailDedicatedIpPoolDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	action := "DedicatedIpPoolDelete"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["Id"] = d.Id()
	request["RegionId"] = client.RegionId

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.RpcPost("Dm", "2015-11-23", action, query, request, true)
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
