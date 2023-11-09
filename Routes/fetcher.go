package Route

import (
	"database/sql"
	"fmt"
	Err "forum/Authentication"
	Com "forum/Communication"
	tools "forum/tools"
	"net/http"
)

/*
GetAll_fromDB connects to database, retrieves from it informations
that will be display in the hime and index page
*/
func GetAll_fromDB(w http.ResponseWriter) {
	// connecting to database
	database.Doc, errd = sql.Open("sqlite3", "forum.db")
	if errd != nil {
		Err.Snippets(w, 500)
		return
	}
	//-------------- retrieving datas ---------------//
	//--1
	errGetPost := postab.GetPost_data(database)
	if errGetPost != nil {
		Err.Snippets(w, 500)
		return
	}
	//--2
	errGetComm := commtab.GetComment_data(database)
	if errGetComm != nil {
		Err.Snippets(w, 500)
		return
	}
	//--3
	errectabcomm := reactab_com.GetReact_comdata(database)
	if errectabcomm != nil {
		Err.Snippets(w, 500)
		return
	}
	//--4
	categos, err := Com.GetPost_categories(database)
	if err != nil {
		fmt.Printf("⚠ ERROR ⚠ : Couldn't get categories data from database\n")
		Err.Snippets(w, 500)
		return
	}
	//--5
	errectab := reactab.Get_reacPosts_data(database)
	if errectab != nil {
		fmt.Printf("⚠ ERROR ⚠ : Couldn't get reaction for display a from database\n")
		Err.Snippets(w, 500)
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

	//storing user's name and profil image in structures
	for i := range postab {
		username, name, surname, errGN := tools.GetName_byID(database, postab[i].UserId)
		Profil, errprof := tools.GetPic_byID(database, postab[i].UserId)

		if errprof != nil || errGN != nil {
			//sending metadata about the error to the servor
			Err.Snippets(w, 500)
			return
		}
		postab[i].Profil = Profil
		postab[i].Username = username
		postab[i].Name = name
		postab[i].Surname = surname
	}

	for i := range commtab {
		username, name, surname, errGN := tools.GetName_byID(database, commtab[i].UserId)
		Profil, errprof := tools.GetPic_byID(database, commtab[i].UserId)

		if errprof != nil || errGN != nil {
			//sending metadata about the error to the servor
			Err.Snippets(w, 500)
			return
		}
		commtab[i].Profil = Profil
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
