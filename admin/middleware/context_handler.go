package middleware

import (
	"github.com/Soneso/lumenshine-backend/admin/db"
	"github.com/Soneso/lumenshine-backend/admin/route"
	"github.com/Soneso/lumenshine-backend/helpers"
	"net/http"
	"strings"

	cerr "github.com/Soneso/lumenshine-backend/icop_error"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

//AdminContext context used in the apis to store some default values
type AdminContext struct {
	RequestID string
	Language  string
	Log       *logrus.Entry
	User      *db.UserDetails
}

//UseAdminContext is used to wrap the gin context
func UseAdminContext(f func(uc *AdminContext, c *gin.Context), handlerName string) gin.HandlerFunc {
	uc := &AdminContext{}

	return func(c *gin.Context) {
		var err error
		uc.Language = c.GetString("language")
		uc.RequestID = c.GetString("request_id")
		uc.Log = helpers.GetDefaultLog(c.GetString("servicename"), uc.RequestID)

		//we read the user and groups on every request
		//we take the id from the context, because the jwt middleware saved that in there
		email := c.GetString("userID")

		uc.User, err = db.GetUserByEmail(email)
		if err != nil {
			uc.Log.WithError(err).WithField("email", email).Error("Could not get user groups")
			c.Abort() //on error we abort the requesthandling
			return
		}

		if uc.User == nil {
			c.JSON(http.StatusForbidden, cerr.NewIcopErrorShort(cerr.UserNotExists, "User does not exist in db"))
			return
		}

		//check that user is active
		if !uc.User.Active {
			c.JSON(http.StatusForbidden, cerr.NewIcopErrorShort(cerr.UserInactive, "User not active"))
			return
		}

		r := route.GetRouteForName(handlerName)
		if r == nil {
			c.JSON(http.StatusForbidden, cerr.NewIcopErrorShort(cerr.NoPermission, "route not found"))
			return
		}
		if r.RequiredGroups != nil && len(r.RequiredGroups) > 0 {
			//check that user has right by checking the groups of the user
			has := false
			for _, group := range uc.User.Groups {
				for _, gr := range r.RequiredGroups {
					if strings.ToLower(group) == strings.ToLower(gr) {
						has = true
						break
					}
				}
				if has {
					break //no need to run another group
				}
			}
			if !has {
				c.JSON(http.StatusForbidden, cerr.NewIcopErrorShort(cerr.NoPermission, "User has no permission"))
				return
			}
		}

		f(uc, c)
	}
}
