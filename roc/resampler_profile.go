package roc

// Resampler profile.
// Affects speed and quality.
// Each resampler backend treats profile in its own way.
//
//go:generate stringer -type ResamplerProfile -trimprefix ResamplerProfile
type ResamplerProfile int

const (
	// Do not perform resampling.
	// Clock drift compensation will be disabled in this case.
	// If in doubt, do not disable resampling.
	ResamplerProfileDisable ResamplerProfile = -1

	// Default profile.
	// Current default is ResamplerProfileMedium.
	ResamplerProfileDefault ResamplerProfile = 0

	// High quality, low speed.
	ResamplerProfileHigh ResamplerProfile = 1

	// Medium quality, medium speed.
	ResamplerProfileMedium ResamplerProfile = 2

	// Low quality, high speed.
	ResamplerProfileLow ResamplerProfile = 3
)
