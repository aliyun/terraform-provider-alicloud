package cms

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/fwadapt"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
)

const alertWebhooksTypeName = "alicloud_cms_alert_webhooks"

var (
	_ datasource.DataSource              = &alertWebhooksDataSource{}
	_ datasource.DataSourceWithConfigure = &alertWebhooksDataSource{}
)

func NewAlertWebhooksDataSource() datasource.DataSource {
	return &alertWebhooksDataSource{}
}

type alertWebhooksDataSource struct {
	fwadapt.DataSourceBase
}

type alertWebhooksModel struct {
	Id               types.String `tfsdk:"id"`
	Ids              types.List   `tfsdk:"ids"`
	NameRegex        types.String `tfsdk:"name_regex"`
	Names            types.List   `tfsdk:"names"`
	AlertWebhookName types.String `tfsdk:"alert_webhook_name"`
	Workspace        types.String `tfsdk:"workspace"`
	OutputFile       types.String `tfsdk:"output_file"`
	Webhooks         types.List   `tfsdk:"webhooks"`
}

type alertWebhookItemModel struct {
	Id               types.String `tfsdk:"id"`
	AlertWebhookId   types.String `tfsdk:"alert_webhook_id"`
	AlertWebhookName types.String `tfsdk:"alert_webhook_name"`
	ContentType      types.String `tfsdk:"content_type"`
	Headers          types.Map    `tfsdk:"headers"`
	Lang             types.String `tfsdk:"lang"`
	Method           types.String `tfsdk:"method"`
	Url              types.String `tfsdk:"url"`
	Workspace        types.String `tfsdk:"workspace"`
}

func (d *alertWebhooksDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = alertWebhooksTypeName
}

func (d *alertWebhooksDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "This data source provides a list of CloudMonitor (Cms) alert webhooks in an Alibaba Cloud account.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "The hash of the returned webhook ids.",
			},
			"ids": schema.ListAttribute{
				ElementType: types.StringType,
				Optional:    true,
				Computed:    true,
				Description: "A list of alert webhook ids used to filter the results (mapped to the webhookIds query parameter).",
			},
			"name_regex": schema.StringAttribute{
				Optional:    true,
				Description: "A regex string to filter results by the alert webhook name.",
			},
			"names": schema.ListAttribute{
				ElementType: types.StringType,
				Computed:    true,
				Description: "A list of alert webhook names.",
			},
			"alert_webhook_name": schema.StringAttribute{
				Optional:    true,
				Description: "The name of the alert webhook to filter by.",
			},
			"workspace": schema.StringAttribute{
				Optional:    true,
				Description: "The workspace to filter by.",
			},
			"output_file": schema.StringAttribute{
				Optional:    true,
				Description: "Save the result to a file.",
			},
		},
		Blocks: map[string]schema.Block{
			// webhooks is a nested BLOCK rather than a nested attribute because the
			// muxed provider exposes protocol v5 (tf5muxserver), and protocol v5 can
			// express a list of objects only as a nested block — a ListNestedAttribute
			// would fail schema conversion with "AttributeName webhooks: protocol
			// version 5 cannot have Attributes set". The SDK v2 half expresses the same
			// shape as a ListNestedBlock (see provider/framework/schema.go convertBlock).
			"webhooks": schema.ListNestedBlock{
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Computed:    true,
							Description: "The ID of the alert webhook.",
						},
						"alert_webhook_id": schema.StringAttribute{
							Computed:    true,
							Description: "The ID of the alert webhook.",
						},
						"alert_webhook_name": schema.StringAttribute{
							Computed:    true,
							Description: "The name of the alert webhook.",
						},
						"content_type": schema.StringAttribute{
							Computed:    true,
							Description: "The content type of the alert webhook request body.",
						},
						"headers": schema.MapAttribute{
							ElementType: types.StringType,
							Computed:    true,
							Description: "The custom request headers of the alert webhook.",
						},
						"lang": schema.StringAttribute{
							Computed:    true,
							Description: "The language of the alert webhook.",
						},
						"method": schema.StringAttribute{
							Computed:    true,
							Description: "The HTTP method of the alert webhook request.",
						},
						"url": schema.StringAttribute{
							Computed:    true,
							Description: "The URL of the alert webhook.",
						},
						"workspace": schema.StringAttribute{
							Computed:    true,
							Description: "The workspace to which the alert webhook belongs.",
						},
					},
				},
				Description: "A list of alert webhooks.",
			},
		},
	}
}

