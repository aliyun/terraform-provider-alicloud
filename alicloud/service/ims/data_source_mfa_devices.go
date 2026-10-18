package ims

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"regexp"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/fwadapt"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

const mfaDevicesTypeName = "alicloud_ram_mfa_devices"

var (
	_ datasource.DataSource              = &mfaDevicesDataSource{}
	_ datasource.DataSourceWithConfigure = &mfaDevicesDataSource{}
)

func NewMfaDevicesDataSource() datasource.DataSource {
	return &mfaDevicesDataSource{}
}

type mfaDevicesDataSource struct {
	fwadapt.DataSourceBase
}

type mfaDevicesDataSourceModel struct {
	Id         types.String               `tfsdk:"id"`
	NameRegex  types.String               `tfsdk:"name_regex"`
	Ids        types.List                 `tfsdk:"ids"`
	OutputFile types.String               `tfsdk:"output_file"`
	MfaDevices []mfaDeviceDataSourceModel `tfsdk:"mfa_devices"`
}

type mfaDeviceDataSourceModel struct {
	SerialNumber         types.String `tfsdk:"serial_number"`
	VirtualMFADeviceName types.String `tfsdk:"virtual_mfa_device_name"`
	ActivateDate         types.String `tfsdk:"activate_date"`
}

func (d *mfaDevicesDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = mfaDevicesTypeName
}

func (d *mfaDevicesDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Provides a list of RAM virtual MFA Devices available to the user.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "The ID of the data source, a hash of the listed MFA device serial numbers.",
			},
			"name_regex": schema.StringAttribute{
				Optional:    true,
				Description: "A regex string to filter results by the MFA device name. The device name is extracted from the serial number.",
			},
			"ids": schema.ListAttribute{
				ElementType: types.StringType,
				Optional:    true,
				Computed:    true,
				Description: "A list of MFA device serial numbers. If specified, only devices with matching serial numbers are returned.",
			},
			"output_file": schema.StringAttribute{
				Optional:    true,
				Description: "File name where to write the data source results in JSON format.",
			},
		},
		// mfa_devices is a computed list of objects. The provider muxes the
		// Framework half onto protocol v5 via tf5muxserver, and protocol v5
		// cannot represent nested attributes (NestedAttributeObject.Attributes
		// is unsupported). A nested block is the protocol-v5 representation of a
		// list of objects, matching the SDK v2 TypeList+*schema.Resource form
		// used by the original resource set, so mfa_devices is declared as a
		// ListNestedBlock rather than a ListNestedAttribute.
		Blocks: map[string]schema.Block{
			"mfa_devices": schema.ListNestedBlock{
				Description: "A list of MFA Devices. Each element contains the following attributes.",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"serial_number": schema.StringAttribute{
							Computed:    true,
							Description: "The serial number of the MFA device.",
						},
						"virtual_mfa_device_name": schema.StringAttribute{
							Computed:    true,
							Description: "The name of the MFA device, extracted from the serial number.",
						},
						"activate_date": schema.StringAttribute{
							Computed:    true,
							Description: "The activation time of the MFA device. Empty if the device has not been activated.",
						},
					},
				},
			},
		},
	}
}

func (d *mfaDevicesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config mfaDevicesDataSourceModel
	resp.Diagnostics.Append(resp.State.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	devices, err := listVirtualMFADevices(d.Client())
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s: calling ListVirtualMFADevices", mfaDevicesTypeName),
			err.Error(),
		)
		return
	}

	var nameRegex *regexp.Regexp
	if !config.NameRegex.IsNull() && config.NameRegex.ValueString() != "" {
		nameRegex, err = regexp.Compile(config.NameRegex.ValueString())
		if err != nil {
			resp.Diagnostics.AddError(
				fmt.Sprintf("Reading %s: compiling name_regex", mfaDevicesTypeName),
				err.Error(),
			)
			return
		}
	}

	idsFilter := map[string]bool{}
	for _, id := range config.Ids.Elements() {
		idsFilter[id.(types.String).ValueString()] = true
	}

	var ids []string
	output := []mfaDeviceDataSourceModel{}
	for _, device := range devices {
		serialNumber := fmt.Sprint(device["SerialNumber"])
		name := mfaDeviceName(serialNumber)

		if nameRegex != nil && !nameRegex.MatchString(name) {
			continue
		}
		if len(idsFilter) > 0 && !idsFilter[serialNumber] {
			continue
		}

		ids = append(ids, serialNumber)
		output = append(output, mfaDeviceDataSourceModel{
			SerialNumber:         types.StringValue(serialNumber),
			VirtualMFADeviceName: types.StringValue(name),
			ActivateDate:         mfaDeviceStr(device, "ActivateDate"),
		})
	}

	config.Id = types.StringValue(mfaDevicesIdHash(ids))
	config.MfaDevices = output

	var idValues []string
	for _, id := range ids {
		idValues = append(idValues, id)
	}
	idsList, diags := types.ListValueFrom(ctx, types.StringType, idValues)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	config.Ids = idsList

	if !config.OutputFile.IsNull() && config.OutputFile.ValueString() != "" {
		if err := writeMfaDevicesToFile(config.OutputFile.ValueString(), output); err != nil {
			resp.Diagnostics.AddError(
				fmt.Sprintf("Reading %s: writing output_file", mfaDevicesTypeName),
				err.Error(),
			)
			return
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &config)...)
}

// mfaDevicesIdHash mirrors the data source id hash used across the provider:
// the helper hashcode of the dash-joined ids.
func mfaDevicesIdHash(ids []string) string {
	var buf []byte
	for _, id := range ids {
		buf = append(buf, []byte(fmt.Sprintf("%s-", id))...)
	}
	return fmt.Sprintf("%d", helper.Hashcode(string(buf)))
}

func writeMfaDevicesToFile(path string, devices []mfaDeviceDataSourceModel) error {
	raw := make([]map[string]string, 0, len(devices))
	for _, d := range devices {
		raw = append(raw, map[string]string{
			"serial_number":           d.SerialNumber.ValueString(),
			"virtual_mfa_device_name": d.VirtualMFADeviceName.ValueString(),
			"activate_date":           d.ActivateDate.ValueString(),
		})
	}
	data, err := json.MarshalIndent(raw, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0o644)
}
