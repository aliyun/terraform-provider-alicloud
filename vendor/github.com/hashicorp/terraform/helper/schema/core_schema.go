package schema

import (
	"fmt"

	"github.com/hashicorp/terraform/configs/configschema"
	"github.com/zclconf/go-cty/cty"
)

// The functions and methods in this file are concerned with the conversion
// of this package's schema model into the slightly-lower-level schema model
// used by Terraform core for configuration parsing.

// CoreConfigSchema lowers the receiver to the schema model expected by
// Terraform core.
//
// This lower-level model has fewer features than the schema in this package,
// describing only the basic structure of configuration and state values we
// expect. The full schemaMap from this package is still required for full
// validation, handling of default values, etc.
//
// This method presumes a schema that passes InternalValidate, and so may
// panic or produce an invalid result if given an invalid schemaMap.
func (m schemaMap) CoreConfigSchema() *configschema.Block {
	return m.coreConfigSchema(true, true)
}

// CoreConfigSchemaForShimming is a variant of CoreConfigSchema that returns
// the schema that should be used when applying our shimming behaviors.
//
// In particular, it ignores the SkipCoreTypeCheck flag on any legacy schemas,
// since the shims live on the SDK side and so they need to see the full
// type information that we'd normally hide from Terraform Core when skipping
// type checking over there.
func (m schemaMap) CoreConfigSchemaForShimming() *configschema.Block {
	return m.coreConfigSchema(true, false)
}

// CoreConfigSchemaWhenShimmed is a variant of CoreConfigSchema that returns
// the schema as it would appear when working with data structures that have
// already been shimmed to the legacy form.
//
// In particular, it ignores the AsSingle flag on any legacy schemas and behaves
// as if they were really lists/sets instead, thus giving a description of
// the shape of the data structure after the AsSingle fixup has been applied.
//
// This should be used with care only in unusual situations where we need to
// work with an already-shimmed value using a new-style schema.
func (m schemaMap) CoreConfigSchemaWhenShimmed() *configschema.Block {
	return m.coreConfigSchema(false, false)
}

func (m schemaMap) coreConfigSchema(asSingle, skipCoreCheck bool) *configschema.Block {
	if len(m) == 0 {
		// We return an actual (empty) object here, rather than a nil,
		// because a nil result would mean that we don't have a schema at
		// all, rather than that we have an empty one.
		return &configschema.Block{}
	}

	ret := &configschema.Block{
		Attributes: map[string]*configschema.Attribute{},
		BlockTypes: map[string]*configschema.NestedBlock{},
	}

	for name, schema := range m {
		if schema.Elem == nil {
			ret.Attributes[name] = schema.coreConfigSchemaAttribute(asSingle, skipCoreCheck)
			continue
		}
		if schema.Type == TypeMap {
			// For TypeMap in particular, it isn't valid for Elem to be a
			// *Resource (since that would be ambiguous in flatmap) and
			// so Elem is treated as a TypeString schema if so. This matches
			// how the field readers treat this situation, for compatibility
			// with configurations targeting Terraform 0.11 and earlier.
			if _, isResource := schema.Elem.(*Resource); isResource {
				sch := *schema // shallow copy
				sch.Elem = &Schema{
					Type: TypeString,
				}
				ret.Attributes[name] = sch.coreConfigSchemaAttribute(asSingle, skipCoreCheck)
				continue
			}
		}
		switch schema.ConfigMode {
		case SchemaConfigModeAttr:
			ret.Attributes[name] = schema.coreConfigSchemaAttribute(asSingle, skipCoreCheck)
		case SchemaConfigModeBlock:
			ret.BlockTypes[name] = schema.coreConfigSchemaBlock(asSingle, skipCoreCheck)
		default: // SchemaConfigModeAuto, or any other invalid value
			if schema.Computed && !schema.Optional {
				// Computed-only schemas are always handled as attributes,
				// because they never appear in configuration.
				ret.Attributes[name] = schema.coreConfigSchemaAttribute(asSingle, skipCoreCheck)
				continue
			}
			switch schema.Elem.(type) {
			case *Schema, ValueType:
				ret.Attributes[name] = schema.coreConfigSchemaAttribute(asSingle, skipCoreCheck)
			case *Resource:
				ret.BlockTypes[name] = schema.coreConfigSchemaBlock(asSingle, skipCoreCheck)
			default:
				// Should never happen for a valid schema
				panic(fmt.Errorf("invalid Schema.Elem %#v; need *Schema or *Resource", schema.Elem))
			}
		}
	}

	return ret
}

