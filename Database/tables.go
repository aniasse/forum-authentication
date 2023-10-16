package db

import (
	"database/sql"
	"fmt"
)

type Db struct {
	Doc *sql.DB
}

func Create_DB() (*sql.DB, error) {

	DB, err := sql.Open("sqlite3", "forum.db")
	if err != nil {
		return DB, err
	}
	return DB, err

}

func (database Db) Tables() {

	//----------------- 1 user table --------------------//
	User := `CREATE TABLE IF NOT EXISTS users (
		id_user TEXT PRIMARY KEY NOT NULL,
		username TEXT NOT NULL,
		name TEXT NOT NULL,
		surname TEXT NOT NULL,
		email TEXT NOT NULL,
		password TEXT NOT NULL,
		pp TEXT,
		PC TEXT,
		usersession TEXT
		);
		`
	_, errUser := database.Doc.Exec(User)
	if errUser != nil {
		fmt.Println("⚠ ERROR with table 'users' ⚠ :", errUser)
		return
	} else {
		fmt.Println("✅ 'users' table has been created in database succesfully")
	}

	//----------------- 2 posts table --------------------//
	Post := `
		CREATE TABLE IF NOT EXISTS posts (
			id_post TEXT PRIMARY KEY NOT NULL,
			user_id TEXT NOT NULL ,
			title TEXT NOT NULL,
			description TEXT DEFAULT "",
			image TEXT DEFAULT "",
			time TEXT NOT NULL,
			date TEXT NOT NULL,
			FOREIGN KEY (user_id) REFERENCES users (id_user)
			);
			`

	_, errPost := database.Doc.Exec(Post)
	if errPost != nil {
		fmt.Println("⚠ ERROR with table 'posts' ⚠ :", errPost)
		return
	} else {
		fmt.Println("✅ 'posts' table has been created in database succesfully")
	}

	//----------------- 3 comment table --------------------//
	Comment := `CREATE TABLE IF NOT EXISTS comments (
		id_comment TEXT PRIMARY KEY,
		user_id TEXT NOT NULL,
		post_id TEXT NOT NULL,
		content TEXT NOT NULL,
		username TEXT NOT NULL,
		surname TEXT NOT NULL,
		name TEXT NOT NULL,
		date TEXT NOT NULL,
		time TEXT NOT NULL,
		FOREIGN KEY(post_id) REFERENCES posts(id_post),
		FOREIGN KEY(user_id) REFERENCES users(id_user)
		
		);
		`
	_, errComment := database.Doc.Exec(Comment)
	if errComment != nil {
		fmt.Println("⚠ ERROR with table 'comment' ⚠ :", errComment)
		return
	}
	fmt.Println("✅ 'comments' table has been created in database succesfully")

	//----------------- 4 comment table --------------------//
	Post_reaction := `CREATE TABLE IF NOT EXISTS post_reactions (
		user_id TEXT NOT NULL,
		post_id TEXT NOT NULL,
		reaction BOOLEAN,
		FOREIGN KEY(post_id) REFERENCES posts(id_post),
		FOREIGN KEY(user_id) REFERENCES users(id_user)
		
		);
		`
	_, err1 := database.Doc.Exec(Post_reaction)
	if err1 != nil {
		fmt.Println("⚠ ERROR with table 'post_reactions' ⚠ :", err1)
		return
	}
	fmt.Println("✅ 'post_reactions' table has been created in database succesfully")

	//----------------- 5 comment table --------------------//
	Comment_reaction := `CREATE TABLE IF NOT EXISTS comment_reactions (
			user_id TEXT NOT NULL,
			comment_id TEXT NOT NULL,
			reaction BOOLEAN,
			FOREIGN KEY(user_id) REFERENCES users(id_user),
			FOREIGN KEY(comment_id) REFERENCES comments(id_comment)
			
			);
			`
	_, err2 := database.Doc.Exec(Comment_reaction)
	if err2 != nil {
		fmt.Println("⚠ ERROR with table 'comment_reactions' ⚠ :", err2)
		return
	}
	fmt.Println("✅ 'comment_reactions' table has been created in database succesfully")

	//----------------- 6 categorie table --------------------//
	Categories := `CREATE TABLE IF NOT EXISTS categories (
		user_id TEXT NOT NULL,
		post_id TEXT NOT NULL,
		category TEXT NOT NULL,
		FOREIGN KEY(user_id) REFERENCES users(id_user),
		FOREIGN KEY(post_id) REFERENCES posts(id_post)
		
		);
		`
	_, errcat := database.Doc.Exec(Categories)
	if errcat != nil {
		fmt.Println("⚠ ERROR with table 'categories' ⚠ :", errcat)
		return
	}

	fmt.Println("✅ 'categories' table has been created in database succesfully")
	//----------------- 7 session table --------------------//
	Session := `CREATE TABLE IF NOT EXISTS sessions (
	    user_id TEXT,
		id_session TEXT,
		expireat TEXT,
		FOREIGN KEY(user_id) REFERENCES users(id_user)
		);
			`
	_, errSession := database.Doc.Exec(Session)
	if errSession != nil {
		fmt.Println("⚠ ERROR with table 'sessions' ⚠ :", errSession)
		return
	} else {
		fmt.Println("✅ 'sessions' table has been created in database succesfully")
	}

}
