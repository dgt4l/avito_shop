package handler

func RegisterRoutes(h *ShopHandler) {
	authRouter := h.e.Group("/api")
	authRouter.POST("/auth", h.AuthUser)
	authRouter.GET("/ping", h.Ping)

	router := h.e.Group("/api", h.AuthMiddleware())
	router.GET("/info", h.GetInfo)
	router.GET("/buy", h.BuyItem)
	router.POST("/sendCoin", h.SendCoin)
}
