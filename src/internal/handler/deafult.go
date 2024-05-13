package handler

import "github.com/labstack/echo/v4"

type IHandler interface {
	RegisterRoutes(basePath string, e *echo.Group)
}
