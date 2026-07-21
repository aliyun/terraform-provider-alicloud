// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
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
)

func resourceAliCloudDmsAirflow() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudDmsAirflowCreate,
		Read:   resourceAliCloudDmsAirflowRead,
		Update: resourceAliCloudDmsAirflowUpdate,
		Delete: resourceAliCloudDmsAirflowDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"airflow_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"airflow_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"app_spec": {
				Type:     schema.TypeString,
				Required: true,
			},
			"dags_dir": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"description": {
				Type:     schema.TypeString,
				Required: true,
			},
			"oss_bucket_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"oss_path": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"plugins_dir": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"region_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"requirement_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"security_group_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"startup_file": {
				Type:     schema.TypeString,
				Optional: true,
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
			"worker_serverless_replicas": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"workspace_id": {
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

func resourceAliCloudDmsAirflowCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "CreateAirflow"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	if v, ok := d.GetOk("workspace_id"); ok {
		request["WorkspaceId"] = v
	}
	request["RegionId"] = client.RegionId
	request["ClientToken"] = buildClientToken(action)

	request["WorkerServerlessReplicas"] = d.Get("worker_serverless_replicas")
	request["Description"] = d.Get("description")
	request["SecurityGroupId"] = d.Get("security_group_id")
	request["VSwitchId"] = d.Get("vswitch_id")
	request["AppSpec"] = d.Get("app_spec")
	request["OssBucketName"] = d.Get("oss_bucket_name")
	if v, ok := d.GetOk("requirement_file"); ok {
		request["RequirementFile"] = v
	}
	request["AirflowName"] = d.Get("airflow_name")
	if v, ok := d.GetOk("startup_file"); ok {
		request["StartupFile"] = v
	}
	request["OssPath"] = d.Get("oss_path")
	request["VpcId"] = d.Get("vpc_id")
	if v, ok := d.GetOk("plugins_dir"); ok {
		request["PluginsDir"] = v
	}
	request["ZoneId"] = d.Get("zone_id")
	if v, ok := d.GetOk("dags_dir"); ok {
		request["DagsDir"] = v
	}
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RpcPost("Dms", "2025-04-14", action, query, request, true)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_dms_airflow", action, AlibabaCloudSdkGoERROR)
	}

	RootWorkspaceIdVar, _ := jsonpath.Get("$.Root.WorkspaceId", response)
	RootAirflowIdVar, _ := jsonpath.Get("$.Root.AirflowId", response)
	d.SetId(fmt.Sprintf("%v:%v", RootWorkspaceIdVar, RootAirflowIdVar))

	return resourceAliCloudDmsAirflowRead(d, meta)
}

func resourceAliCloudDmsAirflowRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	dmsServiceV2 := DmsServiceV2{client}

	objectRaw, err := dmsServiceV2.DescribeDmsAirflow(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_dms_airflow DescribeDmsAirflow Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("airflow_name", objectRaw["AirflowName"])
	d.Set("app_spec", objectRaw["AppSpec"])
	d.Set("dags_dir", objectRaw["DagsDir"])
	d.Set("description", objectRaw["Description"])
	d.Set("oss_bucket_name", objectRaw["OssBucketName"])
	d.Set("oss_path", objectRaw["OssPath"])
	d.Set("plugins_dir", objectRaw["PluginsDir"])
	d.Set("region_id", objectRaw["RegionId"])
	d.Set("requirement_file", objectRaw["RequirementFile"])
	d.Set("security_group_id", objectRaw["SecurityGroupId"])
	d.Set("startup_file", objectRaw["StartupFile"])
	d.Set("vswitch_id", objectRaw["VSwitchId"])
	d.Set("vpc_id", objectRaw["VpcId"])
	d.Set("worker_serverless_replicas", objectRaw["WorkerServerlessReplicas"])
	d.Set("zone_id", objectRaw["ZoneId"])
	d.Set("airflow_id", objectRaw["AirflowId"])
	d.Set("workspace_id", objectRaw["WorkspaceId"])

	return nil
}

func resourceAliCloudDmsAirflowUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	update := false

	var err error
	parts := strings.Split(d.Id(), ":")
	action := "UpdateAirflow"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["WorkspaceId"] = parts[0]
	request["AirflowId"] = parts[1]
	request["RegionId"] = client.RegionId
	request["ClientToken"] = buildClientToken(action)
	if d.HasChange("startup_file") {
		update = true
		request["StartupFile"] = d.Get("startup_file")
	}

	if d.HasChange("worker_serverless_replicas") {
		update = true
	}
	request["WorkerServerlessReplicas"] = d.Get("worker_serverless_replicas")
	if d.HasChange("description") {
		update = true
	}
	request["Description"] = d.Get("description")
	if d.HasChange("plugins_dir") {
		update = true
		request["PluginsDir"] = d.Get("plugins_dir")
	}

	if d.HasChange("app_spec") {
		update = true
	}
	request["AppSpec"] = d.Get("app_spec")
	if d.HasChange("requirement_file") {
		update = true
		request["RequirementFile"] = d.Get("requirement_file")
	}

	if d.HasChange("airflow_name") {
		update = true
	}
	request["AirflowName"] = d.Get("airflow_name")
	if d.HasChange("dags_dir") {
		update = true
		request["DagsDir"] = d.Get("dags_dir")
	}

	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("Dms", "2025-04-14", action, query, request, true)
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

	return resourceAliCloudDmsAirflowRead(d, meta)
}

func resourceAliCloudDmsAirflowDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	parts := strings.Split(d.Id(), ":")
	action := "DeleteAirflow"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["WorkspaceId"] = parts[0]
	request["AirflowId"] = parts[1]
	request["RegionId"] = client.RegionId
	request["ClientToken"] = buildClientToken(action)

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.RpcPost("Dms", "2025-04-14", action, query, request, true)
		request["ClientToken"] = buildClientToken(action)

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
		if IsExpectedErrors(err, []string{"InvalidAirflow.NotFound"}) || NotFoundError(err) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	return nil
}
