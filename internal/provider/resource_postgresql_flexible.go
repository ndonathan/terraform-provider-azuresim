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

var _ resource.Resource = &PostgreSQLFlexibleServerResource{}
var _ resource.ResourceWithConfigure = &PostgreSQLFlexibleServerResource{}

type PostgreSQLFlexibleServerResource struct{ subscriptionID string }

type PostgreSQLFlexibleServerModel struct {
	ID                          types.String `tfsdk:"id"`
	Name                        types.String `tfsdk:"name"`
	ResourceGroupName           types.String `tfsdk:"resource_group_name"`
	Location                    types.String `tfsdk:"location"`
	Version                     types.String `tfsdk:"version"`
	SKUName                     types.String `tfsdk:"sku_name"`
	StorageMB                   types.Int64  `tfsdk:"storage_mb"`
	StorageTier                 types.String `tfsdk:"storage_tier"`
	BackupRetentionDays         types.Int64  `tfsdk:"backup_retention_days"`
	GeoRedundantBackupEnabled   types.Bool   `tfsdk:"geo_redundant_backup_enabled"`
	AdministratorLogin          types.String `tfsdk:"administrator_login"`
	AdministratorPassword       types.String `tfsdk:"administrator_password"`
	DelegatedSubnetID           types.String `tfsdk:"delegated_subnet_id"`
	PrivateDNSZoneID            types.String `tfsdk:"private_dns_zone_id"`
	Zone                        types.String `tfsdk:"zone"`
	HighAvailability            types.List   `tfsdk:"high_availability"`
	PublicNetworkAccessEnabled  types.Bool   `tfsdk:"public_network_access_enabled"`
	FQDN                        types.String `tfsdk:"fqdn"`
	Tags                        types.Map    `tfsdk:"tags"`
}

func NewPostgreSQLFlexibleServerResource() resource.Resource { return &PostgreSQLFlexibleServerResource{} }

func (r *PostgreSQLFlexibleServerResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_postgresql_flexible_server"
}

func (r *PostgreSQLFlexibleServerResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Simulates an Azure PostgreSQL Flexible Server.",
		Attributes: map[string]schema.Attribute{
			"id":   schema.StringAttribute{Computed: true, PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()}, Description: "Server ID."},
			"name": schema.StringAttribute{Required: true, PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()}, Description: "Server name (globally unique)."},
			"resource_group_name": schema.StringAttribute{Required: true, PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()}, Description: "Resource Group."},
			"location":              schema.StringAttribute{Required: true, PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()}, Description: "Azure region."},
			"version":               schema.StringAttribute{Optional: true, Description: "Postgres version (e.g. `14`, `15`, `16`)."},
			"sku_name":              schema.StringAttribute{Required: true, Description: "SKU (e.g. `B_Standard_B1ms`, `GP_Standard_D2s_v3`)."},
			"storage_mb":            schema.Int64Attribute{Optional: true, Description: "Storage in MB."},
			"storage_tier":          schema.StringAttribute{Optional: true, Description: "Storage tier (`P4`, `P6`, `P10`, ...)."},
			"backup_retention_days": schema.Int64Attribute{Optional: true, Description: "Backup retention (7-35)."},
			"geo_redundant_backup_enabled": schema.BoolAttribute{Optional: true, Description: "Enable geo-redundant backup."},
			"administrator_login":          schema.StringAttribute{Optional: true, Description: "Admin login."},
			"administrator_password":       schema.StringAttribute{Optional: true, Sensitive: true, Description: "Admin password."},
			"delegated_subnet_id":          schema.StringAttribute{Optional: true, Description: "Delegated subnet for VNet integration."},
			"private_dns_zone_id":          schema.StringAttribute{Optional: true, Description: "Private DNS zone ID."},
			"zone":                         schema.StringAttribute{Optional: true, Description: "Availability zone."},
			"public_network_access_enabled": schema.BoolAttribute{Optional: true, Description: "Allow public access."},
			"fqdn": schema.StringAttribute{
				Computed: true, Description: "Simulated FQDN: `<name>.postgres.database.azure.com`.",
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"tags": schema.MapAttribute{Optional: true, ElementType: types.StringType, Description: "Tags."},
		},
		Blocks: map[string]schema.Block{
			"high_availability": schema.ListNestedBlock{
				Description: "High-availability config.",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"mode":                      schema.StringAttribute{Required: true, Description: "`SameZone` or `ZoneRedundant`."},
						"standby_availability_zone": schema.StringAttribute{Optional: true, Description: "Standby zone."},
					},
				},
			},
		},
	}
}

func (r *PostgreSQLFlexibleServerResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	r.subscriptionID = req.ProviderData.(string)
}

func (r *PostgreSQLFlexibleServerResource) applyComputed(plan *PostgreSQLFlexibleServerModel) {
	name := plan.Name.ValueString()
	plan.ID = types.StringValue(fmt.Sprintf(
		"/subscriptions/%s/resourceGroups/%s/providers/Microsoft.DBforPostgreSQL/flexibleServers/%s",
		r.subscriptionID, plan.ResourceGroupName.ValueString(), name,
	))
	plan.FQDN = types.StringValue(fmt.Sprintf("%s.postgres.database.azure.com", name))
}

func (r *PostgreSQLFlexibleServerResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan PostgreSQLFlexibleServerModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	r.applyComputed(&plan)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *PostgreSQLFlexibleServerResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state PostgreSQLFlexibleServerModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *PostgreSQLFlexibleServerResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan PostgreSQLFlexibleServerModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	r.applyComputed(&plan)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *PostgreSQLFlexibleServerResource) Delete(_ context.Context, _ resource.DeleteRequest, _ *resource.DeleteResponse) {
}

