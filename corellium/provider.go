package corellium

import (
	"context"
	"fmt"
	"os"

	"github.com/aimoda/go-corellium-api-client"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"terraform-provider-corellium/corellium/pkg/api"
)

// Ensure the implementation satisfies the expected interfaces
var (
	_ provider.Provider = &corelliumProvider{}
)

// New is a helper function to simplify provider server and testing implementation.
func New() provider.Provider {
	return &corelliumProvider{}
}

// corelliumProvider is the provider implementation.
type corelliumProvider struct{}

// Metadata returns the provider type name.
func (p *corelliumProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "corellium"
}

// Schema defines the provider-level schema for configuration data.
func (p *corelliumProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"token": schema.StringAttribute{
				Sensitive: true,
				Required:  true,
			},
		},
	}
}

// hashicupsProviderModel maps provider schema data to a Go type.
// The Terraform Plugin Framework uses Go struct types with 'tfsdk' struct field tags to map schema definitions into Go types with the actual data.
// NOTE: The types within the struct must align with the types in the schema above.
type corelliumProviderModel struct {
	Token types.String `tfsdk:"token"`
}

// Configure prepares a corellium API client for data sources and resources.
func (p *corelliumProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	// Retrieve provider data from configuration
	var config corelliumProviderModel
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// If practitioner provided a configuration value for any of the
	// attributes, it must be a known value.

	if config.Token.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("token"),
			"Unknown Corellium Token",
			"The provider cannot create the Corellium API client as there is an unknown configuration value for the Corellium API token. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the CORELLIUM_TOKEN environment variable.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	// Default values to environment variables, but override
	// with Terraform configuration value if set.

	token := os.Getenv("CORELLIUM_TOKEN")
	if !config.Token.IsNull() {
		token = config.Token.ValueString()
	}

	api.SetAccessToken(token)

	api.SetAccessToken(token)

	// If any of the expected configurations are missing, return
	// errors with provider-specific guidance.

	if token == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("token"),
			"Missing Corellium API Token",
			"The provider cannot create the Corellium API client as there is a missing or empty value for the Corellium API token. "+
				"Set the token value in the configuration or use the CORELLIUM_TOKEN environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	configuration := corellium.NewConfiguration()
	configuration.Host = "moda.enterprise.corellium.com"

	client := corellium.NewAPIClient(configuration)
	r, err := client.StatusApi.V1Ready(ctx).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `StatusApi.V1Ready``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
		resp.Diagnostics.AddError(
			"Unable to Create Corellium API Client",
			"An unexpected error occurred when creating the Corellium API client. "+
				"If the error is not clear, please contact the provider developers.\n\n"+
				"Corellium Client Error: "+err.Error(),
		)
	}

	resp.DataSourceData = client
	resp.ResourceData = client
}

// DataSources defines the data sources implemented in the provider.
func (p *corelliumProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewCorelliumV1ReadyDataSource,
		NewCorelliumV1GetInstancesSource,
		NewCorelliumV1SupportedModelsDataSource,
		NewCorelliumV1ModelSoftwareDataSource,
		NewCorelliumV1RolesDataSource,
	}
}

// Resources defines the resources implemented in the provider.
func (p *corelliumProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewCorelliumV1ImageResource,
		NewCorelliumV1ProjectResource,
	}
}
