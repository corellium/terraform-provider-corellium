package corellium

import (
	"context"

	"github.com/aimoda/go-corellium-api-client"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &v1GetInstancesDataSource{}
	_ datasource.DataSourceWithConfigure = &v1GetInstancesDataSource{}
)

// NewCorelliumDataSource is a helper function to simplify the provider implementation.
func NewCorelliumV1GetInstancesSource() datasource.DataSource {
	return &v1GetInstancesDataSource{}
}

// corelliumDataSource is the data source implementation.
type v1GetInstancesDataSource struct {
	client *corellium.APIClient
}

// corelliumDataSourceModel maps the data source schema data.
type v1GetInstanceModel struct {
}

// Metadata returns the data source type name.
func (d *v1GetInstancesDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_v1ready"
}

// Schema defines the schema for the data source.
func (d *v1GetInstancesDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"status": schema.StringAttribute{
				Computed: true,
			},
		},
	}
}

// Read refreshes the Terraform state with the latest data.
func (d *v1GetInstancesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state v1GetInstanceModel

	status, err := d.client.StatusApi.V1Ready(ctx).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to check Corellium status",
			err.Error(),
		)
		return
	}
	// Map response body to model
	statusState := v1GetInstanceModel{
		Status: types.StringValue(status.Status),
	}
	state.Status = statusState.Status
	// Set state
	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Configure adds the provider configured client to the data source.
func (d *v1GetInstancesDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.client = req.ProviderData.(*corellium.APIClient)
}
