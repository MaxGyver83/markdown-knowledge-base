package api

type CreateDocumentRequest struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

type UpdateDocumentRequest struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}
