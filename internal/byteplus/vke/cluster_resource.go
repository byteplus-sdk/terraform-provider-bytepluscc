// Copyright (c) 2025 Beijing Volcano Engine Technology Co., Ltd.
// SPDX-License-Identifier: MPL-2.0

// Hand-written customization for `bytepluscc_vke_cluster`.
//
// Cloud Control CreateResource for Byteplus::VKE::Cluster can fail with
// InvalidParameter.LoggingConfig in some regions (observed in ap-southeast-3)
// even when Terraform omits logging_config, and also when a disabled/empty
// LoggingConfig is sent. Native VKE CreateCluster succeeds when LoggingConfig
// is omitted.
//
// This override:
//  1. Strips empty / all-disabled LoggingConfig from the Cloud Control TargetState.
//  2. If Cloud Control still fails with InvalidParameter.LoggingConfig, falls
//     back to the native VKE CreateCluster API (same credentials/region).
//
// Update/Delete/Read/Metadata/Schema stay on the auto-generated resource.
// Registered via customresources.Register so cluster_resource_gen.go is untouched.

package vke

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	byteplusvke "github.com/byteplus-sdk/byteplus-go-sdk-v2/service/vke"
	"github.com/byteplus-sdk/terraform-provider-bytepluscc/internal/base"
	"github.com/byteplus-sdk/terraform-provider-bytepluscc/internal/cloudcontrol"
	"github.com/byteplus-sdk/terraform-provider-bytepluscc/internal/customresources"
	"github.com/byteplus-sdk/terraform-provider-bytepluscc/internal/generic"
	tfcloudcontrol "github.com/byteplus-sdk/terraform-provider-bytepluscc/internal/service/cloudcontrol"
	"github.com/byteplus-sdk/terraform-provider-bytepluscc/internal/util"
	hclog "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

const (
	clusterTerraformTypeName    = "bytepluscc_vke_cluster"
	clusterCloudControlTypeName = "Byteplus::VKE::Cluster"
)

// Mirror of WithAttributeNameMap in cluster_resource_gen.go (without synthetic "id").
var clusterTfToCcNameMap = map[string]string{
	"api_server_endpoints":                   "ApiServerEndpoints",
	"api_server_public_access_config":        "ApiServerPublicAccessConfig",
	"api_server_public_access_enabled":       "ApiServerPublicAccessEnabled",
	"audience":                               "Audience",
	"bandwidth":                              "Bandwidth",
	"billing_type":                           "BillingType",
	"cluster_config":                         "ClusterConfig",
	"cluster_id":                             "ClusterId",
	"component_configs":                      "ComponentConfigs",
	"conditions":                             "Conditions",
	"created_time":                           "CreatedTime",
	"creating_count":                         "CreatingCount",
	"delete_protection_enabled":              "DeleteProtectionEnabled",
	"deleting_count":                         "DeletingCount",
	"description":                            "Description",
	"enable_metrics_external_collection":     "EnableMetricsExternalCollection",
	"enabled":                                "Enabled",
	"failed_count":                           "FailedCount",
	"flannel_config":                         "FlannelConfig",
	"ip_family":                              "IpFamily",
	"ipv_4":                                  "Ipv4",
	"irsa_config":                            "IrsaConfig",
	"isp":                                    "Isp",
	"issuer":                                 "Issuer",
	"jwks_url":                               "JwksUrl",
	"key":                                    "Key",
	"kubernetes_version":                     "KubernetesVersion",
	"kubernetes_version_create":              "KubernetesVersionCreate",
	"log_project_id":                         "LogProjectId",
	"log_setups":                             "LogSetups",
	"log_topic_id":                           "LogTopicId",
	"log_ttl":                                "LogTtl",
	"log_type":                               "LogType",
	"logging_config":                         "LoggingConfig",
	"max_pods_per_node":                      "MaxPodsPerNode",
	"message":                                "Message",
	"monitoring_config":                      "MonitoringConfig",
	"name":                                   "Name",
	"node_statistics":                        "NodeStatistics",
	"oidc_trn":                               "OidcTrn",
	"open_id_config_url":                     "OpenIdConfigUrl",
	"phase":                                  "Phase",
	"pod_cidrs":                              "PodCidrs",
	"pod_network_mode":                       "PodNetworkMode",
	"pods_config":                            "PodsConfig",
	"private_ip":                             "PrivateIp",
	"project_name":                           "ProjectName",
	"public_access_network_config":           "PublicAccessNetworkConfig",
	"public_ip":                              "PublicIp",
	"resource_public_access_default_enabled": "ResourcePublicAccessDefaultEnabled",
	"running_count":                          "RunningCount",
	"security_group_ids":                     "SecurityGroupIds",
	"service_cidrsv_4":                       "ServiceCidrsv4",
	"services_config":                        "ServicesConfig",
	"source_region":                          "SourceRegion",
	"status":                                 "Status",
	"subnet_ids":                             "SubnetIds",
	"tags":                                   "Tags",
	"total_count":                            "TotalCount",
	"trunk_eni_enabled":                      "TrunkEniEnabled",
	"type":                                   "Type",
	"updated_time":                           "UpdatedTime",
	"updating_count":                         "UpdatingCount",
	"value":                                  "Value",
	"vpc_cni_config":                         "VpcCniConfig",
	"vpc_id":                                 "VpcId",
	"workspace_id":                           "WorkspaceId",
}

