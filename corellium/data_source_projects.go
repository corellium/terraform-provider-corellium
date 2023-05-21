package corellium

import (
	"context"
	"io"
	"math/big"
	"net/http"

	"github.com/aimoda/go-corellium-api-client"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"terraform-provider-corellium/corellium/pkg/api"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &V1ProjectsDataSource{}
	_ datasource.DataSourceWithConfigure = &V1ProjectsDataSource{}
)

// NewCorelliumV1ProjectsDataSource is a helper function to simplify the provider implementation.
func NewCorelliumV1ProjectsDataSource() datasource.DataSource {
	return &V1ProjectsDataSource{}
}

// V1ProjectsDataSource is the data source implementation.
type V1ProjectsDataSource struct {
	client *corellium.APIClient
}

type V1ProjectsModel struct {
	Id       types.String     `tfsdk:"id"`
	Projects []V1ProjectModel `tfsdk:"projects"`
}

// Metadata returns the data source type name.
func (d *V1ProjectsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_v1projects"
	// TypeName is the name of the data resource type, which must be unique within the provider.
	// This is used to identify the data resource type in state and plan files.
	// i.e: data corellium_v1projects "projects" { ... }
}

// Schema defines the schema for the data source.
func (d *V1ProjectsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
			},
			"projects": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Description: "Project id",
							Computed:    true,
						},
						"name": schema.StringAttribute{
							Description: "Project name",
							Computed:    true,
						},
						"settings": schema.SingleNestedAttribute{
							Description: "Project settings",
							Computed:    true,
							Attributes: map[string]schema.Attribute{
								"version": schema.NumberAttribute{
									Description: "Project version",
									Computed:    true,
								},
								"internet_access": schema.BoolAttribute{
									Description: "Project internet access",
									Computed:    true,
								},
								"dhcp": schema.BoolAttribute{
									Description: "Project dhcp",
									Computed:    true,
								},
							},
						},
						"quotas": schema.SingleNestedAttribute{
							Description: "Project quotas",
							Computed:    true,
							Attributes: map[string]schema.Attribute{
								"name": schema.StringAttribute{
									Description: "Project quota name",
									Computed:    true,
								},
								"cores": schema.NumberAttribute{
									Description: "Project quota cores",
									Computed:    true,
								},
								"instances": schema.NumberAttribute{
									Description: "Project quota instances",
									Computed:    true,
								},
								"ram": schema.NumberAttribute{
									Description: "Project quota ram",
									Computed:    true,
								},
							},
						},
						"users": schema.ListNestedAttribute{
							Description: "Project users",
							Computed:    true,
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"id": schema.StringAttribute{
										Description: "Project user id",
										Computed:    true,
									},
									"name": schema.StringAttribute{
										Description: "Project user name",
										Computed:    true,
									},
									"label": schema.StringAttribute{
										Description: "Project user label",
										Computed:    true,
									},
									"email": schema.StringAttribute{
										Description: "Project user email",
										Computed:    true,
									},
									"role": schema.StringAttribute{
										Description: "Project user role",
										Computed:    true,
									},
								},
							},
						},
						"teams": schema.ListNestedAttribute{
							Description: "Project teams",
							Computed:    true,
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"id": schema.StringAttribute{
										Description: "Project team id",
										Computed:    true,
									},
									"label": schema.StringAttribute{
										Description: "Project team label",
										Computed:    true,
									},
									"role": schema.StringAttribute{
										Description: "Project team role",
										Computed:    true,
									},
								},
							},
						},
						"keys": schema.ListNestedAttribute{
							Description: "Project keys",
							Required:    true,
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"id": schema.StringAttribute{
										Description: "ProjectKey ID",
										Computed:    true,
									},
									"label": schema.StringAttribute{
										Description: "ProjectKey label",
										Required:    true,
									},
									"kind": schema.StringAttribute{
										Description: "ProjectKey kind",
										Required:    true,
										Validators: []validator.String{
											stringvalidator.OneOf("ssh", "adb"),
										},
									},
									"key": schema.StringAttribute{
										Description: "ProjectKey key",
										Required:    true,
									},
									"fingerprint": schema.StringAttribute{
										Description: "ProjectKey fingerprint",
										Computed:    true,
									},
									"created_at": schema.StringAttribute{
										Description: "ProjectKey creation date",
										Computed:    true,
									},
									"updated_at": schema.StringAttribute{
										Description: "ProjectKey last update date",
										Computed:    true,
									},
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
				},
			},
		},
	}
}

