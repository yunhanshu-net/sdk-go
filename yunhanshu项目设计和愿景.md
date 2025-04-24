# 项目概述文档

1. 项目基本信息

- 项目名称：云函数

- 项目背景：

我在腾讯企业it部门工作了一年多，主要是做企业服务类项目研发，用户也是服务于企业用户，期间发现企业内部很多工具类的项目，举例：token卡过期，然后用户需要续期token卡，这时候为了满足续期token卡的需求，需要单独去提需求，然后开发token卡续期接口，部署上线，部署的话负责该业务的团队会部署到自己的服务器集群，这时候有个问题是这个服务本身访问量不大，腾讯11万员工的体量，这种接口每天访问量大部分情况下只有几百，甚至更低，这样的话就导致一个问题，这样一个访问量如此之小的服务却需要单独部署一台机器，或者单独部署一个容器，这是对资源极大的浪费，相类似情况还有很多，邮件组续期，业务邮箱续期，个人资产查询，邮件组添加成员，mnet查询，等等，有类似数不胜数的情况，这种接口，都是分布在不同业务部门，然后部署在不同的服务器实例上（也可能是容器），需要24小时不间断的运行，这种程序每天几十到几百的访问量对服务器是一种极大的浪费，甚至有个证书系统对外提供了几十个api的项目，每天访问量也才几千，这种情况太常见了。





-项目的设计思路

所以我就想到一个思路，这个系统是serviceless的faas系统，因为我在腾讯是用go开发项目的且go比较适合这一部分，所以打算用go去实现，想要设计这样的一个项目，大体由三部分组成，api-server，runcher，sdk-go 这三个项目

整个执行流程是这样，用户：tencent在前端，先创建一个项目假设叫tencentOa(腾讯oa系统)，这时runcher会去root目录下创建tencent并在下面创建一个tencentOa的go项目，然后并初始化该目录，这是项目的目录结构

```Go
 //tencentOa下的目录结构
    //bin
	//	data
    //     -data.db 这里会存放程序的数据库文件
	//	.request
	//	tencent_tencentOa_v1 （可执行程序v1版本，每次新增或者删除api都会重新编译一个新的版本）
	//  tencent_tencentOa_v2 （可执行程序v2版本）
	//version （每次变更代码前都会把代码先保存一份在新版本进行改动，假如更新失败还可以回滚到之前的版本）
	//	-v1 
    //      -go.mod
    //      -main.go
	//		-api
	//  -v2
    //      -go.mod
    //      -main.go
	//		-api
    //         -bookManage
    //             bookCreate.go 这个是注册图书的源代码
    //             ....这里还可以有更多的func
    //
```

