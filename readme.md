# Playground

Playground is a project designed to provide the minimum requirements for working with microservices. It integrates various tools and components to streamline the development process.

## Features

- **Database**: Utilizes PostgreSQL with GORM as the ORM for database operations.
- **Logging**: Offers logging capabilities to Logstash, file, and standard output (stdout) using Zap.
- **Tracing**: Incorporates APM (Application Performance Monitoring) as a tracer for performance monitoring.
- **Framework**: Built on top of Echo, a high-performance Go web framework.
- **Secret Management**: Integrates Vault as a secret manager for secure handling of sensitive information.
- **Command Line Interface**: Employs Cobra for building command-line interfaces.
- **Commands**:
    - `serve`: Command to serve the microservice.
    - `migrate`: Command to perform database migrations.
    - `seed`: Command to seed initial data into the database.

## Usage

To use Playground, follow these steps:

1. **Install Dependencies**: Ensure you have all required dependencies installed. You may refer to the project's documentation for detailed instructions.

2. **Configuration**: Configure Playground according to your requirements, including database settings, logging configurations, tracer settings, etc. Make sure to set up Vault for managing secrets.

3. **Commands**:
    - `serve`: Execute `app serve` to start the microservice.
    - `migrate`: Run `app migrate` to perform database migrations.
    - `seed`: Use `app seed` to seed initial data into the database.

4. **Development**: Start developing your microservices using Playground. Leverage its integrated tools and components to enhance productivity and efficiency.

## Contributing

Contributions to Playground are welcome! If you have any suggestions, bug reports, or feature requests, feel free to open an issue or submit a pull request.

## License

This project is licensed under the [MIT License](LICENSE).
