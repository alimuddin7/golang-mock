# 🚀 GopherMock: The Ultimate Mock API Server & Editor

[![Go Version](https://img.shields.io/badge/Go-1.22+-00ADD8?style=for-the-badge&logo=go)](https://golang.org)
[![Fiber Framework](https://img.shields.io/badge/Fiber-v2-EF4444?style=for-the-badge&logo=fastapi)](https://gofiber.io)
[![Alpine.js](https://img.shields.io/badge/Alpine.js-3.x-8BC0D0?style=for-the-badge&logo=alpine.js)](https://alpinejs.dev)
[![DaisyUI](https://img.shields.io/badge/DaisyUI-v5-5852D6?style=for-the-badge&logo=daisyui)](https://daisyui.com)

**GopherMock** is a powerful, lightweight, and easy-to-use mock server built with Go. It features a sleek web interface that allows you to manage mock configurations, import OpenAPI specifications, and evaluate complex response rules in real-time.

---

## ✨ Features

- 🖥️ **Interactive Web UI**: Manage all your mocks from a beautiful, responsive dashboard built with Tailwind CSS and DaisyUI.
- ⚡ **Conditional Responses**: Define advanced logic with rules (AND/OR) targeting request body, headers, query parameters, and path variables.
- 📥 **OpenAPI/Swagger Import**: Quickly bootstrap your mock server by importing OpenAPI 3.x or Swagger 2.0 specifications.
- 🔄 **Real-time Configuration**: Changes are saved instantly to `configs.json` and applied without restarting the server.
- 📋 **Integrated Logger**: Monitor incoming requests and outgoing responses directly in the console with structured logging.
- 📦 **Bulk Actions**: Duplicate, delete, or bulk-delete configurations with ease.
- 🐳 **Docker Ready**: Includes `Dockerfile` and `docker-compose.yml` for seamless deployment.

---

## 🛠️ Technology Stack

- **Backend**: [Go](https://golang.org) with [Fiber](https://gofiber.io) (Express-inspired web framework).
- **Frontend**: [Alpine.js](https://alpinejs.dev) for lightweight reactivity and [Tailwind CSS](https://tailwindcss.com) + [DaisyUI](https://daisyui.com) for a premium UI.
- **Templating**: Go's native `html/template` engine.
- **Serialization**: JSON for persistence and API communication.

---

## 🚀 Getting Started

### Prerequisites

- [Go 1.22+](https://golang.org/dl/)
- [Docker](https://www.docker.com/) (optional)

### Running Locally

1. **Clone the repository**:
   ```bash
   git clone https://github.com/alimuddin7/gopher-mock.git
   cd gopher-mock
   ```

2. **Install dependencies**:
   ```bash
   go mod download
   ```

3. **Run the application**:
   ```bash
   go run main.go
   ```
   The server will start on `http://localhost:3000`.

### Running with Docker

Use Docker Compose for a quick start:
```bash
docker-compose up -d --build
```

---

## 📖 Usage Guide

### 1. Creating a Mock
Click the **+** button in the sidebar to create a new configuration. Fill in the name, method, path, and default response.

### 2. Adding Conditional Rules
Go to the **Conditional Rules** tab to add logic. For example, return a `400 Bad Request` if a specific field in the request body matches a value.

### 3. Importing OpenAPI
Click the **Import** button in the navbar, upload your specification file, and choose whether to merge or replace existing mocks.

### 4. Bulk Delete
Use the checkboxes in the sidebar to select multiple mocks and click the trash icon in the header to delete them all at once.

---

## 📂 Project Structure

```text
.
├── handler/            # Fiber request handlers
├── service/            # Business logic (OpenAPI parser, Rule engine)
├── model/              # Data structures
├── templates/          # HTML partials and views
├── static/             # Frontend assets
├── configs.json        # Persistent configuration storage
├── main.go             # Entry point
└── Dockerfile          # Container configuration
```

---

## 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

---

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

---

<p align="center">Made with ❤️ by the GopherMock Team</p>
