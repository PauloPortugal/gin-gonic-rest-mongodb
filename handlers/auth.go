package handlers

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/PauloPortugal/gin-gonic-rest-mongodb/datastore"
	"github.com/PauloPortugal/gin-gonic-rest-mongodb/model"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/rs/xid"
	"github.com/spf13/viper"
)

type AuthHandler struct {
	ctx        context.Context
	cfg        *viper.Viper
	db         datastore.Users
	redisStore datastore.Redis
}

func NewAuthHandler(ctx context.Context, cfg *viper.Viper, store *datastore.UsersClient, redisStore *datastore.RedisClient) *AuthHandler {
	return &AuthHandler{
		ctx:        ctx,
		cfg:        cfg,
		db:         store,
		redisStore: redisStore,
	}
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
//         parameters:
//         - name: books_api_token
//           in: cookie
//           description: The user's session cookie/token
//     '400':
//         description: invalid input
//     '401':
//         description: unauthorised
//     '500':
//         description: internal server error
func (h *AuthHandler) SignIn(ctx *gin.Context) {
	var user model.User

	if h.isInvalidUser(ctx, &user) {
		return
	}

	// set session cookie
	sessionToken := xid.New().String()
	session := sessions.Default(ctx)
	session.Set("username", user.Username)
	session.Set("token", sessionToken)
	if err := session.Save(); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "User signed in"})
}

// swagger:operation POST /signout auth signUser
// Signs out a user
// ---
// responses:
//     '200':
//         description: Successful operation
//     '400':
//         description: invalid input
//     '401':
//         description: unauthorised
//     '500':
//         description: internal server error
func (h *AuthHandler) SignOut(ctx *gin.Context) {
	session := sessions.Default(ctx)
	session.Clear()
	session.Delete("username")
	session.Delete("token")
	session.Flashes()
	if err := session.Save(); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "User signed out"})
}

func (h *AuthHandler) isInvalidUser(ctx *gin.Context, reqUser *model.User) bool {
	if err := ctx.ShouldBindJSON(reqUser); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return true
	}

	_, err := h.db.Get(h.ctx, reqUser.Username, reqUser.Password)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid username or password",
		})
		return true
	}

	return false
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