var clusterCcToTfNameMap = func() map[string]string {
	m := make(map[string]string, len(clusterTfToCcNameMap))
	for tf, cc := range clusterTfToCcNameMap {
		m[cc] = tf
	}
	return m
}()

func init() {
	customresources.Register(clusterTerraformTypeName, func(_ context.Context, inner resource.Resource) (resource.Resource, error) {
		return &clusterResourceWithLoggingFix{Resource: inner}, nil
	})
}

type vkeClientProvider interface {
	VKEClient(context.Context) *byteplusvke.VKE
}

type clusterResourceWithLoggingFix struct {
	resource.Resource
	provider tfcloudcontrol.Provider
}

func (r *clusterResourceWithLoggingFix) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if c, ok := r.Resource.(resource.ResourceWithConfigure); ok {
		c.Configure(ctx, req, resp)
	}
	if v := req.ProviderData; v != nil {
		if p, ok := v.(tfcloudcontrol.Provider); ok {
			r.provider = p
		}
	}
}

func (r *clusterResourceWithLoggingFix) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	if i, ok := r.Resource.(resource.ResourceWithImportState); ok {
		i.ImportState(ctx, req, resp)
	}
}

func (r *clusterResourceWithLoggingFix) ConfigValidators(ctx context.Context) []resource.ConfigValidator {
	if c, ok := r.Resource.(resource.ResourceWithConfigValidators); ok {
		return c.ConfigValidators(ctx)
	}
	return nil
}

