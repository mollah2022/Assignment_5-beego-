// @APIVersion 1.0.0
// @Title Expense Tracker API
// @Description Personal Expense Tracker REST API built with Go and Beego v2
// @Contact admin@example.com
// @License Apache 2.0
package routers

import (
	"expense-tracker-api/controllers"

	beego "github.com/beego/beego/v2/server/web"
	"github.com/beego/beego/v2/server/web/context"
)

func init() {
	beego.Get("/api/v1/health", func(ctx *context.Context) {
		ctx.Output.SetStatus(200)
		ctx.Output.JSON(map[string]interface{}{
			"success": true,
			"message": "Server is running",
		}, false, false)
	})

	ns := beego.NewNamespace("/api/v1",
		beego.NSNamespace("/auth",
			beego.NSInclude(&controllers.AuthController{}),
			beego.NSRouter("/register", &controllers.AuthController{}, "post:Register"),
			beego.NSRouter("/login", &controllers.AuthController{}, "post:Login"),
		),
		beego.NSNamespace("/expenses",
			beego.NSInclude(&controllers.ExpenseController{}),
			beego.NSRouter("/summary", &controllers.ExpenseController{}, "get:Summary"),
			beego.NSRouter("/", &controllers.ExpenseController{}, "post:Create;get:List"),
			beego.NSRouter("/:id", &controllers.ExpenseController{}, "get:GetOne;put:Update;delete:Delete"),
		),
	)
	beego.AddNamespace(ns)
}
