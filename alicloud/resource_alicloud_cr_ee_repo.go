package alicloud

import (
	"fmt"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudCrEERepo() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudCrEERepoCreate,
		Read:   resourceAliCloudCrEERepoRead,
		Update: resourceAliCloudCrEERepoUpdate,
		Delete: resourceAliCloudCrEERepoDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"namespace": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"repo_type": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: StringInSlice([]string{RepoTypePublic, RepoTypePrivate}, false),
			},
			"summary": {
				Type:     schema.TypeString,
				Required: true,
			},
			"detail": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"tag_immutability": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"repo_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAliCloudCrEERepoCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	action := "CreateRepository"
	query := make(map[string]interface{})
	request := make(map[string]interface{})

	request["RegionId"] = client.RegionId
	request["InstanceId"] = d.Get("instance_id")
	request["RepoNamespaceName"] = d.Get("namespace")
	request["RepoName"] = d.Get("name")
	request["RepoType"] = d.Get("repo_type")
	request["Summary"] = d.Get("summary")

	if v, ok := d.GetOk("detail"); ok {
		request["Detail"] = v
	}
	if v, ok := d.GetOkExists("tag_immutability"); ok {
		request["TagImmutability"] = v
	}

	var response map[string]interface{}
	var err error
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutCreate)), func() *resource.RetryError {
		response, err = client.RpcPost("cr", "2018-12-01", action, query, request, false)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_cr_ee_repo", action, AlibabaCloudSdkGoERROR)
	}
	if isSuccess, ok := response["IsSuccess"].(bool); ok && !isSuccess {
		return WrapErrorf(fmt.Errorf("%v", response), DefaultErrorMsg, "alicloud_cr_ee_repo", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprintf("%s:%s:%s", request["InstanceId"], request["RepoNamespaceName"], request["RepoName"]))

	return resourceAliCloudCrEERepoRead(d, meta)
}

func resourceAliCloudCrEERepoRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	crServiceV2 := CrServiceV2{client}

	object, err := crServiceV2.DescribeCrEERepo(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("instance_id", object["InstanceId"])
	d.Set("namespace", object["RepoNamespaceName"])
	d.Set("name", object["RepoName"])
	d.Set("repo_type", object["RepoType"])
	d.Set("summary", object["Summary"])
	d.Set("detail", object["Detail"])
	d.Set("tag_immutability", object["TagImmutability"])
	d.Set("repo_id", object["RepoId"])

	return nil
}

func resourceAliCloudCrEERepoUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	action := "UpdateRepository"
	query := make(map[string]interface{})
	request := make(map[string]interface{})

	parts, err := ParseResourceId(d.Id(), 3)
	if err != nil {
		return WrapError(err)
	}

	request["RegionId"] = client.RegionId
	request["InstanceId"] = parts[0]
	request["RepoId"] = d.Get("repo_id")
	request["RepoType"] = d.Get("repo_type")
	request["Summary"] = d.Get("summary")

	update := false
	if d.HasChange("repo_type") {
		update = true
	}
	if d.HasChange("summary") {
		update = true
	}
	if d.HasChange("detail") {
		update = true
	}
	if v, ok := d.GetOk("detail"); ok {
		request["Detail"] = v
	}
	if d.HasChange("tag_immutability") {
		update = true
	}
	if v, ok := d.GetOkExists("tag_immutability"); ok {
		request["TagImmutability"] = v
	}

	if update {
		var response map[string]interface{}
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *resource.RetryError {
			response, err = client.RpcPost("cr", "2018-12-01", action, query, request, false)
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
		if isSuccess, ok := response["IsSuccess"].(bool); ok && !isSuccess {
			return WrapErrorf(fmt.Errorf("%v", response), DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
	}

	return resourceAliCloudCrEERepoRead(d, meta)
}

func resourceAliCloudCrEERepoDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	crService := &CrService{client}

	repoId := d.Get("repo_id").(string)
	_, err := crService.DeleteCrEERepo(d.Id(), repoId)
	if err != nil {
		if NotFoundError(err) {
			return nil
		} else {
			return WrapError(err)
		}
	}

	return WrapError(crService.WaitForCrEERepo(d.Id(), Deleted, DefaultTimeout))
}
