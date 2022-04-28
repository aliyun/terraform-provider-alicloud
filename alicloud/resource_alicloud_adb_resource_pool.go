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

func resourceAlicloudAdbResourcePool() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudAdbResourcePoolCreate,
		Read:   resourceAlicloudAdbResourcePoolRead,
		Update: resourceAlicloudAdbResourcePoolUpdate,
		Delete: resourceAlicloudAdbResourcePoolDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"db_cluster_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"node_num": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"pool_name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringMatch(regexp.MustCompile(`^[[A-Z0-9._-]{1,64}$`), "The name must be `1` to `64` characters in length, and can contain uppercase letters, digits, hyphens (-) and underscores (_)."),
				ForceNew:     true,
			},
			"query_type": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice([]string{"default_type", "batch", "interactive"}, false),
			},
		},
	}
}

func resourceAlicloudAdbResourcePoolCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "CreateDBResourcePool"
	request := make(map[string]interface{})
	conn, err := client.NewAdsClient()
	if err != nil {
		return WrapError(err)
	}
	request["DBClusterId"] = d.Get("db_cluster_id")
	if v, ok := d.GetOk("node_num"); ok {
		request["NodeNum"] = v
	}
	request["PoolName"] = d.Get("pool_name")
	if v, ok := d.GetOk("query_type"); ok {
		request["QueryType"] = v
	}
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-03-15"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_adb_resource_pool", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(request["DBClusterId"], ":", request["PoolName"]))

	return resourceAlicloudAdbResourcePoolRead(d, meta)
}
func resourceAlicloudAdbResourcePoolRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	adbService := AdbService{client}
	object, err := adbService.DescribeAdbResourcePool(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_adb_resource_pool adbService.DescribeAdbResourcePool Failed!!! %s", err)
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
	d.Set("pool_name", parts[1])
	if v, ok := object["NodeNum"]; ok {
		d.Set("node_num", formatInt(v))
	}
	d.Set("query_type", object["QueryType"])
	return nil
}
func resourceAlicloudAdbResourcePoolUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	update := false
	request := map[string]interface{}{
		"DBClusterId": parts[0],
		"PoolName":    parts[1],
	}
	if d.HasChange("node_num") {
		update = true
		if v, ok := d.GetOk("node_num"); ok {
			request["NodeNum"] = v
		}
	}
	if d.HasChange("query_type") {
		update = true
		if v, ok := d.GetOk("query_type"); ok {
			request["QueryType"] = v
		}
	}
	if update {
		action := "ModifyDBResourcePool"
		conn, err := client.NewAdsClient()
		if err != nil {
			return WrapError(err)
		}
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-03-15"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
	return resourceAlicloudAdbResourcePoolRead(d, meta)
}
func resourceAlicloudAdbResourcePoolDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	action := "DeleteDBResourcePool"
	var response map[string]interface{}
	conn, err := client.NewAdsClient()
	if err != nil {
		return WrapError(err)
	}
	request := map[string]interface{}{
		"DBClusterId": parts[0],
		"PoolName":    parts[1],
	}

	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-03-15"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
		if IsExpectedErrors(err, []string{"PoolNotExist"}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}
	return nil
}
