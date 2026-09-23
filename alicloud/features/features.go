// Package features holds the provider arguments that change how a resource behaves, as opposed to
// which account, region or endpoint the provider talks to. It is a leaf package: it imports
// nothing, so a resource can read a toggle without pulling in connectivity.
package features

// Features mirrors the provider's features block.
type Features struct {
	EcsInstance EcsInstance
}

// EcsInstance holds the toggles of the features.ecs_instance block.
type EcsInstance struct {
	// ReplaceOnImageUpdate makes an image_id change plan as a replacement of the instance instead of
	// an in-place replacement of its system disk.
	ReplaceOnImageUpdate bool
}

// Default returns the behaviour of a provider that configures no features block at all. Defaults
// live here as well as in the schema because an absent nested block contributes no schema default.
func Default() Features {
	return Features{
		EcsInstance: EcsInstance{
			ReplaceOnImageUpdate: false,
		},
	}
}
