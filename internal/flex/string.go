// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package flex

import (
	"context"

	"github.com/byteplus-sdk/terraform-provider-bytepluscc/internal/util"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// StringFromFramework converts a Framework String value to a string pointer.
// A null String is converted to a nil string pointer.
func StringFromFramework(_ context.Context, v types.String) *string {
	if v.IsNull() || v.IsUnknown() {
		return nil
	}

	return util.StringPtr(v.ValueString())
}
