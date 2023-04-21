package corellium

import (
	"context"
	"io"
	"math/big"
	"os"

	"github.com/aimoda/go-corellium-api-client"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"terraform-provider-corellium/corellium/pkg/api"
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

// CorelliumV1ImageResource is the resource implementation.
type CorelliumV1ImageResource struct {
	client *corellium.APIClient
}

// V1ImageModel maps the resource schema data.
// https://github.com/aimoda/go-corellium-api-client/blob/main/docs/Image.md
type V1ImageModel struct {
	Status types.String `tfsdk:"status"`
	// Id is the image ID.
	Id types.String `tfsdk:"id"`
	// Name is the image name.
	Name types.String `tfsdk:"name"`
	// Type is the image type.
	// Type can be one of the following: "fwbinary", "kernel", "devicetree", "ramdisk", "loaderfile", "sepfw", "seprom", "bootrom", "llb", "iboot", "ibootdata", "fwpackage", "partition", "backup"
	Type types.String `tfsdk:"type"`
	// Filename is the image filename or path.
	Filename types.String `tfsdk:"filename"`
	// Encapsulated is the image encapsulated flag.
	Encapsulated types.Bool `tfsdk:"encapsulated"`
	// Uniqueid is the image unique ID.
	Uniqueid types.String `tfsdk:"unique_id"`
	// Size is the image size.
	Size types.Number `tfsdk:"size"`
	// Project is the project ID.
	Project types.String `tfsdk:"project"`
	// CreatedAt is the image creation date.
	CreatedAt types.String `tfsdk:"created_at"`
}

// Metadata returns the resource type name.
func (d *CorelliumV1ImageResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_v1image"
	// TypeName is the name of the resource type, which must be unique within the provider.
	// This is used to identify the resource type in state and plan files.
	// i.e: resource corellium_v1image "image" { ... }
}

// Schema defines the schema for the resource.
func (d *CorelliumV1ImageResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"status": schema.StringAttribute{
				Computed: true,
			},
			"id": schema.StringAttribute{
				Description: "Image ID",
				Computed:    true,
			},
			"name": schema.StringAttribute{
				Description: "Image name",
				Required:    true,
			},
			"type": schema.StringAttribute{
				Description: "Image type",
				Required:    true,
				Validators: []validator.String{
					stringvalidator.OneOf("fwbinary", "kernel", "devicetree", "ramdisk", "loaderfile", "sepfw", "seprom", "bootrom", "llb", "iboot", "ibootdata", "fwpackage", "partition", "backup"),
				},
			},
			"filename": schema.StringAttribute{
				Description: "Image filename or path",
				Required:    true,
			},
			"encapsulated": schema.BoolAttribute{
				Description: "Image encapsulated flag",
				Required:    true,
			},
			"unique_id": schema.StringAttribute{
				Description: "Image unique ID",
				Computed:    true,
			},
			"size": schema.NumberAttribute{
				Description: "Image size",
				Computed:    true,
			},
			"project": schema.StringAttribute{
				Description: "Project ID",
				Required:    true,
			},
			"created_at": schema.StringAttribute{
				Description: "Image creation date",
				Computed:    true,
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

	file, err := os.Open(plan.Filename.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error opening image file",
			"Couldn't open the image file: "+err.Error(),
		)
		return
	}

	auth := context.WithValue(ctx, corellium.ContextAccessToken, api.GetAccessToken())
	// auth is the context with the access token, what is required by the API client.
	image, r, err := d.client.ImagesApi.V1CreateImage(auth).
		Encoding("plain").
		Name(plan.Name.ValueString()).
		Type_(plan.Type.ValueString()).
		File(file).
		Encapsulated(plan.Encapsulated.ValueBool()).
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

	plan.Status = types.StringValue(image.GetStatus())
	plan.Id = types.StringValue(image.GetId())
	plan.Name = types.StringValue(image.GetName())
	plan.Type = types.StringValue(image.GetType())
	plan.Filename = types.StringValue(plan.Filename.ValueString())
	plan.Encapsulated = types.BoolValue(plan.Encapsulated.ValueBool())
	plan.Filename = types.StringValue(plan.Filename.ValueString())
	plan.Uniqueid = types.StringValue(image.GetUniqueid())
	plan.Size = types.NumberValue(big.NewFloat(float64(image.GetSize())))
	plan.Project = types.StringValue(image.GetProject())
	plan.CreatedAt = types.StringValue(image.GetCreatedAt().String())

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
	// auth is the context with the access token, what is required by the API client.
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

	state.Status = types.StringValue(image.GetStatus())
	state.Id = types.StringValue(image.GetId())
	state.Name = types.StringValue(image.GetName())
	state.Type = types.StringValue(image.GetType())
	state.Filename = types.StringValue(state.Filename.ValueString())
	state.Encapsulated = types.BoolValue(state.Encapsulated.ValueBool())
	state.Uniqueid = types.StringValue(image.GetUniqueid())
	state.Size = types.NumberValue(big.NewFloat(float64(image.GetSize())))
	state.Project = types.StringValue(image.GetProject())
	state.CreatedAt = types.StringValue(image.GetCreatedAt().String())

	diags = resp.State.Set(ctx, &state)

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (d *CorelliumV1ImageResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// NOTICE: In this case, image cannot be updated, so we just return the current state.
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

// Configure adds the provider configured client to the resource.
func (d *CorelliumV1ImageResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.client = req.ProviderData.(*corellium.APIClient)
}
