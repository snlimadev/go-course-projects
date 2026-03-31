## About the Project

This project is an event booking REST API that stores data in a SQLite database.

First of all, you need to create a `.env` file in the project's root directory
and add the following environment variables:

```env
GO_ENV=development
JWT_ACCESS_SECRET=your_jwt_access_token_secret
JWT_REFRESH_SECRET=your_jwt_refresh_token_secret
JWT_ACCESS_DURATION_MINUTES=60
JWT_REFRESH_DURATION_DAYS=30
```

After running `go run .` to start the server locally, you can explore the API
endpoints via Swagger UI at http://localhost:8080/docs.