# Web Forum (authentication)

## Table of Contents
1. [Description](#description)
2. [Objectives](#objectives)
3. [Technologies-Used](#technologies-used)
4. [Key-Features](#key-features)
5. [File system:](#file-system)
6. [Main-Directories](#main-directories)
7. [Instructions](#instructions)
8. [Authors](#authors)

### Description:
***
This project is a continuation of the web forum application developed using Go and SQLite for data storage. The forum allows users to communicate with each other, associate categories with posts, like and dislike posts adding image and comments, and filter posts based on different categories. User is also able register with his  **G🔴🟡GLE** or **GITHUB 🐱‍👤** account

## Objectives

The main objective of this project is to allow user registering with:

- GOOGLE
- GITHUB

## Technologies Used
***
- **GO**, also called Golang or Go language, is an open source programming language that Google developed. For more details check their [website](https://golang.org)
* **HTML**,The HyperText Markup Language is the standard markup language for documents designed to be displayed in at web browser.
+ **CSS**,Cascading Style Sheets, form a computer language that describes the presentation of HTML documents.  

+ **SQlite**, it is a database engine written in the C programming language. It is not a standalone app; rather, it is a library that software developers embed in their apps. As such, it belongs to the family of embedded databases.

+ **Docker**, Docker is a platform designed to help developers build, share, and run container applications. See more [here](https://www.docker.com)


## Key-Features

### 🔐 Authentication 

- registration by providing their email address, username, and password.
- Registration using **G🔴🟡GLE** or **GITHUB** account
- Creation of login sessions, allowing users to sign in to the forum.
- Use of cookies to manage sessions with an expiration date.

### 💻 Communication 

- Registered users can create posts adding image to a post and comments.
- Categories can be associated with posts.
- All posts and comments are visible to all users.

### 👍🏻👎🏻 Likes and Dislikes 

- Only registered users can like or dislike posts and comments.
- The number of likes and dislikes is visible to all users.

### ♻ Filtering 

- Users can filter posts by categories, created posts, and posts liked by them.

## 🐋 Docker

The application utilizes Docker for managing the development environment.

## file-system  
The project is organized into multiple directories for better source code organization.
```go
.
|
|____📂Authentication  
|    |-----------📄auth_api_tools.go
|    |-----------📄BD.go
|    |-----------📄helpers.go
|    |-----------📄Session_com.go
|
|____📂Communication
|    |-----------📄categories.go
|    |-----------📄comment.go
|    |-----------📄posts.go
|    |-----------📄reaction.go
|    |-----------📄welcome.go
|
|____📂Database
|    |-----------📄commands.go
|    |-----------📄const_db.go
|    |-----------📄Init_db.go
|    |-----------📄tables.go
|
|____📂Handlers
|    |-----------📄handle.go
|
|____📂Models
|    |-----------📄ERD.go
|
|____📂Routes
|    |-----------📄authentication.go
|    |-----------📄communication.go
|    |-----------📄fetcher.go
|    |-----------📄filter.go
|    |-----------📄github.go
|    |-----------📄google.go
|    |-----------📄index.go
|    |-----------📄processing.go
|    |-----------📄profil-cmt.go
|    |-----------📄profil-fav.go
|    |-----------📄profil.go
|    |-----------📄react.go
|    |-----------📄tools_communication.go
|    |-----------📄upload.go
|
|____🗄📂templates
|    |____📂front-tools
|    |    |___📂css
|    |    |   |-----------🎨aboutus.css
|    |    |   |-----------🎨error.css
|    |    |   |-----------🎨logand.css
|    |    |   |-----------🎨style.css
|    |    |
|    |    |___📂images
|    |    |   |-----------📷(front-images...)
|    |    |
|    |    |___📂Js
|    |    |   |-----------📒script.js
|    |    |   |-----------📒upload.js
|    |    |
|    |___📂image_storage
|    |   |-----(💾images uploaded from website...)
|    |    
|    |--------------------------------------------------📝createaccout.html
|    |--------------------------------------------------📝error.html
|    |--------------------------------------------------📝filter_com.html
|    |--------------------------------------------------📝filter_fav.html
|    |--------------------------------------------------📝footer.html
|    |--------------------------------------------------📝head.html
|    |--------------------------------------------------📝home.html
|    |--------------------------------------------------📝home.html
|    |--------------------------------------------------📝index.html
|    |--------------------------------------------------📝main.html
|    |--------------------------------------------------📝navbar.html
|    |--------------------------------------------------📝profil.html
|    |--------------------------------------------------📝register.html
|
|____📂tools
|    |--------------🔧Id_toold.go
|    |--------------🔧standard_funcs.go
|    |--------------🔧validity.go
|
|----🐋Dockerfile
|
|----🛢forum.db
|
|----⚙go.mod
|
|----🔃go.sum
|
|----📝main.go
|
|----📜README.md
|
|----⚙script.sh
|
|----📜task guider.md
.

```

## Main-Directories
`Authentification` directory contains files related to user authentication, including database management and helper functions.

`Communication` directory is dedicated to communication, including the creation of comments and posts.  

`Database` directory includes files related to database management, such as SQL commands, constants, database initialization, and table definitions.

 `handlers` directory may contain request handler functions or other features related to processing user requests.

`models` directory may contain data models or other elements related to the application's domain model.

`Routes` directory may contain route management files to direct user requests to the appropriate functionalities.

`templates` directory may contain HTML templates used to generate web pages.

`tools` directory may contain tools, utilities, or reusable libraries.

Some directories, such as `css`, `images`, and `JS`, are organized into subdirectories for better structuring of frontend resources.  
Feel free to explore each individual subdirectory for more details on its contents and specific purpose.


## Instructions

1. how to install:
   ```bash
   $ git clone https://learn.zone01dakar.sn/git/mthiaw/forum-authentication.git
   $ cd forum-authentication
   ```

2. how to enter the forum:
   ```go
   go run main.go
   ✅ database has been created successfully
   ✅ 'users' table has been created in database succesfully
   ✅ 'posts' table has been created in database succesfully
   ✅ 'comments' table has been created in database succesfully
   ✅ 'post_reactions' table has been created in database succesfully   
   ✅ 'comment_reactions' table has been created in database succesfully
   ✅ 'categories' table has been created in database succesfully       
   ✅ 'sessions' table has been created in database succesfully
   📡----------------------------------------------------📡
   |                                                      |
   | 🌐 Server has started at "http://localhost:8080 🟢" |
   |                                                      |
   📡----------------------------------------------------📡
   ```
   go to the link and follow instructions

3. how to Use Docker to build an image of the application.
   ```bash
   docker build -t forum-app .
   # step 1 : build image
   docker build -t <image_name>
   # step 2 : build and run container with built image
   docker run -d --name <container_name> -p <port>: <container_port> <image_name>
   # enter the container and list items with this command a
   docker exec <container_name> ls -l 
   ```
or execute directly the `script.sh` file
```bash
$sh script.sh
```
  
## Authors:
### front-end team :

- [**A-boubakrine DIALLO** (*aboubakdiallo*)](https://learn.zone01dakar.sn/git/aboubakdiallo)
- [**V-incent Félix NDOUR** (*vindour*)](https://learn.zone01dakar.sn/git/vindour)
### back-end team :
- [**A-dama NIASSE** (*aniasse*)](https://learn.zone01dakar.sn/git/aniasse)
* [**M-asseck THIAW** (*mthiaw - **captain***)](https://learn.zone01dakar.sn/git/mthiaw)
+ [**S-eynabou NIANG** (*sniang*)](https://learn.zone01dakar.sn/git/sniang)

#### *@Licensed by AVAMS👏 team*
