# Go DDD Starter Template

A production-ready Go project template that follows Domain-Driven Design (DDD) principles to create clean, maintainable, and scalable applications.

## Using This Template

1. Click the "Use this template" button at the top of this repository
2. Name your new repository and create it
3. Clone your new repository locally
4. Update the following files with your project details:
   - `go.mod`: Change the module name to your project's module path
   - `.env-dev`: Configure with your local development settings
   - `README.md`: Update with your project-specific information

## Features

- ✨ Domain-Driven Design architecture
- 🔒 Clean separation of concerns
- 🛠 Built-in development tools (hot reload, testing)
- 📦 Docker support for local development
- 🗄 PostgreSQL integration ready
- ⚡️ Modern Go practices and patterns

## Development Workflow

### Prerequisites

- Go 1.21 or later
- Docker and Docker Compose
- Make

### Initial Setup

1. Create your repository from this template
2. Clone your new repository
3. Copy the environment file: `cp .env-dev .env`
4. Configure your environment variables in `.env`
5. Install development tools:
   ```sh
   make setup
   ```
6. Compile the project:
   ```sh
   make compile
   ```

### Running Dependencies

1. Start the database:
   ```sh
   docker-compose up -d
   ```
2. Load environment variables:
   ```sh
   set -a && source .env && set +a
   ```
3. Run migrations:
   ```sh
   make migration/up
   ```

### Development

#### Hot Reload
For development convenience, we use Air for hot reloading:
```sh
air
```

#### Testing
Run the test suite with:
```sh
make test
```

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request. For major changes, please open an issue first to discuss what you would like to change.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
