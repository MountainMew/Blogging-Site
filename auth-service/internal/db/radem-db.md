# Authentication Service - Database Layer

This package is used to create the database connection and manage all database operations for the authentication microservice.

## Folder Structure

The `db` package is organized into the following files:

- **`database.go`**: Contains all the logic for establishing the database connection and the functions for all CRUD (Create, Read, Update, Delete) operations.
- **`schemas.go`**: A manifest file that contains the raw SQL strings for table creation (`CREATE TABLE ...`).
- **`queries.go`**: Contains all the raw SQL query strings for the CRUD operations (`INSERT`, `SELECT`, etc.).

## How to Test

To start the database and test the connection, follow these steps from the `auth-service` root directory.

1.  **Build and run the services:**
    This command will build your Go application, start the Postgres database, and run them together.

    ```bash
    docker-compose up --build
    ```

2.  **Connect to the database container:**
    Open a **second terminal window** and run the following command to get a shell inside the running Postgres container.

    ```bash
    docker exec -it auth-service-db-1 bash
    ```

3.  **Start the `psql` client:**
    Once inside the container, connect to your specific database using the credentials you defined.

    ```bash
    psql -U my_test_user -d my_test_db
    ```

4.  **Check your tables:**
    After your service has run its initialization logic, you can check if the `users` table was created correctly by describing it.
    ```sql
    \d users
    ```
