package core

import "sync"

//当前世界地图边界的参数
const (
	AOI_MIN_X  int = 85
	AOI_MAX_X  int = 410
	AOI_CNTS_X int = 10
	AOI_MIN_Y  int = 75
	AOI_MAX_Y  int = 400
	AOI_CNTS_Y int = 20
)

/*
	当前地图的世界管理模块
*/

type WorldManager struct {
	//当前在线的玩家集合
	Players map[int32]*Player
	//保护数据的锁
	plock sync.RWMutex
	//AOIManager===>当前地图的管理器
	AoiMgr *AOIManager
}

//初始化方法
func NewWorldManager() *WorldManager {
	wm := &WorldManager{
		Players: make(map[int32]*Player),
		AoiMgr:  NewAOIManager(AOI_MIN_X, AOI_MAX_X, AOI_CNTS_X, AOI_MIN_Y, AOI_MAX_Y, AOI_CNTS_Y),
	}

	return wm
}

//对外提供一个全局的世界管理模块指针
var WorldMgrObj *WorldManager

//设立init函数,执行main函数前会执行init函数(前提是要在main函数中导包)
func init() {
	WorldMgrObj = NewWorldManager()
}

//添加一个玩家
func (wm *WorldManager) AddPlayer(player *Player) {
	//枷锁
	wm.plock.Lock()
	//添加玩家
	wm.Players[player.Pid] = player
	//解锁
	wm.plock.Unlock()

	//加入到世界地图中======>通过坐标 将玩家ID--->pid  添加到 当前坐标所定位的格子grid中
	wm.AoiMgr.AddToGridByPos(int(player.Pid), player.X, player.Z)
}

//删除玩家
func (wm *WorldManager) RemovePlayerByPid(pid int32) {
	//从世界鬼管理表中删除
	wm.plock.Lock()

	//通过pid获取player
	player := wm.Players[pid]

	//从世界地图中删除玩家
	wm.AoiMgr.RemoteFromGridbyPos(int(pid), player.X, player.Z)

	//从玩家管理表中删除
	delete(wm.Players, pid)

	wm.plock.Unlock()
}

//通过一个玩家ID得到一个Player对象
func (wm *WorldManager) GetPlayerByPid(pid int32) *Player {
	wm.plock.RLock()
	p := wm.Players[pid]
	wm.plock.RUnlock()

	return p
}

//获取全部玩家在线集合
func (wm *WorldManager) GetAllPlayers() []*Player {
	//定义容器
	players := make([]*Player, 0)
	//加读锁
	wm.plock.RLock()
	//循环获取玩家列表中的玩家信息
	for _, player := range wm.Players {
		players = append(players, player)
	}
	//解锁
	wm.plock.RUnlock()

	return players
}
