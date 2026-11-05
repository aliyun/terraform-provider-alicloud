// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package gpdb

import (
	"fmt"
	"log"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/convert"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"

	"github.com/blues/jsonata-go"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceSupabaseProject() *schema.Resource {
	return &schema.Resource{
		Create: resourceSupabaseProjectCreate,
		Read:   resourceSupabaseProjectRead,
		Update: resourceSupabaseProjectUpdate,
		Delete: resourceSupabaseProjectDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(15 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"account_password": {
				Type:      schema.TypeString,
				Required:  true,
				Sensitive: true,
			},
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"disk_performance_level": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"project_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"project_spec": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"region_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"security_ip_list": {
				Type:     schema.TypeSet,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"storage_size": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"vswitch_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"zone_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceSupabaseProjectCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "CreateSupabaseProject"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["RegionId"] = client.RegionId
	request["ClientToken"] = helper.BuildClientToken(action)

	request["ProjectName"] = d.Get("project_name")
	request["VpcId"] = d.Get("vpc_id")
	request["ProjectSpec"] = d.Get("project_spec")
	request["VSwitchId"] = d.Get("vswitch_id")
	securityIpListJsonPath, err := jsonpath.Get("$", d.Get("security_ip_list"))
	if err == nil {
		request["SecurityIPList"] = convert.ConvertListToCommaSeparate(convert.ConvertToInterfaceArray(securityIpListJsonPath))
	}
	if v, ok := d.GetOkExists("storage_size"); ok {
		request["StorageSize"] = v
	}
	request["ZoneId"] = d.Get("zone_id")
	request["AccountPassword"] = d.Get("account_password")
	if v, ok := d.GetOk("disk_performance_level"); ok {
		request["DiskPerformanceLevel"] = v
	}
	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		response, err = client.RpcPost("gpdb", "2016-05-03", action, query, request, true)
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
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_gpdb_supabase_project", action, errs.AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["ProjectId"]))

	gpdbServiceV2 := NewServiceV2(client)
	stateConf := helper.BuildStateConf([]string{}, []string{"running"}, d.Timeout(schema.TimeoutCreate), 5*time.Minute, gpdbServiceV2.SupabaseProjectStateRefreshFunc(d.Id(), "$.Status", []string{}))
	if jobDetail, err := stateConf.WaitForState(); err != nil {
		return errs.WrapErrorf(err, errs.IdMsg, d.Id(), jobDetail)
	}

	return resourceSupabaseProjectRead(d, meta)
}

func resourceSupabaseProjectRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	gpdbServiceV2 := NewServiceV2(client)

	objectRaw, err := gpdbServiceV2.DescribeSupabaseProject(d.Id())
	if err != nil {
		if !d.IsNewResource() && errs.NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_gpdb_supabase_project DescribeGpdbSupabaseProject Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return errs.WrapError(err)
	}

	d.Set("create_time", objectRaw["CreateTime"])
	d.Set("disk_performance_level", objectRaw["DiskPerformanceLevel"])
	d.Set("project_name", objectRaw["ProjectName"])
	d.Set("project_spec", objectRaw["ProjectSpec"])
	d.Set("region_id", objectRaw["RegionId"])
	d.Set("storage_size", objectRaw["StorageSize"])
	d.Set("vswitch_id", objectRaw["VSwitchId"])
	d.Set("vpc_id", objectRaw["VpcId"])
	d.Set("zone_id", objectRaw["ZoneId"])
	d.Set("status", objectRaw["Status"])

	e := jsonata.MustCompile("$split($.SecurityIpList, \",\")")
	evaluation, _ := e.Eval(objectRaw)
	d.Set("security_ip_list", evaluation)

	return nil
}

func resourceSupabaseProjectUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	update := false
	d.Partial(true)

	var err error
	action := "ResetSupabaseProjectPassword"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["ProjectId"] = d.Id()
	request["RegionId"] = client.RegionId
	if d.HasChange("account_password") {
		update = true
	}
	request["AccountPassword"] = d.Get("account_password")
	if update {
		wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
		err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
			response, err = client.RpcPost("gpdb", "2016-05-03", action, query, request, true)
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
	update = false
	action = "ModifySupabaseProjectSecurityIps"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["ProjectId"] = d.Id()
	request["RegionId"] = client.RegionId
	if d.HasChange("security_ip_list") {
		update = true
	}
	securityIpListJsonPath, err := jsonpath.Get("$", d.Get("security_ip_list"))
	if err == nil {
		request["SecurityIPList"] = convert.ConvertListToCommaSeparate(convert.ConvertToInterfaceArray(securityIpListJsonPath))
	}

	if update {
		wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
		err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
			response, err = client.RpcPost("gpdb", "2016-05-03", action, query, request, true)
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

	d.Partial(false)
	return resourceSupabaseProjectRead(d, meta)
}

func resourceSupabaseProjectDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	action := "DeleteSupabaseProject"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["ProjectId"] = d.Id()
	request["RegionId"] = client.RegionId

	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutDelete), func() *retry.RetryError {
		response, err = client.RpcPost("gpdb", "2016-05-03", action, query, request, true)
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
		if errs.NotFoundError(err) {
			return nil
		}
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
	}

	gpdbServiceV2 := NewServiceV2(client)
	stateConf := helper.BuildStateConf([]string{}, []string{""}, d.Timeout(schema.TimeoutDelete), 5*time.Second, gpdbServiceV2.SupabaseProjectStateRefreshFunc(d.Id(), "$.ProjectId", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return errs.WrapErrorf(err, errs.IdMsg, d.Id())
	}

	return nil
}
