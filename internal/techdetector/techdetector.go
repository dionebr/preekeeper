package techdetector

import (
	wapp "github.com/projectdiscovery/wappalyzergo"
)

// Engine is a small wrapper around the external detection engine.
// This package hides the third-party package name from the rest of the codebase.
type Engine struct {
	eng *wapp.Wappalyze
}

// New creates a new detection engine instance.
func New() (*Engine, error) {
	eng, err := wapp.New()
	if err != nil {
		return nil, err
	}
	return &Engine{eng: eng}, nil
}

// Fingerprint delegates to the external engine's Fingerprint method.
func (e *Engine) Fingerprint(headers map[string][]string, body []byte) map[string]struct{} {
	return e.eng.Fingerprint(headers, body)
}
