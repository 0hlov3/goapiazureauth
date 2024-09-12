package microsoft

import (
	"context"
	"fmt"
	"github.com/0hlov3/goapiazureauth/internal/helpers"
	"github.com/0hlov3/goapiazureauth/internal/models"
	"go.uber.org/zap"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
	"os"
)

func NewAzureEntraID(c *models.BackendConfig) models.BackendConfig {

	c.Log.Info("Initializing Microsoft Azure Entra-ID")
	tenantID := os.Getenv("AZURE_TEST_API_TENANTID")
	azureADEndpoint := fmt.Sprintf("https://login.microsoftonline.com/%s/oauth2/v2.0/token", tenantID)
	clientID := os.Getenv("AZURE_TEST_API_CLIENT_ID")
	clientSecret := os.Getenv("AZURE_TEST_API_CLIENT_SECRET")
	scopes := []string{os.Getenv("AZURE_TEST_API_SCOPE")}

	c.Log.Info(fmt.Sprintf("tenantID: %s", tenantID))
	c.Log.Info(fmt.Sprintf("azureADEndpoint: %s", azureADEndpoint))
	c.Log.Info(fmt.Sprintf("clientID: %s", clientID))
	c.Log.Info(fmt.Sprintf("scopes: %s", scopes))

	if helpers.ContainsEmpty(tenantID, azureADEndpoint, clientID, clientSecret) {
		c.Log.Fatal("One ore more Variables not set.")
	}

	adConfig := clientcredentials.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		TokenURL:     azureADEndpoint,
		Scopes:       scopes,
	}

	return models.BackendConfig{
		Log:          c.Log,
		AzureEntraID: adConfig,
		ApiUrl:       c.ApiUrl,
	}
}

func ConfigNewClient(auth *models.BackendConfig) (*oauth2.Token, error) {
	auth.Log.Info("Creating new Azure OAuth2 JWT")
	oauthConfig := clientcredentials.Config{
		ClientID:     auth.AzureEntraID.ClientID,
		ClientSecret: auth.AzureEntraID.ClientSecret,
		TokenURL:     auth.AzureEntraID.TokenURL,
		Scopes:       auth.AzureEntraID.Scopes,
	}
	token, err := oauthConfig.Token(context.Background())
	if err != nil {
		auth.Log.Fatal("Error getting token", zap.Error(err))
	}
	auth.Log.Info("Successfully created new Azure OAuth2 JWT")
	return token, err
}
