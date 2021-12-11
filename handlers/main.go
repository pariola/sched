package handlers

import (
	"github.com/labstack/echo/v4"
)

// handler
type handler struct {
}

// New
func New() *handler {
	return &handler{}
}

// CreateJob
func (h handler) CreateJob(c echo.Context) error {
	return nil
}

// UpdateJob
func (h handler) UpdateJob(c echo.Context) error {
	return nil
}

// CancelJob
func (h handler) CancelJob(c echo.Context) error {
	return nil
}