func (r *clusterResourceWithLoggingFix) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	cloudControlClient := r.provider.CloudControlAPIClient(ctx)

	desiredState, err := generic.ToCloudControlString(ctx, request.Plan.Schema, request.Plan.Raw, clusterTfToCcNameMap)
	if err != nil {
		response.Diagnostics.Append(generic.DesiredStateErrorDiag("Plan", err))
		return
	}

	targetState := make(map[string]any)
	if err := json.Unmarshal([]byte(desiredState), &targetState); err != nil {
		response.Diagnostics.Append(generic.DesiredStateErrorDiag("Plan", err))
		return
	}

	if sanitizeLoggingConfig(targetState) {
		tflog.Info(ctx, "omitting empty/disabled LoggingConfig from VKE cluster Create TargetState")
	}

	output, err := cloudControlClient.CreateResourceWithContext(ctx, &cloudcontrol.CreateResourceInput{
		TypeName:    util.StringPtr(clusterCloudControlTypeName),
		RegionID:    r.provider.Region(ctx),
		ClientToken: util.StringPtr(util.GenerateToken(32)),
		TargetState: &targetState,
	})
	if err != nil {
		response.Diagnostics.Append(generic.ServiceOperationErrorDiag("Cloud Control API", "CreateResource", err))
		return
	}
	if output == nil || output.OperationStatus == nil {
		response.Diagnostics.Append(generic.ServiceOperationEmptyResultDiag("Cloud Control API", "CreateResource"))
		return
	}

	var event *cloudcontrol.ProgressEvent
	switch *output.OperationStatus {
	case base.SUCCESS:
		e := output.ProgressEvent
		event = &e
	case base.IN_PROGRESS, base.PENDING:
		taskId := ""
		if output.TaskID != nil {
			taskId = *output.TaskID
		}
		tflog.Info(ctx, "Cloud Control API CreateResource waiting task ......  ", map[string]interface{}{
			"TaskID":    hclog.Fmt("%v", taskId),
			"RequestID": hclog.Fmt("%v", output.GetRequestId()),
		})
		event, _, err = tfcloudcontrol.AwaitTask(ctx, cloudControlClient, taskId)
		if err != nil {
			if isInvalidLoggingConfigError(err) {
				tflog.Warn(ctx, "Cloud Control VKE cluster create failed on LoggingConfig; falling back to native CreateCluster", map[string]interface{}{
					"error": err.Error(),
				})
				id, fallbackErr := r.createViaNativeVKE(ctx, targetState)
				if fallbackErr != nil {
					response.Diagnostics.Append(generic.ServiceOperationErrorDiag("VKE API", "CreateCluster", fallbackErr))
					return
				}
				r.finishCreateFromID(ctx, request, response, cloudControlClient, id)
				return
			}
			response.Diagnostics.Append(generic.ServiceOperationErrorDiag("Cloud Control API Failed", "GetTask", err))
			return
		}
	case base.FAILED:
		failErr := fmt.Errorf("invoke create handler failed status,resp:%s ", util.JsonString(output))
		if isInvalidLoggingConfigError(failErr) {
			tflog.Warn(ctx, "Cloud Control VKE cluster create failed on LoggingConfig; falling back to native CreateCluster", map[string]interface{}{
				"error": failErr.Error(),
			})
			id, fallbackErr := r.createViaNativeVKE(ctx, targetState)
			if fallbackErr != nil {
				response.Diagnostics.Append(generic.ServiceOperationErrorDiag("VKE API", "CreateCluster", fallbackErr))
				return
			}
			r.finishCreateFromID(ctx, request, response, cloudControlClient, id)
			return
		}
		response.Diagnostics.Append(generic.ServiceOperationErrorDiag("Cloud Control API Failed", "CreateResource", failErr))
		return
	default:
		response.Diagnostics.Append(generic.ServiceOperationErrorDiag("Cloud Control API Failed", "CreateResource",
			fmt.Errorf("invoke create handler other status,resp:%s ", util.JsonString(output))))
		return
	}

	if event == nil || event.Identifier == nil || *event.Identifier == "" {
		response.Diagnostics.Append(generic.ServiceOperationEmptyResultDiag("Cloud Control API", "CreateResource"))
		return
	}
	r.finishCreateFromID(ctx, request, response, cloudControlClient, *event.Identifier)
}

func (r *clusterResourceWithLoggingFix) createViaNativeVKE(ctx context.Context, targetState map[string]any) (string, error) {
	vp, ok := r.provider.(vkeClientProvider)
	if !ok || vp.VKEClient(ctx) == nil {
		return "", fmt.Errorf("native VKE client is not configured on the provider")
	}
	client := vp.VKEClient(ctx)

	body := prepareNativeCreateClusterInput(targetState)
	tflog.Info(ctx, "creating VKE cluster via native CreateCluster API", map[string]interface{}{
		"has_logging_config": body["LoggingConfig"] != nil,
	})

	out, err := client.CreateClusterCommonWithContext(ctx, &body)
	if err != nil {
		return "", err
	}
	id := extractCreateClusterID(out)
	if id == "" {
		return "", fmt.Errorf("CreateCluster succeeded but returned no cluster Id, resp:%s", util.JsonString(out))
	}
	if err := waitClusterPhase(ctx, client, id, "Running", 120*time.Minute); err != nil {
		return "", err
	}
	return id, nil
}

func extractCreateClusterID(out *map[string]interface{}) string {
	if out == nil {
		return ""
	}
	if id, ok := (*out)["Id"].(string); ok && id != "" {
		return id
	}
	if result, ok := (*out)["Result"].(map[string]interface{}); ok {
		if id, ok := result["Id"].(string); ok {
			return id
		}
	}
	return ""
}

