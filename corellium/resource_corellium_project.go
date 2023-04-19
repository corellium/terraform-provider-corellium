package corellium

import (
	"context"
	"fmt"
	"io"
	"math/big"
	"time"

	"terraform-provider-corellium/corellium/pkg/api"

	"github.com/aimoda/go-corellium-api-client"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
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

// CorelliumV1ProjectResource is the data source implementation.
type CorelliumV1ProjectResource struct {
	client *corellium.APIClient
}

type V1ProjectSettingsModel struct {
	Version        types.Number `tfsdk:"version"`
	InternetAccess types.Bool   `tfsdk:"internet_access"`
	Dhcp           types.Bool   `tfsdk:"dhcp"`
}

type V1ProjectQuotasModel struct {
	Name      types.String `tfsdk:"name"`
	Cores     types.Number `tfsdk:"cores"`
	Instances types.Number `tfsdk:"instances"`
	Ram       types.Number `tfsdk:"ram"`
}

// V1ProjectModel maps the data source schema data.
type V1ProjectModel struct {
	Id        types.String            `tfsdk:"id"`
	Name      types.String            `tfsdk:"name"`
	Settings  *V1ProjectSettingsModel `tfsdk:"settings"`
	Quotas    *V1ProjectQuotasModel   `tfsdk:"quotas"`
	CreatedAt types.String            `tfsdk:"created_at"`
	UpdatedAt types.String            `tfsdk:"updated_at"`
}

// Metadata returns the data source type name.
func (d *CorelliumV1ProjectResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_v1project"
}

// Schema defines the schema for the data source.
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
				Required:    true,
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

	// TODO: Is there a better way to do that?
	BigFloatToFloat32 := func(bf *big.Float) float32 {
		f, _ := bf.Float32()
		return f
	}

	s := corellium.NewProjectSettingsWithDefaults()
	s.SetVersion(BigFloatToFloat32(plan.Settings.Version.ValueBigFloat()))
	s.SetInternetAccess(plan.Settings.InternetAccess.ValueBool())
	s.SetDhcp(plan.Settings.Dhcp.ValueBool())

	q := corellium.NewProjectQuotaWithDefaults()
	q.SetCores(BigFloatToFloat32(plan.Quotas.Cores.ValueBigFloat()))
	q.SetInstances(BigFloatToFloat32(plan.Quotas.Instances.ValueBigFloat()))
	q.SetRam(BigFloatToFloat32(plan.Quotas.Ram.ValueBigFloat()))

	p := corellium.NewProjectWithDefaults()
	p.SetName(plan.Name.ValueString())
	p.SetSettings(*s)
	p.SetQuotas(*q)

	auth := context.WithValue(ctx, corellium.ContextAccessToken, api.GetAccessToken())
	project, r, err := d.client.ProjectsApi.V1CreateProject(auth).Project(*p).Execute()

	fmt.Printf("%+v", project)
	fmt.Println(project.Settings)
	fmt.Println(project.Quotas)
	fmt.Println(r)
	fmt.Println(err)

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

	plan.Id = types.StringValue(project.GetId())
	plan.Settings.Dhcp = types.BoolValue(project.Settings.GetDhcp())
	plan.Quotas.Name = types.StringValue(project.GetName())
	plan.CreatedAt = types.StringValue(time.Now().Format(time.RFC3339))
	plan.UpdatedAt = types.StringValue(time.Now().Format(time.RFC3339))
	/*plan.Settings = &V1ProjectSettingsModel{
		Version:        types.NumberValue(big.NewFloat(float64(project.Settings.GetVersion()))),
		InternetAccess: types.BoolValue(project.Settings.GetInternetAccess()),
		Dhcp:           types.BoolValue(project.Settings.GetDhcp()),
	}
	plan.Quotas = &V1ProjectQuotasModel{
		Name:      types.StringValue(project.GetName()),
		Cores:     types.NumberValue(big.NewFloat(float64(project.Quotas.GetCores()))),
		Instances: types.NumberValue(big.NewFloat(float64(project.Quotas.GetInstances()))),
		Ram:       types.NumberValue(big.NewFloat(float64(project.Quotas.GetRam()))),
	}*/
	/*plan.Settings.Version = types.NumberValue(big.NewFloat(float64(project.Settings.GetVersion())))
	plan.Settings.InternetAccess = types.BoolValue(project.Settings.GetInternetAccess())
	plan.Settings.Dhcp = types.BoolValue(project.Settings.GetDhcp())

	plan.Quotas.Name = types.StringValue(project.GetName())
	plan.Quotas.Cores = types.NumberValue(big.NewFloat(float64(project.Quotas.GetCores())))
	plan.Quotas.Instances = types.NumberValue(big.NewFloat(float64(project.Quotas.GetInstances())))
	plan.Quotas.Ram = types.NumberValue(big.NewFloat(float64(project.Quotas.GetRam())))*/

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

	diags = resp.State.Set(ctx, &V1ProjectModel{
		Id:   types.StringValue(project.GetId()),
		Name: types.StringValue(project.GetName()),
		Settings: &V1ProjectSettingsModel{
			Version:        types.NumberValue(big.NewFloat(float64(project.Settings.GetVersion()))),
			InternetAccess: types.BoolValue(project.Settings.GetInternetAccess()),
			Dhcp:           types.BoolValue(project.Settings.GetDhcp()),
		},
		Quotas: &V1ProjectQuotasModel{
			Name:      types.StringValue(project.GetName()),
			Cores:     types.NumberValue(big.NewFloat(float64(project.Quotas.GetCores()))),
			Instances: types.NumberValue(big.NewFloat(float64(project.Quotas.GetInstances()))),
			Ram:       types.NumberValue(big.NewFloat(float64(project.Quotas.GetRam()))),
		},
	})

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (d *CorelliumV1ProjectResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	/*var state V1ProjectModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	BigFloatToFloat32 := func(bf *big.Float) float32 {
		f, _ := bf.Float32()
		return f
	}

	s := corellium.NewProjectSettingsWithDefaults()
	s.SetDhcp(state.Settings.Dhcp.ValueBool())

	q := corellium.NewProjectQuotaWithDefaults()
	q.SetCores(BigFloatToFloat32(state.Quotas.Cores.ValueBigFloat()))
	q.SetInstances(BigFloatToFloat32(state.Quotas.Instances.ValueBigFloat()))
	q.SetRam(BigFloatToFloat32(state.Quotas.Ram.ValueBigFloat()))

	p := corellium.NewProjectWithDefaults()
	p.SetName(state.Name.ValueString())
	p.SetSettings(*s)
	p.SetQuotas(*q)

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
	state.Settings.Dhcp = types.BoolValue(project.Settings.GetDhcp())
	state.Quotas.Name = types.StringValue(project.GetName())

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}*/
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

// Configure adds the provider configured client to the data source.
func (d *CorelliumV1ProjectResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.client = req.ProviderData.(*corellium.APIClient)
}
