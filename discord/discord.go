package discord

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"net/url"
	"path"
	"strconv"
	"strings"
	"time"

	ua "github.com/eddycjy/fake-useragent"
)

type DiscordClient struct {
	index    int
	token    string
	useProxy bool
	client   http.Client
}

func NewClient(index int, token string, Proxy bool) DiscordClient {
	return DiscordClient{
		index:    index,
		token:    token,
		useProxy: Proxy,
		client:   http.Client{},
	}
}

func (c *DiscordClient) SetProxy(proxyList []string) error {
	if c.useProxy {
		parts := strings.Split(proxyList[c.index], "@")
		logPass := strings.Split(parts[0], ":")
		log.Println(parts)
		log.Println(logPass)

		proxy := &url.URL{
			Scheme: "http",
			User:   url.UserPassword(logPass[0], logPass[1]),
			Host:   parts[1],
		}
		log.Printf("Acc.%v | Setup proxy: %v", c.index+1, proxy)

		c.client.Transport = &http.Transport{
			Proxy: http.ProxyURL(proxy),
		}
	}
	return nil
}

func (c *DiscordClient) SendMessage(channel_id, message string) error {
	url := "https://" + path.Join("discord.com/api/v9/channels/", channel_id, "/messages")

	msg := Message{
		Content:             message,
		Flags:               0,
		Mobile_network_type: "unknown",
		Nonce:               strconv.FormatInt(time.Now().UnixNano()/1e6, 10),
		Tts:                 false,
	}

	b, _ := json.Marshal(msg)

	resp, err := c.newMessage(b, url)
	if err != nil {
		return err
	}
	var response Response
	if err := json.Unmarshal(resp, &response); err != nil {
		return err
	}

	return nil
}

func (c *DiscordClient) newMessage(msg []byte, url string) ([]byte, error) {
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(msg))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", c.token)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", ua.Computer())
	resp, _ := c.client.Do(req)
	if resp.StatusCode != 200 {
		return nil, errors.New("bad request")
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return body, nil
}
