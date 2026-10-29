package ebs

import (
	"fmt"
	"log"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceDedicatedBlockStorageCluster() *schema.Resource {
	return &schema.Resource{
		Create: resourceDedicatedBlockStorageClusterCreate,
		Read:   resourceDedicatedBlockStorageClusterRead,
		Update: resourceDedicatedBlockStorageClusterUpdate,
		Delete: resourceDedicatedBlockStorageClusterDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"available_capacity": {
				Computed: true,
				Type:     schema.TypeString,
			},
			"category": {
				Computed: true,
				Type:     schema.TypeString,
			},
			"create_time": {
				Computed: true,
				Type:     schema.TypeString,
			},
			"dedicated_block_storage_cluster_id": {
				Computed: true,
				Type:     schema.TypeString,
			},
			"dedicated_block_storage_cluster_name": {
				Required: true,
				Type:     schema.TypeString,
			},
			"delivery_capacity": {
				Computed: true,
				Type:     schema.TypeString,
			},
			"description": {
				Optional: true,
				Computed: true,
				Type:     schema.TypeString,
			},
			"expired_time": {
				Computed: true,
				Type:     schema.TypeString,
			},
			"performance_level": {
				Computed: true,
				Type:     schema.TypeString,
			},
			"resource_group_id": {
				Computed: true,
				Type:     schema.TypeString,
			},
			"status": {
				Computed: true,
				Type:     schema.TypeString,
			},
			"supported_category": {
				Computed: true,
				Type:     schema.TypeString,
			},
			"total_capacity": {
				Required: true,
				ForceNew: true,
				Type:     schema.TypeString,
			},
			"type": {
				Required: true,
				ForceNew: true,
				Type:     schema.TypeString,
			},
			"used_capacity": {
				Computed: true,
				Type:     schema.TypeString,
			},
			"zone_id": {
				Required: true,
				ForceNew: true,
				Type:     schema.TypeString,
			},
		},
	}
}

func resourceDedicatedBlockStorageClusterCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ebsService := NewEbsService(client)
	request := map[string]interface{}{
		"RegionId": client.RegionId,
	}
	var err error

	if v, ok := d.GetOk("dedicated_block_storage_cluster_name"); ok {
		request["DbscName"] = v
	}
	if v, ok := d.GetOk("total_capacity"); ok {
		request["Capacity"] = v
	}
	if v, ok := d.GetOk("zone_id"); ok {
		request["Azone"] = v
	}
	request["Type"] = d.Get("type")

	var response map[string]interface{}
	action := "CreateDedicatedBlockStorageCluster"
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutCreate)), func() *retry.RetryError {
		resp, err := client.RpcPost("ebs", "2021-07-30", action, nil, request, false)
		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		response = resp
		logs.AddDebug(action, response, request)
		return nil
	})
	if err != nil {
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_ebs_dedicated_block_storage_cluster", action, errs.AlibabaCloudSdkGoERROR)
	}

	if v, err := jsonpath.Get("$.DbscId", response); err != nil || v == nil {
		return errs.WrapErrorf(err, errs.IdMsg, "alicloud_ebs_dedicated_block_storage_cluster")
	} else {
		d.SetId(fmt.Sprint(v))
	}

	stateConf := helper.BuildStateConf([]string{}, []string{"Running"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, ebsService.EbsDedicatedBlockStorageClusterStateRefreshFunc(d, []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return errs.WrapErrorf(err, errs.IdMsg, d.Id())
	}

	return resourceDedicatedBlockStorageClusterRead(d, meta)
}

func resourceDedicatedBlockStorageClusterRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ebsService := NewEbsService(client)

	object, err := ebsService.DescribeEbsDedicatedBlockStorageCluster(d.Id())
	if err != nil {
		if errs.NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_ebs_dedicated_block_storage_cluster ebsService.DescribeEbsDedicatedBlockStorageCluster Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return errs.WrapError(err)
	}

	d.Set("available_capacity", object["DedicatedBlockStorageClusterCapacity"].(map[string]interface{})["AvailableCapacity"])
	d.Set("category", object["Category"])
	d.Set("create_time", object["CreateTime"])
	d.Set("dedicated_block_storage_cluster_name", object["DedicatedBlockStorageClusterName"])
	d.Set("delivery_capacity", object["DedicatedBlockStorageClusterCapacity"].(map[string]interface{})["DeliveryCapacity"])
	d.Set("description", object["Description"])
	d.Set("expired_time", object["ExpiredTime"])
	d.Set("performance_level", object["PerformanceLevel"])
	d.Set("resource_group_id", object["ResourceGroupId"])
	d.Set("status", object["Status"])
	d.Set("supported_category", object["SupportedCategory"])
	d.Set("total_capacity", object["DedicatedBlockStorageClusterCapacity"].(map[string]interface{})["TotalCapacity"])
	d.Set("type", object["Type"])
	d.Set("used_capacity", object["DedicatedBlockStorageClusterCapacity"].(map[string]interface{})["UsedCapacity"])
	d.Set("zone_id", object["ZoneId"])

	return nil
}

func resourceDedicatedBlockStorageClusterUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	var err error
	d.Partial(true)

	update := false
	request := map[string]interface{}{
		"DbscId":   d.Id(),
		"RegionId": client.RegionId,
	}

	request["DbscName"] = d.Get("dedicated_block_storage_cluster_name")

	if d.HasChange("description") {
		update = true
		if v, ok := d.GetOk("description"); ok {
			request["Description"] = v
		}
	}

	if update {
		action := "ModifyDedicatedBlockStorageClusterAttribute"
		wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
		err = retry.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *retry.RetryError {
			resp, err := client.RpcPost("ebs", "2021-07-30", action, nil, request, true)
			if err != nil {
				if helper.NeedRetry(err) {
					wait()
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
			}
			logs.AddDebug(action, resp, request)
			return nil
		})
		if err != nil {
			return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
		}
	}

	d.Partial(false)
	return resourceDedicatedBlockStorageClusterRead(d, meta)
}

func resourceDedicatedBlockStorageClusterDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[WARN] Cannot destroy resourceAlicloudEbsDedicatedBlockStorageCluster. Terraform will remove this resource from the state file, however resources may remain.")
	return nil
}
