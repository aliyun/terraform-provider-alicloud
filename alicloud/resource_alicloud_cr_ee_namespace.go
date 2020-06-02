package alicloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func resourceAlicloudCrEENamespace() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudCrEENamespaceCreate,
		Read:   resourceAlicloudCrEENamespaceRead,
		Update: resourceAlicloudCrEENamespaceUpdate,
		Delete: resourceAlicloudCrEENamespaceDelete,
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
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(2, 30),
			},
			"auto_create": {
				Type:     schema.TypeBool,
				Required: true,
			},
			"default_visibility": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{RepoTypePublic, RepoTypePrivate}, false),
			},
		},
	}
}

func resourceAlicloudCrEENamespaceCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	crService := &CrService{client}
	instanceId := d.Get("instance_id").(string)
	namespace := d.Get("name").(string)
	autoCreate := d.Get("auto_create").(bool)
	visibility := d.Get("default_visibility").(string)
	_, err := crService.CreateCrEENamespace(instanceId, namespace, autoCreate, visibility)
	if err != nil {
		return WrapError(err)
	}

	d.SetId(crService.GenResourceId(instanceId, namespace))

	return resourceAlicloudCrEENamespaceRead(d, meta)
}

func resourceAlicloudCrEENamespaceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	crService := &CrService{client}
	resp, err := crService.DescribeCrEENamespace(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("instance_id", resp.InstanceId)
	d.Set("name", resp.NamespaceName)
	d.Set("auto_create", resp.AutoCreateRepo)
	d.Set("default_visibility", resp.DefaultRepoType)

	return nil
}

func resourceAlicloudCrEENamespaceUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	crService := &CrService{client}
	instanceId := d.Get("instance_id").(string)
	namespace := d.Get("name").(string)
	if d.HasChanges("auto_create", "default_visibility") {
		autoCreate := d.Get("auto_create").(bool)
		visibility := d.Get("default_visibility").(string)
		_, err := crService.UpdateCrEENamespace(instanceId, namespace, autoCreate, visibility)
		if err != nil {
			return WrapError(err)
		}
	}

	return resourceAlicloudCrEENamespaceRead(d, meta)
}

func resourceAlicloudCrEENamespaceDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	crService := &CrService{client}
	instanceId := d.Get("instance_id").(string)
	namespace := d.Get("name").(string)
	_, err := crService.DeleteCrEENamespace(instanceId, namespace)
	if err != nil {
		if NotFoundError(err) {
			return nil
		} else {
			return WrapError(err)
		}
	}

	return WrapError(crService.WaitForCrEENamespace(instanceId, namespace, Deleted, DefaultTimeout))
}
