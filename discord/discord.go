package discord

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"path"
	"strconv"
	"time"
)

type DiscrodClient struct {
	token  string
	client http.Client
}

func NewClient(token string) DiscrodClient {
	return DiscrodClient{
		token:  token,
		client: http.Client{},
	}
}

func (c *DiscrodClient) SendMessage(channel_id, message string) error {
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

func (c *DiscrodClient) newMessage(msg []byte, url string) ([]byte, error) {
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(msg))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", c.token)
	req.Header.Set("Content-Type", "application/json")
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