// coreConfigSchemaAttribute prepares a configschema.Attribute representation
// of a schema. This is appropriate only for primitives or collections whose
// Elem is an instance of Schema. Use coreConfigSchemaBlock for collections
// whose elem is a whole resource.
func (s *Schema) coreConfigSchemaAttribute(asSingle, skipCoreCheck bool) *configschema.Attribute {
	// The Schema.DefaultFunc capability adds some extra weirdness here since
	// it can be combined with "Required: true" to create a sitution where
	// required-ness is conditional. Terraform Core doesn't share this concept,
	// so we must sniff for this possibility here and conditionally turn
	// off the "Required" flag if it looks like the DefaultFunc is going
	// to provide a value.
	// This is not 100% true to the original interface of DefaultFunc but
	// works well enough for the EnvDefaultFunc and MultiEnvDefaultFunc
	// situations, which are the main cases we care about.
	//
	// Note that this also has a consequence for commands that return schema
	// information for documentation purposes: running those for certain
	// providers will produce different results depending on which environment
	// variables are set. We accept that weirdness in order to keep this
	// interface to core otherwise simple.
	reqd := s.Required
	opt := s.Optional
	if reqd && s.DefaultFunc != nil {
		v, err := s.DefaultFunc()
		// We can't report errors from here, so we'll instead just force
		// "Required" to false and let the provider try calling its
		// DefaultFunc again during the validate step, where it can then
		// return the error.
		if err != nil || (err == nil && v != nil) {
			reqd = false
			opt = true
		}
	}

	return &configschema.Attribute{
		Type:        s.coreConfigSchemaType(asSingle, skipCoreCheck),
		Optional:    opt,
		Required:    reqd,
		Computed:    s.Computed,
		Sensitive:   s.Sensitive,
		Description: s.Description,
	}
}

// coreConfigSchemaBlock prepares a configschema.NestedBlock representation of
// a schema. This is appropriate only for collections whose Elem is an instance
// of Resource, and will panic otherwise.
func (s *Schema) coreConfigSchemaBlock(asSingle, skipCoreCheck bool) *configschema.NestedBlock {
	ret := &configschema.NestedBlock{}
	if nested := schemaMap(s.Elem.(*Resource).Schema).coreConfigSchema(asSingle, skipCoreCheck); nested != nil {
		ret.Block = *nested
	}
	switch s.Type {
	case TypeList:
		ret.Nesting = configschema.NestingList
	case TypeSet:
		ret.Nesting = configschema.NestingSet
	case TypeMap:
		ret.Nesting = configschema.NestingMap
	default:
		// Should never happen for a valid schema
		panic(fmt.Errorf("invalid s.Type %s for s.Elem being resource", s.Type))
	}

	ret.MinItems = s.MinItems
	ret.MaxItems = s.MaxItems

	if s.AsSingle && asSingle {
		// In AsSingle mode, we artifically force a TypeList or TypeSet
		// attribute in the SDK to be treated as a single block by Terraform Core.
		// This must then be fixed up in the shim code (in helper/plugin) so
		// that the SDK still sees the lists or sets it's expecting.
		ret.Nesting = configschema.NestingSingle
	}

	if s.Required && s.MinItems == 0 {
		// configschema doesn't have a "required" representation for nested
		// blocks, but we can fake it by requiring at least one item.
		ret.MinItems = 1
	}
	if s.Optional && s.MinItems > 0 {
		// Historically helper/schema would ignore MinItems if Optional were
		// set, so we must mimic this behavior here to ensure that providers
		// relying on that undocumented behavior can continue to operate as
		// they did before.
		ret.MinItems = 0
	}
	if s.Computed && !s.Optional {
		// MinItems/MaxItems are meaningless for computed nested blocks, since
		// they are never set by the user anyway. This ensures that we'll never
		// generate weird errors about them.
		ret.MinItems = 0
		ret.MaxItems = 0
	}

	return ret
}

// coreConfigSchemaType determines the core config schema type that corresponds
// to a particular schema's type.
func (s *Schema) coreConfigSchemaType(asSingle, skipCoreCheck bool) cty.Type {
	if skipCoreCheck && s.SkipCoreTypeCheck {
		// If we're preparing a schema for Terraform Core and the schema is
		// asking us to skip the Core type-check then we'll tell core that this
		// attribute is dynamically-typed, so it'll just pass through anything
		// and let us validate it on the plugin side.
		return cty.DynamicPseudoType
	}

	switch s.Type {
	case TypeString:
		return cty.String
	case TypeBool:
		return cty.Bool
	case TypeInt, TypeFloat:
		// configschema doesn't distinguish int and float, so helper/schema
		// will deal with this as an additional validation step after
		// configuration has been parsed and decoded.
		return cty.Number
	case TypeList, TypeSet, TypeMap:
		var elemType cty.Type
		switch set := s.Elem.(type) {
		case *Schema:
			elemType = set.coreConfigSchemaType(asSingle, skipCoreCheck)
		case ValueType:
			// This represents a mistake in the provider code, but it's a
			// common one so we'll just shim it.
			elemType = (&Schema{Type: set}).coreConfigSchemaType(asSingle, skipCoreCheck)
		case *Resource:
			// By default we construct a NestedBlock in this case, but this
			// behavior is selected either for computed-only schemas or
			// when ConfigMode is explicitly SchemaConfigModeBlock.
			// See schemaMap.CoreConfigSchema for the exact rules.
			elemType = schemaMap(set.Schema).coreConfigSchema(asSingle, skipCoreCheck).ImpliedType()
		default:
			if set != nil {
				// Should never happen for a valid schema
				panic(fmt.Errorf("invalid Schema.Elem %#v; need *Schema or *Resource", s.Elem))
			}
			// Some pre-existing schemas assume string as default, so we need
			// to be compatible with them.
			elemType = cty.String
		}
		if s.AsSingle && asSingle {
			// In AsSingle mode, we artifically force a TypeList or TypeSet
			// attribute in the SDK to be treated as a single value by Terraform Core.
			// This must then be fixed up in the shim code (in helper/plugin) so
			// that the SDK still sees the lists or sets it's expecting.
			return elemType
		}
		switch s.Type {
		case TypeList:
			return cty.List(elemType)
		case TypeSet:
			return cty.Set(elemType)
		case TypeMap:
			return cty.Map(elemType)
		default:
			// can never get here in practice, due to the case we're inside
			panic("invalid collection type")
		}
	default:
		// should never happen for a valid schema
		panic(fmt.Errorf("invalid Schema.Type %s", s.Type))
	}
}

