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
	"labrpc"
	"sync/atomic"
	"time"
	"math/rand"
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
	Command interface{}
	Term    int
	Index   int
}

// 用于定义Raft中，代表候选人id为空时的none值 Add: tantexian, <my.oschina.net/tantexian> Since: 2017/3/28
const none = -1

// 心跳发送间隔时间，心跳接收超时时间，候选人选举之后随机延时基础时间 (单位毫秒）
const (
	// ElectionTick is the number of Node.Tick invocations that must pass between
	// elections. That is, if a follower does not receive any message from the
	// leader of current term before ElectionTick has elapsed, it will become
	// candidate and start an election. ElectionTick must be greater than
	// HeartbeatTick. We suggest ElectionTick = 10 * HeartbeatTick to avoid
	// unnecessary leader switching.
	ElectionTimeoutTick int = 400
	// HeartbeatTick is the number of Node.Tick invocations that must pass between
	// heartbeats. That is, a leader sends heartbeat messages to maintain its
	// leadership every HeartbeatTick ticks.
	HeartbeatBroadcastTick int = 100
)

// 增加Raft的state Add: tantexian, <my.oschina.net/tantexian> Since: 2017/3/28
type RaftState int

const (
	StateFollower  RaftState = iota // 跟随者
	StateCandidate                  // 候选人
	StateLeader                     // 领导者
)

var raftStateMap = [...]string{
	"StateFollower",
	"StateCandidate",
	"StateLeader",
}

func (raftState RaftState) String() string {
	return raftStateMap[raftState]
}

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

	heartbeatBroadcastTime int // 领导者发送心跳的间隔时间,更多请查看：常量HeartbeatBroadcastTick

	electionTimeout     int // 选举超时时间（base值）,更多请查看：常量ElectionTimeoutTick
	randElectionTimeout int
}

// GetLastLogTerm 获取raft最后日志条目的任期
// Author: tantexian, <my.oschina.net/tantexian>
// Since: 2017/3/28
func (rf *Raft) GetLastLogTerm() (term int) {
	len := len(rf.logs)
	if len == 0 {
		// DPrintf("raft[%v] has no logs\n", rf.me)
		// 如果当前没有日志条目，则任期号为none -1
		term = none
		return term
	}
	term = rf.logs[len-1].Term
	return term
}

