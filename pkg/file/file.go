// Package file | File operation helper method
package file

import (
	"fmt"
	"github.com/disintegration/imaging"
	"github.com/gin-gonic/gin"
	"gohub/pkg/app"
	"gohub/pkg/auth"
	"gohub/pkg/helpers"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
)

// Put Save data to file
func Put(data []byte, to string) error {
	err := os.WriteFile(to, data, 0644)
	if err != nil {
		return err
	}
	return nil
}

// Exists Determine whether the file exists
func Exists(fileToCheck string) bool {
	if _, err := os.Stat(fileToCheck); os.IsNotExist(err) {
		return false
	}
	return true
}

func NameWithoutExtension(fileName string) string {
	return strings.TrimSuffix(fileName, filepath.Ext(fileName))
}

func SaveUploadAvatar(c *gin.Context, file *multipart.FileHeader) (string, error) {
	var avatar string

	publicPath := "public"
	dirName := fmt.Sprintf(
		"/uploads/avatars/%s/%s/",
		app.TimenowInTimezone().Format("2006/01/02"),
		auth.CurrentUID(c),
	)
	if err := os.MkdirAll(publicPath+dirName, 0755); err != nil {
		return avatar, err
	}

	fileName := randomNameFromUploadFile(file)
	avatarPath := publicPath + dirName + fileName
	if err := c.SaveUploadedFile(file, avatarPath); err != nil {
		return avatar, err
	}

	img, err := imaging.Open(avatarPath, imaging.AutoOrientation(true))
	if err != nil {
		return avatar, err
	}
	resizeAvatar := imaging.Thumbnail(img, 256, 256, imaging.Lanczos)
	resizeAvatarName := randomNameFromUploadFile(file)
	resizeAvatarPath := publicPath + dirName + resizeAvatarName
	err = imaging.Save(resizeAvatar, resizeAvatarPath)
	if err != nil {
		return avatar, err
	}

	err = os.Remove(avatarPath)
	if err != nil {
		return avatar, err
	}

	return dirName + resizeAvatarName, nil
}

func randomNameFromUploadFile(file *multipart.FileHeader) string {
	return helpers.RandomString(16) + filepath.Ext(file.Filename)
}
