package handlers

import (
	"backend/db"
	"backend/helpers"
	"backend/models"
	"context"
	"mime/multipart"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type SignupForm struct {
	Username        string                `form:"username" binding:"required"`
	Password        string                `form:"password" binding:"required"`
	ConfirmPassword string                `form:"confirm_password" binding:"required"`
	Avatar          *multipart.FileHeader `form:"avatar"`
}

const maxAvatarSize = 2 << 20 // 2 MiB

func (h *Handler) Signup(c *gin.Context) {
	var form SignupForm

	if err := c.ShouldBind(&form); err != nil {
		c.String(http.StatusBadRequest, "Compila tutti i campi.")
		return
	}

	if form.Password != form.ConfirmPassword {
		c.String(http.StatusBadRequest, "Le password non coincidono.")
		return
	}

	collection := h.DB.Collection(db.UsersCollectionName)
	result := collection.FindOne(context.TODO(), bson.M{"username": form.Username}).Err()
	if result == mongo.ErrNoDocuments {

		user := models.User{
			Username: form.Username,
			Password: helpers.GeneratePasswordHash(form.Password),
			Elo:      0,
			Status:   "offline",
		}

		if form.Avatar != nil {

			if form.Avatar.Size > maxAvatarSize {
				c.String(http.StatusBadRequest, "Immagine troppo grande (max. 2MB)")
				return
			}

			file, err := form.Avatar.Open()
			if err != nil {
				c.String(http.StatusInternalServerError, err.Error())
				return
			}
			buff := make([]byte, 512)
			_, err = file.Read(buff)
			if err != nil {
				c.String(http.StatusInternalServerError, err.Error())
				return
			}

			filetype := http.DetectContentType(buff)
			splitFiletype := strings.Split(filetype, "/")
			mimetype, extension := splitFiletype[0], splitFiletype[1]
			if mimetype == "image" && (extension == "png" || extension == "jpeg" || extension == "gif" || extension == "bmp" || extension == "tiff") {

				//newFileName := uuid.New().String() + "." + extension
				newFileName := primitive.NewObjectID().Hex() + "." + extension

				err = c.SaveUploadedFile(form.Avatar, "avatar/"+newFileName)
				if err != nil {
					c.String(http.StatusInternalServerError, err.Error())
					return
				}
				user.Avatar = newFileName
			} else {
				c.String(http.StatusBadRequest, "L'avatar fornito non è un'immagine valida.")
				return
			}
		}

		_, err := collection.InsertOne(context.TODO(), user)
		if err != nil {
			c.String(http.StatusInternalServerError, err.Error())
			return
		}

	} else if result == nil {
		c.String(http.StatusBadRequest, "L'username inserito è già utilizzato.")
		return
	} else {
		c.String(http.StatusBadRequest, result.Error())
		return
	}

}
