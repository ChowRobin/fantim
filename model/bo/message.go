package bo

import (
	"encoding/json"

	"github.com/ChowRobin/fantim/constant"

	"github.com/ChowRobin/fantim/client"
	"github.com/ChowRobin/fantim/model/vo"
	"github.com/go-redis/redis"
)

// 用户链消息
type UserMessage struct {
	vo.MessageBody
}

func (u *UserMessage) Str() string {
	d, _ := json.Marshal(u)
	return string(d)
}

func (u *UserMessage) Add(key string) (int64, error) {
	_, err := client.RedisConn().ZAdd(key, redis.Z{Score: float64(u.MsgId), Member: u.Str()}).Result()
	if err != nil {
		return 0, err
	}
	idx, err := client.RedisConn().ZCard(key).Result()
	idx -= 1
	if err != nil {
		return idx, err
	}
	go func() {
		_, err = client.RedisConn().Expire(key, constant.UserInboxExpiredTime).Result()
	}()
	return idx, nil
}

func PullMessage(key string, start, stop int64) ([]*UserMessage, error) {
	zList, err := client.RedisConn().ZRangeWithScores(key, start, stop).Result()
	if err != nil {
		return nil, err
	}
	list := make([]*UserMessage, 0, len(zList))
	for _, z := range zList {
		u := &UserMessage{}
		err = json.Unmarshal([]byte(z.Member.(string)), u)
		if err != nil {
			return nil, err
		}
		list = append(list, u)
	}
	return list, nil
}

func UserMessageListToVo(msgList []*UserMessage) []*vo.MessageBody {
	result := make([]*vo.MessageBody, 0, len(msgList))
	for _, msg := range msgList {
		result = append(result, &msg.MessageBody)
	}
	return result
}
