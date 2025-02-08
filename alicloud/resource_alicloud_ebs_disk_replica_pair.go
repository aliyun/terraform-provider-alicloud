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

func resourceAlicloudEbsDiskReplicaPair() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudEbsDiskReplicaPairCreate,
		Read:   resourceAlicloudEbsDiskReplicaPairRead,
		Update: resourceAlicloudEbsDiskReplicaPairUpdate,
		Delete: resourceAlicloudEbsDiskReplicaPairDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"bandwidth": {
				Optional: true,
				ForceNew: true,
				Computed: true,
				Type:     schema.TypeString,
			},
			"create_time": {
				Computed: true,
				Type:     schema.TypeString,
			},
			"description": {
				Optional: true,
				Type:     schema.TypeString,
			},
			"destination_disk_id": {
				Required: true,
				ForceNew: true,
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
			"disk_id": {
				Required: true,
				ForceNew: true,
				Type:     schema.TypeString,
			},
			"pair_name": {
				Optional: true,
				Type:     schema.TypeString,
			},
			"payment_type": {
				Optional: true,
				ForceNew: true,
				Type:     schema.TypeString,
			},
			"period": {
				Optional: true,
				Type:     schema.TypeString,
			},
			"period_unit": {
				Optional: true,
				Type:     schema.TypeString,
			},
			"rpo": {
				Optional: true,
				ForceNew: true,
				Computed: true,
				Type:     schema.TypeString,
			},
			"replica_pair_id": {
				Optional: true,
				Computed: true,
				Type:     schema.TypeString,
			},
			"resource_group_id": {
				Computed: true,
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

func resourceAlicloudEbsDiskReplicaPairCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ebsService := EbsService{client}
	request := map[string]interface{}{
		"RegionId": client.RegionId,
	}
	var err error

	if v, ok := d.GetOk("bandwidth"); ok {
		request["Bandwidth"] = v
	}
	if v, ok := d.GetOk("description"); ok {
		request["Description"] = v
	}
	if v, ok := d.GetOk("destination_disk_id"); ok {
		request["DestinationDiskId"] = v
	}
	if v, ok := d.GetOk("destination_region_id"); ok {
		request["DestinationRegionId"] = v
	}
	if v, ok := d.GetOk("destination_zone_id"); ok {
		request["DestinationZoneId"] = v
	}
	if v, ok := d.GetOk("disk_id"); ok {
		request["DiskId"] = v
	}
	if v, ok := d.GetOk("pair_name"); ok {
		request["PairName"] = v
	}
	if v, ok := d.GetOk("payment_type"); ok {
		request["ChargeType"] = v
	}
	if v, ok := d.GetOk("period"); ok {
		request["Period"] = v
	}
	if v, ok := d.GetOk("period_unit"); ok {
		request["PeriodUnit"] = v
	}
	if v, ok := d.GetOk("rpo"); ok {
		request["RPO"] = v
	}
	if v, ok := d.GetOk("source_zone_id"); ok {
		request["SourceZoneId"] = v
	}

	request["ClientToken"] = buildClientToken("CreateDiskReplicaPair")
	var response map[string]interface{}
	action := "CreateDiskReplicaPair"
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_ebs_disk_replica_pair", action, AlibabaCloudSdkGoERROR)
	}

	if v, err := jsonpath.Get("$.ReplicaPairId", response); err != nil || v == nil {
		return WrapErrorf(err, IdMsg, "alicloud_ebs_disk_replica_pair")
	} else {
		d.SetId(fmt.Sprint(v))
	}
	stateConf := BuildStateConf([]string{}, []string{"created", "initial_syncing", "manual_syncing", "syncing", "normal"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, ebsService.EbsDiskReplicaPairStateRefreshFunc(d, []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}
	return resourceAlicloudEbsDiskReplicaPairRead(d, meta)
}

func resourceAlicloudEbsDiskReplicaPairRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ebsService := EbsService{client}

	object, err := ebsService.DescribeEbsDiskReplicaPair(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_ebs_disk_replica_pair ebsService.DescribeEbsDiskReplicaPair Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	d.Set("bandwidth", object["Bandwidth"])
	d.Set("create_time", object["CreateTime"])
	d.Set("description", object["Description"])
	d.Set("destination_disk_id", object["DestinationDiskId"])
	d.Set("destination_region_id", object["DestinationRegion"])
	d.Set("destination_zone_id", object["DestinationZoneId"])
	d.Set("disk_id", object["SourceDiskId"])
	d.Set("pair_name", object["PairName"])
	d.Set("payment_type", object["ChargeType"])
	d.Set("rpo", object["RPO"])
	d.Set("resource_group_id", object["ResourceGroupId"])
	d.Set("source_zone_id", object["SourceZoneId"])
	d.Set("status", object["Status"])

	return nil
}

func resourceAlicloudEbsDiskReplicaPairUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	var err error
	d.Partial(true)
	update := false
	request := map[string]interface{}{
		"ReplicaPairId": d.Id(),
		"RegionId":      client.RegionId,
	}

	if d.HasChange("description") {
		update = true
		if v, ok := d.GetOk("description"); ok {
			request["Description"] = v
		}
	}
	if d.HasChange("pair_name") {
		update = true
		if v, ok := d.GetOk("pair_name"); ok {
			request["PairName"] = v
		}
	}

	if update {
		action := "ModifyDiskReplicaPair"
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
		d.SetPartial("description")
		d.SetPartial("pair_name")
	}

	d.Partial(false)
	return resourceAlicloudEbsDiskReplicaPairRead(d, meta)
}

func resourceAlicloudEbsDiskReplicaPairDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var err error

	request := map[string]interface{}{

		"ReplicaPairId": d.Id(),
		"RegionId":      client.RegionId,
	}

	request["ClientToken"] = buildClientToken("DeleteDiskReplicaPair")
	action := "DeleteDiskReplicaPair"
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
