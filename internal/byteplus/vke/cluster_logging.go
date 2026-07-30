// Copyright (c) 2025 Beijing Volcano Engine Technology Co., Ltd.
// SPDX-License-Identifier: MPL-2.0

package vke

import (
	"encoding/json"
	"strings"
)

// sanitizeLoggingConfig removes LoggingConfig from a Cloud Control TargetState
// when it would be rejected by the VKE CreateCluster API.
//
// Background: Byteplus Cloud Control create for Byteplus::VKE::Cluster can fail
// with InvalidParameter.LoggingConfig when LoggingConfig is empty, only contains
// disabled log setups, or is otherwise incomplete. Native VKE CreateCluster
// succeeds when LoggingConfig is omitted entirely. Terraform's nested
// log_ttl default (30) also makes it easy to send a disabled Audit setup that
// the API rejects.
//
// Returns true when LoggingConfig was removed from targetState.
func sanitizeLoggingConfig(targetState map[string]any) bool {
	raw, ok := targetState["LoggingConfig"]
	if !ok {
		return false
	}
	if raw == nil {
		delete(targetState, "LoggingConfig")
		return true
	}

	cfg, ok := raw.(map[string]any)
	if !ok {
		delete(targetState, "LoggingConfig")
		return true
	}

	if !loggingConfigHasEnabledSetup(cfg) {
		delete(targetState, "LoggingConfig")
		return true
	}
	return false
}

func loggingConfigHasEnabledSetup(cfg map[string]any) bool {
	setupsRaw, ok := cfg["LogSetups"]
	if !ok || setupsRaw == nil {
		return false
	}

	setups, ok := setupsRaw.([]any)
	if !ok {
		// JSON decode of a Terraform set may still be []any; anything else is unusable.
		return false
	}
	if len(setups) == 0 {
		return false
	}

	for _, item := range setups {
		setup, ok := item.(map[string]any)
		if !ok {
			continue
		}
		enabled, ok := setup["Enabled"].(bool)
		if ok && enabled {
			return true
		}
	}
	return false
}

func isInvalidLoggingConfigError(err error) bool {
	if err == nil {
		return false
	}
	msg := err.Error()
	return strings.Contains(msg, "InvalidParameter.LoggingConfig") ||
		strings.Contains(msg, "LoggingConfig is invalid")
}

// prepareNativeCreateClusterInput clones TargetState into a body suitable for
// the native VKE CreateCluster API (not Cloud Control).
func prepareNativeCreateClusterInput(targetState map[string]any) map[string]any {
	out := cloneMap(targetState)
	sanitizeLoggingConfig(out)

	// Cloud Control write-only create property → native CreateCluster field.
	if v, ok := out["KubernetesVersionCreate"]; ok {
		if _, exists := out["KubernetesVersion"]; !exists {
			out["KubernetesVersion"] = v
		}
		delete(out, "KubernetesVersionCreate")
	}

	// Drop read-only / CC-only fields that must not be sent to CreateCluster.
	for _, key := range []string{
		"ClusterId",
		"CreatedTime",
		"UpdatedTime",
		"Status",
		"Message",
		"Type",
		"NodeStatistics",
		"SourceRegion",
	} {
		delete(out, key)
	}

	return out
}

func cloneMap(in map[string]any) map[string]any {
	b, err := json.Marshal(in)
	if err != nil {
		out := make(map[string]any, len(in))
		for k, v := range in {
			out[k] = v
		}
		return out
	}
	out := make(map[string]any)
	_ = json.Unmarshal(b, &out)
	return out
}
