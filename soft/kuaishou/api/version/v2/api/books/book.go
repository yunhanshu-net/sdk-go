package books

import (
	"github.com/yunhanshu-net/sdk-go/model/request"
	"github.com/yunhanshu-net/sdk-go/runner"
)

type Book struct {
	ID        int     `json:"id" gorm:"primary_key"`
	Name      string  `json:"name"`
	Desc      string  `json:"desc"`
	Price     float32 `json:"price"`
	Author    string  `json:"author"`
	CreatedBy string  `json:"created_by"`
}

func (Book) TableName() string {
	return "book"
}

type GetBookReq struct {
	ID int `json:"id" runner:"desc:图书编号;required:必填;example:1,2,3,4"`
	request.PageInfo
}

type GetBookResp struct {
	ID        int     `json:"id" runner:"desc:图书编号;example:1,2,3,4"`
	Name      string  `json:"name" runner:"desc:图书名称;example:三体"`
	Desc      string  `json:"desc" runner:"desc:介绍;example:刘慈欣的科幻小说"`
	Price     float32 `json:"price" runner:"desc:图书价格;example:199.9"`
	Author    string  `json:"author" runner:"desc:图书作者;example:刘慈欣"`
	CreatedBy string  `json:"created_by" runner:"desc:创建该图书的用户;example:beiluo"`
}

type CreateBookReq struct {
	ID     int     `json:"id" runner:"desc:图书编号;example:1,2,3,4"`
	Name   string  `json:"name" runner:"desc:图书名称;example:三体"`
	Desc   string  `json:"desc" runner:"desc:介绍;example:刘慈欣的科幻小说"`
	Price  float32 `json:"price" runner:"desc:图书价格;example:199.9"`
	Author string  `json:"author" runner:"desc:图书作者;example:刘慈欣"`
}

type CreateBookResp struct {
	ID int `json:"id" runner:"desc:图书编号;example:1,2,3,4"`
}

func init() {

	getBookConfig := &runner.ApiConfig{
		ApiDesc:     "根据编号获取书信息",
		ChineseName: "获取书籍信息",
		Labels:      []string{"图书管理", "图书馆", "学生", "科研"}, //适用场景？
		EnglishName: "getBook",
		Classify:    "图书管理",
		//Tags:        "array,list,数组,集合", //表示元素本身的特性
		Request:  &GetBookReq{}, //这里可以主要是为了
		Response: &GetBookResp{},
		OnPageLoad: func(ctx *runner.HttpContext) error {
			var req GetBookReq
			db, err := ctx.GetAndInitDB("books.db") //这里有的话会连接数据库，没有的话会自动创建sqlite数据库
			if err != nil {
				return err
			}
			db.Where("created_by = ?", ctx.GetUser())
			var books []Book
			return ctx.Response.Table(&books).AutoPaginated(db, &Book{}, &req.PageInfo).Build()

		},
	}

	createBookConfig := &runner.ApiConfig{
		ApiDesc:     "新增一个图书",
		ChineseName: "新增图书",
		Labels:      []string{"图书管理", "图书馆", "学生", "科研"}, //适用场景？
		EnglishName: "createBook",
		Classify:    "图书管理",
		//Tags:        "array,list,数组,集合", //表示元素本身的特性
		Request:  &CreateBookReq{},  //这里可以主要是为了反射时候获取请求字段信息，可以直接一键生成请求参数
		Response: &CreateBookResp{}, //这里可以主要是为了可以通过反射获取相应字段信息，可以直接一键生成响应参数

		OnCreated: func(ctx *runner.HttpContext) error {
			//创建新的api时候的回调函数,新增一个api假如新增了一张user表，
			//可以在这里用gorm的db.AutoMigrate(&User)来创建表，保证新版本的api可以正常使用新增的表
			//这个api只会在我创建这个api的时候执行一次
			db, err := ctx.GetAndInitDB("books.db") //这里有的话会连接数据库，没有的话会自动创建sqlite数据库
			if err != nil {
				return err
			}
			return db.AutoMigrate(&Book{})
		},

		OnPageLoad: func(ctx *runner.HttpContext) error {
			//假如该接口有对应的前端界面，渲染该界面后会调用该函数来加载默认请求数据，
			//比如一个用户订单列表的页面，在点进去页面后会调用该回调
			//此时已经知道是哪个用户的了，然后可以根据用户信息，展示该用户的默认数据。
			//这样就省的用户自己输入用户名然后再点击运行按钮展示出来了

			return nil
		},

		OnFuzzy: func(ctx *runner.HttpContext) error {
			//模糊搜索回调函数，比如搜索用户，可以在这里做一些操作，比如根据用户名模糊搜索用户，然后返回用户列表
			return nil
		},

		BeforeClose: func(ctx *runner.HttpContext) error {
			//程序结束前的回调函数，可以在程序结束前做一些操作，比如上报一些数据?或者其他场景？我暂时还没想到有更好的场景
			return nil
		},
		AfterClose: func(ctx *runner.HttpContext) error {
			//程序结束后的回调函数，可以在程序结束后做一些操作，比如清理某些文件
			return nil
		},

		AfterDelete: func(ctx *runner.HttpContext) error {
			//api删除后触发回调，比如该api删除的话，可以在这里做一些操作，比如删除该api对应的表
			db, err := ctx.GetAndInitDB("books.db") //这里有的话会连接数据库，没有的话会自动创建sqlite数据库
			if err != nil {
				return err
			}
			return db.Migrator().DropTable(&Book{})
		},
	}

	runner.Get("/book/list", GetBookApi, getBookConfig)
	runner.Post("/book/create", CreateBookApi, createBookConfig)

}

func GetBookApi(ctx *runner.HttpContext) error {
	var req GetBookReq
	err := ctx.Request.ShouldBindJSON(&req)
	if err != nil {
		return err
	}
	db, err := ctx.GetAndInitDB("books.db") //这里有的话会连接数据库，没有的话会自动创建sqlite数据库
	if err != nil {
		return err
	}
	var book Book
	db.Where("id=?", req.ID).First(&book)
	resp := GetBookResp{
		ID:     book.ID,
		Name:   book.Name,
		Desc:   book.Desc,
		Price:  book.Price,
		Author: book.Author,
	}
	return ctx.Response.JSON(resp).Build()

}

func CreateBookApi(ctx *runner.HttpContext) error {
	var req CreateBookReq
	err := ctx.Request.ShouldBindJSON(&req)
	if err != nil {
		return err
	}
	db, err := ctx.GetAndInitDB("books.db") //这里有的话会连接数据库，没有的话会自动创建sqlite数据库
	if err != nil {
		return err
	}
	var book Book
	db.Where("id=?", req.ID).First(&book)
	if book.ID != 0 {
		return ctx.Response.FailWithJSON(nil, "图书已存在")
	}
	book = Book{
		ID:        req.ID,
		Name:      req.Name,
		Desc:      req.Desc,
		Price:     req.Price,
		Author:    req.Author,
		CreatedBy: ctx.GetUser(),
	}
	db.Create(&book)
	resp := CreateBookResp{
		ID: book.ID,
	}
	return ctx.Response.JSON(resp).Build()

}
