package Route

import (
	"fmt"
	Err "forum/Authentification"
	db "forum/Database"
	"net/http"
)

// Reactpost_mngmnt handle all posts activities (creation, like and dislike)
func Reactpost_mngmnt(w http.ResponseWriter, r *http.Request, Post_Id string, React string) {
	// Checking if essential parameters are not empty
	if Id_user != "" && Post_Id != "" && React != "" {
		// Initializing react variable based on React string
		react := false
		if React == "true" {
			react = true
		}
		fmt.Println("[INFO] react bool: ", react) //debug

		// Constructing SQL condition string for querying the database
		condition := fmt.Sprintf("WHERE %s = '%s' AND %s = '%s'", db.User_id, Id_user, db.Post_id, Post_Id)
		// Fetching previous reaction data from the database
		prev_react, err := database.GetData(db.Reaction, db.Post_reaction, condition)

		if err == nil {
			PrevR, _ := db.Getelement(prev_react) // extract value
			fmt.Println("previous reaction : ", PrevR)
			if PrevR != "" {
				fmt.Println("already react")
				// ------- Handling cases for different combinations of previous and current reactions ------//
				switch {
				case PrevR == "true" && React == "false":
					fmt.Println("already liked but dislike")
					toset := fmt.Sprintf("%s = %t", db.Reaction, false)
					condition_upd := fmt.Sprintf("WHERE %s = '%s' AND %s = '%s'", db.User_id, Id_user, db.Post_id, Post_Id)
					errupdate := database.UPDATE(db.Post_reaction, toset, condition_upd)
					if errupdate != nil {
						fmt.Printf("⚠ ERROR ⚠ : Couldn't update user's reaction in database\n %s", errupdate)
						Err.Snippets(w, 500)
						return
					}

				case PrevR == "false" && React == "true":
					fmt.Println("already disliked but liked")
					toset := fmt.Sprintf("%s = %t", db.Reaction, true)
					condition_upd := fmt.Sprintf("WHERE %s = '%s' AND %s = '%s'", db.User_id, Id_user, db.Post_id, Post_Id)
					errupdate := database.UPDATE(db.Post_reaction, toset, condition_upd)
					if errupdate != nil {
						fmt.Printf("⚠ ERROR ⚠ : Couldn't update user's reaction in database\n %s", errupdate)
						Err.Snippets(w, 500)
						return
					}

				case PrevR == "false" && React == "false":
					fmt.Println("remove dislike")
					condel := fmt.Sprintf("WHERE %s = '%s' AND %s = '%s'", db.User_id, Id_user, db.Post_id, Post_Id)
					errdel := database.DELETE(db.Post_reaction, condel)
					if errdel != nil {
						fmt.Printf("⚠ ERROR ⚠ : Couldn't delete dislike reaction in database\n %s", errdel)
						Err.Snippets(w, 500)
						return
					}

				case PrevR == "true" && React == "true":
					fmt.Println("remove like")
					condel := fmt.Sprintf("WHERE %s = '%s' AND %s = '%s'", db.User_id, Id_user, db.Post_id, Post_Id)
					errdel := database.DELETE(db.Post_reaction, condel)

					if errdel != nil {
						fmt.Printf("⚠ ERROR ⚠ : Couldn't delete like reaction in database %s \n", errdel)
						Err.Snippets(w, 500)
						return
					}

				}
				// ---------------------------------------------------------------------------------------//

			} else { // it's a new post
				erreac := reactab.React_post(database, Post_Id, Id_user, react)
				if erreac != nil {
					fmt.Printf("⚠ ERROR ⚠ : %s ❌\n", erreac)
					Err.Snippets(w, 500)
					return
				}

				fmt.Println("--------------- 🟢🌐 data sent after reaction -----------------------")
			}

		} else {
			fmt.Println("⚠ ERROR ⚠ : there is no prev while fetching post reaction in reactpost func: ", err)
			return
		}

	} else {
		fmt.Println("error in Reactpost_mngmnt")
	}
}

/*----------------------------------------------------------------------------------------------------------------------------
----------------------------------------------------------------------------------------------------------------------------*/

