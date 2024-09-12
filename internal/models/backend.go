package models

import (
	"go.uber.org/zap"
	"golang.org/x/oauth2/clientcredentials"
)

type BackendConfig struct {
	Log          *zap.Logger
	AzureEntraID clientcredentials.Config
	ApiUrl       string
}
