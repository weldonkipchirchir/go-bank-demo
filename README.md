
# Simple Bank API

This project is a simple RESTful API for managing bank accounts. It allows users to perform operations such as creating accounts, updating balances, transferring funds, and retrieving account information.

## Installation

1. Clone the repository:

    ```bash
    git clone https://github.com/weldonkipchirchir/simple-bank-api.git
    ```

2. Navigate to the project directory:

    ```bash
    cd simple-bank-api
    ```

3. Install dependencies:

    ```bash
    go mod tidy
    ```

4. Create a `app.env` file in the root directory of the project and configure it according to your environment. Below is a sample `.env` file:

    ```plaintext
    DB_DRIVER=postgres
    DB_SOURCE=postgresql://root:mysecretpassword@postgres:5432/simple_bank?sslmode=disable
    SERVER_ADDRESS=0.0.0.0:8000
    ```

    Adjust the values accordingly based on your PostgreSQL configuration.

5. Run the following commands to set up your database:

    - **Create Docker network (if not already created)**:
    
        ```bash
        docker network create bank-network
        ```

    - **Start PostgreSQL Docker container**:

        ```bash
        docker run --name postgres12 -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:12-alpine
        ```

    - **Create the database**:

        ```bash
        docker exec -it postgres12 createdb --username=root --owner=root simple_bank
        ```

6. Install the `migrate` tool if you haven't already:

    ```bash
    go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
    ```

7. Migrate up to apply database migrations:
    DB_URL=postgresql://root:mysecretpassword@postgres:5432/simple_bank?sslmode=disable

    ```bash
    migrate -path db/migration -database "$(DB_URL)" -verbose up
    ```

## Usage

The API provides the following endpoints:

- **Create Account**: `POST /accounts`
- **Get Account**: `GET /accounts/:id`
- **Update Account Balance**: `PATCH /accounts/:id`
- **Transfer Funds**: `POST /transfers`

## Contributing

Contributions are welcome! If you find any issues or have suggestions for improvements, please open an issue or create a pull request.

