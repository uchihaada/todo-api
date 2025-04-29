package todoapi_2_0

import (
	"github.com/gin-gonic/gin"
	"todo-api-2.0/config"
	"todo-api-2.0/routes"
)

func main() {
	config.InitDB()

	r := gin.Default()

	routes.SetUpRoutes(r)

	if err := r.Run(":8080"); err != nil {
		panic(err)
	}
}
