import type { Comment } from './comment.interface'

export interface ArticleDTO {
  id: string
  title: string
  content: string
  category: string
  publisherName: string
  publisherId: number
  createdAt: string
  comments: Comment[]
}
