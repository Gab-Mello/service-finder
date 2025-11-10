package posting

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	authmw "github.com/Gab-Mello/service-finder/internal/http/middleware/auth"
	domain "github.com/Gab-Mello/service-finder/internal/posting"
)

type Handler struct{ svc *domain.Service }

func NewHandler(s *domain.Service) *Handler { return &Handler{svc: s} }

// Create
// @Summary  Criar anúncio (posting)
// @Tags     postings
// @Router   /postings [post]
func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	pid, ok := authmw.UserIDFromContext(r)
	if !ok {
		writeErr(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	var req CreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeErr(w, 400, "invalid json")
		return
	}

	p, err := h.svc.Create(pid, req.Title, req.Description, req.Price, req.Category, req.City, req.District)
	if err != nil {
		writeErr(w, 400, err.Error())
		return
	}
	writeJSON(w, 201, p)
}

// Update
// @Summary  Atualizar anúncio
// @Tags     postings
// @Param    id      path   string  true  "Posting ID"
// @Router   /postings/{id} [patch]
func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	pid, ok := authmw.UserIDFromContext(r)
	if !ok {
		writeErr(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	id := strings.TrimPrefix(r.URL.Path, "/api/v1/postings/")
	var patch map[string]any
	if err := json.NewDecoder(r.Body).Decode(&patch); err != nil {
		writeErr(w, 400, "invalid json")
		return
	}
	p, err := h.svc.Update(pid, id, patch)
	if err != nil {
		writeErr(w, statusFor(err), err.Error())
		return
	}
	writeJSON(w, 200, p)
}

// Archive
// @Summary  Arquivar anúncio
// @Tags     postings
// @Param    id      path   string  true  "Posting ID"
// @Router   /postings/{id}/archive [post]
func (h *Handler) Archive(w http.ResponseWriter, r *http.Request) {
	pid, ok := authmw.UserIDFromContext(r)
	if !ok {
		writeErr(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	id := strings.TrimPrefix(r.URL.Path, "/api/v1/postings/")
	id = strings.TrimSuffix(id, "/archive")

	if err := h.svc.Archive(pid, id); err != nil {
		writeErr(w, statusFor(err), err.Error())
		return
	}
	writeJSON(w, 200, map[string]string{"status": "ok"})
}

// ListMine
// @Summary  Listar meus anúncios
// @Tags     postings
// @Router   /postings/mine [get]
func (h *Handler) ListMine(w http.ResponseWriter, r *http.Request) {
	pid, ok := authmw.UserIDFromContext(r)
	if !ok {
		writeErr(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	list, _ := h.svc.ListMine(pid)
	writeJSON(w, 200, list)
}

// GetPublic
// @Summary  Detalhar anúncio público
// @Tags     postings
// @Param    id  path  string  true  "Posting ID"
// @Router   /postings/{id} [get]
func (h *Handler) GetPublic(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/api/v1/postings/")
	p, err := h.svc.GetPublic(id)
	if err != nil {
		writeErr(w, 404, "not found")
		return
	}
	writeJSON(w, 200, p)
}

// Search
// @Summary  Buscar/Listar anúncios com filtros
// @Tags     postings
// @Router   /postings [get]
func (h *Handler) Search(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	p := domain.SearchParams{
		Query:    q.Get("q"),
		Category: q.Get("category"),
		City:     q.Get("city"),
		District: q.Get("district"),
		Sort:     q.Get("sort"),
		Order:    q.Get("order"),
	}
	if v := q.Get("price_min"); v != "" {
		p.PriceMin = parseI64(v)
	}
	if v := q.Get("price_max"); v != "" {
		p.PriceMax = parseI64(v)
	}
	if v := q.Get("rating_min"); v != "" {
		p.RatingMin = parseF64(v)
	}
	if v := q.Get("limit"); v != "" {
		p.Limit = parseI(v)
	}
	if v := q.Get("offset"); v != "" {
		p.Offset = parseI(v)
	}

	items, next := h.svc.Search(p)
	writeJSON(w, 200, map[string]any{
		"items":       items,
		"next_offset": next,
	})
}

func parseI(s string) int       { i, _ := strconv.Atoi(s); return i }
func parseI64(s string) int64   { i, _ := strconv.ParseInt(s, 10, 64); return i }
func parseF64(s string) float64 { f, _ := strconv.ParseFloat(s, 64); return f }

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}
func writeErr(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}
func statusFor(err error) int {
	switch err {
	case domain.ErrForbidden:
		return 403
	case domain.ErrNotFound:
		return 404
	default:
		return 400
	}
}
