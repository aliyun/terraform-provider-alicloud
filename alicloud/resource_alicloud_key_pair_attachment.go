package alicloud

import (
	"fmt"
	"github.com/denverdino/aliyungo/ecs"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"strings"
	"time"
)

func resourceAlicloudKeyPairAttachment() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudKeyPairAttachmentCreate,
		Read:   resourceAlicloudKeyPairAttachmentRead,
		Delete: resourceAlicloudKeyPairAttachmentDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"key_name": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateKeyPairName,
				ForceNew:     true,
			},
			"instance_ids": &schema.Schema{
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceAlicloudKeyPairAttachmentCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*AliyunClient).ecsconn
	instanceIds := convertListToJsonString(d.Get("instance_ids").(*schema.Set).List())

	args := &ecs.AttachKeyPairArgs{
		RegionId:    getRegion(d, meta),
		KeyPairName: d.Get("key_name").(string),
		InstanceIds: instanceIds,
	}
	err := resource.Retry(5*time.Minute, func() *resource.RetryError {
		if er := conn.AttachKeyPair(args); er != nil {
			if IsExceptedError(er, KeyPairServiceUnavailable) {
				return resource.RetryableError(fmt.Errorf("Attach Key Pair timeout and got an error: %#v.", er))
			}
			return resource.NonRetryableError(fmt.Errorf("Error Attach KeyPair: %#v", er))
		}
		return nil
	})

	if err != nil {
		return err
	}
	d.SetId(d.Get("key_name").(string) + ":" + instanceIds)

	return resourceAlicloudKeyPairAttachmentRead(d, meta)
}

func resourceAlicloudKeyPairAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*AliyunClient).ecsconn
	keyname := strings.Split(d.Id(), ":")[0]
	keypairs, _, err := conn.DescribeKeyPairs(&ecs.DescribeKeyPairsArgs{
		RegionId:    getRegion(d, meta),
		KeyPairName: keyname,
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
		d.Set("instance_ids", d.Get("instance_ids"))
		return nil
	}

	return fmt.Errorf("Unable to find key pair %s in the current account.", keyname)
}

func resourceAlicloudKeyPairAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*AliyunClient)
	keyname := strings.Split(d.Id(), ":")[0]
	instanceIds := strings.Split(d.Id(), ":")[1]

	return resource.Retry(5*time.Minute, func() *resource.RetryError {
		err := client.ecsconn.DetachKeyPair(&ecs.DetachKeyPairArgs{
			RegionId:    getRegion(d, meta),
			KeyPairName: keyname,
			InstanceIds: instanceIds,
		})
		if err != nil {
			return resource.NonRetryableError(fmt.Errorf("Error DetachKeyPair:%#v", err))
		}

		instance_ids, _, err := client.QueryInstancesWithKeyPair(getRegion(d, meta), instanceIds, d.Id())
		if err != nil {
			return resource.NonRetryableError(err)
		}
		if len(instance_ids) > 0 {
			instanceIds = convertListToJsonString(instance_ids)
			return resource.RetryableError(fmt.Errorf("Detach Key Pair timeout and got an error: %#v.", err))
		}

		return nil
	})
}
