package mqsubs

import (
	"context"
	"conversation-relay/pkg/llms"
	"conversation-relay/pkg/mq"
	"conversation-relay/pkg/mqsubs/handlers"
	"conversation-relay/pkg/repo"
	"conversation-relay/pkg/trace"
	"conversation-relay/pkg/types"
)

func CreateSubscribers(pub *mq.MQPub, repo repo.IRepo, llm llms.ILLM, ctx context.Context, span trace.ISpan) {
	if loggerSub := mq.NewSubscriber([]string{types.Logger}, pub.Publisher, span); loggerSub != nil {
		go loggerSub.Listen(ctx, handlers.NewLoggerHandler())
	}

	if genericSub := mq.NewSubscriber([]string{types.GenericHandler}, pub.Publisher, span); genericSub != nil {
		go genericSub.Listen(ctx, handlers.NewGenericHandler(repo, llm))
	}

	if greetSub := mq.NewSubscriber([]string{types.GreetHandler}, pub.Publisher, span); greetSub != nil {
		go greetSub.Listen(ctx, handlers.NewGenericHandler(repo, llm))
	}

	if intentSub := mq.NewSubscriber([]string{types.IntentHandler}, pub.Publisher, span); intentSub != nil {
		go intentSub.Listen(ctx, handlers.NewIntentHandler(repo, llm))
	}
}
