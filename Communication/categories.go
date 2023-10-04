package comm

import (
	"fmt"
	data "forum/Database"
)

type Catego struct {
	UserId   string
	PostId   string
	Category string
}

func GetPost_categories(database data.Db) ([]Catego, error) {
	// getting all posts from the user
	request := fmt.Sprintf("%s, %s, %s", data.Post_id, data.User_id, data.Category)

	rows_value, errow := database.GetData(request, data.Categorie, "") //retrieving datas
	if errow != nil {
		return []Catego{}, errow
	}
	defer rows_value.Close()
	fmt.Println("✔ categories fetched from database")

	var catego_tab []Catego
	
	//storing retrieved datas in local structure
	for rows_value.Next() {
		var temp Catego
		errscan := rows_value.Scan(&temp.PostId, &temp.UserId, &temp.Category)
		if errscan != nil {
			fmt.Println("scanerr ", errscan)
			return []Catego{}, errscan
		}
		catego_tab = append(catego_tab, temp)
	}
	return catego_tab, nil
}
