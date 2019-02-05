package types

type ContextCategory struct {
	ID   int64
	Name string
}

type ContextTopic struct {
	ID    int64
	Title string
}

type ContextOpinion struct {
	ID        int64
	TopicID   int64
	CreatorID int64
}
