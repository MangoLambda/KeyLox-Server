# KeyLox-Server

KeyLox-Server is a server designed for the KeyLox password manager. It allows credentials synchronization between clients, ensuring that your passwords and other sensitive information are securely stored and synchronized across multiple devices.

## Features

- **User Registration**: Register new users with secure password hashing and salting.
- **User Authentication**: Authenticate users securely.
- **Vault Management**: Create and manage vaults for storing credentials.
- **Credentials Synchronization**: Synchronize credentials between multiple clients.
- **RESTful API**: Provides a RESTful API for interacting with the server.
- **Swagger Documentation**: Includes Swagger documentation for easy API exploration.

## Installation

1. **Clone the repository**:
    ```sh
    git clone https://github.com/yourusername/KeyLox-Server.git
    cd KeyLox-Server
    ```

2. **Install dependencies**:
    Ensure you have Go installed. Then, install the required Go packages:
    ```sh
    go mod tidy
    ```

3. **Set up the database**:
    The server uses SQLite for storage. Ensure you have SQLite installed. The database will be created automatically when you run the server.

## Usage

1. **Run the server**:
    ```sh
    go run main.go
    ```

2. **Access the API**:
    The server will be running on `http://localhost:8080`. You can use tools like `curl`, Postman, or your web browser to interact with the API.

3. **Swagger Documentation**:
    Access the Swagger documentation at `http://localhost:8080/swagger/index.html` to explore the API endpoints and their usage.

## API Endpoints

### User Endpoints

- **Register a new user**:
    ```http
    POST /register
    ```
    Request Body:
    ```json
    {
        "username": "example_user",
        "key": "base64_encoded_key",
        "clientSalt": "base64_encoded_salt"
    }
    ```

- **Get user details**:
    ```http
    GET /user/{username}
    ```

### Vault Endpoints

- **Get vault details**:
    ```http
    GET /vault/{username}
    ```

## Contributing

Contributions are welcome! Please fork the repository and submit pull requests.

## License

This project is licensed under the MIT License.

## Acknowledgements

- [Chi](https://github.com/go-chi/chi) for the HTTP router.
- [SQLite](https://www.sqlite.org/) for the database.
- [Swagger](https://swagger.io/) for API documentation.

## Contact

For any questions or suggestions, please open an issue or contact the repository owner.
