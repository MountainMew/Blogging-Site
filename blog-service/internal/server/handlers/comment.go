package handlers

import (
	"blog-service/internal/db/mongo"
	"blog-service/internal/db/postgres"
	"blog-service/internal/server/models"
	"encoding/json"
	"errors"
	"net/http"
	"regexp"
	"strconv"
)

type CommentHandler struct {
	MongoDB    *mongo.Client
	PostgresDB *postgres.Client
}

var (
	CommentDeleteRe = regexp.MustCompile(`/comment/\d+`)
	CommentLikeRe   = regexp.MustCompile(`^/comment/\d+/like$`)
)

func (h *CommentHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch {
	case r.Method == http.MethodDelete && CommentDeleteRe.MatchString(r.URL.Path):
		h.CommentDelete(w, r)
		return
	case r.Method == http.MethodPost && CommentLikeRe.MatchString(r.URL.Path):
		h.CommentLike(w, r)
		return
	case r.Method == http.MethodPost:
		h.CommentCreate(w, r)
		return
	}
}

func (h *CommentHandler) CommentDelete(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (h *CommentHandler) CommentLike(w http.ResponseWriter, r *http.Request) {
	commentID := r.PathValue("id")
	intID, err := strconv.Atoi(commentID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var unauthorizedErr *models.UnauthorizedError
	var dbInternalErr *models.DBInternalError

	err = models.CommentAddRemoveLike(r.Context(), h.PostgresDB, intID)
	switch {
	case errors.As(err, &unauthorizedErr):
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	case errors.As(err, &dbInternalErr):
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *CommentHandler) CommentCreate(w http.ResponseWriter, r *http.Request) {
	var comment models.CommentCreateDTO
	var invalidArticleErr *models.InvalidArticleError
	var paramErr *models.ParamError
	var unauthorizedErr *models.UnauthorizedError

	err := json.NewDecoder(r.Body).Decode(&comment)
	if err != nil {
		http.Error(w, paramErr.Error(), http.StatusBadRequest)
		return
	}

	err = models.CreateComment(r.Context(), h.PostgresDB, h.MongoDB, &comment)
	if err != nil {
		switch {
		case errors.As(err, &invalidArticleErr):
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		case errors.As(err, &paramErr):
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		case errors.As(err, &unauthorizedErr):
			http.Error(w, err.Error(), http.StatusUnauthorized)
		default:
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
}
