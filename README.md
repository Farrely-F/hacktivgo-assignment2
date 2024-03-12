# Simple E-Commerce CRUD API

This is a simple CRUD (Create, Read, Update, Delete) API for an e-commerce system, built using Go (Golang) and the Gin web framework, with PostgreSQL as the database.

## Prerequisites

Before running this application, make sure you have the following installed:

- Go (version 1.13 or higher)
- PostgreSQL
- Postman or similar tool for testing the API endpoints

## Installation

1. Clone this repository to your local machine:

```pwsh
git clone https://github.com/yourusername/e-commerce-crud-api.git
```

2. Navigate to the project directory:

3. Install dependencies:

```pwsh
go mod download
```

4. Create a `.env` file in the root directory and add your PostgreSQL connection details:

```
DB_HOST=localhost
DB_PORT=5432
DB_USER=yourusername
DB_NAME=yourdatabase
DB_PASSWORD=yourpassword
```

5. Run the application:

```pwsh
go run main.go
```

## Usage

The API provides the following endpoints:

- `POST /orders`: Create a new order.
- `GET /orders`: Retrieve all orders.
- `PUT /order/:orderId`: Update an existing order.
- `DELETE /order/:orderId`: Delete an order.