// Reactcmnt_mngmnt handle all posts activities (creation, like and dislike)
func Reactcmnt_mngmnt(w http.ResponseWriter, r *http.Request, Comment_Id string, React string) {
	// Checking if essential parameters are not empty
	if Id_user != "" && Comment_Id != "" && React != "" {
		// Initializing react variable based on React string
		react := false
		if React == "true" {
			react = true
		}
		fmt.Println("[INFO] reactcom bool: ", react)

		// Constructing SQL condition string for querying the database
		condition := fmt.Sprintf("WHERE %s = '%s' AND %s = '%s'", db.User_id, Id_user, db.Comment_id, Comment_Id)

		// Fetching previous reaction data from the database
		prev_react, err := database.GetData(db.Reaction, db.Comment_reaction, condition)

		if err == nil {
			PrevR, _ := db.Getelement(prev_react) //extract value
			fmt.Println("previous reactioncomm : ", PrevR)
			if PrevR != "" {
				fmt.Println("already reacted com")

				// ------- Handling cases for different combinations of previous and current reactions ------//
				switch {
				case PrevR == "true" && React == "false":
					fmt.Println("already liked but dislike comm")
					toset := fmt.Sprintf("%s = %t", db.Reaction, false)
					condition_upd := fmt.Sprintf("WHERE %s = '%s' AND %s = '%s'", db.User_id, Id_user, db.Comment_id, Comment_Id)
					errupdate := database.UPDATE(db.Comment_reaction, toset, condition_upd)
					if errupdate != nil {
						fmt.Printf("⚠ ERROR ⚠ : Couldn't update user's comment reaction in database\n %s", errupdate)
						Err.Snippets(w, 500)
						return
					}

				case PrevR == "false" && React == "true":
					fmt.Println("already disliked but liked comm")
					toset := fmt.Sprintf("%s = %t", db.Reaction, true)
					condition_upd := fmt.Sprintf("WHERE %s = '%s' AND %s = '%s'", db.User_id, Id_user, db.Comment_id, Comment_Id)
					errupdate := database.UPDATE(db.Comment_reaction, toset, condition_upd)
					if errupdate != nil {
						fmt.Printf("⚠ ERROR ⚠ : Couldn't update user's comment reaction in database\n %s", errupdate)
						Err.Snippets(w, 500)
						return
					}

				case PrevR == "false" && React == "false":
					fmt.Println("remove dislike")
					condel := fmt.Sprintf("WHERE %s = '%s' AND %s = '%s'", db.User_id, Id_user, db.Comment_id, Comment_Id)
					errdel := database.DELETE(db.Comment_reaction, condel)
					if errdel != nil {
						fmt.Printf("⚠ ERROR ⚠ : Couldn't delete dislike comment reaction in database\n %s", errdel)
						Err.Snippets(w, 500)
						return
					}

				case PrevR == "true" && React == "true":
					fmt.Println("remove like")
					condel := fmt.Sprintf("WHERE %s = '%s' AND %s = '%s'", db.User_id, Id_user, db.Comment_id, Comment_Id)
					errdel := database.DELETE(db.Comment_reaction, condel)

					if errdel != nil {
						fmt.Printf("⚠ ERROR ⚠ : Couldn't delete like comment reaction in database %s \n", errdel)
						Err.Snippets(w, 500)
						return
					}

				}
				// ---------------------------------------------------------------------------------------//

			} else {
				println("react before insert: ", react)
				erreac := reactab_com.React_comment(database, Id_user, Comment_Id, react)
				if erreac != nil {
					fmt.Printf("⚠ ERROR ⚠ : %s ❌\n", erreac)
					Err.Snippets(w, 500)
					return
				}

				fmt.Println("--------------- 🟢🌐 data sent after reaction -----------------------")
			}

		} else {
			fmt.Println("⚠ ERROR ⚠ : there is no prev while fetching comment reaction in react_comment func: ", err)
			return
		}

	} else {
		fmt.Println("error in Reactcmnt_mngmnt")
	}
}
