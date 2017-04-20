package raft

//
// this is an outline of the API that raft must expose to
// the service (or tester). see comments below for
// each of these functions for more details.
//
// rf = Make(...)
//   create a new Raft server.
// rf.Start(command interface{}) (index, term, isleader)
//   start agreement on a new log entry
// rf.GetState() (term, isLeader)
//   ask a Raft for its current term, and whether it thinks it is leader
// ApplyMsg
//   each time a new entry is committed to the log, each Raft peer
//   should send an ApplyMsg to the service (or tester)
//   in the same server.
//

import "sync"
import (
	"fmt"
	"labrpc"
	"sync/atomic"
	"time"
)

// import "bytes"
// import "encoding/gob"

//
// as each Raft peer becomes aware that successive log entries are
// committed, the peer should send an ApplyMsg to the service (or
// tester) on the same server, via the applyCh passed to Make().
//
type ApplyMsg struct {
	Index       int
	Command     interface{}
	UseSnapshot bool   // ignore for lab2; only used in lab3
	Snapshot    []byte // ignore for lab2; only used in lab3
}

// Log 日志条目数据结构,
// contains command for state machine, and term when entry was received by leader (first index is 1)
// Author: tantexian, <my.oschina.net/tantexian>
// Since: 2017/327
type LogEntry struct {
	command interface{}
	term    int
	index   int
}

// 用于定义Raft中，代表候选人id为空时的none值 Add: tantexian, <my.oschina.net/tantexian> Since: 2017/3/28
const none = -1

// 增加Raft的status Add: tantexian, <my.oschina.net/tantexian> Since: 2017/3/28
type RaftState int

const (
	StateFollower  RaftState = iota // 跟随者
	StateCandidate                  // 候选人
	StateLeader                     // 领导者
)

//
// A Go object implementing a single Raft peer.
//
type Raft struct {
	mu        sync.Mutex          // Lock to protect shared access to this peer's state
	peers     []*labrpc.ClientEnd // RPC end points of all peers
	persister *Persister          // Object to hold this peer's persisted state
	me        int                 // this peer's index into peers[]

	// Your data here (2A, 2B, 2C).
	// Look at the paper's Figure 2 for a description of what
	// state a Raft server must maintain.
	// Add: tantexian, <my.oschina.net/tantexian> Since: 2017/3/27
	// Persistent state on all servers:(Updated on stable storage before responding to RPCs)
	currentTerm int         // 任期（第一次启动初始化为0，后续单调递增）
	votedFor    int         // 当前任期接收选票的候选人id(如果没有则设置为none，用常量-1代替)
	logs        []*LogEntry //log entries; each entry contains command for state machine, and term when entry was received by leader (first index is 1)

	// Volatile state on all servers:
	commitIndex int // index of highest log entry known to be committed (initialized to 0, increases monotonically)
	lastApplied int // index of highest log entry applied to state machine (initialized to 0, increases monotonically)

	// Volatile state on leaders:(Reinitialized after election)
	nextIndex  []int // for each server, index of the next log entry to send to that server (initialized to leader last log index + 1)
	matchIndex []int // for each server, index of highest log entry known to be replicated on server(initialized to 0, increases monotonically)

	state RaftState // 此处用type RaftState int表示跟随者，候选者、领导人

	heartbeatCh chan string // 此处用于接收其他raft的心跳信息
}

// GetLastLogTerm 获取raft最后日志条目的任期
// Author: tantexian, <my.oschina.net/tantexian>
// Since: 2017/3/28
func (rf *Raft) GetLastLogTerm() (term int) {
	len := len(rf.logs)
	if len == 0 {
		fmt.Printf("raft[%v] has no logs", rf.me)
		term = none
		return term
	}
	term = rf.logs[len-1].term
	return term
}

// GetLastLogIndex 获取raft最后日志条目的index
// Author: tantexian, <my.oschina.net/tantexian>
// Since: 2017/3/28
func (rf *Raft) GetLastLogIndex() (index int) {
	len := len(rf.logs)
	if len == 0 {
		fmt.Printf("raft[%v] has no logs\n", rf.me)
		index = none
		return index
	}
	index = rf.logs[len-1].index
	return index
}

