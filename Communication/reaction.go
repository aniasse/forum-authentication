package comm

import (
	"fmt"
	data "forum/Database"
)

type ReactionP struct {
	UserId   string
	PostId   string
	Reaction bool
}

type ReactionC struct {
	UserId    string
	CommentId string
	Reaction  bool
}

type Reacts []ReactionP
type ReactC []ReactionC

func (react_tab *Reacts) Get_reacPosts_data(database data.Db) error {
	// getting all database's reactions
	request1 := fmt.Sprintf("%s, %s, %s", data.User_id, data.Post_id, data.Reaction)
	rows_value, errow := database.GetData(request1, data.Post_reaction, "")
	if errow != nil {
		fmt.Println("⚠ ERROR ⚠ : Couldn't get reaction values from database ❌")
		fmt.Printf("⚠ : %v\n", errow)
		return errow
	}
	fmt.Println("✔ Reactions fetched from database")

	//storing them in a local structure
	temreact_tab := Reacts{}
	for rows_value.Next() {
		var temp ReactionP
		errscan := rows_value.Scan(&temp.UserId, &temp.PostId, &temp.Reaction)
		if errscan != nil {
			fmt.Println("⚠ ERROR ⚠ : Couldn't scan reactions values ❌")
			fmt.Printf("⚠ : %v\n", errscan)
			return errscan
		}
		temreact_tab = append(temreact_tab, temp)
	}
	*react_tab = temreact_tab

	fmt.Println("✅ reactions stored in local 🗄 structure successfully")
	return nil
}

func (react_tab *Reacts) React_post(database data.Db, id_post string, id_user string, react bool) error {
	// insert value in database
	request0 := fmt.Sprintf("(%s, %s, %s)", data.User_id, data.Post_id, data.Reaction)
	values := fmt.Sprintf("('%v', '%v', '%t')", id_user, id_post, react)
	err := database.INSERT(data.Post_reaction, request0, values)
	if err != nil {
		fmt.Printf("⚠ ERROR ⚠ : Couldn't insert reaction on post %v in database ❌\n", id_post)
		fmt.Printf("⚠ : %v\n", err)
		return err
	}
	fmt.Printf("✅ Reaction added to post %s with user %s and type of reaction: %t\n", id_post, id_user, react)

	return nil
}

/****************************************************************************************************************/
/****************************************************************************************************************/

func (react_tab *ReactC) React_comment(database data.Db, id_user string, id_comment string, react bool) error {
	// insert value in database
	request0 := fmt.Sprintf("(%s, %s, %s)", data.User_id, data.Comment_id, data.Reaction)
	values := fmt.Sprintf("('%v', '%v', '%t')", id_user, id_comment, react)
	err := database.INSERT(data.Comment_reaction, request0, values)
	if err != nil {
		fmt.Printf("⚠ ERROR ⚠ : Couldn't insert reaction on comment %v in database ❌\n", id_comment)
		fmt.Printf("⚠ : %v\n", err)
		return err
	}
	fmt.Printf("✅ Comment reaction added to comment %s with user %s and type of reaction: %t\n", id_comment, id_user, react)
	return nil
}

func (react_tab *ReactC) GetReact_comdata(database data.Db) error {
	// getting all database's reactions
	request1 := fmt.Sprintf("%s, %s, %s", data.User_id, data.Comment_id, data.Reaction)
	rows_value, errow := database.GetData(request1, data.Comment_reaction, "")
	if errow != nil {
		fmt.Println("⚠ ERROR ⚠ : Couldn't get comments reaction values from database ❌")
		fmt.Printf("⚠ : %v\n", errow)
		return errow
	}
	fmt.Println("✔ Comments reactions fetched from database")

	//storing them in a local structure
	temreact_tab := ReactC{}
	for rows_value.Next() {
		var temp ReactionC
		errscan := rows_value.Scan(&temp.UserId, &temp.CommentId, &temp.Reaction)
		if errscan != nil {
			fmt.Println("⚠ ERROR ⚠ : Couldn't scan comments reactions values ❌")
			fmt.Printf("⚠ : %v\n", errscan)
			return errscan
		}
		temreact_tab = append(temreact_tab, temp)
	}

	*react_tab = temreact_tab
	fmt.Println("✅ comments reactions stored in local 🗄 structure successfully")

	return nil
}
