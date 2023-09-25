package generators

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vmware-labs/yaml-jsonpath/pkg/yamlpath"
	"gopkg.in/yaml.v3"
)

var y = `---
name: tower
repo: github.com/dashotv/tower
groups:
  downloads:
    path: /downloads
    routes:
      index:
        path: /
      create:
        method: POST
        path: /
      show:
        path: /:id
        params:
          - name: id
            type: string
      update:
        method: PUT
        path: /:id
        params:
          - name: id
            type: string
      setting:
        method: PATCH
        path: /:id
        params:
          - name: id
            type: string
      delete:
        method: DELETE
        path: /:id
        params:
          - name: id
            type: string
      recent:
        path: /recent
      medium:
        path: /:id/medium
        params:
          - name: id
            type: string
      select:
        method: PUT
        path: /:id/select
        params:
          - name: id
            type: string
  upcoming:
    path: /upcoming
    routes:
      index:
        path: /
  episodes:
    path: /episodes
    routes:
      update:
        path: /:id
        method: PUT
        params:
          - name: id
            type: string
      setting:
        path: /:id
        method: PATCH
        params:
          - name: id
            type: string
  series:
    path: /series
    routes:
      index:
        path: /
      create:
        method: POST
        path: /
      show:
        path: /:id
        params:
          - name: id
            type: string
      update:
        method: PUT
        path: /:id
        params:
          - name: id
            type: string
      setting:
        method: PATCH
        path: /:id
        params:
          - name: id
            type: string
      delete:
        method: DELETE
        path: /:id
        params:
          - name: id
            type: string
      seasons:
        path: /:id/seasons
        params:
          - name: id
            type: string
      seasonEpisodesAll:
        path: /:id/seasons/all
        params:
          - name: id
            type: string
      seasonEpisodes:
        path: /:id/seasons/:season
        params:
          - name: id
            type: string
          - name: season
            type: string
      currentSeason:
        path: /:id/currentseason
        params:
          - name: id
            type: string
      paths:
        path: /:id/paths
        params:
          - name: id
            type: string
      watches:
        path: /:id/watches
        params:
          - name: id
            type: string
  movies:
    path: /movies
    routes:
      index:
        path: /
      create:
        method: POST
        path: /
      show:
        path: /:id
        params:
          - name: id
            type: string
      update:
        method: PUT
        path: /:id
        params:
          - name: id
            type: string
      setting:
        method: PATCH
        path: /:id
        params:
          - name: id
            type: string
      delete:
        method: DELETE
        path: /:id
        params:
          - name: id
            type: string
      paths:
        path: /:id/paths
        params:
          - name: id
            type: string
  releases:
    path: /releases
    routes:
      index:
        path: /
      create:
        method: POST
        path: /
      show:
        path: /:id
        params:
          - name: id
            type: string
      update:
        method: PUT
        path: /:id
        params:
          - name: id
            type: string
      setting:
        method: PATCH
        path: /:id
        params:
          - name: id
            type: string
      delete:
        method: DELETE
        path: /:id
        params:
          - name: id
            type: string
  feeds:
    path: /feeds
    rest: true
`

func TestFindRoutesYAML(t *testing.T) {
	path := "/downloads/:id/select"
	parts := strings.Split(path, "/")

	parent := parts[1]
	child := "/" + strings.Join(parts[2:], "/")
	fmt.Printf("parent: %s, child: %s\n", parent, child)

	r := &yaml.Node{}
	err := yaml.Unmarshal([]byte(y), r)
	if err != nil {
		assert.NoError(t, err, "error unmarshalling yaml")
	}

	p, err := yamlpath.NewPath("$..groups.." + parent + "..routes")
	if err != nil {
		assert.NoError(t, err, "error creating yaml path")
	}

	q, err := p.Find(r)
	if err != nil {
		assert.NoError(t, err, "error finding yaml path")
	}

	for _, v := range q {
		// fmt.Printf("found: %#v\n", v)
		fp, err := yamlpath.NewPath("$..path")
		if err != nil {
			assert.NoError(t, err, "error creating yaml path")
		}
		fq, err := fp.Find(v)
		if err != nil {
			assert.NoError(t, err, "error finding yaml path")
		}
		for _, c := range fq {
			if c.Value == child {
				fmt.Printf("parent: %#v\n", v)
				fmt.Printf("child: %#v\n", c)
			}
		}
	}
}
