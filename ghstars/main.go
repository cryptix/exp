package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"sort"
	"time"

	"github.com/cryptix/go/logging"
	"github.com/dustin/go-humanize"
	"github.com/jinzhu/configor"
	"github.com/kr/pretty"
	"github.com/shurcooL/githubql"
	"golang.org/x/oauth2"
)

var (
	log   logging.Interface
	check = logging.CheckFatal
)

var cfgFile = flag.String("config", "my.cfg", "which config file to use")

type Config struct {
	Github struct {
		AccessToken string
	}
	Config struct {
		Email          string
		UpdateInterval string
		Debug          bool
		Update         string
	}
}

func main() {
	logging.SetupLogging(os.Stdout)
	log = logging.Logger(os.Args[0])

	// setup config
	var cfg Config
	check(configor.Load(&cfg, *cfgFile))
	pretty.Println(cfg)

	src := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: cfg.Github.AccessToken},
	)
	httpClient := oauth2.NewClient(context.Background(), src)

	client := githubql.NewClient(httpClient)

	var query struct {
		Viewer struct {
			StarredRepositories struct {
				PageInfo struct {
					EndCursor   githubql.String
					HasNextPage githubql.Boolean
				}
				TotalCount githubql.Int
				Nodes      []repo
			} `graphql:"starredRepositories(first: 100, after: $starCursor)"`
		}
	}

	variables := map[string]interface{}{
		"starCursor": (*githubql.String)(nil),
	}

	var allStars []repo
	var i = 0

	for {
		start := time.Now()
		err := client.Query(context.Background(), &query, variables)
		check(err)

		allStars = append(allStars, query.Viewer.StarredRepositories.Nodes...)
		if !query.Viewer.StarredRepositories.PageInfo.HasNextPage {
			break
		}
		log.Log("page", i, "took", fmt.Sprintf("%v", time.Since(start)))
		i++

		variables["starCursor"] = githubql.NewString(query.Viewer.StarredRepositories.PageInfo.EndCursor)
	}

	kbytes := 0
	sort.Sort(BySize(allStars))
	for i, r := range allStars {
		log.Log("star", i,
			"name", r.NameWithOwner,
			"size", humanize.Bytes(uint64(r.DiskUsage)*1024))
		kbytes += int(r.DiskUsage)
	}
	log.Log("total", len(allStars), "size", humanize.Bytes(uint64(kbytes*1024)))

}

type repo struct {
	Id            githubql.ID
	NameWithOwner githubql.String
	Description   githubql.String
	Url           githubql.URI

	Ref struct {
		Id     githubql.ID
		Target struct {
			Oid githubql.GitObjectID
		}
	} `graphql:"ref(qualifiedName:\"master\")"`
	DiskUsage   githubql.Int
	LicenseInfo struct {
		Id     githubql.ID
		Name   githubql.String
		SpdxId githubql.String
	}
}

type BySize []repo

func (a BySize) Len() int           { return len(a) }
func (a BySize) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a BySize) Less(i, j int) bool { return a[i].DiskUsage < a[j].DiskUsage }
