package confetti

import (
	"testing"
)

// TestNewDefLoader tests if NewDefLoader returns a valid ILoader.
func TestNewDefLoader(t *testing.T) {
	loader := NewDefLoader()
	if loader == nil {
		t.Errorf("Expected loader to be non-nil, but it is nil.")
		return
	}
}

// TestNewLoader tests if the NewLoader returns a valid ILoader.
func TestNewLoader(t *testing.T) {
	loader := NewLoader(LoaderOptions{})
	if loader == nil {
		t.Errorf("Expected loader to be non-nil, but it is nil.")
		return
	}
}

// TestNewFlagger tests if newFlagger returns a valid iFlagger.
func TestNewFlagger(t *testing.T) {
	flagger := newFlagger(defaultLoaderOptions)
	if flagger == nil {
		t.Errorf("Expected flagger to be non-nil, but it is nil.")
		return
	}
}

// TestNewResolver tests if newResolver returns a valid iResolver.
func TestNewResolver(t *testing.T) {
	resolver := newResolver(defaultLoaderOptions)
	if resolver == nil {
		t.Errorf("Expected resolver to be non-nil, but it is nil.")
		return
	}
}
