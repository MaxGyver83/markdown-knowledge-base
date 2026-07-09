# Markdown Knowledge Base

- Backend: REST API written in Go (using go-chi for routing)
- Frontend: HMTL + Javascript + CSS

## How to test

Start backend (REST API):

```sh
cd backend/markdown-api
HTTP_PORT=8085 go run ./cmd/server
```

Start frontend (web server):

```sh
cd frontend
python3 -m http.server 3000
```

or (if available):

```sh
cd frontend
serve
```

(`serve` can be installed with `sudo npm install -g serve`)
