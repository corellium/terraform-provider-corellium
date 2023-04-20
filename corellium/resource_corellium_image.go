package corellium

import (
	"context"
	"io"
	"math/big"

	"terraform-provider-corellium/corellium/pkg/api"

	"github.com/aimoda/go-corellium-api-client"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource              = &CorelliumV1ImageResource{}
	_ resource.ResourceWithConfigure = &CorelliumV1ImageResource{}
)

// NewCorelliumV1ImageResource is a helper function to simplify the provider implementation.
func NewCorelliumV1ImageResource() resource.Resource {
	return &CorelliumV1ImageResource{}
}

// CorelliumV1ImageResource is the data source implementation.
type CorelliumV1ImageResource struct {
	client *corellium.APIClient
}

// V1ImageModel maps the data source schema data.
type V1ImageModel struct {
	Status    types.String `tfsdk:"status"`
	Id        types.String `tfsdk:"id"`
	Name      types.String `tfsdk:"name"`
	Type      types.String `tfsdk:"type"`
	Filename  types.String `tfsdk:"filename"`
	Uniqueid  types.String `tfsdk:"unique_id"`
	Size      types.Number `tfsdk:"size"`
	Project   types.String `tfsdk:"project"`
	CreatedAt types.String `tfsdk:"created_at"`
	UpdateAt  types.String `tfsdk:"updated_at"`
}

// Metadata returns the data source type name.
func (d *CorelliumV1ImageResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_v1image"
}

// Schema defines the schema for the data source.
func (d *CorelliumV1ImageResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"status": schema.StringAttribute{
				Computed: true,
			},
			"id": schema.StringAttribute{
				Computed: true,
			},
			"name": schema.StringAttribute{
				Description: "Image name",
				Required:    true,
			},
			"type": schema.StringAttribute{
				Description: "Image type",
				Required:    true,
			},
			"filename": schema.StringAttribute{
				Computed: true,
			},
			"unique_id": schema.StringAttribute{
				Computed: true,
			},
			"size": schema.NumberAttribute{
				Computed: true,
			},
			"project": schema.StringAttribute{
				Description: "Project ID",
				Required:    true,
			},
			"created_at": schema.StringAttribute{
				Computed: true,
			},
			"updated_at": schema.StringAttribute{
				Computed: true,
			},
		},
	}
}

// Create creates the resource and sets the initial Terraform state.
func (d *CorelliumV1ImageResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan V1ImageModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	auth := context.WithValue(ctx, corellium.ContextAccessToken, api.GetAccessToken())
	image, r, err := d.client.ImagesApi.V1CreateImage(auth).
		Encoding("plain").
		Name(plan.Name.ValueString()).
		Type_(plan.Type.ValueString()).
		Project(plan.Project.ValueString()).
		Execute()
	if err != nil {
		b, err := io.ReadAll(r.Body)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error creating image",
				"Coudn't read the response body: "+err.Error(),
			)
			return
		}

		resp.Diagnostics.AddError(
			"Error creating image",
			"An unexpected error was encountered trying to create the image:\n\n"+string(b),
		)
		return
	}

	// NOTICE: we don't really need to reset all states here, but I consider it safer for now.
	plan.Status = types.StringValue(image.GetStatus())
	plan.Id = types.StringValue(image.GetId())
	plan.Name = types.StringValue(image.GetName())
	plan.Type = types.StringValue(image.GetType())
	plan.Filename = types.StringValue(image.GetFilename())
	plan.Uniqueid = types.StringValue(image.GetUniqueid())
	plan.Size = types.NumberValue(big.NewFloat(float64(image.GetSize())))
	plan.Project = types.StringValue(image.GetProject())
	plan.CreatedAt = types.StringValue(image.GetCreatedAt().String())
	plan.UpdateAt = types.StringValue(image.GetCreatedAt().String())

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read refreshes the Terraform state with the latest data.
func (d *CorelliumV1ImageResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state V1ImageModel

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	auth := context.WithValue(ctx, corellium.ContextAccessToken, api.GetAccessToken())
	image, r, err := d.client.ImagesApi.V1GetImage(auth, state.Id.ValueString()).Execute()
	if err != nil {
		b, err := io.ReadAll(r.Body)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error reading image",
				"Coudn't read the response body: "+err.Error(),
			)
			return
		}

		resp.Diagnostics.AddError(
			"Unable to read image",
			"An unexpected error was encountered trying to read the image:\n\n"+string(b))
		return
	}

	diags = resp.State.Set(ctx, &V1ImageModel{
		Status:    types.StringValue(image.Status),
		Id:        types.StringValue(image.GetId()),
		Name:      types.StringValue(image.GetName()),
		Type:      types.StringValue(image.GetType()),
		Filename:  types.StringValue(image.GetFilename()),
		Uniqueid:  types.StringValue(image.GetUniqueid()),
		Size:      types.NumberValue(big.NewFloat(float64(image.GetSize()))),
		Project:   types.StringValue(image.GetProject()),
		CreatedAt: types.StringValue(image.GetCreatedAt().String()),
		UpdateAt:  types.StringValue(image.GetCreatedAt().String()),
	})

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (d *CorelliumV1ImageResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// NOTICE: Can a image be updated?
}

// Delete deletes the resource and removes the Terraform state on success.
func (d *CorelliumV1ImageResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state V1ImageModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	auth := context.WithValue(ctx, corellium.ContextAccessToken, api.GetAccessToken())
	r, err := d.client.ImagesApi.V1DeleteImage(auth, state.Id.ValueString()).Execute()
	if err != nil {
		b, err := io.ReadAll(r.Body)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error deleting image",
				"Coudn't read the response body: "+err.Error(),
			)
			return
		}

		resp.Diagnostics.AddError(
			"Unable to delete image",
			"An unexpected error was encountered trying to delete the image:\n\n"+string(b))
		return
	}
}

// Configure adds the provider configured client to the data source.
func (d *CorelliumV1ImageResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.client = req.ProviderData.(*corellium.APIClient)
}
