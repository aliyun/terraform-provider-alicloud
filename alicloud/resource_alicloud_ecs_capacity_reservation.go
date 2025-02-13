package alicloud

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceAlicloudEcsCapacityReservation() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudEcsCapacityReservationCreate,
		Read:   resourceAlicloudEcsCapacityReservationRead,
		Update: resourceAlicloudEcsCapacityReservationUpdate,
		Delete: resourceAlicloudEcsCapacityReservationDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"capacity_reservation_name": {
				Optional: true,
				Computed: true,
				Type:     schema.TypeString,
			},
			"description": {
				Optional: true,
				Type:     schema.TypeString,
			},
			"dry_run": {
				Optional: true,
				Type:     schema.TypeBool,
			},
			"end_time": {
				Optional: true,
				Computed: true,
				Type:     schema.TypeString,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					diff := d.Get("end_time_type").(string) == "Unlimited"
					if diff {
						return diff
					}
					if old != "" && new != "" && strings.HasPrefix(new, strings.Trim(old, "Z")) {
						diff = true
					}
					return diff
				},
			},
			"end_time_type": {
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice([]string{"Limited", "Unlimited"}, false),
				Type:         schema.TypeString,
			},
			"match_criteria": {
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"Open", "Target"}, false),
				Type:         schema.TypeString,
			},
			"payment_type": {
				Computed: true,
				Type:     schema.TypeString,
			},
			"platform": {
				Optional:     true,
				ForceNew:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice([]string{"linux", "windows"}, false),
				Type:         schema.TypeString,
			},
			"resource_group_id": {
				Optional: true,
				ForceNew: true,
				Type:     schema.TypeString,
			},
			"start_time": {
				Computed: true,
				Type:     schema.TypeString,
			},
			"start_time_type": {
				Computed: true,
				Type:     schema.TypeString,
			},
			"status": {
				Computed: true,
				Type:     schema.TypeString,
			},
			"tags": tagsSchema(),
			"time_slot": {
				Computed: true,
				Type:     schema.TypeString,
			},
			"instance_amount": {
				Required: true,
				Type:     schema.TypeInt,
			},
			"instance_type": {
				Required: true,
				ForceNew: true,
				Type:     schema.TypeString,
			},
			"zone_ids": {
				Required: true,
				ForceNew: true,
				Type:     schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func resourceAlicloudEcsCapacityReservationCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ecsService := EcsService{client}
	request := map[string]interface{}{
		"RegionId": client.RegionId,
	}
	var err error

	if v, ok := d.GetOk("capacity_reservation_name"); ok {
		request["PrivatePoolOptions.Name"] = v
	}
	if v, ok := d.GetOk("description"); ok {
		request["Description"] = v
	}
	if v, ok := d.GetOk("end_time"); ok {
		request["EndTime"] = v
	}
	if v, ok := d.GetOk("end_time_type"); ok {
		request["EndTimeType"] = v
	}
	if v, ok := d.GetOk("instance_amount"); ok {
		request["InstanceAmount"] = v
	}
	if v, ok := d.GetOk("instance_type"); ok {
		request["InstanceType"] = v
	}
	if v, ok := d.GetOk("match_criteria"); ok {
		request["PrivatePoolOptions.MatchCriteria"] = v
	}
	if v, ok := d.GetOk("platform"); ok {
		request["Platform"] = v
	}
	if v, ok := d.GetOk("resource_group_id"); ok {
		request["ResourceGroupId"] = v
	}
	if v, ok := d.GetOk("tags"); ok {
		count := 1
		for key, value := range v.(map[string]interface{}) {
			request[fmt.Sprintf("Tag.%d.Key", count)] = key
			request[fmt.Sprintf("Tag.%d.Value", count)] = value
			count++
		}
	}
	if v, ok := d.GetOk("zone_ids"); ok {
		request["ZoneId"] = v.([]interface{})
	}

	request["ClientToken"] = buildClientToken("CreateCapacityReservation")
	var response map[string]interface{}
	action := "CreateCapacityReservation"
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutCreate)), func() *resource.RetryError {
		resp, err := client.RpcPost("Ecs", "2014-05-26", action, nil, request, true)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_ecs_capacity_reservation", action, AlibabaCloudSdkGoERROR)
	}

	if v, err := jsonpath.Get("$.PrivatePoolOptionsId", response); err != nil || v == nil {
		return WrapErrorf(err, IdMsg, "alicloud_ecs_capacity_reservation")
	} else {
		d.SetId(fmt.Sprint(v))
	}
	stateConf := BuildStateConf([]string{}, []string{"Active"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, ecsService.EcsCapacityReservationStateRefreshFunc(d, []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}
	return resourceAlicloudEcsCapacityReservationRead(d, meta)
}

func resourceAlicloudEcsCapacityReservationRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ecsService := EcsService{client}

	object, err := ecsService.DescribeEcsCapacityReservation(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_ecs_capacity_reservation ecsService.DescribeEcsCapacityReservation Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	d.Set("capacity_reservation_name", object["PrivatePoolOptionsName"])
	d.Set("description", object["Description"])
	d.Set("end_time", object["EndTime"])
	d.Set("end_time_type", object["EndTimeType"])
	d.Set("match_criteria", object["PrivatePoolOptionsMatchCriteria"])
	d.Set("payment_type", object["InstanceChargeType"])
	d.Set("platform", object["Platform"])
	d.Set("resource_group_id", object["ResourceGroupId"])
	d.Set("start_time", object["StartTime"])
	d.Set("start_time_type", object["StartTimeType"])
	d.Set("status", object["Status"])

	if v, ok := object["AllocatedResources"]; ok {
		allocatedResources := v.(map[string]interface{})
		if v, ok := allocatedResources["AllocatedResource"]; ok && len(v.([]interface{})) > 0 {
			allocatedResourceMap := v.([]interface{})[0].(map[string]interface{})
			d.Set("instance_type", allocatedResourceMap["InstanceType"])
			d.Set("instance_amount", allocatedResourceMap["TotalAmount"])
			d.Set("zone_ids", []string{fmt.Sprint(allocatedResourceMap["zoneId"])})
		}
	}

	if v, ok := object["Tags"].(map[string]interface{}); ok {
		d.Set("tags", tagsToMap(v["Tag"]))
	}
	d.Set("time_slot", object["TimeSlot"])

	return nil
}

func resourceAlicloudEcsCapacityReservationUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	ecsService := EcsService{client}
	var err error
	update := false

	if d.HasChange("tags") {
		if err := ecsService.SetResourceTags(d, "capacityreservation"); err != nil {
			return WrapError(err)
		}
		d.SetPartial("tags")
	}

	request := map[string]interface{}{
		"RegionId": client.RegionId,
	}

	if d.HasChange("description") {
		update = true
		if v, ok := d.GetOk("description"); ok {
			request["Description"] = v
		}
	}
	if d.HasChange("end_time") {
		update = true
		if v, ok := d.GetOk("end_time"); ok {
			request["EndTime"] = v
		}
	}
	if d.HasChange("end_time_type") {
		update = true
		if v, ok := d.GetOk("end_time_type"); ok {
			request["EndTimeType"] = v
		}
	}

	if d.HasChange("instance_amount") {
		update = true
		request["InstanceAmount"] = d.Get("instance_amount")
	}
	if d.HasChange("platform") {
		update = true
		if v, ok := d.GetOk("platform"); ok {
			request["Platform"] = v
		}
	}
	request["PrivatePoolOptions.Id"] = d.Id()

	if d.HasChange("capacity_reservation_name") {
		update = true
		if v, ok := d.GetOk("capacity_reservation_name"); ok {
			request["PrivatePoolOptions.Name"] = v
		}
	}

	if update {
		action := "ModifyCapacityReservation"
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *resource.RetryError {
			resp, err := client.RpcPost("Ecs", "2014-05-26", action, nil, request, false)
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
	stateConf := BuildStateConf([]string{}, []string{"Active"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, ecsService.EcsCapacityReservationStateRefreshFunc(d, []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAlicloudEcsCapacityReservationRead(d, meta)
}

func resourceAlicloudEcsCapacityReservationDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var err error

	ecsService := EcsService{client}

	request := map[string]interface{}{
		"RegionId":              client.RegionId,
		"PrivatePoolOptions.Id": d.Id(),
	}

	if v, ok := d.GetOk("dry_run"); ok {
		request["DryRun"] = v
	}

	action := "ReleaseCapacityReservation"
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutDelete)), func() *resource.RetryError {
		resp, err := client.RpcPost("Ecs", "2014-05-26", action, nil, request, false)
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
		if IsExpectedErrors(err, []string{"Invalid.PrivatePoolOptions.Ids"}) || NotFoundError(err) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	stateConf := BuildStateConf([]string{}, []string{}, d.Timeout(schema.TimeoutDelete), 5*time.Second, ecsService.EcsCapacityReservationStateRefreshFunc(d, []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}
	return nil
}
