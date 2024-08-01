package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	"github.com/rs/cors"
)

type Comment struct {
	ID      int    `json:"id"`
	Content string `json:"content"`
}

type Post struct {
	ID        int       `json:"id"`
	Content   string    `json:"content"`
	Comments  []Comment `json:"comments"`
	Likes     int       `json:"likes"`
	Dislikes  int       `json:"dislikes"`
	ShareLink string    `json:"share_link,omitempty"`
}

type PostFlow struct {
	posts  map[int]Post
	nextID int
	mu     sync.Mutex
}

func NewPostFlow() *PostFlow {
	return &PostFlow{
		posts:  make(map[int]Post),
		nextID: 1,
	}
}

func (pf *PostFlow) AddPost(content string) int {
	pf.mu.Lock()
	defer pf.mu.Unlock()
	post := Post{
		ID:       pf.nextID,
		Content:  content,
		Comments: []Comment{}, // Initialize comments as an empty array
	}
	pf.posts[pf.nextID] = post
	pf.nextID++
	return post.ID
}

func (pf *PostFlow) AddComment(postID int, content string) {
	pf.mu.Lock()
	defer pf.mu.Unlock()
	comment := Comment{
		ID:      len(pf.posts[postID].Comments) + 1,
		Content: content,
	}
	post := pf.posts[postID]
	post.Comments = append(post.Comments, comment)
	pf.posts[postID] = post
}

func (pf *PostFlow) LikePost(postID int) {
	pf.mu.Lock()
	defer pf.mu.Unlock()
	post := pf.posts[postID]
	post.Likes++
	pf.posts[postID] = post
}

func (pf *PostFlow) DislikePost(postID int) {
	pf.mu.Lock()
	defer pf.mu.Unlock()
	post := pf.posts[postID]
	post.Dislikes++
	pf.posts[postID] = post
}

func (pf *PostFlow) SharePost(postID int) string {
	pf.mu.Lock()
	defer pf.mu.Unlock()
	post := pf.posts[postID]
	post.ShareLink = fmt.Sprintf("https://postflow.com/post/%d", postID)
	pf.posts[postID] = post
	return post.ShareLink
}

func (pf *PostFlow) GetPosts(w http.ResponseWriter, r *http.Request) {
	pf.mu.Lock()
	defer pf.mu.Unlock()
	posts := make([]Post, 0, len(pf.posts))
	for _, post := range pf.posts {
		posts = append(posts, post)
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(posts)
}

func (pf *PostFlow) addPostHandler(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Content string `json:"content"`
	}
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	postID := pf.AddPost(data.Content)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]int{"post_id": postID})
}

func (pf *PostFlow) addCommentHandler(w http.ResponseWriter, r *http.Request) {
	var data struct {
		PostID  int    `json:"post_id"`
		Content string `json:"content"`
	}
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	pf.AddComment(data.PostID, data.Content)
	w.WriteHeader(http.StatusCreated)
}

func (pf *PostFlow) likePostHandler(w http.ResponseWriter, r *http.Request) {
	var data struct {
		PostID int `json:"post_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	pf.LikePost(data.PostID)
	w.WriteHeader(http.StatusOK)
}

func (pf *PostFlow) dislikePostHandler(w http.ResponseWriter, r *http.Request) {
	var data struct {
		PostID int `json:"post_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	pf.DislikePost(data.PostID)
	w.WriteHeader(http.StatusOK)
}

func (pf *PostFlow) sharePostHandler(w http.ResponseWriter, r *http.Request) {
	var data struct {
		PostID int `json:"post_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	link := pf.SharePost(data.PostID)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"share_link": link})
}

func main() {
	pf := NewPostFlow()

	corsHandler := cors.Default().Handler(http.DefaultServeMux)

	http.HandleFunc("/add_post", pf.addPostHandler)
	http.HandleFunc("/add_comment", pf.addCommentHandler)
	http.HandleFunc("/like_post", pf.likePostHandler)
	http.HandleFunc("/dislike_post", pf.dislikePostHandler)
	http.HandleFunc("/share_post", pf.sharePostHandler)
	http.HandleFunc("/posts", pf.GetPosts)

	fmt.Println("Server is running on port 8080...")
	http.ListenAndServe(":8080", corsHandler)
}
