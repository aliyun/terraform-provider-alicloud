package alicloud

import (
	"fmt"
	"io/ioutil"
	"strings"
	"time"

	"os"

	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceAlicloudKeyPair() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudKeyPairCreate,
		Read:   resourceAlicloudKeyPairRead,
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
				ValidateFunc:  validateKeyPairName,
				ConflictsWith: []string{"key_name_prefix"},
			},
			"key_name_prefix": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validateKeyPairPrefix,
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
		args := ecs.CreateImportKeyPairRequest()
		args.KeyPairName = keyName
		args.PublicKeyBody = publicKey.(string)
		raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
			return ecsClient.ImportKeyPair(args)
		})
		if err != nil {
			return fmt.Errorf("Error Import KeyPair: %s", err)
		}
		keypair, _ := raw.(*ecs.ImportKeyPairResponse)
		d.SetId(keypair.KeyPairName)
	} else {
		args := ecs.CreateCreateKeyPairRequest()
		args.KeyPairName = keyName
		raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
			return ecsClient.CreateKeyPair(args)
		})
		if err != nil {
			return fmt.Errorf("Error Create KeyPair: %s", err)
		}
		keypair, _ := raw.(*ecs.CreateKeyPairResponse)
		d.SetId(keypair.KeyPairName)
		if file, ok := d.GetOk("key_file"); ok {
			ioutil.WriteFile(file.(string), []byte(keypair.PrivateKeyBody), 0600)
			os.Chmod(file.(string), 0400)
		}
	}

	return resourceAlicloudKeyPairRead(d, meta)
}

func resourceAlicloudKeyPairRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ecsService := EcsService{client}

	keypair, err := ecsService.DescribeKeyPair(d.Id())
	if err != nil {
		if NotFoundError(err) || IsExceptedError(err, KeyPairNotFound) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error Retrieving KeyPair: %s", err)
	}
	d.Set("key_name", keypair.KeyPairName)
	d.Set("fingerprint", keypair.KeyPairFingerPrint)
	return nil
}

func resourceAlicloudKeyPairDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ecsService := EcsService{client}

	deldArgs := ecs.CreateDeleteKeyPairsRequest()
	deldArgs.KeyPairNames = convertListToJsonString(append(make([]interface{}, 0, 1), d.Id()))

	return resource.Retry(5*time.Minute, func() *resource.RetryError {

		_, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
			return ecsClient.DeleteKeyPairs(deldArgs)
		})
		if err != nil {
			if IsExceptedError(err, KeyPairNotFound) {
				return nil
			}
		}

		_, err = ecsService.DescribeKeyPair(d.Id())
		if err != nil {
			if NotFoundError(err) || IsExceptedError(err, KeyPairNotFound) {
				return nil
			}
		}
		return resource.RetryableError(fmt.Errorf("Delete Key Pair timeout and got an error: %#v.", err))
	})
}
