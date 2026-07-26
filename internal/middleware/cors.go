package middleware

import (
	"time"

	"github.com/gin-contrib/cors"
)

func CORS(
	origins []string,
) cors.Config {

	return cors.Config{

		AllowOrigins: origins,

		AllowMethods: []string{
			"GET",
			"POST",
			"PUT",
			"PATCH",
			"DELETE",
			"OPTIONS",
		},

		AllowHeaders: []string{
			"Origin",
			"Content-Type",
			"Authorization",
		},

		AllowCredentials: true,

		MaxAge: 12 * time.Hour,
	}
}
