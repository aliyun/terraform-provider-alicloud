package alicloud

import (
	"fmt"
	"log"
	"regexp"
	"time"

	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceAlicloudAdbResourceGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudAdbResourceGroupCreate,
		Read:   resourceAlicloudAdbResourceGroupRead,
		Update: resourceAlicloudAdbResourceGroupUpdate,
		Delete: resourceAlicloudAdbResourceGroupDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"create_time": {
				Computed: true,
				Type:     schema.TypeString,
			},
			"update_time": {
				Computed: true,
				Type:     schema.TypeString,
			},
			"db_cluster_id": {
				Required: true,
				ForceNew: true,
				Type:     schema.TypeString,
			},
			"group_name": {
				Required:     true,
				ForceNew:     true,
				Type:         schema.TypeString,
				ValidateFunc: validation.StringMatch(regexp.MustCompile(`^[A-Z0-9_]{2,30}$`), "The group name must be 2 to 30 characters in length, and can contain upper case letters, digits, and underscore(_)."),
			},
			"group_type": {
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"batch", "default_type", "interactive"}, false),
				Type:         schema.TypeString,
			},
			"node_num": {
				Optional: true,
				Computed: true,
				Type:     schema.TypeInt,
			},
			"user": {
				Computed: true,
				Type:     schema.TypeString,
			},
		},
	}
}

func resourceAlicloudAdbResourceGroupCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	request := make(map[string]interface{})
	conn, err := client.NewAdsClient()
	if err != nil {
		return WrapError(err)
	}

	if v, ok := d.GetOk("db_cluster_id"); ok {
		request["DBClusterId"] = v
	}
	if v, ok := d.GetOk("group_name"); ok {
		request["GroupName"] = v
	}
	if v, ok := d.GetOk("group_type"); ok {
		request["GroupType"] = v
	}
	if v, ok := d.GetOk("node_num"); ok {
		request["NodeNum"] = v
	}

	var response map[string]interface{}
	action := "CreateDBResourceGroup"
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutCreate)), func() *resource.RetryError {
		resp, err := conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-03-15"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_adb_resource_group", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(request["DBClusterId"], ":", request["GroupName"]))

	return resourceAlicloudAdbResourceGroupRead(d, meta)
}

func resourceAlicloudAdbResourceGroupRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	adbService := AdbService{client}

	object, err := adbService.DescribeAdbResourceGroup(d.Id())
	if err != nil {
		if NotFoundError(err) {
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
	d.Set("group_name", parts[1])
	d.Set("create_time", object["CreateTime"])
	d.Set("update_time", object["UpdateTime"])
	d.Set("group_type", object["GroupType"])
	d.Set("node_num", object["NodeNum"])
	d.Set("user", object["GroupUsers"])

	return nil
}

func resourceAlicloudAdbResourceGroupUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	conn, err := client.NewAdsClient()
	if err != nil {
		return WrapError(err)
	}
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	update := false
	request := map[string]interface{}{
		"DBClusterId": parts[0],
		"GroupName":   parts[1],
	}

	if !d.IsNewResource() && d.HasChange("node_num") {
		update = true
		if v, ok := d.GetOk("node_num"); ok {
			request["NodeNum"] = v
		}
	}

	if update {
		action := "ModifyDBResourceGroup"
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *resource.RetryError {
			resp, err := conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-03-15"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			addDebug(action, resp, request)
			return nil
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
	}

	return resourceAlicloudAdbResourceGroupRead(d, meta)
}

func resourceAlicloudAdbResourceGroupDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	conn, err := client.NewAdsClient()
	if err != nil {
		return WrapError(err)
	}
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}

	request := map[string]interface{}{
		"DBClusterId": parts[0], "GroupName": parts[1],
	}

	action := "DeleteDBResourceGroup"
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutDelete)), func() *resource.RetryError {
		resp, err := conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-03-15"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, resp, request)
		return nil
	})
	if err != nil {
		if IsExpectedErrors(err, []string{"InvalidDBCluster.NotFound"}) || NotFoundError(err) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}
	return nil
}
