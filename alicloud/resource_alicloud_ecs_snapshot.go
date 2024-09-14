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

func resourceAliCloudEcsSnapshot() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudEcsSnapshotCreate,
		Read:   resourceAliCloudEcsSnapshotRead,
		Update: resourceAliCloudEcsSnapshotUpdate,
		Delete: resourceAliCloudEcsSnapshotDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(2 * time.Minute),
			Update: schema.DefaultTimeout(2 * time.Minute),
			Delete: schema.DefaultTimeout(2 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"disk_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"category": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Computed:     true,
				ValidateFunc: StringInSlice([]string{"standard", "flash"}, false),
			},
			"retention_days": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"resource_group_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"snapshot_name": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"name"},
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"force": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"tags": tagsSchema(),
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"snapshot_name"},
				Deprecated:    "Field `name` has been deprecated from provider version 1.120.0. New field `snapshot_name` instead.",
			},
			"instant_access": {
				Type:       schema.TypeBool,
				Optional:   true,
				Deprecated: "Field `instant_access` has been deprecated from provider version 1.231.0.",
			},
			"instant_access_retention_days": {
				Type:       schema.TypeInt,
				Optional:   true,
				Deprecated: "Field `instant_access_retention_days` has been deprecated from provider version 1.231.0.",
			},
		},
	}
}

func resourceAliCloudEcsSnapshotCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ecsService := EcsService{client}
	var response map[string]interface{}
	action := "CreateSnapshot"
	request := make(map[string]interface{})
	conn, err := client.NewEcsClient()
	if err != nil {
		return WrapError(err)
	}

	request["ClientToken"] = buildClientToken("CreateSnapshot")
	request["DiskId"] = d.Get("disk_id")

	if v, ok := d.GetOk("category"); ok {
		request["Category"] = v
	}

	if v, ok := d.GetOkExists("retention_days"); ok {
		request["RetentionDays"] = v
	}

	if v, ok := d.GetOk("resource_group_id"); ok {
		request["ResourceGroupId"] = v
	}

	if v, ok := d.GetOk("snapshot_name"); ok {
		request["SnapshotName"] = v
	} else if v, ok := d.GetOk("name"); ok {
		request["SnapshotName"] = v
	}

	if v, ok := d.GetOk("description"); ok {
		request["Description"] = v
	}

	if v, ok := d.GetOk("tags"); ok {
		count := 1
		for key, value := range v.(map[string]interface{}) {
			request[fmt.Sprintf("Tag.%d.Key", count)] = key
			request[fmt.Sprintf("Tag.%d.Value", count)] = value
			count++
		}
	}

	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutCreate)), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2014-05-26"), StringPointer("AK"), nil, request, &runtime)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_ecs_snapshot", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["SnapshotId"]))

	stateConf := BuildStateConf([]string{}, []string{"accomplished"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, ecsService.EcsSnapshotStateRefreshFunc(d.Id(), []string{"failed"}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAliCloudEcsSnapshotRead(d, meta)
}

func resourceAliCloudEcsSnapshotRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ecsService := EcsService{client}

	object, err := ecsService.DescribeEcsSnapshot(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_ecs_snapshot ecsService.DescribeEcsSnapshot Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("disk_id", object["SourceDiskId"])
	d.Set("category", object["Category"])
	d.Set("retention_days", formatInt(object["RetentionDays"]))
	d.Set("resource_group_id", object["ResourceGroupId"])
	d.Set("snapshot_name", object["SnapshotName"])
	d.Set("description", object["Description"])
	d.Set("status", object["Status"])
	d.Set("name", object["SnapshotName"])

	tags, err := ecsService.ListTagResources(d.Id(), "snapshot")
	if err != nil {
		return WrapError(err)
	}

	d.Set("tags", tagsToMap(tags))

	return nil
}

func resourceAliCloudEcsSnapshotUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ecsService := EcsService{client}
	var response map[string]interface{}
	d.Partial(true)

	update := false
	request := map[string]interface{}{
		"SnapshotId": d.Id(),
	}

	if d.HasChange("retention_days") {
		update = true

		if v, ok := d.GetOkExists("retention_days"); ok {
			request["RetentionDays"] = v
		}
	}

	if d.HasChange("snapshot_name") {
		update = true

		if v, ok := d.GetOk("snapshot_name"); ok {
			request["SnapshotName"] = v
		}
	}

	if d.HasChange("description") {
		update = true
	}
	if v, ok := d.GetOk("description"); ok {
		request["Description"] = v
	}

	if d.HasChange("name") {
		update = true

		if v, ok := d.GetOk("name"); ok {
			request["SnapshotName"] = v
		}
	}

	if update {
		action := "ModifySnapshotAttribute"
		conn, err := client.NewEcsClient()
		if err != nil {
			return WrapError(err)
		}

		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *resource.RetryError {
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

		stateConf := BuildStateConf([]string{}, []string{"accomplished"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, ecsService.EcsSnapshotStateRefreshFunc(d.Id(), []string{"failed"}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}

		d.SetPartial("retention_days")
		d.SetPartial("snapshot_name")
		d.SetPartial("description")
		d.SetPartial("name")
	}

	if d.HasChange("tags") {
		if err := ecsService.SetResourceTags(d, "snapshot"); err != nil {
			return WrapError(err)
		}

		d.SetPartial("tags")
	}

	d.Partial(false)

	return resourceAliCloudEcsSnapshotRead(d, meta)
}

func resourceAliCloudEcsSnapshotDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ecsService := EcsService{client}
	action := "DeleteSnapshot"
	var response map[string]interface{}

	conn, err := client.NewEcsClient()
	if err != nil {
		return WrapError(err)
	}

	request := map[string]interface{}{
		"SnapshotId": d.Id(),
	}

	if v, ok := d.GetOkExists("force"); ok {
		request["Force"] = v
	}

	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutDelete)), func() *resource.RetryError {
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
		if IsExpectedErrors(err, []string{"InvalidSnapshotId.NotFound"}) || NotFoundError(err) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	stateConf := BuildStateConf([]string{}, []string{}, d.Timeout(schema.TimeoutDelete), 5*time.Second, ecsService.EcsSnapshotStateRefreshFunc(d.Id(), []string{"failed"}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return nil
}
