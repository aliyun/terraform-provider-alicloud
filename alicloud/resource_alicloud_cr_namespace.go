package alicloud

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
	"time"
)

func resourceAlicloudCRNamespace() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudCRNamespaceCreate,
		Read:   resourceAlicloudCRNamespaceRead,
		Update: resourceAlicloudCRNamespaceUpdate,
		Delete: resourceAlicloudCRNamespaceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateContainerRegistryNamespaceName,
			},
			"auto_create": {
				Type:     schema.TypeBool,
				Required: true,
			},
			"default_visibility": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateAllowedStringValue([]string{"PUBLIC", "PRIVATE"}),
			},
		},
	}
}

func resourceAlicloudCRNamespaceCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	crService := CrService{client}

	namespaceName := d.Get("name").(string)

	_, err := crService.CreateNamespace(namespaceName)
	if err != nil {
		return fmt.Errorf("creating namespace got an error: %#v", err)
	}
	d.SetId(namespaceName)

	_, err = crService.UpdateNamespace(namespaceName, d.Get("auto_create").(bool), d.Get("default_visibility").(string))
	if err != nil {
		return fmt.Errorf("updating namespace got an error: %#v", err)
	}

	return resourceAlicloudCRNamespaceUpdate(d, meta)
}

func resourceAlicloudCRNamespaceUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	crService := CrService{client}

	d.Partial(true)
	namespaceName := d.Get("name").(string)

	if d.HasChange("auto_create") || d.HasChange("default_visibility") {
		_, err := crService.UpdateNamespace(namespaceName, d.Get("auto_create").(bool), d.Get("default_visibility").(string))
		if err != nil {
			return fmt.Errorf("updating namespace got an error: %#v", err)
		}
		d.SetPartial("auto_create")
		d.SetPartial("default_visibility")
	}
	d.Partial(false)

	return resourceAlicloudCRNamespaceRead(d, meta)
}

func resourceAlicloudCRNamespaceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	crService := CrService{client}

	resp, err := crService.GetNamespace(d.Id())
	if err != nil {
		return fmt.Errorf("getting namespace got an error: %#v", err)
	}

	d.Set("name", resp.Data.Namespace.Namespace)
	d.Set("auto_create", resp.Data.Namespace.AutoCreate)
	d.Set("default_visibility", resp.Data.Namespace.DefaultVisibility)

	return nil
}

func resourceAlicloudCRNamespaceDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	crService := CrService{client}

	namespaceName := d.Get("name").(string)

	return resource.Retry(5*time.Minute, func() *resource.RetryError {
		_, err := crService.DeleteNamespace(namespaceName)
		if err != nil {
			if NotFoundError(err) || IsExceptedError(err, ErrorNamespaceNotExist) {
				return nil
			}
			return resource.RetryableError(err)
		}
		return nil
	})
}
