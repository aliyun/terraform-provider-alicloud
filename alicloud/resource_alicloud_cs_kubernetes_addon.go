package alicloud

import (
	"encoding/json"
	"fmt"
	"reflect"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

const ResourceAlicloudCSKubernetesAddon = "resourceAlicloudCSKubernetesAddon"

func resourceAlicloudCSKubernetesAddon() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudCSKubernetesAddonCreate,
		Read:   resourceAlicloudCSKubernetesAddonRead,
		Update: resourceAlicloudCSKubernetesAddonUpdate,
		Delete: resourceAlicloudCSKubernetesAddonDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"cluster_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"version": {
				Type:     schema.TypeString,
				Required: true,
			},
			"next_version": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"can_upgrade": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"required": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"config": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func resourceAlicloudCSKubernetesAddonRead(d *schema.ResourceData, meta interface{}) error {
	client, err := meta.(*connectivity.AliyunClient).NewRoaCsClient()
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, ResourceAlicloudCSKubernetesAddon, "InitializeClient", err)
	}
	csClient := CsClient{client}

	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}

	addon, err := csClient.DescribeCsKubernetesAddon(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, ResourceAlicloudCSKubernetesAddon, "DescribeCsKubernetesAddon", err)
	}

	d.Set("cluster_id", parts[0])
	d.Set("name", addon.ComponentName)
	d.Set("version", addon.Version)
	d.Set("next_version", addon.NextVersion)
	d.Set("can_upgrade", addon.CanUpgrade)
	d.Set("required", addon.Required)
	if addon.Config != "" {
		d.Set("config", addon.Config)
	}

	return nil
}

func resourceAlicloudCSKubernetesAddonCreate(d *schema.ResourceData, meta interface{}) error {
	clusterId := d.Get("cluster_id").(string)
	name := d.Get("name").(string)

	client, err := meta.(*connectivity.AliyunClient).NewRoaCsClient()
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, ResourceAlicloudCSKubernetesAddon, "InitializeClient", err)
	}
	csClient := CsClient{client}

	status, err := csClient.DescribeCsKubernetesAddonStatus(clusterId, name)
	if err != nil && !NotFoundError(err) {
		return WrapErrorf(err, DefaultErrorMsg, ResourceAlicloudCSKubernetesAddon, "DescribeCsKubernetesAddonStatus", err)
	}

	changed := false
	if status.Version == "" {
		// Addon not installed
		changed = true
		err := csClient.installAddon(d)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, ResourceAlicloudCSKubernetesAddon, "installAddon", err)
		}
	} else {
		// Addon has been installed
		ok, err := needUpgrade(d, status)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, ResourceAlicloudCSKubernetesAddon, "upgradeAddon", err)
		}
		if ok {
			changed = true
			err := csClient.upgradeAddon(d)
			if err != nil {
				return WrapErrorf(err, DefaultErrorMsg, ResourceAlicloudCSKubernetesAddon, "upgradeAddon", err)
			}
		}
	}

	d.SetId(fmt.Sprintf("%s%s%s", clusterId, COLON_SEPARATED, name))

	stateConf := BuildStateConf([]string{"running", "Running", "Upgrading", "Pause"}, []string{"Success", "NoUpgrade"}, d.Timeout(schema.TimeoutCreate), 10*time.Second, csClient.CsKubernetesAddonStateRefreshFunc(clusterId, name, []string{"Failed", "Canceled"}))
	if changed {
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, DefaultErrorMsg, ResourceAlicloudCSKubernetesAddon, "WaitForSuccessAfterCreate", d.Id())
		}
	}
	return resourceAlicloudCSKubernetesAddonRead(d, meta)
}

func resourceAlicloudCSKubernetesAddonUpdate(d *schema.ResourceData, meta interface{}) error {
	client, err := meta.(*connectivity.AliyunClient).NewRoaCsClient()
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, ResourceAlicloudCSKubernetesAddon, "InitializeClient", err)
	}
	csClient := CsClient{client}

	clusterId := d.Get("cluster_id").(string)
	name := d.Get("name").(string)

	if d.HasChange("version") {
		err := csClient.upgradeAddon(d)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, ResourceAlicloudCSKubernetesAddon, "upgradeAddon", err)
		}
	} else if d.HasChange("config") {
		err := csClient.updateAddonConfig(d)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, ResourceAlicloudCSKubernetesAddon, "upgradeAddonConfig", err)
		}
	}
	stateConf := BuildStateConf([]string{"running", "Running", "Upgrading", "Pause"}, []string{"Success", "NoUpgrade"}, d.Timeout(schema.TimeoutCreate), 10*time.Second, csClient.CsKubernetesAddonStateRefreshFunc(clusterId, name, []string{"Failed", "Canceled"}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, DefaultErrorMsg, ResourceAlicloudCSKubernetesAddon, "WaitForSuccessAfterUpdate", d.Id())
	}

	return resourceAlicloudCSKubernetesAddonRead(d, meta)
}

func resourceAlicloudCSKubernetesAddonDelete(d *schema.ResourceData, meta interface{}) error {
	clusterId := d.Get("cluster_id").(string)
	name := d.Get("name").(string)

	client, err := meta.(*connectivity.AliyunClient).NewRoaCsClient()
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, ResourceAlicloudCSKubernetesAddon, "InitializeClient", err)
	}
	csClient := CsClient{client}

	err = csClient.uninstallAddon(d)
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, ResourceAlicloudCSKubernetesAddon, "uninstallAddon", err)
	}

	stateConf := BuildStateConf([]string{"Running"}, []string{"Deleted"}, d.Timeout(schema.TimeoutDelete), 10*time.Second, csClient.CsKubernetesAddonExistRefreshFunc(clusterId, name))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, DefaultErrorMsg, ResourceAlicloudCSKubernetesAddon, "WaitForSuccessAfterDelete", d.Id())
	}

	return nil
}

func needUpgrade(d *schema.ResourceData, c *Component) (bool, error) {
	// Is version changed
	version := d.Get("version").(string)
	if c.Version != "" && c.Version != version {
		return true, nil
	}

	// Is config changed
	if _, ok := d.GetOk("config"); ok {
		localConfig := map[string]interface{}{}
		err := json.Unmarshal([]byte(d.Get("config").(string)), &localConfig)
		if err != nil {
			return false, WrapErrorf(err, DefaultErrorMsg, ResourceAlicloudCSKubernetesAddon, "json config parse error", err)
		}

		if c.Config == "" {
			c.Config = "{}"
		}
		remoteConfig := map[string]interface{}{}
		err = json.Unmarshal([]byte(c.Config), &remoteConfig)
		if err != nil {
			return false, WrapErrorf(err, DefaultErrorMsg, ResourceAlicloudCSKubernetesAddon, "json config parse error", err)
		}
		return !reflect.DeepEqual(localConfig, remoteConfig), nil
	}
	return false, nil
}