// return currentTerm and whether this server
// believes it is the leader.
func (rf *Raft) GetState() (int, bool) {

	var term int
	var isleader bool
	// Your code here (2A).
	term = rf.currentTerm
	if rf.state == StateLeader {
		isleader = true
	}
	return term, isleader
}

//
// save Raft's persistent state to stable storage,
// where it can later be retrieved after a crash and restart.
// see paper's Figure 2 for a description of what should be persistent.
//
func (rf *Raft) persist() {
	// Your code here (2C).
	// Example:
	// w := new(bytes.Buffer)
	// e := gob.NewEncoder(w)
	// e.Encode(rf.xxx)
	// e.Encode(rf.yyy)
	// data := w.Bytes()
	// rf.persister.SaveRaftState(data)
}

//
// restore previously persisted state.
//
func (rf *Raft) readPersist(data []byte) {
	// Your code here (2C).
	// Example:
	// r := bytes.NewBuffer(data)
	// d := gob.NewDecoder(r)
	// d.Decode(&rf.xxx)
	// d.Decode(&rf.yyy)
	if data == nil || len(data) < 1 {
		// bootstrap without any state?
		return
	}
}

//
// example RequestVote RPC arguments structure.
// field names must start with capital letters!
//
type RequestVoteArgs struct {
	// Your data here (2A, 2B).
	// Add: tantexian, <my.oschina.net/tantexian> Since: 2017/3/27
	Term         int // 候选人的任期
	CandidateId  int // 候选人id
	LastLogIndex int // 候选人最后日志条目索引
	LastLogTerm  int // 候选人最后日志条目任期
}

//
// example RequestVote RPC reply structure.
// field names must start with capital letters!
//
type RequestVoteReply struct {
	// Your data here (2A).
	// Add: tantexian, <my.oschina.net/tantexian> Since: 2017/3/27
	Term        int  // 当前的任期（用于候选人更新自己的任期）
	VoteGranted bool // true代表候选人接收到选票
}

//
// example RequestVote RPC handler.
//
func (rf *Raft) RequestVote(args *RequestVoteArgs, reply *RequestVoteReply) {
	// Your code here (2A, 2B).
	reply.Term = rf.currentTerm
	if args.Term < rf.currentTerm {
		reply.VoteGranted = false
	}

	if rf.votedFor == -1 || // 如果当前候选人votedFor为-1
			rf.votedFor == args.CandidateId || // 或者等于当前请求投票的候选人
			(args.LastLogTerm == rf.GetLastLogTerm() && args.LastLogIndex == rf.GetLastLogIndex()) { //或者当前请求投票的候选人的日志和rf的日志一样新，则投票
		reply.VoteGranted = true
	}
}

//
// example code to send a RequestVote RPC to a server.
// server is the index of the target server in rf.peers[].
// expects RPC arguments in args.
// fills in *reply with RPC reply, so caller should
// pass &reply.
// the types of the args and reply passed to Call() must be
// the same as the types of the arguments declared in the
// handler function (including whether they are pointers).
//
// The labrpc package simulates a lossy network, in which servers
// may be unreachable, and in which requests and replies may be lost.
// Call() sends a request and waits for a reply. If a reply arrives
// within a timeout interval, Call() returns true; otherwise
// Call() returns false. Thus Call() may not return for a while.
// A false return can be caused by a dead server, a live server that
// can't be reached, a lost request, or a lost reply.
//
// Call() is guaranteed to return (perhaps after a delay) *except* if the
// handler function on the server side does not return.  Thus there
// is no need to implement your own timeouts around Call().
//
// look at the comments in ../labrpc/labrpc.go for more details.
//
// if you're having trouble getting RPC to work, check that you've
// capitalized all field names in structs passed over RPC, and
// that the caller passes the address of the reply struct with &, not
// the struct itself.
//
func (rf *Raft) sendRequestVote(server int, args *RequestVoteArgs, reply *RequestVoteReply) bool {
	ok := rf.peers[server].Call("Raft.RequestVote", args, reply)
	return ok
}

