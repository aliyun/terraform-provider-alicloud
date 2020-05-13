package alicloud

import (
	"fmt"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/kms"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
	"log"
)

func resourceAlicloudKmsKeyVersion() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudKmsKeyVersionCreate,
		Read:   resourceAlicloudKmsKeyVersionRead,
		Delete: resourceAlicloudKmsKeyVersionDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"creation_date": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"key_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"key_version_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAlicloudKmsKeyVersionCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	request := kms.CreateCreateKeyVersionRequest()
	request.KeyId = d.Get("key_id").(string)
	raw, err := client.WithKmsClient(func(kmsClient *kms.Client) (interface{}, error) {
		return kmsClient.CreateKeyVersion(request)
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_kms_key_version", request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw)
	response, _ := raw.(*kms.CreateKeyVersionResponse)
	d.SetId(fmt.Sprintf("%v:%v", response.KeyVersion.KeyId, response.KeyVersion.KeyVersionId))

	return resourceAlicloudKmsKeyVersionRead(d, meta)
}
func resourceAlicloudKmsKeyVersionRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	kmsService := KmsService{client}
	object, err := kmsService.DescribeKmsKeyVersion(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	d.Set("key_id", parts[0])
	d.Set("creation_date", object.KeyVersion.CreationDate)
	return nil
}
func resourceAlicloudKmsKeyVersionDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[WARN] Cannot destroy resourceAlicloudKmsKeyVersion. Terraform will remove this resource from the state file, however resources may remain.")
	return nil
}