然后这个项目下是可以创建树形结构的服务目录，这个树形结构有两种，一种是package，一种是func，举例：我在tencent这个项目下创建了一个package叫calcManage(计算管理系统），然后这是一个项目（可以在下面继续创建任意数量的package和func），假如我想要在这个package下创建一个calc的package，我可以在这个目录下创建 add(添加计算记录)函数，这时我需要做的就是描述需求，然后我们的智能体会根据我们自己的sdk-go去生成代码（一个go文件），我们可以审阅一下代码（不懂技术的不审核也没关系，可以直接发布后验证），然后发布该func，这时我们的api-server收到发布请求，会解析出代码，把代码和该代码对应的package位置创建文件，以计算记录管理为例，会生成以下文件 tencent/calcManage/version/下一个版本/api/calc/add.go 

其实创建package的同时会也会创建一个init_.go ，然后runcher会自动把init_.go在main.go引入来初始化这个package

```Go
package calc

func Init()  {}
```

下面是add.go的代码

```Go
package calc

import (
	"github.com/sirupsen/logrus"
	"github.com/yunhanshu-net/sdk-go/model/request"
	"github.com/yunhanshu-net/sdk-go/model/response"
	"github.com/yunhanshu-net/sdk-go/runner"
)

var dbName = "calc.db"

type Calc struct {
	ID       int    `gorm:"primaryKey;autoIncrement" runner:"code:id;name:id"`
	A        int    `json:"a" runner:"code:a;name:a"`
	B        int    `json:"b" runner:"code:b;name:b"`
	C        int    `json:"c" runner:"code:c;name:c"`
	Receiver string `json:"receiver" runner:"code:receiver;name:receiver"`
	Code     string `json:"code" runner:"code:code;name:code"`
}

func init() {
	addConfig := &runner.ApiConfig{
		Tags:        "数据管理;数据分析;记录管理",
		EnglishName: "calcAdd",
		ChineseName: "添加计算记录",
		ApiDesc:     "这里可以描述的详细一点",
		UseTables:   []interface{}{&Calc{}}, //这里会在注册这个api的时候自动创建相关的表
		OnPageLoad: func(ctx *runner.Context) (resetRequest interface{}, resp interface{}, err error) {
			return &AddReq{Receiver: ctx.GetUsername()}, nil, nil
		},
		OnInputValidate: func(ctx *runner.Context, req *request.OnInputValidate) (*response.OnInputValidate, error) {
			msg := ""
			if req.Key == "code" {
				if len(req.Value) > 64 {
					msg = "最长不能超过64个字符"
				}
				//其他判断......
			}
			return &response.OnInputValidate{Msg: msg}, nil
		},
	}

	runner.Get("/calc/add", Add, addConfig)
}

type AddReq struct {
	Receiver string `json:"receiver"`
	A        int    `json:"a" form:"a"`
	B        int    `json:"b" form:"b"`
	Code     string `json:"code" form:"code"`
}

type AddResp struct {
	ID     int `json:"id"`
	Result int `json:"result"`
}

// Add 拿这个处理函数举例，ctx是固定参数， req *AddReq是用户自定义的参数，根据接口请求参数自己定义，resp response.Response是固定参数，用户可以根据这个返回自己的json数据
func Add(ctx *runner.Context, req *AddReq, resp response.Response) error {
	db := ctx.MustGetOrInitDB(dbName)
	res := Calc{A: req.A, B: req.B, C: req.A + req.B} //这里模拟处理逻辑
	err := db.Model(&Calc{}).Create(&res).Error
	if err != nil {
		logrus.Errorf("Add err:%s", err.Error())
		return err
	}
	return resp.JSON(res).Build()
}

```



如此一个添加计算记录的功能就完成了，这个代码完全可以让自己微调好的大模型帮我生成，连写都不用写，



另外这个是我假设新增的另一个get(获取计算列表)的func接口

```Go
package calc

import (
	"github.com/yunhanshu-net/sdk-go/model/request"
	"github.com/yunhanshu-net/sdk-go/model/response"
	"github.com/yunhanshu-net/sdk-go/runner"
	"strconv"
)

type GetReq struct {
	ID int `json:"id" form:"id"`
	*request.PageInfo
}

func init() {
	getConfig := &runner.ApiConfig{
		ChineseName: "获取计算记录",
		EnglishName: "calcGet",
		ApiDesc:     "这里可以描述的详细一点",
		Tags:        "数据管理;数据分析;记录管理",
		OnApiCreated: func(ctx *runner.Context, req *request.OnApiCreated) error {
			db := runner.MustGetOrInitDB(dbName) //这里会返回*gorm.DB
			return db.AutoMigrate(&Calc{})
		},
		AfterApiDeleted: func(ctx *runner.Context, req *request.AfterApiDeleted) error {
			//视情况决定是否要删除表
			return runner.MustGetOrInitDB(dbName).Migrator().DropTable(&Calc{})
		},

		OnPageLoad: func(ctx *runner.Context) (resetRequest interface{}, resp interface{}, err error) {
			//这里返回用户的个人信息
			return &AddReq{A: 1, B: 2, Receiver: ctx.GetUsername()}, nil, nil
		},

		OnInputFuzzy: func(ctx *runner.Context, req *request.OnInputFuzzy) (*response.OnInputFuzzy, error) {
			var values []string
			if req.Key == "a" { //
				db := ctx.MustGetOrInitDB(dbName)
				var calcs []Calc
				db.Model(&Calc{}).Where("a  like ?", "%"+req.Value+"%").Limit(10).Find(&calcs)
				for _, calc := range calcs {
					values = append(values, strconv.Itoa(calc.A))
				}

			}
			return &response.OnInputFuzzy{Values: values}, nil
		},
	}
	runner.Get("/calc/get", Get, getConfig)
}

func Get(ctx *runner.Context, req *GetReq, resp response.Response) error {
	db := ctx.MustGetOrInitDB(dbName)
	var res []Calc
	db.Where("id > ?", req.ID)

	//这里会返回table类型的数据，前端可以直接渲染成element的表格进行展示
	//AutoPaginated 会自动把查询到的数据挂载到res上，自动添加分页等等
	return resp.Table(&res).AutoPaginated(db, &Calc{}, req.PageInfo).Build()
}

```

这个func生成代码后，由于之前已经创建了add函数已经init过这个package了，runcher会创建个新版本，copy代码，编译成新版本，然后把新版本变更的api发给apiserver，更新api，并且更改版本指向，这时候就可以实现软件升级，同时在腾讯oa系统可以看到新加的功能，然后这个函数就是一个可以在前端点击运行的页面程序了，后端的go属性结构目录结构和前端的服务目录（属性结构）是保持一致的，后端是英文名，前端是自己设置的中文名称，这样前端可以进入tencent这个项目下可以看到树形结构信息，并可以看到哪些是可以调用的，可以点到add(添加计算记录),这个页面会有几个对应请求参数的输入框，同时还有一个运行按钮，当输入完请求参数后，点击运行即可把响应参数渲染到界面，整个sdk可以控制渲染的返回是常规的（标题1：值1，标题2：值2），或者/文件/echarts/table/.....其他扩展类型，整个流程基本上是这样，不知道你能不能理解？

默认的api包括

调用指定接口的回调函数

建立长连接（在高并发情况下可以和runcher建立长连接，这个可以保证能处理海量请求）

获取整个项目的所有API信息（api的路径，api的请求参数，api的响应参数，api的method，api的描述，api相关的回调函数）







sdk-go的职责：基于这个sdk写业务代码，

另外我设计了自己的runner 的tag，这些都是在request和response的结构体上用的，用来描述字段的渲染信息，以下是tag的含义（仅供参考，还在优化中）

```Go
通用：
code：
  解释：表单的英文key，默认用json的名称
  示例：base_arr
name：
  解释：表单的名称。
  示例: 原数组
desc: 
  解释：参数的描述信息，这里可以详细描述。
  示例：原始的数据，对比的样本数据
required: 
  解释：是否必填参数。
  示例：必填
example: 
  解释：示例值。
  示例：1,2,3,4
default：
  解释：默认值
  示例：zhangsan,lisi
type：
  解释：数据类型，string/number/bool/float/file。
  示例：string
widget:
  解释：会把该字段渲染成什么组件，input（文本输入框,请求参数默认是输入框，响应参数默认是文本展示框）/checkbox(多选框)/radio(单选框)/select(下拉框)/switch(开关)/slider(滑块)/file(文件上传组件)
  示例：input


文本框包含以下
  placeholder：
    解释：占位符。
    示例：请输入原数组，默认用,号分割
  fuzzy
    解释：是否支持模糊查询
  max：
    解释：最长字符数量
    示例：100
  min
    解释：最短字符数量
    示例：1
  

多选框包含以下





```



基于sdk-go编译的可执行程序runcher可以像执行cmd那样执行这个sdk-go编写的runner程序，会把路由对应成cmd的不同参数，就像这样，

```Go
tencent_tencentOa_v1 /calc/get ./request.json
```

相当于调用了tencent_tencentOa_v1 的http的/calc/get路由，method写在request.json文件中，当用户从页面发起一个请求时候，我可以通过api-server的业务系统从数据库获取接口的信息，然后通过发送nats消息，runcher引擎会监听nats消息，收到消息后再去调用指定的runner，因为项目是多租户项目，每个用户上传的runner都会放在自己的目录下，调用指定runner的原理是，用golang直接调用可执行程序，通过：```Go指定用户目录/指定服务/可执行程序名称_版本(runner) router 路径/request_uuid.json```来调用执行runner，我先来说一下这个原理，可执行程序名字就不用说了吧,router就是和http的路由一样的作用，用来区分不同的接口功能，类似：/api/email/send或者/api/email/list 等等，request.json是调用请求接口时候携带请求参数，runcher引擎会把前端调用接口携带的请求参数写到一个文件，把文件路径/request_uuid.json 当请求参数传递，这时候就可以唤醒启动可执行程序并处理该接口对应的函数了，要用sdk写这个程序的原因是因为不用sdk无法根据我的需求去解析对应的指令和请求参数，sdk还制定了统一的返回数据格式，所以这是要用sdk写程序的原因，



输出的话runner项目会通过fmt.Println来打印输出数据<Response>{这里是json数据，这里json数据被装箱了，用户json返回在内层data，外层还包括了运行时的元数据，比如函数运行耗时哦，运行内存占用，等等}</Response><log>这里是日志数据</log>返回的数据会被runcher引擎解析，然后通过nats的RespondMsg响应给api-server可以直接返回给前端，这样的好处有一些，就是可以把runcher和runner部署在本地服务器上，也就是所谓的混合云，api-server在公有云只负责接口查询和消息传递，不存在什么计算逻辑，复杂的高性能计算都放在了本地，你感觉这个架构如何？后续可以扩展更多标签来处理，

api-server：业务系统（golang）管理所有api（runner）的地址包括api的请求参数，请求方式，响应参数，接口版本等等
runcher：调度执行引擎（golang）可以调度执行各种runner

sdk-go：实现调度相关的协议通过sdk写的程序统称为 runner





其实我现在正在做的是基础设施层面，只要sdk完善了，我的目标是对接ai大模型，自己微调，让大模型基于我的sdk自动生成代码，然后自动发布程序，最终的目标是，用户只需要在界面上描述需求，大模型基于我的sdk自动生成后端代码，后端代码可以携带出前端的参数和渲染信息，能直接渲染成用户界面，这样的话是可以实现无代码平台的，用户不需要关注sdk这个概念，这个sdk只是针对那些程序员的，他们可能需要开发一些非常高级或者业务逻辑非常复杂的应用，他们不需要ai或者写的代码比ai更高级，然后搞成压缩包，然后上传到我们平台，我们自动帮他注册功能页面，对于不懂技术的也没关系，只要会描述问题即可，我的目的是降低应用的开发成本，就像作为一个科研人员，他不懂代码，但是他清楚自己的需求，他只要说出：我需要一个可以生成指定位数斐波那契数列的函数，然后大模型自动生成go代码，自动部署，然后可以根据函数的api路径自动在页面添加功能，然后基于请求参数和响应参数自动渲染用户界面，然后你可以参考calc和 crm calc是我自己写的，crm是ai生成的，假如生成一个可以在页面输入指定位数输出数列的斐波那契数列的功能是非常简单的，这样科研人员根本不需要懂一点代码，这样可以把这个平台的用户扩展到各行各业，你是个普通的办公人员，或者是个企业用户，或者是个零售店的老板，都无所谓，只要能描述出你的需求，我都可以生成代码和界面，最重要的一点是这个项目是以cmd可执行程序为单位的项目，我在web页面先创建一个项目起个名字，然后后端对应一个go的项目，然后一个项目有不同的服务目录，我可以添加服务目录，一个服务目录对应后端go的一个package，比如json package那么这个json的package在后端的表现是一个json的package，在前端的表现是一个json树结构的节点，json下可以有json2csv （json转换csv）在后端的表现是json2csv.go的文件，在前端的表现是json节点下的一个子节点json转换csv的功能，这个会有fx标识，意味着可以直接运行，然后前端的目录结构和后端的代码是完全一一对应的，这样还有一个更伟大的构想是，用户做的功能可以直接让别人基于package的形式去fork到自己的项目里，比如我做了一个图片转换的package，别人觉得好用，同时希望数据能够存储在他的项目里（多租户），只这样后端只需要把原始文件的package fork 到fork用户的指定目录，然后编译，这时候期间用到的数据均（sqlite）会在fork用户的目录记录，实现了数据隔离，代码复用，同时如果fork用户想要迭代的话也不会影响老用户，每次编译都会存储旧版本的可执行程序，这样非常方便回滚，user_runner_version 这个是可执行程序的命名规范，所以后续这个平台我的愿景是这样的

其实我这个项目想做成那种类似工具平台，当创建应用的门槛变低了，那么有想法的人就能创建出天马行空的应用，大家可以在这个平台进行办公，每天上班先打开，我的yunhanshu平台，如果我是在某某公司上班的打工人我会选择打开xxx的某某公司工作常用工具，然后里面都是提升我办公效率的工具，例如：token管理，token生成，批量发送邮件消息，这些可能在没有工具函数的情况下需要自己写代码去实现，但是现在可能只需要描述一下需求即可生成，然后如果你是科研人员，或者大学生，你会打开清华大学科研常用函数/物理学院/求斐波那契数列前n项和，亦或者你是事业单位？北京市图书馆：图书借阅管理，图书查询？等等，我觉得好像可以覆盖很多场景。



