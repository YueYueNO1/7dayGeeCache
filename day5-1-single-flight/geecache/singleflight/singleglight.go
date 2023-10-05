package singleflight

import "sync"

// 正在进行中、或已经结束的请求
type call struct {
	wg  sync.WaitGroup //避免重入
	val interface{}
	err error
}

// 是singleflight的主数据结构，管理不同key的请求（call）
type Group struct {
	mu sync.Mutex
	m  map[string]*call
}

/*
接收两个参数，key 和 fn。
作用：针对相同的key，无论Do被调用多少次，函数fn都只会被调用一次，等待fn调用结束了，返回返回值或错误
！并发的gorountine不需要消息传递，非常适合sync.WaitGroup
*/
func (g *Group) Do(key string, fn func() (interface{}, error)) (interface{}, error) {
	g.mu.Lock()     //确保同时只有一个gorountine可以访问Group
	if g.m == nil { //没有，新建
		g.m = make(map[string]*call)
	}

	if c, ok := g.m[key]; ok { //判断有没有重复，重复就等
		g.mu.Unlock() //上面lock了
		c.wg.Wait()   //阻塞等待，直到call结构体中的等待组计数器归零
		return c.val, c.err
	}

	c := new(call) //不重复key，新建
	c.wg.Add(1)    //用于向等待组(WaitGroup)添加一个等待任务
	g.m[key] = c   //将其添加到Group的映射表中，以便其他goroutine可以共享该结果
	g.mu.Unlock()  //允许其他gorountine可以共享该结果

	c.val, c.err = fn() //执行函数fn的操作，返回两个返回值
	c.wg.Done()         //等待组计数器减一

	g.mu.Lock()      //锁定group
	delete(g.m, key) //删除映射表中对应的key，以清理已经完成的缓存
	g.mu.Unlock()

	return c.val, c.err //返回call结构体中存储的结果

}
