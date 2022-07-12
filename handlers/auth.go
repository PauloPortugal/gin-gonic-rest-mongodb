package handlers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/PauloPortugal/gin-gonic-rest-mongodb/model"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

type AuthHandler struct {
	cfg *viper.Viper
}

func NewAuthHandler(cfg *viper.Viper) *AuthHandler {
	return &AuthHandler{cfg: cfg}
}

// swagger:operation POST /signin auth signUser
// Signs in a user
// ---
// consumes:
// - application/json
// produces:
// - application/json
// parameters:
// - name: credentials
//   in: body
//   name: User
//   description: The user credentials
//   schema:
//         "$ref": "#/definitions/User"
// responses:
//     '200':
//         description: Successful operation
//         schema:
//           type: array
//           items:
//                "$ref": "#/definitions/User"
//     '400':
//         description: invalid input
//     '401':
//         description: unauthorised
//     '500':
//         description: internal server error
func (h *AuthHandler) SignInHandler(ctx *gin.Context) {
	var user model.User

	if h.isInvalidUser(ctx, &user) {
		return
	}

	jwtOutput, err := h.generateJWT(user)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, jwtOutput)
}

func (h *AuthHandler) generateJWT(user model.User) (model.JWTOutput, error) {
	expirationTime := time.Now().Add(10 * time.Minute)

	claims := &model.Claims{
		Username: user.Username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(fmt.Sprint(h.cfg.Get("auth.secret"))))
	if err != nil {
		return model.JWTOutput{}, err
	}

	jwtOutput := model.JWTOutput{
		Token:   tokenString,
		Expires: expirationTime,
	}
	return jwtOutput, nil
}

func (h *AuthHandler) isInvalidUser(ctx *gin.Context, user *model.User) bool {
	if err := ctx.ShouldBindJSON(user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return true
	}

	//TODO: remove this hardcoded check
	if user.Username != "admin" || user.Password != "password" {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid username or password",
		})
		return true
	}
	return false
}
