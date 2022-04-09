package main

import (
	"context"
	"encoding/json"
	"fmt"
	"hash/fnv"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sort"
	"strings"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

const (
	hellotalk = "https://interview.hellotalk8.com"
)

var ch = make(chan struct{})

type Params struct {
	Key   string `json:"key"`
	Token string `json:"token"`
}

type Resp struct {
	Info string `json:"info"`
}

type RateLimitedClient struct {
	client      *http.Client
	RateLimiter *rate.Limiter
}

func NewClient(rl *rate.Limiter) *RateLimitedClient {
	c := &RateLimitedClient{
		client:      http.DefaultClient,
		RateLimiter: rl,
	}
	return c
}

func (c *RateLimitedClient) Do(req *http.Request) (*http.Response, error) {
	ctx := context.Background()
	err := c.RateLimiter.Wait(ctx)
	if err != nil {
		return nil, err
	}
	resp, err := c.client.Do(req)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	return resp, nil
}

func NewPost(key, token string) *http.Request {
	req, _ := http.NewRequest(http.MethodPost, hellotalk, nil)
	q := req.URL.Query()
	q.Add("key", key)
	q.Add("token", token)
	req.URL.RawQuery = q.Encode()

	return req
}

func GetKeyAndIP(c *http.Client) (string, string, string) {
	resp, err := c.Get(hellotalk)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	params := Params{}
	err = json.Unmarshal(body, &params)
	if err != nil {
		log.Printf("Reading body failed: %s", err)
		return "", "", ""
	}

	return params.Key, params.Token, resp.Header.Get("X-Rate-Limit-Request-Forwarded-For")
}

func start(c *RateLimitedClient, logger *log.Logger) {
	var wg sync.WaitGroup
	key, token, _ := GetKeyAndIP(c.client)
	fmt.Println(key, token)
	res := Permutation(key)
	lr := len(res)
	var reqs []*http.Request
	for _, v := range res {
		reqs = append(reqs, NewPost(v, token))
	}

	for i := 0; i < lr; i++ {
		go func(i int) {
			wg.Add(1)
			resp, err := c.Do(reqs[i])
			logger.Println(resp.Request.URL, resp.StatusCode)
			if err != nil {
				logger.Fatalln(err.Error())
			}
			defer resp.Body.Close()
			body, err := ioutil.ReadAll(resp.Body)
			response := Resp{}
			_ = json.Unmarshal(body, response)
			if (resp.StatusCode == 410 || resp.StatusCode == 400) && tokenExpired(response.Info) { // token 过期
				ch <- struct{}{}
			}
			wg.Done()
			if resp.StatusCode != 200 {
				return
			}
			logger.Println(body)
			ch <- struct{}{}
		}(i)
	}

	go func() {
		wg.Wait()
		// ch <- struct{}{}
	}()

	<-ch
}

func main() {
	f, err := os.OpenFile("hellotalk.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println(err)
	}
	defer f.Close()
	logger := log.New(f, "", log.LstdFlags)

	rl := rate.NewLimiter(rate.Every(1*time.Second), 500)
	c := NewClient(rl)
	start(c, logger)

	for {
		select {
		case <-ch:
			start(c, logger)
		}
	}
}

func Permutation(s string) (ans []string) {
	t := []byte(s)
	sort.Slice(t, func(i, j int) bool {
		return t[i] < t[j]
	})
	n := len(t)
	perm := make([]byte, 0, n)
	vis := make([]bool, n)
	var backtrack func(int)
	backtrack = func(i int) {
		if i == n {
			ans = append(ans, string(perm))
			return
		}
		for j, b := range vis {
			if b {
				continue
			}
			if j > 0 && !vis[j-1] && t[j-1] == t[j] {
				continue
			}
			vis[j] = true
			perm = append(perm, t[j])
			backtrack(i + 1)
			perm = perm[:len(perm)-1]
			vis[j] = false
		}
	}
	backtrack(0)
	return
}

func hash(s string) uint32 {
	h := fnv.New32a()
	h.Write([]byte(s))
	return h.Sum32()
}

func tokenExpired(s string) bool {
	return strings.Contains(s, "token")
}
