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

var _ resource.Resource = &KeyVaultResource{}
var _ resource.ResourceWithConfigure = &KeyVaultResource{}

type KeyVaultResource struct {
	subscriptionID string
}

type KeyVaultModel struct {
	ID                           types.String `tfsdk:"id"`
	Name                         types.String `tfsdk:"name"`
	ResourceGroupName            types.String `tfsdk:"resource_group_name"`
	Location                     types.String `tfsdk:"location"`
	TenantID                     types.String `tfsdk:"tenant_id"`
	SKUName                      types.String `tfsdk:"sku_name"`
	EnabledForDeployment         types.Bool   `tfsdk:"enabled_for_deployment"`
	EnabledForDiskEncryption     types.Bool   `tfsdk:"enabled_for_disk_encryption"`
	EnabledForTemplateDeployment types.Bool   `tfsdk:"enabled_for_template_deployment"`
	EnableRBACAuthorization      types.Bool   `tfsdk:"enable_rbac_authorization"`
	PurgeProtectionEnabled       types.Bool   `tfsdk:"purge_protection_enabled"`
	SoftDeleteRetentionDays      types.Int64  `tfsdk:"soft_delete_retention_days"`
	PublicNetworkAccessEnabled   types.Bool   `tfsdk:"public_network_access_enabled"`
	AccessPolicy                 types.List   `tfsdk:"access_policy"`
	VaultURI                     types.String `tfsdk:"vault_uri"`
	Tags                         types.Map    `tfsdk:"tags"`
}

func NewKeyVaultResource() resource.Resource {
	return &KeyVaultResource{}
}

func (r *KeyVaultResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_key_vault"
}

func (r *KeyVaultResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Simulates an Azure Key Vault.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true, Description: "The simulated Key Vault ID.",
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"name": schema.StringAttribute{
				Required: true, Description: "Name of the Key Vault (must be globally unique in real Azure).",
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"resource_group_name": schema.StringAttribute{
				Required: true, Description: "Name of the Resource Group.",
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"location": schema.StringAttribute{
				Required: true, Description: "Azure region.",
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"tenant_id": schema.StringAttribute{
				Required: true, Description: "Azure AD tenant ID.",
			},
			"sku_name": schema.StringAttribute{
				Required: true, Description: "`standard` or `premium`.",
			},
			"enabled_for_deployment":          schema.BoolAttribute{Optional: true, Description: "Allow VMs to retrieve secrets."},
			"enabled_for_disk_encryption":     schema.BoolAttribute{Optional: true, Description: "Allow Disk Encryption to retrieve secrets and unwrap keys."},
			"enabled_for_template_deployment": schema.BoolAttribute{Optional: true, Description: "Allow Resource Manager to retrieve secrets."},
			"enable_rbac_authorization":       schema.BoolAttribute{Optional: true, Description: "Use RBAC instead of access policies."},
			"purge_protection_enabled":        schema.BoolAttribute{Optional: true, Description: "Enable purge protection."},
			"soft_delete_retention_days": schema.Int64Attribute{
				Optional: true, Computed: true,
				Description: "Soft-delete retention (7-90 days). Defaults to 90.",
			},
			"public_network_access_enabled": schema.BoolAttribute{
				Optional: true, Computed: true,
				Description: "Whether public network access is enabled. Defaults to `true`.",
			},
			"vault_uri": schema.StringAttribute{
				Computed: true, Description: "Simulated vault URI (`https://<name>.vault.azure.net/`).",
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"tags": schema.MapAttribute{
				Optional: true, ElementType: types.StringType, Description: "Tags.",
			},
		},
		Blocks: map[string]schema.Block{
			"access_policy": schema.ListNestedBlock{
				Description: "Access policy entry. Mutually exclusive with `enable_rbac_authorization = true`.",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"tenant_id":      schema.StringAttribute{Required: true, Description: "Azure AD tenant ID."},
						"object_id":      schema.StringAttribute{Required: true, Description: "Principal object ID."},
						"application_id": schema.StringAttribute{Optional: true, Description: "Application ID (for service principal access)."},
						"key_permissions": schema.ListAttribute{
							Optional: true, ElementType: types.StringType,
							Description: "Permitted key operations.",
						},
						"secret_permissions": schema.ListAttribute{
							Optional: true, ElementType: types.StringType,
							Description: "Permitted secret operations.",
						},
						"certificate_permissions": schema.ListAttribute{
							Optional: true, ElementType: types.StringType,
							Description: "Permitted certificate operations.",
						},
						"storage_permissions": schema.ListAttribute{
							Optional: true, ElementType: types.StringType,
							Description: "Permitted storage operations.",
						},
					},
				},
			},
		},
	}
}

func (r *KeyVaultResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	r.subscriptionID = req.ProviderData.(string)
}

func (r *KeyVaultResource) applyComputed(plan *KeyVaultModel) {
	name := plan.Name.ValueString()
	plan.ID = types.StringValue(fmt.Sprintf(
		"/subscriptions/%s/resourceGroups/%s/providers/Microsoft.KeyVault/vaults/%s",
		r.subscriptionID, plan.ResourceGroupName.ValueString(), name,
	))
	plan.VaultURI = types.StringValue(fmt.Sprintf("https://%s.vault.azure.net/", name))
	if plan.SoftDeleteRetentionDays.IsNull() || plan.SoftDeleteRetentionDays.IsUnknown() {
		plan.SoftDeleteRetentionDays = types.Int64Value(90)
	}
	if plan.PublicNetworkAccessEnabled.IsNull() || plan.PublicNetworkAccessEnabled.IsUnknown() {
		plan.PublicNetworkAccessEnabled = types.BoolValue(true)
	}
}

func (r *KeyVaultResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan KeyVaultModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	r.applyComputed(&plan)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *KeyVaultResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state KeyVaultModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *KeyVaultResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan KeyVaultModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	var state KeyVaultModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	plan.ID = state.ID
	plan.VaultURI = state.VaultURI
	r.applyComputed(&plan)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *KeyVaultResource) Delete(_ context.Context, _ resource.DeleteRequest, _ *resource.DeleteResponse) {
}
