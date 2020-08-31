package alicloud

import (
	"fmt"
	"log"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/actiontrail"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceAlicloudActiontrailTrail() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudActiontrailTrailCreate,
		Read:   resourceAlicloudActiontrailTrailRead,
		Update: resourceAlicloudActiontrailTrailUpdate,
		Delete: resourceAlicloudActiontrailTrailDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(1 * time.Minute),
			Update: schema.DefaultTimeout(1 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"event_rw": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"All", "Read", "Write"}, false),
				Default:      "Write",
			},
			"is_organization_trail": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
			},
			"mns_topic_arn": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"oss_bucket_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"oss_key_prefix": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"role_name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"sls_project_arn": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"sls_write_role_arn": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"status": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"Enable", "Disable"}, false),
				Default:      "Enable",
			},
			"trail_name": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ForceNew:      true,
				ConflictsWith: []string{"name"},
			},
			"name": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ForceNew:      true,
				Deprecated:    "Field 'name' has been deprecated from version 1.95.0. Use 'trail_name' instead.",
				ConflictsWith: []string{"trail_name"},
			},
			"trail_region": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"All", "cn-beijing", "cn-hangzhou"}, false),
				Default:      "All",
			},
		},
	}
}

func resourceAlicloudActiontrailTrailCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	actiontrailService := ActiontrailService{client}

	request := actiontrail.CreateCreateTrailRequest()
	if v, ok := d.GetOk("event_rw"); ok {
		request.EventRW = v.(string)
	}

	if v, ok := d.GetOkExists("is_organization_trail"); ok {
		request.IsOrganizationTrail = requests.NewBoolean(v.(bool))
	}

	if v, ok := d.GetOk("mns_topic_arn"); ok {
		request.MnsTopicArn = v.(string)
	}

	if v, ok := d.GetOk("oss_bucket_name"); ok {
		request.OssBucketName = v.(string)
	}

	if v, ok := d.GetOk("oss_key_prefix"); ok {
		request.OssKeyPrefix = v.(string)
	}

	if v, ok := d.GetOk("role_name"); ok {
		request.RoleName = v.(string)
	}

	if v, ok := d.GetOk("sls_project_arn"); ok {
		request.SlsProjectArn = v.(string)
	}

	if v, ok := d.GetOk("sls_write_role_arn"); ok {
		request.SlsWriteRoleArn = v.(string)
	}

	if v, ok := d.GetOk("trail_name"); ok {
		request.Name = v.(string)
	} else if v, ok := d.GetOk("name"); ok {
		request.Name = v.(string)
	} else {
		return WrapError(Error(`[ERROR] Argument "name" or "trail_name" must be set one!`))
	}

	if v, ok := d.GetOk("trail_region"); ok {
		request.TrailRegion = v.(string)
	}

	wait := incrementalWait(3*time.Second, 10*time.Second)
	err := resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		raw, err := client.WithActiontrailClient(func(actiontrailClient *actiontrail.Client) (interface{}, error) {
			return actiontrailClient.CreateTrail(request)
		})
		if err != nil {
			if IsExpectedErrors(err, []string{"InsufficientBucketPolicyException"}) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(request.GetActionName(), raw)
		response, _ := raw.(*actiontrail.CreateTrailResponse)
		d.SetId(fmt.Sprintf("%v", response.Name))
		return nil
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_actiontrail_trail", request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	stateConf := BuildStateConf([]string{}, []string{"Fresh"}, d.Timeout(schema.TimeoutCreate), 10*time.Second, actiontrailService.ActiontrailTrailStateRefreshFunc(d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAlicloudActiontrailTrailUpdate(d, meta)
}
func resourceAlicloudActiontrailTrailRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	actiontrailService := ActiontrailService{client}
	object, err := actiontrailService.DescribeActiontrailTrail(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_actiontrail_trail actiontrailService.DescribeActiontrailTrail Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("trail_name", d.Id())
	d.Set("name", d.Id())
	d.Set("event_rw", object.EventRW)
	d.Set("oss_bucket_name", object.OssBucketName)
	d.Set("oss_key_prefix", object.OssKeyPrefix)
	d.Set("role_name", object.RoleName)
	d.Set("sls_project_arn", object.SlsProjectArn)
	d.Set("sls_write_role_arn", object.SlsWriteRoleArn)
	d.Set("status", object.Status)
	d.Set("trail_region", object.TrailRegion)
	return nil
}
func resourceAlicloudActiontrailTrailUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	actiontrailService := ActiontrailService{client}
	d.Partial(true)

	update := false
	request := actiontrail.CreateUpdateTrailRequest()
	request.Name = d.Id()
	if !d.IsNewResource() && d.HasChange("event_rw") {
		update = true
		request.EventRW = d.Get("event_rw").(string)
	}
	if !d.IsNewResource() && d.HasChange("mns_topic_arn") {
		update = true
		request.MnsTopicArn = d.Get("mns_topic_arn").(string)
	}
	if !d.IsNewResource() && d.HasChange("oss_bucket_name") {
		update = true
		request.OssBucketName = d.Get("oss_bucket_name").(string)
	}
	if !d.IsNewResource() && d.HasChange("oss_key_prefix") {
		update = true
		request.OssKeyPrefix = d.Get("oss_key_prefix").(string)
	}
	request.RegionId = client.RegionId
	if !d.IsNewResource() && d.HasChange("role_name") {
		update = true
		request.RoleName = d.Get("role_name").(string)
	}
	if !d.IsNewResource() && d.HasChange("sls_project_arn") {
		update = true
		request.SlsProjectArn = d.Get("sls_project_arn").(string)
	}
	if !d.IsNewResource() && d.HasChange("sls_write_role_arn") {
		update = true
		request.SlsWriteRoleArn = d.Get("sls_write_role_arn").(string)
	}
	if !d.IsNewResource() && d.HasChange("trail_region") {
		update = true
		request.TrailRegion = d.Get("trail_region").(string)
	}
	if update {

		if v, ok := d.GetOk("sls_project_arn"); ok {
			request.SlsProjectArn = v.(string)
		}
		if v, ok := d.GetOk("sls_write_role_arn"); ok {
			request.SlsWriteRoleArn = v.(string)
		}
		if v, ok := d.GetOk("oss_bucket_name"); ok {
			request.OssBucketName = v.(string)
		}

		wait := incrementalWait(3*time.Second, 5*time.Second)
		err := resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			raw, err := client.WithActiontrailClient(func(actiontrailClient *actiontrail.Client) (interface{}, error) {
				return actiontrailClient.UpdateTrail(request)
			})
			if err != nil {
				if IsExpectedErrors(err, []string{"InsufficientBucketPolicyException"}) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			addDebug(request.GetActionName(), raw)
			return nil
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		d.SetPartial("event_rw")
		d.SetPartial("mns_topic_arn")
		d.SetPartial("oss_bucket_name")
		d.SetPartial("oss_key_prefix")
		d.SetPartial("role_name")
		d.SetPartial("sls_project_arn")
		d.SetPartial("sls_write_role_arn")
		d.SetPartial("trail_region")
	}
	if d.HasChange("status") {
		object, err := actiontrailService.DescribeActiontrailTrail(d.Id())
		if err != nil {
			return WrapError(err)
		}
		target := d.Get("status").(string)
		if object.Status != target {
			if target == "Disable" {
				request := actiontrail.CreateStopLoggingRequest()
				request.Name = d.Id()
				wait := incrementalWait(3*time.Second, 5*time.Second)
				err := resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
					raw, err := client.WithActiontrailClient(func(actiontrailClient *actiontrail.Client) (interface{}, error) {
						return actiontrailClient.StopLogging(request)
					})
					if err != nil {
						if IsExpectedErrors(err, []string{"InsufficientBucketPolicyException"}) {
							wait()
							return resource.RetryableError(err)
						}
						return resource.NonRetryableError(err)
					}
					addDebug(request.GetActionName(), raw)
					return nil
				})
				if err != nil {
					return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
				}
				stateConf := BuildStateConf([]string{}, []string{"Disable"}, d.Timeout(schema.TimeoutUpdate), 10*time.Second, actiontrailService.ActiontrailTrailStateRefreshFunc(d.Id(), []string{}))
				if _, err := stateConf.WaitForState(); err != nil {
					return WrapErrorf(err, IdMsg, d.Id())
				}
			}
			if target == "Enable" {
				request := actiontrail.CreateStartLoggingRequest()
				request.Name = d.Id()
				wait := incrementalWait(3*time.Second, 5*time.Second)
				err := resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
					raw, err := client.WithActiontrailClient(func(actiontrailClient *actiontrail.Client) (interface{}, error) {
						return actiontrailClient.StartLogging(request)
					})
					if err != nil {
						if IsExpectedErrors(err, []string{"InsufficientBucketPolicyException"}) {
							wait()
							return resource.RetryableError(err)
						}
						return resource.NonRetryableError(err)
					}
					addDebug(request.GetActionName(), raw)
					return nil
				})
				if err != nil {
					return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
				}
				stateConf := BuildStateConf([]string{}, []string{"Enable"}, d.Timeout(schema.TimeoutUpdate), 10*time.Second, actiontrailService.ActiontrailTrailStateRefreshFunc(d.Id(), []string{}))
				if _, err := stateConf.WaitForState(); err != nil {
					return WrapErrorf(err, IdMsg, d.Id())
				}
			}
			d.SetPartial("status")
		}
	}
	d.Partial(false)
	return resourceAlicloudActiontrailTrailRead(d, meta)
}
func resourceAlicloudActiontrailTrailDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	request := actiontrail.CreateDeleteTrailRequest()
	request.Name = d.Id()
	raw, err := client.WithActiontrailClient(func(actiontrailClient *actiontrail.Client) (interface{}, error) {
		return actiontrailClient.DeleteTrail(request)
	})
	addDebug(request.GetActionName(), raw)
	if err != nil {
		if IsExpectedErrors(err, []string{"TrailNotFoundException"}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	return nil
}
