package azure

import (
	"net/http"
	"net/http/httputil"

	"github.com/pkg/errors"

	"github.com/sirupsen/logrus"

	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2020-11-01/network"
	"github.com/Azure/go-autorest/autorest"

	azure_auth "github.com/naemono/go-cloud-actions/pkg/auth/azure"
)

// Config is an azure peering config
type Config struct {
	azure_auth.AuthConfig
	Logger *logrus.Entry
}

// Client is an azure peering client
type Client struct {
	Config
	vnpClient         network.VirtualNetworkPeeringsClient
	vnpAutorestClient autorest.Client
}

// New will return a new azure peering client
func New(conf Config) (*Client, error) {
	vnpc, err := azure_auth.NewVirtualNetworkPeeringsClient(conf.AuthConfig)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get new virtual network peerings client")
	}
	client := &Client{
		Config:            conf,
		vnpClient:         vnpc,
		vnpAutorestClient: vnpc.Client,
	}
	if client.Logger == nil {
		client.Logger = logrus.NewEntry(logrus.New())
		client.Logger.Logger.SetLevel(logrus.InfoLevel)
		client.Logger.Logger.SetFormatter(&logrus.JSONFormatter{})
	}
	return client, nil
}

func withRequestLogging() autorest.SendDecorator {
	return func(s autorest.Sender) autorest.Sender {
		return autorest.SenderFunc(func(r *http.Request) (*http.Response, error) {
			// dump request to wire format
			if dump, err := httputil.DumpRequestOut(r, true); err == nil {
				logrus.Printf("[DEBUG] AzureRM Request: \n%s\n", dump)
			} else {
				// fallback to basic message
				logrus.Printf("[DEBUG] AzureRM Request: %s to %s\n", r.Method, r.URL)
			}
			resp, err := s.Do(r)
			if resp != nil {
				// dump response to wire format
				if dump, err := httputil.DumpResponse(resp, true); err == nil {
					logrus.Printf("[DEBUG] AzureRM Response for %s: \n%s\n", r.URL, dump)
				} else {
					// fallback to basic message
					logrus.Printf("[DEBUG] AzureRM Response: %s for %s\n", resp.Status, r.URL)
				}
			} else {
				logrus.Printf("[DEBUG] Request to %s completed with no response", r.URL)
			}
			return resp, err
		})
	}
}
