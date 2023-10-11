package Route

import (
	"fmt"
	Err "forum/Authentification"
	db "forum/Database"
	tools "forum/tools"
	"net/http"
	"strings"
)

// CreateP_mngmnt handles user's post activity
func CreateP_mngmnt(w http.ResponseWriter, r *http.Request, categorie []string, title string, content string, image string, redirect string) {

	idPost_toReplace, errpost := postab.Create_post(database, Id_user, categorie, title, content, image)
	if errpost != nil {
		fmt.Printf("⚠ ERROR ⚠ : Couldn't create post from user %s ❌\n", Id_user)
		Err.Snippets(w, 500)
		return
	}
	fmt.Println("post created with content ", content)

	//*Getting the id post according to the content for html relative link
	//---formatting content to escape special chars
	if content == "" {
		content = idPost_toReplace
	} else {
		content = strings.ReplaceAll(content, "'", "2@c86cb3")
		content = strings.ReplaceAll(content, "`", "2#c86cb3")
	}

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
		Err.Snippets(w, 500)
		return
	}
}

// ReplyC_mngmnt handles user's comment reply activity
func ReplyC_mngmnt(w http.ResponseWriter, r *http.Request, Id_post string, Id_comment string, Id_user string, replycomm string) {

	com_owner_username, errGN := tools.GetName_bycomment(database, Id_comment)
	if errGN != nil {
		//sending metadata about the error to the servor
		Err.Snippets(w, 500)
		return
	}
	
	fmt.Println("name comm ", com_owner_username)
	reply := fmt.Sprintf("@%v %v", com_owner_username, replycomm)

	errcomm := commtab.Create_comment(database, Id_user, Id_post, reply)
	if errcomm != nil {
		fmt.Printf("⚠ ERROR ⚠ : Couldn't create reply to comment %s , in on post %s from user %s ❌\n", Id_comment, Id_post, Id_user)
		Err.Snippets(w, 500)
		return
	}

}
