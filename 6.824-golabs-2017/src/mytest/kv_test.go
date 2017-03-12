package mytest

import (
	"sync"
	"net/rpc"
	"net"
	"fmt"
	"testing"
)

type PutArgs struct {
	Key string
	Val string
}

type GetArgs struct {
	Key string
}

type GetReply struct {
	Val string
}

type PutReply struct {
	Val string
}

type KV struct {
	lock  sync.Mutex
	kvMap map[string]string
}

func (kv *KV) Get(getArgs *GetArgs, getReply *GetReply) (error) {
	kv.lock.Lock()
	getReply.Val = kv.kvMap[getArgs.Key]
	kv.lock.Unlock()
	return nil
}

func (kv *KV) Put(putArgs *PutArgs, putReply *PutReply) (error) {
	kv.lock.Lock()
	kv.kvMap[putArgs.Key] = putArgs.Val
	putReply.Val = putArgs.Val
	kv.lock.Unlock()
	return nil
}

func server() {
	kv := new(KV)
	kv.kvMap = map[string]string{}
	server := rpc.NewServer()
	server.Register(kv)
	listen, err := net.Listen("tcp", "127.0.0.1:51234")
	if err != nil {
		fmt.Printf("listen error:%v", err)
	}
	go func() {
		for {
			conn, err := listen.Accept()
			if err == nil {
				go server.ServeConn(conn)
			} else {
				// 连接有错误则跳出循环，结束
				break
			}
		}
		listen.Close()
		fmt.Printf("Server done\n")
	}()
}

func client() {
	client, err := rpc.Dial("tcp", "127.0.0.1:51234")
	if err != nil {
		fmt.Printf("client Dial:%v", err)
	}
	putArgs := &PutArgs{"a", "a"}
	putReply := &PutReply{}
	client.Call("KV.Put", putArgs, putReply)

	getArgs := &GetArgs{"a"}
	getReply := &GetReply{""}
	callErr := client.Call("KV.Get", getArgs, getReply)
	if callErr == nil {
		fmt.Printf("getReply == %v\n", getReply.Val)

	} else {
		fmt.Printf("callErr: %v\n", callErr)
	}
}

func TestRpc(t *testing.T) {
	server()
	client()
}

