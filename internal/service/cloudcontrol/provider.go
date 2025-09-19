// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package cloudcontrol

import (
	"context"

	"github.com/byteplus/terraform-provider-bytepluscc/internal/cloudcontrol"
)

type Provider interface {
	CloudControlAPIClient(context.Context) *cloudcontrol.CloudControl

	Region(ctx context.Context) string

	RegisterLogger(ctx context.Context) context.Context
}
