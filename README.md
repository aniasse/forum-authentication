# Web Forum

## Table of Contents
1. [Objectives](#objectives)
2. [Technologies-Used](#technologies-Used)
3. [Key-Features](#key-Features)
4. [Docker-Image](#docker-Usage)
5. [Main-Directories](#main-Directories)
6. [Instructions](#instructions)


This project is a continuation of the web forum application developed using Go and SQLite for data storage. The forum allows users to communicate with each other, associate categories with posts, like and dislike posts adding image and comments, and filter posts based on different categories.

## Objectives

The main objectives of this project include:

- Allowing user registration and authentication.
- Adding image to a post
- Enabling users to create posts and comments.
- Associating categories with posts.
- Managing likes and dislikes of posts and comments.
- Providing a mechanism to filter posts by categories, created posts, and posts liked by the logged-in user.

## Technologies Used

- Go: Used for backend development.
- SQLite: Used for managing the database.
- HTML/CSS: Used for creating the user interface.
- Docker: Used for containerizing the application.

## Key Features

### Authentication

- Users can register by providing their email address, username, and password.
- Creation of login sessions allows users to sign in to the forum.
- Use of cookies to manage sessions with an expiration date.

### Communication

- Registered users can create posts adding image to a post and comments.
- Categories can be associated with posts.
- All posts and comments are visible to all users.

### Likes and Dislikes

- Only registered users can like or dislike posts and comments.
- The number of likes and dislikes is visible to all users.

### Filtering

- Users can filter posts by categories, created posts, and posts liked by them.

## Docker Usage

The application utilizes Docker for managing the development environment. You can create a Docker image to run the application.

# Project Directory Structure
The project is organized into multiple directories for better source code organization.

      ├── Authentification 

         |   ├── BD.go

         |   ├── helpers.go

      ├──    Communication 

         |   ├── categories.go

         |   ├── comment.go

         |   ├── posts.go

         |   ├── reaction.go

         |   ├── welcome.go

      ├── Handlers

         |   ├── handle.go

      ├── Database 

         |   ├── commands.go

         |   ├── const_db.go

         |   ├── Init_db.go

         |   ├── tables.go

      ├── templates

         |   ├── front-tools

         |   ├── image_storage

         |   ├── createacount.html

         |   ├── error.html

         |   ├── footer.html

         |   ├── head.html

         |   ├── image_storage

         |   ├── home.html

         |   ├── index.html

         |   ├── main.html

         |   ├── navbar.hmtl

         |   ├── profil.html

         |   ├── register.html

      ├── tools

         |   ├── Id_tools.go

         |   ├── uuid.go

         |   ├── Validity.go

         |   ├── standard_funcs.go


      ├── models 

         |   ├── ERD.png

      ├── Routes 

         |   ├── authentification.go

         |   ├── communication.go

         |   ├── filter.go

         |   ├── index.go

         |   ├── react.go

         |   ├── tools_communication.go

         |   ├──  upload.go


## Main-Directories

### Authentification

The    `Authentification` directory contains files related to user authentication, including database management and helper functions.

### Communication

The `Communication` directory is dedicated to communication, including the creation of comments and posts.

### Database

The `Database` directory includes files related to database management, such as SQL commands, constants, database initialization, and table definitions.

### handlers

The `handlers` directory may contain request handler functions or other features related to processing user requests.

### models

The `models` directory may contain data models or other elements related to the application's domain model.

### Routes

The `Routes` directory may contain route management files to direct user requests to the appropriate functionalities.

### templates

The `templates` directory may contain HTML templates used to generate web pages.

### tools

The `tools` directory may contain tools, utilities, or reusable libraries.

## Subdirectories

Some directories, such as `css`, `images`, and `JS`, are organized into subdirectories for better structuring of frontend resources.

Feel free to explore each individual subdirectory for more details on its contents and specific purpose.

This directory structure is designed for efficient organization of source code and project resources. Please refer to individual subdirectories for more information on their content and usage.

## Instructions

1.
***
A little intro about how to install:
```
$ Clone the project from https://learn.zone01dakar.sn/git/aboubakdiallo/forum-image-upload.git
$ cd forum
```

2. Use Docker to build an image of the application.

   ```bash
   docker build -t forum-app .



*@Licensed by team AVAMS*
