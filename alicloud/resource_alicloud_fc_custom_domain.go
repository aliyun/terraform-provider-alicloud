package alicloud

import (
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	"github.com/aliyun/fc-go-sdk"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAlicloudFCCustomDomain() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudFCCustomDomainCreate,
		Read:   resourceAlicloudFCCustomDomainRead,
		Update: resourceAlicloudFCCustomDomainUpdate,
		Delete: resourceAlicloudFCCustomDomainDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"protocol": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"HTTP", "HTTP,HTTPS"}, false),
			},
			"route_config": {
				Type: schema.TypeSet,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"path": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringLenBetween(1, 1024),
						},
						"service": {
							Type:     schema.TypeString,
							Required: true,
						},
						"function": {
							Type:     schema.TypeString,
							Required: true,
						},
						"qualifier": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
				Optional: true,
				MinItems: 1,
			},
			"create": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"last_modified": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAlicloudFCCustomDomainCreate(d *schema.ResourceData, meta interface{}) error {
	// even after dns record creation completed, fc still can't resolve it, so sleep to wait for sync
	time.Sleep(3 * time.Second)
	client := meta.(*connectivity.AliyunClient)

	domainName := d.Get("name").(string)

	var routeConfig *fc.RouteConfig

	if routes := buildRoutes(d); len(routes) != 0 {
		routeConfig = &fc.RouteConfig{
			Routes: routes,
		}
	}

	request := &fc.CreateCustomDomainInput{
		DomainName:  StringPointer(domainName),
		Protocol:    StringPointer(d.Get("protocol").(string)),
		RouteConfig: routeConfig,
	}
	var response *fc.CreateCustomDomainOutput
	var requestInfo *fc.Client
	if err := resource.Retry(2*time.Minute, func() *resource.RetryError {
		raw, err := client.WithFcClient(func(fcClient *fc.Client) (interface{}, error) {
			requestInfo = fcClient
			return fcClient.CreateCustomDomain(request)
		})
		if err != nil {
			if IsExpectedErrors(err, []string{"AccessDenied"}) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug("CreateCustomDomain", raw, requestInfo, request)
		response, _ = raw.(*fc.CreateCustomDomainOutput)
		return nil

	}); err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_fc_custom_domain", "CreateCustomDomain", FcGoSdk)
	}

	d.SetId(*response.DomainName)

	return resourceAlicloudFCCustomDomainRead(d, meta)
}

func resourceAlicloudFCCustomDomainRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	fcService := FcService{client}

	customDomain, err := fcService.DescribeFcCustomDomain(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("name", customDomain.DomainName)
	d.Set("protocol", customDomain.Protocol)
	d.Set("create", customDomain.CreatedTime)
	d.Set("last_modified", customDomain.LastModifiedTime)

	var s []map[string]interface{}
	if customDomain.RouteConfig != nil {
		for _, pc := range customDomain.RouteConfig.Routes {
			s = append(s, map[string]interface{}{
				"path":      pc.Path,
				"service":   pc.ServiceName,
				"function":  pc.FunctionName,
				"qualifier": pc.Qualifier,
			})
		}
	}
	d.Set("route_config", s)

	return nil
}

func resourceAlicloudFCCustomDomainUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	updated := false
	object := fc.UpdateCustomDomainObject{}

	if d.HasChange("protocol") {
		object.Protocol = StringPointer(d.Get("protocol").(string))
		updated = true
	}
	if d.HasChange("route_config") {
		object.RouteConfig = &fc.RouteConfig{
			Routes: buildRoutes(d),
		}
		updated = true
	}
	updateInput := &fc.UpdateCustomDomainInput{
		DomainName:               StringPointer(d.Id()),
		UpdateCustomDomainObject: object,
	}

	if updated {
		var requestInfo *fc.Client
		raw, err := client.WithFcClient(func(fcClient *fc.Client) (interface{}, error) {
			requestInfo = fcClient
			return fcClient.UpdateCustomDomain(updateInput)
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), "UpdateCustomDomain", FcGoSdk)
		}
		addDebug("UpdateCustomDomain", raw, requestInfo, updateInput)
	}

	return resourceAlicloudFCCustomDomainRead(d, meta)
}

func resourceAlicloudFCCustomDomainDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	fcService := FcService{client}

	request := &fc.DeleteCustomDomainInput{
		DomainName: StringPointer(d.Id()),
	}
	var requestInfo *fc.Client
	raw, err := client.WithFcClient(func(fcClient *fc.Client) (interface{}, error) {
		requestInfo = fcClient
		return fcClient.DeleteCustomDomain(request)
	})
	if err != nil {
		if IsExpectedErrors(err, []string{"CustomDomainNotFound", "DomainNameNotFound"}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), "DeleteCustomDomain", FcGoSdk)
	}
	addDebug("DeleteCustomDomain", raw, requestInfo, request)
	return WrapError(fcService.WaitForFcCustomDomain(d.Id(), Deleted, DefaultTimeoutMedium))
}

func buildRoutes(d *schema.ResourceData) (result []fc.PathConfig) {
	if routes, ok := d.GetOk("route_config"); ok {
		for _, r := range routes.(*schema.Set).List() {
			v := r.(map[string]interface{})
			pathConfig := fc.PathConfig{
				Path:         StringPointer(v["path"].(string)),
				ServiceName:  StringPointer(v["service"].(string)),
				FunctionName: StringPointer(v["function"].(string)),
			}
			if qualifier, ok := v["qualifier"]; ok {
				pathConfig.Qualifier = StringPointer(qualifier.(string))
			}
			result = append(result, pathConfig)
		}
	}
	return
}
