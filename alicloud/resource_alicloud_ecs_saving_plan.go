// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"log"
	"time"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudEcsSavingPlan() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudEcsSavingPlanCreate,
		Read:   resourceAliCloudEcsSavingPlanRead,
		Update: resourceAliCloudEcsSavingPlanUpdate,
		Delete: resourceAliCloudEcsSavingPlanDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"committed_amount": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"instance_family": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"offering_type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"payment_type": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"period": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
			},
			"plan_type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"purchase_method": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "family",
			},
		},
	}
}

func resourceAliCloudEcsSavingPlanCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	action := "PurchaseSavingPlanOffering"
	var request map[string]interface{}
	var response map[string]interface{}
	conn, err := client.NewEcsClient()
	if err != nil {
		return WrapError(err)
	}
	request = make(map[string]interface{})
	request["RegionId"] = client.RegionId
	request["ClientToken"] = buildClientToken(action)

	request["CommittedAmount"] = d.Get("committed_amount")
	if v, ok := d.GetOk("period"); ok {
		request["Period"] = v
	}
	request["PlanType"] = d.Get("plan_type")
	request["OfferingType"] = d.Get("offering_type")
	if v, ok := d.GetOk("instance_family"); ok {
		request["InstanceFamily"] = v
	}
	if v, ok := d.GetOk("instance_family_set"); ok {
		request["InstanceFamilySet"] = v
	}
	if v, ok := d.GetOk("period_unit"); ok {
		request["PeriodUnit"] = v
	}
	if v, ok := d.GetOk("saving_plan_name"); ok {
		request["Name"] = v
	}
	if v, ok := d.GetOk("purchase_method"); ok {
		request["PurchaseMethod"] = v
	}
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2016-03-14"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
		request["ClientToken"] = buildClientToken(action)

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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_ecs_saving_plan", action, AlibabaCloudSdkGoERROR)
	}

	id, _ := jsonpath.Get("$.SavingPlanIdSets[0]", response)
	d.SetId(fmt.Sprint(id))

	return resourceAliCloudEcsSavingPlanRead(d, meta)
}

func resourceAliCloudEcsSavingPlanRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ecsServiceV2 := EcsServiceV2{client}

	objectRaw, err := ecsServiceV2.DescribeEcsSavingPlan(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_ecs_saving_plan DescribeEcsSavingPlan Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("committed_amount", objectRaw["PoolValue"])
	d.Set("create_time", objectRaw["StartTimestamp"])
	d.Set("instance_family", objectRaw["InstanceFamily"])
	d.Set("offering_type", objectRaw["PayMode"])
	d.Set("period", convertEcsDataItemsCycleResponse(objectRaw["Cycle"]))
	d.Set("plan_type", objectRaw["SavingsType"])

	return nil
}

func resourceAliCloudEcsSavingPlanUpdate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Cannot update resource Alicloud Resource Saving Plan.")
	return nil
}

func resourceAliCloudEcsSavingPlanDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[WARN] Cannot destroy resource AliCloud Resource Saving Plan. Terraform will remove this resource from the state file, however resources may remain.")
	return nil
}

func convertEcsDataItemsCycleResponse(source interface{}) interface{} {
	switch source {
	case "1:Year":
		return "1"
	}
	return source
}
