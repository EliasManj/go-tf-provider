package provider

import (
	"context"
	"os"
	"terraform-provider-movies/client"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var (
	_ provider.Provider = &moviesProvider{}
)

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &moviesProvider{
			version: version,
		}
	}
}

type moviesProvider struct {
	version string
}

type moviesModel struct {
	Host types.String `tfsdk:"host"`
	Port types.String `tfsdk:"port"`
}

// Configure implements provider.Provider.
func (*moviesProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {

	tflog.Info(ctx, "Configuring Movies Client")

	var config moviesModel
	req.Config.Get(ctx, &config)

	if config.Host.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("host"),
			"Unknown Movies API Host",
			"The provider cannot create the Movies API client as there is an unknown configuration value for the HashiCups API host. "+
				"Provide the value or use the MOVIES_HOST environment variable.",
		)
		resp.Diagnostics.AddAttributeError(
			path.Root("port"),
			"Unknown Movies API Port",
			"The provider cannot create the Movies API client as there is an unknown configuration value for the HashiCups API host. "+
				"Provide the value or use the MOVIES_PORT environment variable.",
		)
	}

	host := os.Getenv("MOVIES_HOST")
	port := os.Getenv("MOVIES_PORT")

	if !config.Host.IsNull() {
		host = config.Host.ValueString()
	}

	if !config.Port.IsNull() {
		port = config.Port.ValueString()
	}

	ctx = tflog.SetField(ctx, "movies_host", host)
	ctx = tflog.SetField(ctx, "movies_port", port)

	tflog.Debug(ctx, "Creating Movies client")

	url := "http://" + host + ":" + port

	client, err := client.NewClient(url)

	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Create Movies API Client", err.Error(),
		)
		return
	}

	// Make the HashiCups client available during DataSource and Resource
	// type Configure methods.
	resp.DataSourceData = client
	resp.ResourceData = client

	tflog.Info(ctx, "Configured Movies client", map[string]any{"success": true})
}

// DataSources implements provider.Provider.
func (*moviesProvider) DataSources(context.Context) []func() datasource.DataSource {
	//panic("unimplemented data soruces")
	return nil
}

// Metadata implements provider.Provider.
func (p *moviesProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "movies"
	resp.Version = p.version
}

// Resources implements provider.Provider.
func (*moviesProvider) Resources(context.Context) []func() resource.Resource {
	//panic("unimplemented resources")
	return []func() resource.Resource{
		NewMovieResource,
	}
}

// Schema implements provider.Provider.
func (*moviesProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Interact with movies API",
		Attributes: map[string]schema.Attribute{
			"host": schema.StringAttribute{
				Description: "URI for Movies API. May also be provided via MOVIES_HOST environment variable.",
				Optional:    true,
			},
			"port": schema.StringAttribute{
				Description: "Port for Movies API. May also be provided via MOVIES_PORT environment variable.",
				Optional:    true,
			},
		},
	}
}
