package corellium

import (
	"context"
	"io"
	"math/big"
	"time"

	"github.com/aimoda/go-corellium-api-client"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"terraform-provider-corellium/corellium/pkg/api"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource              = &CorelliumV1ProjectResource{}
	_ resource.ResourceWithConfigure = &CorelliumV1ProjectResource{}
)

// NewCorelliumV1ProjectResource is a helper function to simplify the provider implementation.
func NewCorelliumV1ProjectResource() resource.Resource {
	return &CorelliumV1ProjectResource{}
}

// CorelliumV1ProjectResource is the resource implementation.
type CorelliumV1ProjectResource struct {
	client *corellium.APIClient
}

type V1ProjectSettingsModel struct {
	// Version is the project version.
	Version types.Number `tfsdk:"version"`
	// InternetAccess is a boolean that defines if the project has Internet access.
	InternetAccess types.Bool `tfsdk:"internet_access"`
	// Dhcp is a boolean that defines if the project has DHCP enabled.
	Dhcp types.Bool `tfsdk:"dhcp"`
}

type V1ProjectQuotasModel struct {
	// Name is the project name.
	Name types.String `tfsdk:"name"`
	// Core is the project cores quota.
	Cores types.Number `tfsdk:"cores"`
	// Instances is the project instances quota.
	Instances types.Number `tfsdk:"instances"`
	// Ram is the project RAM quota.
	Ram types.Number `tfsdk:"ram"`
}

// V1ProjectModel maps the resource schema data.
// https://github.com/aimoda/go-corellium-api-client/blob/main/docs/Project.md
type V1ProjectModel struct {
	// Id is the project ID.
	// The project ID is a uuid, universally unique identifier.
	Id types.String `tfsdk:"id"`
	// Name is the project name.
	Name types.String `tfsdk:"name"`
	// Settings is the project settings.
	Settings *V1ProjectSettingsModel `tfsdk:"settings"`
	// Quotas is the project quotas.
	Quotas *V1ProjectQuotasModel `tfsdk:"quotas"`
	// CreatedAt is the project creation date.
	CreatedAt types.String `tfsdk:"created_at"`
	// UpdatedAt is the project last update date.
	UpdatedAt types.String `tfsdk:"updated_at"`
}

// Metadata returns the resource type name.
func (d *CorelliumV1ProjectResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_v1project"
	// TypeName is the name of the resource type, which must be unique within the provider.
	// This is used to identify the resource type in state and plan files.
	// i.e: resource corellium_v1project "project" { ... }
}

// Schema defines the schema for the resource.
func (d *CorelliumV1ProjectResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Project id",
				Computed:    true,
			},
			"name": schema.StringAttribute{
				Description: "Project name",
				Required:    true,
			},
			"settings": schema.SingleNestedAttribute{
				Description: "Project settings",
				Required:    true,
				Attributes: map[string]schema.Attribute{
					"version": schema.NumberAttribute{
						Description: "Project version",
						Required:    true,
					},
					"internet_access": schema.BoolAttribute{
						Description: "Project internet access",
						Required:    true,
					},
					"dhcp": schema.BoolAttribute{
						Description: "Project dhcp",
						Required:    true,
					},
				},
			},
			"quotas": schema.SingleNestedAttribute{
				Description: "Project quotas",
				Optional:    true,
				Attributes: map[string]schema.Attribute{
					"name": schema.StringAttribute{
						Description: "Project quota name",
						Computed:    true,
					},
					"cores": schema.NumberAttribute{
						Description: "Project quota cores",
						Required:    true,
					},
					"instances": schema.NumberAttribute{
						Description: "Project quota instances",
						Required:    true,
					},
					"ram": schema.NumberAttribute{
						Description: "Project quota ram",
						Required:    true,
					},
				},
			},
			"created_at": schema.StringAttribute{
				Description: "Project created at",
				Computed:    true,
			},
			"updated_at": schema.StringAttribute{
				Description: "Project updated at",
				Optional:    true, // TODO: Check if the `Optional` flag is required.
				Computed:    true,
			},
		},
	}
}

