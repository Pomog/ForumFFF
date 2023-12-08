# ForumFFF
Online forum for MMORPG fans and friends!

## Project Overview

This project is one of the starting tasks from the Kood/Jõhvi Retraining Program. The description can be found [here](https://github.com/01-edu/public/tree/master/subjects/forum).

## Features

This is a Web Application written in Golang, utilizing only standard [Go libraries](https://pkg.go.dev/std) and the following external packages:
- `github.com/google/uuid v1.4.0`
- `github.com/mattn/go-sqlite3 v1.14.18`

### Database
- We are using the SQLite driver. However, our application implements The Repository pattern, and to switch databases, the `DatabaseInt` interface (located in the `repository` package) should be implemented for database systems.

### Template Rendering
-

###  HTTP server `routes` - HTTP request multiplexer
- sets up a basic HTTP server with route handlers for static files and various application endpoints, using the http.ServeMux as the multiplexer.
- The handlers are defined in the handler package, and the `routes` function is responsible for configuring the routing logic for the application.

## Unsolved Issues
1. Not optimized requests to the Database. Sometimes there are several requests per function or method.
2. No Middleware.
3. User's passwords are stored as strings in the Database.

## SQL schema
<img src="https://github.com/Pomog/ForumFFF/blob/25c4eb9089759d55ed4141969cdd4ca707d5ceca/SQL_schema.jpg?raw=true" alt="example" style="width:50%;">

## Authors
- [Denys Verves](https://github.com/TartuDen)
- [Yurii Panasiuk](https://github.com/pomog)
