package alicloud

import (
	"fmt"
	"log"
	"regexp"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAlicloudEcsSnapshotGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudEcsSnapshotGroupCreate,
		Read:   resourceAlicloudEcsSnapshotGroupRead,
		Update: resourceAlicloudEcsSnapshotGroupUpdate,
		Delete: resourceAlicloudEcsSnapshotGroupDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(1 * time.Minute),
			Update: schema.DefaultTimeout(1 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.All(validation.StringLenBetween(2, 256), validation.StringDoesNotMatch(regexp.MustCompile(`(^http://.*)|(^https://.*)`), "It must be `2` to `256` characters in length and cannot start with `https://` or `https://`.")),
			},
			"disk_id": {
				Type:          schema.TypeList,
				Optional:      true,
				MaxItems:      16,
				Elem:          &schema.Schema{Type: schema.TypeString},
				ConflictsWith: []string{"exclude_disk_id"},
			},
			"exclude_disk_id": {
				Type:          schema.TypeList,
				Optional:      true,
				MaxItems:      16,
				Elem:          &schema.Schema{Type: schema.TypeString},
				ConflictsWith: []string{"disk_id"},
			},
			"instance_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"instant_access": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"instant_access_retention_days": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntBetween(1, 65535),
			},
			"snapshot_group_name": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.All(validation.StringDoesNotMatch(regexp.MustCompile(`(^http://.*)|(^https://.*)`), "It cannot begin with \"http://\", \"https://\"."), validation.StringMatch(regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9_.:-]{1,127}$`), "It must be `2` to `128` characters in length, The description must start with a letter, and can contain letters, digits, underscores (_), and hyphens (-).")),
			},
			"resource_group_id": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tags": tagsSchema(),
		},
	}
}

func resourceAlicloudEcsSnapshotGroupCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "CreateSnapshotGroup"
	request := make(map[string]interface{})
	conn, err := client.NewEcsClient()
	if err != nil {
		return WrapError(err)
	}
	if v, ok := d.GetOk("description"); ok {
		request["Description"] = v
	}
	if v, ok := d.GetOk("disk_id"); ok {
		request["DiskId"] = v
	}
	if v, ok := d.GetOk("exclude_disk_id"); ok {
		request["ExcludeDiskId"] = v
	}
	if v, ok := d.GetOk("instance_id"); ok {
		request["InstanceId"] = v
	}
	if v, ok := d.GetOkExists("instant_access"); ok {
		request["InstantAccess"] = v
	}
	if v, ok := d.GetOkExists("resource_group_id"); ok {
		request["ResourceGroupId"] = v
	}
	if v, ok := d.GetOk("instant_access_retention_days"); ok {
		request["InstantAccessRetentionDays"] = v
	}
	request["RegionId"] = client.RegionId
	if v, ok := d.GetOk("snapshot_group_name"); ok {
		request["Name"] = v
	}
	if v, ok := d.GetOk("tags"); ok {
		count := 1
		for key, value := range v.(map[string]interface{}) {
			request[fmt.Sprintf("Tag.%d.Key", count)] = key
			request[fmt.Sprintf("Tag.%d.Value", count)] = value
			count++
		}
	}
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2014-05-26"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_ecs_snapshot_group", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["SnapshotGroupId"]))
	ecsService := EcsService{client}
	stateConf := BuildStateConf([]string{}, []string{"accomplished"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, ecsService.EcsSnapshotGroupStateRefreshFunc(d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAlicloudEcsSnapshotGroupRead(d, meta)
}
func resourceAlicloudEcsSnapshotGroupRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ecsService := EcsService{client}
	object, err := ecsService.DescribeEcsSnapshotGroup(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_ecs_snapshot_group ecsService.DescribeEcsSnapshotGroup Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	d.Set("description", object["Description"])
	d.Set("instance_id", object["InstanceId"])
	d.Set("snapshot_group_name", object["Name"])
	d.Set("resource_group_id", object["ResourceGroupId"])
	d.Set("status", object["Status"])
	if v, ok := object["Tags"].(map[string]interface{}); ok {
		d.Set("tags", tagsToMap(v["Tag"]))
	}

	return nil
}
func resourceAlicloudEcsSnapshotGroupUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	conn, err := client.NewEcsClient()
	if err != nil {
		return WrapError(err)
	}
	ecsService := EcsService{client}
	var response map[string]interface{}
	d.Partial(true)
	update := false
	request := map[string]interface{}{
		"SnapshotGroupId": d.Id(),
	}
	request["RegionId"] = client.RegionId
	if d.HasChange("description") {
		update = true
		if v, ok := d.GetOk("description"); ok {
			request["Description"] = v
		}
	}
	if d.HasChange("snapshot_group_name") {
		update = true
		if v, ok := d.GetOk("snapshot_group_name"); ok {
			request["Name"] = v
		}
	}
	if update {
		action := "ModifySnapshotGroup"
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2014-05-26"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
		d.SetPartial("snapshot_group_name")
	}
	if d.HasChange("tags") {
		if err := ecsService.SetResourceTags(d, "snapshotgroup"); err != nil {
			return WrapError(err)
		}
		d.SetPartial("tags")
	}
	d.Partial(false)
	return resourceAlicloudEcsSnapshotGroupRead(d, meta)
}
func resourceAlicloudEcsSnapshotGroupDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "DeleteSnapshotGroup"
	var response map[string]interface{}
	conn, err := client.NewEcsClient()
	if err != nil {
		return WrapError(err)
	}
	request := map[string]interface{}{
		"SnapshotGroupId": d.Id(),
		"RegionId":        client.RegionId,
	}

	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2014-05-26"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
