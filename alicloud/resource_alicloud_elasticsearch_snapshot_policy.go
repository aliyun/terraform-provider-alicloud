package alicloud

import (
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAliCloudElasticsearchSnapshotPolicy() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudElasticsearchSnapshotPolicyCreate,
		Read:   resourceAliCloudElasticsearchSnapshotPolicyRead,
		Update: resourceAliCloudElasticsearchSnapshotPolicyUpdate,
		Delete: resourceAliCloudElasticsearchSnapshotPolicyDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"enable": {
				Type:     schema.TypeBool,
				Required: true,
			},
			"quartz_regex": {
				Type:     schema.TypeString,
				Optional: true,
				// quartz_regex only takes effect when enable is true. Keeping it in the
				// configuration while enable is false would otherwise produce a perpetual diff.
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return !d.Get("enable").(bool)
				},
			},
		},
	}
}

func resourceAliCloudElasticsearchSnapshotPolicyCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	enable := d.Get("enable").(bool)
	quartzRegex := d.Get("quartz_regex").(string)
	if enable && quartzRegex == "" {
		return WrapError(Error("'quartz_regex' is required when 'enable' is set to true"))
	}

	// The ID is set only after the initial update succeeds. Otherwise a failed
	// create (e.g. non-existent instance) leaves a tainted resource in state and
	// every following refresh keeps failing on the instance-scoped API.
	if err := updateElasticsearchSnapshotSetting(client, d.Get("instance_id").(string), enable, quartzRegex, d.Timeout(schema.TimeoutCreate)); err != nil {
		return err
	}

	d.SetId(d.Get("instance_id").(string))

	return resourceAliCloudElasticsearchSnapshotPolicyRead(d, meta)
}

func resourceAliCloudElasticsearchSnapshotPolicyRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	elasticsearchService := ElasticsearchService{client}

	object, err := elasticsearchService.DescribeSnapshotSetting(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("instance_id", d.Id())
	d.Set("enable", snapshotSettingAttribute(object, "Enable", "enable"))
	d.Set("quartz_regex", snapshotSettingAttribute(object, "QuartzRegex", "quartzRegex"))

	return nil
}

func resourceAliCloudElasticsearchSnapshotPolicyUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	if d.HasChange("enable") || d.HasChange("quartz_regex") {
		enable := d.Get("enable").(bool)
		quartzRegex := d.Get("quartz_regex").(string)
		if enable && quartzRegex == "" {
			return WrapError(Error("'quartz_regex' is required when 'enable' is set to true"))
		}
		if err := updateElasticsearchSnapshotSetting(client, d.Id(), enable, quartzRegex, d.Timeout(schema.TimeoutUpdate)); err != nil {
			return err
		}
	}

	return resourceAliCloudElasticsearchSnapshotPolicyRead(d, meta)
}

func resourceAliCloudElasticsearchSnapshotPolicyDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	// There is no dedicated delete API for the snapshot setting. The delete action
	// only disables the automatic snapshot creation and keeps the instance and the
	// existing snapshot data untouched.
	err := updateElasticsearchSnapshotSetting(client, d.Id(), false, "", d.Timeout(schema.TimeoutDelete))
	if err != nil {
		if NotFoundError(err) {
			return nil
		}
		return err
	}

	return nil
}
