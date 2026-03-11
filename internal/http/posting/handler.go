package posting

import (
	"encoding/json"
	"net/http"
	"strconv"

	authmw "github.com/Gab-Mello/service-finder/internal/http/middleware/auth"
	"github.com/Gab-Mello/service-finder/internal/http/response"
	domain "github.com/Gab-Mello/service-finder/internal/posting"
)

const basePath = "/api/v1/postings/"

type Handler struct{ svc *domain.Service }

func NewHandler(s *domain.Service) *Handler { return &Handler{svc: s} }

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	pid, ok := authmw.UserIDFromContext(r)
	if !ok {
		response.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	var req CreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid json")
		return
	}

	p, err := h.svc.Create(pid, req.Title, req.Description, req.Price, req.Category, req.City, req.District)
	if err != nil {
		response.Error(w, statusFor(err), err.Error())
		return
	}
	response.JSON(w, http.StatusCreated, p)
}

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	pid, ok := authmw.UserIDFromContext(r)
	if !ok {
		response.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	id := response.PathParam(r.URL.Path, basePath, "")
	var patch map[string]any
	if err := json.NewDecoder(r.Body).Decode(&patch); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid json")
		return
	}
	p, err := h.svc.Update(pid, id, patch)
	if err != nil {
		response.Error(w, statusFor(err), err.Error())
		return
	}
	response.JSON(w, http.StatusOK, p)
}

func (h *Handler) Archive(w http.ResponseWriter, r *http.Request) {
	pid, ok := authmw.UserIDFromContext(r)
	if !ok {
		response.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	id := response.PathParam(r.URL.Path, basePath, "/archive")

	if err := h.svc.Archive(pid, id); err != nil {
		response.Error(w, statusFor(err), err.Error())
		return
	}
	response.JSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (h *Handler) ListMine(w http.ResponseWriter, r *http.Request) {
	pid, ok := authmw.UserIDFromContext(r)
	if !ok {
		response.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	list, err := h.svc.ListMine(pid)
	if err != nil {
		response.InternalError(w, err)
		return
	}
	response.JSON(w, http.StatusOK, list)
}

func (h *Handler) GetPublic(w http.ResponseWriter, r *http.Request) {
	id := response.PathParam(r.URL.Path, basePath, "")
	p, err := h.svc.GetPublic(id)
	if err != nil {
		response.Error(w, http.StatusNotFound, "not found")
		return
	}
	response.JSON(w, http.StatusOK, p)
}

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
	response.JSON(w, http.StatusOK, map[string]any{
		"items":       items,
		"next_offset": next,
	})
}

func parseI(s string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		return 0
	}
	return i
}

func parseI64(s string) int64 {
	i, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return 0
	}
	return i
}

func parseF64(s string) float64 {
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0
	}
	return f
}

func statusFor(err error) int {
	switch err {
	case domain.ErrForbidden:
		return http.StatusForbidden
	case domain.ErrNotFound:
		return http.StatusNotFound
	case domain.ErrInvalidFields:
		return http.StatusBadRequest
	default:
		return http.StatusInternalServerError
	}
}
