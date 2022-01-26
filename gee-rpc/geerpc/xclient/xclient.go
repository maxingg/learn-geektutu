package xclient

import (
	"context"
	"geektutu/gee-rpc/geerpc"
	"io"
	"reflect"
	"sync"
)

type XClient struct {
	d       Discovery
	mode    SelectMode
	opt     *geerpc.Option
	mu      sync.Mutex
	clients map[string]*geerpc.Client
}

func (xc *XClient) Close() error {
	xc.mu.Lock()
	defer xc.mu.Unlock()
	for key, client := range xc.clients {
		_ = client.Close()
		delete(xc.clients, key)
	}
	return nil
}

func NewXClient(d Discovery, mode SelectMode, opt *geerpc.Option) *XClient {
	return &XClient{
		d:       d,
		mode:    mode,
		opt:     opt,
		clients: make(map[string]*geerpc.Client),
	}
}

var _ io.Closer = (*XClient)(nil)

func (xc *XClient) Call(ctx context.Context, serviceMethod string, args, reply interface{}) error {
	rpcAddr, err := xc.d.Get(xc.mode)
	if err != nil {
		return err
	}
	return xc.call(rpcAddr, ctx, serviceMethod, args, reply)
}

func (xc *XClient) call(addr string, ctx context.Context, method string, args, reply interface{}) error {
	client, err := xc.dial(addr)
	if err != nil {
		return err
	}
	return client.Call(ctx, method, args, reply)
}

func (xc *XClient) dial(addr string) (*geerpc.Client, error) {
	xc.mu.Lock()
	defer xc.mu.Unlock()
	client, ok := xc.clients[addr]
	if ok && !client.IsAvailable() {
		_ = client.Close()
		delete(xc.clients, addr)
		client = nil
	}
	if client == nil {
		var err error
		client, err = geerpc.XDial(addr, xc.opt)
		if err != nil {
			return nil, err
		}
		xc.clients[addr] = client
	}
	return client, nil
}

func (xc *XClient) Broadcast(ctx context.Context, serviceMethod string, args, reply interface{}) error {
	servers, err := xc.d.GetAll()
	if err != nil {
		return err
	}
	var wg sync.WaitGroup
	var mu sync.Mutex
	var e error
	// reply 若为nil，reply无需设置值
	replyDone := reply == nil
	ctx, cancel := context.WithCancel(ctx)
	for _, rpcAddr := range servers {
		wg.Add(1)
		go func(rpcAddr string) {
			defer wg.Done()
			// 广播是并发，如果reply需要填值，每一个调用需要一个克隆对象
			var cloneReply interface{}
			if reply != nil {
				cloneReply = reflect.New(reflect.Indirect(reflect.ValueOf(reply)).Type()).Interface()
			}
			err := xc.call(rpcAddr, ctx, serviceMethod, args, cloneReply)
			mu.Lock()
			// 快速捕获到第一个报错并停止
			if err != nil && e == nil {
				e = err
				cancel()
			}
			if err == nil && !replyDone {
				reflect.ValueOf(reply).Elem().Set(reflect.ValueOf(cloneReply).Elem())
				replyDone = true
			}
			mu.Unlock()
		}(rpcAddr)
	}
	wg.Wait()
	return e
}
