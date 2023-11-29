RESTful API for a blog platform.

## Installation

`go get -u github.com/gorilla/mux`

`go get -u github.com/rs/zerolog`


Using gorilla/mux as a router in the API has several advantages:
- Middleware Support 
gorilla/mux supports middleware, allowing you to plug in additional functionality like logging, authentication, and CORS handling into your request handling pipeline. This makes it easy to add cross-cutting concerns without cluttering the business logic.

- Subrouters 
It provides the functionality to create subrouters, which are particularly useful for organizing routes into groups, like API versioning or separating public and private API endpoints. This helps in maintaining a clean and well-structured codebase.

## Makefile usage
The provided Makefile simplifies the process of building and running the Blog API application, especially in a Docker environment. Here's a breakdown of its functionality:

Environment Variables: The Makefile includes an .env file, ensuring that any environment-specific variables are loaded.

1. Build and Run with Docker
- docker-build
This command builds a Docker image for the application. It uses the specified IMAGE_NAME and TAG, and builds the image based on the Dockerfile located at DOCKERFILE_PATH. The build context is set to BUILD_CONTEXT.

- docker-run
This command runs the Docker container in an interactive terminal mode and automatically removes the container after it's stopped. It uses the previously built image with the specified IMAGE_NAME and TAG.

2. Local Development
- api
This command runs the application locally using Go. It's useful for development and debugging purposes. 
- tests
This command runs all the tests in the application. It's essential for ensuring code quality and functionality.

This Makefile is designed to streamline the development and deployment process, making it easier to build, test, and run the application in different environments. It's particularly useful for maintaining consistency in build and deployment processes across different machines and environments.

## Endpoints

POST /login: Authenticate an author.

POST /api/posts: Create a new post.

GET /api/posts/{id}: Retrieve a specific post.

GET /api/posts: Retrieve all posts for an author.

PUT /api/posts/{id}: Update a specific post.

DELETE /api/posts/{id}: Delete a specific post.

## API Services

1. PostsService
    Handles operations related to blog posts:

    - CreatePosts
    - GetAllPosts
    - UpdatePosts
    - GetPosts
    - DeletePosts

2. AuthorsService
    Manages author authentication:
    - ValidAuthor

## API security

### Endpoint:
- POST /login

### Request
Send a POST request to the /login endpoint with the author's credentials or Admin use that can access all posts:

`{"username": "admin", "password": "admin"}`

`{"username": "Author 1", "password": "password1"}`

### Response
Upon successful authentication, the server responds with a JWT in the response body:

`{"token": "YOUR_TOKEN"}`

### Using the Token
This token must be included in the Authorization header of subsequent API requests to access protected endpoints. The header format is as follows:

`Authorization: Bearer YOUR_TOKEN`

### Second not on authentication

In the provided API server implementation, the use of JWT (JSON Web Token) for authentication is primarily for demonstration purposes and may not adhere to all best practices for secure token management, particularly regarding the security key used for token generation and validation.

Here are some key points to consider:

- Security Key Importance
JWT tokens are typically signed with a secret key to ensure their integrity and authenticity. In a demonstration or development environment, the security key might be simplified or hardcoded. However, this approach is not secure for production environments.

- Key Management in Production 
In a production setting, the security key should be a strong, randomly generated value that is kept secure and confidential. It should not be hardcoded in the source code or exposed in version control systems. Environment variables or secure key management systems are often used to handle this.

- Demonstration Focus
The primary goal in a demonstration context is often to showcase the functionality and workflow of using JWTs for authentication, rather than focusing on the intricacies of secure key management and token security.

- Vulnerability to Security Risks
Using a simplistic or publicly known key makes the token vulnerable to tampering and forgery. An attacker with knowledge of the key can generate their own valid tokens, gaining unauthorized access to protected resources.

- Best Practices for Production
For a production-grade system, it's crucial to implement robust security measures, including:

    - Using strong, unique keys for signing tokens.
    - Regularly rotating these keys.
    - Implementing additional token validation checks, like verifying the token issuer (iss) and the intended audience (aud).
    - Ensuring secure token storage and transmission, typically over HTTPS.
    - Handling token expiration and renewal securely.