// AppendEntriesArgs AppendEntries函数的RPC请求参数
// Author: tantexian, <my.oschina.net/tantexian>
// Since: 2017/3/28
type AppendEntriesArgs struct {
	Term         int         // 领导人的任期
	LeaderId     int         // 领导人id，so follower can redirect clients
	PrevLogIndex int         // index of log entry immediately preceding new ones
	PrevLogTerm  int         // term of prevLogIndex entry
	Entries      []*LogEntry // log entries to store (empty for heartbeat; may send more than one for efficiency)
	LeaderCommit int         // leader’s commitIndex
}

// AppendEntriesReply AppendEntries函数的RPC回复参数
// Author: tantexian, <my.oschina.net/tantexian>
// Since: 2017/3/28
type AppendEntriesReply struct {
	Term    int  // currentTerm, for leader to update itself
	Success bool // true if follower contained entry matching prevLogIndex and prevLogTerm
}

// AppendEntries AppendEntries RPC handler.
// Author: tantexian, <my.oschina.net/tantexian>
// Since: 2017/3/29
func (rf *Raft) AppendEntries(args *AppendEntriesArgs, reply *AppendEntriesReply) {
	reply.Term = rf.currentTerm
	if args.Entries == nil {
		rf.heartbeatCh <- time.Now().String()
		fmt.Printf("raft:%v heartbeat fill to heartbeatCh.\n", rf.me)
	}
	if rf.GetLastLogIndex() == args.PrevLogIndex && rf.GetLastLogTerm() == args.PrevLogTerm {
		reply.Success = true
	}
	if args.Term < rf.currentTerm {
		reply.Success = false
	}
	length := len(rf.logs)
	if length >= args.PrevLogIndex && rf.logs[args.PrevLogIndex].term != args.PrevLogTerm {
		reply.Success = false
	}
	entries := args.Entries
	if entries != nil {
		for i := 0; i < len(entries); i++ {
			for j := len(rf.logs); j > 0; j-- {
				if rf.logs[j].index == entries[i].index && rf.logs[j].term != entries[i].term {
					rf.logs = rf.logs[0: len(rf.logs)-2]
				}
			}
		}
		for i := 0; i < len(entries); i++ {
			rf.logs = append(rf.logs, entries[i])
		}
	}
	if args.LeaderCommit > rf.commitIndex {
		index1 := entries[len(entries)-1].index
		if index1 < args.LeaderCommit {
			rf.commitIndex = index1
		} else {
			rf.commitIndex = args.LeaderCommit
		}
	}

}

// AppendEntries send a RequestVote RPC to a server.添加日志条目，Invoked by leader to replicate log entries (§5.3); also used as heartbeat (§5.2).
// Author: tantexian, <my.oschina.net/tantexian>
// Since: 2017/3/28
func (rf *Raft) sendAppendEntries(server int, args *AppendEntriesArgs, reply *AppendEntriesReply) bool {
	ok := rf.peers[server].Call("Raft.AppendEntries", args, reply)
	return ok
}

//
// the service using Raft (e.g. a k/v server) wants to start
// agreement on the next command to be appended to Raft's log. if this
// server isn't the leader, returns false. otherwise start the
// agreement and return immediately. there is no guarantee that this
// command will ever be committed to the Raft log, since the leader
// may fail or lose an election.
//
// the first return value is the index that the command will appear at
// if it's ever committed. the second return value is the current
// term. the third return value is true if this server believes it is
// the leader.
// 询问Raft启动一个进程将命令添加到副本log中
func (rf *Raft) Start(command interface{}) (int, int, bool) {
	index := -1
	term := -1
	isLeader := true

	// Your code here (2B).

	return index, term, isLeader
}

