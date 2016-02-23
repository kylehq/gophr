package main

import (
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func HandleImageNew(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	// Display New Image Form\
	RenderTemplate(w, r, "images/new", nil)
}

func HandleImageCreate(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	if r.FormValue("url") != "" {
		HandleImageCreateFromURL(w, r)
		return
	}
	HandleImageCreateFromFile(w, r)
}

func HandleImageCreateFromURL(w http.ResponseWriter, r *http.Request) {
	user := RequestUser(r)
	image := NewImage(user)
	image.Description = r.FormValue("description")
	err := image.CreateFromURL(r.FormValue("url"))
	if err != nil {
		if IsValidationError(err) {
			RenderTemplate(w, r, "images/new", map[string]interface{}{
				"Error":    err,
				"ImageURL": r.FormValue("url"),
				"Image":    image,
			})
			return
		}
		panic(err)
	}
	http.Redirect(w, r, "/?flash=Image+Uploaded+Successfully", http.StatusFound)
}

func HandleImageShow(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	image, err := globalImageStore.Find(params.ByName("imageID"))
	if err != nil {
		panic(err)
	}

	// 404
	if image == nil {
		http.NotFound(w, r)
		return
	}

	user, err := globalUserStore.Find(image.UserID)
	if err != nil {
		panic(err)
	}

	if user == nil {
		panic(fmt.Errorf("Could not find user %s", image.UserID))
	}

	RenderTemplate(w, r, "images/show", map[string]interface{}{
		"Image": image,
		"User":  user,
	})
}
