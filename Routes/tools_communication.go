package Route

import (
	"database/sql"
	"fmt"
	Com "forum/Communication"
	db "forum/Database"
	tools "forum/tools"
	"html/template"
	"net/http"
	"strings"
)

/*
Display_mngmnt connects to database, retrieves from it informations
that will be display in the hime and index page
*/
func Display_mngmnt(w http.ResponseWriter, r *http.Request) {
	// connecting to database
	database.Doc, errd = sql.Open("sqlite3", "forum.db")
	if errd != nil {
		w.WriteHeader(http.StatusInternalServerError)
		error_file := template.Must(template.ParseFiles("templates/error.html"))
		error_file.Execute(w, "500")
		return
	}
	//-------------- retrieving datas ---------------//
	//--1
	errGetPost := postab.GetPost_data(database)
	if errGetPost != nil {
		w.WriteHeader(http.StatusInternalServerError)
		error_file := template.Must(template.ParseFiles("templates/error.html"))
		error_file.Execute(w, "500")
		return
	}
	//--2
	errGetComm := commtab.GetComment_data(database)
	if errGetComm != nil {
		w.WriteHeader(http.StatusInternalServerError)
		error_file := template.Must(template.ParseFiles("templates/error.html"))
		error_file.Execute(w, "500")
		return
	}
	//--3
	errectabcomm := reactab_com.GetReact_comdata(database)
	if errectabcomm != nil {
		fmt.Printf("⚠ ERROR ⚠ : Couldn't get comments reaction for display from database\n")
		w.WriteHeader(http.StatusInternalServerError)
		error_file := template.Must(template.ParseFiles("templates/error.html"))
		error_file.Execute(w, "500")
		return
	}
	//--4
	categos, err := Com.GetPost_categories(database)
	if err != nil {
		fmt.Printf("⚠ ERROR ⚠ : Couldn't get categories data from database\n")
		w.WriteHeader(http.StatusInternalServerError)
		error_file := template.Must(template.ParseFiles("templates/error.html"))
		error_file.Execute(w, "500")
		return
	}
	//--5
	errectab := reactab.Get_reacPosts_data(database)
	if errectab != nil {
		fmt.Printf("⚠ ERROR ⚠ : Couldn't get reaction for display a from database\n")
		w.WriteHeader(http.StatusInternalServerError)
		error_file := template.Must(template.ParseFiles("templates/error.html"))
		error_file.Execute(w, "500")
		return
	}
	//--------------------------------------------------------------------//
	// storing the session's id
	for i := range postab {
		postab[i].SessionId = Id_user
	}
	for i := range commtab {
		commtab[i].SessionId = Id_user
	}

	//storing user's name in structures
	for i := range postab {
		username, name, surname := tools.GetName_byID(database, postab[i].UserId)
		postab[i].Username = username
		postab[i].Name = name
		postab[i].Surname = surname
	}

	for i := range commtab {
		username, name, surname := tools.GetName_byID(database, commtab[i].UserId)
		commtab[i].Username = username
		commtab[i].Name = name
		commtab[i].Surname = surname
	}

	//storing the reactions in corresponding comments
	for i := range commtab {
		for j := range reactab_com {
			if commtab[i].CommentId == reactab_com[j].CommentId {
				switch reactab_com[j].Reaction {
				case true:
					commtab[i].Likecomm = append(commtab[i].Likecomm, "true")
					if commtab[i].SessionId == reactab_com[j].UserId {
						commtab[i].SessionReact = "true"
					}

				case false:
					commtab[i].Dislikecomm = append(commtab[i].Dislikecomm, "false")
					if commtab[i].SessionId == reactab_com[j].UserId {
						commtab[i].SessionReact = "false"
					}
				}
			}
		}
	}

	//storing the comments in corresponding posts
	for i := range postab {
		for j := range commtab {
			if postab[i].PostId == commtab[j].PostId {
				postab[i].Comment_tab = append(postab[i].Comment_tab, commtab[j])
			}
		}
	}

	//storing the categories in corresponding posts
	for i := range postab {
		for j := range categos {
			if postab[i].PostId == categos[j].PostId {
				postab[i].Categorie = append(postab[i].Categorie, categos[j].Category)
			}
		}
	}

	//storing the reactions in corresponding posts
	for i := range postab {
		for j := range reactab {
			if postab[i].PostId == reactab[j].PostId {
				switch reactab[j].Reaction {
				case true:
					postab[i].Like = append(postab[i].Like, "true")
					if postab[i].SessionId == reactab[j].UserId {
						postab[i].SessionReact = "true"
					}

				case false:
					postab[i].Dislike = append(postab[i].Dislike, "false")
					if postab[i].SessionId == reactab[j].UserId {
						postab[i].SessionReact = "false"
					}
				}
			}
		}
	}

}

// CreateP_mngmnt handles user's post activity
func CreateP_mngmnt(w http.ResponseWriter, r *http.Request, categorie []string, title string, content string, image string, redirect string) {

	errpost := postab.Create_post(database, Id_user, categorie, title, content, image)
	if errpost != nil {
		fmt.Printf("⚠ ERROR ⚠ : Couldn't create post from user %s ❌\n", Id_user)
		w.WriteHeader(http.StatusInternalServerError)
		error_file := template.Must(template.ParseFiles("templates/error.html"))
		error_file.Execute(w, "500")
		return
	}
	fmt.Println("post created with content ", content)

	//*Getting the id post according to the content for html relative link
	//---formatting content to escape special chars
	content = strings.ReplaceAll(content, "'", "2@c86cb3")
	content = strings.ReplaceAll(content, "`", "2#c86cb3")
	//---fetching id post in database
	condition := fmt.Sprintf("WHERE %s = '%s'", db.Description, content)
	Idpost, err1 := database.GetData(db.Id_post, db.Post, condition)
	Idpost_got, err2 := db.Getelement(Idpost)
	if err1 != nil && err2 != nil {
		http.Redirect(w, r, redirect+"#"+Idpost_got, http.StatusSeeOther)
	} else { //no id found in database, post creation encountered a problem
		http.Redirect(w, r, redirect, http.StatusSeeOther)
	}
}

// CreateC_mngmnt handles user's comment activity
func CreateC_mngmnt(w http.ResponseWriter, r *http.Request, Id_post string, newcomment string) {
	errcomm := commtab.Create_comment(database, Id_user, Id_post, newcomment)
	if errcomm != nil {
		fmt.Printf("⚠ ERROR ⚠ : Couldn't create comment in post %s from user %s ❌\n", Id_post, Id_user)
		w.WriteHeader(http.StatusInternalServerError)
		error_file := template.Must(template.ParseFiles("templates/error.html"))
		error_file.Execute(w, "500")
		return
	}
}

// ReplyC_mngmnt handles user's comment reply activity
func ReplyC_mngmnt(w http.ResponseWriter, r *http.Request, Id_post string, Id_comment string, Id_user string, replycomm string) {

	com_owner_username := tools.GetName_bycomment(database, Id_comment)
	fmt.Println("name comm ", com_owner_username)
	reply := fmt.Sprintf("@%v %v", com_owner_username, replycomm)

	errcomm := commtab.Create_comment(database, Id_user, Id_post, reply)
	if errcomm != nil {
		fmt.Printf("⚠ ERROR ⚠ : Couldn't create reply to comment %s , in on post %s from user %s ❌\n", Id_comment, Id_post, Id_user)
		w.WriteHeader(http.StatusInternalServerError)
		error_file := template.Must(template.ParseFiles("templates/error.html"))
		error_file.Execute(w, "500")
		return
	}

}
