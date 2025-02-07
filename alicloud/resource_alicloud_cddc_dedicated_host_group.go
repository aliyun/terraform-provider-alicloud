package alicloud

import (
	"fmt"
	"log"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAlicloudCddcDedicatedHostGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudCddcDedicatedHostGroupCreate,
		Read:   resourceAlicloudCddcDedicatedHostGroupRead,
		Update: resourceAlicloudCddcDedicatedHostGroupUpdate,
		Delete: resourceAlicloudCddcDedicatedHostGroupDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"allocation_policy": {
				Type:         schema.TypeString,
				Computed:     true,
				Optional:     true,
				ValidateFunc: StringInSlice([]string{"Evenly", "Intensively"}, false),
			},
			"cpu_allocation_ratio": {
				Type:         schema.TypeInt,
				Computed:     true,
				Optional:     true,
				ValidateFunc: IntBetween(100, 300),
			},
			"dedicated_host_group_desc": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"disk_allocation_ratio": {
				Type:         schema.TypeInt,
				Computed:     true,
				Optional:     true,
				ValidateFunc: IntBetween(100, 300),
			},
			"engine": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: StringInSlice([]string{"Redis", "SQLServer", "MySQL", "PostgreSQL", "MongoDB", "alisql", "tair", "mssql"}, false),
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
				ValidateFunc: IntBetween(0, 100),
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

func resourceAlicloudCddcDedicatedHostGroupCreate(d *schema.ResourceData, meta interface{}) error {
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
		request["OpenPermission"] = convertCddcDedicatedHostGroupOpenPermissionRequest(v.(bool))
	}
	request["Engine"] = d.Get("engine")
	if v, ok := d.GetOk("disk_allocation_ratio"); ok {
		//if d.Get("engine").(string) == "SQLServer" && v.(int) > 100 {
		//	return WrapError(fmt.Errorf("disk_allocation_ratio needs to be less than 100 under the SQLServer"))
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
	request["ClientToken"] = buildClientToken("CreateDedicatedHostGroup")
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RpcPost("cddc", "2020-03-20", action, nil, request, true)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_cddc_dedicated_host_group", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["DedicatedHostGroupId"]))

	return resourceAlicloudCddcDedicatedHostGroupRead(d, meta)
}
func resourceAlicloudCddcDedicatedHostGroupRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cddcService := CddcService{client}
	object, err := cddcService.DescribeCddcDedicatedHostGroup(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_cddc_dedicated_host_group cddcService.DescribeCddcDedicatedHostGroup Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	d.Set("allocation_policy", object["AllocationPolicy"])
	if v, ok := object["CpuAllocationRatio"]; ok && fmt.Sprint(v) != "0" {
		d.Set("cpu_allocation_ratio", formatInt(v))
	}
	d.Set("dedicated_host_group_desc", object["DedicatedHostGroupDesc"])
	if v, ok := object["DiskAllocationRatio"]; ok && fmt.Sprint(v) != "0" {
		d.Set("disk_allocation_ratio", formatInt(v))
	}

	d.Set("engine", convertCddcDedicatedHostGroupEngineResponse(object["Engine"]))
	d.Set("host_replace_policy", object["HostReplacePolicy"])
	if v, ok := object["MemAllocationRatio"]; ok && fmt.Sprint(v) != "0" {
		d.Set("mem_allocation_ratio", formatInt(v))
	}
	d.Set("vpc_id", object["VPCId"])
	d.Set("open_permission", convertCddcDedicatedHostGroupOpenPermissionResponse(formatInt(object["OpenPermission"])))
	return nil
}
func resourceAlicloudCddcDedicatedHostGroupUpdate(d *schema.ResourceData, meta interface{}) error {
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
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("cddc", "2020-03-20", action, nil, request, false)
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
	return resourceAlicloudCddcDedicatedHostGroupRead(d, meta)
}
func resourceAlicloudCddcDedicatedHostGroupDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "DeleteDedicatedHostGroup"
	var response map[string]interface{}
	var err error
	request := map[string]interface{}{
		"DedicatedHostGroupId": d.Id(),
	}

	request["RegionId"] = client.RegionId
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.RpcPost("cddc", "2020-03-20", action, nil, request, false)
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
	return nil
}

func convertCddcDedicatedHostGroupEngineResponse(engine interface{}) string {
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

func convertCddcDedicatedHostGroupOpenPermissionRequest(source interface{}) interface{} {
	switch source {
	case true:
		return 3
	case false:
		return 0
	}
	return 3
}

func convertCddcDedicatedHostGroupOpenPermissionResponse(source interface{}) interface{} {
	switch source {
	case 0:
		return false
	case 3:
		return true
	}
	return false
}
