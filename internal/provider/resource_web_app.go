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

// Web Apps (Linux and Windows) share most of their schema. We model them as
// two separate resources to mirror the AzureRM split, with a common Go struct.

type WebAppModel struct {
	ID                  types.String `tfsdk:"id"`
	Name                types.String `tfsdk:"name"`
	ResourceGroupName   types.String `tfsdk:"resource_group_name"`
	Location            types.String `tfsdk:"location"`
	ServicePlanID       types.String `tfsdk:"service_plan_id"`
	HTTPSOnly           types.Bool   `tfsdk:"https_only"`
	ClientAffinityEnabled types.Bool `tfsdk:"client_affinity_enabled"`
	AppSettings         types.Map    `tfsdk:"app_settings"`
	SiteConfig          types.List   `tfsdk:"site_config"`
	DefaultHostname     types.String `tfsdk:"default_hostname"`
	OutboundIPAddresses types.String `tfsdk:"outbound_ip_addresses"`
	Tags                types.Map    `tfsdk:"tags"`
}

func webAppAttributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Computed: true, Description: "The Web App ID.",
			PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
		},
		"name": schema.StringAttribute{
			Required: true, Description: "Web App name (must be globally unique in real Azure).",
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
		"service_plan_id": schema.StringAttribute{
			Required: true, Description: "Parent App Service Plan ID.",
		},
		"https_only":              schema.BoolAttribute{Optional: true, Description: "Force HTTPS."},
		"client_affinity_enabled": schema.BoolAttribute{Optional: true, Description: "Enable session affinity."},
		"app_settings": schema.MapAttribute{
			Optional: true, ElementType: types.StringType, Description: "Environment variables.",
		},
		"default_hostname": schema.StringAttribute{
			Computed: true, Description: "Simulated default hostname (`<name>.azurewebsites.net`).",
			PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
		},
		"outbound_ip_addresses": schema.StringAttribute{
			Computed: true, Description: "Simulated outbound IPs (comma-separated).",
			PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
		},
		"tags": schema.MapAttribute{
			Optional: true, ElementType: types.StringType, Description: "Tags.",
		},
	}
}

func webAppBlocks() map[string]schema.Block {
	return map[string]schema.Block{
		"site_config": schema.ListNestedBlock{
			Description: "Site configuration.",
			NestedObject: schema.NestedBlockObject{
				Attributes: map[string]schema.Attribute{
					"always_on":                schema.BoolAttribute{Optional: true, Description: "Keep the app always on."},
					"ftps_state":               schema.StringAttribute{Optional: true, Description: "`AllAllowed`, `FtpsOnly`, or `Disabled`."},
					"http2_enabled":            schema.BoolAttribute{Optional: true, Description: "Enable HTTP/2."},
					"minimum_tls_version":      schema.StringAttribute{Optional: true, Description: "`1.0`, `1.1`, `1.2`, or `1.3`."},
					"websockets_enabled":       schema.BoolAttribute{Optional: true, Description: "Enable WebSockets."},
					"health_check_path":        schema.StringAttribute{Optional: true, Description: "Health check path."},
					"vnet_route_all_enabled":   schema.BoolAttribute{Optional: true, Description: "Route all outbound through the VNet."},
					"linux_fx_version":         schema.StringAttribute{Optional: true, Description: "Linux runtime stack (e.g. `NODE|18-lts`)."},
					"windows_fx_version":       schema.StringAttribute{Optional: true, Description: "Windows runtime stack."},
					"app_command_line":         schema.StringAttribute{Optional: true, Description: "Startup command."},
				},
			},
		},
	}
}

func webAppApplyComputed(plan *WebAppModel, subscriptionID string) {
	name := plan.Name.ValueString()
	plan.ID = types.StringValue(fmt.Sprintf(
		"/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Web/sites/%s",
		subscriptionID, plan.ResourceGroupName.ValueString(), name,
	))
	plan.DefaultHostname = types.StringValue(fmt.Sprintf("%s.azurewebsites.net", name))
	plan.OutboundIPAddresses = types.StringValue("203.0.113.10,203.0.113.11,203.0.113.12,203.0.113.13")
}

// --- Linux Web App ---

var _ resource.Resource = &LinuxWebAppResource{}
var _ resource.ResourceWithConfigure = &LinuxWebAppResource{}

type LinuxWebAppResource struct{ subscriptionID string }

func NewLinuxWebAppResource() resource.Resource { return &LinuxWebAppResource{} }

func (r *LinuxWebAppResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_linux_web_app"
}

func (r *LinuxWebAppResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Simulates an Azure Linux Web App.",
		Attributes:  webAppAttributes(),
		Blocks:      webAppBlocks(),
	}
}

func (r *LinuxWebAppResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	r.subscriptionID = req.ProviderData.(string)
}

func (r *LinuxWebAppResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan WebAppModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	webAppApplyComputed(&plan, r.subscriptionID)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *LinuxWebAppResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state WebAppModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *LinuxWebAppResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan WebAppModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	webAppApplyComputed(&plan, r.subscriptionID)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *LinuxWebAppResource) Delete(_ context.Context, _ resource.DeleteRequest, _ *resource.DeleteResponse) {
}

// --- Windows Web App ---

var _ resource.Resource = &WindowsWebAppResource{}
var _ resource.ResourceWithConfigure = &WindowsWebAppResource{}

type WindowsWebAppResource struct{ subscriptionID string }

func NewWindowsWebAppResource() resource.Resource { return &WindowsWebAppResource{} }

func (r *WindowsWebAppResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_windows_web_app"
}

func (r *WindowsWebAppResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Simulates an Azure Windows Web App.",
		Attributes:  webAppAttributes(),
		Blocks:      webAppBlocks(),
	}
}

func (r *WindowsWebAppResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	r.subscriptionID = req.ProviderData.(string)
}

func (r *WindowsWebAppResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan WebAppModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	webAppApplyComputed(&plan, r.subscriptionID)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *WindowsWebAppResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state WebAppModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *WindowsWebAppResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan WebAppModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	webAppApplyComputed(&plan, r.subscriptionID)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *WindowsWebAppResource) Delete(_ context.Context, _ resource.DeleteRequest, _ *resource.DeleteResponse) {
}
