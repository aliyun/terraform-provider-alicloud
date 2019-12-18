package alicloud

import (
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/cbn"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func resourceAlicloudCenInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudCenInstanceCreate,
		Read:   resourceAlicloudCenInstanceRead,
		Update: resourceAlicloudCenInstanceUpdate,
		Delete: resourceAlicloudCenInstanceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(3 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(2, 128),
			},
			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(2, 256),
			},
		},
	}
}

func resourceAlicloudCenInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cenService := CenService{client}
	request := cbn.CreateCreateCenRequest()
	request.Name = d.Get("name").(string)
	request.Description = d.Get("description").(string)

	var response *cbn.CreateCenResponse
	wait := incrementalWait(3*time.Second, 5*time.Second)

	err := resource.Retry(5*time.Minute, func() *resource.RetryError {
		req := *request
		raw, err := client.WithCenClient(func(cbnClient *cbn.Client) (interface{}, error) {
			return cbnClient.CreateCen(&req)
		})
		if err != nil {
			if IsExceptedErrors(err, []string{OperationBlocking, UnknownError}) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		response, _ = raw.(*cbn.CreateCenResponse)
		return nil
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_cen_instance", request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), response)
	d.SetId(response.CenId)

	stateConf := BuildStateConf([]string{"Creating"}, []string{"Active"}, d.Timeout(schema.TimeoutCreate), 3*time.Second, cenService.CenInstanceStateRefreshFunc(d.Id(), []string{"Deleting"}))

	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAlicloudCenInstanceRead(d, meta)
}

func resourceAlicloudCenInstanceRead(d *schema.ResourceData, meta interface{}) error {
	cenService := CenService{meta.(*connectivity.AliyunClient)}
	object, err := cenService.DescribeCenInstance(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("name", object.Name)
	d.Set("description", object.Description)

	return nil
}

func resourceAlicloudCenInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	update := false
	request := cbn.CreateModifyCenAttributeRequest()
	request.CenId = d.Id()

	if d.HasChange("name") {
		request.Name = d.Get("name").(string)
		update = true
	}

	if d.HasChange("description") {
		request.Description = d.Get("description").(string)
		update = true
	}

	if update {
		client := meta.(*connectivity.AliyunClient)
		_, err := client.WithCenClient(func(cbnClient *cbn.Client) (interface{}, error) {
			return cbnClient.ModifyCenAttribute(request)
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
	}

	return resourceAlicloudCenInstanceRead(d, meta)
}

func resourceAlicloudCenInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cenService := CenService{client}
	request := cbn.CreateDeleteCenRequest()
	request.CenId = d.Id()

	raw, err := client.WithCenClient(func(cbnClient *cbn.Client) (interface{}, error) {
		return cbnClient.DeleteCen(request)
	})

	if err != nil {
		if IsExceptedError(err, ParameterCenInstanceIdNotExist) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)

	stateConf := BuildStateConf([]string{"Creating", "Active", "Deleting"}, []string{}, d.Timeout(schema.TimeoutDelete), 3*time.Second, cenService.CenInstanceStateRefreshFunc(d.Id(), []string{}))

	if _, err = stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}
	return nil
}
