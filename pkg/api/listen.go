package api

import (
	"conversation-relay/pkg/llms"
	"conversation-relay/pkg/repo"
	"conversation-relay/pkg/trace"
	"conversation-relay/pkg/twilio"
	"conversation-relay/pkg/ws"
	"net/http"

	"github.com/gorilla/mux"
)

type Api struct {
	tracer trace.ITracer
	hub    *ws.Hub
	port   string
	llm    llms.ILLM
	repo   repo.IRepo
}

func NewApi(tracer trace.ITracer, hub *ws.Hub, llm llms.ILLM, port string, repo repo.IRepo) *Api {
	return &Api{
		tracer: tracer,
		hub:    hub,
		port:   port,
		llm:    llm,
		repo:   repo,
	}
}

func (a *Api) Listen() error {
	r := mux.NewRouter()

	r.HandleFunc("/v1.0/{configId}/twiml", func(w http.ResponseWriter, r *http.Request) {
		span := a.tracer.Start("twiml")
		defer span.Finish()
		span.Info("twiml request received")
		vars := mux.Vars(r)
		configId := vars["configId"]
		r.ParseForm()
		accSid := r.FormValue("AccountSid")
		callSid := r.FormValue("CallSid")
		span.Info("twiml::request body", "accSid", accSid, "configId", configId, "callSid", callSid)
		t := twilio.NewTwiml(span)
		twimlStr, err := t.CreateConversationRelayTwiml(accSid, configId, r.Host)
		if err != nil {
			span.Error("twiml::error creating twiml", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		span.Info("twiml::response body", "twiml", twimlStr)
		w.Header().Set("Content-Type", "application/xml")
		w.Write([]byte(twimlStr))
	})

	r.HandleFunc("/v1.0/{configId}/ws", func(w http.ResponseWriter, r *http.Request) {
		span := a.tracer.Start("twilio")
		defer span.Finish()
		span.Info("twilio websocket request received")
		ws.StartTwilioHandler(a.hub, w, r)
	})

	r.HandleFunc("/v1.0/{accSid}/{configId}/queue", func(w http.ResponseWriter, r *http.Request) {
		span := a.tracer.Start("queue")
		defer span.Finish()
		vars := mux.Vars(r)
		configId := vars["configId"]
		accSid := vars["accSid"]
		t := twilio.NewTwiml(span)
		config, _ := a.repo.GetAccountConfig(accSid, configId)
		twiml := t.EnqueueCall(config, r.Host)
		span.Info("QueueTwiml::response body", "twiml", twiml)
		w.Header().Set("Content-Type", "application/xml")
		w.Write([]byte(twiml))
	})

	err := http.ListenAndServe(a.port, r)
	return err
}
