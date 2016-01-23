# nodenet
一个分布式，异步，流式的消息框架
	
## 设计思想

>一个服务，可以向客户提供不同的服务(API)。

>一个图，定义消息流经的组件。

>一个组件，可以实处理多种消息.

>一个组，包函多个同类型的组件。

>一个服务，对应一条具体的消息格式和一个图。
	
## 配置
参考[XIM](https://github.com/liuhengloveyou/xim/blob/master/example/nodenet.conf.sample)的配置

	{
	"components":[ #一个系统里有很多组件
		{
       "name": "access-127.0.0.1.5001", #组件名
       "intype": "tcp", #组件输入口类型
		"inconf": { #组件输入口配置
	      "url": "127.0.0.1:5001",
	      "timeout": 3}
		},
	
		{
			"name": "tgroup-127.0.0.1.6001",
        	"intype": "tcp",
			"inconf": {
				"url": "127.0.0.1:6001",
	      		"timeout": 3}
		},
	
		{
			"name": "state-127.0.0.1.6002",
        	"intype": "tcp",
			"inconf": {
	   	   "url": "127.0.0.1:6002",
	   	   "timeout": 3}
		},

		{
			"name": "forward-127.0.0.1.6003",
        	"intype": "tcp",
			"inconf": {
	      		"url": "127.0.0.1:6003",
	      		"timeout": 3}
		}],
	
	"groups":[ # 组
		{
			"name": "tgroup", # 组名
			"dispense":"hash", # 分发策略
			"members":["tgroup-127.0.0.1.6001", "tgroup-127.0.0.1.6001"]
		},{
			"name": "state",
			"dispense":"hash",
			"members":["state-127.0.0.1.6002"]
		},{
			"name": "forward",
			"dispense":"hash",
			"members":["forward-127.0.0.1.6003"]
		}],

	"graphs": { # 图
		"tgroup":["tgroup"], 
		"state":["state"],
		"forward":["forward"]
		}
	}


## 开始使用

1. 定义业务消息类型,并注册到nodenet系统中
	
	参考[XIM消息定义](https://github.com/liuhengloveyou/xim/blob/master/common/message.go)
	
		import (
			"github.com/liuhengloveyou/nodenet"
		)
		
		... ...
		
		func init() {
			nodenet.RegisterMessageType(MessageLogin{}) //注册消息类型		
		}
		
		... ... 
		
		// 长连接登入
		type MessageLogin struct {
			Userid         string // 用户ID
			ClientType     string // 客户端类型
			AccessName     string // 接入节点名
			AccessSession  string // 接入节点会话ID
			ConfirmMessage int64  // 确认的消息标识
			UpdateTime     int64  // 状态更新时间
		}

	
2. 定义消息处理函数,并注册到nodenet系统中

	参考[XIM状态逻辑](https://github.com/liuhengloveyou/xim/blob/master/logic/state.go)
	
		import (
			"github.com/liuhengloveyou/nodenet"
		)
	
		func init() {
			nodenet.RegisterWorker("UerLogin", common.MessageLogin{}, UerLogin)
		}
		
		... ...
		
		
		func UerLogin(data interface{}) (result interface{}, err error){
			var msg = data.(common.MessageLogin)
			
			... ...
			
3. 配置服务进程要启动的组件

	参考[XIM状辑服务配置](https://github.com/liuhengloveyou/xim/blob/master/example/logic.conf.simple)
	
		{
		... ...

		"nodes":[        #服务进程要启动的组件
			{
			"name":"tgroup-127.0.0.1.6001",        #组件名
			"works":{        # 组件可以处理的消息类型，和对应的处理函数
				"common.MessageTGLogin":"TempGroupLogin",
				"common.MessageForward":"TempGroupSend"}
			},
			{
			"name":"state-127.0.0.1.6002",
			"works":{
				"common.MessageLogin":"UerLogin",
				"common.MessageLogout":"UerLogout"}
			},
			{
			"name":"forward-127.0.0.1.6003",
			"works":{
				"common.MessageForward":"ForwardMessage"}
			}
		],
	
		... ...
	
		}			
3. 依配置信息启动

	参考[XIM逻辑服务](https://github.com/liuhengloveyou/xim/blob/master/logic/logic.go)

		import (
			"github.com/liuhengloveyou/nodenet"
		)
		... ...

		func initNodenet(fn string) error {
			if e := nodenet.BuildFromConfig(fn); e != nil {
				return e
			}

			for i := 0; i < len(common.LogicConf.Nodes); i++ {
				name := common.LogicConf.Nodes[i].Name
				mynodes[name] = nodenet.GetComponentByName(name)
				if mynodes[name] == nil {
					return fmt.Errorf("No node: %v.", name)
				}

				for k, v := range common.LogicConf.Nodes[i].Works {
					t, w := nodenet.GetMessageTypeByName(k),nodenet.GetWorkerByName(v)
					if t == nil {
						return fmt.Errorf("No message registerd: %s", k)
					}
					if w == nil {
						return fmt.Errorf("No worker registerd: %s", v)
					}
					if reflect.TypeOf(t) != w.Message {
						return fmt.Errorf("worker can't recive message: %v %v", w, k)
					}
					mynodes[name].RegisterHandler(t, w.Handler)
				}

				go mynodes[name].Run()
			}
		
			... ...	
		}


是不是很方便省事？:)


呃。。。忘了要发消息：

	... ...
	
	msg := &common.MessageLogin{
		Userid: user.Userid,
		ClientType: user.Client,
		AccessName: common.AccessConf.NodeName,
		AccessSession: sess.Id("")} // 业务系统消息
		
	g := nodenet.GetGraphByName(common.LOGIC_STATE) //图
	
	cMsg := nodenet.NewMessage(common.GID.ID(),NodeName,g,msg) // 消息
	cMsg.DispenseKey = user.Userid //分发策略键

	// 发送
	if e = nodenet.SendMsgToNext(cMsg); e != nil {
		log.Errorln("SendMsgToNext ERR:", e.Error())
		... ...
	}