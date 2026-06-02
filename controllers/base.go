package controllers

import beego "github.com/beego/beego/v2/server/web"

// BaseController provides shared JSON response helpers.
type BaseController struct {
	beego.Controller
}

// SendSuccess writes a 2xx JSON response with an optional data payload.
func (c *BaseController) SendSuccess(status int, message string, data interface{}) {
	c.Ctx.Output.SetStatus(status)
	resp := map[string]interface{}{
		"success": true,
		"message": message,
	}
	if data != nil {
		resp["data"] = data
	}
	c.Data["json"] = resp
	c.ServeJSON()
}

// SendError writes a 4xx/5xx JSON error response.
func (c *BaseController) SendError(status int, message string) {
	c.Ctx.Output.SetStatus(status)
	c.Data["json"] = map[string]interface{}{
		"success": false,
		"message": message,
	}
	c.ServeJSON()
}