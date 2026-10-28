package ackone

import (
	"log"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceMembershipAttachment() *schema.Resource {
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
		Create: resourceMembershipAttachmentCreate,
		Read:   resourceMembershipAttachmentRead,
		Delete: resourceMembershipAttachmentDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
}

func resourceMembershipAttachmentCreate(d *schema.ResourceData, meta interface{}) error {
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

	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		response, err = client.RpcPost("adcp", "2022-01-01", action, nil, request, true)

		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(action, response, request)

	if err != nil {
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_ack_one_membership_attachment", action, errs.AlibabaCloudSdkGoERROR)
	}

	ackOneService := NewService(client)
	stateConf := helper.BuildStateConf([]string{}, []string{"running"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, ackOneService.ClusterStateRefreshFunc(clusterId, "$.ClusterInfo.State", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return errs.WrapErrorf(err, errs.IdMsg, resourceId)
	}

	_, err = ackOneService.DescribeMembershipAttachment(resourceId)
	if err != nil {
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_ack_one_membership_attachment", action, errs.AlibabaCloudSdkGoERROR)
	}

	d.SetId(resourceId)
	return nil
}

func resourceMembershipAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ackOneService := NewService(client)

	objectRaw, err := ackOneService.DescribeMembershipAttachment(d.Id())
	if err != nil {
		if !d.IsNewResource() && errs.NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_ack_one_membership_attachment DescribeMembershipAttachment Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return errs.WrapError(err)
	}

	d.Set("cluster_id", objectRaw["cluster_id"])
	d.Set("sub_cluster_id", objectRaw["sub_cluster_id"])
	return nil
}

func resourceMembershipAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "DetachClusterFromHub"
	var request map[string]interface{}
	var response map[string]interface{}
	var err error
	request = make(map[string]interface{})
	clusterId := d.Get("cluster_id").(string)
	request["ClusterId"] = d.Get("cluster_id")
	request["ClusterIds"] = "[\"" + d.Get("sub_cluster_id").(string) + "\"]"

	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutDelete), func() *retry.RetryError {
		response, err = client.RpcPost("adcp", "2022-01-01", action, nil, request, true)

		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(action, response, request)

	if err != nil {
		if errs.NotFoundError(err) {
			return nil
		}
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
	}

	ackOneService := NewService(client)
	stateConf := helper.BuildStateConf([]string{}, []string{"running"}, d.Timeout(schema.TimeoutDelete), 5*time.Second, ackOneService.ClusterStateRefreshFunc(clusterId, "$.ClusterInfo.State", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return errs.WrapErrorf(err, errs.IdMsg, d.Id())
	}

	return nil
}
