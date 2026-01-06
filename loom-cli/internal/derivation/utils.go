package derivation

import (
	"path/filepath"
	"strings"
)

// =============================================================================
// Path Utilities
// =============================================================================

// detectLayerFromPath determines the layer (l0, l1, l2, l3) from file path
func detectLayerFromPath(path string) string {
	pathNorm := filepath.ToSlash(path)

	if strings.Contains(pathNorm, "/l3/") || strings.Contains(pathNorm, "/l3-") {
		return "l3"
	}
	if strings.Contains(pathNorm, "/l2/") || strings.Contains(pathNorm, "/l2-") {
		return "l2"
	}
	if strings.Contains(pathNorm, "/l1/") || strings.Contains(pathNorm, "/l1-") {
		return "l1"
	}
	if strings.Contains(pathNorm, "/l0/") || strings.Contains(pathNorm, "/l0-") {
		return "l0"
	}

	// Try to detect from filename patterns
	baseName := strings.ToLower(filepath.Base(path))
	if strings.HasPrefix(baseName, "user-stor") || strings.HasPrefix(baseName, "nfr") {
		return "l0"
	}
	if strings.HasPrefix(baseName, "acceptance") || strings.HasPrefix(baseName, "business-rule") ||
		strings.HasPrefix(baseName, "domain-model") || strings.HasPrefix(baseName, "bounded-context") {
		return "l1"
	}
	if strings.HasPrefix(baseName, "tech-spec") || strings.HasPrefix(baseName, "interface-contract") ||
		strings.HasPrefix(baseName, "aggregate") || strings.HasPrefix(baseName, "sequence") {
		return "l2"
	}
	if strings.HasPrefix(baseName, "test-case") || strings.HasPrefix(baseName, "openapi") ||
		strings.HasPrefix(baseName, "skeleton") || strings.HasPrefix(baseName, "ticket") {
		return "l3"
	}

	// Default to l1 for spec files without clear layer
	return "l1"
}

// =============================================================================
// ID Utilities
// =============================================================================

// GetLayerFromID determines the layer from an artifact ID prefix
func GetLayerFromID(id string) string {
	// L0 prefixes
	if strings.HasPrefix(id, "US-") || strings.HasPrefix(id, "NFR-") {
		return "l0"
	}

	// L1 prefixes
	if strings.HasPrefix(id, "AC-") || strings.HasPrefix(id, "BR-") ||
		strings.HasPrefix(id, "ENT-") || strings.HasPrefix(id, "VO-") ||
		strings.HasPrefix(id, "BC-") {
		return "l1"
	}

	// L2 prefixes
	if strings.HasPrefix(id, "TS-") || strings.HasPrefix(id, "IC-") ||
		strings.HasPrefix(id, "AGG-") || strings.HasPrefix(id, "SEQ-") ||
		strings.HasPrefix(id, "DT-") {
		return "l2"
	}

	// L3 prefixes
	if strings.HasPrefix(id, "TC-") || strings.HasPrefix(id, "API-") ||
		strings.HasPrefix(id, "EVT-") || strings.HasPrefix(id, "CMD-") ||
		strings.HasPrefix(id, "TKT-") || strings.HasPrefix(id, "SKL-") {
		return "l3"
	}

	return "unknown"
}

// =============================================================================
// String Utilities
// =============================================================================

// containsAny checks if s contains any of the substrings
func containsAny(s string, substrs ...string) bool {
	for _, substr := range substrs {
		if strings.Contains(s, substr) {
			return true
		}
	}
	return false
}

// uniqueStrings returns unique strings from a slice
func uniqueStrings(input []string) []string {
	seen := make(map[string]bool)
	result := make([]string, 0, len(input))

	for _, s := range input {
		if !seen[s] {
			seen[s] = true
			result = append(result, s)
		}
	}

	return result
}

// filterStrings returns strings that match a predicate
func filterStrings(input []string, predicate func(string) bool) []string {
	result := make([]string, 0)
	for _, s := range input {
		if predicate(s) {
			result = append(result, s)
		}
	}
	return result
}
