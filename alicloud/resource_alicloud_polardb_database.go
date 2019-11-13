package alicloud

import (
	"fmt"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/polardb"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func resourceAlicloudPolarDBDatabase() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudPolarDBDatabaseCreate,
		Read:   resourceAlicloudPolarDBDatabaseRead,
		Update: resourceAlicloudPolarDBDatabaseUpdate,
		Delete: resourceAlicloudPolarDBDatabaseDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"cluster_id": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},

			"name": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},

			"character_set": {
				Type:         schema.TypeString,
				ValidateFunc: validation.StringInSlice(POLAR_DB_CHARACTER_SET_NAME, false),
				Optional:     true,
				Default:      "utf8",
				ForceNew:     true,
			},

			"account_name": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"account_privilege": {
				Type:         schema.TypeString,
				ValidateFunc: validation.StringInSlice(POLAR_DB_ACCOUNT_PRIVILEGE_NAME, false),
				Default:      "ReadWrite",
				Optional:     true,
			},

			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceAlicloudPolarDBDatabaseCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	polarDBService := PolarDBService{client}
	request := polardb.CreateCreateDatabaseRequest()
	request.RegionId = client.RegionId
	request.DBClusterId = d.Get("cluster_id").(string)
	request.DBName = d.Get("name").(string)
	request.CharacterSetName = d.Get("character_set").(string)
	request.AccountName = d.Get("account_name").(string)
	request.AccountPrivilege = d.Get("account_privilege").(string)

	if v, ok := d.GetOk("description"); ok && v.(string) != "" {
		request.DBDescription = v.(string)
	}

	if inst, err := polarDBService.DescribePolarDBCluster(request.DBClusterId); err != nil {
		return WrapError(err)
	} else if inst.Engine == string(PostgreSQL) {
		return WrapError(Error("At present, it does not support creating 'PostgreSQL' database. Please login DB instance to create."))
	}

	err := resource.Retry(5*time.Minute, func() *resource.RetryError {
		raw, err := client.WithPolarDBClient(func(polarDBClient *polardb.Client) (interface{}, error) {
			return polarDBClient.CreateDatabase(request)
		})
		if err != nil {
			if IsExceptedErrors(err, OperationDeniedDBStatus) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		return nil
	})

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprintf("%s%s%s", request.DBClusterId, COLON_SEPARATED, request.DBName))

	return resourceAlicloudPolarDBDatabaseRead(d, meta)
}

func resourceAlicloudPolarDBDatabaseRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	polarDBService := PolarDBService{client}
	object, err := polarDBService.DescribePolarDBDatabase(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return nil
	}

	d.Set("cluster_id", parts[0])
	d.Set("name", object.DBName)
	d.Set("character_set", object.CharacterSetName)
	d.Set("description", object.DBDescription)

	return nil
}

func resourceAlicloudPolarDBDatabaseUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	if d.HasChange("description") && !d.IsNewResource() {
		parts, err := ParseResourceId(d.Id(), 2)
		if err != nil {
			return WrapError(err)
		}
		request := polardb.CreateModifyDBDescriptionRequest()
		request.RegionId = client.RegionId
		request.DBClusterId = parts[0]
		request.DBName = parts[1]
		request.DBDescription = d.Get("description").(string)
		var raw interface{}
		raw, err = client.WithPolarDBClient(func(polarDBClient *polardb.Client) (interface{}, error) {
			return polarDBClient.ModifyDBDescription(request)
		})
		if err != nil {
			return WrapError(err)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	}
	return resourceAlicloudPolarDBDatabaseRead(d, meta)
}

func resourceAlicloudPolarDBDatabaseDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	polarDBService := PolarDBService{client}
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	request := polardb.CreateDeleteDatabaseRequest()
	request.RegionId = client.RegionId
	request.DBClusterId = parts[0]
	request.DBName = parts[1]
	// wait instance status is running before deleting database
	if err := polarDBService.WaitForPolarDBInstance(parts[0], Running, 1800); err != nil {
		return WrapError(err)
	}
	raw, err := client.WithPolarDBClient(func(polarDBClient *polardb.Client) (interface{}, error) {
		return polarDBClient.DeleteDatabase(request)
	})
	if err != nil {
		if IsExceptedErrors(err, []string{InvalidDBClusterIdNotFound, InvalidDBClusterNameNotFound}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	return WrapError(polarDBService.WaitForPolarDBDatabase(d.Id(), Deleted, DefaultTimeoutMedium))
}
