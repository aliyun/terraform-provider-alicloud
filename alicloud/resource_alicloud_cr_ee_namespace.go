package alicloud

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"strconv"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/cr_ee"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudCrEENamespace() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudCrEENamespaceCreate,
		Read:   resourceAliCloudCrEENamespaceRead,
		Update: resourceAliCloudCrEENamespaceUpdate,
		Delete: resourceAliCloudCrEENamespaceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"auto_create": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"default_visibility": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: StringInSlice([]string{RepoTypePublic, RepoTypePrivate}, false),
			},
		},
	}
}

func resourceAliCloudCrEENamespaceCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	crService := &CrService{client}
	request := cr_ee.CreateCreateNamespaceRequest()
	request.RegionId = crService.client.RegionId
	request.InstanceId = d.Get("instance_id").(string)
	request.NamespaceName = d.Get("name").(string)

	if v, ok := d.GetOkExists("auto_create"); ok {
		request.AutoCreateRepo = requests.Boolean(strconv.FormatBool(v.(bool)))
	}

	if v, ok := d.GetOk("default_visibility"); ok {
		request.DefaultRepoType = v.(string)
	}

	var raw interface{}
	var err error
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutCreate)), func() *resource.RetryError {
		raw, err = crService.client.WithCrEEClient(func(creeClient *cr_ee.Client) (interface{}, error) {
			return creeClient.CreateNamespace(request)
		})
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_cr_ee_namespace", request.GetActionName(), AlibabaCloudSdkGoERROR)
	}

	response, _ := raw.(*cr_ee.CreateNamespaceResponse)
	if !response.CreateNamespaceIsSuccess {
		return WrapErrorf(fmt.Errorf("%v", response), DefaultErrorMsg, "alicloud_cr_ee_namespace", request.GetActionName(), AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprintf("%s:%s", request.InstanceId, request.NamespaceName))

	return resourceAliCloudCrEENamespaceRead(d, meta)
}

func resourceAliCloudCrEENamespaceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	crService := &CrService{client}

	object, err := crService.DescribeCrEENamespace(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("instance_id", object.InstanceId)
	d.Set("name", object.NamespaceName)
	d.Set("auto_create", object.AutoCreateRepo)
	d.Set("default_visibility", object.DefaultRepoType)

	return nil
}

func resourceAliCloudCrEENamespaceUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	crService := &CrService{client}
	response := &cr_ee.UpdateNamespaceResponse{}

	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}

	update := false
	request := cr_ee.CreateUpdateNamespaceRequest()
	request.RegionId = crService.client.RegionId
	request.InstanceId = parts[0]
	request.NamespaceName = parts[1]

	if d.HasChange("auto_create") {
		update = true
	}
	if v, ok := d.GetOkExists("auto_create"); ok {
		request.AutoCreateRepo = requests.Boolean(strconv.FormatBool(v.(bool)))
	}

	if d.HasChange("default_visibility") {
		update = true
	}
	if v, ok := d.GetOk("default_visibility"); ok {
		request.DefaultRepoType = v.(string)
	}

	if update {
		var raw interface{}
		var err error
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *resource.RetryError {
			raw, err = crService.client.WithCrEEClient(func(creeClient *cr_ee.Client) (interface{}, error) {
				return creeClient.UpdateNamespace(request)
			})
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)

		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
		}

		response, _ = raw.(*cr_ee.UpdateNamespaceResponse)
		if !response.UpdateNamespaceIsSuccess {
			return WrapErrorf(fmt.Errorf("%v", response), DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
	}

	return resourceAliCloudCrEENamespaceRead(d, meta)
}

func resourceAliCloudCrEENamespaceDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	crService := &CrService{client}

	_, err := crService.DeleteCrEENamespace(d.Id())
	if err != nil {
		if NotFoundError(err) {
			return nil
		} else {
			return WrapError(err)
		}
	}

	return WrapError(crService.WaitForCrEENamespace(d.Id(), Deleted, DefaultTimeout))
}
