package apis

import (
	"review/zinx/net"
	"review/zinx/ziface"
	"review/mmo_game_server/pb"
	"github.com/golang/protobuf/proto"
	"fmt"
	"review/mmo_game_server/core"
)

//世界聊天  路由业务
type WorldChat struct {
	net.BaseRouter
}

//重写方法
func (wc *WorldChat) Handle(request ziface.IRequest) {
	//1 解析客户传递进来的protobuf数据
	proto_msg := &pb.Talk{}

	if err := proto.Unmarshal(request.GetMsg().GetMsgData(), proto_msg); err != nil {
		fmt.Println("Talk message unmarshal error ", err)
		return
	}

	//通过获取链接属性,获取当前玩家的pid
	pid, err := request.GetConnection().GetProperty("pid")
	if err!=nil{
		fmt.Println("get pid err :",err)
		return 
	}
	
	//通过当前pid获取对应的player属性
	player :=  core.WorldMgrObj.GetPlayerByPid(pid.(int32))

	// 当前的聊天数据广播给全部的在线玩家
	//当前玩家的windows客户端发送过来的消息
	player.SendTalkMsgToAll(proto_msg.GetContent())

}
