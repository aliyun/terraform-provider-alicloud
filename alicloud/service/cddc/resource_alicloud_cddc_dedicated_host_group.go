package cddc

import (
	"fmt"
	"log"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/convert"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/validators"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceDedicatedHostGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceDedicatedHostGroupCreate,
		Read:   resourceDedicatedHostGroupRead,
		Update: resourceDedicatedHostGroupUpdate,
		Delete: resourceDedicatedHostGroupDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		DeprecationMessage: "This resource has been deprecated since v1.225.1 and will be removed in the future.",
		Schema: map[string]*schema.Schema{
			"allocation_policy": {
				Type:         schema.TypeString,
				Computed:     true,
				Optional:     true,
				ValidateFunc: validators.StringInSlice([]string{"Evenly", "Intensively"}, false),
			},
			"cpu_allocation_ratio": {
				Type:         schema.TypeInt,
				Computed:     true,
				Optional:     true,
				ValidateFunc: validators.IntBetween(100, 300),
			},
			"dedicated_host_group_desc": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"disk_allocation_ratio": {
				Type:         schema.TypeInt,
				Computed:     true,
				Optional:     true,
				ValidateFunc: validators.IntBetween(100, 300),
			},
			"engine": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validators.StringInSlice([]string{"Redis", "SQLServer", "MySQL", "PostgreSQL", "MongoDB", "alisql", "tair", "mssql"}, false),
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return d.Id() != "" && new == "SQLServer"
				},
			},
			"host_replace_policy": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			"mem_allocation_ratio": {
				Type:         schema.TypeInt,
				Computed:     true,
				Optional:     true,
				ValidateFunc: validators.IntBetween(0, 100),
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"open_permission": {
				Type:     schema.TypeBool,
				Computed: true,
				Optional: true,
				ForceNew: true,
			},
		},
	}
}

func resourceDedicatedHostGroupCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "CreateDedicatedHostGroup"
	request := make(map[string]interface{})
	var err error
	if v, ok := d.GetOk("allocation_policy"); ok {
		request["AllocationPolicy"] = v
	}
	if v, ok := d.GetOk("cpu_allocation_ratio"); ok {
		request["CpuAllocationRatio"] = v
	}
	if v, ok := d.GetOk("dedicated_host_group_desc"); ok {
		request["DedicatedHostGroupDesc"] = v
	}
	if v, ok := d.GetOk("open_permission"); ok {
		request["OpenPermission"] = convertDedicatedHostGroupOpenPermissionRequest(v.(bool))
	}
	request["Engine"] = d.Get("engine")
	if v, ok := d.GetOk("disk_allocation_ratio"); ok {
		//if d.Get("engine").(string) == "SQLServer" && v.(int) > 100 {
		//	return errs.WrapError(fmt.Errorf("disk_allocation_ratio needs to be less than 100 under the SQLServer"))
		//}
		request["DiskAllocationRatio"] = v
	}
	if v, ok := d.GetOk("host_replace_policy"); ok {
		request["HostReplacePolicy"] = v
	}
	if v, ok := d.GetOk("mem_allocation_ratio"); ok {
		request["MemAllocationRatio"] = v
	}
	request["RegionId"] = client.RegionId
	request["VPCId"] = d.Get("vpc_id")
	request["ClientToken"] = helper.BuildClientToken("CreateDedicatedHostGroup")
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		response, err = client.RpcPost("cddc", "2020-03-20", action, nil, request, true)
		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(action, response, request)
	if err != nil {
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_cddc_dedicated_host_group", action, errs.AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["DedicatedHostGroupId"]))

	return resourceDedicatedHostGroupRead(d, meta)
}
func resourceDedicatedHostGroupRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cddcService := NewService(client)
	object, err := cddcService.DescribeDedicatedHostGroup(d.Id())
	if err != nil {
		if errs.NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_cddc_dedicated_host_group cddcService.DescribeDedicatedHostGroup Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return errs.WrapError(err)
	}
	d.Set("allocation_policy", object["AllocationPolicy"])
	if v, ok := object["CpuAllocationRatio"]; ok && fmt.Sprint(v) != "0" {
		cpu_allocation_ratioVal, err := convert.FormatInt(v)
		if err != nil {
			return err
		}
		d.Set("cpu_allocation_ratio", cpu_allocation_ratioVal)
	}
	d.Set("dedicated_host_group_desc", object["DedicatedHostGroupDesc"])
	if v, ok := object["DiskAllocationRatio"]; ok && fmt.Sprint(v) != "0" {
		disk_allocation_ratioVal, err := convert.FormatInt(v)
		if err != nil {
			return err
		}
		d.Set("disk_allocation_ratio", disk_allocation_ratioVal)
	}

	d.Set("engine", convertDedicatedHostGroupEngineResponse(object["Engine"]))
	d.Set("host_replace_policy", object["HostReplacePolicy"])
	if v, ok := object["MemAllocationRatio"]; ok && fmt.Sprint(v) != "0" {
		mem_allocation_ratioVal, err := convert.FormatInt(v)
		if err != nil {
			return err
		}
		d.Set("mem_allocation_ratio", mem_allocation_ratioVal)
	}
	d.Set("vpc_id", object["VPCId"])
	openPermission, err := convert.FormatInt(object["OpenPermission"])
	if err != nil {
		return err
	}
	d.Set("open_permission", convertDedicatedHostGroupOpenPermissionResponse(openPermission))
	return nil
}
func resourceDedicatedHostGroupUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	var err error
	update := false
	request := map[string]interface{}{
		"DedicatedHostGroupId": d.Id(),
	}
	request["RegionId"] = client.RegionId
	if d.HasChange("allocation_policy") {
		update = true
		if v, ok := d.GetOk("allocation_policy"); ok {
			request["AllocationPolicy"] = v
		}
	}
	if d.HasChange("cpu_allocation_ratio") {
		update = true
		if v, ok := d.GetOk("cpu_allocation_ratio"); ok {
			request["CpuAllocationRatio"] = v
		}
	}
	if d.HasChange("dedicated_host_group_desc") {
		update = true
		if v, ok := d.GetOk("dedicated_host_group_desc"); ok {
			request["DedicatedHostGroupDesc"] = v
		}
	}
	if d.HasChange("disk_allocation_ratio") {
		update = true
		if v, ok := d.GetOk("disk_allocation_ratio"); ok {
			request["DiskAllocationRatio"] = v
		}
	}
	if d.HasChange("host_replace_policy") {
		update = true
		if v, ok := d.GetOk("host_replace_policy"); ok {
			request["HostReplacePolicy"] = v
		}
	}
	if d.HasChange("mem_allocation_ratio") {
		update = true
		if v, ok := d.GetOk("mem_allocation_ratio"); ok {
			request["MemAllocationRatio"] = v
		}
	}
	if update {
		action := "ModifyDedicatedHostGroupAttribute"
		wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
		err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
			response, err = client.RpcPost("cddc", "2020-03-20", action, nil, request, false)
			if err != nil {
				if helper.NeedRetry(err) {
					wait()
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
			}
			return nil
		})
		logs.AddDebug(action, response, request)
		if err != nil {
			return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
		}
	}
	return resourceDedicatedHostGroupRead(d, meta)
}
func resourceDedicatedHostGroupDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "DeleteDedicatedHostGroup"
	var response map[string]interface{}
	var err error
	request := map[string]interface{}{
		"DedicatedHostGroupId": d.Id(),
	}

	request["RegionId"] = client.RegionId
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutDelete), func() *retry.RetryError {
		response, err = client.RpcPost("cddc", "2020-03-20", action, nil, request, false)
		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(action, response, request)
	if err != nil {
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
	}
	return nil
}

func convertDedicatedHostGroupEngineResponse(engine interface{}) string {
	switch engine {
	case "mysql":
		engine = "MySQL"
	case "redis":
		engine = "Redis"
	case "pgsql":
		engine = "PostgreSQL"
	case "mongodb":
		engine = "MongoDB"
	}
	return fmt.Sprint(engine)
}

func convertDedicatedHostGroupOpenPermissionRequest(source interface{}) interface{} {
	switch source {
	case true:
		return 3
	case false:
		return 0
	}
	return 3
}

func convertDedicatedHostGroupOpenPermissionResponse(source interface{}) interface{} {
	switch source {
	case 0:
		return false
	case 3:
		return true
	}
	return false
}
