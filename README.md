# easy-app 
a easy way to create golang web app  

```
var app = NewApp(cfg)

	// web层拦截
	app.Controller.RegisterMidleWare("prefix", func(c *gin.Context) {
		fmt.Println("before 调用", c.Request.RequestURI)
		c.Next()
		fmt.Println("after 调用", c.Request.RequestURI)
	})

	//业务层拦截
	app.AddBeforeLogicFilter(func(name string, args []interface{}) {
		fmt.Printf("------------拦截所有 方法 和参数--------- 方法名:%s , 参数：%v \n", name, args)
	})

	app.AddAfterLogicFilter(func(name string, result []interface{}) {
		fmt.Printf("------------拦截所有 方法 和返回值-------- 方法名:%s ,返回值：%v \n", name, result)
	})

	// dispatcher
	controller.InitDispatcherHandler(app)

	//logical
	service.InitLogicHandler(app)

	app.Start()
```

