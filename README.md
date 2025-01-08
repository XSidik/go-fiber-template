# go-fiber-template
This is an API template developed using the GoFiber framework in Golang, with PostgreSQL as the database, JWT for authentication, and Redis for caching.

## Getting Started

### Prerequisites
- Go 1.23 or higher
- Docker
- Docker Compose

### Running the Application

1. **Clone the repository:**
    ```sh
    git clone https://github.com/XSidik/go-fiber-template.git
    cd go-fiber-template
    ```

2. **Set up environment variables:**
    Create a `.env` file in the root directory and copy data from .env.example, udpate the value based on you need

3. **Run PostgreSQL and Redis using Docker Compose:**
    ```sh
    docker-compose -f docker-compose-postgreSQL.yml up -d
    docker-compose -f docker-compose-redis.yml up -d
    ```

4. **Run the application:**
    ```sh
    go run main.go
    ```

### API Documentation
For detailed API documentation, you can see it soon (I will use swagger)
