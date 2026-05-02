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

var _ resource.Resource = &StorageBlobResource{}
var _ resource.ResourceWithConfigure = &StorageBlobResource{}

type StorageBlobResource struct{ subscriptionID string }

type StorageBlobModel struct {
	ID                     types.String `tfsdk:"id"`
	Name                   types.String `tfsdk:"name"`
	StorageAccountName     types.String `tfsdk:"storage_account_name"`
	StorageContainerName   types.String `tfsdk:"storage_container_name"`
	Type                   types.String `tfsdk:"type"`
	Size                   types.Int64  `tfsdk:"size"`
	ContentType            types.String `tfsdk:"content_type"`
	ContentMD5             types.String `tfsdk:"content_md5"`
	Source                 types.String `tfsdk:"source"`
	SourceContent          types.String `tfsdk:"source_content"`
	SourceURI              types.String `tfsdk:"source_uri"`
	URL                    types.String `tfsdk:"url"`
	AccessTier             types.String `tfsdk:"access_tier"`
	CacheControl           types.String `tfsdk:"cache_control"`
	Metadata               types.Map    `tfsdk:"metadata"`
}

func NewStorageBlobResource() resource.Resource { return &StorageBlobResource{} }

func (r *StorageBlobResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_storage_blob"
}

func (r *StorageBlobResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Simulates an Azure Storage Blob.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{Computed: true, PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()}, Description: "Blob URL (used as ID)."},
			"name": schema.StringAttribute{Required: true, PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()}, Description: "Blob name."},
			"storage_account_name":   schema.StringAttribute{Required: true, PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()}, Description: "Parent storage account name."},
			"storage_container_name": schema.StringAttribute{Required: true, PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()}, Description: "Parent container name."},
			"type":                   schema.StringAttribute{Required: true, Description: "`BlockBlob`, `PageBlob`, or `AppendBlob`."},
			"size":                   schema.Int64Attribute{Optional: true, Description: "Blob size in bytes (PageBlob only)."},
			"content_type":           schema.StringAttribute{Optional: true, Description: "Content type."},
			"content_md5":            schema.StringAttribute{Optional: true, Description: "MD5 hash."},
			"source":                 schema.StringAttribute{Optional: true, Description: "Local file path."},
			"source_content":         schema.StringAttribute{Optional: true, Description: "Inline blob content."},
			"source_uri":             schema.StringAttribute{Optional: true, Description: "URI to copy from."},
			"url": schema.StringAttribute{
				Computed: true, Description: "Public blob URL.",
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"access_tier":   schema.StringAttribute{Optional: true, Description: "`Hot`, `Cool`, `Cold`, or `Archive`."},
			"cache_control": schema.StringAttribute{Optional: true, Description: "Cache-Control header."},
			"metadata":      schema.MapAttribute{Optional: true, ElementType: types.StringType, Description: "Blob metadata."},
		},
	}
}

func (r *StorageBlobResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	r.subscriptionID = req.ProviderData.(string)
}

func (r *StorageBlobResource) applyComputed(plan *StorageBlobModel) {
	url := fmt.Sprintf("https://%s.blob.core.windows.net/%s/%s",
		plan.StorageAccountName.ValueString(),
		plan.StorageContainerName.ValueString(),
		plan.Name.ValueString())
	plan.ID = types.StringValue(url)
	plan.URL = types.StringValue(url)
}

func (r *StorageBlobResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan StorageBlobModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	r.applyComputed(&plan)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *StorageBlobResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state StorageBlobModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *StorageBlobResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan StorageBlobModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	r.applyComputed(&plan)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *StorageBlobResource) Delete(_ context.Context, _ resource.DeleteRequest, _ *resource.DeleteResponse) {
}
