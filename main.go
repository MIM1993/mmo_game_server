package main

import (
	"review/zinx/net"
	"review/zinx/ziface"
	"fmt"
	"review/mmo_game_server/core"
	"review/mmo_game_server/apis"
)

//定义在客户端断开后的钩子函数
func OnConnectionLost(conn ziface.IConnection) {
	//客户端已经下线
	fmt.Println("client connid = ", conn.GetConnID(), "logout")

	//得到当前下线的玩家pid
	pid, _ := conn.GetProperty("pid")

	//获取对应的玩家对象
	player := core.WorldMgrObj.GetPlayerByPid(pid.(int32))

	//下线业务
	player.OffLine()

}

//定义在客户端链接后的钩子函数
func OnConnectionAdd(conn ziface.IConnection) {
	fmt.Println("conn Add ")

	//创建一个玩家
	p := core.NewPlayer(conn)

	//给客户端发送一个msgID:1
	p.ReturnPid()

	//给客户端发送一个msgID:200
	p.ReturnPlayerPosition()

	//上线成功
	//将玩家对象添加到世界管理模块中
	core.WorldMgrObj.AddPlayer(p)

	//给conn添加一个属性
	conn.SetProperty("pid", p.Pid)

	//同步周围玩家,告知他们当前玩家已经上线,广播当前玩家位置
	p.SyncSurrounding()

	//打印一句话
	fmt.Println("----->player ID = ", p.Pid, "Onilne ... ", "player num = ", len(core.WorldMgrObj.Players))
}

func main() {
	s := net.NewServer("MMO Game Server")
	//注册一些 链接创建/销毁的 Hook钩子函数
	s.AddOnConnStart(OnConnectionAdd)
	s.AddOnConnStop(OnConnectionLost)

	//添加业务 建立路由
	s.AddRouter(2, &apis.WorldChat{})
	s.AddRouter(3, &apis.Move{})

	//注册一些路由业务
	s.Serve()
}
