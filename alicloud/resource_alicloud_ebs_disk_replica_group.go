package alicloud

import (
	"fmt"
	"log"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAlicloudEbsDiskReplicaGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudEbsDiskReplicaGroupCreate,
		Read:   resourceAlicloudEbsDiskReplicaGroupRead,
		Update: resourceAlicloudEbsDiskReplicaGroupUpdate,
		Delete: resourceAlicloudEbsDiskReplicaGroupDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(1 * time.Minute),
			Update: schema.DefaultTimeout(1 * time.Minute),
			Delete: schema.DefaultTimeout(1 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"description": {
				Optional: true,
				Type:     schema.TypeString,
			},
			"destination_region_id": {
				Required: true,
				ForceNew: true,
				Type:     schema.TypeString,
			},
			"destination_zone_id": {
				Required: true,
				ForceNew: true,
				Type:     schema.TypeString,
			},
			"group_name": {
				Optional: true,
				Type:     schema.TypeString,
			},
			"rpo": {
				Optional: true,
				ForceNew: true,
				Computed: true,
				Type:     schema.TypeInt,
			},
			"source_region_id": {
				Required: true,
				ForceNew: true,
				Type:     schema.TypeString,
			},
			"source_zone_id": {
				Required: true,
				ForceNew: true,
				Type:     schema.TypeString,
			},
			"status": {
				Computed: true,
				Type:     schema.TypeString,
			},
		},
	}
}

func resourceAlicloudEbsDiskReplicaGroupCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ebsService := EbsService{client}
	request := make(map[string]interface{})
	var err error

	if v, ok := d.GetOk("description"); ok {
		request["Description"] = v
	}
	if v, ok := d.GetOk("destination_region_id"); ok {
		request["DestinationRegionId"] = v
	}
	if v, ok := d.GetOk("destination_zone_id"); ok {
		request["DestinationZoneId"] = v
	}
	if v, ok := d.GetOk("group_name"); ok {
		request["GroupName"] = v
	}
	if v, ok := d.GetOk("rpo"); ok {
		request["RPO"] = v
	}
	if v, ok := d.GetOk("source_region_id"); ok {
		request["RegionId"] = v
	}

	if v, ok := d.GetOk("source_zone_id"); ok {
		request["SourceZoneId"] = v
	}

	request["ClientToken"] = buildClientToken("CreateDiskReplicaGroup")
	var response map[string]interface{}
	action := "CreateDiskReplicaGroup"
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutCreate)), func() *resource.RetryError {
		resp, err := client.RpcPost("ebs", "2021-07-30", action, nil, request, true)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		response = resp
		addDebug(action, response, request)
		return nil
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_ebs_disk_replica_group", action, AlibabaCloudSdkGoERROR)
	}

	if v, err := jsonpath.Get("$.ReplicaGroupId", response); err != nil || v == nil {
		return WrapErrorf(err, IdMsg, "alicloud_ebs_disk_replica_group")
	} else {
		d.SetId(fmt.Sprint(v))
	}
	stateConf := BuildStateConf([]string{}, []string{"created"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, ebsService.EbsDiskReplicaGroupStateRefreshFunc(d, []string{"create_failed", "invalid"}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}
	return resourceAlicloudEbsDiskReplicaGroupRead(d, meta)
}

func resourceAlicloudEbsDiskReplicaGroupRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ebsService := EbsService{client}

	object, err := ebsService.DescribeEbsDiskReplicaGroup(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_ebs_disk_replica_group ebsService.DescribeEbsDiskReplicaGroup Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("description", object["Description"])
	d.Set("destination_region_id", object["DestinationRegionId"])
	d.Set("destination_zone_id", object["DestinationZoneId"])
	d.Set("group_name", object["GroupName"])
	d.Set("rpo", object["RPO"])
	d.Set("source_region_id", object["SourceRegionId"])
	d.Set("source_zone_id", object["SourceZoneId"])
	d.Set("status", object["Status"])
	return nil
}

func resourceAlicloudEbsDiskReplicaGroupUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var err error
	update := false
	request := map[string]interface{}{
		"ReplicaGroupId": d.Id(),
		"RegionId":       client.RegionId,
	}

	if d.HasChange("description") {
		update = true
		if v, ok := d.GetOk("description"); ok {
			request["Description"] = v
		}
	}
	if d.HasChange("group_name") {
		update = true
		if v, ok := d.GetOk("group_name"); ok {
			request["GroupName"] = v
		}
	}

	if update {
		action := "ModifyDiskReplicaGroup"
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *resource.RetryError {
			resp, err := client.RpcPost("ebs", "2021-07-30", action, nil, request, true)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			addDebug(action, resp, request)
			return nil
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
	}

	return resourceAlicloudEbsDiskReplicaGroupRead(d, meta)
}

func resourceAlicloudEbsDiskReplicaGroupDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var err error

	request := map[string]interface{}{
		"ReplicaGroupId": d.Id(),
		"RegionId":       client.RegionId,
	}

	request["ClientToken"] = buildClientToken("DeleteDiskReplicaGroup")
	action := "DeleteDiskReplicaGroup"
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutDelete)), func() *resource.RetryError {
		resp, err := client.RpcPost("ebs", "2021-07-30", action, nil, request, true)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, resp, request)
		return nil
	})
	if err != nil {
		if NotFoundError(err) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}
	return nil
}
