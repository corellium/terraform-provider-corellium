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

// hashicupsProviderModel maps provider schema data to a Go type.
// The Terraform Plugin Framework uses Go struct types with 'tfsdk' struct field tags to map schema definitions into Go types with the actual data.
// NOTE: The types within the struct must align with the types in the schema above.
type corelliumProviderModel struct {
	Token types.String `tfsdk:"token"`
	Host  types.String `tfsdk:"host"`
}

// Metadata returns the provider type name.
func (p *corelliumProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "corellium"
}

// Schema defines the provider-level schema for configuration data.
func (p *corelliumProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"token": schema.StringAttribute{
				Description: "The Corellium API token. This can also be set via the CORELLIUM_API_TOKEN environment variable.",
				Sensitive:   true,
				Required:    true,
			},
			"host": schema.StringAttribute{
				Description: "The Corellium API host. This can also be set via the CORELLIUM_API_HOST environment variable.",
				Optional:    true,
			},
		},
	}
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

	// If any of the expected configurations are missing, return
	// errors with provider-specific guidance.

	if config.Token.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("token"),
			"Unknown Corellium Token",
			"The provider cannot create the Corellium API client as there is an unknown configuration value for the Corellium API token. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the CORELLIUM_API_TOKEN environment variable.",
		)
	}

	// Default values to environment variables, but override
	// with Terraform configuration value if set.
	token := os.Getenv("CORELLIUM_API_TOKEN")
	if !(config.Token.IsNull()) {
		if !(config.Token.ValueString() == "") {
			// Handle the case when config.Token is empty replace with ENV var
			token = config.Token.ValueString()
		}
	}

	if token == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("token"),
			"Missing Corellium API Token",
			"The provider cannot create the Corellium API client as there is a missing or empty value for the Corellium API token. "+
				"Set the token value in the configuration or use the CORELLIUM_API_TOKEN environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	api.SetAccessToken(token)

	if resp.Diagnostics.HasError() {
		return
	}

	// NOTICE: here it is an implementation of the override of the host value from the configuration or environment
	// variable. When the host value is not set in the configuration, the default value is used. However, if the
	// CORELLIUM_API_HOST environment variable is set, it will override the default value, but IT WILL NOT override
	// the value set in the configuration. This is because the environment variable is not a Terraform configuration
	// value, but a provider-specific value. This is a good example of how to handle provider-specific values that
	// are not Terraform configuration values.

	// host default value.
	host := "app.corellium.com"

	if h := os.Getenv("CORELLIUM_API_HOST"); h != "" {
		resp.Diagnostics.AddWarning(
			"Using CORELLIUM_API_HOST environment variable",
			"The CORELLIUM_API_HOST environment variable is being used to override the host value set in the configuration.",
		)

		host = h
	}

	if !config.Host.IsNull() && !(config.Host.ValueString() == "") {
		host = config.Host.ValueString()
	}

	if resp.Diagnostics.HasError() {
		return
	}

	configuration := corellium.NewConfiguration()
	configuration.Host = host

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
		NewCorelliumV1ProjectsDataSource,
	}
}

// Resources defines the resources implemented in the provider.
func (p *corelliumProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewCorelliumV1ImageResource,
		NewCorelliumV1ProjectResource,
		NewCorelliumV1TeamResource,
		NewCorelliumV1UserResource,
		NewCorelliumV1SnapshotResource,
		NewCorelliumV1InstanceResource,
		NewCorelliumV1WebPlayerResource,
	}
}
