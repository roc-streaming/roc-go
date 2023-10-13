package roc

// Resampler backend.
// Affects speed and quality.
// Some backends may be disabled at build time.
//
//go:generate stringer -type ResamplerBackend -trimprefix ResamplerBackend
type ResamplerBackend int

const (
	// Default backend.
	// Depends on what was enabled at build time.
	ResamplerBackendDefault ResamplerBackend = 0

	// Slow built-in resampler.
	// Always available.
	ResamplerBackendBuiltin ResamplerBackend = 1

	// Fast good-quality resampler from SpeexDSP.
	// May be disabled at build time.
	ResamplerBackendSpeex ResamplerBackend = 2
)
