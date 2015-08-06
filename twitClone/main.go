package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/ChimeraCoder/anaconda"
	"github.com/baobabus/gcfg"
	"github.com/cryptix/go/logging"
)

var log = logging.Logger("twitClone")

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
	}
}

func main() {
	logging.SetupLogging(nil)

	// setup config
	var cfg Config
	logging.CheckFatal(gcfg.ReadFileInto(&cfg, *cfgFile))

	// setup anaconda
	anaconda.SetConsumerKey(cfg.Anaconda.ConsumerKey)
	anaconda.SetConsumerSecret(cfg.Anaconda.ConsumerSecret)
	api := anaconda.NewTwitterApi(cfg.Anaconda.AccessToken, cfg.Anaconda.AccessTokenSecret)
	//api.SetLogger(anaconda.BasicLogger)

	ok, err := api.VerifyCredentials()
	logging.CheckFatal(err)
	log.Info("Credentials:", ok)

	// Mechanical stuff
	//rand.Seed(time.Now().UnixNano())
	//root := context.Background()
	errc := make(chan error)

	go func() {
		errc <- interrupt()
	}()

	// Pull: my feed
	//go func() {
	//dur, err := time.ParseDuration(cfg.Config.UpdateInterval)
	//logging.CheckFatal(err)
	//	t := time.Tick(dur)
	//	for now := range t {
	//		// pulls MY timeline
	//		tweets, err := api.GetUserTimeline(nil)
	//		if err != nil {
	//			errc <- err
	//			return
	//		}
	//		logging.CheckFatal(err)
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
	go cloneStream("user", api.UserStream(nil), errc)

	log.Fatal("global error:", <-errc)
}

func cloneStream(name string, s anaconda.Stream, errc chan<- error) {
	go func() {
		v := <-s.Quit
		errc <- fmt.Errorf("stream[%s] quit: %v", name, v)
	}()
	for msg := range s.C {
		switch t := msg.(type) {

		case anaconda.FriendsList:
			log.WithField("count", len(t)).Infof("friends[%s]", name)

		case anaconda.Tweet:
			log.WithFields(map[string]interface{}{
				"user":    t.User.ScreenName,
				"userid":  t.User.Id,
				"tweetid": t.Id,
			}).Infof("tweet[%s] %s", name, t.Text)

		case anaconda.EventTweet:
			var tweet anaconda.Tweet = *t.TargetObject
			log.WithFields(map[string]interface{}{
				"user":    tweet.User.ScreenName,
				"userid":  tweet.User.Id,
				"tweetid": tweet.Id,
			}).Infof("Eventtweet[%s] %s", name, tweet.Text)

		case anaconda.StatusDeletionNotice:
			log.WithFields(map[string]interface{}{
				"userid":  t.UserId,
				"tweetid": t.Id,
			}).Warningf("tweet[%s] deleted", name)

		default:
			log.Error("uncased msg %T %+v", t, msg)
		}
	}
}

func interrupt() error {
	c := make(chan os.Signal)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	return fmt.Errorf("%s", <-c)
}
