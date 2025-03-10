// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"log"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/blues/jsonata-go"
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
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"arn": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"encodedsaml_metadata_document": {
				Type:     schema.TypeString,
				Required: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return ramSAMLProviderDiffSuppressFunc(old, new)
				},
			},
			"saml_provider_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
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

	action := "CreateSAMLProvider"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	if v, ok := d.GetOk("saml_provider_name"); ok {
		request["SAMLProviderName"] = v
	}

	if v, ok := d.GetOk("description"); ok {
		request["Description"] = v
	}
	request["EncodedSAMLMetadataDocument"] = d.Get("encodedsaml_metadata_document")
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RpcPost("Ims", "2019-08-15", action, query, request, true)
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

	id, _ := jsonpath.Get("$.SAMLProvider.SAMLProviderName", response)
	d.SetId(fmt.Sprint(id))

	return resourceAliCloudRamSamlProviderRead(d, meta)
}

func resourceAliCloudRamSamlProviderRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ramServiceV2 := RamServiceV2{client}

	objectRaw, err := ramServiceV2.DescribeRamSamlProvider(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_ram_saml_provider DescribeRamSamlProvider Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("arn", objectRaw["Arn"])
	d.Set("description", objectRaw["Description"])
	d.Set("update_date", objectRaw["UpdateDate"])
	d.Set("saml_provider_name", objectRaw["SAMLProviderName"])

	e := jsonata.MustCompile("$replace($.EncodedSAMLMetadataDocument, \"\n\", \"\")")
	evaluation, _ := e.Eval(objectRaw)
	d.Set("encodedsaml_metadata_document", evaluation)

	d.Set("saml_provider_name", d.Id())

	return nil
}

func resourceAliCloudRamSamlProviderUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	update := false

	var err error
	action := "UpdateSAMLProvider"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["SAMLProviderName"] = d.Id()

	if d.HasChange("description") {
		update = true
		request["NewDescription"] = d.Get("description")
	}

	if d.HasChange("encodedsaml_metadata_document") {
		update = true
	}
	request["NewEncodedSAMLMetadataDocument"] = d.Get("encodedsaml_metadata_document")
	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("Ims", "2019-08-15", action, query, request, true)
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
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["SAMLProviderName"] = d.Id()

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.RpcPost("Ims", "2019-08-15", action, query, request, true)

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
		if IsExpectedErrors(err, []string{"EntityNotExist.SAMLProvider"}) || NotFoundError(err) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	return nil
}
