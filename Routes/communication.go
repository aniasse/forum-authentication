package Route

import (
	"fmt"
	"html/template"
	"net/http"

	Err "forum/Authentification"
	Com "forum/Communication"
	db "forum/Database"
	tools "forum/tools"
)

type Res struct {
	CurrentN  string
	CurrentSN string
	CurrentUN string
	Postab    Com.Posts
}

var (
	postab      Com.Posts    // posts local storage
	commtab     Com.Comments // comments local storage
	reactab     Com.Reacts   //posts reactions local storage
	reactab_com Com.ReactC   // comments reactions local storage
	database    db.Db        //database local storage
	errd        error        // manage errors
)

var Id_user string

// "ca1eaeb1-4bf9-4507-a5f6-281d927dda5a"

/*
Communications handles user's posts, comments and reactions
it can only be reached by using method POST or GET
*/

func Communication(w http.ResponseWriter, r *http.Request, Id string, redirect string, route_info string) {
	Id_user = Id
	//!--checking the http request
	if r.Method != "POST" && r.Method != "GET" {
		fmt.Printf("⚠ ERROR ⚠ : cannot access to that page by with mode other than GET & POST must log out to reach it ❌")
		Err.Snippets(w, 405)
		return
	}

	Display_mngmnt(w, r) //display all values in the forum database
	fmt.Println(postab)
	if len(postab) == 0 {
		errwel := postab.Welcome_user(database, Id_user)
		if errwel != nil {
			fmt.Printf("⚠ ERRWEL ⚠ :%s ❌", errwel)
			Err.Snippets(w, 500)
			return
		}
	} else {
		errdelwel:=postab.DeleteWelcome_user(database, Id_user)
		if errdelwel != nil {
			fmt.Printf("⚠ ERRDELWEL ⚠ :%s ❌", errdelwel)
			Err.Snippets(w, 500)
			return
		}
	}

	//?------------ client sent a request-----------------
	if r.Method == "POST" {
		//--------retrieving form values ----------
		fmt.Println("--------------------------------------------")
		fmt.Println("             " + route_info + "Form values" + "                  ")
		fmt.Println("--------------------------------------------")

		//--ID
		Id_post := r.FormValue("postid")
		fmt.Println("[INFO] ID post: ", Id_post) //debug

		Id_postR := r.FormValue("Rpostid")
		fmt.Println("[INFO] ID postREc: ", Id_postR) //debug

		Id_comment := r.FormValue("comId")
		fmt.Println("[INFO] ID comment: ", Id_comment) //debug

		Id_commentR := r.FormValue("Rcomid")
		fmt.Println("[INFO] ID commentR: ", Id_commentR) //debug

		//-----title
		Title := r.FormValue("title")
		fmt.Println("[INFO] Post title: ", Title) //debug

		//---text content
		content := r.FormValue("post_content")
		fmt.Println("[INFO] content: ", content) //debug

		newcomment := r.FormValue("newcomment")
		fmt.Println("[INFO] comment: ", newcomment) //debug

		replycomm := r.FormValue("replycomm")
		fmt.Println("[INFO] reply comment: ", replycomm) //debug

		//-------------------------image's link----------------------------
		Image, errimage := Upload_mngmnt(w, r)
		fmt.Println("[INFO] Post image link: ", Image) //debug
		//---------------------------------------------------------------

		//----Reactions
		React := r.FormValue("react")
		fmt.Println("[INFO] react: ", React) //debug

		Reactcomm := r.FormValue("reactcomm")
		fmt.Println("[INFO] reactcomm: ", Reactcomm) //debug

		//-----submit buttons
		Subpost := r.FormValue("subpost")
		fmt.Println("[INFO] subpost: ", Subpost) //debug

		Subcomm := r.FormValue("subcomm")
		fmt.Println("[INFO] subcomm: ", Subcomm) //debug

		subreply := r.FormValue("subreply")
		fmt.Println("[INFO] subreply: ", subreply) //debug

		//------categories
		education := r.FormValue("education")
		sport := r.FormValue("sport")
		art_culture := r.FormValue("art_culture")
		cinema := r.FormValue("cinema")
		health := r.FormValue("health")
		others := r.FormValue("others")

		categorie := []string{education, sport, art_culture, cinema, health, others}
		var tempc []string
		for _, v := range categorie {
			if v != "" {
				tempc = append(tempc, v)
			}
		}
		categorie = tempc
		fmt.Println("[INFO] categorie: ", categorie) //debug

		fmt.Println("--------------------------------------------")
		//-----------end of retrieving form value----------

		switch {

		//*-create post case:
		case Id_user != "" && Subpost != "":
			//verifying the request method
			if r.Method != "POST" {
				fmt.Printf("⚠ ERROR ⚠ : cannot access to that page by with mode other than POST ❌")
				Err.Snippets(w, 400)
				return
			}
			//checking Id_user validity
			if tools.IsnotExist_user(Id_user, database) {
				Err.Snippets(w, 400)
				return
			}
			//checking Title's validity
			if Title == "" {
				fmt.Printf("⚠ ERROR ⚠ : Couldn't create post from user %s due to empty title ❌\n", Id_user)
				Err.Snippets(w, 400)
				return
			}
			//checking content's validity
			if content == "" {
				fmt.Printf("⚠ ERROR ⚠ : Couldn't create post from user %s due to empty content ❌\n", Id_user)
				Err.Snippets(w, 400)
				return
			}
			//checking categore's validity
			if len(categorie) < 1 { //user did not select a categorie
				fmt.Printf("⚠ ERROR ⚠ : Couldn't create post from user %s due to missing category❌\n", Id_user)
				Err.Snippets(w, 400)
				return
			}

			if tools.IsInvalid(content) || tools.IsInvalid(Title) || len(Title) > 25 { //found only spaces,newlines in the input or chars number limit exceeded
				fmt.Printf("⚠ ERROR ⚠ : Couldn't create post from user %s due to invalid input ❌\n", Id_user)
				Err.Snippets(w, 400)
				return
			}

			if errimage != nil {
				fmt.Printf("⚠ ERROR ⚠ : Couldn't create post from user %s, error encoutered while uploading image\n%s ❌\n", Id_user, errimage)
				Err.Snippets(w, 400)
				return
			}
			CreateP_mngmnt(w, r, categorie, content, Title, Image, redirect)

			// create comment case:
		case Id_user != "" && Subcomm != "" && Id_post != "":
			//!--checking Id_user and Id_post validity
			if tools.IsnotExist_user(Id_user, database) || tools.IsnotExist_Post(Id_post, database) {
				Err.Snippets(w, 400)
				return
			}

			//!--checking if the comment is empty
			if newcomment == "" {
				fmt.Printf("⚠ ERROR ⚠ : Couldn't create comment from user %s due to empty content ❌\n", Id_user)
				Err.Snippets(w, 400)
				return
			}

			//!--checking the comment validity
			if tools.IsInvalid(newcomment) { //found only spaces or newlines in the input
				fmt.Printf("⚠ ERROR ⚠ : Couldn't create comment in post %s from user %s due to invalid input ❌\n", Id_post, Id_user)
				Err.Snippets(w, 400)
				return
			}

			if r.Method != "POST" {
				fmt.Printf("⚠ ERROR ⚠ : cannot access to that page by with mode other than POST ❌")
				Err.Snippets(w, 405)
				return
			}
			CreateC_mngmnt(w, r, Id_post, newcomment)
			http.Redirect(w, r, redirect+"#"+Id_post, http.StatusSeeOther)

			//*reply comment case:
		case Id_user != "" && Id_post != "" && Id_comment != "" && subreply != "":
			//!--checking Id_user, Id_post and Id_comment validity
			if tools.IsnotExist_user(Id_user, database) || tools.IsnotExist_Post(Id_post, database) || tools.IsnotExist_Comment(Id_comment, database) {
				Err.Snippets(w, 400)
				return
			}

			//!--checking if the comment is empty
			if replycomm == "" {
				fmt.Printf("⚠ ERROR ⚠ : Couldn't create comment reply from user %s due to empty content ❌\n", Id_user)
				Err.Snippets(w, 400)
				return
			}

			//!--checking the comment validity
			if tools.IsInvalid(replycomm) { //found only spaces or newlines in the input
				fmt.Printf("⚠ ERROR ⚠ : Couldn't create comment in post %s from user %s due to invalid input ❌\n", Id_post, Id_user)
				Err.Snippets(w, 400)
				return
			}

			if r.Method != "POST" {
				fmt.Printf("⚠ ERROR ⚠ : cannot access to that page by with mode other than POST ❌")
				Err.Snippets(w, 405)
				return
			}
			ReplyC_mngmnt(w, r, Id_post, Id_comment, Id_user, replycomm)
			http.Redirect(w, r, redirect+"#"+Id_post, http.StatusSeeOther)

			//* reactpost case:
		case Id_user != "" && Id_postR != "" && React != "":
			//!--checking id_user and id_post validity
			if tools.IsnotExist_user(Id_user, database) || tools.IsnotExist_Post(Id_postR, database) {
				Err.Snippets(w, 400)
				return
			}

			if r.Method != "POST" {
				fmt.Printf("⚠ ERROR ⚠ : cannot access to that page by with mode other than POST ❌")
				Err.Snippets(w, 405)
				return
			}
			Reactpost_mngmnt(w, r, Id_postR, React)
			http.Redirect(w, r, redirect+"#"+Id_postR, http.StatusSeeOther) //refreshing the page after data processing

			//*reactcomment case
		case Id_user != "" && Id_commentR != "" && Reactcomm != "":
			//!--checking id_user and id_post validity
			if tools.IsnotExist_user(Id_user, database) || tools.IsnotExist_Comment(Id_commentR, database) {
				Err.Snippets(w, 400)
				return
			}

			if r.Method != "POST" {
				fmt.Printf("⚠ ERROR ⚠ : cannot access to that page by with mode other than POST ❌")
				Err.Snippets(w, 405)
				return
			}
			Reactcmnt_mngmnt(w, r, Id_commentR, Reactcomm)
			http.Redirect(w, r, redirect+"#"+Id_commentR, http.StatusSeeOther) //refreshing the page after data processing

			//default: just display datas

		} // end switch case

	} //?------------ end of request treatment-----------------

	file, errf := template.ParseFiles("templates/home.html")
	if errf != nil {
		//sending metadata about the error to the servor
		fmt.Printf("⚠ ERROR ⚠ parsing --> %v\n", errf)
		Err.Snippets(w, 500)
		return
	}

	// user's name
	current_username, current_surname, current_name := tools.GetName_byID(database, Id_user)
	//struct to execute
	final := Res{
		CurrentUN: current_username,
		CurrentSN: current_surname,
		CurrentN:  current_name,
		Postab:    postab,
	}

	//sending data to html
	errexc := file.Execute(w, final)
	if errexc != nil {
		//sending metadata about the error to the servor
		fmt.Printf("⚠ ERROR ⚠ executing file --> %v\n", errexc)
		return
	}
	fmt.Println("--------------- 🟢🌐 home data sent -----------------------") //debug

}