// GetLastLogIndex 获取raft最后日志条目的index
// Author: tantexian, <my.oschina.net/tantexian>
// Since: 2017/3/28
func (rf *Raft) GetLastLogIndex() (index int) {
	len := len(rf.logs)
	if len == 0 {
		// DPrintf("raft[%v] has no logs\n", rf.me)
		// 如果当前没有日志条目，则索引为none -1
		index = none
		return index
	}
	index = rf.logs[len-1].Index
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
	//DPrintf("## RequestVote start %+v , args == %+v replay == %+v\n", rf, args, reply)
	// Your code here (2A, 2B).

	// 接收选票时候，发现对方的Term大于自己Term,则设置自己的Term，且切换为StateFollower
	if args.Term > rf.currentTerm {
		DPrintf("[SWITCHSTATE: ->StateFollower]: raft[%v] receive RequestVote other raft Term(%v) > self Term(%v).\n", rf.me, args.Term, rf.currentTerm)
		rf.mu.Lock()
		rf.currentTerm = args.Term
		rf.state = StateFollower
		rf.votedFor = none
		rf.mu.Unlock()
	}

	reply.Term = rf.currentTerm
	// Reply false if term < currentTerm (§5.1)
	if args.Term < rf.currentTerm {
		reply.VoteGranted = false
		// DPrintf("** fails RequestVote  [%v to %v term %v] args == %v reply==%v\n", rf.me, args.CandidateId, rf.currentTerm, args, reply)
		return
	}

	// 如果当前候选人votedFor为none 或者votedFor等于当前请求投票的候选人
	// If votedFor is null or candidateId, and candidate’s log is at
	// least as up-to-date as receiver’s log, grant vote (§5.2, §5.4)
	if (rf.votedFor == none || rf.votedFor == args.CandidateId) &&
	//且当前请求投票的候选人的日志和rf的日志(至少)一样新，则投票
			(args.LastLogTerm >= rf.GetLastLogTerm() && args.LastLogIndex >= rf.GetLastLogIndex()) {
		rf.mu.Lock()
		rf.votedFor = args.CandidateId
		rf.mu.Unlock()
		reply.VoteGranted = true
		DPrintf("[VOTEING]: success RequestVote [%v to %v term %v] args == %v reply==%v\n", rf.me, args.CandidateId, rf.currentTerm, args, reply)
		return
	} else {
		// DPrintf("** fails else RequestVote [%v to %v term %v] args == %v reply==%v\n", rf.me, args.CandidateId, rf.currentTerm, args, reply)
		reply.VoteGranted = false
		return
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
func (rf *Raft) AppendEntries(args AppendEntriesArgs, reply *AppendEntriesReply) {
	//DPrintf("args====%+v\n", args)
	// currentTerm, for leader to update itself
	reply.Term = rf.currentTerm

	// 1. Reply false if term < currentTerm (§5.1)
	if args.Term < rf.currentTerm {
		DPrintf("[rf=%+v] AppendEntries args.Term %v < rf.currentTerm %v\n", rf.me, args.Term, rf.currentTerm)
		reply.Success = false
		return
	} else { // 接收RPC日志时候，发现对方的Term大于等于自己Term,则设置自己的Term，且切换为StateFollower
		DPrintf("[SWITCHSTATE: ->StateFollower]: raft[%v] receive AppendEntries other raft Term(%v) > self Term(%v).\n", rf.me, args.Term, rf.currentTerm)
		rf.mu.Lock()
		rf.currentTerm = args.Term
		rf.state = StateFollower
		rf.votedFor = none
		rf.mu.Unlock()
	}

	// 如果Entries为空，则表示心跳信息
	if len(args.Entries) == 0 {
		// DPrintf("send heartbeat to %v\n", rf.me)
		rf.heartbeatCh <- time.Now().String()
		return // TODO：心跳信息则直接返回？
	}

	// 2. Reply false if log doesn’t contain an entry at prevLogIndex
	// whose term matches prevLogTerm (§5.3)
	length := len(rf.logs)
	if length > args.PrevLogIndex && rf.logs[args.PrevLogIndex].Term != args.PrevLogTerm {
		reply.Success = false
		return
	}

	// 3. If an existing entry conflicts with a new one (same index
	// but different terms), delete the existing entry and all that
	// follow it (§5.3)
	// 在论文图8中有解释该不一致情况出现场景
	// 此处一定能保证此处entries肯定不为空，为空则在心跳逻辑，return了，
	// 且contained entry matching prevLogIndex and prevLogTerm
	logEntries := args.Entries
	startIndex := 0

	for i := 0; i < len(logEntries); i++ {
		// 从prevLogIndex的后面一个开始验证
		for j := args.PrevLogIndex + 1; j < length; j++ {
			// 如果有冲突，则删除冲突及之后的所有数据
			if rf.logs[j].Index == logEntries[i].Index && rf.logs[j].Term != logEntries[i].Term {
				rf.logs = rf.logs[0: j-1]
				startIndex = i
				goto appendFlag
			}
		}
	}

	// 4. Append any new entries not already in the log
appendFlag:
	for i := startIndex; i < len(logEntries); i++ {
		rf.logs = append(rf.logs, logEntries[i])
	}

	// 5. If leaderCommit > commitIndex, set commitIndex =
	// min(leaderCommit, index of last new entry)
	if args.LeaderCommit > rf.commitIndex {
		index1 := logEntries[len(logEntries)-1].Index
		if index1 < args.LeaderCommit {
			rf.commitIndex = index1
		} else {
			rf.commitIndex = args.LeaderCommit
		}
	}

	// true if follower contained entry matching
	// prevLogIndex and prevLogTerm
	if rf.GetLastLogIndex() == args.PrevLogIndex && rf.GetLastLogTerm() == args.PrevLogTerm {
		reply.Success = true
		return
	}
}

// AppendEntries send a RequestVote RPC to a server.添加日志条目，
// Invoked by leader to replicate log entries (§5.3); also used as heartbeat (§5.2).
// Author: tantexian, <my.oschina.net/tantexian>
// Since: 2017/3/28
func (rf *Raft) sendAppendEntries(server int, args AppendEntriesArgs, reply *AppendEntriesReply) bool {
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
	if rf.state != StateLeader {
		isLeader = false
	}

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

// 当获得超过半数选票时，立马切换为leader
func (rf *Raft) gatherVotesChangeToLeader(voteNum int32) bool {
	success := true
	if rf.state != StateLeader {
		if int(voteNum) > len(rf.peers)/2 {
			rf.mu.Lock()
			rf.state = StateLeader
			rf.mu.Unlock()
			DPrintf("[SWITCHSTATE: ->StateLeader]: raft[%v] has gather majority votes[%v], term[%v], peers[%v].\n", rf.me, voteNum, rf.currentTerm, len(rf.peers))
		} else {
			success = false
		}
	}
	return success
}

func (rf *Raft) resetRandElectionTimeout() {
	// 第三种可能的结果是候选人既没有赢得选举也没有输：如果有多个跟随者同时成为候选人，
	// 那么选票可能会被瓜分以至于没有候选人可以赢得大多数人的支持。当这种情况发生的时候，
	// 每一个候选人都会超时，然后通过增加当前任期号来开始一轮新的选举。然而，没有其他机制的话，
	// 选票可能会被无限的重复瓜分。
	// 为了阻止选票起初就被瓜分，选举超时时间是从一个固定的区间（例如 150-300毫秒）随机选择。
	randTimeOut := rf.electionTimeout + rand.Int()%rf.electionTimeout
	rf.randElectionTimeout = randTimeOut
}

/*func (rf *Raft) resetCurrentTerm(term int) {
	// If RPC request or response contains term T > currentTerm:
	// set currentTerm = T, convert to follower (§5.1)
	if term > rf.currentTerm {
		rf.currentTerm = term
		if rf.state != StateFollower {
			rf.mu.Lock()
			rf.state = StateFollower
			rf.mu.Unlock()
		}
	}
}*/

// 处理状态变化及选举 Add: tantexian, <my.oschina.net/tantexian> Since: 2017/3/29
func (rf *Raft) handleElection() {
	for {
		switch rf.state {
		case StateFollower:
			rf.resetRandElectionTimeout()
			select {
			case heartbeat := <-rf.heartbeatCh:
				// rf.votedFor = none // TODO:正常在StateFollower状态，votedFor== none？？？
				DPrintf("[HEARTBEAT]: raft[%v] receive heartbeat %v.\n", rf.me, heartbeat)

			// 如果超过随机超时时间（随机）未收到心跳信息，则切换状态为候选人
			case <-time.After(time.Duration(rf.randElectionTimeout) * time.Millisecond):
				rf.mu.Lock()
				rf.state = StateCandidate
				rf.mu.Unlock()
				DPrintf("[SWITCHSTATE: ->StateCandidate]: raft[%v] does not receive heartbeat.\n", rf.me)
			}
		case StateCandidate: // TODO: 候选者收到心跳？？？
			rf.mu.Lock()
			rf.votedFor = rf.me
			// 任期加1
			rf.currentTerm++
			rf.mu.Unlock()

			var voteNum int32 = 0

			requestVoteArgs := &RequestVoteArgs{
				Term:         rf.currentTerm,
				CandidateId:  rf.me,
				LastLogIndex: rf.GetLastLogIndex(),
				LastLogTerm:  rf.GetLastLogTerm()}
			// 并行向其他服务器发送请求投票RPCs
			rfLen := len(rf.peers)

			quorumCh := make(chan bool)

			for i := 0; i < rfLen; i++ {
				if i == rf.me { // 不需要向自己发送请求投票RPCs
					atomic.AddInt32(&voteNum, 1) // 给自己投票
					continue
				}
				go func(server int) {
					requestVoteReply := &RequestVoteReply{}
					ok := rf.sendRequestVote(server, requestVoteArgs, requestVoteReply)
					// 如果候选者从发送选举投票时候，发现对方的Term大于自己Term,则设置自己的Term，且切换为StateFollower
					if requestVoteReply.Term > rf.currentTerm {
						DPrintf("[SWITCHSTATE: ->StateCandidate]: raft[%v] receive raft[%v] Term(%v) > self Term(%v).\n", rf.me, server, requestVoteReply.Term, rf.currentTerm)
						rf.mu.Lock()
						rf.currentTerm = requestVoteReply.Term
						rf.state = StateFollower
						rf.mu.Unlock()
					}

					if ok && requestVoteReply.VoteGranted {
						atomic.AddInt32(&voteNum, 1)
						if rf.gatherVotesChangeToLeader(voteNum) {
							// 此处使用select的default来实现，判断chan是否已经塞满数据了
							// 如果之前goroutine已经投票超过半数，则quorumCh<-true被执行过一次（则被塞满了数据，如果不使用default后续将一直阻塞）
							// 使用default后，那么在后续执行quorumCh <- true时候，则会被阻塞，从来走default逻辑，达到不阻塞
							select {
							case quorumCh <- true:
							default:
								// 此处代表之间已经获取投票成功了，后续goroutine不能阻塞在quorumCh <- true，应该直接放过
								return
							}

						}
					}
				}(i)
			}
			// 如果此时不为StateCandidate，则不需要进入下一轮选举操作
			if rf.state != StateCandidate {
				break
			}

			// 第三种可能的结果是候选人既没有赢得选举也没有输：如果有多个跟随者同时成为候选人，
			// 那么选票可能会被瓜分以至于没有候选人可以赢得大多数人的支持。当这种情况发生的时候，
			// 每一个候选人都会超时，然后通过增加当前任期号来开始一轮新的选举。然而，没有其他机制的话，
			// 选票可能会被无限的重复瓜分。

			// 为了阻止选票起初就被瓜分，选举超时时间是从一个固定的区间（例如 150-300毫秒）随机选择。
			// 处理选举超时时间还未获取对应的quorum票数
			select {
			case <-quorumCh:
				DPrintf("raft[%v] get quorum votes success.", rf.me)
			case <-time.After(time.Duration(rf.randElectionTimeout) * time.Millisecond):
				// 超时进入下一轮选举
				DPrintf("raft[%v] get quorum votes timeout, start next election.", rf.me)
				rf.resetRandElectionTimeout()
			}

			// 如果自己当选为leader或者进行下一轮候选人选举，votedFor都应该设置为none
			rf.mu.Lock()
			rf.votedFor = none
			rf.mu.Unlock()
		case StateLeader:
			rfLen := len(rf.peers)
			for i := 0; i < rfLen; i++ {
				if i == rf.me { // 不需要向自己发送心跳消息
					continue
				}

				appendEntriesArgs := AppendEntriesArgs{
					Term:     rf.currentTerm,
					LeaderId: rf.me,
					Entries:  nil}

				go func(server int) {
					appendEntriesReply := &AppendEntriesReply{}
					rf.sendAppendEntries(server, appendEntriesArgs, appendEntriesReply)
					// 如果领导者发送心跳时候，发现对方的Term大于自己Term,则设置自己的Term，且切换为StateFollower
					if appendEntriesReply.Term > rf.currentTerm {
						DPrintf("[SWITCHSTATE: ->StateFollower]: raft[%v] receive raft[%v] Term(%v) > self Term(%v), and switch to StateFollower.\n", rf.me, server, appendEntriesReply.Term, rf.currentTerm)
						rf.mu.Lock()
						rf.currentTerm = appendEntriesReply.Term
						rf.state = StateFollower
						rf.mu.Unlock()
					}
				}(i)
			}
			time.Sleep(time.Duration(rf.heartbeatBroadcastTime) * time.Millisecond)
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
	rf.heartbeatCh = make(chan string)
	rf.heartbeatBroadcastTime = HeartbeatBroadcastTick
	//rf.heartbeatTimeOut = HeartbeatTimeOutTick
	rf.electionTimeout = ElectionTimeoutTick
	rf.resetRandElectionTimeout()

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
