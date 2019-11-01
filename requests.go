/*this package will do all the operations regarding

the http request methods*/
package http

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"reflect"
	"strings"

	uuid "github.com/satori/go.uuid"
)

////this will get any image sent with
///sent in the body of a request given
///the imazge key
func FormFile(r *http.Request, path, url string) (string, *ErrResponse) {
	_, image, err := ReadRequestFile(r, "image", path, url)
	if err != nil {
		log.Println(err)
		return "", ErrInvalidRequest
	}

	return image, nil

}

///////get the file from the data
/////filename is the key the formfile was sent with
////path is where you want to store the file
////base url is the base http you want to be accesssing the file with
////example of baseUrl="https://ndunyu.co.ke/images"
func ReadRequestFile(r *http.Request, filename string, storagePath string, BaseUrl string) (string, string, error) {
	_ = r.ParseMultipartForm(32 << 20)
	file, handler, err := r.FormFile(filename)

	//ex, err := os.Executable()
	if err != nil {
		log.Println(err)
		return "", "", err

	}
	defer file.Close()
	if _, err := os.Stat(storagePath); os.IsNotExist(err) {
		_ = os.MkdirAll(storagePath, os.ModePerm)
	}

	imageName := uuid.NewV4().String() + filepath.Ext(handler.Filename)

	imagePath := storagePath + imageName

	f, err := os.OpenFile(imagePath, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return "", "", err

	}

	result := strings.Replace(imagePath, storagePath, "", -1)

	imageUrl := BaseUrl + result
	//image_name =   handler.Filename
	defer f.Close()
	_, _ = io.Copy(f, file)
	return imageName, imageUrl, nil

}


//takes in a pointer and reads to it the request body sent
func RequestBody(r *http.Request, item interface{}) *ErrResponse {
	if reflect.ValueOf(item).Kind() != reflect.Ptr {
		return ErrInvalidRequest

	}
	err := json.NewDecoder(r.Body).Decode(&item)
	if err != nil {
		///sentry.CaptureException(err)
		return ErrInvalidRequest
	}
	defer r.Body.Close()
	return nil

}