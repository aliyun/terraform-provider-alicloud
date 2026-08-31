// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"log"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudSimpleApplicationServerDisk() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudSimpleApplicationServerDiskCreate,
		Read:   resourceAliCloudSimpleApplicationServerDiskRead,
		Update: resourceAliCloudSimpleApplicationServerDiskUpdate,
		Delete: resourceAliCloudSimpleApplicationServerDiskDelete,
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
			"disk_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"disk_size": {
				Type:         schema.TypeInt,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: IntBetween(0, 16380),
			},
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"region_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"remark": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceAliCloudSimpleApplicationServerDiskCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "CreateDisk"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["RegionId"] = client.RegionId
	request["ClientToken"] = buildClientToken(action)

	request["DiskSize"] = d.Get("disk_size")
	request["InstanceId"] = d.Get("instance_id")
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RpcPost("SWAS-OPEN", "2020-06-01", action, query, request, true)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_simple_application_server_disk", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["DiskId"]))

	return resourceAliCloudSimpleApplicationServerDiskUpdate(d, meta)
}

func resourceAliCloudSimpleApplicationServerDiskRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	simpleApplicationServerServiceV2 := SimpleApplicationServerServiceV2{client}

	objectRaw, err := simpleApplicationServerServiceV2.DescribeSimpleApplicationServerDisk(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_simple_application_server_disk DescribeSimpleApplicationServerDisk Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("create_time", objectRaw["CreationTime"])
	d.Set("disk_name", objectRaw["DiskName"])
	d.Set("disk_size", objectRaw["Size"])
	d.Set("instance_id", objectRaw["InstanceId"])
	d.Set("region_id", objectRaw["RegionId"])
	d.Set("remark", objectRaw["Remark"])

	return nil
}

func resourceAliCloudSimpleApplicationServerDiskUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	update := false

	var err error
	action := "UpdateDiskAttribute"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["DiskId"] = d.Id()
	request["RegionId"] = client.RegionId
	request["ClientToken"] = buildClientToken(action)
	if d.HasChange("remark") {
		update = true
	}
	request["Remark"] = d.Get("remark")
	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("SWAS-OPEN", "2020-06-01", action, query, request, true)
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

	return resourceAliCloudSimpleApplicationServerDiskRead(d, meta)
}

func resourceAliCloudSimpleApplicationServerDiskDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[WARN] Cannot destroy resource AliCloud Resource Disk. Terraform will remove this resource from the state file, however resources may remain.")
	return nil
}
