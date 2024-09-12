# JWT Authentication in Golang: Securing APIs with Azure

This project demonstrates how to secure APIs using JWT (JSON Web Token) authentication in Golang, leveraging Azure Active Directory for enhanced security. It accompanies this blog article, where you'll find detailed explanations and code examples.

## Features
 - Token-based authentication using JWT.
 - Integration with Azure Active Directory.
 - Easy-to-follow setup for local testing and development.
## Requirements
 - [go 1.23.1](https://go.dev/dl/)
 - [Azure Account with Entra ID](https://portal.azure.com/)

### Oprtional
 - [httpie](https://httpie.io) - A user-friendly command-line HTTP client for testing.
 - [jq](https://jqlang.github.io/jq/) - A lightweight and flexible command-line JSON processor.

## Setup Instructions
### 1. Install Dependencies
Install the necessary library using the command:
```shell
go get github.com/go-resty/resty/v2 
```

### 2. Configuration
You'll need to retrieve the following from Azure:

 - Tenant ID
 - Client ID
 - Client Secret 

Once you have these, follow these steps:

1. Copy the example environment setup file:
   ```shell
   cp setenv-example.sh setenv.sh
   ```

2. Fill in the environment variables in setenv.sh with your Azure credentials.

### 3. Start the API
1. Export your environment variables:
   ```shell
   source setenv.sh
   ```

2. Run the API:
   ```shell
   go run cmd/goapiazureauth/main.go
   ```

### 4. Test the API
You can test the API using [httpie](https://httpie.io):

- Check the status endpoint:
  ```shell
  http 127.0.0.1:8081/status
  ```
#### Expected Response:
```json
{
    "status": "OK"
}
```

- Try to access the items endpoint without a token:
```shell
http 127.0.0.1:8081/items
```

### Expected Response:
```shell
Authorization header missing
```

- Export your environment variables:
```shell
source setenv.sh
```

 - Obtain a token and test again:
   1. Run this command to fetch a token:
      ```shell
      token=$(https --form POST "https://login.microsoftonline.com/${AZURE_TEST_API_TENANTID}/oauth2/v2.0/token" "grant_type=client_credentials" "client_id=${AZURE_TEST_API_CLIENT_ID}" "client_secret=${AZURE_TEST_API_CLIENT_SECRET}" "scope=${AZURE_TEST_API_SCOPE}" |jq -r '.access_token')
      ```
   2. Access the `items` endpoint with the token:
      ```shell
      http -A bearer -a ${token} 127.0.0.1:8081/items
      ```

### Expected Response:
```shell
[
    {
        "id": "782e1c86-df9c-43a9-ae18-d072a18dd87c",
        "item": "Sureki Zealot's Insignia",
        "level": 639
    },
    {
        "id": "0c1caf8f-c843-4ca4-b0ed-d1a7843a0fb5",
        "item": "Seal of the Poisoned Pact",
        "level": 639
    },
    {
        "id": "76f07f71-c12a-4844-b2b8-3ef9bbe1aac1",
        "item": "Spymaster's Web",
        "level": 639
    }
]
```

## 5. Start the Backend
1. Export your environment variables again:
   ```shell
   source setenv.sh
   ```

2. Run the backend in a new terminal:
   ```shell
   go run cmd/backend/main.go
   ```

3. Test the backend without authentication:
   ```shell
   http 127.0.0.1:8082/items
   ```

### Expected Response:
```json
[
   {
      "id": "dc5b964d-0f6c-4420-80c8-b818a3825875",
      "item": "Sureki Zealot's Insignia",
      "level": 639
   },
   {
      "id": "ad0a6198-fa14-41b1-a487-fa0f565ec41a",
      "item": "Seal of the Poisoned Pact",
      "level": 639
   },
   {
      "id": "73b896dd-0800-4381-a558-1d1c0c3cf64d",
      "item": "Spymaster's Web",
      "level": 639
   }
]
```

## Conclusion
This is a demo project aimed at showcasing JWT-based authentication in a Golang API with Azure. Contributions are welcome! If you have any questions, feel free to open an issue so we can improve the project together.