// CoreConfigSchema is a convenient shortcut for calling CoreConfigSchema on
// the resource's schema. CoreConfigSchema adds the implicitly required "id"
// attribute for top level resources if it doesn't exist.
func (r *Resource) CoreConfigSchema() *configschema.Block {
	return r.coreConfigSchema(true, true)
}

// CoreConfigSchemaForShimming is a variant of CoreConfigSchema that returns
// the schema that should be used to apply shims on the SDK side.
//
// In particular, it ignores the SkipCoreTypeCheck flag on any legacy schemas
// and uses the real type information instead.
func (r *Resource) CoreConfigSchemaForShimming() *configschema.Block {
	return r.coreConfigSchema(true, false)
}

// CoreConfigSchemaWhenShimmed is a variant of CoreConfigSchema that returns
// the schema as it would appear when working with data structures that have
// already been shimmed to the legacy form.
//
// In particular, it ignores the AsSingle flag on any legacy schemas and behaves
// as if they were really lists/sets instead, thus giving a description of
// the shape of the data structure after the AsSingle fixup has been applied.
//
// This should be used with care only in unusual situations where we need to
// work with an already-shimmed value using a new-style schema.
func (r *Resource) CoreConfigSchemaWhenShimmed() *configschema.Block {
	return r.coreConfigSchema(false, false)
}

func (r *Resource) coreConfigSchema(asSingle, skipCoreCheck bool) *configschema.Block {
	block := schemaMap(r.Schema).coreConfigSchema(asSingle, skipCoreCheck)

	if block.Attributes == nil {
		block.Attributes = map[string]*configschema.Attribute{}
	}

	// Add the implicitly required "id" field if it doesn't exist
	if block.Attributes["id"] == nil {
		block.Attributes["id"] = &configschema.Attribute{
			Type:     cty.String,
			Optional: true,
			Computed: true,
		}
	}

	_, timeoutsAttr := block.Attributes[TimeoutsConfigKey]
	_, timeoutsBlock := block.BlockTypes[TimeoutsConfigKey]

	// Insert configured timeout values into the schema, as long as the schema
	// didn't define anything else by that name.
	if r.Timeouts != nil && !timeoutsAttr && !timeoutsBlock {
		timeouts := configschema.Block{
			Attributes: map[string]*configschema.Attribute{},
		}

		if r.Timeouts.Create != nil {
			timeouts.Attributes[TimeoutCreate] = &configschema.Attribute{
				Type:     cty.String,
				Optional: true,
			}
		}

		if r.Timeouts.Read != nil {
			timeouts.Attributes[TimeoutRead] = &configschema.Attribute{
				Type:     cty.String,
				Optional: true,
			}
		}

		if r.Timeouts.Update != nil {
			timeouts.Attributes[TimeoutUpdate] = &configschema.Attribute{
				Type:     cty.String,
				Optional: true,
			}
		}

		if r.Timeouts.Delete != nil {
			timeouts.Attributes[TimeoutDelete] = &configschema.Attribute{
				Type:     cty.String,
				Optional: true,
			}
		}

		if r.Timeouts.Default != nil {
			timeouts.Attributes[TimeoutDefault] = &configschema.Attribute{
				Type:     cty.String,
				Optional: true,
			}
		}

		block.BlockTypes[TimeoutsConfigKey] = &configschema.NestedBlock{
			Nesting: configschema.NestingSingle,
			Block:   timeouts,
		}
	}

	return block
}

// CoreConfigSchema is a convenient shortcut for calling CoreConfigSchema
// on the backends's schema.
func (r *Backend) CoreConfigSchema() *configschema.Block {
	return schemaMap(r.Schema).CoreConfigSchema()
}
