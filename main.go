package main

import (
	"flag"
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/SevereCloud/vksdk/v2/api"
	"github.com/SevereCloud/vksdk/v2/object"
	"github.com/go-co-op/gocron"
	"log"
	"strings"
	"time"
)

type config struct {
	ServiceToken string
	MessageToken string

	Admins []int
	Groups []int

	Period string

	Message string
}

var (
	conf      config
	counts    [10]int
	serviceVK *api.VK
	botVK     *api.VK
)

func main() {
	configFile := flag.String("config", "config.toml", "path to the config file")
	flag.Parse()

	var err error
	if _, err = toml.DecodeFile(*configFile, &conf); err != nil {
		log.Fatalln(err)
	}

	serviceVK = api.NewVK(conf.ServiceToken)
	botVK = api.NewVK(conf.MessageToken)

	s := gocron.NewScheduler(time.UTC)
	_, err = s.Every(conf.Period).Do(check)
	if err != nil {
		log.Fatal(err)
	}
	s.StartAsync()
	s.StartBlocking()
}

func check() {
	for i, groupId := range conf.Groups {
		group := get(groupId)
		count := group.Count
		previewCount := counts[i]

		counts[i] = count
		if previewCount == 0 || count <= previewCount {
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
		log.Fatal(err)
	}
}

func get(groupId int) api.WallGetResponse {
	group, err := serviceVK.WallGet(api.Params{
		"count":    2,
		"owner_id": -groupId,
	})
	if err != nil {
		log.Fatal(err)
	}

	return group
}
