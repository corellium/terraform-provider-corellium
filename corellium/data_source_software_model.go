package corellium

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"

	"terraform-provider-corellium/corellium/pkg/api"

	"github.com/aimoda/go-corellium-api-client"
	"github.com/hashicorp/go-uuid"
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
	ID             types.String      `tfsdk:"id"`
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
	Size        types.Int64  `tfsdk:"size"`
	Unique_Id   types.String `tfsdk:"unique_id"`
	Upload_Date types.String `tfsdk:"upload_date"`
	URL         types.String `tfsdk:"url"`
	Version     types.String `tfsdk:"version"`
	// TODO: Waiting for updated bindings from David
	// Metadata Metadata `tfsdk:"metadata"`
}

// type Metadata struct {
// 	// TODO: Waiting for updated bindings from David
// }

// Metadata returns the data source type name.
func (d *V1ModelSoftwareDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_v1modelsoftware"
}

// Schema defines the schema for the data source.
func (d *V1ModelSoftwareDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
			},
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
						"size": schema.Int64Attribute{
							Optional: true,
						},
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
						// "metadata": schema.SingleNestedAttribute{
						// 	Description: "Project quotas",
						// 	Optional:    true,
						// 	Attributes:  map[string]schema.Attribute{
						// 		// TODO Waiting for updated bindings from David
						// 	},
						// },
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
	// software, _, err := d.client.ModelsApi.V1GetModelSoftware(auth, state.Model.ValueString()).Execute()
	// Replace apiUrl with the actual API URL. Workaround for endpoint.
	url := strings.Join([]string{"https://", os.Getenv("CORELLIUM_API_HOST"), "/api"}, "")
	customSoftware, err := V1GetModelSoftwareManual(auth, url, state.Model.ValueString())
	// if err != nil && software != nil {
	// 	resp.Diagnostics.AddError(
	// 		"Error getting model software for model: "+state.Model.ValueString()+" Build ID: "+software[0].GetBuildid(),
	// 		err.Error(),
	// 	)
	// 	return
	// } else
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to fetch software for model: "+state.Model.ValueString(),
			err.Error(),
		)
		return
	}

	// Map response body to model
	for _, s := range customSoftware {
		// metadata := make([]Metadata, len(s.Metadata))
		// for i, md := range s.Metadata {
		// 	metadata[i] = V1Metadata{}
		// }
		softwareState := V1SoftwareModel{
			API_Version:    types.StringValue(s.APIVersion),
			Android_Flavor: types.StringValue(s.AndroidFlavor),
			Build_Id:       types.StringValue(s.BuildID),
			Filename:       types.StringValue(s.Filename),
			Md5_Sum:        types.StringValue(s.Md5Sum),
			Orig_Url:       types.StringValue(s.OrigURL),
			Release_Date:   types.StringValue(s.ReleaseDate),
			Sha1_Sum:       types.StringValue(s.Sha1Sum),
			Sha256_Sum:     types.StringValue(s.Sha256Sum),
			Size:           types.Int64Value(s.Size),
			Unique_Id:      types.StringValue(s.UniqueId),
			Upload_Date:    types.StringValue(s.Uploaddate),
			URL:            types.StringValue(s.URL),
			Version:        types.StringValue(s.Version),
			// Metadata:       types.StringValue(s.Metadata),

		}
		// Intentional placeholder ID
		// As stated in the docs, The testing framework requires an id attribute to be present in every data source and resource. In order to run tests on data sources and resources that do not have their own ID, you must implement an ID field with a placeholder value.
		id, err := uuid.GenerateUUID()
		if err != nil {
			resp.Diagnostics.AddError(
				"Error generating UUID",
				"An unexpected error was encountered trying to generate the ID:\n\n"+err.Error(),
			)
			return
		}

		state.ID = types.StringValue(id)
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

type CustomFirmware struct {
	APIVersion    string `json:"api_version"`
	AndroidFlavor string `json:"android_flavor"`
	BuildID       string `json:"buildid"`
	Filename      string `json:"filename"`
	ID            string `json:"id"`
	Md5Sum        string `json:"md5sum"`
	OrigURL       string `json:"orig_url"`
	ReleaseDate   string `json:"releasedate"`
	Sha1Sum       string `json:"sha1sum"`
	Sha256Sum     string `json:"sha256sum"`
	Size          int64  `json:"size"`
	UniqueId      string `json:"uniqueid"`
	Uploaddate    string `json:"uploaddate"`
	URL           string `json:"url"`
	Version       string `json:"version"`
	// Metadata      []CustomMetadata `json:"metadata"`
}

// type CustomMetadata struct {
// 	// TODO Waiting for updated bindings from David
// }

// Workaround for Corellium ModelsAPI. This API is not currently working as expected. (returning values such as Size that are larger than Int32 can handle)
func V1GetModelSoftwareManual(ctx context.Context, url string, model string) ([]CustomFirmware, error) {
	client := &http.Client{}

	req, err := http.NewRequestWithContext(ctx, "GET", url+"/v1/models/"+model+"/software", nil)
	if err != nil {
		return nil, err
	}

	// Get access token from context and add it to the request header
	accessToken, ok := ctx.Value(corellium.ContextAccessToken).(string)
	if !ok {
		return nil, errors.New("access token not found in context")
	}
	req.Header.Add("Authorization", "Bearer "+accessToken)

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusForbidden {
		return nil, errors.New("access token is invalid")
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error fetching firmware data: %s", resp.Status)
	}

	var firmwares []CustomFirmware
	err = json.NewDecoder(resp.Body).Decode(&firmwares)
	if err != nil {
		return nil, err
	}

	return firmwares, nil
}
