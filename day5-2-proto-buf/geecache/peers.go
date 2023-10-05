package geecache

import pb "geecache/geecachepb"

type PeerPicker interface {
	//根据传入的key选择相应节点
	PickPeer(key string) (peer PeerGetter, ok bool)
}

// 对应上述流程的HTTP客户端
type PeerGetter interface {
	//用于从对应的group中查找缓存值。
	Get(in *pb.Request, out *pb.Response) error
}
