package alicloud

import (
	"fmt"
	"strings"
	"time"

	"github.com/denverdino/aliyungo/cs"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceAlicloudCSApplication() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudCSApplicationCreate,
		Read:   resourceAlicloudCSApplicationRead,
		Update: resourceAlicloudCSApplicationUpdate,
		Delete: resourceAlicloudCSApplicationDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"cluster_name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"name": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateContainerAppName,
			},
			"description": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return !d.HasChange("version")
				},
			},
			"template": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				StateFunc: func(v interface{}) string {
					yaml, _ := normalizeYamlString(v)
					return yaml
				},
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					if !d.HasChange("version") {
						return true
					}
					equal, _ := CompareYmalTemplateAreEquivalent(old, new)
					return equal
				},
				ValidateFunc: validateYamlString,
			},
			"version": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "1.0",
			},
			"environment": &schema.Schema{
				Type:     schema.TypeMap,
				Optional: true,
				Elem:     schema.TypeString,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return !d.HasChange("version")
				},
			},
			"latest_image": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return !d.HasChange("version")
				},
			},
			"blue_green": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return d.IsNewResource()
				},
			},
			"blue_green_confirm": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return !d.Get("blue_green").(bool)
				},
			},
			"services": &schema.Schema{
				Type: schema.TypeList,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"version": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
				Computed: true,
			},
			"default_domain": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAlicloudCSApplicationCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*AliyunClient)
	clusterName := Trim(d.Get("cluster_name").(string))
	conn, err := client.GetApplicationClientByClusterName(clusterName)

	args := &cs.ProjectCreationArgs{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
		Template:    d.Get("template").(string),
		Version:     d.Get("version").(string),
		LatestImage: d.Get("latest_image").(bool),
	}
	if environment, ok := d.GetOk("environment"); ok {
		env := make(map[string]string)
		for k, v := range environment.(map[string]interface{}) {
			env[k] = v.(string)
		}
		args.Environment = env
	}
	invoker := NewInvoker()
	if err := invoker.Run(func() error {
		return conn.CreateProject(args)
	}); err != nil {
		return fmt.Errorf("Creating container application got an error: %#v", err)
	}

	d.SetId(fmt.Sprintf("%s%s%s", clusterName, COLON_SEPARATED, args.Name))

	if err = client.WaitForContainerApplication(clusterName, args.Name, Running, DefaultTimeoutMedium); err != nil {
		return fmt.Errorf("Waitting for container application %#v got an error: %#v", cs.Running, err)
	}

	return resourceAlicloudCSApplicationRead(d, meta)
}

func resourceAlicloudCSApplicationRead(d *schema.ResourceData, meta interface{}) error {
	parts := strings.Split(d.Id(), COLON_SEPARATED)

	application, err := meta.(*AliyunClient).DescribeContainerApplication(parts[0], parts[1])

	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return err
	}

	d.Set("cluster_name", parts[0])
	d.Set("name", application.Name)
	d.Set("description", application.Description)
	d.Set("version", application.Version)
	d.Set("template", application.Template)
	env := make(map[string]string)
	for k, v := range application.Environment {
		if k == "COMPOSE_PROJECT_NAME" || k == "ACS_PROJECT_VERSION" {
			continue
		}
		if k == "ACS_DEFAULT_DOMAIN" {
			d.Set("default_domain", v)
			continue
		}
		env[k] = v
	}
	d.Set("environment", env)
	var services []map[string]interface{}
	for _, s := range application.Services {
		mapping := map[string]interface{}{
			"id":      s.ID,
			"name":    s.Name,
			"status":  s.CurrentState,
			"version": s.Version,
		}
		services = append(services, mapping)
	}
	d.Set("services", services)

	return nil
}

func resourceAlicloudCSApplicationUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*AliyunClient)
	parts := strings.Split(d.Id(), COLON_SEPARATED)
	conn, err := client.GetApplicationClientByClusterName(parts[0])
	if err != nil {
		return err
	}
	invoker := NewInvoker()
	args := &cs.ProjectUpdationArgs{
		Name:        parts[1],
		Description: d.Get("description").(string),
		Template:    d.Get("template").(string),
		Version:     d.Get("version").(string),
		LatestImage: d.Get("latest_image").(bool),
	}

	if environment, ok := d.GetOk("environment"); ok {
		env := make(map[string]string)
		for k, v := range environment.(map[string]interface{}) {
			env[k] = v.(string)
		}
		args.Environment = env
	}

	blue_green := d.Get("blue_green").(bool)
	if blue_green {
		args.UpdateMethod = "blue-green"
	}

	d.Partial(true)
	update := false
	if d.HasChange("description") {
		update = true
		d.SetPartial("description")
	}

	if d.HasChange("template") {
		update = true
		d.SetPartial("template")
	}

	if d.HasChange("environment") {
		update = true
		d.SetPartial("environment")
	}

	if d.HasChange("version") {
		update = true
		d.SetPartial("version")
	}

	if d.HasChange("latest_image") {
		update = true
	}

	if d.HasChange("blue_green") {
		update = true
		d.SetPartial("blue_green")
	}

	if !d.HasChange("version") && !blue_green {
		if err := conn.RollBackBlueGreenProject(parts[1], true); err != nil {
			return fmt.Errorf("Rollbacking container application blue-green got an error: %#v", err)
		}
	} else if update {
		for {
			if err := invoker.Run(func() error {
				return conn.UpdateProject(args)
			}); err != nil {
				if IsExceptedError(err, ApplicationConfirmConflict) {
					if err := invoker.Run(func() error {
						return conn.RollBackBlueGreenProject(parts[1], true)
					}); err != nil {
						return fmt.Errorf("Rollbacking container application blue-green got an error: %#v", err)
					}
					if err := client.WaitForContainerApplication(parts[0], parts[1], Running, DefaultTimeoutMedium); err != nil {
						return fmt.Errorf("After rolling back blue-green project, waitting for container application %#v got an error: %#v", Running, err)
					}
					continue
				}
				return fmt.Errorf("Updating container application got an error: %#v", err)
			} else {
				break
			}
		}
	}

	if d.Get("blue_green_confirm").(bool) {
		if err := invoker.Run(func() error {
			return conn.ConfirmBlueGreenProject(parts[1], true)
		}); err != nil {
			return fmt.Errorf("Confirmming container application blue-green got an error: %#v", err)
		}
	}

	if err := client.WaitForContainerApplication(parts[0], parts[1], Running, DefaultTimeoutMedium); err != nil {
		return fmt.Errorf("After updating, waitting for container application %#v got an error: %#v", Running, err)
	}

	d.Partial(false)

	return resourceAlicloudCSApplicationRead(d, meta)
}

func resourceAlicloudCSApplicationDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*AliyunClient)
	parts := strings.Split(d.Id(), COLON_SEPARATED)
	conn, err := client.GetApplicationClientByClusterName(parts[0])
	if err != nil {
		if NotFoundError(err) {
			return nil
		}
		return err
	}

	appName := parts[1]
	invoker := NewInvoker()

	return resource.Retry(3*time.Minute, func() *resource.RetryError {
		err := invoker.Run(func() error {
			return conn.DeleteProject(appName, true, false)
		})
		if err != nil {
			if IsExceptedError(err, ApplicationNotFound) {
				return nil
			}
			if !IsExceptedError(err, ApplicationErrorIgnore) && !IsExceptedError(err, AliyunGoClientFailure) {
				return resource.NonRetryableError(fmt.Errorf("Deleting container application %s got an error: %#v.", appName, err))
			}
		}

		var project cs.GetProjectResponse
		if err := invoker.Run(func() error {
			resp, e := conn.GetProject(appName)
			if e != nil {
				return e
			}
			project = resp
			return nil
		}); err != nil {
			if IsExceptedError(err, ApplicationNotFound) || IsExceptedError(err, ApplicationErrorIgnore) {
				return nil
			}
			return resource.NonRetryableError(fmt.Errorf("Getting container application %s got an error: %#v.", appName, err))
		}
		if project.Name == "" {
			return nil
		}

		return resource.RetryableError(fmt.Errorf("Deleting container application %s timeout: %#v.", appName, err))
	})
}
