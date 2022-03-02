package cache

import pb "cache/mycache/cachepb"

// PeerPicker 实现根据key选择相应节点
type PeerPicker interface {
	PickPeer(key string) (peer PeerGetter, ok bool)
}

// PeerGetter 实现根据group和key返回缓存值
type PeerGetter interface {
	Get(in *pb.Request, out *pb.Response) error
}