// webhookItemObjectType mirrors the webhooks nested attribute object so a
// types.List built from alertWebhookItemModel values matches the schema.
func webhookItemObjectType() attr.Type {
	return types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"id":                 types.StringType,
			"alert_webhook_id":   types.StringType,
			"alert_webhook_name": types.StringType,
			"content_type":       types.StringType,
			"headers":            types.MapType{ElemType: types.StringType},
			"lang":               types.StringType,
			"method":             types.StringType,
			"url":                types.StringType,
			"workspace":          types.StringType,
		},
	}
}

func (d *alertWebhooksDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var plan alertWebhooksModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	client := d.Client()

	var nameRegex *regexp.Regexp
	if !plan.NameRegex.IsNull() {
		r, err := regexp.Compile(plan.NameRegex.ValueString())
		if err != nil {
			resp.Diagnostics.AddError("Compiling name_regex", err.Error())
			return
		}
		nameRegex = r
	}

	idsSet := map[string]struct{}{}
	var ids []string
	if !plan.Ids.IsNull() {
		resp.Diagnostics.Append(plan.Ids.ElementsAs(ctx, &ids, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
		for _, id := range ids {
			idsSet[id] = struct{}{}
		}
	}

	action := "/webhooks"

	// ListAlertWebhooks exposes webhookIds as a server-side filter that scopes
	// the response to the requested ids (a direct lookup, the same path the
	// resource Read uses after Create). When the caller supplies ids, we send
	// webhookIds so the data source finds a freshly-created webhook even when
	// the bare list view for the account has not indexed it yet — the bare
	// GET /webhooks index is observed to omit a webhook created without a
	// workspace, while webhookIds returns it immediately. name and workspace
	// are likewise server-side filters; name_regex remains client-side over
	// whatever the server returns. A filtered query that still comes back empty
	// is retried briefly for propagation before returning an empty result; an
	// unfiltered query is returned as-is (including empty), so querying an
	// account with no webhooks does not block.
	filtered := nameRegex != nil || len(idsSet) > 0 ||
		(!plan.AlertWebhookName.IsNull() && plan.AlertWebhookName.ValueString() != "") ||
		(!plan.Workspace.IsNull() && plan.Workspace.ValueString() != "")

	var objects []map[string]interface{}
	var propagationExhausted bool
	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err := retry.Retry(2*time.Minute, func() *retry.RetryError {
		objects = nil
		propagationExhausted = false
		query := map[string]*string{
			"pageNumber": ptr("1"),
			"pageSize":   ptr(fmt.Sprintf("%d", pageSizeLarge)),
		}
		if len(ids) > 0 {
			webhookIdsJson, jerr := json.Marshal(ids)
			if jerr != nil {
				return retry.NonRetryableError(jerr)
			}
			query["webhookIds"] = ptr(string(webhookIdsJson))
		}
		if !plan.AlertWebhookName.IsNull() && plan.AlertWebhookName.ValueString() != "" {
			query["name"] = ptr(plan.AlertWebhookName.ValueString())
		}
		if !plan.Workspace.IsNull() && plan.Workspace.ValueString() != "" {
			query["workspace"] = ptr(plan.Workspace.ValueString())
		}
		for {
			response, e := client.RoaGet(cmsProduct, cmsAPIVersion, action, query, nil, nil)
			if e != nil {
				if helper.NeedRetry(e) {
					wait()
					return retry.RetryableError(e)
				}
				return retry.NonRetryableError(e)
			}
			logs.AddDebug(action, response, query)

			raw, gerr := jsonpath.Get("$.webhooks[*]", response)
			if gerr != nil {
				if response["webhooks"] == nil {
					break
				}
				return retry.NonRetryableError(errs.WrapErrorf(gerr, errs.FailedGetAttributeMsg, "$.webhooks[*]", "$.webhooks[*]", response))
			}
			items, _ := raw.([]interface{})

			totalCount := 0
			if t, _ := jsonpath.Get("$.total", response); t != nil {
				if n, perr := strconv.Atoi(fmt.Sprint(t)); perr == nil {
					totalCount = n
				}
			}
			currentPage, _ := strconv.Atoi(*query["pageNumber"])

			for _, v := range items {
				item, ok := v.(map[string]interface{})
				if !ok {
					continue
				}
				webhookName := fmt.Sprint(item["name"])
				if nameRegex != nil && !nameRegex.MatchString(webhookName) {
					continue
				}
				if len(idsSet) > 0 {
					if _, ok := idsSet[fmt.Sprint(item["webhookId"])]; !ok {
						continue
					}
				}
				objects = append(objects, item)
			}

			if len(items) < pageSizeLarge || currentPage*pageSizeLarge >= totalCount {
				break
			}
			query["pageNumber"] = ptr(fmt.Sprintf("%d", currentPage+1))
		}
		if filtered && len(objects) == 0 {
			propagationExhausted = true
			wait()
			return retry.RetryableError(fmt.Errorf("waiting for alert webhooks matching the filter to propagate into the ListAlertWebhooks index"))
		}
		return nil
	})
	if err != nil && !propagationExhausted {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s: calling ListAlertWebhooks", alertWebhooksTypeName),
			errs.WrapErrorf(err, errs.DataDefaultErrorMsg, alertWebhooksTypeName, action, errs.AlibabaCloudSdkGoERROR).Error(),
		)
		return
	}

	ids = make([]string, 0, len(objects))
	names := make([]string, 0, len(objects))
	items := make([]alertWebhookItemModel, 0, len(objects))
	for _, object := range objects {
		webhookId := fmt.Sprint(object["webhookId"])
		webhookName := fmt.Sprint(object["name"])
		ids = append(ids, webhookId)
		names = append(names, webhookName)
		items = append(items, alertWebhookItemModel{
			Id:               types.StringValue(webhookId),
			AlertWebhookId:   types.StringValue(webhookId),
			AlertWebhookName: types.StringValue(webhookName),
			ContentType:      stringValue(object["contentType"]),
			Headers:          headersToState(ctx, object["headers"]),
			Lang:             stringValue(object["lang"]),
			Method:           stringValue(object["method"]),
			Url:              stringValue(object["url"]),
			Workspace:        stringValue(object["workspace"]),
		})
	}

	idList, diags := types.ListValueFrom(ctx, types.StringType, ids)
	resp.Diagnostics.Append(diags...)
	nameList, diags := types.ListValueFrom(ctx, types.StringType, names)
	resp.Diagnostics.Append(diags...)
	webhookList, diags := types.ListValueFrom(ctx, webhookItemObjectType(), items)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	state := alertWebhooksModel{
		Id:               types.StringValue(dataSourceIDsHash(ids)),
		Ids:              idList,
		NameRegex:        plan.NameRegex,
		Names:            nameList,
		AlertWebhookName: plan.AlertWebhookName,
		Workspace:        plan.Workspace,
		OutputFile:       plan.OutputFile,
		Webhooks:         webhookList,
	}

	if !plan.OutputFile.IsNull() && plan.OutputFile.ValueString() != "" {
		if err := writeToFile(plan.OutputFile.ValueString(), objects); err != nil {
			resp.Diagnostics.AddError("Writing output_file", err.Error())
			return
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// stringValue returns a null String for a missing/empty API value, matching
// the inherited PR Read behavior.
func stringValue(v interface{}) types.String {
	if v == nil {
		return types.StringNull()
	}
	if s, ok := v.(string); ok && s != "" {
		return types.StringValue(s)
	}
	s := fmt.Sprint(v)
	if s == "" || s == "<nil>" {
		return types.StringNull()
	}
	return types.StringValue(s)
}

// dataSourceIDsHash mirrors alicloud.dataResourceIdHash without importing the
// flat layout's root package.
func dataSourceIDsHash(ids []string) string {
	b, _ := json.Marshal(ids)
	return fmt.Sprintf("%x", b)
}

func writeToFile(filePath string, data interface{}) error {
	b, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(filePath, b, 0o644)
}
