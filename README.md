# FFForum
Online forum for MMORPG fans and friends!
![Forum Home Page](static/readme_images/image.png)

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

# Project Overview

This project encompasses a server-based application built using Go programming language. The structure and components of the application are as follows:

## Main Function (`main()`)

- Starts the server.
- Initializes the database connection.
- Sets up routes and handlers.
- Listens on a specified port for incoming HTTP requests.

## Database Initialization (`run()`)

- Establishes a database connection and sets up necessary tables.
- Registers custom types for encoding.
- Initializes repositories and handlers.

## Routing (`routes()`)

- Defines the application routes and corresponding handler functions.
- Handles serving static files and HTML templates.

## Handlers (`LoginHandler`, `RegisterHandler`, `HomeHandler`, `ThemeHandler`, `ErrorPage`)

- Perform various tasks based on HTTP requests:
  - Handling user login and registration.
  - Rendering templates.
  - Processing form data and performing validations.
  - Interacting with the database (creating users, threads, posts, etc.).
  - Redirecting users based on certain conditions.

## Form Handling (`Form` and associated methods)

- Represents HTML form data.
- Contains methods for form validation and error handling.
- Checks for required fields, validates field lengths and formats (like email format or password length), and more.

## Configuration (`AppConfig`)

- Holds application configuration settings.
- Manages caching, logging, and production mode status.
- Stores template caches, loggers, and other configuration data.

## Models (`User`, `Thread`, `Post`, `Votes`)

- Define structures representing entities in the application (users, threads, posts, etc.).
- Provide data structures for storing and working with information retrieved from the database.

## Template Rendering (`RendererTemplate`, `CreateTemplateCache`)

- Renders HTML templates using Go's `html/template` package.
- Caches templates for efficient rendering.
- Implements default data and error handling for rendering templates.

## Database Repository (`Repository`, `SqliteBDRepo`, `DatabaseInt`)

- Abstracts database interactions.
- Provides methods for database operations such as user presence check, user creation, thread creation, post creation, etc.
- Implements interfaces for different database interactions and functionality.

