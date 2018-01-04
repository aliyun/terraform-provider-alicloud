package alicloud

import (
	"fmt"
	"github.com/denverdino/aliyungo/ecs"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"io/ioutil"
	"strings"
	"time"
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
			"key_name": &schema.Schema{
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ForceNew:      true,
				ValidateFunc:  validateKeyPairName,
				ConflictsWith: []string{"key_name_prefix"},
			},
			"key_name_prefix": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validateKeyPairPrefix,
			},
			"public_key": &schema.Schema{
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
			"key_file": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"finger_print": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAlicloudKeyPairCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*AliyunClient).ecsconn

	var keyName string
	if v, ok := d.GetOk("key_name"); ok {
		keyName = v.(string)
	} else if v, ok := d.GetOk("key_name_prefix"); ok {
		keyName = resource.PrefixedUniqueId(v.(string))
	} else {
		keyName = resource.UniqueId()
	}

	if publicKey, ok := d.GetOk("public_key"); ok {
		keypair, err := conn.ImportKeyPair(&ecs.ImportKeyPairArgs{
			RegionId:      getRegion(d, meta),
			KeyPairName:   keyName,
			PublicKeyBody: publicKey.(string),
		})
		if err != nil {
			return fmt.Errorf("Error Import KeyPair: %s", err)
		}

		d.SetId(keypair.KeyPairName)
	} else {
		keypair, err := conn.CreateKeyPair(&ecs.CreateKeyPairArgs{
			RegionId:    getRegion(d, meta),
			KeyPairName: keyName,
		})
		if err != nil {
			return fmt.Errorf("Error Create KeyPair: %s", err)
		}

		d.SetId(keypair.KeyPairName)
		if file, ok := d.GetOk("key_file"); ok {
			ioutil.WriteFile(file.(string), []byte(keypair.PrivateKeyBody), 400)
		}
	}

	return resourceAlicloudKeyPairRead(d, meta)
}

func resourceAlicloudKeyPairRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*AliyunClient).ecsconn

	keypairs, _, err := conn.DescribeKeyPairs(&ecs.DescribeKeyPairsArgs{
		RegionId:    getRegion(d, meta),
		KeyPairName: d.Id(),
	})
	if err != nil {
		if IsExceptedError(err, KeyPairNotFound) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error Retrieving KeyPair: %s", err)
	}

	if len(keypairs) > 0 {
		d.Set("key_name", keypairs[0].KeyPairName)
		d.Set("fingerprint", keypairs[0].KeyPairFingerPrint)
		return nil
	}

	return fmt.Errorf("Unable to find key pair within: %#v", keypairs)
}

func resourceAlicloudKeyPairDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*AliyunClient)

	instance_ids, _, err := client.QueryInstancesWithKeyPair(getRegion(d, meta), "", d.Id())
	if err != nil {
		return err
	}
	detachArgs := &ecs.DetachKeyPairArgs{
		RegionId:    getRegion(d, meta),
		KeyPairName: d.Id(),
	}

	return resource.Retry(5*time.Minute, func() *resource.RetryError {

		// Detach keypair from its all instances before removing it.
		if len(instance_ids) > 0 {
			detachArgs.InstanceIds = convertListToJsonString(instance_ids)
			if err := client.ecsconn.DetachKeyPair(detachArgs); err != nil {
				return resource.NonRetryableError(fmt.Errorf("Error DetachKeyPair:%#v", err))
			}
		}
		instance_ids, _, err = client.QueryInstancesWithKeyPair(getRegion(d, meta), "", d.Id())
		if err != nil {
			return resource.NonRetryableError(err)
		}
		if len(instance_ids) > 0 {
			return resource.RetryableError(fmt.Errorf("Delete Key Pair timeout and got an error: %#v.", err))
		}

		err := client.ecsconn.DeleteKeyPairs(&ecs.DeleteKeyPairsArgs{
			RegionId:     getRegion(d, meta),
			KeyPairNames: convertListToJsonString(append(make([]interface{}, 0, 1), d.Id())),
		})
		if err != nil {
			if IsExceptedError(err, KeyPairNotFound) {
				return nil
			}
		}

		keypairs, _, err := client.ecsconn.DescribeKeyPairs(&ecs.DescribeKeyPairsArgs{
			RegionId:    getRegion(d, meta),
			KeyPairName: d.Id(),
		})
		if len(keypairs) > 0 {
			return resource.RetryableError(fmt.Errorf("Delete Key Pair timeout and got an error: %#v.", err))
		}

		return nil
	})
}
