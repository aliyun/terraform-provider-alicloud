package alicloud

import (
	"fmt"
	"log"
	"time"

	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAlicloudOceanBaseDatabase() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudOceanBaseDatabaseCreate,
		Read:   resourceAlicloudOceanBaseDatabaseRead,
		Update: resourceAlicloudOceanBaseDatabaseUpdate,
		Delete: resourceAlicloudOceanBaseDatabaseDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(1 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(1 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"collation": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"database_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"encoding": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tenant_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"users": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceAlicloudOceanBaseDatabaseCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "CreateDatabase"
	request := make(map[string]interface{})
	conn, err := client.NewOceanbaseClient()
	if err != nil {
		return WrapError(err)
	}
	if v, ok := d.GetOk("collation"); ok {
		request["Collation"] = v
	}
	request["DatabaseName"] = d.Get("database_name")
	if v, ok := d.GetOk("description"); ok {
		request["Description"] = v
	}
	request["Encoding"] = d.Get("encoding")
	if v, ok := d.GetOk("instance_id"); ok {
		request["InstanceId"] = v
	}
	request["RegionId"] = client.RegionId
	request["TenantId"] = d.Get("tenant_id")
	request["ClientToken"] = buildClientToken("CreateDatabase")
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutCreate)), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-09-01"), StringPointer("AK"), nil, request, &runtime)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_ocean_base_database", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(request["TenantId"], ":", request["DatabaseName"]))

	return resourceAlicloudOceanBaseDatabaseUpdate(d, meta)
}
func resourceAlicloudOceanBaseDatabaseRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	oceanBaseProService := OceanBaseProService{client}
	object, err := oceanBaseProService.DescribeOceanBaseDatabase(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_ocean_base_database oceanBaseProService.DescribeOceanBaseDatabase Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	d.Set("database_name", parts[1])
	d.Set("tenant_id", parts[0])
	d.Set("collation", object["Collation"])
	d.Set("description", object["Description"])
	d.Set("encoding", object["Encoding"])
	d.Set("status", object["Status"])
	return nil
}
func resourceAlicloudOceanBaseDatabaseUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	d.Partial(true)

	update := false
	request := map[string]interface{}{
		"DatabaseName": parts[1],
		"TenantId":     parts[0],
	}
	if !d.IsNewResource() && d.HasChange("description") {
		update = true
	}
	if v, ok := d.GetOk("description"); ok {
		request["Description"] = v
	}

	request["InstanceId"] = d.Get("instance_id")
	request["RegionId"] = client.RegionId
	if update {
		action := "ModifyDatabaseDescription"
		conn, err := client.NewOceanbaseClient()
		if err != nil {
			return WrapError(err)
		}
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-09-01"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
		d.SetPartial("description")
		d.SetPartial("instance_id")
	}
	update = false
	modifyDatabaseUserRolesReq := map[string]interface{}{
		"DatabaseName": parts[1],
		"TenantId":     parts[0],
	}
	if d.HasChange("users") {
		update = true
	}
	if v, ok := d.GetOk("users"); ok {
		modifyDatabaseUserRolesReq["Users"] = v
	}

	modifyDatabaseUserRolesReq["InstanceId"] = d.Get("instance_id")
	modifyDatabaseUserRolesReq["RegionId"] = client.RegionId
	if update {
		action := "ModifyDatabaseUserRoles"
		conn, err := client.NewOceanbaseClient()
		if err != nil {
			return WrapError(err)
		}
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-09-01"), StringPointer("AK"), nil, modifyDatabaseUserRolesReq, &util.RuntimeOptions{})
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		addDebug(action, response, modifyDatabaseUserRolesReq)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
		d.SetPartial("users")
		d.SetPartial("instance_id")
	}
	d.Partial(false)
	return resourceAlicloudOceanBaseDatabaseRead(d, meta)
}
func resourceAlicloudOceanBaseDatabaseDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	oceanBaseProService := OceanBaseProService{client}
	action := "DeleteDatabases"
	var response map[string]interface{}
	conn, err := client.NewOceanbaseClient()
	if err != nil {
		return WrapError(err)
	}
	request := map[string]interface{}{
		"TenantId":      parts[0],
		"DatabaseNames": "[\"" + parts[1] + "\"]",
	}

	if v, ok := d.GetOk("instance_id"); ok {
		request["InstanceId"] = v
	}
	request["RegionId"] = client.RegionId
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutDelete)), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-09-01"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
	stateConf := BuildStateConf([]string{}, []string{}, d.Timeout(schema.TimeoutDelete), 5*time.Second, oceanBaseProService.OceanBaseDatabaseStateRefreshFunc(d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}
	return nil
}
