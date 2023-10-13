package Route

import (
	"errors"
	"fmt"
	"forum/tools"
	"io"
	"net/http"
	"os"
	"strings"
)

func Upload_mngmnt(w http.ResponseWriter, r *http.Request, size int, formFile string) (string, error) {
	//*checking the file 's size
	if r.Method == "POST" {
		maxsize := size * 1024 * 1024
		err := r.ParseMultipartForm(int64(maxsize))
		if err != nil {
			return "", errors.New("❌ could not allocted memory due to empty file in form")
		}

		file, header, err := r.FormFile(formFile)
		if err != nil { //!empty value sent wwhile submitting form
			fmt.Println("🚫 empty image")
			return "", nil
		}
		defer file.Close()

		if header.Size > int64(maxsize) { // Check if file size is greater than 5 MB
			fmt.Println("⚠ Image exceeds 20MB")
			return "", errors.New("file size exceeds  20MB limit")
		}
		fmt.Println("✅ image size checked")

		//*creating a copy of the uploaded in the server
		//!--checking extension validity
		if !tools.ValidExtension(strings.ToLower(header.Filename)) {
			fmt.Println("⚠ Wrong image extension")
			return "", errors.New("invalid extension")
		}
		uploaded, err := os.Create("templates/image_storage/" + header.Filename)
		if err != nil {
			fmt.Println("⚠ wrong image path")
			return "", err
		}

		defer uploaded.Close()

		//*Copying the uploaded file's content in the local one
		if _, err := io.Copy(uploaded, file); err != nil {
			fmt.Println("⚠ couldn't copy image in local")
			return "", err
		}

		return header.Filename, nil
	}
	return "", nil

}

func UploadImageUser(w http.ResponseWriter, r *http.Request, id string) {
	if r.Method != "POST" && r.Method != "GET" {
		fmt.Println("not allowed")
	}
	//restriction pour les différents chemins

	if r.Method == "POST" {
		fmt.Println("it's post")
		imageProfil, errProfil := Upload_mngmnt(w, r, 1, "profileImage")
		imageCover, errCover := Upload_mngmnt(w, r, 1, "murImage")
		if errProfil != nil || errCover != nil {
			fmt.Println("erreur cover ou profil", errProfil, errCover)
			return
		}

		if imageProfil != "" {
			imageProfil = "/static/image_storage/" + imageProfil
			errorUpdate := database.UPDATE("users", "pp='"+imageProfil+"'", "WHERE id_user='"+id+"'")
			if errorUpdate != nil {
				fmt.Println(errorUpdate)

			}
		}
		if imageCover != "" {
			imageCover = "/static/image_storage/" + imageCover
			errorCover := database.UPDATE("users", "pc='"+imageCover+"'", "WHERE id_user='"+id+"'")
			if errorCover != nil {
				fmt.Println(errorCover)
			}
		}

	} else if r.Method == "GET" {
		fmt.Println("affiche l'objet")
	}
	fmt.Println("i see u")
}
