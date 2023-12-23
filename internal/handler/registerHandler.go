package handler

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/Pomog/ForumFFF/internal/forms"
	"github.com/Pomog/ForumFFF/internal/models"
	"github.com/Pomog/ForumFFF/internal/renderer"
)

// RegisterHandler handles both GET and POST requests for the registration page.
func (m *Repository) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		var emptyRegistration models.User
		data := make(map[string]interface{})
		data["registrationData"] = emptyRegistration
		renderer.RendererTemplate(w, "register.page.html", &models.TemplateData{
			Form: forms.NewForm(nil),
			Data: data,
		})

	} else if r.Method == http.MethodPost {
		// Parse the form data, including files Need to Set Upper limit for DATA
		err := r.ParseMultipartForm(1 << 20)
		if err != nil {
			setErrorAndRedirect(w, r, dbErrorUserPresent, "/error-page")
			return
		}

		// Create a User struct with data from the HTTP request form
		registrationData := models.User{
			FirstName: r.FormValue("firstName"),
			LastName:  r.FormValue("lastName"),
			UserName:  r.FormValue("nickName"),
			Email:     strings.ToLower(r.FormValue("emailRegistr")),
			Password:  r.FormValue("passwordReg"),
			Picture:   r.FormValue("avatar"),
		}

		// Create a new form instance based on the HTTP request's PostForm
		form := forms.NewForm(r.PostForm)

		// Validation checks for required fields and their specific formats and lengths
		form.Required("firstName", "lastName", "nickName", "emailRegistr", "passwordReg")
		form.First_LastName_Min_Max_Len("firstName", 3, 12, r)
		form.First_LastName_Min_Max_Len("lastName", 3, 12, r)
		form.First_LastName_Min_Max_Len("nickName", 3, 12, r)
		form.EmailFormat("emailRegistr", r)
		form.First_LastName_Min_Max_Len("emailRegistr", 10, 30, r)
		form.PassFormat("passwordReg", 6, 15, r)

		// Check if the form data is valid; if not, render the home page with error messages
		if !form.Valid() {
			data := make(map[string]interface{})
			data["registrationData"] = registrationData
			renderer.RendererTemplate(w, "register.page.html", &models.TemplateData{
				Form: form,
				Data: data,
			})
			return
		}

		// Check if User is Present in the DB, ERR should be handled
		userAlreadyExist, err := m.DB.UserPresent(registrationData.UserName, registrationData.Email)
		if err != nil {
			setErrorAndRedirect(w, r, userAlreadyExistsMsg, "/error-page")
			return
		}

		if userAlreadyExist {
			setErrorAndRedirect(w, r, "User with such Email OR NickName Already Exist", "/error-page")
		} else {
			// Get the file from the form data
			file, handler, errFileGet := r.FormFile("avatar")
			if errFileGet != nil {
				setErrorAndRedirect(w, r, fileReceivingErrorMsg, "/error-page")
				return
			}
			defer file.Close()

			// Validate file size (1 MB limit)
			if handler.Size > 1<<20 {
				form.Errors.Add("avatar", "File size should be below 1 MB")
				data := make(map[string]interface{})
				data["registrationData"] = registrationData
				renderer.RendererTemplate(w, "register.page.html", &models.TemplateData{
					Form: form,
					Data: data,
				})
				return
			}

			// Validate file type (must be an image)
			contentType := handler.Header.Get("Content-Type")
			if contentType != "image/jpeg" && contentType != "image/png" && contentType != "image/gif" {
				form.Errors.Add("avatar", "Invalid file type. Only JPEG, PNG, and GIF images are allowed")
				data := make(map[string]interface{})
				data["registrationData"] = registrationData
				renderer.RendererTemplate(w, "register.page.html", &models.TemplateData{
					Form: form,
					Data: data,
				})
				return
			}

			// Create a new file in the "static/ava" directory
			newFilePath := filepath.Join("static/ava", handler.Filename)
			newFile, errFileCreate := os.Create(newFilePath)
			if errFileCreate != nil {
				setErrorAndRedirect(w, r, fileCreatingErrorMsg, "/error-page")
				return
			}
			defer newFile.Close()

			// Copy the uploaded file to the new file
			_, err = io.Copy(newFile, file)
			if err != nil {
				setErrorAndRedirect(w, r, fileSavingErrorMsg, "/error-page")
				return
			}

			registrationData.Picture = path.Join("/", newFilePath)

			err := m.DB.CreateUser(registrationData)
			if err != nil {
				setErrorAndRedirect(w, r, "DB Error func CreateUser", "/error-page")
				return
			}

			message := fmt.Sprintf("User %s is registered", registrationData.UserName)
			fmt.Println(message)
			// helper.SendEmail(m.App.ServerEmail, message)

			http.Redirect(w, r, "/login", http.StatusSeeOther)
		}

	} else {
		http.Error(w, "No such method", http.StatusMethodNotAllowed)
		return
	}

}
