// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudAdbResourceGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudAdbResourceGroupCreate,
		Read:   resourceAliCloudAdbResourceGroupRead,
		Update: resourceAliCloudAdbResourceGroupUpdate,
		Delete: resourceAliCloudAdbResourceGroupDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(31 * time.Minute),
			Update: schema.DefaultTimeout(31 * time.Minute),
			Delete: schema.DefaultTimeout(31 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"cluster_mode": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"cluster_size_resource": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"db_cluster_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"engine": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"engine_params": {
				Type:     schema.TypeMap,
				Optional: true,
			},
			"group_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					if old != "" && new != "" && old != new {
						return strings.ToUpper(old) == strings.ToUpper(new)
					}
					return false
				},
			},
			"group_type": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: StringInSlice([]string{"interactive", "batch", "job", "default_type"}, false),
			},
			"max_cluster_count": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"max_compute_resource": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"min_cluster_count": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"min_compute_resource": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"node_num": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"update_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"user": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"users": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"connection_string": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"port": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAliCloudAdbResourceGroupCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "CreateDBResourceGroup"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["GroupName"] = d.Get("group_name")
	request["DBClusterId"] = d.Get("db_cluster_id")
	request["ClientToken"] = buildClientToken(action)

	if v, ok := d.GetOkExists("max_cluster_count"); ok {
		request["MaxClusterCount"] = v
	}
	if v, ok := d.GetOk("max_compute_resource"); ok {
		request["MaxComputeResource"] = v
	}
	if v, ok := d.GetOkExists("min_cluster_count"); ok {
		request["MinClusterCount"] = v
	}
	if v, ok := d.GetOk("engine_params"); ok {
		engineParamsJson, err := convertMaptoJsonString(v.(map[string]interface{}))
		if err != nil {
			return WrapError(err)
		}

		request["EngineParams"] = engineParamsJson
	}
	if v, ok := d.GetOk("group_type"); ok {
		request["GroupType"] = v
	}
	if v, ok := d.GetOk("cluster_mode"); ok {
		request["ClusterMode"] = v
	}
	if v, ok := d.GetOk("engine"); ok {
		request["Engine"] = v
	}
	if v, ok := d.GetOk("cluster_size_resource"); ok {
		request["ClusterSizeResource"] = v
	}
	if v, ok := d.GetOk("min_compute_resource"); ok {
		request["MinComputeResource"] = v
	}
	if v, ok := d.GetOkExists("node_num"); ok {
		request["NodeNum"] = v
	}
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RpcPost("adb", "2019-03-15", action, query, request, true)
		if err != nil {
			if IsExpectedErrors(err, []string{"ResourceNotEnough", "ACS.ServerError"}) || NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_adb_resource_group", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprintf("%v:%v", request["DBClusterId"], request["GroupName"]))

	adbServiceV2 := AdbServiceV2{client}
	stateConf := BuildStateConf([]string{}, []string{"Running"}, d.Timeout(schema.TimeoutCreate), 60*time.Second, adbServiceV2.AdbResourceGroupStateRefreshFunc(d.Id(), "Status", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAliCloudAdbResourceGroupUpdate(d, meta)
}

func resourceAliCloudAdbResourceGroupRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	adbServiceV2 := AdbServiceV2{client}

	objectRaw, err := adbServiceV2.DescribeAdbResourceGroup(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_adb_resource_group DescribeAdbResourceGroup Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("cluster_mode", objectRaw["ClusterMode"])
	d.Set("cluster_size_resource", objectRaw["ClusterSizeResource"])
	d.Set("create_time", objectRaw["CreateTime"])
	d.Set("engine", objectRaw["Engine"])
	d.Set("engine_params", objectRaw["EngineParams"])
	d.Set("group_type", objectRaw["GroupType"])
	d.Set("max_cluster_count", objectRaw["MaxClusterCount"])
	d.Set("max_compute_resource", objectRaw["MaxComputeResource"])
	d.Set("min_cluster_count", objectRaw["MinClusterCount"])
	d.Set("min_compute_resource", objectRaw["MinComputeResource"])
	d.Set("node_num", objectRaw["NodeNum"])
	d.Set("status", objectRaw["Status"])
	d.Set("update_time", objectRaw["UpdateTime"])
	d.Set("user", objectRaw["GroupUsers"])
	d.Set("group_name", objectRaw["GroupName"])
	d.Set("connection_string", objectRaw["ConnectionString"])
	d.Set("port", objectRaw["Port"])

	groupUserListRaw := make([]interface{}, 0)
	if objectRaw["GroupUserList"] != nil {
		groupUserListRaw = objectRaw["GroupUserList"].([]interface{})
	}

	d.Set("users", groupUserListRaw)

	parts := strings.Split(d.Id(), ":")
	d.Set("db_cluster_id", parts[0])

	return nil
}

func resourceAliCloudAdbResourceGroupUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	update := false

	var err error
	parts := strings.Split(d.Id(), ":")
	action := "ModifyDBResourceGroup"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["GroupName"] = parts[1]
	request["DBClusterId"] = parts[0]

	request["ClientToken"] = buildClientToken(action)
	if !d.IsNewResource() && d.HasChange("max_cluster_count") {
		update = true

		if v, ok := d.GetOkExists("max_cluster_count"); ok {
			request["MaxClusterCount"] = v
		}
	}

	if !d.IsNewResource() && d.HasChange("max_compute_resource") {
		update = true

		if v, ok := d.GetOk("max_compute_resource"); ok {
			request["MaxComputeResource"] = v
		}
	}

	if !d.IsNewResource() && d.HasChange("min_cluster_count") {
		update = true

		if v, ok := d.GetOkExists("min_cluster_count"); ok {
			request["MinClusterCount"] = v
		}
	}

	if !d.IsNewResource() && d.HasChange("engine_params") {
		update = true

		if v, ok := d.GetOk("engine_params"); ok {
			engineParamsJson, err := convertMaptoJsonString(v.(map[string]interface{}))
			if err != nil {
				return WrapError(err)
			}

			request["EngineParams"] = engineParamsJson
		}
	}

	if !d.IsNewResource() && d.HasChange("group_type") {
		update = true

		if v, ok := d.GetOk("group_type"); ok {
			request["GroupType"] = v
		}
	}

	if !d.IsNewResource() && d.HasChange("cluster_mode") {
		update = true

		if v, ok := d.GetOk("cluster_mode"); ok {
			request["ClusterMode"] = v
		}
	}

	if !d.IsNewResource() && d.HasChange("cluster_size_resource") {
		update = true

		if v, ok := d.GetOk("cluster_size_resource"); ok {
			request["ClusterSizeResource"] = v
		}
	}

	if !d.IsNewResource() && d.HasChange("min_compute_resource") {
		update = true

		if v, ok := d.GetOk("min_compute_resource"); ok {
			request["MinComputeResource"] = v
		}
	}

	if !d.IsNewResource() && d.HasChange("node_num") {
		update = true

		if v, ok := d.GetOkExists("node_num"); ok {
			request["NodeNum"] = v
		}
	}

	if d.HasChange("users") {
		update = true

		if v, ok := d.GetOk("users"); ok {
			request["PoolUserList"] = convertListToJsonString(v.(*schema.Set).List())
		}
	}

	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("adb", "2019-03-15", action, query, request, true)
			if err != nil {
				if IsExpectedErrors(err, []string{"ResourceNotEnough"}) || NeedRetry(err) {
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
		adbServiceV2 := AdbServiceV2{client}
		stateConf := BuildStateConf([]string{}, []string{"Running"}, d.Timeout(schema.TimeoutUpdate), 30*time.Second, adbServiceV2.AdbResourceGroupStateRefreshFunc(d.Id(), "Status", []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
	}

	return resourceAliCloudAdbResourceGroupRead(d, meta)
}

func resourceAliCloudAdbResourceGroupDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	parts := strings.Split(d.Id(), ":")
	action := "DeleteDBResourceGroup"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["GroupName"] = parts[1]
	request["DBClusterId"] = parts[0]

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.RpcPost("adb", "2019-03-15", action, query, request, true)

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
		if IsExpectedErrors(err, []string{"InvalidDBCluster.NotFound", "InvalidDBClusterId.NotFound"}) || NotFoundError(err) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	adbServiceV2 := AdbServiceV2{client}
	stateConf := BuildStateConf([]string{}, []string{}, d.Timeout(schema.TimeoutDelete), 30*time.Second, adbServiceV2.AdbResourceGroupStateRefreshFunc(d.Id(), "Status", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return nil
}