// Create creates the resource and sets the initial Terraform state.
func (d *CorelliumV1ProjectResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan V1ProjectModel

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	BigFloatToFloat32 := func(bf *big.Float) float32 {
		f, _ := bf.Float32()
		return f
	}

	p := corellium.NewProjectWithDefaults()
	p.SetName(plan.Name.ValueString())

	if plan.Settings != nil {
		s := corellium.NewProjectSettingsWithDefaults()
		s.SetVersion(BigFloatToFloat32(plan.Settings.Version.ValueBigFloat()))
		s.SetInternetAccess(plan.Settings.InternetAccess.ValueBool())
		s.SetDhcp(plan.Settings.Dhcp.ValueBool())

		p.SetSettings(*s)
	}

	if plan.Quotas != nil {
		q := corellium.NewProjectQuotaWithDefaults()
		q.SetCores(BigFloatToFloat32(plan.Quotas.Cores.ValueBigFloat()))
		q.SetInstances(BigFloatToFloat32(plan.Quotas.Instances.ValueBigFloat()))
		q.SetRam(BigFloatToFloat32(plan.Quotas.Ram.ValueBigFloat()))

		p.SetQuotas(*q)
	}

	auth := context.WithValue(ctx, corellium.ContextAccessToken, api.GetAccessToken())
	created, r, err := d.client.ProjectsApi.V1CreateProject(auth).Project(*p).Execute()
	if err != nil {
		b, err := io.ReadAll(r.Body)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error creating project",
				"Coudn't read the response body: "+err.Error(),
			)
			return
		}

		resp.Diagnostics.AddError(
			"Error creating project",
			"An unexpected error was encountered trying to create the project:\n\n"+string(b),
		)
		return
	}

	project, r, err := d.client.ProjectsApi.V1GetProject(auth, created.GetId()).Execute()
	if err != nil {
		b, err := io.ReadAll(r.Body)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error get the creatd project",
				"Coudn't read the response body: "+err.Error(),
			)
			return
		}

		resp.Diagnostics.AddError(
			"Error to get created project",
			"An unexpected error was encountered trying to create the project:\n\n"+string(b),
		)
		return
	}

	plan.Id = types.StringValue(project.GetId())
	plan.Name = types.StringValue(project.GetName())

	plan.Settings.Version = types.NumberValue(big.NewFloat(float64(project.Settings.GetVersion())))
	plan.Settings.InternetAccess = types.BoolValue(project.Settings.GetInternetAccess())
	plan.Settings.Dhcp = types.BoolValue(project.Settings.GetDhcp())

	plan.Quotas.Name = types.StringValue(project.GetName())
	plan.Quotas.Cores = types.NumberValue(big.NewFloat(float64(project.Quotas.GetCores())))
	plan.Quotas.Instances = types.NumberValue(big.NewFloat(float64(project.Quotas.GetInstances())))
	plan.Quotas.Ram = types.NumberValue(big.NewFloat(float64(project.Quotas.GetRam())))

	plan.CreatedAt = types.StringValue(time.Now().Format(time.RFC3339))
	plan.UpdatedAt = types.StringValue(time.Now().Format(time.RFC3339))

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read refreshes the Terraform state with the latest data.
func (d *CorelliumV1ProjectResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state V1ProjectModel

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	auth := context.WithValue(ctx, corellium.ContextAccessToken, api.GetAccessToken())
	project, r, err := d.client.ProjectsApi.V1GetProject(auth, state.Id.ValueString()).Execute()
	if err != nil {
		b, err := io.ReadAll(r.Body)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error reading project",
				"Coudn't read the response body: "+err.Error(),
			)
			return
		}

		resp.Diagnostics.AddError(
			"Unable to read project",
			"An unexpected error was encountered trying to read the image:\n\n"+string(b))
		return
	}

	state.Id = types.StringValue(project.GetId())
	state.Name = types.StringValue(project.GetName())

	state.Settings.Version = types.NumberValue(big.NewFloat(float64(project.Settings.GetVersion())))
	state.Settings.InternetAccess = types.BoolValue(project.Settings.GetInternetAccess())
	state.Settings.Dhcp = types.BoolValue(project.Settings.GetDhcp())

	state.Quotas.Name = types.StringValue(project.GetName())
	state.Quotas.Cores = types.NumberValue(big.NewFloat(float64(project.Quotas.GetCores())))
	state.Quotas.Instances = types.NumberValue(big.NewFloat(float64(project.Quotas.GetInstances())))
	state.Quotas.Ram = types.NumberValue(big.NewFloat(float64(project.Quotas.GetRam())))

	state.CreatedAt = types.StringValue(time.Now().Format(time.RFC3339))
	state.UpdatedAt = types.StringValue(time.Now().Format(time.RFC3339))

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (d *CorelliumV1ProjectResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var state V1ProjectModel
	// state is the current state of the resource.

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var plan V1ProjectModel
	// plan is the proposed new state of the resource.

	diags = req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	BigFloatToFloat32 := func(bf *big.Float) float32 {
		f, _ := bf.Float32()
		return f
	}

	p := corellium.NewProject(state.Id.ValueString())
	p.SetName(plan.Name.ValueString())

	if plan.Settings != nil {
		s := corellium.NewProjectSettings()
		s.SetVersion(BigFloatToFloat32(plan.Settings.Version.ValueBigFloat()))
		s.SetInternetAccess(plan.Settings.InternetAccess.ValueBool())
		s.SetDhcp(plan.Settings.Dhcp.ValueBool())

		p.SetSettings(*s)
	}

	if plan.Quotas != nil {
		q := corellium.NewProjectQuota()
		q.SetCores(BigFloatToFloat32(plan.Quotas.Cores.ValueBigFloat()))
		q.SetInstances(BigFloatToFloat32(plan.Quotas.Instances.ValueBigFloat()))
		q.SetRam(BigFloatToFloat32(plan.Quotas.Ram.ValueBigFloat()))

		p.SetQuotas(*q)
	}

	auth := context.WithValue(ctx, corellium.ContextAccessToken, api.GetAccessToken())
	project, r, err := d.client.ProjectsApi.V1UpdateProject(auth, state.Id.ValueString()).Project(*p).Execute()
	if err != nil {
		b, err := io.ReadAll(r.Body)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error updating project",
				"Coudn't read the response body: "+err.Error(),
			)
			return
		}

		resp.Diagnostics.AddError(
			"Error updating project",
			"An unexpected error was encountered trying to update the project:\n\n"+string(b),
		)
		return
	}

	state.Id = types.StringValue(project.GetId())
	state.Name = types.StringValue(project.GetName())

	state.Settings.Version = types.NumberValue(big.NewFloat(float64(project.Settings.GetVersion())))
	state.Settings.InternetAccess = types.BoolValue(project.Settings.GetInternetAccess())
	state.Settings.Dhcp = types.BoolValue(project.Settings.GetDhcp())

	state.Quotas.Name = types.StringValue(project.GetName())
	state.Quotas.Cores = types.NumberValue(big.NewFloat(float64(project.Quotas.GetCores())))
	state.Quotas.Instances = types.NumberValue(big.NewFloat(float64(project.Quotas.GetInstances())))
	state.Quotas.Ram = types.NumberValue(big.NewFloat(float64(project.Quotas.GetRam())))

	state.UpdatedAt = types.StringValue(time.Now().Format(time.RFC3339))

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success.
func (d *CorelliumV1ProjectResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state V1ProjectModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	auth := context.WithValue(ctx, corellium.ContextAccessToken, api.GetAccessToken())
	r, err := d.client.ProjectsApi.V1DeleteProject(auth, state.Id.ValueString()).Execute()
	if err != nil {
		b, err := io.ReadAll(r.Body)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error deleting project",
				"Coudn't read the response body: "+err.Error(),
			)
			return
		}

		resp.Diagnostics.AddError(
			"Unable to delete project",
			"An unexpected error was encountered trying to delete the project:\n\n"+string(b))
		return
	}
}

// Configure adds the provider configured client to the resource.
func (d *CorelliumV1ProjectResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.client = req.ProviderData.(*corellium.APIClient)
}
