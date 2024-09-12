package models

import (
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type ApiConfig struct {
	Log   *zap.Logger
	Azure Azure
}

type Azure struct {
	TenantID string `json:"tenantId"`
	Scope    string `json:"scope"`
}

type Health struct {
	Status Status `json:"status"`
}

type Status string

const StatusOk Status = "OK"

type ItemList struct {
	Id    uuid.UUID `json:"id"`
	Item  string    `json:"item"`
	Level int16     `json:"level"`
}
