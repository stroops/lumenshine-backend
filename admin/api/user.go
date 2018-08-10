package api

import (
	"github.com/Soneso/lumenshine-backend/admin/db"
	mw "github.com/Soneso/lumenshine-backend/admin/middleware"
	"github.com/Soneso/lumenshine-backend/admin/models"
	cerr "github.com/Soneso/lumenshine-backend/icop_error"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/labstack/echo"

	"github.com/Soneso/lumenshine-backend/admin/route"
	"github.com/Soneso/lumenshine-backend/helpers"

	"golang.org/x/crypto/bcrypt"
)

const (
	//UserRoutePrefix is the prefix for the user group. We need this in order to get all the routes for this base url
	UserRoutePrefix = "user"
)

//init setup all the routes for the users handling
func init() {
	route.AddRoute("GET", "/details/:id", User, []string{"Users", "Administrators"}, "user_details", UserRoutePrefix)
	route.AddRoute("GET", "/list", Users, []string{"Administrators"}, "user_list", UserRoutePrefix)
	route.AddRoute("GET", "/user_data", UserData, []string{}, "user_data", UserRoutePrefix)
	route.AddRoute("POST", "/register", Register, []string{"Administrators"}, "user_register", UserRoutePrefix)
	route.AddRoute("POST", "/edit", EditUser, []string{"Users", "Administrators"}, "user_edit", UserRoutePrefix)
	route.AddRoute("POST", "/activate", Activate, []string{"Administrators"}, "user_activate", UserRoutePrefix)
	route.AddRoute("POST", "/setgroups", SetGroups, []string{"Administrators"}, "user_setgroups", UserRoutePrefix)
}

//AddUserRoutes adds all the routes for the user handling
func AddUserRoutes(rg *gin.RouterGroup) {
	for _, r := range route.GetRoutesForPrefix(UserRoutePrefix) {
		f := r.HandlerFunc.(func(uc *mw.AdminContext, c *gin.Context))
		rg.Handle(r.Method, r.Prefix+r.Path, mw.UseAdminContext(f, r.Name))
	}
}

// UserAuthData of the logged in user
type UserAuthData struct {
	UserID      int      `json:"id"`
	Email       string   `json:"email"`
	FirstName   string   `json:"firstname"`
	LastName    string   `json:"lastname"`
	Phone       string   `json:"phone"`
	Active      bool     `json:"active"`
	Groups      []string `json:"groups"`
	IsAdmin     bool     `json:"is_admin"`
	IsDeveloper bool     `json:"is_developer"`
}

//UserData returns the current logged in user info
func UserData(uc *mw.AdminContext, c *gin.Context) {

	u := uc.User

	var userData UserAuthData
	userData.UserID = u.ID
	userData.Email = u.Email
	userData.FirstName = u.FirstName
	userData.LastName = u.LastName
	userData.Phone = u.Phone
	userData.Active = u.Active
	userData.Groups = u.Groups
	userData.IsAdmin = helpers.StringInSliceI("Administrators", userData.Groups)
	userData.IsDeveloper = helpers.StringInSliceI("Developers", userData.Groups)

	c.JSON(http.StatusOK, echo.Map{
		"data": userData,
	})
}

//UserResponse holds the user reponse info
type UserResponse struct {
	ID        int      `json:"id"`
	Email     string   `json:"email"`
	FirstName string   `json:"firstname"`
	LastName  string   `json:"lastname"`
	Phone     string   `json:"phone"`
	Active    bool     `json:"active"`
	Groups    []string `json:"groups"`
}

func getUpdatedBy(c *gin.Context) string {
	return c.GetString("userID")
}

//RegisterRequest new user information
type RegisterRequest struct {
	Email     string   `form:"email" json:"email" validate:"required,icop_email"`
	Password  string   `form:"password" json:"password" validate:"required"`
	FirstName string   `form:"firstname" json:"firstname" validate:"required"`
	LastName  string   `form:"lastname" json:"lastname" validate:"required"`
	Phone     string   `form:"phone" json:"phone" validate:"required,icop_phone"`
	Active    bool     `form:"active" json:"active"`
	Groups    []string `form:"groups" json:"groups" validate:"required"`
}

//RegisterResponse after registration
type RegisterResponse struct {
	ID     int    `form:"id" json:"id"`
	Email  string `form:"email" json:"email"`
	Active bool   `form:"active" json:"active"`
}

