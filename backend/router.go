package main

import (
	"github.com/kingwrcy/moments/handler"
	"github.com/kingwrcy/moments/vo"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/samber/do/v2"
	echoSwagger "github.com/swaggo/echo-swagger"
)

func setupRouter(injector do.Injector) {
	userHandler := handler.NewUserHandler(injector)
	memoHandler := handler.NewMemoHandler(injector)
	commentHandler := handler.NewCommentHandler(injector)
	sycConfigHandler := handler.NewSysConfigHandler(injector)
	fileHandler := handler.NewFileHandler(injector)
	tagHandler := handler.NewTagHandler(injector)
	rssHandler := handler.NewRssHandler(injector)
	e := do.MustInvoke[*echo.Echo](injector)
	cfg := do.MustInvoke[*vo.AppConfig](injector)

	apiGroup := e.Group("/api")

	userGroup := apiGroup.Group("/user")
	userGroup.POST("/login", userHandler.Login)
	userGroup.POST("/reg", userHandler.Reg)
	userGroup.POST("/profile", userHandler.Profile)
	userGroup.POST("/profile/:username", userHandler.ProfileForUser)
	userGroup.POST("/saveProfile", userHandler.SaveProfile)

	memoGroup := apiGroup.Group("/memo")
	memoGroup.POST("/list", memoHandler.ListMemos)
	memoGroup.POST("/save", memoHandler.SaveMemo)
	memoGroup.POST("/remove", memoHandler.RemoveMemo)
	memoGroup.POST("/like", memoHandler.LikeMemo)
	memoGroup.POST("/get", memoHandler.GetMemo)
	memoGroup.POST("/setPinned", memoHandler.SetPinned)
	memoGroup.POST("/getFaviconAndTitle", memoHandler.GetFaviconAndTitle)
	memoGroup.POST("/getDoubanMovieInfo", memoHandler.GetDoubanMovieInfo)
	memoGroup.POST("/getDoubanBookInfo", memoHandler.GetDoubanBookInfo)

	commentGroup := apiGroup.Group("/comment")
	commentGroup.POST("/add", commentHandler.AddComment)
	commentGroup.POST("/remove", commentHandler.RemoveComment)

	sycConfigGroup := apiGroup.Group("/sysConfig")
	sycConfigGroup.POST("/save", sycConfigHandler.SaveConfig)
	sycConfigGroup.POST("/get", sycConfigHandler.GetConfig)
	sycConfigGroup.POST("/getFull", sycConfigHandler.GetFullConfig)

	tagGroup := apiGroup.Group("/tag")
	tagGroup.POST("/list", tagHandler.List)

	fileGroup := apiGroup.Group("/file")
	fileGroup.POST("/exist", fileHandler.Exist)
	fileGroup.POST("/upload", fileHandler.Upload)
	fileGroup.POST("/clean", fileHandler.Clean)
	fileGroup.POST("/s3PreSigned", fileHandler.S3PreSigned)

	uploadGroup := e.Group("/upload")
	// 安全头:防止上传的文件被浏览器当作 HTML/脚本执行(存储型 XSS)。
	// nosniff 禁止 MIME 嗅探,CSP sandbox 禁用脚本执行,X-Frame-Options 防点击劫持。
	uploadGroup.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			h := c.Response().Header()
			h.Set("X-Content-Type-Options", "nosniff")
			h.Set("Content-Security-Policy", "default-src 'none'; sandbox; img-src 'self'; media-src 'self'")
			h.Set("X-Frame-Options", "SAMEORIGIN")
			return next(c)
		}
	})
	uploadGroup.Use(middleware.StaticWithConfig(middleware.StaticConfig{
		Root:       cfg.UploadDir,
		HTML5:      false,
		IgnoreBase: true,
		Browse:     false,
	}))

	rssGroup := e.Group("/rss")
	rssGroup.GET("", rssHandler.GetRss)

	if cfg.EnableSwagger {
		e.GET("/swagger/*", echoSwagger.WrapHandler)
	}

}
