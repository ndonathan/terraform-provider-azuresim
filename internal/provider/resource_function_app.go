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

type FunctionAppModel struct {
	ID                       types.String `tfsdk:"id"`
	Name                     types.String `tfsdk:"name"`
	ResourceGroupName        types.String `tfsdk:"resource_group_name"`
	Location                 types.String `tfsdk:"location"`
	ServicePlanID            types.String `tfsdk:"service_plan_id"`
	StorageAccountName       types.String `tfsdk:"storage_account_name"`
	StorageAccountAccessKey  types.String `tfsdk:"storage_account_access_key"`
	StorageUsesMSI           types.Bool   `tfsdk:"storage_uses_managed_identity"`
	HTTPSOnly                types.Bool   `tfsdk:"https_only"`
	FunctionsExtensionVersion types.String `tfsdk:"functions_extension_version"`
	BuiltinLoggingEnabled    types.Bool   `tfsdk:"builtin_logging_enabled"`
	AppSettings              types.Map    `tfsdk:"app_settings"`
	DefaultHostname          types.String `tfsdk:"default_hostname"`
	OutboundIPAddresses      types.String `tfsdk:"outbound_ip_addresses"`
	Tags                     types.Map    `tfsdk:"tags"`
}

func functionAppAttributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Computed: true, Description: "The Function App ID.",
			PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
		},
		"name": schema.StringAttribute{
			Required: true, Description: "Function App name (must be globally unique in real Azure).",
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
		"service_plan_id":               schema.StringAttribute{Required: true, Description: "Parent App Service Plan ID."},
		"storage_account_name":          schema.StringAttribute{Required: true, Description: "Backing Storage Account name."},
		"storage_account_access_key":    schema.StringAttribute{Optional: true, Sensitive: true, Description: "Backing Storage Account access key. Mutually exclusive with `storage_uses_managed_identity = true`."},
		"storage_uses_managed_identity": schema.BoolAttribute{Optional: true, Description: "Use the system-assigned identity to access storage."},
		"https_only":                    schema.BoolAttribute{Optional: true, Description: "Force HTTPS."},
		"functions_extension_version":   schema.StringAttribute{Optional: true, Description: "Functions runtime version (e.g. `~4`)."},
		"builtin_logging_enabled":       schema.BoolAttribute{Optional: true, Description: "Enable Application Insights built-in logging."},
		"app_settings": schema.MapAttribute{
			Optional: true, ElementType: types.StringType, Description: "Environment variables.",
		},
		"default_hostname": schema.StringAttribute{
			Computed: true, Description: "Simulated default hostname.",
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

func functionAppApplyComputed(plan *FunctionAppModel, subscriptionID string) {
	name := plan.Name.ValueString()
	plan.ID = types.StringValue(fmt.Sprintf(
		"/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Web/sites/%s",
		subscriptionID, plan.ResourceGroupName.ValueString(), name,
	))
	plan.DefaultHostname = types.StringValue(fmt.Sprintf("%s.azurewebsites.net", name))
	plan.OutboundIPAddresses = types.StringValue("203.0.113.20,203.0.113.21,203.0.113.22,203.0.113.23")
}

// --- Linux Function App ---

var _ resource.Resource = &LinuxFunctionAppResource{}
var _ resource.ResourceWithConfigure = &LinuxFunctionAppResource{}

type LinuxFunctionAppResource struct{ subscriptionID string }

func NewLinuxFunctionAppResource() resource.Resource { return &LinuxFunctionAppResource{} }

func (r *LinuxFunctionAppResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_linux_function_app"
}

func (r *LinuxFunctionAppResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{Description: "Simulates an Azure Linux Function App.", Attributes: functionAppAttributes()}
}

func (r *LinuxFunctionAppResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	r.subscriptionID = req.ProviderData.(string)
}

func (r *LinuxFunctionAppResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan FunctionAppModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	functionAppApplyComputed(&plan, r.subscriptionID)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *LinuxFunctionAppResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state FunctionAppModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *LinuxFunctionAppResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan FunctionAppModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	functionAppApplyComputed(&plan, r.subscriptionID)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *LinuxFunctionAppResource) Delete(_ context.Context, _ resource.DeleteRequest, _ *resource.DeleteResponse) {
}

// --- Windows Function App ---

var _ resource.Resource = &WindowsFunctionAppResource{}
var _ resource.ResourceWithConfigure = &WindowsFunctionAppResource{}

type WindowsFunctionAppResource struct{ subscriptionID string }

func NewWindowsFunctionAppResource() resource.Resource { return &WindowsFunctionAppResource{} }

func (r *WindowsFunctionAppResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_windows_function_app"
}

func (r *WindowsFunctionAppResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{Description: "Simulates an Azure Windows Function App.", Attributes: functionAppAttributes()}
}

func (r *WindowsFunctionAppResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	r.subscriptionID = req.ProviderData.(string)
}

func (r *WindowsFunctionAppResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan FunctionAppModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	functionAppApplyComputed(&plan, r.subscriptionID)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *WindowsFunctionAppResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state FunctionAppModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *WindowsFunctionAppResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan FunctionAppModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	functionAppApplyComputed(&plan, r.subscriptionID)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *WindowsFunctionAppResource) Delete(_ context.Context, _ resource.DeleteRequest, _ *resource.DeleteResponse) {
}
