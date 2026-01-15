# Rate Limiter Playground

A highly interactive Rate Limiter visualization tool built with Go and vanilla HTML/CSS/JS. This project demonstrates various rate limiting algorithms in action, featuring a real-time dashboard to visualize allowed vs. denied requests.

## Features

- **Multiple Algorithms**:
  - Token Bucket
  - Leaky Bucket
  - Fixed Window Counter
  - Sliding Window Log
  - Sliding Window Counter
  - Concurrent Request Limiter
- **Interactive UI**:
  - Real-time visualization of requests as moving dots.
  - Live statistics (Allowed vs. Denied).
  - Configurable parameters for each algorithm.
  - "Auto-Spike" button to simulate traffic bursts.
- **Backend**: Robust Go implementation of rate limiting logic.

## Project Structure

```
├── cmd/
│   ├── server/         # Web server entry point
│   └── simulation/     # CLI simulation tool (local only)
├── pkg/
│   └── ratelimit/      # Core rate limiting algorithms
├── public/             # Frontend assets (HTML/CSS/JS)
├── docker/             # Docker configuration
└── go.mod              # Go module definition
```

## Getting Started

### Prerequisites

- [Go 1.25+](https://go.dev/dl/) installed.
- (Optional) [Docker](https://www.docker.com/) for containerized deployment.

### Running Locally (Go)

1. **Clone the repository**:
   ```bash
   git clone https://github.com/YOUR_USERNAME/ratelimiter.git
   cd ratelimiter
   ```

2. **Install dependencies**:
   ```bash
   go mod download
   ```

3. **Run the server**:
   ```bash
   go run cmd/server/main.go
   ```

4. **Access the application**:
   Open your browser and navigate to [http://localhost:8080](http://localhost:8080).

### Running with Docker

1. **Build the image**:
   ```bash
   docker build -t ratelimiter .
   ```

2. **Run the container**:
   ```bash
   docker run -p 8080:8080 ratelimiter
   ```

3. **Access the application**:
   Visit [http://localhost:8080](http://localhost:8080).

## Deployment

### Deploy to Render

1. Create a new **Web Service** on [Render](https://render.com/).
2. Connect your GitHub repository.
3. Select **Docker** as the Runtime.
4. Render will automatically build and deploy your application.

## Usage

1. Select an algorithm from the dropdown menu.
2. Adjust the parameters (e.g., Request Rate, Capacity) and click **Apply Configuration**.
3. Use the **Single Request** button to test manually or **Auto-Spike** to simulate load.
4. Watch the visualization area:
   - **Green Dots**: Allowed requests.
   - **Red Dots**: Denied requests.

## License

MIT
