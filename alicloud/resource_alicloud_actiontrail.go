package alicloud

import (
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
				Type:         schema.TypeString,
				Optional:     true,
				Default:      EventWrite,
				ValidateFunc: validateActiontrailEventrw,
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

	if v, ok := d.GetOk("event_rw"); ok && v.(string) != "" {
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

	d.SetId(response.Name)
	if err := trailService.startActionTrail(d.Id()); err != nil {
		return WrapError(err)
	}
	if err := trailService.WaitForActionTrail(d.Id(), Enable, DefaultTimeout); err != nil {
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
	d.Set("event_rw", object.EventRW)
	d.Set("role_name", object.RoleName)
	d.Set("oss_bucket_name", object.OssBucketName)
	d.Set("oss_key_prefix", object.OssKeyPrefix)
	d.Set("sls_project_arn", object.SlsProjectArn)
	d.Set("sls_write_role_arn", object.SlsWriteRoleArn)

	return nil
}

func resourceAlicloudActiontrailUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	request := actiontrail.CreateUpdateTrailRequest()
	request.Name = d.Id()
	request.OssBucketName = d.Get("oss_bucket_name").(string)
	request.RoleName = d.Get("role_name").(string)
	request.EventRW = d.Get("event_rw").(string)

	//Product problem fields need to be added " "
	if d.HasChange("sls_write_role_arn") || d.HasChange("sls_project_arn") {
		slsProjectArn := d.Get("sls_project_arn").(string)
		slsWriteRoleArn := d.Get("sls_write_role_arn").(string)
		if len(slsProjectArn) == 0 {
			slsProjectArn = " "
		}
		if len(slsWriteRoleArn) == 0 {
			slsWriteRoleArn = " "
		}
		request.SlsProjectArn = slsProjectArn
		request.SlsWriteRoleArn = slsWriteRoleArn
	}

	if d.HasChange("oss_key_prefix") {
		request.OssKeyPrefix = d.Get("oss_key_prefix").(string) + " "
	}

	err := resource.Retry(5*time.Minute, func() *resource.RetryError {
		raw, err := client.WithActionTrailClient(func(actiontrailClient *actiontrail.Client) (interface{}, error) {
			return actiontrailClient.UpdateTrail(request)
		})
		if err != nil {
			if IsExceptedErrors(err, []string{InsufficientBucketPolicyException, TrailNeedRamAuthorize}) {
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

	return resourceAlicloudActiontrailRead(d, meta)
}

func resourceAlicloudActiontrailDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	trailService := ActionTrailService{client}
	request := actiontrail.CreateDeleteTrailRequest()
	request.Name = d.Id()

	raw, err := client.WithActionTrailClient(func(actiontrailClient *actiontrail.Client) (interface{}, error) {
		return actiontrailClient.DeleteTrail(request)
	})

	if err != nil {
		if IsExceptedErrors(err, []string{InvalidVpcIDNotFound, ForbiddenVpcNotFound, InvalidTrailNotFound}) {
			return nil
		}
		return WrapErrorf(err, DefaultTimeoutMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
	}

	addDebug(request.GetActionName(), raw)

	return WrapError(trailService.WaitForActionTrail(d.Id(), Deleted, DefaultTimeout))

}
