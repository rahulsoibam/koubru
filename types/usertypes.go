package types

import (
	"errors"

	"github.com/rahulsoibam/koubru/utils"
)

type NewTopic struct {
	Title      string   `json:"title"`
	Details    string   `json:"details"`
	Categories [3]int64 `json:"categories"`
}

func (t *NewTopic) Validate() error {
	if len(t.Title) < 8 || len(t.Title) > 80 {
		return errors.New("Should be less than 80 and more than 8 characters")
	}
	if len(t.Details) > 1024 {
		return errors.New("Details should be less than 1024 characters")
	}
	return nil
}

type NewCategory struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

func (c *NewCategory) Validate() error {
	// TODO
	if len(c.Name) < 3 || len(c.Name) > 32 {
		return errors.New("Should be less than 32 and more than 3 characters")
	}
	return nil
}

type NewOpinion struct {
	TopicID     int64  `json:"topic_id"`
	Reaction    string `json:"reaction"`
	IsAnonymous bool   `json:"is_anonymous"`
	Mp4         string `json:"mp4"`
}

func (no *NewOpinion) Validate() error {
	if no.TopicID == 0 {
		return errors.New("Topic ID is required.")
	}
	if !utils.IsValidOpinion(no.Reaction) {
		return errors.New("Not a valid reaction.")
	}
	if no.Mp4 == "" {
		return errors.New("Problem with upload.")
	}
	return nil
}

type NewReply struct {
	TopicID     int64
	Reaction    string
	IsAnonymous bool
	Mp4         string
	ParentID    int64
}

func (nr *NewReply) Validate() error {
	if nr.TopicID == 0 {
		return errors.New("Topic ID is required.")
	}
	if !utils.IsValidOpinion(nr.Reaction) {
		return errors.New("Not a valid opinion")
	}
	if nr.Mp4 == "" {
		return errors.New("Problem with upload.")
	}
	return nil
}
