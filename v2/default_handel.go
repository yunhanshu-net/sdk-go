package v2

func env(ctx *Context) {
	ctx.Response.OKWithJSON(map[string]string{"version": "1.0", "lang": "go"})
}

func ping(ctx *Context) {
	ctx.Response.OKWithJSON(map[string]string{"ping": "pong"})
}
