package state

// Commiter 应用层各个状态同步接口
type Commiter interface {
	Commit() ([]byte, error)
}
