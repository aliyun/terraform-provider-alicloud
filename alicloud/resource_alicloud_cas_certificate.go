package alicloud

import (
	"fmt"
	"io/ioutil"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/cas"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func resourceAlicloudCasCertificate() *schema.Resource {
	return &schema.Resource{
		Create:   resourceAlicloudCasCreate,
		Read:     resourceAlicloudCasRead,
		Delete:   resourceAlicloudCasDelete,
		Importer: &schema.ResourceImporter{State: schema.ImportStatePassthrough},

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateCasName,
			},
			"cert": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"key": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceAlicloudCasCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	args := cas.CreateCreateUserCertificateRequest()
	if v, ok := d.GetOk("name"); ok {
		args.Name = v.(string)
	}

	if v, ok := d.GetOk("cert"); ok {
		b, err := ioutil.ReadFile(v.(string))
		if err != nil {
			return WrapErrorf(err, DataDefaultErrorMsg, "cas", args.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		args.Cert = string(b)
	}

	if v, ok := d.GetOk("key"); ok {
		b, err := ioutil.ReadFile(v.(string))
		if err != nil {
			return WrapErrorf(err, DataDefaultErrorMsg, "vpcs", args.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		args.Key = string(b)
	}

	raw, _ := client.WithCasClient(func(casClient *cas.Client) (interface{}, error) {
		return casClient.CreateUserCertificate(args)
	})

	response, _ := raw.(*cas.CreateUserCertificateResponse)
	d.SetId(string(response.CertId))
	return resourceAlicloudCasRead(d, meta)
}

func resourceAlicloudCasRead(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceAlicloudCasDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	args := cas.CreateDeleteUserCertificateRequest()
	args.CertId = requests.Integer(d.Id())

	return resource.Retry(5*time.Minute, func() *resource.RetryError {
		_, err := client.WithCasClient(func(casClient *cas.Client) (interface{}, error) {
			return casClient.DeleteUserCertificate(args)
		})
		if err != nil {
			if IsExceptedError(err, CertNotExist) {
				return nil
			}
			return resource.NonRetryableError(fmt.Errorf("Error deleting cert %s: %#v", d.Id(), err))
		}
		return nil
	})
}
