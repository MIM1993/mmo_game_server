package apis

import (
	"review/zinx/net"
	"review/zinx/ziface"
	"review/mmo_game_server/pb"
	"github.com/golang/protobuf/proto"
	"fmt"
	"review/mmo_game_server/core"
)

//坐标更新  路由业务
type Move struct {
	net.BaseRouter
}

//处理业务
func (m *Move) Handle(request ziface.IRequest) {
	//解析当前客户段发送过来的消息  proto协议  msgId=3
	proto_msg := &pb.Position{} //定义容器
	proto.Unmarshal(request.GetMsg().GetMsgData(), proto_msg)

	//通过链接属性,获取当前玩家的ID
	pid, _ := request.GetConnection().GetProperty("pid")

	fmt.Println("player id =", pid.(int32), "move -->", proto_msg.X, ",", proto_msg.Z, ",", proto_msg.V)

	//通过pid得到当前的玩家得对象
	player := core.WorldMgrObj.GetPlayerByPid(pid.(int32))

	//将玩家当前位置 发送给周围其他玩家

	//player.UpdatePosition(player.X, player.Y, player.Z, player.V)
	player.UpdatePosition(proto_msg.X, proto_msg.Y, proto_msg.Z, proto_msg.V)

}
