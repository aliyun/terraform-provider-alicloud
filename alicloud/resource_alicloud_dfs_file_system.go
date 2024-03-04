// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"log"
	"time"

	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudDfsFileSystem() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudDfsFileSystemCreate,
		Read:   resourceAliCloudDfsFileSystemRead,
		Update: resourceAliCloudDfsFileSystemUpdate,
		Delete: resourceAliCloudDfsFileSystemDelete,
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
				ValidateFunc: StringInSlice([]string{"LRS", "ZRS"}, true),
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
				ValidateFunc: StringInSlice([]string{"HDFS"}, true),
			},
			"provisioned_throughput_in_mi_bps": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: IntBetween(0, 1024),
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					if v, ok := d.GetOk("throughput_mode"); ok && v.(string) == "Provisioned" {
						return false
					}
					return true
				},
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
				ValidateFunc: StringInSlice([]string{"STANDARD", "PERFORMANCE"}, true),
			},
			"throughput_mode": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: StringInSlice([]string{"Standard", "Provisioned"}, true),
			},
			"zone_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceAliCloudDfsFileSystemCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "CreateFileSystem"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	conn, err := client.NewDfsClient()
	if err != nil {
		return WrapError(err)
	}
	request = make(map[string]interface{})

	request["ZoneId"] = d.Get("zone_id")
	request["InputRegionId"] = client.RegionId
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
	if v, ok := d.GetOk("provisioned_throughput_in_mi_bps"); ok {
		request["ProvisionedThroughputInMiBps"] = v
	}
	if v, ok := d.GetOk("storage_set_name"); ok {
		request["StorageSetName"] = v
	}
	if v, ok := d.GetOk("partition_number"); ok {
		request["PartitionNumber"] = v
	}
	if v, ok := d.GetOk("data_redundancy_type"); ok {
		request["DataRedundancyType"] = v
	}
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2018-06-20"), StringPointer("AK"), query, request, &runtime)

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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_dfs_file_system", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["FileSystemId"]))

	return resourceAliCloudDfsFileSystemRead(d, meta)
}

func resourceAliCloudDfsFileSystemRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	dfsServiceV2 := DfsServiceV2{client}

	objectRaw, err := dfsServiceV2.DescribeDfsFileSystem(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_dfs_file_system DescribeDfsFileSystem Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("create_time", objectRaw["CreateTime"])
	d.Set("description", objectRaw["Description"])
	d.Set("file_system_name", objectRaw["FileSystemName"])
	d.Set("protocol_type", objectRaw["ProtocolType"])
	d.Set("provisioned_throughput_in_mi_bps", formatInt(objectRaw["ProvisionedThroughputInMiBps"]))
	d.Set("space_capacity", objectRaw["SpaceCapacity"])
	d.Set("storage_type", objectRaw["StorageType"])
	d.Set("throughput_mode", objectRaw["ThroughputMode"])
	d.Set("zone_id", objectRaw["ZoneId"])

	return nil
}

func resourceAliCloudDfsFileSystemUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	update := false
	action := "ModifyFileSystem"
	conn, err := client.NewDfsClient()
	if err != nil {
		return WrapError(err)
	}
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	query["FileSystemId"] = d.Id()
	query["InputRegionId"] = client.RegionId
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
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2018-06-20"), StringPointer("AK"), query, request, &runtime)

			if err != nil {
				if NeedRetry(err) || IsExpectedErrors(err, []string{"FileSystem.ModifyThroughputModeTooFrequent"}) {
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

	return resourceAliCloudDfsFileSystemRead(d, meta)
}

func resourceAliCloudDfsFileSystemDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	action := "DeleteFileSystem"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	conn, err := client.NewDfsClient()
	if err != nil {
		return WrapError(err)
	}
	request = make(map[string]interface{})
	query["FileSystemId"] = d.Id()
	query["InputRegionId"] = client.RegionId

	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2018-06-20"), StringPointer("AK"), query, request, &runtime)

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
		if IsExpectedErrors(err, []string{"InvalidParameter.FileSystemNotFound"}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	return nil
}
