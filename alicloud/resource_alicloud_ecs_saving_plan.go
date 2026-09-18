package alicloud

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceAlicloudEcsSavingPlan() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudEcsSavingPlanCreate,
		Read:   resourceAlicloudEcsSavingPlanRead,
		Update: resourceAlicloudEcsSavingPlanUpdate,
		Delete: resourceAlicloudEcsSavingPlanDelete,
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
			"plan_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Default:      "EcsCompute",
				ValidateFunc: validation.StringInSlice([]string{"EcsCompute", "General"}, false),
			},
			"offering_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Default:      "AllUpfront",
				ValidateFunc: validation.StringInSlice([]string{"AllUpforward", "HalfUpfront", "NoUpfront"}, false),
			},
			"purchase_method": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Default:      "ByInstanceFamily",
				ValidateFunc: validation.StringInSlice([]string{"ByInstanceFamily", "ByInstanceFamilySet"}, false),
			},
			"charge_type": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Default:  "PrePaid",
			},
			"period": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
			},
			"period_unit": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"instance_family": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"instance_family_set": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"saving_plan_name": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"start_time": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"region_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"saving_plan_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"payment_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAlicloudEcsSavingPlanCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	request := map[string]interface{}{
		"RegionId": client.RegionId,
	}

	if v, ok := d.GetOk("committed_amount"); ok {
		request["CommittedAmount"] = v
	}
	if v, ok := d.GetOk("plan_type"); ok {
		request["PlanType"] = v
	}
	if v, ok := d.GetOk("offering_type"); ok {
		request["OfferingType"] = v
	}
	if v, ok := d.GetOk("purchase_method"); ok {
		request["PurchaseMethod"] = v
	}
	if v, ok := d.GetOk("charge_type"); ok {
		request["ChargeType"] = v
	}
	if v, ok := d.GetOk("period"); ok {
		request["Period"] = v
	}
	if v, ok := d.GetOk("period_unit"); ok {
		request["PeriodUnit"] = v
	}
	if v, ok := d.GetOk("instance_family"); ok {
		request["InstanceFamily"] = v
	}
	if v, ok := d.GetOk("instance_family_set"); ok {
		request["InstanceFamilySet"] = v
	}
	if v, ok := d.GetOk("saving_plan_name"); ok {
		request["Name"] = v
	}
	if v, ok := d.GetOk("description"); ok {
		request["Description"] = v
	}
	if v, ok := d.GetOk("start_time"); ok {
		request["StartTime"] = v
	}

	request["ClientToken"] = buildClientToken("PurchaseSavingPlanOffering")
	action := "PurchaseSavingPlanOffering"
	var response map[string]interface{}
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err := resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutCreate)), func() *resource.RetryError {
		resp, err := client.RpcPost("Ecs", "2016-03-14", action, nil, request, true)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_ecs_saving_plan", action, AlibabaCloudSdkGoERROR)
	}

	v, err := jsonpath.Get("$.SavingPlanIdSets", response)
	if err != nil || v == nil {
		return WrapErrorf(err, IdMsg, "alicloud_ecs_saving_plan")
	}
	idSets, ok := v.([]interface{})
	if !ok || len(idSets) == 0 {
		return WrapErrorf(fmt.Errorf("SavingPlanIdSets is empty"), IdMsg, "alicloud_ecs_saving_plan")
	}
	d.SetId(fmt.Sprint(idSets[0]))

	return resourceAlicloudEcsSavingPlanRead(d, meta)
}

func resourceAlicloudEcsSavingPlanRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ecsService := EcsService{client}

	object, err := ecsService.DescribeEcsSavingPlan(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_ecs_saving_plan ecsService.DescribeEcsSavingPlan Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("saving_plan_id", d.Id())
	if v, ok := object["StartTimestamp"]; ok {
		d.Set("create_time", v)
	}
	if v, ok := object["Region"]; ok {
		d.Set("region_id", v)
	}
	if v, ok := object["PoolValue"]; ok {
		d.Set("committed_amount", fmt.Sprint(v))
	}
	if v, ok := object["InstanceFamily"]; ok {
		d.Set("instance_family", v)
	}
	if v, ok := object["PayMode"]; ok {
		d.Set("offering_type", v)
	}
	if v, ok := object["SavingsType"]; ok {
		d.Set("plan_type", v)
	}
	if v, ok := object["StartTime"]; ok {
		d.Set("start_time", v)
	}
	if v, ok := object["Cycle"]; ok {
		if cycle, ok := v.(string); ok && cycle != "" {
			if period, err := parseSavingPlanCycle(cycle); err == nil {
				d.Set("period", period)
			}
		}
	}
	if v, ok := object["PaymentType"]; ok {
		d.Set("payment_type", v)
	}

	return nil
}

func resourceAlicloudEcsSavingPlanUpdate(d *schema.ResourceData, meta interface{}) error {
	// Saving Plan does not support update; all attributes are ForceNew.
	return resourceAlicloudEcsSavingPlanRead(d, meta)
}

func resourceAlicloudEcsSavingPlanDelete(d *schema.ResourceData, meta interface{}) error {
	// Saving Plans are financial commitments and cannot be cancelled via API.
	// Remove the resource from Terraform state without destroying the cloud resource.
	log.Printf("[DEBUG] Resource alicloud_ecs_saving_plan Delete: removing from state only. The Saving Plan resource %s remains on the cloud and must be cancelled manually if needed.", d.Id())
	d.SetId("")
	return nil
}

func parseSavingPlanCycle(cycle string) (int, error) {
	// Cycle values returned by QuerySavingsPlansInstance follow the format "N:Unit" (e.g. "1:Year").
	parts := strings.SplitN(cycle, ":", 2)
	return strconv.Atoi(parts[0])
}
