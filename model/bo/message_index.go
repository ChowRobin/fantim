package bo

import (
	"bytes"
	"context"
	"encoding/json"
	"time"

	"github.com/ChowRobin/fantim/client"
	"github.com/elastic/go-elasticsearch/v8/esapi"
)

// 消息索引
type MessageIndex struct {
	MsgId          int64  `json:"msg_id"`
	ConversationId string `json:"conversation_id"`
	Content        string `json:"content"`
}

type SearchResult struct {
	ScrollId string `json:"_scroll_id"`
	Hits     struct {
		Hits []struct {
			Source struct {
				MessageIndex
			} `json:"_source"`
		} `json:"hits"`
	} `json:"hits"`
}

func (m *MessageIndex) Create(ctx context.Context) error {
	idxStr, err := json.Marshal(m)
	if err != nil {
		return err
	}

	resp, err := esapi.IndexRequest{
		Index: "fantim",
		Body:  bytes.NewReader(idxStr),
	}.Do(ctx, client.EsClient)
	defer func() {
		if resp != nil && resp.Body != nil {
			_ = resp.Body.Close()
		}
	}()
	if err != nil {
		return err
	}

	return nil
}

func SearchMessageIndex(c context.Context, conversationId, key, cursor string, count int) (result []MessageIndex, newCursor string, err error) {
	var resp *esapi.Response
	defer func() {
		if resp != nil && resp.Body != nil {
			_ = resp.Body.Close()
		}
	}()
	var reqMap map[string]interface{}
	if cursor == "" {
		reqMap = map[string]interface{}{
			"query": map[string]interface{}{
				"bool": map[string]interface{}{
					"filter": []interface{}{
						map[string]interface{}{
							"term": map[string]interface{}{
								"conversation_id": conversationId,
							},
						},
						map[string]interface{}{
							"term": map[string]interface{}{
								"content": key,
							},
						},
					},
				},
			},
			"sort": map[string]interface{}{
				"msg_id": map[string]interface{}{
					"order": "desc",
				},
			},
			"size": count,
		}
		reqStr, _ := json.Marshal(reqMap)
		req := esapi.SearchRequest{
			Body:   bytes.NewReader(reqStr),
			Scroll: time.Minute * 3,
			Pretty: true,
		}
		resp, err = req.Do(c, client.EsClient)

	} else if cursor != "" {
		req := esapi.ScrollRequest{
			ScrollID: cursor,
			Scroll:   time.Minute * 3,
			Pretty:   true,
		}
		resp, err = req.Do(c, client.EsClient)
	}

	if resp == nil {
		return
	}

	respStruct := &SearchResult{}
	buf := new(bytes.Buffer)
	_, _ = buf.ReadFrom(resp.Body)
	_ = json.Unmarshal(buf.Bytes(), respStruct)

	newCursor = respStruct.ScrollId
	for _, hit := range respStruct.Hits.Hits {
		result = append(result, hit.Source.MessageIndex)
	}

	return
}
