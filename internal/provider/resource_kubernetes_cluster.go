package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ resource.Resource = &KubernetesClusterResource{}
var _ resource.ResourceWithConfigure = &KubernetesClusterResource{}

type KubernetesClusterResource struct {
	subscriptionID string
}

type KubernetesClusterModel struct {
	ID                       types.String `tfsdk:"id"`
	Name                     types.String `tfsdk:"name"`
	ResourceGroupName        types.String `tfsdk:"resource_group_name"`
	Location                 types.String `tfsdk:"location"`
	DNSPrefix                types.String `tfsdk:"dns_prefix"`
	KubernetesVersion        types.String `tfsdk:"kubernetes_version"`
	SKUTier                  types.String `tfsdk:"sku_tier"`
	PrivateClusterEnabled    types.Bool   `tfsdk:"private_cluster_enabled"`
	NodeResourceGroup        types.String `tfsdk:"node_resource_group"`
	FQDN                     types.String `tfsdk:"fqdn"`
	PrivateFQDN              types.String `tfsdk:"private_fqdn"`
	DefaultNodePool          types.List   `tfsdk:"default_node_pool"`
	Identity                 types.List   `tfsdk:"identity"`
	NetworkProfile           types.List   `tfsdk:"network_profile"`
	Tags                     types.Map    `tfsdk:"tags"`
}

func NewKubernetesClusterResource() resource.Resource { return &KubernetesClusterResource{} }

func (r *KubernetesClusterResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_kubernetes_cluster"
}

func (r *KubernetesClusterResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Simulates an Azure Kubernetes Service (AKS) cluster.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true, Description: "The AKS Cluster ID.",
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"name": schema.StringAttribute{
				Required: true, Description: "Cluster name.",
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"resource_group_name": schema.StringAttribute{
				Required: true, Description: "Resource Group.",
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"location": schema.StringAttribute{
				Required: true, Description: "Azure region.",
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"dns_prefix": schema.StringAttribute{
				Required: true, Description: "DNS prefix.",
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"kubernetes_version":      schema.StringAttribute{Optional: true, Description: "Kubernetes version (e.g. `1.30.0`)."},
			"sku_tier":                schema.StringAttribute{Optional: true, Description: "`Free`, `Standard`, or `Premium`."},
			"private_cluster_enabled": schema.BoolAttribute{Optional: true, Description: "Enable private cluster."},
			"node_resource_group": schema.StringAttribute{
				Computed: true, Description: "Node resource group (`MC_<rg>_<name>_<location>`).",
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"fqdn": schema.StringAttribute{
				Computed: true, Description: "Public API server FQDN.",
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"private_fqdn": schema.StringAttribute{
				Computed: true, Description: "Private API server FQDN (when `private_cluster_enabled = true`).",
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"tags": schema.MapAttribute{Optional: true, ElementType: types.StringType, Description: "Tags."},
		},
		Blocks: map[string]schema.Block{
			"default_node_pool": schema.ListNestedBlock{
				Description: "Default system node pool.",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"name":           schema.StringAttribute{Required: true, Description: "Node pool name."},
						"vm_size":        schema.StringAttribute{Required: true, Description: "VM size (e.g. `Standard_D2s_v5`)."},
						"node_count":     schema.Int64Attribute{Optional: true, Description: "Initial node count."},
						"min_count":      schema.Int64Attribute{Optional: true, Description: "Min nodes for autoscale."},
						"max_count":      schema.Int64Attribute{Optional: true, Description: "Max nodes for autoscale."},
						"vnet_subnet_id": schema.StringAttribute{Optional: true, Description: "Subnet ID."},
						"os_disk_size_gb": schema.Int64Attribute{Optional: true, Description: "OS disk size in GB."},
						"os_disk_type":   schema.StringAttribute{Optional: true, Description: "`Managed` or `Ephemeral`."},
						"zones":          schema.ListAttribute{Optional: true, ElementType: types.StringType, Description: "Availability zones."},
					},
				},
			},
			"identity": schema.ListNestedBlock{
				Description: "Cluster managed identity.",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"type":         schema.StringAttribute{Required: true, Description: "`SystemAssigned` or `UserAssigned`."},
						"identity_ids": schema.ListAttribute{Optional: true, ElementType: types.StringType, Description: "User-assigned identity IDs."},
					},
				},
			},
			"network_profile": schema.ListNestedBlock{
				Description: "Network profile.",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"network_plugin":     schema.StringAttribute{Required: true, Description: "`azure`, `kubenet`, or `none`."},
						"network_policy":     schema.StringAttribute{Optional: true, Description: "`calico`, `azure`, or `cilium`."},
						"service_cidr":       schema.StringAttribute{Optional: true, Description: "Service CIDR."},
						"dns_service_ip":     schema.StringAttribute{Optional: true, Description: "DNS service IP."},
						"load_balancer_sku":  schema.StringAttribute{Optional: true, Description: "`basic` or `standard`."},
					},
				},
			},
		},
	}
}

func (r *KubernetesClusterResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	r.subscriptionID = req.ProviderData.(string)
}

func (r *KubernetesClusterResource) applyComputed(plan *KubernetesClusterModel) {
	name := plan.Name.ValueString()
	rg := plan.ResourceGroupName.ValueString()
	loc := plan.Location.ValueString()
	dns := plan.DNSPrefix.ValueString()

	plan.ID = types.StringValue(fmt.Sprintf(
		"/subscriptions/%s/resourceGroups/%s/providers/Microsoft.ContainerService/managedClusters/%s",
		r.subscriptionID, rg, name,
	))
	plan.NodeResourceGroup = types.StringValue(fmt.Sprintf("MC_%s_%s_%s", rg, name, loc))
	plan.FQDN = types.StringValue(fmt.Sprintf("%s-deadbeef.hcp.%s.azmk8s.io", dns, loc))
	plan.PrivateFQDN = types.StringValue(fmt.Sprintf("%s-deadbeef.privatelink.%s.azmk8s.io", dns, loc))
}

func (r *KubernetesClusterResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan KubernetesClusterModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	r.applyComputed(&plan)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *KubernetesClusterResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state KubernetesClusterModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *KubernetesClusterResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan KubernetesClusterModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	r.applyComputed(&plan)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *KubernetesClusterResource) Delete(_ context.Context, _ resource.DeleteRequest, _ *resource.DeleteResponse) {
}
