// This interface now matches the Go 'Like' struct.
export interface Like {
  id: number
  commentId: number
  userId: number
}

// This interface now matches the Go 'Comment' struct and includes likes.
export interface Comment {
  id: number
  userId: number
  articleId: string
  content: string
  createdAt: string // Go's time.Time marshals to a string in JSON
  commenterName?: string // Optional: In a real app, you'd join this data on the backend.
  likes: number // An array of likes associated with the comment
  liked: boolean
}
