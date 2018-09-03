package alicloud

import (
	"fmt"
	"strings"

	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/cbn"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceAlicloudCen() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudCenCreate,
		Read:   resourceAlicloudCenRead,
		Update: resourceAlicloudCenUpdate,
		Delete: resourceAlicloudCenDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
					value := v.(string)
					if len(value) < 2 || len(value) > 128 {
						errors = append(errors, fmt.Errorf("%s cannot be shorter than 2 characters or longer than 128 characters", k))
					}

					if strings.HasPrefix(value, "http://") || strings.HasPrefix(value, "https://") {
						errors = append(errors, fmt.Errorf("%s cannot starts with http:// or https://", k))
					}

					return
				},
			},
			"description": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
					value := v.(string)
					if len(value) < 2 || len(value) > 256 {
						errors = append(errors, fmt.Errorf("%s cannot be shorter than 2 characters or longer than 256 characters", k))
					}

					if strings.HasPrefix(value, "http://") || strings.HasPrefix(value, "https://") {
						errors = append(errors, fmt.Errorf("%s cannot starts with http:// or https://", k))
					}

					return
				},
			},
		},
	}
}

func resourceAlicloudCenCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*AliyunClient)

	var cen *cbn.CreateCenResponse
	err := resource.Retry(3*time.Minute, func() *resource.RetryError {
		args, err := buildAliCloudCenArgs(d, meta)
		if err != nil {
			return resource.NonRetryableError(fmt.Errorf("Building CreateVpcRequest got an error: %#v", err))
		}
		resp, err := client.cenconn.CreateCen(args)
		if err != nil {
			if IsExceptedError(err, CenQuotaExceeded) {
				return resource.NonRetryableError(fmt.Errorf("The number of CEN exceeds the maximum value that can be created."))
			}
			if IsExceptedError(err, UnknownError) {
				return resource.RetryableError(fmt.Errorf("Create vpc timeout and got an error: %#v.", err))
			}
			return resource.NonRetryableError(err)
		}
		cen = resp
		return nil
	})
	if err != nil {
		return fmt.Errorf("Create cen got an error :%#v", err)
	}

	d.SetId(cen.CenId)

	err = client.WaitForCen(d.Id(), Active, 60)
	if err != nil {
		return fmt.Errorf("Timeout when WaitForCenAvailable")
	}

	return resourceAlicloudCenUpdate(d, meta)
}

func resourceAlicloudCenRead(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*AliyunClient)
	resp, err := client.DescribeCen(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return err
	}

	d.Set("name", resp.Name)
	d.Set("description", resp.Description)

	return nil
}

func resourceAlicloudCenUpdate(d *schema.ResourceData, meta interface{}) error {

	d.Partial(true)

	attributeUpdate := false
	request := cbn.CreateModifyCenAttributeRequest()
	request.CenId = d.Id()

	if d.HasChange("name") {
		d.SetPartial("name")
		request.Name = d.Get("name").(string)

		attributeUpdate = true
	}

	if d.HasChange("description") {
		d.SetPartial("description")
		request.Description = d.Get("description").(string)

		attributeUpdate = true
	}

	if attributeUpdate {
		if _, err := meta.(*AliyunClient).cenconn.ModifyCenAttribute(request); err != nil {
			return err
		}
	}

	d.Partial(false)

	return resourceAlicloudCenRead(d, meta)
}

func resourceAlicloudCenDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*AliyunClient)
	request := cbn.CreateDeleteCenRequest()
	request.CenId = d.Id()

	return resource.Retry(5*time.Minute, func() *resource.RetryError {
		_, err := client.cenconn.DeleteCen(request)

		if err != nil {
			if IsExceptedError(err, ParameterCenInstanceIdNotExist) {
				return nil
			}
			return resource.RetryableError(fmt.Errorf("Delete CEN timeout and got an error: %#v.", err))
		}

		if _, err := client.DescribeCen(d.Id()); err != nil {
			if NotFoundError(err) {
				return nil
			}
			return resource.NonRetryableError(err)
		}

		return nil
	})
}

func buildAliCloudCenArgs(d *schema.ResourceData, meta interface{}) (*cbn.CreateCenRequest, error) {
	request := cbn.CreateCreateCenRequest()

	if v := d.Get("name").(string); v != "" {
		request.Name = v
	}

	if v := d.Get("description").(string); v != "" {
		request.Description = v
	}

	return request, nil
}
