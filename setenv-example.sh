#!/usr/bin/env sh

# General
export AZURE_TEST_API_LOGLEVEL=Debug
export AZURE_TEST_API_TENANTID=""

# API
export AZURE_TEST_API_AUD="api://example.jwt.application.auth"

# Backend
export AZURE_TEST_API_CLIENT_ID=""
export AZURE_TEST_API_CLIENT_SECRET=""
export AZURE_TEST_API_SCOPE="api://example.jwt.application.auth/.default"
export AZURE_TEST_API_ENDPOINT="http://localhost:8081/items"