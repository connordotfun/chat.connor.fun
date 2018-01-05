package prehandlers

import "github.com/labstack/echo"


func HerokuHTTPSRedirect(next echo.HandlerFunc) echo.HandlerFunc {

	return func(c echo.Context) error {

		req := c.Request()
		host := req.Host
		uri := req.RequestURI
		if req.Header.Get("x-forwarded-proto") != "https" {
			if host[:3] == "www" {
				return c.Redirect(301, "https://"+host[4:]+uri)
			}
			return c.Redirect(301, "https://"+host+uri)
		}
		return next(c)
	}
}
