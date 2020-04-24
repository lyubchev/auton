package main

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
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

	r.Use(
		middleware.RequestID,
		middleware.RedirectSlashes,
		middleware.RealIP,
		middleware.Logger,
		middleware.Recoverer,
		Security,
		cors.New(cors.Options{
			AllowedOrigins:   []string{"*"},
			AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
			AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
			ExposedHeaders:   []string{"Link"},
			AllowCredentials: true,
		}).Handler,
		middleware.Heartbeat("/ping"),
	)
	r.Use(render.SetContentType(render.ContentTypeJSON))
	r.Get("/{videoID}", w.analyze)

	return w
}

func (web *Web) analyze(w http.ResponseWriter, r *http.Request) {
	videoID := chi.URLParam(r, "videoID")
	maxComments, err := strconv.Atoi((r.URL.Query().Get("max")))
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, "Max comments must be a number of type int")
		return
	}

	comments, err := web.youtubeClient.GetComments(videoID, youtube.OrderRelevance, maxComments)
	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, "Couldn't fetch comments from youtube video with id"+videoID)
		return
	}

	tones, err := AnalyzeCommentsTone(comments, web.ibmClient)
	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, "Couldn't analyze comments from youtube video with id"+videoID)
		return
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, tones)
}

const (
	xFrameOptions                = "X-Frame-Options"
	xFrameOptionsValue           = "DENY"
	xContentTypeOptions          = "X-Content-Type-Options"
	xContentTypeOptionsValue     = "nosniff"
	xssProtection                = "X-XSS-Protection"
	xssProtectionValue           = "1; mode=block"
	strictTransportSecurity      = "Strict-Transport-Security"                    // details https://blog.bracelab.com/achieving-perfect-ssl-labs-score-with-go + https://developer.mozilla.org/en-US/docs/Web/Security/HTTP_strict_transport_security
	strictTransportSecurityValue = "max-age=31536000; includeSubDomains; preload" // 31536000 = just shy of 12 months
)

func Security(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add(xFrameOptions, xFrameOptionsValue)
		w.Header().Add(xContentTypeOptions, xContentTypeOptionsValue)
		w.Header().Add(xssProtection, xssProtectionValue)
		w.Header().Add(strictTransportSecurity, strictTransportSecurityValue)

		next.ServeHTTP(w, r)
	})
}
