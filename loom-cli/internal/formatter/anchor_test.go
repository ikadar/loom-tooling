package formatter

import (
	"testing"
)

func TestToAnchor(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "simple uppercase ID",
			input:    "AC-ORD-001",
			expected: "ac-ord-001",
		},
		{
			name:     "test case ID",
			input:    "TC-AC-CUST-001-P01",
			expected: "tc-ac-cust-001-p01",
		},
		{
			name:     "already lowercase",
			input:    "br-ord-001",
			expected: "br-ord-001",
		},
		{
			name:     "mixed case",
			input:    "Ts-BrOrd-001",
			expected: "ts-brord-001",
		},
		{
			name:     "empty string",
			input:    "",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ToAnchor(tt.input)
			if result != tt.expected {
				t.Errorf("ToAnchor(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestToLink(t *testing.T) {
	tests := []struct {
		name     string
		id       string
		file     string
		expected string
	}{
		{
			name:     "AC reference",
			id:       "AC-ORD-001",
			file:     "../l1/acceptance-criteria.md",
			expected: "[AC-ORD-001](../l1/acceptance-criteria.md#ac-ord-001)",
		},
		{
			name:     "BR reference",
			id:       "BR-CUST-002",
			file:     "../l1/business-rules.md",
			expected: "[BR-CUST-002](../l1/business-rules.md#br-cust-002)",
		},
		{
			name:     "test case reference",
			id:       "TC-AC-ORD-001-P01",
			file:     "./test-cases.md",
			expected: "[TC-AC-ORD-001-P01](./test-cases.md#tc-ac-ord-001-p01)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ToLink(tt.id, tt.file)
			if result != tt.expected {
				t.Errorf("ToLink(%q, %q) = %q, want %q", tt.id, tt.file, result, tt.expected)
			}
		})
	}
}

func TestFormatHeader(t *testing.T) {
	result := FormatHeader("Test Document", "2024-01-15T10:00:00Z")

	// Check title
	if !contains(result, "# Test Document") {
		t.Error("Expected title '# Test Document' in header")
	}

	// Check timestamp
	if !contains(result, "Generated: 2024-01-15T10:00:00Z") {
		t.Error("Expected timestamp in header")
	}

	// Check separator
	if !contains(result, "---") {
		t.Error("Expected separator '---' in header")
	}
}

func TestFormatSectionHeader(t *testing.T) {
	tests := []struct {
		name     string
		level    int
		id       string
		title    string
		wantH    string
		wantID   string
	}{
		{
			name:   "level 2 header",
			level:  2,
			id:     "TS-ORD-001",
			title:  "Order Validation",
			wantH:  "## TS-ORD-001 – Order Validation",
			wantID: "{#ts-ord-001}",
		},
		{
			name:   "level 3 header",
			level:  3,
			id:     "IC-CUST-001",
			title:  "Customer Service",
			wantH:  "### IC-CUST-001 – Customer Service",
			wantID: "{#ic-cust-001}",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FormatSectionHeader(tt.level, tt.id, tt.title)

			if !contains(result, tt.wantH) {
				t.Errorf("Expected header %q in result, got %q", tt.wantH, result)
			}

			if !contains(result, tt.wantID) {
				t.Errorf("Expected anchor %q in result, got %q", tt.wantID, result)
			}
		})
	}
}

// Helper function
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 ||
		(len(s) > 0 && len(substr) > 0 && findSubstring(s, substr)))
}

func findSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
