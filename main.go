package main

import (
	"fmt"
	"io"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	//日志记录到当前目录下development.log文件
	f, _ := os.Create("development.log")
	// 同时保留了控制台输出
	gin.DefaultWriter = io.MultiWriter(f, os.Stdout)

	// 创建一个默认的路由器
	r := gin.Default()

	r.Use(userMiddleware()) // 注册用户中间件 这里注册后对所有的路由都能获取到当前用户信息
	// 注册一个hello路由
	r.GET("/hello", func(c *gin.Context) {
		// 向客户端返回hello world
		c.String(200, "hello world")
	})

	// 响应json的hello路由
	r.GET("/hellojson", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"code":    200,
			"message": "hello world",
		})
	})

	// 响应html页面
	r.LoadHTMLGlob("templates/*") // 加载模板文件
	//r.LoadHTMLFiles("templates/template1.html", "templates/template2.html")
	r.GET("/index.html", func(c *gin.Context) {
		c.HTML(200, "index.html", gin.H{
			"title": "dmy", // 传递给模板的数据
		})
	})

	// POST请求
	r.POST("/login", func(c *gin.Context) {
		name := c.PostForm("name")         // 获取表单数据
		password := c.PostForm("password") // 获取表单数据

		type any map[string]interface{}
		c.JSON(200, gin.H{
			"code":    200,
			"message": "login success",
			"data":    any{"name": name, "password": password}, // 传递给客户端的数据
		})
	})

	// PUT更新请求
	r.PUT("/user/:id", func(c *gin.Context) {
		id := c.Param("id") // 获取路径参数
		name := c.PostForm("name")
		password := c.PostForm("password")

		type any map[string]interface{}
		c.JSON(200, gin.H{
			"code":    200,
			"message": "update success",
			"data":    any{"id": id, "name": name, "password": password},
		})
	})

	// DELETE删除请求
	r.DELETE("/user/:id", func(c *gin.Context) {
		id := c.Param("id") // 获取路径参数
		c.JSON(200, gin.H{
			"code":    200,
			"message": "delete success",
			"data":    id,
		})
	})

	//========================== 参数部分 ==========================//
	// 查询参数
	// curl "http://localhost:8080/query?name=dmy&age=20&ids=1&ids=2&ids=3"
	// 返回结果：{"code":200,"data":{"age":"20","ids":["1","2","3"],"name":"dmy"},"message":"query success"}
	r.GET("/query", func(c *gin.Context) {
		name := c.Query("name")            // 获取查询参数
		age := c.DefaultQuery("age", "18") // 获取查询参数，如果不存在则返回默认值
		ids := c.QueryArray("ids")         // 获取查询参数数组

		fmt.Printf("%#v\n", ids)
		c.JSON(200, gin.H{
			"code":    200,
			"message": "query success",
			"data":    gin.H{"name": name, "age": age, "ids": ids},
		})
	})

	type User struct {
		Name string `form:"name" json:"name" xml:"name" binding:"required"`
		Age  int    `form:"age" json:"age" xml:"age" binding:"required"`
	}

	// curl -X POST -H "Content-Type: application/json" -d '{"name": "dmy", "age": 20}' http://localhost:8080/users
	r.POST("/users", func(c *gin.Context) {
		var user User
		// 自动决定绑定类型，默认是JSON绑定 也可以是xml
		if err := c.ShouldBind(&user); err == nil {
			fmt.Println(user.Name, user.Age)
			c.JSON(200, gin.H{
				"code":    200,
				"message": "user created",
				"data":    user,
			})
		} else {
			c.JSON(400, gin.H{
				"code":    400,
				"message": "invalid request",
				"error":   err.Error(),
			})
		}
	})

	// 表单参数
	//  curl -X POST  -d "ids=1&ids=2&ids=3" -d "name=dmy&password=456" -d "account[id]=111&account[name]=dmy" http://localhost:8080/form
	// 返回结果： {"code":200,"data":{"account":{"id":"111","name":"dmy"},"ids":["1","2","3"],"name":"dmy","password":"456"},"message":"form success"}
	r.POST("/form", func(c *gin.Context) {
		name := c.PostForm("name")
		password := c.DefaultPostForm("password", "123456") // 获取表单参数，如果不存在则返回默认值
		ids := c.PostFormArray("ids")                       // 获取表单参数数组
		// 表单map
		account := c.PostFormMap("account") // account[id]=111&account[name]=dmy

		fmt.Println(name, password)
		c.JSON(200, gin.H{
			"code":    200,
			"message": "form success",
			"data":    gin.H{"name": name, "password": password, "ids": ids, "account": account},
		})
	})

	// 文件上传
	// curl -X POST -F "file=@/Users/dongmingyan/Desktop/test.txt" http://localhost:8080/upload
	r.POST("/upload", func(c *gin.Context) {
		file, _ := c.FormFile("file") // 获取上传的文件

		c.SaveUploadedFile(file, "./uploads/"+file.Filename) // 保存文件到指定路径
		c.JSON(200, gin.H{
			"code":    200,
			"message": "upload success",
			"data":    file.Filename,
		})
	})

	//========================== 路由组 ==========================//
	apiGroup := r.Group("/api")
	{
		// api/v1路由组
		v1Group := apiGroup.Group("/v1")
		{
			v1Group.GET("/users", func(c *gin.Context) {
				c.JSON(200, gin.H{"code": 200, "message": "v1 users"})
			})
		}

		// api/v2路由组
		v2Group := apiGroup.Group("/v2")
		{
			v2Group.GET("/users", func(c *gin.Context) {
				c.JSON(200, gin.H{"code": 200, "message": "v2 users"})
			})
		}
	}

	// 测试当前用户信息
	r.GET("/current_user", func(c *gin.Context) {
		// 获取当前用户信息
		currentUser, exists := c.Get("currentUser")

		if exists {
			c.JSON(200, gin.H{
				"code":         200,
				"message":      "current user",
				"current_user": currentUser,
			})
		} else {
			c.JSON(401, gin.H{
				"code":    401,
				"message": "unauthorized",
			})
		}
	})

	r.Run() // 启动服务
}

// 定义用户中间件
func userMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 做一些用户的获取和验证操作

		currentUser := "dmy"              // 获取当前用户信息
		c.Set("currentUser", currentUser) // 设置当前用户信息到上下文

		c.Next() // 继续处理请求
		//c.JSON(404, gin.H{"code": 404, "message": "not found"})
		//c.Abort() // 中止请求 终止前要写响应不然啥也没有
	}
}
