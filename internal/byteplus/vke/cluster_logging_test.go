// Copyright (c) 2025 Beijing Volcano Engine Technology Co., Ltd.
// SPDX-License-Identifier: MPL-2.0

package vke

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestSanitizeLoggingConfig(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		input   map[string]any
		want    map[string]any
		removed bool
	}{
		{
			name:    "absent",
			input:   map[string]any{"Name": "c1"},
			want:    map[string]any{"Name": "c1"},
			removed: false,
		},
		{
			name:    "null",
			input:   map[string]any{"Name": "c1", "LoggingConfig": nil},
			want:    map[string]any{"Name": "c1"},
			removed: true,
		},
		{
			name: "empty object",
			input: map[string]any{
				"Name":          "c1",
				"LoggingConfig": map[string]any{},
			},
			want:    map[string]any{"Name": "c1"},
			removed: true,
		},
		{
			name: "empty log setups",
			input: map[string]any{
				"Name": "c1",
				"LoggingConfig": map[string]any{
					"LogSetups": []any{},
				},
			},
			want:    map[string]any{"Name": "c1"},
			removed: true,
		},
		{
			name: "disabled audit only",
			input: map[string]any{
				"Name": "c1",
				"LoggingConfig": map[string]any{
					"LogSetups": []any{
						map[string]any{
							"LogType": "Audit",
							"Enabled": false,
							"LogTtl":  float64(30),
						},
					},
				},
			},
			want:    map[string]any{"Name": "c1"},
			removed: true,
		},
		{
			name: "enabled audit kept",
			input: map[string]any{
				"Name": "c1",
				"LoggingConfig": map[string]any{
					"LogSetups": []any{
						map[string]any{
							"LogType": "Audit",
							"Enabled": true,
							"LogTtl":  float64(30),
						},
					},
				},
			},
			want: map[string]any{
				"Name": "c1",
				"LoggingConfig": map[string]any{
					"LogSetups": []any{
						map[string]any{
							"LogType": "Audit",
							"Enabled": true,
							"LogTtl":  float64(30),
						},
					},
				},
			},
			removed: false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			gotRemoved := sanitizeLoggingConfig(tt.input)
			if gotRemoved != tt.removed {
				t.Fatalf("removed = %v, want %v", gotRemoved, tt.removed)
			}
			if diff := cmp.Diff(tt.want, tt.input); diff != "" {
				t.Fatalf("targetState mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestPrepareNativeCreateClusterInput(t *testing.T) {
	t.Parallel()

	in := map[string]any{
		"Name":                    "demo",
		"KubernetesVersionCreate": "1.32",
		"ClusterId":               "should-drop",
		"Status":                  map[string]any{"Phase": "Creating"},
		"LoggingConfig": map[string]any{
			"LogSetups": []any{
				map[string]any{"LogType": "Audit", "Enabled": false, "LogTtl": float64(30)},
			},
		},
		"ClusterConfig": map[string]any{
			"SubnetIds": []any{"subnet-1"},
		},
	}

	got := prepareNativeCreateClusterInput(in)
	want := map[string]any{
		"Name":              "demo",
		"KubernetesVersion": "1.32",
		"ClusterConfig": map[string]any{
			"SubnetIds": []any{"subnet-1"},
		},
	}
	if diff := cmp.Diff(want, got); diff != "" {
		t.Fatalf("prepareNativeCreateClusterInput mismatch (-want +got):\n%s", diff)
	}
	// Input must not be mutated.
	if _, ok := in["KubernetesVersionCreate"]; !ok {
		t.Fatal("expected input KubernetesVersionCreate to remain")
	}
}

func TestIsInvalidLoggingConfigError(t *testing.T) {
	t.Parallel()

	if !isInvalidLoggingConfigError(errString("InvalidParameter.LoggingConfig: The specified parameter LoggingConfig is invalid.")) {
		t.Fatal("expected match")
	}
	if isInvalidLoggingConfigError(errString("InvalidParameter.Name")) {
		t.Fatal("did not expect match")
	}
	if isInvalidLoggingConfigError(nil) {
		t.Fatal("nil should be false")
	}
}

type errString string

func (e errString) Error() string { return string(e) }
