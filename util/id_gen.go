package util

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/ChowRobin/fantim/constant"
)

/*
组成：0(1 bit) | timestamp in milli second (41 bit) | machine id (10 bit) | index (12 bit)
每毫秒最多生成4096个id，集群机器最多1024台
*/

var (
	sf    *Snowflake
	mutex *sync.Mutex
)

func init() {
	mutex = &sync.Mutex{}
	sf = &Snowflake{}
	// 分布式场景下要用唯一机器id init
	err := sf.Init(1)
	if err != nil {
		panic(err)
	}
}

type Snowflake struct {
	lastTimestamp int64
	index         int16
	machId        int16
}

func (s *Snowflake) Init(id int16) error {
	if id > 0xfff {
		return errors.New("illegal machine id")
	}

	s.machId = id
	s.lastTimestamp = time.Now().UnixNano() / 1e6
	s.index = 0
	return nil
}

func (s *Snowflake) GetId() (int64, error) {
	mutex.Lock()
	defer mutex.Unlock()
	curTimestamp := time.Now().UnixNano() / 1e6
	if curTimestamp == s.lastTimestamp {
		s.index++
		if s.index > 0xfff {
			s.index = 0xfff
			return -1, errors.New("out of range")
		}
	} else {
		s.index = 0
		s.lastTimestamp = curTimestamp
	}
	return (0x1ffffffffff&s.lastTimestamp)<<22 + int64(0xff<<10) + int64(0xfff&s.index), nil
}

func GenId() int64 {
	id, _ := sf.GetId()
	return id
}

func GenConversationId(conversationType int32, sender, receiver int64) string {
	switch conversationType {
	case constant.ConversationTypeSingle:
		if sender < receiver {
			return fmt.Sprintf(constant.ConversationIdPatternSingle, sender, receiver)
		} else {
			return fmt.Sprintf(constant.ConversationIdPatternSingle, receiver, sender)
		}
	case constant.ConversationTypeGroup:
		return fmt.Sprintf(constant.ConversationIdPatternGroup, sender, receiver)
	}
}
