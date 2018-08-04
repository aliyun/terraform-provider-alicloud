package alicloud

import (
	"fmt"
	"strings"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
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
	keyname := d.Get("key_name").(string)
	instanceIds := d.Get("instance_ids").(*schema.Set).List()

	if err := meta.(*AliyunClient).AttachKeyPair(keyname, instanceIds); err != nil {
		return err
	}
	d.SetId(keyname + ":" + convertListToJsonString(instanceIds))

	return resourceAlicloudKeyPairAttachmentRead(d, meta)
}

func resourceAlicloudKeyPairAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*AliyunClient)
	keyname := strings.Split(d.Id(), ":")[0]
	keypair, err := client.DescribeKeyPair(keyname)

	if err != nil {
		if NotFoundError(err) || IsExceptedError(err, KeyPairNotFound) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error Retrieving KeyPair: %s", err)
	}

	d.Set("key_name", keypair.KeyPairName)
	if ids, ok := d.GetOk("instance_ids"); ok {
		d.Set("instance_ids", ids)
	} else {
		ids, _, err := client.QueryInstancesWithKeyPair("", keyname)
		if err != nil {
			return fmt.Errorf("Describe instances by keypair %s got an error: %#v.", keyname, err)
		}
		d.Set("instance_ids", ids)
	}
	return nil
}

func resourceAlicloudKeyPairAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*AliyunClient)
	keyname := strings.Split(d.Id(), ":")[0]
	instanceIds := strings.Split(d.Id(), ":")[1]

	req := ecs.CreateDetachKeyPairRequest()
	req.KeyPairName = keyname

	return resource.Retry(5*time.Minute, func() *resource.RetryError {
		req.InstanceIds = instanceIds
		_, err := client.ecsconn.DetachKeyPair(req)
		if err != nil {
			return resource.NonRetryableError(fmt.Errorf("Error DetachKeyPair:%#v", err))
		}

		instance_ids, _, err := client.QueryInstancesWithKeyPair(instanceIds, d.Id())
		if err != nil {
			return resource.NonRetryableError(err)
		}
		if len(instance_ids) > 0 {
			var ids []interface{}
			for _, id := range instance_ids {
				ids = append(ids, id)
			}
			instanceIds = convertListToJsonString(ids)
			return resource.RetryableError(fmt.Errorf("Detach Key Pair timeout and got an error: %#v.", err))
		}

		return nil
	})
}