//Register creates new use in the db
func Register(uc *mw.AdminContext, c *gin.Context) {
	var rr RegisterRequest
	if err := c.Bind(&rr); err != nil {
		c.JSON(http.StatusBadRequest, cerr.LogAndReturnError(uc.Log, err, cerr.ValidBadInputData, cerr.BindError))
		return
	}

	if valid, validErrors := cerr.ValidateStruct(uc.Log, rr); !valid {
		c.JSON(http.StatusBadRequest, validErrors)
		return
	}

	pwd, err := bcrypt.GenerateFromPassword([]byte(rr.Password), 14)
	if err != nil {
		c.JSON(http.StatusInternalServerError, cerr.LogAndReturnError(uc.Log, err, "Error encrypting password", cerr.GeneralError))
		return
	}

	user := &models.AdminUser{
		Email:     rr.Email,
		Password:  string(pwd),
		Forename:  rr.FirstName,
		Lastname:  rr.LastName,
		Phone:     rr.Phone,
		Active:    rr.Active,
		UpdatedBy: getUpdatedBy(c)}

	err = db.RegisterUser(user, rr.Groups)
	if err != nil {
		c.JSON(http.StatusInternalServerError, cerr.LogAndReturnError(uc.Log, err, "Error registering user", cerr.GeneralError))
		return
	}

	c.JSON(http.StatusOK, &RegisterResponse{
		ID:     user.ID,
		Email:  user.Email,
		Active: user.Active})
}

//ActivateRequest for activate/deactivate
type ActivateRequest struct {
	ID     int  `form:"id" json:"id" validate:"required"`
	Active bool `form:"active" json:"active"`
}

//ActivateResponse after activation/deactivation
type ActivateResponse struct {
	ID     int  `form:"id" json:"id"`
	Active bool `form:"active" json:"active"`
}

//Activate - activates or deactivates a user
func Activate(uc *mw.AdminContext, c *gin.Context) {
	var ar ActivateRequest
	if err := c.Bind(&ar); err != nil {
		c.JSON(http.StatusBadRequest, cerr.LogAndReturnError(uc.Log, err, cerr.ValidBadInputData, cerr.BindError))
		return
	}

	if valid, validErrors := cerr.ValidateStruct(uc.Log, ar); !valid {
		c.JSON(http.StatusBadRequest, validErrors)
		return
	}

	activeAdministrators, err := db.AllActiveAdministrators()
	if err != nil {
		c.JSON(http.StatusInternalServerError, cerr.LogAndReturnError(uc.Log, err, "Could not read the users", cerr.GeneralError))
		return
	}

	if len(activeAdministrators) == 1 {
		if activeAdministrators[0].ID == ar.ID && !ar.Active {
			c.JSON(http.StatusBadRequest, cerr.LogAndReturnError(uc.Log, err, "At least an administrator must be active", cerr.GeneralError))
			return
		}
	}

	err = db.ActivateUser(ar.ID, ar.Active, getUpdatedBy(c))
	if err != nil {
		c.JSON(http.StatusInternalServerError, cerr.LogAndReturnError(uc.Log, err, "Error activating user", cerr.GeneralError))
		return
	}

	c.JSON(http.StatusOK, &ActivateResponse{
		ID:     ar.ID,
		Active: ar.Active})
}

//UsersResponse holds the user list
type UsersResponse struct {
	Users []UserResponse `form:"users" json:"users"`
}

//Users reads the list of users from the db
func Users(uc *mw.AdminContext, c *gin.Context) {
	dbUsers, err := db.AllUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, cerr.LogAndReturnError(uc.Log, err, "Error reading users", cerr.GeneralError))
		return
	}

	users := make([]UserResponse, len(dbUsers))
	for index, user := range dbUsers {
		groups := make([]string, 0)
		for _, g := range user.R.UserAdminUsergroups {
			groups = append(groups, g.R.Group.Name)
		}
		users[index] = UserResponse{
			ID:        user.ID,
			Email:     user.Email,
			FirstName: user.Forename,
			LastName:  user.Lastname,
			Phone:     user.Phone,
			Active:    user.Active,
			Groups:    groups}
	}

	c.JSON(http.StatusOK, &UsersResponse{Users: users})
}

//User gets the specified user by id
func User(uc *mw.AdminContext, c *gin.Context) {
	idParam := c.Param("id")

	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, cerr.LogAndReturnError(uc.Log, err, "Invalid user id specified", cerr.GeneralError))
		return
	}

	user, err := db.GetUserByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, cerr.LogAndReturnError(uc.Log, err, "Could not read user", cerr.GeneralError))
		return
	}

	if user == nil {
		c.JSON(http.StatusBadRequest, cerr.LogAndReturnError(uc.Log, err, "Cannot find user in db", cerr.GeneralError))
		return
	}

	c.JSON(http.StatusOK,
		&UserResponse{
			ID:        user.ID,
			Email:     user.Email,
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Phone:     user.Phone,
			Active:    user.Active,
			Groups:    user.Groups})
}

