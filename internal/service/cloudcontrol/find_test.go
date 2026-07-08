// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package cloudcontrol

import (
	"encoding/json"
	"testing"

	"github.com/byteplus-sdk/terraform-provider-bytepluscc/internal/cloudcontrol"
	"github.com/volcengine/volcengine-go-sdk/volcengine"
)

func TestNormalizeResourceDescriptionPreservesUserTagFields(t *testing.T) {
	properties := `{
		"Name": "test-resource",
		"Tags": [
			{"Key": "env", "Value": "test", "Type": "CUSTOM"},
			{"Key": "sys:owner", "Value": "platform", "Type": "SYSTEM"}
		]
	}`

	desc := &cloudcontrol.ResourceDescriptionForGetResourceOutput{
		Properties: volcengine.String(properties),
	}

	if err := NormalizeResourceDescription(desc); err != nil {
		t.Fatalf("NormalizeResourceDescription() error = %v", err)
	}

	var got map[string]interface{}
	if err := json.Unmarshal([]byte(*desc.Properties), &got); err != nil {
		t.Fatalf("failed to unmarshal normalized properties: %v", err)
	}

	rawTags, ok := got["Tags"].([]interface{})
	if !ok {
		t.Fatalf("Tags type = %T, want []interface{}", got["Tags"])
	}
	if len(rawTags) != 1 {
		t.Fatalf("len(Tags) = %d, want 1", len(rawTags))
	}

	tag, ok := rawTags[0].(map[string]interface{})
	if !ok {
		t.Fatalf("Tags[0] type = %T, want map[string]interface{}", rawTags[0])
	}
	if tag["Key"] != "env" {
		t.Fatalf("Tags[0].Key = %v, want env", tag["Key"])
	}
	if tag["Value"] != "test" {
		t.Fatalf("Tags[0].Value = %v, want test", tag["Value"])
	}
	if tag["Type"] != "CUSTOM" {
		t.Fatalf("Tags[0].Type = %v, want CUSTOM", tag["Type"])
	}
}
