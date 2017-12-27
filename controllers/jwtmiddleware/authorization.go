package jwtmiddleware

import (
	"github.com/labstack/echo"
	"net/http"
	"github.com/aaronaaeng/chat.connor.fun/model"
	"github.com/aaronaaeng/chat.connor.fun/db/roles"
	"github.com/aaronaaeng/chat.connor.fun/controllers"
)

func doAuthorization(next echo.HandlerFunc, claims *Claims, c echo.Context) error {
	permissions := model.NewPermissionSet()
	var principleRole *model.Role
	if claims != nil { //there are authenticated claims
		err := claims.Valid()
		if err != nil {
			c.JSON(http.StatusUnauthorized, invalidTokenResponse)
		}
		if claims.User.Username != "" { //this is very hacky
			userRoles, err := roles.Repo.GetUserRoles(claims.User.Id)
			if err != nil {
				return c.JSON(http.StatusInternalServerError, controllers.Response{
					Error: &controllers.ResponseError{Type: "ROLES_ACCESS_FAILED", Message: err.Error()},
					Data: nil,
				})
			}

			for _, role := range userRoles {
				if role.Name == "admin" {
					principleRole = role //TODO: make this system better
				}
				permissions.Add(role.Permissions...)
			}
		} else {
			anon := model.Roles.GetRole("anon_user")
			permissions.Add(anon.Permissions...)
			principleRole = &anon
		}

		if claims.Permissions != nil { //cached or extra permissions
			permissions.Add(claims.Permissions...)
		}
	} else {
		anon := model.Roles.GetRole("anon_user")
		permissions.Add(anon.Permissions...)
		principleRole = &anon
	}

	if principleRole.Name == "admin" {
		return next(c)
	}
	if principleRole.Name == "banned" {
		return c.JSON(http.StatusForbidden, controllers.Response{
			Error: &controllers.ResponseError{Type: "BANNED", Message: "User banned"},
			Data: nil,
		})
	}

	if isAuthorized(permissions, c.Request()) {
		return next(c)
	} else {
		return c.JSON(http.StatusForbidden, controllers.Response{
			Error: &controllers.ResponseError{Type: "UNAUTHORIZED", Message: "Cannot access resource"},
			Data: nil,
		})
	}
}

func isAuthorized(permissionSet *model.PermissionSet, r *http.Request) bool { //TODO: this will be rrreally slow
	method := r.Method
	resourcePath := r.URL.Path

	permissions := permissionSet.Permissions()
	for _, permission := range permissions {
		if permission.IsPermitted(method, resourcePath) {
			return true
		}
	}
	return false
}
