package lib

import (
	"fmt"
	"time"

	"github.com/wsxiaoys/terminal/color"
)

var layout string = "2006-01-02T15:04:05"

type DoneAt struct {
	buildedAt     *time.Time
	startedAt     *time.Time
	reloadedAt    *time.Time
	postrotatedAt *time.Time
}

func (it *DoneAt) BuildedAt() *time.Time     { return it.buildedAt }
func (it *DoneAt) StartedAt() *time.Time     { return it.startedAt }
func (it *DoneAt) ReloadedAt() *time.Time    { return it.reloadedAt }
func (it *DoneAt) PostrotatedAt() *time.Time { return it.postrotatedAt }

func (it *DoneAt) UpdateReloadedAt() *DoneAt {
	reloadedAt := time.Now()
	it.reloadedAt = &reloadedAt
	return it
}

func (it *DoneAt) UpdatePostrotatedAt() *DoneAt {
	postrotatedAt := time.Now()
	it.postrotatedAt = &postrotatedAt
	return it
}

func (it *DoneAt) Print() *DoneAt {
	fmt.Printf("builded-at: %s\n", color.Sprintf("@g%s", it.BuildedAt()))
	fmt.Printf("started-at: %s\n", color.Sprintf("@g%s", it.StartedAt()))
	if it.ReloadedAt() != nil {
		fmt.Printf("reloaded-at: %s\n", color.Sprintf("@g%s", it.ReloadedAt()))
	}
	if it.PostrotatedAt() != nil {
		fmt.Printf("postrotated-at: %s\n", color.Sprintf("@g%s", it.PostrotatedAt()))
	}
	return it
}

func NewDefaultDoneAt() *DoneAt {
	it, _ := NewDoneAt("")
	return it
}

func NewDoneAt(buildDate string) (*DoneAt, error) {
	it := new(DoneAt)

	if buildDate == "" {
		buildAt := time.Now()
		it.buildedAt = &buildAt
	} else {
		if buildedAt, e := time.Parse(layout, buildDate); e != nil {
			return nil, fmt.Errorf("time.Parse(...)(layout='%s', buildDate='%s') failed: %s", layout, buildDate, e.Error())
		} else {
			it.buildedAt = &buildedAt
		}
	}

	startedAt := time.Now()
	it.startedAt = &startedAt

	return it, nil
}
