package comm

import (
	"errors"
	"fmt"
	data "forum/Database"
	tools "forum/tools"
	"strings"
)

func (Post_tab *Posts) Welcome_user(database data.Db, id_user string) error {
	if id_user != "" {

		date, time := tools.Time() //date and time
		categorie := []string{"education", "sport", "art & culture", "cinema", "health", "others"}
		// inserting value in database
		//-- formatting value's special chars
		value := `
Welcome to our forum!
As the first user, you are a pioneer. Feel free to explore topics that interest you and share your ideas.
Whether it's science, technology, art, or more, your voice is valuable. 
Contribute to make this community vibrant and enriching for everyone. We look forward to your contributions.
Thank you for being part of this adventure with us!

❌ PS: THIS POST WILL DISAPPEAR AFTER YOU CREATE ONE ❌
		  `
		value = strings.ReplaceAll(value, "'", "2@c86cb3")
		//-- formatting title's special chars
		title := "FIRST USER 🎉"

		//-- formatting image link's special chars
		image := "welcome.jpg"

		request0 := fmt.Sprintf("(%s, %s,%s, %s, %s, %s, %s)", data.Id_post, data.User_id, data.Title, data.Description, data.Image, data.Time, data.Date)
		values := fmt.Sprintf("('%s', '%s', '%s', '%s','%s', '%s', '%s')", "avamspost", "avams", title, value, image, time, date)
		err := database.INSERT(data.Post, request0, values)
		if err != nil {
			fmt.Println("⚠ ERROR ⚠ : Couldn't insert welcome post in database ❌")
			return err
		}

		request1 := fmt.Sprintf("(%s, %s,%s, %s, %s, %s, %s,%s, %s)", data.Id_user, data.Username, data.Name, data.Surname, data.Email, data.Password, data.Usersession, data.Pp, data.Pc)
		values1 := fmt.Sprintf("('%s', '%s', '%s', '%s','%s', '%s', '%s','%s', '%s')", "avams", "avams_team", "TEAM", "AVAMS", "avams@avams.com", "12345678", "sessiontestforwelcome", "../static/front-tools/images/profil.jpeg","../static/front-tools/images/mur.png")
		erruser := database.INSERT("users", request1, values1)
		if erruser != nil {
			fmt.Println("⚠ ERROR ⚠ : Couldn't insert avams profil in database ❌")
			return erruser
		}

		//inserting categories
		for i := range categorie {
			request0c := fmt.Sprintf("(%s, %s, %s)", data.Post_id, data.User_id, data.Category)
			valuesc := fmt.Sprintf("('%s', '%s','%s')", "avamspost", "avams", categorie[i])
			err = database.INSERT(data.Categorie, request0c, valuesc)
			if err != nil {
				fmt.Println("⚠ ERROR ⚠ : Couldn't insert 'welcome' categories in database ❌")
				return err
			}

		}

		fmt.Println("✅ welcome post has been created successfully")

		Post_tab = &Posts{}
		Post_tab.GetPost_data(database)
	}
	return nil
}

func (Post_tab *Posts) DeleteWelcome_user(database data.Db, id_user string) (error, bool) {
	if id_user != "" {
		//----------   checking the existence of the post   ------------------
		condicheck := "WHERE id_post = \"avamspost\""
		check, errcheck := database.Exist(data.Id_post, data.Post, condicheck)
		if errcheck != nil {
			fmt.Println("⚠ ERROR ⚠ : Couldn't get data from database in 'DeleteWelcome_user function' ❌")
			return errcheck, check
		}

		if !check {
			return nil, false
		}
		//-----------------------  end of checking   -----------------------

		condition := "WHERE id_post = \"avamspost\""
		errdel := database.DELETE("posts", condition)
		if errdel != nil {
			fmt.Println("⚠ ERROR ⚠ : Couldn't delete welcome post in database ❌")
			return errdel, false
		}

		condition1 := "WHERE id_user = \"avams\""
		errdel1 := database.DELETE("users", condition1)
		if errdel1 != nil {
			fmt.Println("⚠ ERROR ⚠ : Couldn't delete avams profil in database ❌")
			return errdel1, false
		}
		return nil, true
	}
	return errors.New("cannot delete without an Id user"), false
}
