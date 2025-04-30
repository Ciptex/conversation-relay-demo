package main

import (
	"context"
	"conversation-relay/pkg/api"
	"conversation-relay/pkg/llms"
	"conversation-relay/pkg/mq"
	"conversation-relay/pkg/mqsubs"
	"conversation-relay/pkg/repo"
	"conversation-relay/pkg/trace"
	"conversation-relay/pkg/ws"
	"log"

	"github.com/joho/godotenv"
)

func main() {
	tracer, _ := trace.CreateGlobalTracer(trace.Console, "conversation-relay", trace.DEV)
	span := tracer.Start("main")
	defer span.Finish()

	godotenv.Load()
	repo := repo.NewDB()

	hub := ws.NewHub(repo)
	go hub.Listen()
	llm := llms.CreateLLMModel()

	pub := mq.NewPublisher(hub.VoiceResponse, span)
	go pub.Listen()
	hub.SetMQPublisher(pub.Publisher)

	ctx := context.TODO()
	mqsubs.CreateSubscribers(pub, repo, llm, ctx, span)

	span.Info("start conversation relay handler...", "port", "8080")
	api := api.NewApi(tracer, hub, llm, ":8080", repo)
	err := api.Listen()
	if err != nil {
		span.Error("api listent error", err)
		log.Fatal("API error: ", err)
	}
}
