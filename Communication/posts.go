package comm

import (
	"fmt"
	"strings"

	"github.com/gofrs/uuid/v5"

	data "forum/Database"
	tools "forum/tools"
)

type Post struct {
	SessionId    string
	Profil       string
	Username     string
	Name         string
	Surname      string
	Title        string
	Content      string
	ImageLink    string
	PostId       string
	UserId       string
	Time         string
	Date         string
	Categorie    []string
	SessionReact string   // check ether the post is liked or disliked by the current user
	Like         []string // number of likes
	Dislike      []string //numer of dislkes
	Comment_tab  Comments
}

type Posts []Post

/*
	Create_post is a method which allows creates a post in the forum database by the userID, postID and the content.

Then it store the database informations in a local structure
*/
func (Post_tab *Posts) GetPost_data(database data.Db) error {
	// getting all posts from the user
	request := fmt.Sprintf("%s, %s, %s, %s, %s, %s, %s", data.Id_post, data.User_id, data.Title, data.Description, data.Image, data.Time, data.Date)

	rows_value, errow := database.GetData(request, data.Post, "") //retrieving datas
	if errow != nil {
		return errow
	}
	defer rows_value.Close()
	fmt.Println("✔ posts fetched from database")

	var temPost_tab Posts
	//storing retrieved datas in local structure
	for rows_value.Next() {
		var temp Post
		errscan := rows_value.Scan(&temp.PostId, &temp.UserId, &temp.Title, &temp.Content, &temp.ImageLink, &temp.Time, &temp.Date)
		if errscan != nil {
			fmt.Println("⚠ GetPost_data scan err ⚠ :", errscan)
			return errscan
		}
		//--formatting content's special chars
		temp.Content = strings.ReplaceAll(temp.Content, "2@c86cb3", "'")
		temp.Content = strings.ReplaceAll(temp.Content, "2#c86cb3", "`")

		//--formatting title's special chars
		temp.Title = strings.ReplaceAll(temp.Title, "2@c86cb3", "'")
		temp.Title = strings.ReplaceAll(temp.Title, "2#c86cb3", "`")

		//--formatting image link special chars
		temp.ImageLink = strings.ReplaceAll(temp.ImageLink, "2@c86cb3", "'")
		temp.ImageLink = strings.ReplaceAll(temp.ImageLink, "2#c86cb3", "`")

		temPost_tab = append(temPost_tab, temp)
	}
	fmt.Println("✅ posts 📊 has been stored in local 🗄 structure successfully")

	//reversing the tab
	for i, j := 0, len(temPost_tab)-1; i < j; i, j = i+1, j-1 {
		temPost_tab[i], temPost_tab[j] = temPost_tab[j], temPost_tab[i]
	}
	*Post_tab = temPost_tab
	return nil
}

func (Post_tab *Posts) Create_post(database data.Db, id_user string, categorie []string, value string, title string, image string) (string, error) {
	id_post, errp := uuid.NewV4() //id

	if id_user != "" && title != "" {
		//generating postID, date and time
		if errp != nil {
			fmt.Println("❌ Create_post ⚠ ERROR ⚠ : couldn't generate a unique post id")
			return "", errp
		}
		date, time := tools.Time() //date and time

		// inserting value in database
		//-- formatting value's special chars
		if value == "" {
			value = "0nbo6vda5l2udefa-v7a6i6l9a-b4lbefe-9ac6"
		} else {
			value = strings.ReplaceAll(value, "'", "2@c86cb3")
			value = strings.ReplaceAll(value, "`", "2#c86cb3")
		}

		//-- formatting title's special chars
		title = strings.ReplaceAll(title, "'", "2@c86cb3")
		title = strings.ReplaceAll(title, "`", "2#c86cb3")

		//-- formatting image link's special chars
		image = strings.ReplaceAll(image, "'", "2@c86cb3")
		image = strings.ReplaceAll(image, "`", "2#c86cb3")

		request0 := fmt.Sprintf("(%s, %s,%s, %s, %s, %s, %s)", data.Id_post, data.User_id, data.Title, data.Description, data.Image, data.Time, data.Date)
		values := fmt.Sprintf("('%s', '%s', '%s', '%s','%s', '%s', '%s')", id_post.String(), id_user, title, value, image, time, date)
		err := database.INSERT(data.Post, request0, values)
		if err != nil {
			fmt.Println("⚠ ERROR ⚠ : Couldn't insert post in database ❌")
			fmt.Printf("⚠ : %v\n", err)
			return "", err
		}

		//inserting categories
		for i := range categorie {
			request0c := fmt.Sprintf("(%s, %s, %s)", data.Post_id, data.User_id, data.Category)
			valuesc := fmt.Sprintf("('%s', '%s','%s')", id_post.String(), id_user, categorie[i])
			err = database.INSERT(data.Categorie, request0c, valuesc)
			if err != nil {
				fmt.Println("⚠ ERROR ⚠ : Couldn't insert categories in database ❌")
				fmt.Printf("⚠ : %v\n", err)
				return "", err
			}

		}

		fmt.Println("✅ post has been created successfully")

		Post_tab = &Posts{}
		Post_tab.GetPost_data(database)
	}

	return id_post.String(), nil
}
