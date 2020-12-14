package alicloud

import (
	"io/ioutil"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	"os"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAlicloudKeyPair() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudKeyPairCreate,
		Read:   resourceAlicloudKeyPairRead,
		Update: resourceAlicloudKeyPairUpdate,
		Delete: resourceAlicloudKeyPairDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"key_name": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ForceNew:      true,
				ValidateFunc:  validation.StringLenBetween(2, 128),
				ConflictsWith: []string{"key_name_prefix"},
			},
			"key_name_prefix": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(0, 100),
			},
			"resource_group_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"public_key": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				StateFunc: func(v interface{}) string {
					switch v.(type) {
					case string:
						return strings.TrimSpace(v.(string))
					default:
						return ""
					}
				},
			},
			"key_file": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"finger_print": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tags": tagsSchema(),
		},
	}
}

func resourceAlicloudKeyPairCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	var keyName string
	if v, ok := d.GetOk("key_name"); ok {
		keyName = v.(string)
	} else if v, ok := d.GetOk("key_name_prefix"); ok {
		keyName = resource.PrefixedUniqueId(v.(string))
	} else {
		keyName = resource.UniqueId()
	}

	if publicKey, ok := d.GetOk("public_key"); ok {
		request := ecs.CreateImportKeyPairRequest()
		request.RegionId = client.RegionId
		request.KeyPairName = keyName
		request.PublicKeyBody = publicKey.(string)
		if rsg, ok := d.GetOk("resource_group_id"); ok && rsg.(string) != "" {
			request.ResourceGroupId = rsg.(string)
		}
		raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
			return ecsClient.ImportKeyPair(request)
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, "alicloud_key_pair", request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		object, _ := raw.(*ecs.ImportKeyPairResponse)
		d.SetId(object.KeyPairName)
	} else {
		request := ecs.CreateCreateKeyPairRequest()
		request.RegionId = client.RegionId
		request.KeyPairName = keyName
		if v, ok := d.GetOk("resource_group_id"); ok && v.(string) != "" {
			request.ResourceGroupId = v.(string)
		}
		raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
			return ecsClient.CreateKeyPair(request)
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, "alicloud_key_pair", request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		keyPair, _ := raw.(*ecs.CreateKeyPairResponse)
		d.SetId(keyPair.KeyPairName)
		if file, ok := d.GetOk("key_file"); ok {
			ioutil.WriteFile(file.(string), []byte(keyPair.PrivateKeyBody), 0600)
			os.Chmod(file.(string), 0400)
		}
	}

	return resourceAlicloudKeyPairUpdate(d, meta)
}
func resourceAlicloudKeyPairUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	err := setTags(client, TagResourceKeypair, d)
	if err != nil {
		return WrapError(err)
	}
	return resourceAlicloudKeyPairRead(d, meta)
}

func resourceAlicloudKeyPairRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ecsService := EcsService{client}

	keyPair, err := ecsService.DescribeKeyPair(d.Id())
	if err != nil {
		if NotFoundError(err) || IsExpectedErrors(err, []string{"InvalidKeyPair.NotFound"}) {
			log.Printf("[DEBUG] Resource alicloud_key_pair ecsService.DescribeKeyPair Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	d.Set("key_name", keyPair.KeyPairName)
	d.Set("resource_group_id", keyPair.ResourceGroupId)
	d.Set("finger_print", keyPair.KeyPairFingerPrint)
	tags, err := ecsService.ListTagResources(d.Id(), "keypair")
	if err != nil {
		return WrapError(err)
	} else {
		d.Set("tags", tagsToMap(tags))
	}
	return nil
}

func resourceAlicloudKeyPairDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ecsService := EcsService{client}

	request := ecs.CreateDeleteKeyPairsRequest()
	request.RegionId = client.RegionId
	request.KeyPairNames = convertListToJsonString(append(make([]interface{}, 0, 1), d.Id()))

	err := resource.Retry(5*time.Minute, func() *resource.RetryError {
		raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
			return ecsClient.DeleteKeyPairs(request)
		})
		if err != nil {
			if IsExpectedErrors(err, []string{"InvalidKeyPair.NotFound"}) {
				return nil
			}
			return resource.RetryableError(err)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		return nil
	})
	if err != nil {
		WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	return WrapError(ecsService.WaitForKeyPair(d.Id(), Deleted, DefaultTimeoutMedium))
}
