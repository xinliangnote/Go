## 项目介绍

[Gin 路由配置](https://github.com/xinliangnote/Go/blob/master/01-Gin框架/02-路由配置.md)

## 配置

```
func InitRouter(r *gin.Engine)  {

	r.GET("/sn", SignDemo)

	// v1 版本
	GroupV1 := r.Group("/v1")
	{
		GroupV1.Any("/product/add", v1.AddProduct)
		GroupV1.Any("/member/add", v1.AddMember)
	}

	// v2 版本
	GroupV2 := r.Group("/v2", common.VerifySign)
	{
		GroupV2.Any("/product/add", v2.AddProduct)
		GroupV2.Any("/member/add", v2.AddMember)
	}
}
```

## 运行

**下载源码后，请先执行 `dep ensure` 下载依赖包！**

## 效果图

![](https://github.com/xinliangnote/Go/blob/master/01-Gin框架/images/02-路由配置/2_go_1.png)

![](https://github.com/xinliangnote/Go/blob/master/01-Gin框架/images/02-路由配置/2_go_2.png)

![](https://github.com/xinliangnote/Go/blob/master/01-Gin框架/images/02-路由配置/2_go_3.png)

![](https://github.com/xinliangnote/Go/blob/master/01-Gin框架/images/02-路由配置/2_go_4.png)
