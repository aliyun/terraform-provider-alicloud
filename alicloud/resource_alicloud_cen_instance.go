package alicloud

import (
	"fmt"
	"log"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/cbn"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
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
			Create: schema.DefaultTimeout(6 * time.Minute),
			Delete: schema.DefaultTimeout(6 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"cen_instance_name": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"name"},
			},
			"name": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				Deprecated:    "Field 'name' has been deprecated from version 1.98.0. Use 'cen_instance_name' instead.",
				ConflictsWith: []string{"cen_instance_name"},
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"protection_level": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "REDUCED",
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tags": tagsSchema(),
		},
	}
}

func resourceAlicloudCenInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cbnService := CbnService{client}

	request := cbn.CreateCreateCenRequest()
	if v, ok := d.GetOk("cen_instance_name"); ok {
		request.Name = v.(string)
	} else if v, ok := d.GetOk("name"); ok {
		request.Name = v.(string)
	}

	if v, ok := d.GetOk("description"); ok {
		request.Description = v.(string)
	}

	if v, ok := d.GetOk("protection_level"); ok {
		request.ProtectionLevel = v.(string)
	}

	wait := incrementalWait(3*time.Second, 3*time.Second)
	err := resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		raw, err := client.WithCbnClient(func(cbnClient *cbn.Client) (interface{}, error) {
			return cbnClient.CreateCen(request)
		})
		if err != nil {
			if IsExpectedErrors(err, []string{"Operation.Blocking"}) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(request.GetActionName(), raw)
		response, _ := raw.(*cbn.CreateCenResponse)
		d.SetId(fmt.Sprintf("%v", response.CenId))
		return nil
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_cen_instance", request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	stateConf := BuildStateConf([]string{}, []string{"Active"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, cbnService.CenInstanceStateRefreshFunc(d.Id(), []string{"Deleting"}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAlicloudCenInstanceUpdate(d, meta)
}
func resourceAlicloudCenInstanceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cbnService := CbnService{client}
	object, err := cbnService.DescribeCenInstance(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_cen_instance cbnService.DescribeCenInstance Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	d.Set("cen_instance_name", object.Name)
	d.Set("name", object.Name)
	d.Set("description", object.Description)
	d.Set("protection_level", object.ProtectionLevel)
	d.Set("status", object.Status)

	tags := make(map[string]string)
	for _, t := range object.Tags.Tag {
		tags[t.Key] = t.Value
	}
	d.Set("tags", tags)
	return nil
}
func resourceAlicloudCenInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cbnService := CbnService{client}
	d.Partial(true)

	if d.HasChange("tags") {
		if err := cbnService.SetResourceTags(d, "cen"); err != nil {
			return WrapError(err)
		}
		d.SetPartial("tags")
	}
	update := false
	request := cbn.CreateModifyCenAttributeRequest()
	request.CenId = d.Id()
	if !d.IsNewResource() && d.HasChange("cen_instance_name") {
		update = true
		request.Name = d.Get("cen_instance_name").(string)
	}
	if !d.IsNewResource() && d.HasChange("name") {
		update = true
		request.Name = d.Get("name").(string)
	}
	if !d.IsNewResource() && d.HasChange("description") {
		update = true
		request.Description = d.Get("description").(string)
	}
	if !d.IsNewResource() && d.HasChange("name") {
		update = true
		request.Name = d.Get("name").(string)
	}
	if !d.IsNewResource() && d.HasChange("protection_level") {
		update = true
		request.ProtectionLevel = d.Get("protection_level").(string)
	}
	if update {
		raw, err := client.WithCbnClient(func(cbnClient *cbn.Client) (interface{}, error) {
			return cbnClient.ModifyCenAttribute(request)
		})
		addDebug(request.GetActionName(), raw)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		d.SetPartial("cen_instance_name")
		d.SetPartial("name")
		d.SetPartial("description")
		d.SetPartial("protection_level")
	}
	d.Partial(false)
	return resourceAlicloudCenInstanceRead(d, meta)
}
func resourceAlicloudCenInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cbnService := CbnService{client}
	request := cbn.CreateDeleteCenRequest()
	request.CenId = d.Id()
	raw, err := client.WithCbnClient(func(cbnClient *cbn.Client) (interface{}, error) {
		return cbnClient.DeleteCen(request)
	})
	addDebug(request.GetActionName(), raw)
	if err != nil {
		if IsExpectedErrors(err, []string{"ParameterCenInstanceId"}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	stateConf := BuildStateConf([]string{}, []string{}, d.Timeout(schema.TimeoutDelete), 5*time.Second, cbnService.CenInstanceStateRefreshFunc(d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}
	return nil
}