//
// the tester calls Kill() when a Raft instance won't
// be needed again. you are not required to do anything
// in Kill(), but it might be convenient to (for example)
// turn off debug output from this instance.
//
func (rf *Raft) Kill() {
	// Your code here, if desired.
}

// 处理状态变化及选举 Add: tantexian, <my.oschina.net/tantexian> Since: 2017/3/29
func (rf *Raft) handleElection() {
	for {
		switch rf.state {
		case StateFollower:
			select {
			case heartbeat := <-rf.heartbeatCh:
				fmt.Printf("raft[%v] receive heartbeat %v.\n", rf.me, heartbeat)
			// 如果超过500毫秒未收到心跳信息，则切换状态为候选人
			case <-time.After(500 * time.Millisecond):
				rf.currentTerm++
				rf.votedFor = rf.me
				rf.state = StateCandidate
				fmt.Printf("raft[%v] does not receive heartbeat, and switch to StateCandidate.\n", rf.me)
			}
		case StateCandidate:
			var voteNum int32 = 0
			requestVoteArgs := &RequestVoteArgs{
				Term:         rf.currentTerm,
				CandidateId:  rf.me,
				LastLogIndex: rf.GetLastLogIndex(),
				LastLogTerm:  rf.GetLastLogTerm()}
			// 并行向其他服务器发送请求投票RPCs
			rfLen := len(rf.peers)
			for i := 0; i < rfLen; i++ {
				wg := sync.WaitGroup{}
				wg.Add(rfLen)

				if i == rf.me { // 不需要向自己发送请求投票RPCs
					wg.Done()
					continue
				}
				go func(server int) {
					requestVoteReply := &RequestVoteReply{}
					ok := rf.sendRequestVote(server, requestVoteArgs, requestVoteReply)
					if ok && requestVoteReply.VoteGranted {
						atomic.AddInt32(&voteNum, 1)
					}
					wg.Done()
				}(i)
				// 当获得超过半数选票时，立马切换为leader
				if rf.state != StateLeader {
					if int(voteNum) > rfLen/2 {
						rf.state = StateLeader
						fmt.Printf("raft[%v] has get majority votes, and switch to StateLeader.\n", rf.me)
					}
				}
				wg.Wait()
			}
		case StateLeader:
			for {
				rfLen := len(rf.peers)
				for i := 0; i < rfLen; i++ {
					appendEntriesArgs := &AppendEntriesArgs{
						Term:     rf.currentTerm,
						LeaderId: rf.me,
						Entries:  nil}
					appendEntriesReply := &AppendEntriesReply{}

					go func(server int) {
						rf.sendAppendEntries(server, appendEntriesArgs, appendEntriesReply)
					}(i)
				}
				time.Sleep(10 * time.Millisecond)
			}

		}
	}
}

//
// the service or tester wants to create a Raft server. the ports
// of all the Raft servers (including this one) are in peers[]. this
// server's port is peers[me]. all the servers' peers[] arrays
// have the same order. persister is a place for this server to
// save its persistent state, and also initially holds the most
// recent saved state, if any. applyCh is a channel on which the
// tester or service expects Raft to send ApplyMsg messages.
// Make() must return quickly, so it should start goroutines
// for any long-running work.
//
func Make(peers []*labrpc.ClientEnd, me int,
		persister *Persister, applyCh chan ApplyMsg) *Raft {
	rf := &Raft{}
	rf.peers = peers
	rf.persister = persister
	rf.me = me

	// Your initialization code here (2A, 2B, 2C).
	// Add: tantexian, <my.oschina.net/tantexian> Since: 2017/3/27
	rf.state = StateFollower // 初始状态为跟随者
	rf.currentTerm = 0
	rf.votedFor = none
	rf.logs = make([]*LogEntry, 0)

	rf.commitIndex = 0
	rf.lastApplied = 0

	// 使用goroutine，处理状态变化及选举
	go rf.handleElection()

	// initialize from state persisted before a crash
	rf.readPersist(persister.ReadRaftState())

	return rf
}
