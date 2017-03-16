package mytest

import (
	"fmt"
	"net"
	"net/rpc"
	"strconv"
	"sync"
	"testing"
	"time"
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

type GetRequest struct {
	getArgs  *GetArgs
	getReply *GetReply
}

func (kv *KV) Get(getArgs *GetArgs, getReply *GetReply) error {
	kv.lock.Lock()
	getReply.Val = kv.kvMap[getArgs.Key]
	kv.lock.Unlock()
	return nil
}

func (kv *KV) Put(putArgs *PutArgs, putReply *PutReply) error {
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

func client(key string, val string) {
	client, err := rpc.Dial("tcp", "127.0.0.1:51234")
	if err != nil {
		fmt.Printf("client Dial:%v", err)
	}
	putArgs := &PutArgs{key, val}
	putReply := &PutReply{}
	client.Call("KV.Put", putArgs, putReply)

	getArgs := &GetArgs{key}
	getReply := &GetReply{""}
	callErr := client.Call("KV.Get", getArgs, getReply)
	if callErr == nil {
		fmt.Printf("getReply == %v\n", getReply.Val)

	} else {
		fmt.Printf("callErr: %v\n", callErr)
	}
}

/**
 * 由于linux默认连接数为1024，因此需要重新设置
 * echo 'ulimit -n 655350' >> /etc/profile
 * echo "$USERNAME hard nofile 655350" >> /etc/security/limits.conf
 * @author tantexian(https://my.oschina.net/tantexian/blog)
 * @since 2017/3/13
 * @params
 */
func TestRpc(t *testing.T) {
	server()
	wg := sync.WaitGroup{}
	for i := 0; i < 6400; i++ {
		wg.Add(1)
		go func(index int) {
			defer wg.Done()
			client(strconv.Itoa(index), time.Now().String())
		}(i)
	}
	wg.Wait()
}
