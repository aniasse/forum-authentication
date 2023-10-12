package comm

import (
	"fmt"
	data "forum/Database"
	tools "forum/tools"
	"strings"

	"github.com/gofrs/uuid/v5"
)

type Comment struct {
	SessionId    string
	Profil       string
	PostId       string
	UserId       string
	CommentId    string
	Username     string
	Name         string
	Surname      string
	Content      string
	Time         string
	Date         string
	SessionReact string
	Likecomm     []string
	Dislikecomm  []string
}

type Comments []Comment

// GetComment_data retrieves all datas related to comments from databse
func (Comm_tab *Comments) GetComment_data(database data.Db) error {
	request := fmt.Sprintf("%s, %s, %s, %s, %s, %s", data.Post_id, data.User_id, data.Id_comment, data.Content, data.Time, data.Date)

	rows_value, errow := database.GetData(request, data.Comment, "") //retrieving datas
	if errow != nil {
		fmt.Println("⚠ ERROR ⚠ : Couldn't get comment values from database ❌")
		fmt.Printf("⚠ : %v\n", errow)
		return errow
	}
	fmt.Println("✔ comments fetched from database")

	var temCom Comments
	//storing retrieved datas in local structure
	for rows_value.Next() {
		var temp Comment
		errscan := rows_value.Scan(&temp.PostId, &temp.UserId, &temp.CommentId, &temp.Content, &temp.Time, &temp.Date)
		if errscan != nil {
			fmt.Println("⚠ ERROR ⚠ : Couldn't scan comments values ❌")
			fmt.Printf("⚠ : %v\n", errscan)
			return errscan
		}
		temp.Content = strings.ReplaceAll(temp.Content, "2@c86cb3", "'")
		temp.Content = strings.ReplaceAll(temp.Content, "2#c86cb3", "`")
		temCom = append(temCom, temp)
	}

	*Comm_tab = temCom

	fmt.Println("✅ comments 📊 has been stored in local 🗄 structure successfully")

	return nil
}

/*
	Create_comment is a method which allows creates a comment in the forum database by the userID, postID , commentID and the content.

Then it store the database informations in a local structure
*/
func (Comm_tab *Comments) Create_comment(database data.Db, id_user string, id_post string, Content string) error {
	//generating commentID, date and time
	id_comment, errp := uuid.NewV4() //id
	if errp != nil {
		fmt.Println("❌ Create_comment ⚠ ERROR ⚠ : couldn't generate a unique comment id")
		return errp
	}
	date, time := tools.Time() //date and time
	username, surname, name, errGN := tools.GetName_byID(database, id_user)
	if errGN != nil {
		return errGN
	}

	// inserting value in database
	Content = strings.ReplaceAll(Content, "'", "2@c86cb3")
	Content = strings.ReplaceAll(Content, "`", "2#c86cb3")
	request0 := fmt.Sprintf("(%s, %s, %s, %s, %s, %s, %s, %s, %s)", data.Post_id, data.User_id, data.Id_comment, data.Username, data.Surname, data.Name, data.Content, data.Time, data.Date)
	values := fmt.Sprintf("('%v', '%v', '%v', '%s', '%s', '%s', '%s', '%s', '%s')", id_post, id_user, id_comment, username, surname, name, Content, time, date)
	err := database.INSERT(data.Comment, request0, values)
	if err != nil {
		fmt.Println("⚠ ERROR ⚠ : Couldn't insert comment in database ❌")
		fmt.Printf("⚠ : %v\n", err)
		return err
	}

	fmt.Printf("✅ comment %s has been added to database successfully\n", id_comment.String())

	return nil
}
