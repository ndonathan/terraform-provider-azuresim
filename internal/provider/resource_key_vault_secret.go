package provider

import (
	"context"
	"crypto/sha1"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ resource.Resource = &KeyVaultSecretResource{}
var _ resource.ResourceWithConfigure = &KeyVaultSecretResource{}

type KeyVaultSecretResource struct {
	subscriptionID string
}

type KeyVaultSecretModel struct {
	ID              types.String `tfsdk:"id"`
	Name            types.String `tfsdk:"name"`
	KeyVaultID      types.String `tfsdk:"key_vault_id"`
	Value           types.String `tfsdk:"value"`
	ContentType     types.String `tfsdk:"content_type"`
	NotBeforeDate   types.String `tfsdk:"not_before_date"`
	ExpirationDate  types.String `tfsdk:"expiration_date"`
	Version         types.String `tfsdk:"version"`
	VersionlessID   types.String `tfsdk:"versionless_id"`
	ResourceID      types.String `tfsdk:"resource_id"`
	Tags            types.Map    `tfsdk:"tags"`
}

func NewKeyVaultSecretResource() resource.Resource {
	return &KeyVaultSecretResource{}
}

func (r *KeyVaultSecretResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_key_vault_secret"
}

func (r *KeyVaultSecretResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Simulates an Azure Key Vault Secret.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true, Description: "Versioned data-plane URI.",
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"name": schema.StringAttribute{
				Required: true, Description: "Secret name.",
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"key_vault_id": schema.StringAttribute{
				Required: true, Description: "Parent Key Vault ID.",
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"value": schema.StringAttribute{
				Required: true, Sensitive: true, Description: "Secret value.",
			},
			"content_type":    schema.StringAttribute{Optional: true, Description: "Free-form content type tag."},
			"not_before_date": schema.StringAttribute{Optional: true, Description: "RFC 3339 timestamp (e.g. `2026-01-01T00:00:00Z`)."},
			"expiration_date": schema.StringAttribute{Optional: true, Description: "RFC 3339 timestamp."},
			"version": schema.StringAttribute{
				Computed: true, Description: "Secret version (simulated, deterministic).",
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
			"tags": schema.MapAttribute{
				Optional: true, ElementType: types.StringType, Description: "Tags.",
			},
		},
	}
}

func (r *KeyVaultSecretResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	r.subscriptionID = req.ProviderData.(string)
}

// vaultNameFromID extracts the vault name from a Key Vault resource ID.
// Returns an empty string if the ID doesn't match the expected format.
func vaultNameFromID(id string) string {
	const marker = "/providers/Microsoft.KeyVault/vaults/"
	idx := strings.Index(id, marker)
	if idx < 0 {
		return ""
	}
	return strings.SplitN(id[idx+len(marker):], "/", 2)[0]
}

// simulatedKVVersion produces a hex string that mirrors the 32-char Azure version layout.
func simulatedKVVersion(seed string) string {
	sum := sha1.Sum([]byte(seed))
	return fmt.Sprintf("%x%x", sum[0:10], sum[10:16])[:32]
}

func (r *KeyVaultSecretResource) applyComputed(plan *KeyVaultSecretModel) {
	vaultName := vaultNameFromID(plan.KeyVaultID.ValueString())
	name := plan.Name.ValueString()
	version := simulatedKVVersion("secret/" + vaultName + "/" + name + "/" + plan.Value.ValueString())

	plan.Version = types.StringValue(version)
	plan.ID = types.StringValue(fmt.Sprintf("https://%s.vault.azure.net/secrets/%s/%s", vaultName, name, version))
	plan.VersionlessID = types.StringValue(fmt.Sprintf("https://%s.vault.azure.net/secrets/%s", vaultName, name))
	plan.ResourceID = types.StringValue(fmt.Sprintf("%s/secrets/%s/%s", plan.KeyVaultID.ValueString(), name, version))
}

func (r *KeyVaultSecretResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan KeyVaultSecretModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	r.applyComputed(&plan)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *KeyVaultSecretResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state KeyVaultSecretModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *KeyVaultSecretResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan KeyVaultSecretModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	r.applyComputed(&plan)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *KeyVaultSecretResource) Delete(_ context.Context, _ resource.DeleteRequest, _ *resource.DeleteResponse) {
}
