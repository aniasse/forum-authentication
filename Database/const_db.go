package db

//this file allow the call of constants below in the sql/command functions

const (

	/************************/
	/*   Universal keys    */
	/*   --------------   */
	User_id = "user_id"
	Post_id = "post_id"

	/*************************/
	/*      users table     */
	/*    ------------     */
	User        = "users"
	Id_user     = "id_user"
	Username    = "username"
	Name        = "name"
	Surname     = "surname"
	Email       = "email"
	Password    = "password"
	Pp          = "pp"
	Pc          = "pc"
	Usersession = "usersession"

	/************************/
	/*      post table     */
	/*    ------------    */
	Post    = "posts"
	Id_post = "id_post"
	//User_id = "user_id"
	Title       = "title"
	Image       = "image"
	Description = "description"
	Categorie   = "categories"
	Hashtag     = "hashtag"
	Time        = "time"
	Date        = "date"

	/************************/
	/* post_reaction table */
	/* ------------------ */
	Post_reaction = "post_reactions"
	//User_id = "user_id"
	//Post_id ="post_id"
	Reaction = "reaction"

	/***************************/
	/* comment_reaction table */
	/* --------------------- */
	Comment_reaction = "comment_reactions"
	//User_id = "user_id"
	Comment_id = "comment_id"
	//Reaction = "reaction"

	/************************/
	/*    comment table    */
	/*    --------------  */
	Comment = "comments"
	//User_id = "user_id"
	//Post_id ="post_id"
	Id_comment = "id_comment"
	Content    = "content"
	// Time        = "time"
	// Date        = "date"

	/************************/
	/*    category table    */
	/*    --------------  */
	Category = "category"
)
