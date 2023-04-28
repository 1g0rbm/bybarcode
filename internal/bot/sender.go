package bot

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

type Sender struct {
	token  string
	apiUrl string
	client http.Client
}

func NewSender(token string, apiUrl string) *Sender {
	client := http.Client{
		Timeout: 10 * time.Second,
	}
	return &Sender{
		token:  token,
		apiUrl: apiUrl,
		client: client,
	}
}

func (s Sender) SendMessage(chatId int64, msg string) error {
	actionUrl := url.URL{
		Scheme: "https",
		Host:   s.apiUrl,
		Path:   fmt.Sprintf("/bot%s/sendMessage", s.token),
	}

	val := url.Values{}
	val.Set("chat_id", strconv.FormatInt(chatId, 10))
	val.Set("text", msg)
	body := []byte(val.Encode())

	fmt.Println(actionUrl.String())

	req, err := http.NewRequest(http.MethodPost, actionUrl.String(), bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	response, err := s.client.Do(req.WithContext(ctx))
	if err != nil {
		return err
	}

	defer func(Body io.ReadCloser) {
		err = Body.Close()
	}(response.Body)

	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("send message to bot finished with error")
	}

	return err
}
