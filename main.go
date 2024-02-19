package main

import (
	"flag"
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/SevereCloud/vksdk/v2/api"
	"github.com/SevereCloud/vksdk/v2/object"
	"github.com/go-co-op/gocron"
	log "github.com/inconshreveable/log15"
	"strings"
	"time"
)

type config struct {
	ServiceToken string
	MessageToken string

	Admins []uint
	Groups []int

	Period string

	Message string
}

var (
	conf      config
	countList []int
	serviceVK *api.VK
	botVK     *api.VK
)

var logger = log.New("service", "vkonity")

func main() {
	configFile := flag.String("config", "config.toml", "path to the config file")
	flag.Parse()

	var err error
	if _, err = toml.DecodeFile(*configFile, &conf); err != nil {
		panic(err)
	}

	serviceVK = api.NewVK(conf.ServiceToken)
	botVK = api.NewVK(conf.MessageToken)

	s := gocron.NewScheduler(time.UTC)
	_, err = s.Every(conf.Period).Do(check)
	if err != nil {
		panic(err)
	}
	s.StartAsync()
	s.StartBlocking()
}

func check() {
	for i, groupId := range conf.Groups {
		group := get(groupId)
		count := group.Count

		// If not current group in list
		if len(countList) < i+1 {
			countList = append(countList, count)
			continue
		}

		previewCount := countList[i]
		countList[i] = count
		if count <= previewCount {
			continue
		}

		post := getPost(group.Items)
		send(generateMessageText(post), getAttachment(post))
	}
}

func getAttachment(post object.WallWallpost) string {
	attachments := post.Attachments
	if len(attachments) < 1 {
		return ""
	}

	var attachmentsParsed [10]string
	for i, attachment := range attachments {
		attachmentsParsed[i] = attachment.Photo.ToAttachment()
	}

	return strings.Join(attachmentsParsed[:], ",")
}

func generateMessageText(post object.WallWallpost) string {
	groupId := post.OwnerID
	postId := post.ID
	text := post.Text

	return fmt.Sprintf(conf.Message, -groupId, groupId, postId, text)
}

func getPost(items []object.WallWallpost) (post object.WallWallpost) {
	post = items[0]

	if post.IsPinned {
		post = items[1]
	}

	return post
}

func send(text string, attachment string) {
	_, err := botVK.MessagesSendPeerIDs(api.Params{
		"user_ids":   conf.Admins,
		"message":    text,
		"random_id":  0,
		"attachment": attachment,
	})
	if err != nil {
		logger.Error("Failed to send message", "admins", conf.Admins, "text", text, "attachment", attachment)
	}
}

func get(groupId int) api.WallGetResponse {
	group, err := serviceVK.WallGet(api.Params{
		"count":    2,
		"owner_id": -groupId,
	})
	if err != nil {
		logger.Error("Failed to get posts", "owner_id", -groupId)
	}

	return group
}
