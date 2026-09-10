package alicloud

import (
	"fmt"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAlicloudPolarDBGlobalDatabaseNetwork() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudPolarDBGlobalDatabaseNetworkCreate,
		Read:   resourceAlicloudPolarDBGlobalDatabaseNetworkRead,
		Update: resourceAlicloudPolarDBGlobalDatabaseNetworkUpdate,
		Delete: resourceAlicloudPolarDBGlobalDatabaseNetworkDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(3 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"db_cluster_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAlicloudPolarDBGlobalDatabaseNetworkCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	polarDBService := PolarDBService{client}
	var response map[string]interface{}
	action := "CreateGlobalDatabaseNetwork"
	request := make(map[string]interface{})
	var err error

	request["DBClusterId"] = d.Get("db_cluster_id")

	if v, ok := d.GetOk("description"); ok {
		request["GDNDescription"] = v
	}

	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RpcPost("polardb", "2017-08-01", action, nil, request, false)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_polardb_global_database_network", action, AlibabaCloudSdkGoERROR)
	}
	d.SetId(fmt.Sprint(response["GDNId"]))

	stateConf := BuildStateConf([]string{"creating"}, []string{"active"}, d.Timeout(schema.TimeoutCreate), 5*time.Minute, polarDBService.PolarDBGlobalDatabaseNetworkRefreshFunc(d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAlicloudPolarDBGlobalDatabaseNetworkRead(d, meta)
}

func resourceAlicloudPolarDBGlobalDatabaseNetworkRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	polarDBService := PolarDBService{client}
	object, err := polarDBService.DescribePolarDBGlobalDatabaseNetwork(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	dBClusterId := object["DBClusters"].([]interface{})[0].(map[string]interface{})["DBClusterId"]
	d.Set("db_cluster_id", dBClusterId)
	d.Set("description", object["GDNDescription"])
	d.Set("status", object["GDNStatus"])

	return nil
}

func resourceAlicloudPolarDBGlobalDatabaseNetworkUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	var err error
	update := false
	request := map[string]interface{}{
		"GDNId": d.Id(),
	}

	if d.HasChange("description") {
		update = true
	}
	if v, ok := d.GetOk("description"); ok {
		request["GDNDescription"] = v
	}

	if update {
		action := "ModifyGlobalDatabaseNetwork"
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("polardb", "2017-08-01", action, nil, request, false)
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
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
	}

	return resourceAlicloudPolarDBGlobalDatabaseNetworkRead(d, meta)
}

func resourceAlicloudPolarDBGlobalDatabaseNetworkDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	polarDBService := PolarDBService{client}
	action := "DeleteGlobalDatabaseNetwork"
	var response map[string]interface{}
	var err error
	request := map[string]interface{}{
		"GDNId": d.Id(),
	}

	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.RpcPost("polardb", "2017-08-01", action, nil, request, false)
		if err != nil {
			if NeedRetry(err) || IsExpectedErrors(err, []string{"MemberNumber.NotSupport"}) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	stateConf := BuildStateConf([]string{"deleting"}, []string{}, d.Timeout(schema.TimeoutDelete), 3*time.Minute, polarDBService.PolarDBGlobalDatabaseNetworkRefreshFunc(d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return nil
}
