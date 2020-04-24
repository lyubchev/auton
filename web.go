package main

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/impzero/auton/lib/ibm"
	"github.com/impzero/auton/lib/youtube"
)

// Web is a web struct
type Web struct {
	Router        *chi.Mux
	youtubeClient *youtube.Client
	ibmClient     *ibm.Client
}

// NewWeb is a constructor for Web
func NewWeb(ytClient *youtube.Client, ibmClient *ibm.Client) *Web {
	r := chi.NewRouter()
	w := &Web{Router: r, youtubeClient: ytClient, ibmClient: ibmClient}

	r.Use(render.SetContentType(render.ContentTypeJSON))
	r.Get("/{videoID}", w.analyze)

	return w
}

func (web *Web) analyze(w http.ResponseWriter, r *http.Request) {
	videoID := chi.URLParam(r, "videoID")
	maxComments, err := strconv.Atoi((r.URL.Query().Get("max")))
	if err != nil {

		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, "Max comments must be a number")
	}

	comments, err := web.youtubeClient.GetComments(videoID, youtube.OrderRelevance, maxComments)
	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, "Couldn't fetch comments from youtube video with id"+videoID)
	}

	tones, err := AnalyzeCommentsTone(comments, web.ibmClient)
	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, "Couldn't analyze comments from youtube video with id"+videoID)
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, tones)
}
