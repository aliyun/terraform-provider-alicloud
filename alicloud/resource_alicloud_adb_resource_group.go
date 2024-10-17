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
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(1 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"db_cluster_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
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
				ValidateFunc: StringInSlice([]string{"batch", "default_type", "interactive"}, false),
			},
			"node_num": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"users": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"user": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"update_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAliCloudAdbResourceGroupCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	var err error
	action := "CreateDBResourceGroup"
	request := make(map[string]interface{})
	request["DBClusterId"] = d.Get("db_cluster_id")
	request["GroupName"] = d.Get("group_name")

	if v, ok := d.GetOk("group_type"); ok {
		request["GroupType"] = v
	}

	if v, ok := d.GetOkExists("node_num"); ok {
		request["NodeNum"] = v
	}

	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutCreate)), func() *resource.RetryError {
		response, err = client.RpcPost("adb", "2019-03-15", action, nil, request, false)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_adb_resource_group", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprintf("%v:%v", request["DBClusterId"], request["GroupName"]))

	return resourceAliCloudAdbResourceGroupUpdate(d, meta)
}

func resourceAliCloudAdbResourceGroupRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	adbService := AdbService{client}

	object, err := adbService.DescribeAdbResourceGroup(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_adb_resource_group adbService.DescribeAdbResourceGroup Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}

	d.Set("db_cluster_id", parts[0])
	d.Set("group_name", object["GroupName"])
	d.Set("group_type", object["GroupType"])
	d.Set("node_num", object["NodeNum"])
	d.Set("users", object["GroupUserList"])
	d.Set("user", object["GroupUsers"])
	d.Set("create_time", object["CreateTime"])
	d.Set("update_time", object["UpdateTime"])

	return nil
}

func resourceAliCloudAdbResourceGroupUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	update := false

	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}

	request := map[string]interface{}{
		"DBClusterId": parts[0],
		"GroupName":   parts[1],
	}

	if !d.IsNewResource() && d.HasChange("group_type") {
		update = true
	}
	if v, ok := d.GetOk("group_type"); ok {
		request["GroupType"] = v
	}

	if !d.IsNewResource() && d.HasChange("node_num") {
		update = true
	}
	if v, ok := d.GetOkExists("node_num"); ok {
		request["NodeNum"] = v
	}

	if d.HasChange("users") {
		update = true

		if v, ok := d.GetOk("users"); ok {
			request["PoolUserList"] = convertListToJsonString(v.(*schema.Set).List())
		}
	}

	if update {
		action := "ModifyDBResourceGroup"
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *resource.RetryError {
			response, err = client.RpcPost("adb", "2019-03-15", action, nil, request, false)
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
	}

	return resourceAliCloudAdbResourceGroupRead(d, meta)
}

func resourceAliCloudAdbResourceGroupDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "DeleteDBResourceGroup"
	var response map[string]interface{}
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}

	request := map[string]interface{}{
		"DBClusterId": parts[0],
		"GroupName":   parts[1],
	}

	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutDelete)), func() *resource.RetryError {
		response, err = client.RpcPost("adb", "2019-03-15", action, nil, request, false)
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
		if IsExpectedErrors(err, []string{"InvalidDBCluster.NotFound"}) || NotFoundError(err) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	return nil
}
