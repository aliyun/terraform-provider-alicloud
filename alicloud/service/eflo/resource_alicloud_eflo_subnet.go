package eflo

import (
	"fmt"
	"log"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func ResourceSubnet() *schema.Resource {
	return &schema.Resource{
		Create: resourceSubnetCreate,
		Read:   resourceSubnetRead,
		Update: resourceSubnetUpdate,
		Delete: resourceSubnetDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"cidr": {
				Required: true,
				ForceNew: true,
				Type:     schema.TypeString,
			},
			"create_time": {
				Computed: true,
				Type:     schema.TypeString,
			},
			"gmt_modified": {
				Computed: true,
				Type:     schema.TypeString,
			},
			"message": {
				Computed: true,
				Type:     schema.TypeString,
			},
			"resource_group_id": {
				Computed: true,
				Type:     schema.TypeString,
			},
			"status": {
				Computed: true,
				Type:     schema.TypeString,
			},
			"subnet_id": {
				Computed: true,
				Type:     schema.TypeString,
			},
			"subnet_name": {
				Required: true,
				Type:     schema.TypeString,
			},
			"type": {
				Optional:     true,
				ForceNew:     true,
				Type:         schema.TypeString,
				ValidateFunc: validation.StringInSlice([]string{"OOB", "LB"}, false),
			},
			"vpd_id": {
				Required: true,
				ForceNew: true,
				Type:     schema.TypeString,
			},
			"zone_id": {
				Required: true,
				ForceNew: true,
				Type:     schema.TypeString,
			},
		},
	}
}

func resourceSubnetCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	efloService := EfloService{client}
	request := make(map[string]interface{})
	var err error

	if v, ok := d.GetOk("vpd_id"); ok {
		request["VpdId"] = v
	}
	request["SubnetName"] = d.Get("subnet_name")
	request["ZoneId"] = d.Get("zone_id")
	if v, ok := d.GetOk("cidr"); ok {
		request["Cidr"] = v
	}
	if v, ok := d.GetOk("type"); ok {
		request["Type"] = v
	}

	request["ClientToken"] = helper.BuildClientToken("CreateSubnet")
	var response map[string]interface{}
	action := "CreateSubnet"
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutCreate)), func() *retry.RetryError {
		resp, err := client.RpcPost("eflo", "2022-05-30", action, nil, request, false)
		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		response = resp
		logs.AddDebug(action, response, request)
		return nil
	})
	if err != nil {
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_eflo_subnet", action, errs.AlibabaCloudSdkGoERROR)
	}
	subnetIdValue, err := jsonpath.Get("$.Content.SubnetId", response)
	if err != nil || subnetIdValue == nil {
		return errs.WrapErrorf(err, errs.IdMsg, "alicloud_eflo_subnet")
	}

	d.SetId(fmt.Sprint(request["VpdId"], ":", subnetIdValue))
	stateConf := helper.BuildStateConf([]string{}, []string{"Available", "Waiting4AssignPort"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, efloService.EfloSubnetStateRefreshFunc(d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return errs.WrapErrorf(err, errs.IdMsg, d.Id())
	}
	return resourceSubnetRead(d, meta)
}

func resourceSubnetRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	efloService := EfloService{client}

	object, err := efloService.DescribeEfloSubnet(d.Id())
	if err != nil {
		if !d.IsNewResource() && errs.NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_eflo_subnet efloService.DescribeEfloSubnet Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return errs.WrapError(err)
	}
	parts, err := helper.ParseResourceId(d.Id(), 2)
	if err != nil {
		return errs.WrapError(err)
	}
	d.Set("vpd_id", parts[0])
	d.Set("subnet_id", parts[1])
	d.Set("cidr", object["Cidr"])
	d.Set("create_time", object["CreateTime"])
	d.Set("gmt_modified", object["GmtModified"])
	d.Set("message", object["Message"])
	d.Set("resource_group_id", object["ResourceGroupId"])
	d.Set("status", object["Status"])
	d.Set("subnet_name", object["SubnetName"])
	d.Set("type", object["Type"])
	d.Set("zone_id", object["ZoneId"])

	return nil
}

func resourceSubnetUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	var err error
	parts, err := helper.ParseResourceId(d.Id(), 2)
	if err != nil {
		return errs.WrapError(err)
	}
	update := false
	request := map[string]interface{}{
		"VpdId":    parts[0],
		"SubnetId": parts[1],
		"RegionId": client.RegionId,
	}

	if d.HasChange("subnet_name") {
		update = true
		request["SubnetName"] = d.Get("subnet_name")
	}
	request["ZoneId"] = d.Get("zone_id")

	if update {
		action := "UpdateSubnet"
		wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
		err = retry.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *retry.RetryError {
			resp, err := client.RpcPost("eflo", "2022-05-30", action, nil, request, false)
			if err != nil {
				if helper.NeedRetry(err) {
					wait()
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
			}
			logs.AddDebug(action, resp, request)
			return nil
		})
		if err != nil {
			return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
		}
	}

	return resourceSubnetRead(d, meta)
}

func resourceSubnetDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	efloService := EfloService{client}
	var response map[string]interface{}
	var err error
	parts, err := helper.ParseResourceId(d.Id(), 2)
	if err != nil {
		return errs.WrapError(err)
	}

	request := map[string]interface{}{
		"VpdId":    parts[0],
		"SubnetId": parts[1],
		"ZoneId":   d.Get("zone_id"),
	}

	action := "DeleteSubnet"
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutDelete)), func() *retry.RetryError {
		response, err = client.RpcPost("eflo", "2022-05-30", action, nil, request, false)
		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		logs.AddDebug(action, response, request)
		return nil
	})
	if err != nil {
		if errs.NotFoundError(err) || errs.IsExpectedErrors(err, []string{"1003"}) {
			return nil
		}
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
	}
	stateConf := helper.BuildStateConf([]string{}, []string{}, d.Timeout(schema.TimeoutDelete), 5*time.Second, efloService.EfloSubnetStateRefreshFunc(d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return errs.WrapErrorf(err, errs.IdMsg, d.Id())
	}
	return nil
}
