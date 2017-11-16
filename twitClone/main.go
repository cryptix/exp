package main

import (
	"flag"
	"fmt"
	"net/url"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/ChimeraCoder/anaconda"
	"github.com/cryptix/go/logging"
	"github.com/jinzhu/configor"
	"github.com/kr/pretty"
	"github.com/pkg/errors"
)

var (
	log   logging.Interface
	check = logging.CheckFatal
)

var cfgFile = flag.String("config", "my.cfg", "which config file to use")

type Config struct {
	Anaconda struct {
		ConsumerKey       string
		ConsumerSecret    string
		AccessToken       string
		AccessTokenSecret string
	}
	Config struct {
		Email          string
		UpdateInterval string
		Debug          bool
		Update         string
		Lists          string
	}
}

func main() {
	logging.SetupLogging(nil)
	log = logging.Logger("twitClone")

	// setup config
	var cfg Config
	check(configor.Load(&cfg, *cfgFile))
	pretty.Println(cfg)

	// setup anaconda
	anaconda.SetConsumerKey(cfg.Anaconda.ConsumerKey)
	anaconda.SetConsumerSecret(cfg.Anaconda.ConsumerSecret)
	api := anaconda.NewTwitterApi(cfg.Anaconda.AccessToken, cfg.Anaconda.AccessTokenSecret)
	//api.SetLogger(anaconda.BasicLogger)

	ok, err := api.VerifyCredentials()
	check(err)
	log.Log("event", "info", "login", ok)

	ud, err := time.ParseDuration(cfg.Config.Update)
	check(err)

	// Mechanical stuff
	//rand.Seed(time.Now().UnixNano())
	//root := context.Background()
	errc := make(chan error)

	go func() {
		errc <- interrupt()
	}()

	for _, l := range strings.Split(cfg.Config.Lists, ",") {
		id, err := strconv.ParseInt(l, 10, 64)
		check(err)
		list, err := api.GetList(id, url.Values{})
		check(err)
		pretty.Println(list)

		go func(l anaconda.List) {
			c := time.Tick(ud)
			for _ = range c {

				tweets, err := api.GetListTweets(list.Id, true, url.Values{})
				check(errors.Wrapf(err, "faild to get tweets for list %s", list.Name))
				log.Log("event", "new tweets", "count", len(tweets))
				if len(tweets) > 0 {
					t := tweets[0]
					log.Log("event", "info", "type", "tweet-onlist", "list", list.Name,
						"user", t.User.ScreenName,
						"userid", t.User.Id,
						"tweetid", t.Id,
						"text", t.Text)
				}
			}
		}(list)
	}

	// Pull: my feed
	//go func() {
	//dur, err := time.ParseDuration(cfg.Config.UpdateInterval)
	//check(err)
	//	t := time.Tick(dur)
	//	for now := range t {
	//		// pulls MY timeline
	//		tweets, err := api.GetUserTimeline(nil)
	//		if err != nil {
	//			errc <- err
	//			return
	//		}
	//		check(err)
	//		for _, tweet := range tweets {
	//			log.WithFields(map[string]interface{}{
	//				"user":    tweet.User.ScreenName,
	//				"userid":  tweet.User.Id,
	//				"tweetid": tweet.Id,
	//			}).Infof("[tweet] %s", tweet.Text)
	//		}
	//		log.WithFields(map[string]interface{}{
	//			"cnt":  len(tweets),
	//			"took": time.Since(now),
	//		}).Info("cycle done")
	//	}
	//}()

	// Stream: my timeline
	//go cloneStream("site", api.SiteStream(nil), errc)
	//go cloneStream("user", api.UserStream(nil), errc)

	check(<-errc)
}

func cloneStream(name string, s *anaconda.Stream, errc chan<- error) {
	for msg := range s.C {
		switch t := msg.(type) {

		case anaconda.FriendsList:
			log.Log("event", "stats", "friends", len(t))

		case anaconda.Tweet:
			log.Log("event", "info", "type", "tweet",
				"user", t.User.ScreenName,
				"userid", t.User.Id,
				"tweetid", t.Id,
				"text", t.Text)

		case anaconda.EventTweet:
			var tweet anaconda.Tweet = *t.TargetObject
			log.Log("event", "info", "type", "eventTweet",
				"user", tweet.User.ScreenName,
				"userid", tweet.User.Id,
				"tweetid", tweet.Id,
				"text", tweet.Text)

		case anaconda.StatusDeletionNotice:
			log.Log("event", "info", "type", "deleted tweet",
				"userid", t.UserId,
				"tweetid", t.Id)

		default:
			log.Log("event", "error", "msg", "uncased msg", "type", fmt.Sprintf("%T %+v", t, msg))
		}
	}
}

func interrupt() error {
	c := make(chan os.Signal)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	return fmt.Errorf("%s", <-c)
}
