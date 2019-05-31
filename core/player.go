package core

import (
	"review/zinx/ziface"
	"sync"
	"math/rand"
	"github.com/golang/protobuf/proto"
	"fmt"
	"review/mmo_game_server/pb"
)

//定义玩家结构体
type Player struct {
	//当前玩家的ID
	Pid int32
	//当前玩家的链接
	Conn ziface.IConnection
	//平面的X轴坐标
	X float32
	//玩家高度
	Y float32
	//平面的y轴坐标
	Z float32
	//玩家脸的朝向
	V float32
}

//play生成器

//生成玩家ID生成器
var PidGen int32 = 1
//保护锁
var IdLock sync.Mutex

//初始化玩家方法
func NewPlayer(conn ziface.IConnection) *Player {
	//分配玩家ID
	//加锁
	IdLock.Lock()
	id := PidGen
	PidGen++
	IdLock.Unlock()

	//创建一个玩家对象
	p := &Player{
		Pid:  id,
		Conn: conn,
		X:    float32(160 + rand.Intn(10)),
		Y:    0,
		Z:    float32(140 + rand.Intn(10)),
		V:    0, //玩家角度
	}

	return p
}

//玩家可以和对端客户端发送消息
func (p *Player) sendMsg(msgID uint32, port_struct proto.Message) error {
	//将porto结构体,转换 二进制数据
	binary_proto_data, err := proto.Marshal(port_struct)
	if err != nil {
		fmt.Println("Marshal err : ", err)
		return err
	}

	//再调用zinx原生的connecton.Send（msgID, 二进制数据）
	if err := p.Conn.Send(msgID, binary_proto_data); err != nil {
		fmt.Println("Player send err :", err)
		return err
	}

	return nil
}

/*
 服务器给客户端发送玩家初始ID
 */
func (p *Player) ReturnPid() {
	//定义msg:ID  proto 数据结构
	proto_msg := &pb.SyncPid{
		Pid: p.Pid,
	}

	//发送给客户端
	err := p.sendMsg(1, proto_msg)
	if err != nil {
		fmt.Println("sendMsg err :", err)
		return
	}
}

//服务器给客户端发送一个玩家的初始化位置信息
func (p *Player) ReturnPlayerPosition() {
	//组建MsgID:200消息
	proto_msg := &pb.BroadCast{
		Pid: p.Pid,
		Tp:  2,
		Data: &pb.BroadCast_P{
			P: &pb.Position{
				X: p.X,
				Y: p.Y,
				Z: p.Z,
				V: p.V,
			},
		},
	}

	//将消息,发送给客户端
	p.sendMsg(200, proto_msg)
}

//将聊天数据广播给全部的在线玩家
func (p *Player) SendTalkMsgToAll(content string) {
	//定义一个广播消息的数据类型
	proto_msg := &pb.BroadCast{
		Pid: p.Pid,
		Tp:  1,
		Data: &pb.BroadCast_Content{
			Content: content,
		},
	}

	//获取的在线玩家
	players := WorldMgrObj.GetAllPlayers()

	//循环得到所有玩家
	for _, player := range players {
		player.sendMsg(200, proto_msg)
	}
}

//获取当前玩家周围九宫格的玩家
func (p *Player) GetSurroundingPlayers() []*Player {
	//获取当前玩家周围九宫格的玩家pid
	pids := WorldMgrObj.AoiMgr.GetSurroundPIDsByPos(p.X, p.Z)
	//fmt.Println("Surrounding player = ", pids)

	players := make([]*Player, 0, len(pids))

	//循环pids获取其中的pid,获取对应的玩家对象
	for _, pid := range pids {
		players = append(players, WorldMgrObj.GetPlayerByPid(int32(pid)))
	}

	return players
}

//将自己的信息发送给周围玩家
func (p *Player) SyncSurrounding() {
	//获取当前玩家周围九宫格的玩家有那些
	players := p.GetSurroundingPlayers()
	//players := WorldMgrObj.GetAllPlayers()

	//构建一个广播消息200， 循环全部players 分别给player对应的客户端发送200消息（让其他玩家看见当前玩家）
	proto_msg := &pb.BroadCast{
		Pid: p.Pid,
		Tp:  2,
		Data: &pb.BroadCast_P{
			P: &pb.Position{
				X: p.X,
				Y: p.Y,
				Z: p.Z,
				V: p.V,
			},
		},
	}

	//将当前玩家id和位置发送到周边的玩家
	for _, player := range players {
		player.sendMsg(200, proto_msg)
	}

	//将其他玩家告诉当前玩家  （让当前玩家能够看见周边玩家的坐标）
	//构建一个202消息  players的信息 告知当前玩家 p.send(202, ... )
	//得到当前周边玩家的player集合message Player
	players_proto_msg := make([]*pb.Player, 0, len(players))
	for _, player := range players {
		//制作一个message Player 消息
		p := &pb.Player{
			Pid: player.Pid,
			P: &pb.Position{
				X: player.X,
				Y: player.Y,
				Z: player.Z,
				V: player.V,
			},
		}

		players_proto_msg = append(players_proto_msg, p)
	}

	//创建一个message SyncPlayers
	syncPlayer_proto_msg := &pb.SyncPlayers{
		Ps: players_proto_msg[:],
	}

	//将当前的周边全部玩家消息 发送给当前客户端
	p.sendMsg(202, syncPlayer_proto_msg)
}

//更新广播当前玩家的最新消息
func (p *Player) UpdatePosition(x, y, z, v float32) {
	//需要当前玩家的最新坐标
	p.X = x
	p.Y = y
	p.Z = z
	p.V = v

	//组建广播 proto协议 msgId:200  TP= 4
	proto_msg := &pb.BroadCast{
		Pid: p.Pid,
		Tp:  4,
		Data: &pb.BroadCast_P{
			P: &pb.Position{
				X: p.X,
				Y: p.Y,
				Z: p.Z,
				V: p.V,
			},
		},
	}

	//获取当前玩家周围九宫格内的玩家
	players := p.GetSurroundingPlayers()

	//依次调用players 将位置信息发送过去
	for _, player := range players {
		//每个玩家都会给各自的 client客户端发送200消息
		player.sendMsg(200, proto_msg)
	}
}

//下线业务函数
func (p *Player) OffLine() {
	//得到当前玩家周围九宫格有那些玩家
	players := p.GetSurroundingPlayers()

	//制作一个消息MsgId:201
	proto_msg := &pb.SyncPid{
		Pid: p.Pid,
	}

	//给周边玩家广播消息 msgId: 201
	for _, player := range players {
		//客户端会从视野中将当前玩家删除
		player.sendMsg(201, proto_msg)
	}

	//将该下线玩家从世界管理器移除
	WorldMgrObj.RemovePlayerByPid(p.Pid)

	//将该下线玩家从世界地图删除
	WorldMgrObj.AoiMgr.RemoteFromGridbyPos(int(p.Pid), p.X, p.Z)

}
