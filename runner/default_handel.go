package runner

func env(ctx *HttpContext) {
	ctx.Response.JSON(map[string]string{"version": "1.0", "lang": "go"}).Build()
}

func ping(ctx *HttpContext) {
	ctx.Response.JSON(map[string]string{"ping": "pong"}).Build()
}
