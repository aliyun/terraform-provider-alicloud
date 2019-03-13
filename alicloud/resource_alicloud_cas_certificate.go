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
			return WrapError(err)
		}
		args.Cert = string(b)
	}

	if v, ok := d.GetOk("key"); ok {
		b, err := ioutil.ReadFile(v.(string))
		if err != nil {
			return WrapError(err)
		}
		args.Key = string(b)
	}

	raw, err := client.WithCasClient(func(casClient *cas.Client) (interface{}, error) {
		return casClient.CreateUserCertificate(args)
	})

	if err != nil {
		return WrapError(err)
	}

	response, _ := raw.(*cas.CreateUserCertificateResponse)
	d.SetId(string(response.CertId))
	return resourceAlicloudCasRead(d, meta)
}

func resourceAlicloudCasRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	casService := &CasService{client: client}
	cert, err := casService.DescribeCas(d.Id())
	if err != nil {
		return WrapError(err)
	}

	d.Set("name", cert.Name)

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
