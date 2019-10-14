package handler

import (
	"context"
	"encoding/json"
	"github.com/go-chi/chi"
	shortlinkService "github.com/jasondeutsch/shorty/internal/link/service"
	statsService "github.com/jasondeutsch/shorty/internal/stats/service"
	"gopkg.in/matryer/respond.v1"
	"net/http"
	"time"
)

type Handler struct {
	router      chi.Router
	linkService shortlinkService.Service
	statsService statsService.Service
}

func NewHandler(linkService shortlinkService.Service, statsService statsService.Service) *Handler {
	h := &Handler{
		router: chi.NewMux(),
		linkService: linkService,
		statsService: statsService,
	}

	h.router.With(ShortLinkCtx).Get("/{slug}", h.RedirectFromShortLink)
	h.router.Post("/links", h.CreateShortLink)
	h.router.With(ShortLinkCtx).Get("/links/{slug}/stats/", h.GetShortLinkStats)

	return h
}

func (h *Handler) RedirectFromShortLink(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")

	shortlink, err := h.linkService.ReadOne(slug)
	if err != nil {
		respond.With(w,r, http.StatusBadRequest, err)
		return
	}

	if shortlink == nil {
		respond.WithStatus(w, r, http.StatusNotFound)
		return
	}

	http.Redirect(w, r, shortlink.Destination, http.StatusSeeOther)

	view := h.statsService.NewView(shortlink.Slug, time.Now())
	err = h.statsService.SaveView(view)
}


func (h *Handler) CreateShortLink(w http.ResponseWriter, r *http.Request) {
	type createRequest struct {
		Destination string `json:"destination"`
	}
	var req createRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		respond.With(w, r, http.StatusBadRequest, err)
		return
	}

	link, err := h.linkService.Create(req.Destination)
	if err != nil {
		respond.With(w,r, http.StatusInternalServerError, err)
		return
	}

	respond.With(w,r, http.StatusCreated, link)
}

func (h * Handler) GetShortLinkStats(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")

	snapshot, err := h.statsService.GetSnapshot(slug)
	// TODO: handle no result
	if err != nil {
		respond.With(w,r, http.StatusInternalServerError, err)
		return
	}

	resp := snapShotResposne{
		Slug: snapshot.Slug,
		Count24Hours: snapshot.Count24Hours,
		CountOneWeek: snapshot.CountOneWeek,
		CountAllTime: snapshot.CountAllTime,
	}

	respond.With(w,r, http.StatusOK, resp)
}

type snapShotResposne struct {
	Slug string  `json:"slug"`
	Count24Hours int  `json:"count_24_hours"`
	CountOneWeek int  `json:"count_one_week"`
	CountAllTime  int`json:"count_all_time"`
}

func Cow(w http.ResponseWriter, r *http.Request) {
	cow := `
 _____
< 404 >
 -----
        \   ^__^
         \  (oo)\_______
            (__)\       )\/\
                ||----w |
                ||     ||
`

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte(cow))
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.router.ServeHTTP(w, r)
}

func ShortLinkCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		slug := chi.URLParam(r, "slug")
		if slug == "" {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		ctx := context.WithValue(r.Context(), "slug", slug)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

