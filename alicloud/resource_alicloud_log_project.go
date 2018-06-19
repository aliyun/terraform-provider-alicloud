package alicloud

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceAlicloudLogProject() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudLogProjectCreate,
		Read:   resourceAlicloudLogProjectRead,
		//Update: resourceAlicloudLogProjectUpdate,
		Delete: resourceAlicloudLogProjectDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"description": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
		},
	}
}

func resourceAlicloudLogProjectCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*AliyunClient)

	project, err := client.logconn.CreateProject(d.Get("name").(string), d.Get("description").(string))
	if err != nil {
		return fmt.Errorf("CreateProject got an error: %#v.", err)
	}

	d.SetId(project.Name)

	return resourceAlicloudLogProjectRead(d, meta)
}

func resourceAlicloudLogProjectRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*AliyunClient)

	project, err := client.logconn.GetProject(d.Id())
	if err != nil {
		if IsExceptedError(err, ProjectNotExist) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("GetProject got an error: %#v.", err)
	}
	d.Set("name", project.Name)
	d.Set("description", project.Description)

	return nil
}

//func resourceAlicloudLogProjectUpdate(d *schema.ResourceData, meta interface{}) error {
//	client := meta.(*AliyunClient)
//
//	d.Partial(true)
//
//	if d.HasChange("description") {
//		if err := client.logconn.UpdateProject(d.Id(), d.Get("description").(string)); err != nil {
//			return fmt.Errorf("UpdateProject got an error: %#v", err)
//		}
//		d.SetPartial("description")
//	}
//
//	d.Partial(false)
//
//	return resourceAlicloudLogProjectRead(d, meta)
//}

func resourceAlicloudLogProjectDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*AliyunClient).logconn

	return resource.Retry(3*time.Minute, func() *resource.RetryError {
		if err := conn.DeleteProject(d.Id()); err != nil {
			if !IsExceptedErrors(err, []string{ProjectNotExist}) {
				return resource.NonRetryableError(fmt.Errorf("Deleting log project got an error: %#v", err))
			}
		}

		exist, err := conn.CheckProjectExist(d.Id())
		if err != nil {
			return resource.NonRetryableError(fmt.Errorf("While deleting log project, checking project existing got an error: %#v.", err))
		}
		if !exist {
			return nil
		}

		return resource.RetryableError(fmt.Errorf("Deleting log project %s timeout.", d.Id()))
	})
}
