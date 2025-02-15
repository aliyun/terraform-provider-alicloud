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

func resourceAlicloudEcsElasticityAssurance() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudEcsElasticityAssuranceCreate,
		Read:   resourceAlicloudEcsElasticityAssuranceRead,
		Update: resourceAlicloudEcsElasticityAssuranceUpdate,
		Delete: resourceAlicloudEcsElasticityAssuranceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"assurance_times": {
				ForceNew:     true,
				Optional:     true,
				Computed:     true,
				Type:         schema.TypeString,
				ValidateFunc: validation.StringInSlice([]string{"Unlimited"}, false),
			},
			"description": {
				Optional: true,
				Type:     schema.TypeString,
			},
			"elasticity_assurance_id": {
				Computed: true,
				Type:     schema.TypeString,
			},
			"end_time": {
				Computed: true,
				Type:     schema.TypeString,
			},
			"instance_amount": {
				Required:     true,
				ForceNew:     true,
				Type:         schema.TypeInt,
				ValidateFunc: validation.IntBetween(1, 1000),
			},
			"instance_charge_type": {
				Computed: true,
				Type:     schema.TypeString,
			},
			"instance_type": {
				Required: true,
				ForceNew: true,
				MaxItems: 1,
				Type:     schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"period": {
				Optional:     true,
				ValidateFunc: validation.IntInSlice([]int{1, 10, 11, 12, 2, 3, 4, 5, 6, 7, 8, 9}),
				Type:         schema.TypeInt,
			},
			"period_unit": {
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"Month", "Year"}, false),
				Type:         schema.TypeString,
			},
			"private_pool_options_match_criteria": {
				Optional:     true,
				ForceNew:     true,
				Computed:     true,
				Type:         schema.TypeString,
				ValidateFunc: validation.StringInSlice([]string{"Open", "Target"}, false),
			},
			"private_pool_options_name": {
				Optional: true,
				ForceNew: true,
				Computed: true,
				Type:     schema.TypeString,
			},
			"resource_group_id": {
				Optional: true,
				ForceNew: true,
				Type:     schema.TypeString,
			},
			"start_time": {
				Optional: true,
				Computed: true,
				ForceNew: true,
				Type:     schema.TypeString,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return old != "" && new != "" && strings.HasPrefix(new, strings.Trim(old, "00Z"))
				},
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
			"used_assurance_times": {
				Computed: true,
				Type:     schema.TypeInt,
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

func resourceAlicloudEcsElasticityAssuranceCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ecsService := EcsService{client}
	request := map[string]interface{}{
		"RegionId": client.RegionId,
	}
	var err error

	if v, ok := d.GetOk("assurance_times"); ok {
		request["AssuranceTimes"] = v
	}
	if v, ok := d.GetOk("description"); ok {
		request["Description"] = v
	}
	if v, ok := d.GetOk("instance_amount"); ok {
		request["InstanceAmount"] = v
	}
	if v, ok := d.GetOk("instance_type"); ok {
		request["InstanceType"] = v.([]interface{})
	}
	if v, ok := d.GetOk("period"); ok {
		request["Period"] = v
	}
	if v, ok := d.GetOk("period_unit"); ok {
		request["PeriodUnit"] = v
	}
	if v, ok := d.GetOk("private_pool_options_match_criteria"); ok {
		request["PrivatePoolOptions.MatchCriteria"] = v
	}
	if v, ok := d.GetOk("private_pool_options_name"); ok {
		request["PrivatePoolOptions.Name"] = v
	}
	if v, ok := d.GetOk("resource_group_id"); ok {
		request["ResourceGroupId"] = v
	}
	if v, ok := d.GetOk("start_time"); ok {
		request["StartTime"] = v
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

	request["ClientToken"] = buildClientToken("CreateElasticityAssurance")
	var response map[string]interface{}
	action := "CreateElasticityAssurance"
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_ecs_elasticity_assurance", action, AlibabaCloudSdkGoERROR)
	}

	if v, err := jsonpath.Get("$.PrivatePoolOptionsId", response); err != nil || v == nil {
		return WrapErrorf(err, IdMsg, "alicloud_ecs_elasticity_assurance")
	} else {
		d.SetId(fmt.Sprint(v))
	}
	stateConf := BuildStateConf([]string{}, []string{"Active", "Prepared"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, ecsService.EcsElasticityAssuranceStateRefreshFunc(d, []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}
	return resourceAlicloudEcsElasticityAssuranceRead(d, meta)
}

func resourceAlicloudEcsElasticityAssuranceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ecsService := EcsService{client}

	object, err := ecsService.DescribeEcsElasticityAssurance(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_ecs_elasticity_assurance ecsService.DescribeEcsElasticityAssurance Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	d.Set("description", object["Description"])
	d.Set("elasticity_assurance_id", object["PrivatePoolOptionsId"])
	d.Set("end_time", object["EndTime"])
	d.Set("instance_charge_type", object["InstanceChargeType"])
	d.Set("private_pool_options_match_criteria", object["PrivatePoolOptionsMatchCriteria"])
	d.Set("private_pool_options_name", object["PrivatePoolOptionsName"])
	d.Set("resource_group_id", object["ResourceGroupId"])
	d.Set("start_time", object["StartTime"])
	d.Set("start_time_type", object["StartTimeType"])
	d.Set("status", object["Status"])
	if v, ok := object["Tags"].(map[string]interface{}); ok {
		d.Set("tags", tagsToMap(v["Tag"]))
	}
	if v, ok := object["AllocatedResources"]; ok {
		allocatedResources := v.(map[string]interface{})
		if v, ok := allocatedResources["AllocatedResource"]; ok && len(v.([]interface{})) > 0 {
			allocatedResourceMap := v.([]interface{})[0].(map[string]interface{})
			d.Set("instance_type", []string{fmt.Sprint(allocatedResourceMap["InstanceType"])})
			d.Set("instance_amount", allocatedResourceMap["TotalAmount"])
			d.Set("zone_ids", []string{fmt.Sprint(allocatedResourceMap["zoneId"])})
		}
	}

	d.Set("assurance_times", object["TotalAssuranceTimes"])
	d.Set("used_assurance_times", object["UsedAssuranceTimes"])

	return nil
}

func resourceAlicloudEcsElasticityAssuranceUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	ecsService := EcsService{client}
	var err error
	update := false
	request := map[string]interface{}{
		"RegionId":              client.RegionId,
		"PrivatePoolOptions.Id": d.Id(),
	}

	if d.HasChange("tags") {
		if err := ecsService.SetResourceTags(d, "elasticityassurance"); err != nil {
			return WrapError(err)
		}
		d.SetPartial("tags")
	}

	if d.HasChange("description") {
		update = true
		if v, ok := d.GetOk("description"); ok {
			request["Description"] = v
		}
	}

	if d.HasChange("private_pool_options_name") {
		update = true
		if v, ok := d.GetOk("private_pool_options_name"); ok {
			request["PrivatePoolOptions.Name"] = v
		}
	}

	if update {
		action := "ModifyElasticityAssurance"
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

	return resourceAlicloudEcsElasticityAssuranceRead(d, meta)
}

func resourceAlicloudEcsElasticityAssuranceDelete(d *schema.ResourceData, meta interface{}) error {
	return nil
}
