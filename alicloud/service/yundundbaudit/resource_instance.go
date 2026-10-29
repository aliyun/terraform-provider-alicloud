package yundundbaudit

import (
	"fmt"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/tags"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	log "github.com/sirupsen/logrus"
)

func ResourceInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceInstanceCreate,
		Read:   resourceInstanceRead,
		Update: resourceInstanceUpdate,
		Delete: resourceInstanceDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(20 * time.Minute),
			Update: schema.DefaultTimeout(20 * time.Minute),
		},

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"description": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 64),
			},
			"plan_code": {
				Type:     schema.TypeString,
				Required: true,
			},
			"period": {
				Type:         schema.TypeInt,
				ValidateFunc: validation.IntInSlice([]int{1, 3, 6, 12, 24, 36}),
				Required:     true,
			},
			"vswitch_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"tags": tags.TagsSchema(),

			"resource_group_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "CreateInstance"
	var request map[string]interface{}
	var response map[string]interface{}
	var err error
	var endpoint string
	query := make(map[string]interface{})
	request = make(map[string]interface{})
	request["ClientToken"] = helper.BuildClientToken(action)
	request["ProductCode"] = "dbaudit"
	request["SubscriptionType"] = "Subscription"
	request["Period"] = d.Get("period")
	parameterMapList := make([]map[string]interface{}, 0)
	parameterMapList = append(parameterMapList, map[string]interface{}{
		"Code":  "NetworkType",
		"Value": "vpc",
	})
	parameterMapList = append(parameterMapList, map[string]interface{}{
		"Code":  "SeriesCode",
		"Value": "alpha",
	})

	parameterMapList = append(parameterMapList, map[string]interface{}{
		"Code":  "PlanCode",
		"Value": d.Get("plan_code").(string),
	})
	parameterMapList = append(parameterMapList, map[string]interface{}{
		"Code":  "RegionId",
		"Value": client.RegionId,
	})
	request["Parameter"] = parameterMapList
	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		response, err = client.RpcPostWithEndpoint("BssOpenApi", "2017-12-14", action, query, request, true, endpoint)
		request["ClientToken"] = helper.BuildClientToken(action)
		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			if !client.IsInternationalAccount() && errs.IsExpectedErrors(err, []string{"NotApplicable"}) {
				endpoint = connectivity.BssOpenAPIEndpointInternational
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		logs.AddDebug(action, response, request)
		return nil
	})

	if err != nil {
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_yundun_dbaudit_instance", action, errs.AlibabaCloudSdkGoERROR)
	}

	id, _ := jsonpath.Get("$.Data.InstanceId", response)
	d.SetId(fmt.Sprint(id))

	dbauditService := NewService(client)

	// check RAM policy
	dbauditService.ProcessRolePolicy()
	// wait for order complete
	stateConf := helper.BuildStateConf([]string{}, []string{"PENDING"}, d.Timeout(schema.TimeoutCreate), 20*time.Second, dbauditService.DbauditInstanceRefreshFunc(d.Id(), []string{"UPGRADING", "UPGRADE_FAILED", "CREATE_FAILED"}))
	if _, err := stateConf.WaitForState(); err != nil {
		return errs.WrapErrorf(err, errs.IdMsg, d.Id())
	}
	// start instance
	if err := dbauditService.StartDbauditInstance(d.Id(), d.Get("vswitch_id").(string)); err != nil {
		return errs.WrapError(err)
	}
	// wait for pending
	stateConf = helper.BuildStateConf([]string{"PENDING", "CREATING"}, []string{"RUNNING"}, d.Timeout(schema.TimeoutCreate), 20*time.Second, dbauditService.DbauditInstanceRefreshFunc(d.Id(), []string{"UPGRADING", "UPGRADE_FAILED", "CREATE_FAILED"}))
	if _, err := stateConf.WaitForState(); err != nil {
		return errs.WrapErrorf(err, errs.IdMsg, d.Id())
	}
	return resourceInstanceUpdate(d, meta)
}

func resourceInstanceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	dbauditService := NewService(client)
	instance, err := dbauditService.DescribeYundunDbauditInstance(d.Id())
	if err != nil {
		if !d.IsNewResource() && errs.NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_yundun_dbaudit_instance DescribeInstanceAttribute Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return errs.WrapError(err)
	}
	d.Set("description", instance["Description"])
	d.Set("plan_code", instance["LicenseCode"])
	d.Set("vswitch_id", instance["VswitchId"])

	tagsMaps, err := dbauditService.DescribeTags(d.Id(), "INSTANCE")
	if err != nil {
		return errs.WrapError(err)
	}
	d.Set("tags", tags.TagsToMap(tagsMaps))
	return nil
}

func resourceInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	dbauditService := NewService(client)

	if d.HasChange("tags") {
		if err := dbauditService.setInstanceTags(d, "INSTANCE"); err != nil {
			return errs.WrapError(err)
		}
	}

	if d.HasChange("description") {
		if err := dbauditService.UpdateDbauditInstanceDescription(d.Id(), d.Get("description").(string)); err != nil {
			return errs.WrapError(err)
		}
	}

	if d.HasChange("resource_group_id") {
		if err := dbauditService.UpdateResourceGroup(d.Id(), d.Get("resource_group_id").(string)); err != nil {
			return errs.WrapError(err)
		}
	}

	if d.IsNewResource() {
		return resourceInstanceRead(d, meta)
	}

	if d.HasChange("plan_code") {
		if err := dbauditService.UpdateInstanceSpec("plan_code", "PlanCode", d); err != nil {
			return errs.WrapError(err)
		}
		stateConf := helper.BuildStateConf([]string{"UPGRADING"}, []string{"PENDING", "RUNNING"}, d.Timeout(schema.TimeoutUpdate), 20*time.Second, dbauditService.DbauditInstanceRefreshFunc(d.Id(), []string{"CREATING", "UPGRADE_FAILED", "CREATE_FAILED"}))
		if _, err := stateConf.WaitForState(); err != nil {
			return errs.WrapErrorf(err, errs.IdMsg, d.Id())
		}
	}

	// wait for order complete
	return resourceInstanceRead(d, meta)
}

func resourceInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	return errs.WrapError(errs.Error("'Subscription' instance cannot be destroyed: please unsubscribe it in the console or wait for it to expire. To remove it from the state only, run 'terraform state rm'."))
}
