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

func resourceAlicloudFCAlias() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudFCAliasCreate,
		Read:   resourceAlicloudFCAliasRead,
		Update: resourceAlicloudFCAliasUpdate,
		Delete: resourceAlicloudFCAliasDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"service": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"name": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(1, 128),
			},
			"version_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(1, 128),
			},
			"additional_weight": {
				Type: schema.TypeMap,
				Elem: &schema.Schema{
					Type: schema.TypeFloat,
				},
				Optional: true,
			},
		},
	}
}

func resourceAlicloudFCAliasCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	serviceName := d.Get("service").(string)
	var name string
	if v, ok := d.GetOk("name"); ok {
		name = v.(string)
	} else {
		name = resource.UniqueId()
	}

	object := fc.AliasCreateObject{
		AliasName:               StringPointer(name),
		VersionID:               StringPointer(d.Get("version_id").(string)),
		AdditionalVersionWeight: make(map[string]float64),
	}
	if v, ok := d.GetOk("description"); ok && v.(string) != "" {
		object.Description = StringPointer(v.(string))
	}
	if v, ok := d.GetOk("additional_weight"); ok {
		for version, weight := range v.(map[string]interface{}) {
			object.AdditionalVersionWeight[version] = weight.(float64)
		}
	}
	request := &fc.CreateAliasInput{
		ServiceName:       StringPointer(serviceName),
		AliasCreateObject: object,
	}
	var response *fc.CreateAliasOutput
	var requestInfo *fc.Client
	if err := resource.Retry(2*time.Minute, func() *resource.RetryError {
		raw, err := client.WithFcClient(func(fcClient *fc.Client) (interface{}, error) {
			requestInfo = fcClient
			return fcClient.CreateAlias(request)
		})
		if err != nil {
			if IsExpectedErrors(err, []string{"AccessDenied"}) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug("CreateAlias", raw, requestInfo, request)
		response, _ = raw.(*fc.CreateAliasOutput)
		return nil

	}); err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_fc_alias", "CreateAlias", FcGoSdk)
	}

	d.SetId(fmt.Sprintf("%s%s%s", serviceName, COLON_SEPARATED, *response.AliasName))

	return resourceAlicloudFCAliasRead(d, meta)
}

func resourceAlicloudFCAliasRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	fcService := FcService{client}

	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	alias, err := fcService.DescribeFcAlias(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("service", parts[0])
	d.Set("name", alias.AliasName)
	d.Set("version_id", alias.VersionID)
	d.Set("description", alias.Description)
	d.Set("additional_weight", alias.AdditionalVersionWeight)

	return nil
}

func resourceAlicloudFCAliasUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	updated := false
	updateInput := &fc.UpdateAliasInput{}

	if d.HasChange("version_id") {
		updateInput.VersionID = StringPointer(d.Get("version_id").(string))
		updated = true
	}
	if d.HasChange("description") {
		updateInput.Description = StringPointer(d.Get("description").(string))
		updated = true
	}
	if d.HasChange("additional_weight") {
		updateInput.AdditionalVersionWeight = make(map[string]float64)
		if weights, ok := d.GetOk("additional_weight"); ok {
			for version, weight := range weights.(map[string]interface{}) {
				updateInput.AdditionalVersionWeight[version] = weight.(float64)
			}
		}
		updated = true
	}

	if updated {
		parts, err := ParseResourceId(d.Id(), 2)
		if err != nil {
			return WrapError(err)
		}
		updateInput.ServiceName = StringPointer(parts[0])
		updateInput.AliasName = StringPointer(parts[1])
		var requestInfo *fc.Client
		raw, err := client.WithFcClient(func(fcClient *fc.Client) (interface{}, error) {
			requestInfo = fcClient
			return fcClient.UpdateAlias(updateInput)
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), "UpdateAlias", FcGoSdk)
		}
		addDebug("UpdateAlias", raw, requestInfo, updateInput)
	}

	return resourceAlicloudFCAliasRead(d, meta)
}

func resourceAlicloudFCAliasDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	fcService := FcService{client}

	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	request := &fc.DeleteAliasInput{
		ServiceName: StringPointer(parts[0]),
		AliasName:   StringPointer(parts[1]),
	}
	var requestInfo *fc.Client
	raw, err := client.WithFcClient(func(fcClient *fc.Client) (interface{}, error) {
		requestInfo = fcClient
		return fcClient.DeleteAlias(request)
	})
	if err != nil {
		if IsExpectedErrors(err, []string{"ServiceNotFound", "AliasNotFound"}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), "DeleteAlias", FcGoSdk)
	}
	addDebug("DeleteAlias", raw, requestInfo, request)
	return WrapError(fcService.WaitForFcAlias(d.Id(), Deleted, DefaultTimeoutMedium))
}
