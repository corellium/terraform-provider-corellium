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
	_ datasource.DataSource              = &V1ModelSoftwareDataSource{}
	_ datasource.DataSourceWithConfigure = &V1ModelSoftwareDataSource{}
)

// NewCorelliumDataSource is a helper function to simplify the provider implementation.
func NewCorelliumV1ModelSoftwareDataSource() datasource.DataSource {
	return &V1ModelSoftwareDataSource{}
}

// corelliumDataSource is the data source implementation.
type V1ModelSoftwareDataSource struct {
	client *corellium.APIClient
}

type V1SoftwareDataSourceModel struct {
	Model          types.String      `tfsdk:"model"`
	Model_Software []V1SoftwareModel `tfsdk:"model_software"`
}

type V1SoftwareModel struct {
	API_Version    types.String `tfsdk:"api_version"`
	Android_Flavor types.String `tfsdk:"android_flavor"`
	Build_Id       types.String `tfsdk:"build_id"`
	Filename       types.String `tfsdk:"filename"`
	Md5_Sum        types.String `tfsdk:"md5_sum"`
	Orig_Url       types.String `tfsdk:"orig_url"`
	Release_Date   types.String `tfsdk:"release_date"`
	Sha1_Sum       types.String `tfsdk:"sha1_sum"`
	Sha256_Sum     types.String `tfsdk:"sha256_sum"`
	// TODO: Waiting on Corellium to fix this int32 issue. Data is mostly within int64 range. But GetSize() is  a int32 function.
	// Size           types.Int64  `tfsdk:"size"`
	Unique_Id   types.String `tfsdk:"unique_id"`
	Upload_Date types.String `tfsdk:"upload_date"`
	URL         types.String `tfsdk:"url"`
	Version     types.String `tfsdk:"version"`
	// TODO: Waiting for updated bindings from David
	// Metadata Metadata `tfsdk:"metadata"`
}

// TODO: Waiting for updated bindings from David
// type Metadata struct {

// }

// Metadata returns the data source type name.
func (d *V1ModelSoftwareDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_v1modelsoftware"
}

// Schema defines the schema for the data source.
func (d *V1ModelSoftwareDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"model": schema.StringAttribute{
				Required: true,
			},
			"model_software": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"api_version": schema.StringAttribute{
							Description: "Android only API version",
							Optional:    true,
						},
						"android_flavor": schema.StringAttribute{
							Description: "Android only flavor",
							Optional:    true,
						},
						"build_id": schema.StringAttribute{
							Optional: true,
						},
						"filename": schema.StringAttribute{
							Optional: true,
						},
						"md5_sum": schema.StringAttribute{
							Optional: true,
						},
						"orig_url": schema.StringAttribute{
							Description: "URL firmware is available at from vendor",
							Optional:    true,
						},
						"release_date": schema.StringAttribute{
							Description: "Release Date",
							Optional:    true,
						},
						"sha1_sum": schema.StringAttribute{
							Optional: true,
						},
						"sha256_sum": schema.StringAttribute{
							Optional: true,
						},
						// TODO: Waiting on Corellium to fix this int32 issue. Data is mostly within int64 range. But GetSize() is  a int32 function.
						// "size": schema.Int64Attribute{
						// 	Optional: true,
						// },
						"unique_id": schema.StringAttribute{
							Optional: true,
						},
						"upload_date": schema.StringAttribute{
							Description: "Date uploaded",
							Optional:    true,
						},
						"url": schema.StringAttribute{
							Description: "URL firmware is available at",
							Optional:    true,
						},
						"version": schema.StringAttribute{
							Optional: true,
						},
						// TODO Waiting for updated bindings from David
						// Metadata
					},
				},
			},
		},
	}
}

// Read refreshes the Terraform state with the latest data.
func (d *V1ModelSoftwareDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state V1SoftwareDataSourceModel
	// Get model from config
	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	auth := context.WithValue(ctx, corellium.ContextAccessToken, api.GetAccessToken())
	software, _, err := d.client.ModelsApi.V1GetModelSoftware(auth, state.Model.ValueString()).Execute()
	if err != nil && software != nil {
		resp.Diagnostics.AddError(
			"Error getting model software for model: "+state.Model.ValueString()+" Build ID: "+software[0].GetBuildid(),
			err.Error(),
		)
		return
	} else if err != nil && software == nil {
		resp.Diagnostics.AddError(
			"Failed to fetch software for model: "+state.Model.ValueString(),
			err.Error(),
		)
		return
	}

	// Map response body to model
	for _, s := range software {
		softwareState := V1SoftwareModel{
			API_Version:    types.StringValue(s.GetAPIVersion()),
			Android_Flavor: types.StringValue(s.GetAndroidFlavor()),
			Build_Id:       types.StringValue(s.GetBuildid()),
			Filename:       types.StringValue(s.GetFilename()),
			Md5_Sum:        types.StringValue(s.GetMd5sum()),
			Orig_Url:       types.StringValue(s.GetOrigUrl()),
			Release_Date:   types.StringValue(s.GetReleasedate().String()),
			Sha1_Sum:       types.StringValue(s.GetSha1sum()),
			Sha256_Sum:     types.StringValue(s.GetSha256sum()),
			// TODO: Waiting on Corellium to fix this int32 issue. Data is mostly within int64 range. But GetSize() is  a int32 function.
			// Size:           types.Int64Value(int64(s.GetSize())),
			Unique_Id:   types.StringValue(s.GetUniqueId()),
			Upload_Date: types.StringValue(s.GetUploaddate().String()),
			URL:         types.StringValue(s.GetUrl()),
			Version:     types.StringValue(s.GetVersion()),
			// TODO Waiting for updated bindings from David
			// Metadata
		}

		state.Model_Software = append(state.Model_Software, softwareState)
	}

	// Set state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Configure adds the provider configured client to the data source.
func (d *V1ModelSoftwareDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.client = req.ProviderData.(*corellium.APIClient)
}