// --- MySQL Flexible Server (mirrors PostgreSQL) ---

var _ resource.Resource = &MySQLFlexibleServerResource{}
var _ resource.ResourceWithConfigure = &MySQLFlexibleServerResource{}

type MySQLFlexibleServerResource struct{ subscriptionID string }

type MySQLFlexibleServerModel struct {
	ID                          types.String `tfsdk:"id"`
	Name                        types.String `tfsdk:"name"`
	ResourceGroupName           types.String `tfsdk:"resource_group_name"`
	Location                    types.String `tfsdk:"location"`
	Version                     types.String `tfsdk:"version"`
	SKUName                     types.String `tfsdk:"sku_name"`
	BackupRetentionDays         types.Int64  `tfsdk:"backup_retention_days"`
	GeoRedundantBackupEnabled   types.Bool   `tfsdk:"geo_redundant_backup_enabled"`
	AdministratorLogin          types.String `tfsdk:"administrator_login"`
	AdministratorPassword       types.String `tfsdk:"administrator_password"`
	DelegatedSubnetID           types.String `tfsdk:"delegated_subnet_id"`
	PrivateDNSZoneID            types.String `tfsdk:"private_dns_zone_id"`
	Zone                        types.String `tfsdk:"zone"`
	HighAvailability            types.List   `tfsdk:"high_availability"`
	Storage                     types.List   `tfsdk:"storage"`
	FQDN                        types.String `tfsdk:"fqdn"`
	Tags                        types.Map    `tfsdk:"tags"`
}

func NewMySQLFlexibleServerResource() resource.Resource { return &MySQLFlexibleServerResource{} }

func (r *MySQLFlexibleServerResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_mysql_flexible_server"
}

func (r *MySQLFlexibleServerResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Simulates an Azure MySQL Flexible Server.",
		Attributes: map[string]schema.Attribute{
			"id":   schema.StringAttribute{Computed: true, PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()}, Description: "Server ID."},
			"name": schema.StringAttribute{Required: true, PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()}, Description: "Server name (globally unique)."},
			"resource_group_name": schema.StringAttribute{Required: true, PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()}, Description: "Resource Group."},
			"location":            schema.StringAttribute{Required: true, PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()}, Description: "Azure region."},
			"version":             schema.StringAttribute{Optional: true, Description: "MySQL version (e.g. `5.7`, `8.0.21`)."},
			"sku_name":            schema.StringAttribute{Required: true, Description: "SKU."},
			"backup_retention_days": schema.Int64Attribute{Optional: true, Description: "Backup retention (1-35)."},
			"geo_redundant_backup_enabled": schema.BoolAttribute{Optional: true, Description: "Enable geo-redundant backup."},
			"administrator_login":          schema.StringAttribute{Optional: true, Description: "Admin login."},
			"administrator_password":       schema.StringAttribute{Optional: true, Sensitive: true, Description: "Admin password."},
			"delegated_subnet_id":          schema.StringAttribute{Optional: true, Description: "Delegated subnet."},
			"private_dns_zone_id":          schema.StringAttribute{Optional: true, Description: "Private DNS zone ID."},
			"zone":                         schema.StringAttribute{Optional: true, Description: "Availability zone."},
			"fqdn": schema.StringAttribute{
				Computed: true, Description: "Simulated FQDN: `<name>.mysql.database.azure.com`.",
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"tags": schema.MapAttribute{Optional: true, ElementType: types.StringType, Description: "Tags."},
		},
		Blocks: map[string]schema.Block{
			"high_availability": schema.ListNestedBlock{
				Description: "High-availability config.",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"mode":                      schema.StringAttribute{Required: true, Description: "`SameZone` or `ZoneRedundant`."},
						"standby_availability_zone": schema.StringAttribute{Optional: true, Description: "Standby zone."},
					},
				},
			},
			"storage": schema.ListNestedBlock{
				Description: "Storage configuration.",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"size_gb":            schema.Int64Attribute{Optional: true, Description: "Storage size in GB."},
						"iops":               schema.Int64Attribute{Optional: true, Description: "Storage IOPS."},
						"auto_grow_enabled":  schema.BoolAttribute{Optional: true, Description: "Enable auto-grow."},
						"io_scaling_enabled": schema.BoolAttribute{Optional: true, Description: "Enable IO scaling."},
					},
				},
			},
		},
	}
}

func (r *MySQLFlexibleServerResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	r.subscriptionID = req.ProviderData.(string)
}

func (r *MySQLFlexibleServerResource) applyComputed(plan *MySQLFlexibleServerModel) {
	name := plan.Name.ValueString()
	plan.ID = types.StringValue(fmt.Sprintf(
		"/subscriptions/%s/resourceGroups/%s/providers/Microsoft.DBforMySQL/flexibleServers/%s",
		r.subscriptionID, plan.ResourceGroupName.ValueString(), name,
	))
	plan.FQDN = types.StringValue(fmt.Sprintf("%s.mysql.database.azure.com", name))
}

func (r *MySQLFlexibleServerResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan MySQLFlexibleServerModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	r.applyComputed(&plan)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *MySQLFlexibleServerResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state MySQLFlexibleServerModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *MySQLFlexibleServerResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan MySQLFlexibleServerModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	r.applyComputed(&plan)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *MySQLFlexibleServerResource) Delete(_ context.Context, _ resource.DeleteRequest, _ *resource.DeleteResponse) {
}