//EditRequest - edits user information
type EditRequest struct {
	ID        int    `form:"id" json:"id" validate:"required"`
	Password  string `form:"password" json:"password"`
	FirstName string `form:"firstname" json:"firstname"`
	LastName  string `form:"lastname" json:"lastname"`
	Phone     string `form:"phone" json:"phone"`
	Active    bool   `form:"active" json:"active"`
}

//EditUser - edits the user's properties
func EditUser(uc *mw.AdminContext, c *gin.Context) {
	var editRequest EditRequest
	if err := c.Bind(&editRequest); err != nil {
		c.JSON(http.StatusBadRequest, cerr.LogAndReturnError(uc.Log, err, cerr.ValidBadInputData, cerr.BindError))
		return
	}

	if valid, validErrors := cerr.ValidateStruct(uc.Log, editRequest); !valid {
		c.JSON(http.StatusBadRequest, validErrors)
		return
	}

	user := models.AdminUser{ID: editRequest.ID}

	if editRequest.Password != "" {
		pwd, err := bcrypt.GenerateFromPassword([]byte(editRequest.Password), 14)
		if err != nil {
			c.JSON(http.StatusInternalServerError, cerr.LogAndReturnError(uc.Log, err, "Error encrypting password", cerr.GeneralError))
			return
		}
		user.Password = string(pwd)
	}
	user.Forename = editRequest.FirstName
	user.Lastname = editRequest.LastName
	user.Phone = editRequest.Phone
	user.Active = editRequest.Active

	err := db.UpdateUser(user, getUpdatedBy(c))
	if err != nil {
		c.JSON(http.StatusInternalServerError, cerr.LogAndReturnError(uc.Log, err, "Error updating user", cerr.GeneralError))
		return
	}

	updatedUser, err := db.GetUserByID(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, cerr.LogAndReturnError(uc.Log, err, "Could not get current user", cerr.GeneralError))
		return
	}

	c.JSON(http.StatusOK, &UserResponse{
		ID:        updatedUser.ID,
		Email:     updatedUser.Email,
		FirstName: updatedUser.FirstName,
		LastName:  updatedUser.LastName,
		Phone:     updatedUser.Phone,
		Active:    updatedUser.Active,
		Groups:    updatedUser.Groups})
}

//SetGroupsRequest information
type SetGroupsRequest struct {
	ID     int      `form:"id" json:"id" validate:"required"`
	Groups []string `form:"groups" json:"groups" validate:"required"`
}

//SetGroups updates the groups for the specified user
func SetGroups(uc *mw.AdminContext, c *gin.Context) {
	var setGroupRequest SetGroupsRequest
	if err := c.Bind(&setGroupRequest); err != nil {
		c.JSON(http.StatusBadRequest, cerr.LogAndReturnError(uc.Log, err, cerr.ValidBadInputData, cerr.BindError))
		return
	}

	if valid, validErrors := cerr.ValidateStruct(uc.Log, setGroupRequest); !valid {
		c.JSON(http.StatusBadRequest, validErrors)
		return
	}

	activeAdministrators, err := db.AllActiveAdministrators()
	if err != nil {
		c.JSON(http.StatusInternalServerError, cerr.LogAndReturnError(uc.Log, err, "Could not read the users", cerr.GeneralError))
		return
	}

	if len(activeAdministrators) == 1 {
		if activeAdministrators[0].ID == setGroupRequest.ID {
			has := false
			for _, group := range setGroupRequest.Groups {
				if strings.EqualFold(group, "administrators") {
					has = true
					break
				}
			}
			if !has {
				c.JSON(http.StatusBadRequest, cerr.LogAndReturnError(uc.Log, err, "At least an administrator must be active", cerr.GeneralError))
				return
			}
		}
	}

	if err := db.SetGroups(setGroupRequest.ID, setGroupRequest.Groups); err != nil {
		c.JSON(http.StatusBadRequest, cerr.LogAndReturnError(uc.Log, err, "Could not set groups", cerr.BindError))
		return
	}

	updatedUser, err := db.GetUserByID(setGroupRequest.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, cerr.LogAndReturnError(uc.Log, err, "Could not get current user", cerr.GeneralError))
		return
	}

	c.JSON(http.StatusOK, &UserResponse{
		ID:        updatedUser.ID,
		Email:     updatedUser.Email,
		FirstName: updatedUser.FirstName,
		LastName:  updatedUser.LastName,
		Phone:     updatedUser.Phone,
		Active:    updatedUser.Active,
		Groups:    updatedUser.Groups})
}
