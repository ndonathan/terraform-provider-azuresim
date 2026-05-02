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

var _ resource.Resource = &KeyVaultKeyResource{}
var _ resource.ResourceWithConfigure = &KeyVaultKeyResource{}

type KeyVaultKeyResource struct {
	subscriptionID string
}

type KeyVaultKeyModel struct {
	ID             types.String `tfsdk:"id"`
	Name           types.String `tfsdk:"name"`
	KeyVaultID     types.String `tfsdk:"key_vault_id"`
	KeyType        types.String `tfsdk:"key_type"`
	KeySize        types.Int64  `tfsdk:"key_size"`
	Curve          types.String `tfsdk:"curve"`
	KeyOpts        types.List   `tfsdk:"key_opts"`
	NotBeforeDate  types.String `tfsdk:"not_before_date"`
	ExpirationDate types.String `tfsdk:"expiration_date"`
	Version        types.String `tfsdk:"version"`
	VersionlessID  types.String `tfsdk:"versionless_id"`
	ResourceID     types.String `tfsdk:"resource_id"`
	PublicKeyPEM   types.String `tfsdk:"public_key_pem"`
	Tags           types.Map    `tfsdk:"tags"`
}

func NewKeyVaultKeyResource() resource.Resource {
	return &KeyVaultKeyResource{}
}

func (r *KeyVaultKeyResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_key_vault_key"
}

func (r *KeyVaultKeyResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Simulates an Azure Key Vault Key.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true, Description: "Versioned data-plane URI.",
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"name": schema.StringAttribute{
				Required: true, Description: "Key name.",
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"key_vault_id": schema.StringAttribute{
				Required: true, Description: "Parent Key Vault ID.",
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"key_type": schema.StringAttribute{
				Required: true, Description: "`EC`, `EC-HSM`, `RSA`, or `RSA-HSM`.",
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"key_size":        schema.Int64Attribute{Optional: true, Description: "RSA key size (e.g. 2048, 3072, 4096)."},
			"curve":           schema.StringAttribute{Optional: true, Description: "EC curve (e.g. `P-256`, `P-384`, `P-521`, `P-256K`)."},
			"key_opts":        schema.ListAttribute{Required: true, ElementType: types.StringType, Description: "Permitted operations (`encrypt`, `decrypt`, `sign`, `verify`, `wrapKey`, `unwrapKey`)."},
			"not_before_date": schema.StringAttribute{Optional: true, Description: "RFC 3339 timestamp."},
			"expiration_date": schema.StringAttribute{Optional: true, Description: "RFC 3339 timestamp."},
			"version": schema.StringAttribute{
				Computed: true, Description: "Key version (simulated).",
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"versionless_id": schema.StringAttribute{
				Computed: true, Description: "Data-plane URI without version.",
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"resource_id": schema.StringAttribute{
				Computed: true, Description: "Resource Manager ID.",
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"public_key_pem": schema.StringAttribute{
				Computed: true, Description: "Placeholder PEM-encoded public key (not a real key).",
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"tags": schema.MapAttribute{
				Optional: true, ElementType: types.StringType, Description: "Tags.",
			},
		},
	}
}

func (r *KeyVaultKeyResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	r.subscriptionID = req.ProviderData.(string)
}

func (r *KeyVaultKeyResource) applyComputed(plan *KeyVaultKeyModel) {
	vaultName := vaultNameFromID(plan.KeyVaultID.ValueString())
	name := plan.Name.ValueString()
	version := simulatedKVVersion("key/" + vaultName + "/" + name + "/" + plan.KeyType.ValueString())

	plan.Version = types.StringValue(version)
	plan.ID = types.StringValue(fmt.Sprintf("https://%s.vault.azure.net/keys/%s/%s", vaultName, name, version))
	plan.VersionlessID = types.StringValue(fmt.Sprintf("https://%s.vault.azure.net/keys/%s", vaultName, name))
	plan.ResourceID = types.StringValue(fmt.Sprintf("%s/keys/%s/%s", plan.KeyVaultID.ValueString(), name, version))
	plan.PublicKeyPEM = types.StringValue("-----BEGIN PUBLIC KEY-----\nU0lNVUxBVEVE\n-----END PUBLIC KEY-----\n")
}

func (r *KeyVaultKeyResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan KeyVaultKeyModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	r.applyComputed(&plan)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *KeyVaultKeyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state KeyVaultKeyModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *KeyVaultKeyResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan KeyVaultKeyModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	r.applyComputed(&plan)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *KeyVaultKeyResource) Delete(_ context.Context, _ resource.DeleteRequest, _ *resource.DeleteResponse) {
}
