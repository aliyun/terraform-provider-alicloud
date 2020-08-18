package alicloud

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	"github.com/aliyun/fc-go-sdk"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAlicloudFCVersion() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudFCVersionCreate,
		Read:   resourceAlicloudFCVersionRead,
		Update: resourceAlicloudFCVersionUpdate,
		Delete: resourceAlicloudFCVersionDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"service": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(1, 128),
			},
			"version_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"create": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"last_modified": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAlicloudFCVersionCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	serviceName := d.Get("service").(string)

	request := fc.NewPublishServiceVersionInput(serviceName)
	if v, ok := d.GetOk("description"); ok && v.(string) != "" {
		request.WithDescription(v.(string))
	}

	var response *fc.PublishServiceVersionOutput
	var requestInfo *fc.Client
	if err := resource.Retry(2*time.Minute, func() *resource.RetryError {
		raw, err := client.WithFcClient(func(fcClient *fc.Client) (interface{}, error) {
			requestInfo = fcClient
			return fcClient.PublishServiceVersion(request)
		})
		if err != nil {
			if IsExpectedErrors(err, []string{"AccessDenied"}) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug("PublishServiceVersion", raw, requestInfo, request)
		response, _ = raw.(*fc.PublishServiceVersionOutput)
		return nil
	}); err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_fc_version", "PublishServiceVersion", FcGoSdk)
	}

	d.SetId(fmt.Sprintf("%s%s%s", serviceName, COLON_SEPARATED, *response.VersionID))

	return resourceAlicloudFCVersionRead(d, meta)
}

func resourceAlicloudFCVersionRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	fcService := FcService{client}

	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	versions, err := fcService.DescribeFcVersion(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("service", parts[0])
	d.Set("description", versions.Versions[0].Description)
	d.Set("version_id", versions.Versions[0].VersionID)
	d.Set("create", versions.Versions[0].CreatedTime)
	d.Set("last_modified", versions.Versions[0].LastModifiedTime)

	return nil
}

func resourceAlicloudFCVersionUpdate(d *schema.ResourceData, meta interface{}) error {
	return WrapError(Error("version is readonly"))
}

func resourceAlicloudFCVersionDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	fcService := FcService{client}

	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	request := fc.NewDeleteServiceVersionInput(parts[0], parts[1])
	var requestInfo *fc.Client
	raw, err := client.WithFcClient(func(fcClient *fc.Client) (interface{}, error) {
		requestInfo = fcClient
		return fcClient.DeleteServiceVersion(request)
	})
	if err != nil {
		if IsExpectedErrors(err, []string{"ServiceNotFound", "VersionNotFound"}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), "DeleteServiceVersion", FcGoSdk)
	}
	addDebug("DeleteServiceVersion", raw, requestInfo, request)
	return WrapError(fcService.WaitForFcVersion(d.Id(), Deleted, DefaultTimeoutMedium))
}
