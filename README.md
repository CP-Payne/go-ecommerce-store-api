# Ecommerce Store API

## Overview
The **Ecommerce Store API** provides a set of functionalities required for managing an online store. This API enables user authentication (registration and login), shopping cart operations, payments (using PayPal), product listings,  and reviews. It is designed to be a backend system that integrates with any frontend framework.
## API Postman Documentation
Detailed API documentation can be found on Postman, providing comprehensive details about each endpoint:
- **Postman Documentation:** [Ecommerce Store Documentation](https://documenter.getpostman.com/view/21095392/2sAXxLDEks)
- **Postman Collection:** [Postman collection](./postman/API%20Postman%20Collection.json)
## Installation
Follow these steps to set up the Ecommerce Store API locally:
### Prerequisites
Ensure that the following tools are installed on your machine:
- **Go** (Version 1.23 or later)
- **Git** for cloning the repository
- **PostgreSQL** (or use Docker for database setup)
### Cloning the Repository

Clone the repository to your local machine:

```bash
git clone https://github.com/CP-Payne/go-ecommerce-store-api
cd go-ecommerce-store-api
```
### Installing Dependencies
The project uses `go.mod` for managing dependencies. Install the required dependencies:

```bash
go mod download
```
This will fetch and install all the necessary Go packages and modules.
### Environment Setup
You need to configure the environment variables before running the API. Create a `.env` file in the root directory with the following values:

```plaintext
# Database environment variables (used in docker-compose)
POSTGRES_DB=<db_name>
POSTGRES_USER=<db_user>
POSTGRES_PASSWORD=<db_password>
POSTGRES_HOST=<db_host>
POSTGRES_PORT=<db_port>

# Api environment variables
PORT=3000
JWT_SECRET=<jwt_secret>

# PAYPAL API CREDS
PAYPAL_CLIENT=<paypal_client_id>
PAYPAL_SECRET=<paypal_secret>
```
- **POSTGRES variables**: Replace these with your PostgreSQL database credentials. If you don't have a PostgreSQL setup, you can use Docker (see the "Database Setup" section below).
- **JWT_SECRET**: A secret key used for signing JSON Web Tokens (JWT).
- **PayPal credentials**: Obtain your PayPal Client ID and Secret by creating a developer account on PayPal (see [Get Started with PayPal REST APIs](https://developer.paypal.com/api/rest/?_ga=2.150971572.368875705.1720450729-1774217071.1701640500&_gac=1.82635492.1720023622.Cj0KCQjw7ZO0BhDYARIsAFttkCgWb0D7wzz0Xq70uhuDYTv5e8bPDEwnDYKG8Gavy5V6iIaMfCL4y7IaAoW1EALw_wcB#link-getclientidandclientsecret))
### Database Setup

#### Using Docker for PostgreSQL
If you prefer using Docker for the database setup, navigate to the `/scripts` directory and run the following command to start a PostgreSQL container:
```bash
cd ./scripts
docker-compose up -d
```

#### Setting Up the Database Schema
Once your database is up, you'll need to create the necessary database schema. This can be done manually or via migration tools like `goose`. To automate the migration process using `goose`, run:
```bash
cd ./sql/schema
goose postgres "postgres://username:password@host:port/database" up
```
Alternatively, you can run the SQL files located in `/sql/schema/` manually.
#### Populating Test Data
To populate the database with test data, particularly product data, run the SQL script located in:
```bash
./sql/test_data/products.sql
```

### Running the Server

After completing the setup, you can start the API server by running the following command from the root of the project:
```bash
go run ./cmd/api/main.go
```

The API server will run on the port specified in the `.env` file (default is `3000`). You can access it via `http://localhost:<port>`.

## Tools and Technologies Used
- **Golang** (v1.23)
- **Goose** (for database migrations)
- **SQLC** (for generating type-safe SQL queries)
- **PostgreSQL** (database)
- **PayPal REST API** (for payment processing)

## Future Enhancements
Here are some potential enhancements that could be added to the API:

- **Admin Dashboard**: Create an administrative interface for managing products, orders, and user accounts.
- **Sales Analytics**: Integrate a dashboard to monitor and analyse sales data.
- **Refund Functionality**: Add support for issuing refunds through the PayPal API.
- **Multi-Payment Gateway Support**: Integrate additional payment processors like Stripe or Square for more flexibility.