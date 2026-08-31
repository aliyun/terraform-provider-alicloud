package alicloud

import (
	"log"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudAckOneMembershipAttachment() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"cluster_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "ID of the ACK One fleet cluster",
			},
			"sub_cluster_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "ID of the ACK cluster that needs to be managed by ACK One fleet",
			},
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(25 * time.Minute),
			Delete: schema.DefaultTimeout(25 * time.Minute),
		},
		Create: resourceAliCloudAckOneMembershipAttachmentCreate,
		Read:   resourceAliCloudAckOneMembershipAttachmentRead,
		Delete: resourceAliCloudAckOneMembershipAttachmentDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
}

func resourceAliCloudAckOneMembershipAttachmentCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	action := "AttachClusterToHub"
	var request map[string]interface{}
	var response map[string]interface{}
	var err error
	request = make(map[string]interface{})
	clusterId := d.Get("cluster_id").(string)
	subClusterId := d.Get("sub_cluster_id").(string)
	resourceId := clusterId + ":" + subClusterId
	request["ClusterId"] = clusterId
	request["ClusterIds"] = "[\"" + subClusterId + "\"]"

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RpcPost("adcp", "2022-01-01", action, nil, request, true)

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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_ack_one_membership_attachment", action, AlibabaCloudSdkGoERROR)
	}

	ackOneServiceV2 := AckOneServiceV2{client}
	stateConf := BuildStateConf([]string{}, []string{"running"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, ackOneServiceV2.AckOneClusterStateRefreshFunc(clusterId, "$.ClusterInfo.State", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, resourceId)
	}

	_, err = ackOneServiceV2.DescribeAckOneMembershipAttachment(resourceId)
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_ack_one_membership_attachment", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(resourceId)
	return nil
}

func resourceAliCloudAckOneMembershipAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ackOneServiceV2 := AckOneServiceV2{client}

	objectRaw, err := ackOneServiceV2.DescribeAckOneMembershipAttachment(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_ack_one_membership_attachment DescribeAckOneMembershipAttachment Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("cluster_id", objectRaw["cluster_id"])
	d.Set("sub_cluster_id", objectRaw["sub_cluster_id"])
	return nil
}

func resourceAliCloudAckOneMembershipAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "DetachClusterFromHub"
	var request map[string]interface{}
	var response map[string]interface{}
	var err error
	request = make(map[string]interface{})
	clusterId := d.Get("cluster_id").(string)
	request["ClusterId"] = d.Get("cluster_id")
	request["ClusterIds"] = "[\"" + d.Get("sub_cluster_id").(string) + "\"]"

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.RpcPost("adcp", "2022-01-01", action, nil, request, true)

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
		if NotFoundError(err) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	ackOneServiceV2 := AckOneServiceV2{client}
	stateConf := BuildStateConf([]string{}, []string{"running"}, d.Timeout(schema.TimeoutDelete), 5*time.Second, ackOneServiceV2.AckOneClusterStateRefreshFunc(clusterId, "$.ClusterInfo.State", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return nil
}
