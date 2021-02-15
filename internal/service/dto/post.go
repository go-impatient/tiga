package dto

// PostRequest 请求数据结果
type PostRequest struct {
	Title   string `json:"title,omitempty"`
	Content string `json:"content,omitempty"`
}

// PostResponse 返回数据结构
type PostResponse struct {
	Id int64 `json:"id,omitempty"`
}

