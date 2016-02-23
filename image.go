package main

import (
	"io"
	"mime"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

const imageIDLength = 10

type Image struct {
	ID          string
	UserID      string
	Name        string
	Location    string
	Size        int64
	CreatedAt   time.Time
	Description string
}

type ImageStore interface {
	Save(image *Image) error
	Find(id string) (*Image, error)
	FindAll(offset int) ([]Image, error)
	FindAllByUser(user *User, offset int) ([]Image, error)
}

func NewImage(user *User) *Image {
	return &Image{
		ID:        GenerateID("img", imageIDLength),
		UserID:    user.ID,
		CreatedAt: time.Now(),
	}
}

// A map of accepted mime types and their file extension
var mimeExtensions = map[string]string{
	"image/png":  ".png",
	"image/jpeg": ".jpg",
	"image/gif":  ".gif",
}

func (image *Image) CreateFromURL(imageURL string) error {
	// Get the response from the URL
	response, err := http.Get(imageURL)
	if err != nil {
		return err
	}
	// Make sure we have a response
	if response.StatusCode != http.StatusOK {
		return errImageURLInvalid
	}
	defer response.Body.Close()

	// Ascertain the type of file we downloaded
	mimeType, _, err := mime.ParseMediaType(response.Header.Get("Content-Type"))
	if err != nil {
		return errInvalidImageType
	}

	// Get an extension for the file
	ext, valid := mimeExtensions[mimeType]
	if !valid {
		return errInvalidImageType
	}

	// Get a name from the URL
	image.Name = filepath.Base(imageURL)
	image.Location = image.ID + ext

	// Open a file at target location
	savedFile, err := os.Create("./data/images/" + image.Location)
	if err != nil {
		return err
	}
	defer savedFile.Close()

	// Copy the entire response to the output flie
	size, err := io.Copy(savedFile, response.Body)
	if err != nil {
		return err
	}

	// The returned value from io.Copy is the number of bytes copied
	image.Size = size

	// Save our image to the store
	return globalImageStore.Save(image)
}

func HandleImageCreateFromFile(w http.ResponseWriter, r *http.Request) {
	user := RequestUser(r)
	image := NewImage(user)
	image.Description = r.FormValue("description")
	file, headers, err := r.FormFile("file")

	// No file was uploaded.
	if file == nil {
		RenderTemplate(w, r, "images/new", map[string]interface{}{
			"Error": errNoImage,
			"Image": image,
		})
		return
	}

	// A file was uploaded, but an error occurred
	if err != nil {
		panic(err)
	}
	defer file.Close()

	// Check we have a valid image file type to save
	contentType := headers.Header.Get("Content-Type")
	if _, valid := mimeExtensions[contentType]; !valid {
		RenderTemplate(w, r, "images/new", map[string]interface{}{
			"Error": errInvalidImageType,
			"Image": image,
		})
		return
	}

	err = image.CreateFromFile(file, headers)
	if err != nil {
		RenderTemplate(w, r, "images/new", map[string]interface{}{
			"Error": err,
			"Image": image,
		})
		return
	}
	http.Redirect(w, r, "/?flash=Image+Uploaded+Successfully", http.StatusFound)
}

func (image *Image) CreateFromFile(file multipart.File, headers *multipart.FileHeader) error {
	// Move our file to an appropriate place, with an appropriate name
	image.Name = headers.Filename
	image.Location = image.ID + filepath.Ext(image.Name)

	// Open a file at the target location
	savedFile, err := os.Create("./data/images/" + image.Location)
	if err != nil {
		return err
	}
	defer savedFile.Close()

	// Copy the uploaded file to the target location
	size, err := io.Copy(savedFile, file)
	if err != nil {
		return err
	}

	image.Size = size

	// Save the image to the database
	return globalImageStore.Save(image)
}

func (image *Image) StaticRoute() string {
	return "/im/" + image.Location
}

func (image *Image) ShowRoute() string {
	return "/image/" + image.ID
}
