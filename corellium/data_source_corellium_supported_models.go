package corellium

import (
	"context"
	"terraform-provider-corellium/corellium/pkg/api"

	"github.com/aimoda/go-corellium-api-client"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &V1SupportedModelsDataSource{}
	_ datasource.DataSourceWithConfigure = &V1SupportedModelsDataSource{}
)

// NewCorelliumDataSource is a helper function to simplify the provider implementation.
func NewCorelliumV1SupportedModelsDataSource() datasource.DataSource {
	return &V1SupportedModelsDataSource{}
}

// corelliumDataSource is the data source implementation.
type V1SupportedModelsDataSource struct {
	client *corellium.APIClient
}

type V1SupportedModelsDataSourceModel struct {
	ID               types.String           `tfsdk:"id"`
	Supported_Models []V1SupportModelsModel `tfsdk:"supported_models"`
}

type V1SupportModelsModel struct {
	Type        types.String `tfsdk:"type"`
	Name        types.String `tfsdk:"name"`
	Model       types.String `tfsdk:"model"`
	Flavor      types.String `tfsdk:"flavor"`
	Description types.String `tfsdk:"description"`
	BoardConfig types.String `tfsdk:"board_config"`
	Platform    types.String `tfsdk:"platform"`
	CpId        types.Int64  `tfsdk:"cp_id"`
	BdId        types.Int64  `tfsdk:"bd_id"`
	Peripherals types.Bool   `tfsdk:"peripherals"`
	//TODO: Waiting for updated bindings from David
	// Quotas      Quotas       `tfsdk:"quotas"`
}

//TODO: Waiting for updated bindings from David
// type Quotas struct {
// 	Cpus  types.Number `tfsdk:"cpus"`
// 	Cores types.Number `tfsdk:"cores"`
// }

// Metadata returns the data source type name.
func (d *V1SupportedModelsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_v1supportedmodels"
}

// Schema defines the schema for the data source.
func (d *V1SupportedModelsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
			},
			"supported_models": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"type": schema.StringAttribute{
							Required: true,
						},
						"name": schema.StringAttribute{
							Required: true,
						},
						"model": schema.StringAttribute{
							Required: true,
						},
						"flavor": schema.StringAttribute{
							Required: true,
						},
						"description": schema.StringAttribute{
							Optional: true,
						},
						"board_config": schema.StringAttribute{
							Optional: true,
						},
						"platform": schema.StringAttribute{
							Optional: true,
						},
						"cp_id": schema.Int64Attribute{
							Optional: true,
						},
						"bd_id": schema.Int64Attribute{
							Optional: true,
						},
						"peripherals": schema.BoolAttribute{
							Optional: true,
						},
						//TODO: Waiting for updated bindings from David
						// "quotas": schema.SingleNestedAttribute{
						// 	Required: true,
						// 	Attributes: map[string]schema.Attribute{
						// 		"cpus": schema.NumberAttribute{
						// 			Required: true,
						// 		},
						// 		"cores": schema.NumberAttribute{
						// 			Required: true,
						// 		},
						// 	},
						// },
					},
				},
			},
		},
	}
}

// Read refreshes the Terraform state with the latest data.
func (d *V1SupportedModelsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state V1SupportedModelsDataSourceModel

	auth := context.WithValue(ctx, corellium.ContextAccessToken, api.GetAccessToken())
	models, _, err := d.client.ModelsApi.V1GetModels(auth).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to fetch Corellium Supported Models",
			err.Error(),
		)
		return
	}
	// Map response body to model
	for _, model := range models {
		modelState := V1SupportModelsModel{
			Type:        types.StringValue(model.GetType()),
			Name:        types.StringValue(model.GetName()),
			Model:       types.StringValue(model.GetModel()),
			Flavor:      types.StringValue(model.GetFlavor()),
			Description: types.StringValue(model.GetDescription()),
			BoardConfig: types.StringValue(model.GetBoardConfig()),
			Platform:    types.StringValue(model.GetPlatform()),
			CpId:        types.Int64Value(int64(model.GetCpid())),
			BdId:        types.Int64Value(int64(model.GetBdid())),
			Peripherals: types.BoolValue(model.GetPeripherals()),
			//TODO: Waiting for updated bindings from David
			// Quotas: Quotas{
			// 	Cpus:  types.NumberValue(model),
			// 	Cores: types.NumberValue(model.GetCores()),
			// },
		}

		state.Supported_Models = append(state.Supported_Models, modelState)
	}

	// Intentional placeholder ID
	// As stated in the docs, The testing framework requires an id attribute to be present in every data source and resource. In 	order to run tests on data sources and resources that do not have their own ID, you must implement an ID field with a placeholder value.
	state.ID = types.StringValue("placeholder")

	// Set state
	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Configure adds the provider configured client to the data source.
func (d *V1SupportedModelsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.client = req.ProviderData.(*corellium.APIClient)
}
