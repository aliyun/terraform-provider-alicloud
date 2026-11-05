// Package dfs. This file is generated automatically. Please do not modify it manually, thank you!
package dfs

import (
	"fmt"
	"log"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/convert"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/validators"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceFileSystem() *schema.Resource {
	return &schema.Resource{
		Create: resourceFileSystemCreate,
		Read:   resourceFileSystemRead,
		Update: resourceFileSystemUpdate,
		Delete: resourceFileSystemDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"data_redundancy_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validators.StringInSlice([]string{"LRS", "ZRS"}, false),
			},
			"dedicated_cluster_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"file_system_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"partition_number": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"protocol_type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validators.StringInSlice([]string{"HDFS", "PANGU"}, false),
			},
			"provisioned_throughput_in_mi_bps": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validators.IntBetween(0, 1024),
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					if v, ok := d.GetOk("throughput_mode"); ok && v.(string) == "Provisioned" {
						return false
					}
					return true
				},
			},
			"region_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"space_capacity": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"storage_set_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"storage_type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validators.StringInSlice([]string{"STANDARD", "PERFORMANCE"}, false),
			},
			"throughput_mode": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validators.StringInSlice([]string{"Standard", "Provisioned"}, false),
			},
			"zone_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
		},
	}
}

func resourceFileSystemCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "CreateFileSystem"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["InputRegionId"] = client.RegionId

	if v, ok := d.GetOk("zone_id"); ok {
		request["ZoneId"] = v
	}
	request["ProtocolType"] = d.Get("protocol_type")
	request["StorageType"] = d.Get("storage_type")
	if v, ok := d.GetOk("description"); ok {
		request["Description"] = v
	}
	request["FileSystemName"] = d.Get("file_system_name")
	request["SpaceCapacity"] = d.Get("space_capacity")
	if v, ok := d.GetOk("throughput_mode"); ok {
		request["ThroughputMode"] = v
	}
	if v, ok := d.GetOkExists("provisioned_throughput_in_mi_bps"); ok && v.(int) > 0 {
		request["ProvisionedThroughputInMiBps"] = v
	}
	if v, ok := d.GetOk("storage_set_name"); ok {
		request["StorageSetName"] = v
	}
	if v, ok := d.GetOkExists("partition_number"); ok {
		request["PartitionNumber"] = v
	}
	if v, ok := d.GetOk("data_redundancy_type"); ok {
		request["DataRedundancyType"] = v
	}
	if v, ok := d.GetOk("dedicated_cluster_id"); ok {
		request["DedicatedClusterId"] = v
	}
	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		response, err = client.RpcPost("DFS", "2018-06-20", action, query, request, true)
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
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_dfs_file_system", action, errs.AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["FileSystemId"]))

	return resourceFileSystemRead(d, meta)
}

func resourceFileSystemRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	dfsServiceV2 := ServiceV2{client}

	objectRaw, err := dfsServiceV2.DescribeDfsFileSystem(d.Id())
	if err != nil {
		if !d.IsNewResource() && errs.NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_dfs_file_system DescribeDfsFileSystem Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return errs.WrapError(err)
	}

	d.Set("create_time", objectRaw["CreateTime"])
	d.Set("description", objectRaw["Description"])
	d.Set("file_system_name", objectRaw["FileSystemName"])
	d.Set("protocol_type", objectRaw["ProtocolType"])
	provisioned_throughput_in_mi_bpsVal, err := convert.FormatInt(objectRaw["ProvisionedThroughputInMiBps"])
	if err != nil {
		return err
	}
	d.Set("provisioned_throughput_in_mi_bps", provisioned_throughput_in_mi_bpsVal)
	d.Set("region_id", objectRaw["RegionId"])
	d.Set("space_capacity", objectRaw["SpaceCapacity"])
	d.Set("storage_type", objectRaw["StorageType"])
	d.Set("throughput_mode", objectRaw["ThroughputMode"])
	d.Set("zone_id", objectRaw["ZoneId"])

	return nil
}

func resourceFileSystemUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	update := false

	var err error
	action := "ModifyFileSystem"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["FileSystemId"] = d.Id()
	request["InputRegionId"] = client.RegionId
	if d.HasChange("file_system_name") {
		update = true
	}
	request["FileSystemName"] = d.Get("file_system_name")
	if d.HasChange("space_capacity") {
		update = true
	}
	request["SpaceCapacity"] = d.Get("space_capacity")
	if d.HasChange("description") {
		update = true
		request["Description"] = d.Get("description")
	}

	if d.HasChange("throughput_mode") {
		update = true
	}
	if v, ok := d.GetOk("throughput_mode"); ok {
		request["ThroughputMode"] = v
	}

	if d.HasChange("provisioned_throughput_in_mi_bps") {
		update = true
	}
	if v, ok := d.GetOk("provisioned_throughput_in_mi_bps"); ok {
		request["ProvisionedThroughputInMiBps"] = v
	}

	if update {
		wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
		err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
			response, err = client.RpcPost("DFS", "2018-06-20", action, query, request, true)
			if err != nil {
				if helper.NeedRetry(err) || errs.IsExpectedErrors(err, []string{"FileSystem.ModifyThroughputModeTooFrequent"}) {
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

	return resourceFileSystemRead(d, meta)
}

func resourceFileSystemDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	action := "DeleteFileSystem"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["FileSystemId"] = d.Id()
	request["InputRegionId"] = client.RegionId

	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutDelete), func() *retry.RetryError {
		response, err = client.RpcPost("DFS", "2018-06-20", action, query, request, true)

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
		if errs.IsExpectedErrors(err, []string{"InvalidParameter.FileSystemNotFound"}) || errs.NotFoundError(err) {
			return nil
		}
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
	}

	return nil
}
