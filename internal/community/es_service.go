package community

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"strconv"
	"strings"

	espkg "github.com/PuppyNote/puppynote-server-golang/pkg/elasticsearch"
	"github.com/elastic/go-elasticsearch/v8/esapi"
)

const esIndex = "community_posts"

type esDocument struct {
	ID       int64    `json:"id"`
	Content  string   `json:"content"`
	Hashtags []string `json:"hashtags"`
}

func indexPost(postID int64, content string, hashtags []string) {
	if espkg.Client == nil {
		return
	}
	doc := esDocument{ID: postID, Content: content, Hashtags: hashtags}
	body, err := json.Marshal(doc)
	if err != nil {
		log.Printf("[ES] 인덱스 직렬화 실패 postID=%d: %v", postID, err)
		return
	}

	req := esapi.IndexRequest{
		Index:      esIndex,
		DocumentID: strconv.FormatInt(postID, 10),
		Body:       bytes.NewReader(body),
		Refresh:    "true",
	}
	res, err := req.Do(context.Background(), espkg.Client)
	if err != nil {
		log.Printf("[ES] 인덱싱 실패 postID=%d: %v", postID, err)
		return
	}
	defer res.Body.Close()
	io.ReadAll(res.Body)
}

func deletePostFromES(postID int64) {
	if espkg.Client == nil {
		return
	}
	req := esapi.DeleteRequest{
		Index:      esIndex,
		DocumentID: strconv.FormatInt(postID, 10),
	}
	res, err := req.Do(context.Background(), espkg.Client)
	if err != nil {
		log.Printf("[ES] 삭제 실패 postID=%d: %v", postID, err)
		return
	}
	defer res.Body.Close()
	io.ReadAll(res.Body)
}

func searchPosts(keyword string, from, size int) ([]int64, int64, error) {
	if espkg.Client == nil {
		return nil, 0, nil
	}

	query := fmt.Sprintf(`{
		"from": %d, "size": %d,
		"query": {
			"multi_match": {
				"query": %q,
				"fields": ["content^2", "hashtags"]
			}
		}
	}`, from, size, keyword)

	res, err := espkg.Client.Search(
		espkg.Client.Search.WithIndex(esIndex),
		espkg.Client.Search.WithBody(strings.NewReader(query)),
	)
	if err != nil {
		return nil, 0, err
	}
	defer res.Body.Close()

	var result struct {
		Hits struct {
			Total struct {
				Value int64 `json:"value"`
			} `json:"total"`
			Hits []struct {
				Source esDocument `json:"_source"`
			} `json:"hits"`
		} `json:"hits"`
	}
	if err = json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, 0, err
	}

	ids := make([]int64, 0, len(result.Hits.Hits))
	for _, h := range result.Hits.Hits {
		ids = append(ids, h.Source.ID)
	}
	return ids, result.Hits.Total.Value, nil
}
