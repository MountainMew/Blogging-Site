package handlers

import (
	"blog-service/internal/db/mongo"
	"encoding/json"
	"log"
	"net/http"
	"regexp"
	"strconv"
)

const pagelimit = 10

type BlogHandler struct {
	Mongo *mongo.Client
}

var (
	BlogHandlerGetRe  = regexp.MustCompile(`/blog$`)
	BlogByPublisherRe = regexp.MustCompile(`/blog/by-publisher$`)
)

func (h *BlogHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Only handle GET requests
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	switch {
	case BlogByPublisherRe.MatchString(r.URL.Path):
		h.getArticlesByPublisher(w, r)

	case BlogHandlerGetRe.MatchString(r.URL.Path):
		h.getBlogPage(w, r)

	default:
		http.NotFound(w, r)
	}
}

func (h *BlogHandler) getBlogPage(w http.ResponseWriter, r *http.Request) {
	page := int64(1)
	limit := int64(pagelimit)

	pageStr := r.URL.Query().Get("page")
	if pageStr != "" {
		p, err := strconv.ParseInt(pageStr, 10, 64)
		if err == nil && p > 0 {
			page = p
		}
	}

	articles, err := h.Mongo.GetArticlesByEngagement(page, limit)
	if err != nil {
		log.Printf("Error retrieving articles: %v", err)
		http.Error(w, "Failed to retrieve articles", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(articles); err != nil {
		log.Printf("Error encoding response: %v", err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

// getArticlesByPublisher fetches articles using the 'id' query parameter.
func (h *BlogHandler) getArticlesByPublisher(w http.ResponseWriter, r *http.Request) {
	publisherIDStr := r.URL.Query().Get("id")
	if publisherIDStr == "" {
		http.Error(w, "Missing 'id' query parameter", http.StatusBadRequest)
		return
	}

	publisherID, err := strconv.Atoi(publisherIDStr)
	if err != nil {
		http.Error(w, "Invalid 'id' query parameter, must be an integer", http.StatusBadRequest)
		return
	}

	articles, err := h.Mongo.FindArticlesByPublisherID(r.Context(), publisherID)
	if err != nil {
		log.Printf("Error retrieving articles by publisher ID %d: %v", publisherID, err)
		http.Error(w, "Failed to retrieve articles", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(articles); err != nil {
		log.Printf("Error encoding response: %v", err)
	}
}
