package provider

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// --- MSSQL Server ---

var _ resource.Resource = &MSSQLServerResource{}
var _ resource.ResourceWithConfigure = &MSSQLServerResource{}

type MSSQLServerResource struct {
	subscriptionID string
}

type MSSQLServerModel struct {
	ID                          types.String `tfsdk:"id"`
	Name                        types.String `tfsdk:"name"`
	ResourceGroupName           types.String `tfsdk:"resource_group_name"`
	Location                    types.String `tfsdk:"location"`
	Version                     types.String `tfsdk:"version"`
	AdministratorLogin          types.String `tfsdk:"administrator_login"`
	AdministratorLoginPassword  types.String `tfsdk:"administrator_login_password"`
	MinimumTLSVersion           types.String `tfsdk:"minimum_tls_version"`
	PublicNetworkAccessEnabled  types.Bool   `tfsdk:"public_network_access_enabled"`
	OutboundNetworkRestrictionEnabled types.Bool `tfsdk:"outbound_network_restriction_enabled"`
	FullyQualifiedDomainName    types.String `tfsdk:"fully_qualified_domain_name"`
	Tags                        types.Map    `tfsdk:"tags"`
}

func NewMSSQLServerResource() resource.Resource {
	return &MSSQLServerResource{}
}

func (r *MSSQLServerResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_mssql_server"
}

func (r *MSSQLServerResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Simulates an Azure SQL Server.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true, Description: "The SQL Server ID.",
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"name": schema.StringAttribute{
				Required: true, Description: "Server name (must be globally unique).",
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
			"version": schema.StringAttribute{
				Required: true, Description: "`2.0` (v11) or `12.0` (v12).",
			},
			"administrator_login": schema.StringAttribute{
				Optional: true, Description: "Server admin login (required if not using AAD-only auth).",
			},
			"administrator_login_password": schema.StringAttribute{
				Optional: true, Sensitive: true, Description: "Server admin password.",
			},
			"minimum_tls_version":                  schema.StringAttribute{Optional: true, Description: "`1.0`, `1.1`, `1.2`, or `1.3`."},
			"public_network_access_enabled":        schema.BoolAttribute{Optional: true, Description: "Whether public network access is enabled."},
			"outbound_network_restriction_enabled": schema.BoolAttribute{Optional: true, Description: "Restrict outbound network traffic."},
			"fully_qualified_domain_name": schema.StringAttribute{
				Computed: true, Description: "Simulated FQDN: `<name>.database.windows.net`.",
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"tags": schema.MapAttribute{
				Optional: true, ElementType: types.StringType, Description: "Tags.",
			},
		},
	}
}

func (r *MSSQLServerResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	r.subscriptionID = req.ProviderData.(string)
}

func (r *MSSQLServerResource) applyComputed(plan *MSSQLServerModel) {
	name := plan.Name.ValueString()
	plan.ID = types.StringValue(fmt.Sprintf(
		"/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Sql/servers/%s",
		r.subscriptionID, plan.ResourceGroupName.ValueString(), name,
	))
	plan.FullyQualifiedDomainName = types.StringValue(fmt.Sprintf("%s.database.windows.net", name))
}

func (r *MSSQLServerResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan MSSQLServerModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	r.applyComputed(&plan)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *MSSQLServerResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state MSSQLServerModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *MSSQLServerResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan MSSQLServerModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	r.applyComputed(&plan)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *MSSQLServerResource) Delete(_ context.Context, _ resource.DeleteRequest, _ *resource.DeleteResponse) {
}

// --- MSSQL Database ---

var _ resource.Resource = &MSSQLDatabaseResource{}
var _ resource.ResourceWithConfigure = &MSSQLDatabaseResource{}

type MSSQLDatabaseResource struct {
	subscriptionID string
}

type MSSQLDatabaseModel struct {
	ID                  types.String `tfsdk:"id"`
	Name                types.String `tfsdk:"name"`
	ServerID            types.String `tfsdk:"server_id"`
	SKUName             types.String `tfsdk:"sku_name"`
	Collation           types.String `tfsdk:"collation"`
	MaxSizeGB           types.Int64  `tfsdk:"max_size_gb"`
	ZoneRedundant       types.Bool   `tfsdk:"zone_redundant"`
	GeoBackupEnabled    types.Bool   `tfsdk:"geo_backup_enabled"`
	StorageAccountType  types.String `tfsdk:"storage_account_type"`
	ReadScale           types.Bool   `tfsdk:"read_scale"`
	LicenseType         types.String `tfsdk:"license_type"`
	Tags                types.Map    `tfsdk:"tags"`
}

func NewMSSQLDatabaseResource() resource.Resource {
	return &MSSQLDatabaseResource{}
}

func (r *MSSQLDatabaseResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_mssql_database"
}

func (r *MSSQLDatabaseResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Simulates an Azure SQL Database.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true, Description: "The SQL Database ID.",
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"name": schema.StringAttribute{
				Required: true, Description: "Database name.",
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"server_id": schema.StringAttribute{
				Required: true, Description: "Parent SQL Server ID.",
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"sku_name":             schema.StringAttribute{Optional: true, Description: "SKU (e.g. `Basic`, `S0`, `P1`, `GP_Gen5_2`, `BC_Gen5_4`)."},
			"collation":            schema.StringAttribute{Optional: true, Description: "Collation (e.g. `SQL_Latin1_General_CP1_CI_AS`)."},
			"max_size_gb":          schema.Int64Attribute{Optional: true, Description: "Max database size in GB."},
			"zone_redundant":       schema.BoolAttribute{Optional: true, Description: "Enable zone redundancy."},
			"geo_backup_enabled":   schema.BoolAttribute{Optional: true, Description: "Enable geo-redundant backups."},
			"storage_account_type": schema.StringAttribute{Optional: true, Description: "`Geo`, `GeoZone`, `Local`, or `Zone`."},
			"read_scale":           schema.BoolAttribute{Optional: true, Description: "Enable read-scale (Premium/Business Critical)."},
			"license_type":         schema.StringAttribute{Optional: true, Description: "`LicenseIncluded` or `BasePrice`."},
			"tags": schema.MapAttribute{
				Optional: true, ElementType: types.StringType, Description: "Tags.",
			},
		},
	}
}

func (r *MSSQLDatabaseResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	r.subscriptionID = req.ProviderData.(string)
}

func (r *MSSQLDatabaseResource) applyComputed(plan *MSSQLDatabaseModel) {
	serverID := strings.TrimRight(plan.ServerID.ValueString(), "/")
	plan.ID = types.StringValue(fmt.Sprintf("%s/databases/%s", serverID, plan.Name.ValueString()))
}

func (r *MSSQLDatabaseResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan MSSQLDatabaseModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	r.applyComputed(&plan)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *MSSQLDatabaseResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state MSSQLDatabaseModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *MSSQLDatabaseResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan MSSQLDatabaseModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	r.applyComputed(&plan)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *MSSQLDatabaseResource) Delete(_ context.Context, _ resource.DeleteRequest, _ *resource.DeleteResponse) {
}
