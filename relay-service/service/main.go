package main

import (
	"github.com/gin-gonic/gin"
	"github.com/playnb/link"
	"github.com/playnb/shooting-game/pb"
	"github.com/playnb/util/log"
	"time"
)

var serverOpt = &link.ServerOption{
	Addr:            ":1234",
	MaxMsgLen:       0,
	MaxConnNum:      0,
	PendingWriteNum: 0,
	HTTPTimeout:     0,
	CertFile:        "",
	KeyFile:         "",
	RelativePath:    "/relay",
	GinLogger:       log.GinLogger(),
}

type RelayTestAgent struct {
	agent         *link.Agent
	operationList []*pb.T_Operation
	frameTime     uint32
	frameStep     uint32
	frameIndex    uint32
}

func relayTest(agent *link.Agent) {
	rta := &RelayTestAgent{}
	rta.agent = agent
	rta.frameTime = 20
	rta.frameStep = 5
	rta.frameIndex = 0

	log.Trace("echo创建连接: %d", agent.GetUniqueID())
	agent.OnClose = func() {
		log.Trace("echo断开连接: %d", agent.GetUniqueID())
	}

	loop := true
	ticker := time.NewTicker(time.Duration(uint32(time.Millisecond) * rta.frameTime * rta.frameStep))
	defer ticker.Stop()

	for loop {
		select {
		case <-ticker.C:
			rta.frameIndex += rta.frameStep

			cmd := &pb.RelayOperation{}
			cmd.FrameIndex = rta.frameIndex
			cmd.FrameStep = rta.frameStep
			cmd.FrameTime = rta.frameTime
			{
				u := &pb.T_UserOperation{}
				u.UID = 10001
				u.Operation = append(u.Operation, rta.operationList...)
				cmd.Users = append(cmd.Users, u)
			}
			agent.WriteMsg(pb.ParseToBuff(cmd))
			rta.operationList = nil

			log.Trace("Relay")

		case data, ok := <-agent.ReadChan():
			if !ok {
				loop = false
				break
			}
			pb.ParseFromBuff(rta, data)
		}
	}

	log.Trace("relayTest return")
}

func main() {
	log.InitPanic("../tmp")
	log.Init(log.DefaultLogger("../tmp", "run"))
	defer log.Flush()

	log.Trace("relay service")
	gin.SetMode(gin.ReleaseMode)

	{
		pb.HandleIndex.UploadOperation = func(i interface{}, cmd *pb.UploadOperation) {
			rta := i.(*RelayTestAgent)
			if len(cmd.Operation) > 0 {
				//log.Trace("%v", cmd.Operation)
			}
			rta.operationList = append(rta.operationList, cmd.Operation...)
		}
	}

	{
		serv := &link.WSServer{}
		serv.OnAccept = func(agent *link.Agent) {
			go relayTest(agent)
		}
		serv.Start(serverOpt)
		log.Trace("WSServer启动 relayTest")
	}

	for {
		time.Sleep(time.Second)
	}
}
