package jwtmiddleware

import (
	"github.com/labstack/echo"
	"net/http"
	"github.com/aaronaaeng/chat.connor.fun/model"
	"github.com/aaronaaeng/chat.connor.fun/controllers/auth"
	"github.com/aaronaaeng/chat.connor.fun/context"
	"github.com/aaronaaeng/chat.connor.fun/db"
	"github.com/satori/go.uuid"
)

func unpackClaims(claims *auth.Claims) (model.User, []model.Permission, error) {
	if claims != nil {
		err := claims.Valid()
		if err != nil {
			return model.User{}, nil, err
		}
		user := claims.User
		permissions := claims.Permissions
		return user, permissions, nil
	}
	return model.User{}, make([]model.Permission, 0), nil
}

func getRoles(user *model.User, repo db.RolesRepository) ([]model.Role, error) {
	if user.Id == uuid.Nil {
		anon := model.Roles.GetRole(model.RoleAnon)
		return []model.Role{anon}, nil
	} else {
		return repo.GetUserRoles(user.Id)
	}
}

func doAuthorization(next echo.HandlerFunc, claims *auth.Claims, c echo.Context, repo db.RolesRepository) error {
	ac := c.(context.AuthorizedContext)
	permissions := model.NewPermissionSet()

	user, extraPermissions, err := unpackClaims(claims)
	if err != nil {
		return c.JSON(http.StatusBadRequest, model.NewErrorResponse("INVALID_CLAIMS"))
	}
	ac.SetRequestor(user)
	permissions.AddAll(extraPermissions...)

	roles, err := getRoles(&user, repo)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, model.NewErrorResponse("AUTH_FAILED"))
	}

	for _, r := range roles {
		if r.Override == "ALLOW_ALL" {
			ac.SetAccessCode(model.NewFullAccessCode())
			return next(c)
		} else if r.Override == "BANNED" {
			return c.JSON(http.StatusForbidden, model.NewErrorResponse("BANNED"))
		}
		permissions.AddAll(r.Permissions...)
	}

	authorized, code := isAuthorized(permissions, c.Request())
	if authorized {
		ac.SetAccessCode(code)
		return next(ac)
	}
	return c.JSON(http.StatusForbidden, model.NewErrorResponse("NOT_AUTHORIZED"))
}

func isAuthorized(permissionSet *model.PermissionSet, r *http.Request) (bool, model.AccessCode) { //TODO: this will be rrreally slow
	method := r.Method
	resourcePath := r.URL.Path

	permissions := permissionSet.Permissions()
	for _, permission := range permissions {
		if permission.IsPermitted(method, resourcePath) {
			return true, permission.Code
		}
	}
	return false, 0
}
