package controller

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/kalpit-sharma-dev/chat-service/src/models"
	"github.com/kalpit-sharma-dev/chat-service/src/service"
)

type ReelController struct {
	Service *service.ReelService
}

func (c *ReelController) UploadReel(w http.ResponseWriter, r *http.Request) {
	// Extract data from request and create a Reel object
	var reel models.Reel
	err := json.NewDecoder(r.Body).Decode(&reel)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Call the service to save the reel
	if err := c.Service.CreateReel(&reel); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (c *ReelController) FetchReels(w http.ResponseWriter, r *http.Request) {
	lastIDStr := r.URL.Query().Get("last_id")
	lastID, _ := strconv.Atoi(lastIDStr)

	reels, err := c.Service.FetchReels(lastID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(reels)
}

// Implement other handlers like LikeReel, CommentReel, etc.

func (c *ReelController) LikeReel(w http.ResponseWriter, r *http.Request) {
	reelID, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		http.Error(w, "Invalid reel ID", http.StatusBadRequest)
		return
	}

	userID := getUserIDFromRequest(r) // Implement this function based on your authentication system
	if err := c.Service.LikeReel(userID, reelID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (c *ReelController) UnlikeReel(w http.ResponseWriter, r *http.Request) {
	reelID, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		http.Error(w, "Invalid reel ID", http.StatusBadRequest)
		return
	}

	userID := getUserIDFromRequest(r) // Implement this function based on your authentication system
	if err := c.Service.UnlikeReel(userID, reelID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (c *ReelController) CommentOnReel(w http.ResponseWriter, r *http.Request) {
	reelID, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		http.Error(w, "Invalid reel ID", http.StatusBadRequest)
		return
	}

	var comment models.Comment
	err = json.NewDecoder(r.Body).Decode(&comment)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	comment.ReelID = reelID
	comment.CreatedAt = time.Now()

	if err := c.Service.CommentOnReel(&comment); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (c *ReelController) GetCommentsForReel(w http.ResponseWriter, r *http.Request) {
	reelID, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		http.Error(w, "Invalid reel ID", http.StatusBadRequest)
		return
	}

	comments, err := c.Service.GetCommentsForReel(reelID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(comments)
}

// Utility function to extract user ID from request (example placeholder)
func getUserIDFromRequest(r *http.Request) int {
	// This is just a placeholder. Replace with your actual user authentication logic.
	return 1
}
