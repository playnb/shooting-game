package pb

import (
	"encoding/binary"
	"errors"
	"github.com/gogo/protobuf/proto"
	"github.com/playnb/util"
)

type ICanParse interface {
	ParseIndex() uint32

	proto.Marshaler
	proto.Sizer
	MarshalTo(data []byte) (n int, err error)
}

var HandleIndex = &handleIndex{}

func ParseToBuff(p ICanParse) util.BuffData {
	size := p.Size()
	d := util.DefaultPool().Get(size + 4)
	p.MarshalTo(d.GetPayload()[4:])
	binary.BigEndian.PutUint32(d.GetPayload()[:4], p.ParseIndex())
	return d
}

///////////////////////
//Auto Generated Code
///////////////////////

func (*UploadOperation) ParseIndex() uint32 {
	return uint32(Index_Index_UploadOperation)
}

func (*RelayOperation) ParseIndex() uint32 {
	return uint32(Index_Index_RelayOperation)
}

func ParseFromBuff(agent interface{}, data util.BuffData) error {
	index := binary.BigEndian.Uint32(data.GetPayload())
	data.ChangeIndex(4)
	defer data.Release()
	switch Index(index) {
	case Index_Index_UploadOperation:
		cmd := &UploadOperation{}
		err := cmd.Unmarshal(data.GetPayload())
		if err != nil {
			return err
		}
		if HandleIndex.UploadOperation != nil {
			HandleIndex.UploadOperation(agent, cmd)
		}
	case Index_Index_RelayOperation:
		cmd := &RelayOperation{}
		err := cmd.Unmarshal(data.GetPayload())
		if err != nil {
			return err
		}
		if HandleIndex.RelayOperation != nil {
			HandleIndex.RelayOperation(agent, cmd)
		}
	default:
		return errors.New("没有相应的消息")
	}
	return nil
}

type handleIndex struct {
	UploadOperation func(interface{}, *UploadOperation)
	RelayOperation  func(interface{}, *RelayOperation)
}
