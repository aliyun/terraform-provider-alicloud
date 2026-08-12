package framework

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	fwschema "github.com/hashicorp/terraform-plugin-framework/provider/schema"
	sdkschema "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// ProviderSchema derives the terraform-plugin-framework provider schema from the
// terraform-plugin-sdk/v2 provider schema map.
//
// tf5muxserver refuses to serve providers whose provider schemas differ once
// converted to tfprotov5 (see tf5muxserver.schemaEquals), so the two schemas must
// stay byte-identical over the protocol. The SDK v2 schema has 27 top-level
// attributes and an "endpoints" block that alone carries ~160 fields and grows with
// every new product, so the framework side is derived from it rather than
// hand-written: a hand-written mirror would break GetProviderSchema for every user
// the first time someone added an endpoint to only one of the two copies.
//
// Constructs the framework cannot express identically produce an error diagnostic
// instead of a silently divergent schema. That turns future drift into a failing
// unit test rather than a runtime "Invalid Provider Server Combination" for users.
func ProviderSchema(sdk map[string]*sdkschema.Schema) (fwschema.Schema, diag.Diagnostics) {
	var diags diag.Diagnostics

	attributes, blocks := convertSchemaMap(sdk, "", &diags)

	return fwschema.Schema{
		Attributes: attributes,
		Blocks:     blocks,
	}, diags
}

// convertSchemaMap splits an SDK v2 schema map into framework attributes and blocks,
// mirroring the attribute/block routing in schemaMap.CoreConfigSchema.
func convertSchemaMap(sdk map[string]*sdkschema.Schema, prefix string, diags *diag.Diagnostics) (map[string]fwschema.Attribute, map[string]fwschema.Block) {
	attributes := make(map[string]fwschema.Attribute, len(sdk))
	blocks := make(map[string]fwschema.Block)

	for name, s := range sdk {
		path := name
		if prefix != "" {
			path = prefix + "." + name
		}

		if s.ConfigMode != sdkschema.SchemaConfigModeAuto {
			addUnsupportedError(diags, path, "ConfigMode is set; only SchemaConfigModeAuto is supported")
			continue
		}

		if s.Elem == nil {
			if attribute, ok := convertAttribute(s, path, diags); ok {
				attributes[name] = attribute
			}
			continue
		}

		elem, ok := s.Elem.(*sdkschema.Resource)
		if !ok {
			addUnsupportedError(diags, path, fmt.Sprintf("Elem is %T; only *schema.Resource is supported", s.Elem))
			continue
		}

		if block, ok := convertBlock(s, elem, path, diags); ok {
			blocks[name] = block
		}
	}

	return attributes, blocks
}

// convertAttribute converts a primitive SDK v2 schema into a framework attribute.
func convertAttribute(s *sdkschema.Schema, path string, diags *diag.Diagnostics) (fwschema.Attribute, bool) {
	// Framework provider schema attributes have no Computed or WriteOnly field, so
	// either one would convert to a proto attribute that differs from the SDK v2 one.
	if s.Computed {
		addUnsupportedError(diags, path, "Computed is not expressible in a framework provider schema")
		return nil, false
	}
	if s.WriteOnly {
		addUnsupportedError(diags, path, "WriteOnly is not expressible in a framework provider schema")
		return nil, false
	}

	required, optional := requiredOptional(s)

	// Description is deliberately used in place of MarkdownDescription: SDK v2 emits
	// tfprotov5.StringKindPlain unless the global schema.DescriptionKind is switched,
	// and MarkdownDescription would emit StringKindMarkdown and fail the mux
	// comparison.
	switch s.Type {
	case sdkschema.TypeString:
		return fwschema.StringAttribute{
			Required:           required,
			Optional:           optional,
			Sensitive:          s.Sensitive,
			Description:        s.Description,
			DeprecationMessage: s.Deprecated,
		}, true
	case sdkschema.TypeInt:
		return fwschema.Int64Attribute{
			Required:           required,
			Optional:           optional,
			Sensitive:          s.Sensitive,
			Description:        s.Description,
			DeprecationMessage: s.Deprecated,
		}, true
	case sdkschema.TypeBool:
		return fwschema.BoolAttribute{
			Required:           required,
			Optional:           optional,
			Sensitive:          s.Sensitive,
			Description:        s.Description,
			DeprecationMessage: s.Deprecated,
		}, true
	}

	addUnsupportedError(diags, path, fmt.Sprintf("type %s is not supported; only TypeString, TypeInt and TypeBool are", s.Type))

	return nil, false
}

// convertBlock converts an SDK v2 schema whose Elem is a *schema.Resource into a
// framework nested block.
//
// MinItems and MaxItems are dropped: the framework cannot express them, and
// tf5muxserver ignores both fields when comparing schemas.
func convertBlock(s *sdkschema.Schema, elem *sdkschema.Resource, path string, diags *diag.Diagnostics) (fwschema.Block, bool) {
	// coreConfigSchemaBlock lifts the outer schema's description and deprecation
	// onto the nested block itself, so they belong on the block here rather than
	// anywhere inside NestedObject.
	nested := fwschema.NestedBlockObject{}
	nested.Attributes, nested.Blocks = convertSchemaMap(elem.Schema, path, diags)

	switch s.Type {
	case sdkschema.TypeList:
		return fwschema.ListNestedBlock{
			NestedObject:       nested,
			Description:        s.Description,
			DeprecationMessage: s.Deprecated,
		}, true
	case sdkschema.TypeSet:
		return fwschema.SetNestedBlock{
			NestedObject:       nested,
			Description:        s.Description,
			DeprecationMessage: s.Deprecated,
		}, true
	}

	addUnsupportedError(diags, path, fmt.Sprintf("type %s with a *schema.Resource Elem is not supported; only TypeList and TypeSet are", s.Type))

	return nil, false
}

// requiredOptional mirrors the required-ness demotion in
// (*schema.Schema).coreConfigSchemaAttribute: SDK v2 reports a Required attribute
// that also has a DefaultFunc as Optional whenever that DefaultFunc yields a value,
// because Terraform core has no concept of conditional required-ness. The framework
// has no equivalent, so the same sniff has to be applied here or the two schemas
// disagree — for example assume_role.role_arn flips to Optional as soon as
// ALICLOUD_ASSUME_ROLE_ARN is set in the environment.
func requiredOptional(s *sdkschema.Schema) (required, optional bool) {
	required, optional = s.Required, s.Optional

	if required && s.DefaultFunc != nil {
		v, err := s.DefaultFunc()
		if err != nil || v != nil {
			required, optional = false, true
		}
	}

	return required, optional
}

func addUnsupportedError(diags *diag.Diagnostics, path, detail string) {
	diags.AddError(
		"Unsupported Provider Schema Construct",
		fmt.Sprintf("The framework provider schema could not be derived from the SDK v2 provider schema at %q: %s.\n\n"+
			"This is a bug in the provider: alicloud/provider/framework/schema.go needs to learn this construct, "+
			"otherwise the muxed provider servers would expose differing provider schemas.", path, detail),
	)
}
