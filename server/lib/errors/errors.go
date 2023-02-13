package errors

import (
	models "degrens/panel/models"
)

var BodyParsingFailed models.RouteErrorMessage = models.RouteErrorMessage{
	Title:       "Server error",
	Description: "An error occurred while reading your incoming request",
}

var Unauthorized models.RouteErrorMessage = models.RouteErrorMessage{
	Title:       "Authorization error",
	Description: "Failed to get valid authorization information from the request",
}
