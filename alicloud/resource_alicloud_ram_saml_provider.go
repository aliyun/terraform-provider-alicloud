package alicloud

import (
	"fmt"
	"log"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudRamSamlProvider() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudRamSamlProviderCreate,
		Read:   resourceAliCloudRamSamlProviderRead,
		Update: resourceAliCloudRamSamlProviderUpdate,
		Delete: resourceAliCloudRamSamlProviderDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"saml_provider_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"encodedsaml_metadata_document": {
				Type:     schema.TypeString,
				Required: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return ramSAMLProviderDiffSuppressFunc(old, new)
				},
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"arn": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"update_date": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAliCloudRamSamlProviderCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "CreateSAMLProvider"
	request := make(map[string]interface{})
	var err error

	request["SAMLProviderName"] = d.Get("saml_provider_name")
	request["EncodedSAMLMetadataDocument"] = d.Get("encodedsaml_metadata_document")

	if v, ok := d.GetOk("description"); ok {
		request["Description"] = v
	}

	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutCreate)), func() *resource.RetryError {
		response, err = client.RpcPost("Ims", "2019-08-15", action, nil, request, true)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_ram_saml_provider", action, AlibabaCloudSdkGoERROR)
	}

	if resp, err := jsonpath.Get("$.SAMLProvider", response); err != nil || resp == nil {
		return WrapErrorf(err, IdMsg, "alicloud_ram_saml_provider")
	} else {
		samlProviderName := resp.(map[string]interface{})["SAMLProviderName"]
		d.SetId(fmt.Sprint(samlProviderName))
	}

	return resourceAliCloudRamSamlProviderRead(d, meta)
}

func resourceAliCloudRamSamlProviderRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	imsService := ImsService{client}

	object, err := imsService.DescribeRamSamlProvider(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_ram_saml_provider imsService.DescribeRamSamlProvider Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("saml_provider_name", object["SAMLProviderName"])
	d.Set("encodedsaml_metadata_document", object["EncodedSAMLMetadataDocument"])
	d.Set("description", object["Description"])
	d.Set("arn", object["Arn"])
	d.Set("update_date", object["UpdateDate"])

	return nil
}

func resourceAliCloudRamSamlProviderUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	var err error
	update := false

	request := map[string]interface{}{
		"SAMLProviderName": d.Id(),
	}

	if d.HasChange("encodedsaml_metadata_document") {
		update = true
	}
	request["NewEncodedSAMLMetadataDocument"] = d.Get("encodedsaml_metadata_document")

	if d.HasChange("description") {
		update = true
	}
	if v, ok := d.GetOk("description"); ok {
		request["NewDescription"] = v
	}

	if update {
		action := "UpdateSAMLProvider"
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *resource.RetryError {
			response, err = client.RpcPost("Ims", "2019-08-15", action, nil, request, true)
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
	}

	return resourceAliCloudRamSamlProviderRead(d, meta)
}

func resourceAliCloudRamSamlProviderDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "DeleteSAMLProvider"
	var response map[string]interface{}
	var err error

	request := map[string]interface{}{
		"SAMLProviderName": d.Id(),
	}

	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutDelete)), func() *resource.RetryError {
		response, err = client.RpcPost("Ims", "2019-08-15", action, nil, request, true)
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
		if IsExpectedErrors(err, []string{"EntityNotExist.SAMLProviderError"}) || NotFoundError(err) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	return nil
}