func waitClusterPhase(ctx context.Context, client *byteplusvke.VKE, id, wantPhase string, timeout time.Duration) error {
	deadline := time.Now().Add(timeout)
	for {
		if err := ctx.Err(); err != nil {
			return err
		}
		body := map[string]interface{}{
			"Filter": map[string]interface{}{
				"Ids": []string{id},
			},
		}
		out, err := client.ListClustersCommonWithContext(ctx, &body)
		if err != nil {
			return fmt.Errorf("ListClusters while waiting for %s: %w", wantPhase, err)
		}
		phase := clusterPhaseFromList(out, id)
		switch phase {
		case wantPhase:
			return nil
		case "Failed", "Error":
			return fmt.Errorf("cluster %s entered phase %s while waiting for %s", id, phase, wantPhase)
		}
		if time.Now().After(deadline) {
			return fmt.Errorf("timed out waiting for cluster %s phase %s (last phase=%q)", id, wantPhase, phase)
		}
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(15 * time.Second):
		}
	}
}

func clusterPhaseFromList(out *map[string]interface{}, id string) string {
	if out == nil {
		return ""
	}
	items := clusterItemsFromList(out)
	for _, item := range items {
		itemID, _ := item["Id"].(string)
		if itemID != id && id != "" {
			// Some responses may omit Id match quirks; still accept single item.
			if len(items) != 1 {
				continue
			}
		}
		if status, ok := item["Status"].(map[string]interface{}); ok {
			if phase, ok := status["Phase"].(string); ok {
				return phase
			}
		}
	}
	return ""
}

func clusterItemsFromList(out *map[string]interface{}) []map[string]interface{} {
	var rawItems any
	if items, ok := (*out)["Items"]; ok {
		rawItems = items
	} else if result, ok := (*out)["Result"].(map[string]interface{}); ok {
		rawItems = result["Items"]
	}
	list, ok := rawItems.([]any)
	if !ok {
		return nil
	}
	outItems := make([]map[string]interface{}, 0, len(list))
	for _, item := range list {
		if m, ok := item.(map[string]interface{}); ok {
			outItems = append(outItems, m)
		}
	}
	return outItems
}

func (r *clusterResourceWithLoggingFix) finishCreateFromID(
	ctx context.Context,
	request resource.CreateRequest,
	response *resource.CreateResponse,
	cloudControlClient *cloudcontrol.CloudControl,
	id string,
) {
	response.State.Raw = request.Plan.Raw
	if diags := response.State.SetAttribute(ctx, path.Root("id"), id); diags.HasError() {
		response.Diagnostics.Append(diags...)
		return
	}
	// Prefer ClusterId attribute when present in schema.
	_ = response.State.SetAttribute(ctx, path.Root("cluster_id"), id)

	unknowns, err := generic.UnknownValuePaths(ctx, response.State.Raw)
	if err != nil {
		response.Diagnostics.AddError("Creation Of Terraform State Unsuccessful", err.Error())
		return
	}
	if len(unknowns) == 0 {
		return
	}

	description, err := tfcloudcontrol.FindResourceByTypeNameAndID(ctx, cloudControlClient, r.provider.Region(ctx), clusterCloudControlTypeName, id)
	if err != nil {
		// Native create may race CC GetResource briefly; retry a few times.
		var lastErr error = err
		for i := 0; i < 12; i++ {
			time.Sleep(10 * time.Second)
			description, err = tfcloudcontrol.FindResourceByTypeNameAndID(ctx, cloudControlClient, r.provider.Region(ctx), clusterCloudControlTypeName, id)
			if err == nil {
				lastErr = nil
				break
			}
			lastErr = err
			if !strings.Contains(err.Error(), "NotFound") && !strings.Contains(err.Error(), "not found") {
				break
			}
		}
		if lastErr != nil {
			response.Diagnostics.Append(generic.ServiceOperationErrorDiag("Cloud Control API", "GetResource", lastErr))
			return
		}
	}
	if description == nil {
		response.Diagnostics.Append(generic.ServiceOperationEmptyResultDiag("Cloud Control API", "GetResource"))
		return
	}
	resourceModel := util.ToString(description.ResourceDescription.Properties)
	if err := generic.SetUnknownValuesFromResourceModel(ctx, &response.State, unknowns, resourceModel, clusterCcToTfNameMap); err != nil {
		response.Diagnostics.AddError("Creation Of Terraform State Unsuccessful", err.Error())
		return
	}
}
