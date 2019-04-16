package alicloud

import (
	"fmt"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/actiontrail"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func resourceAlicloudActiontrail() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudActiontrailCreate,
		Read:   resourceAlicloudActiontrailRead,
		Update: resourceAlicloudActiontrailUpdate,
		Delete: resourceAlicloudActiontrailDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"event_rw": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"oss_bucket_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"role_name": {
				Type:             schema.TypeString,
				Required:         true,
				DiffSuppressFunc: actiontrailRoleNmaeDiffSuppressFunc,
			},
			"oss_key_prefix": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"sls_project_arn": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"sls_write_role_arn": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceAlicloudActiontrailCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	trailService := ActionTrailService{client}
	request := actiontrail.CreateCreateTrailRequest()

	request.Name = d.Get("name").(string)
	request.OssBucketName = d.Get("oss_bucket_name").(string)
	request.RoleName = d.Get("role_name").(string)

	if v, ok := d.GetOk("even_rw"); ok && v.(string) != "" {
		request.EventRW = v.(string)
	}
	if v, ok := d.GetOk("oss_bucket_name"); ok && v.(string) != "" {
		request.OssBucketName = v.(string)
	}
	if v, ok := d.GetOk("role_name"); ok && v.(string) != "" {
		request.RoleName = v.(string)
	}
	if v, ok := d.GetOk("oss_key_prefix"); ok && v.(string) != "" {
		request.OssKeyPrefix = v.(string)
	}
	if v, ok := d.GetOk("sls_project_arn"); ok && v.(string) != "" {
		request.SlsProjectArn = v.(string)
	}
	if v, ok := d.GetOk("sls_write_role_arn"); ok && v.(string) != "" {
		request.SlsWriteRoleArn = v.(string)
	}

	var raw interface{}
	var err error
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		raw, err = client.WithActionTrailClient(func(actiontrailClient *actiontrail.Client) (interface{}, error) {
			return actiontrailClient.CreateTrail(request)
		})
		if err != nil {
			if IsExceptedErrors(err, []string{InsufficientBucketPolicyException}) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(request.GetActionName(), raw)
		return nil
	})

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_actiontrail", request.GetActionName(), AlibabaCloudSdkGoERROR)
	}

	response, _ := raw.(*actiontrail.CreateTrailResponse)
	if response == nil {
		return WrapError(fmt.Errorf("CreateActionTrail got a nil response: %#v", response))
	}

	d.SetId(response.Name)

	if err := trailService.WaitForActionTrail(d.Id(), DefaultTimeout); err != nil {
		return WrapError(err)
	}

	return resourceAlicloudActiontrailRead(d, meta)
}

func resourceAlicloudActiontrailRead(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	trailService := ActionTrailService{client}
	object, err := trailService.DescribeActionTrail(d.Id())

	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("name", object.Name)
	d.Set("even_rw", object.EventRW)
	d.Set("role_name", object.RoleName)
	d.Set("oss_bucket_name", object.OssBucketName)
	d.Set("oss_key_prefix", object.OssKeyPrefix)
	d.Set("sls_project_arn", object.SlsProjectArn)
	d.Set("sls_write_role_arn", object.SlsWriteRoleArn)

	return nil
}

func resourceAlicloudActiontrailUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	update := false
	request := actiontrail.CreateUpdateTrailRequest()
	request.Name = d.Id()
	request.OssBucketName = d.Get("oss_bucket_name").(string)
	request.RoleName = d.Get("role_name").(string)
	request.SlsProjectArn = d.Get("sls_project_arn").(string)
	request.SlsWriteRoleArn = d.Get("sls_write_role_arn").(string)

	if d.HasChange("even_rw") {
		request.EventRW = d.Get("even_rw").(string)
		update = true
	}
	if d.HasChange("oss_bucket_name") {
		update = true
	}
	if d.HasChange("role_name") {
		update = true
	}
	if d.HasChange("oss_key_prefix") {
		request.OssKeyPrefix = d.Get("oss_key_prefix").(string)
		update = true
	}
	if d.HasChange("sls_project_arn") {
		update = true
	}
	if d.HasChange("sls_write_role_arn") {
		update = true
	}
	if update {

		err := resource.Retry(5*time.Minute, func() *resource.RetryError {
			raw, err := client.WithActionTrailClient(func(actiontrailClient *actiontrail.Client) (interface{}, error) {
				return actiontrailClient.UpdateTrail(request)
			})
			if err != nil {
				if IsExceptedErrors(err, []string{InsufficientBucketPolicyException}) {
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
	}

	return resourceAlicloudActiontrailRead(d, meta)
}

func resourceAlicloudActiontrailDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	trailService := ActionTrailService{client}
	request := actiontrail.CreateDeleteTrailRequest()
	request.Name = d.Id()
	return resource.Retry(5*time.Minute, func() *resource.RetryError {
		raw, err := client.WithActionTrailClient(func(actiontrailClient *actiontrail.Client) (interface{}, error) {
			return actiontrailClient.DeleteTrail(request)
		})

		if err != nil {
			if IsExceptedErrors(err, []string{InvalidVpcIDNotFound, ForbiddenVpcNotFound}) {
				return nil
			}
			return resource.RetryableError(WrapErrorf(err, DefaultTimeoutMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR))
		}

		addDebug(request.GetActionName(), raw)

		if _, err := trailService.DescribeActionTrail(d.Id()); err != nil {
			if NotFoundError(err) {
				return nil
			}
			return resource.NonRetryableError(WrapError(err))
		}

		return resource.RetryableError(WrapErrorf(err, DefaultTimeoutMsg, d.Id(), request.GetActionName(), ProviderERROR))
	})
}