// Read refreshes the Terraform state with the latest data.
func (d *V1ProjectsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state V1ProjectsModel

	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	auth := context.WithValue(ctx, corellium.ContextAccessToken, api.GetAccessToken())
	projects, r, err := d.client.ProjectsApi.V1GetProjects(auth).Execute()
	if err != nil {
		if r.StatusCode == http.StatusForbidden {
			resp.Diagnostics.AddError(
				"Error getting projects",
				"You don't have permission to get the projects",
			)
			return
		}

		b, err := io.ReadAll(r.Body)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error getting projects",
				"Coudn't read the response body: "+err.Error(),
			)
			return
		}

		resp.Diagnostics.AddError(
			"Error getting projects",
			"An unexpected error was encountered trying to get the projects\n\n"+string(b),
		)
		return
	}

	id, err := uuid.GenerateUUID()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error generating UUID",
			"An unexpected error was encountered trying to generate UUID\n\n"+err.Error(),
		)
		return
	}
	state.Id = types.StringValue(id)
	state.Projects = make([]V1ProjectModel, len(projects))
	for i, project := range projects {
		state.Projects[i].Id = types.StringValue(project.GetId())
		state.Projects[i].Name = types.StringValue(project.GetName())

		state.Projects[i].Settings = &V1ProjectSettingsModel{
			Version:        types.NumberValue(big.NewFloat(float64(project.Settings.GetVersion()))),
			InternetAccess: types.BoolValue(project.Settings.GetInternetAccess()),
			Dhcp:           types.BoolValue(project.Settings.GetDhcp()),
		}

		state.Projects[i].Quotas = &V1ProjectQuotasModel{
			Name:      types.StringValue(project.GetName()),
			Cores:     types.NumberValue(big.NewFloat(float64(project.Quotas.GetCores()))),
			Instances: types.NumberValue(big.NewFloat(float64(project.Quotas.GetInstances()))),
			Ram:       types.NumberValue(big.NewFloat(float64(project.Quotas.GetRam()))),
		}

		diags = resp.State.Set(ctx, &state)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}

		projectKeys, r, err := d.client.ProjectsApi.V1GetProjectKeys(auth, project.Id).Execute()
		if err != nil {
			b, err := io.ReadAll(r.Body)
			if err != nil {
				resp.Diagnostics.AddError(
					"Error reading the project keys",
					"Coudn't read the response body: "+err.Error(),
				)
				return
			}

			resp.Diagnostics.AddError(
				"Unable to read project keys",
				"An unexpected error was encountered trying to read the project keys from the project:\n\n"+string(b))
			return
		}

		state.Projects[i].Keys = make([]V1ProjectKeyModel, len(projectKeys))
		for j, key := range projectKeys {
			state.Projects[i].Keys[j].Id = types.StringValue(key.GetIdentifier())
			state.Projects[i].Keys[j].Label = types.StringValue("")
			state.Projects[i].Keys[j].Kind = types.StringValue(key.GetKind())
			state.Projects[i].Keys[j].Key = types.StringValue(key.GetKey())
			state.Projects[i].Keys[j].Fingerprint = types.StringValue(key.GetFingerprint())
			state.Projects[i].Keys[j].CreatedAt = types.StringValue(key.GetCreatedAt().String())
			state.Projects[i].Keys[j].UpdatedAt = types.StringValue(key.GetUpdatedAt().String())
		}

		diags = resp.State.Set(ctx, &state)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Configure adds the provider configured client to the data source.
func (d *V1ProjectsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.client = req.ProviderData.(*corellium.APIClient)
}
