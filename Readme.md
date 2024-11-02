# Storytelling Game Backend

This is the backend for a real-time multiplayer storytelling game where players collaboratively write a story, one line at a time. This backend is developed using Golang and uses WebSockets for real-time communication between players.

## Features

- **Create and Join Rooms**: Players can create storytelling rooms or join existing ones using a unique room ID.
- **Real-Time Collaboration**: Players can see updates to the story in real time.
- **Turn-Based Storytelling**: Players take turns to add lines to the story, with the current player notified via WebSocket.
- **Game Management**: Host starts the game, and all players are notified when their turn arrives or when someone leaves.
- **WebSocket Integration**: Used for real-time updates, ensuring smooth gameplay.
- **Cross-Origin Resource Sharing (CORS)**: Allows frontend interaction from different domains.

## Getting Started

### Prerequisites

- **Golang**: Make sure [Golang](https://golang.org/doc/install) is installed (v1.18+ recommended).
- **Docker (optional)**: If you want to run the application in a Docker container.
- **Environment Variables**: Set up the following environment variable:
  - `PORT`: Port for the server to listen on (default is `8080`).

### Installation

1. **Clone the repository**:
   ```bash
   git clone https://github.com/your-username/storytelling-backend.git
   cd storytelling-backend
   ```

2. **Install dependencies**:
   ```bash
   go mod tidy
   ```

3. **Set up environment variables**:
   Create a `.env` file in the root of the project and specify the variables. For example:
   ```plaintext
   PORT=8080
   ```

4. **Run the application**:
   ```bash
   go run cmd/main.go
   ```

5. **Using Docker** (Optional):
   ```bash
   docker build -t storytelling-backend .
   docker run -p 8080:8080 storytelling-backend
   ```

### Folder Structure

- `cmd/`: Contains main server setup and routing.
- `config/`: Handles configuration loading.
- `internal/api/`: API handlers and WebSocket handlers.
- `internal/game/`: Game logic for room and player management.
- `internal/models/`: Structures and logic for player connections and rooms.
- `pkg/utils/`: Utility functions, including generating unique room IDs.

## API Endpoints

### HTTP Routes

| Method | Route                  | Description                  |
|--------|-------------------------|------------------------------|
| POST   | `/create-room`          | Creates a new room           |
| POST   | `/join-room`            | Joins an existing room       |
| POST   | `/start-game/{room_id}` | Starts the game in a room    |
| POST   | `/submit-line`          | Adds a line to the story     |
| GET    | `/get-story`            | Retrieves the current story  |
| GET    | `/ws`                   | WebSocket connection for real-time updates |

### Example Request

**Creating a Room**
```bash
curl -X POST http://localhost:8080/create-room \\
-H \"Content-Type: application/json\" \\
-d '{
    \"story_name\": \"A New Adventure\",
    \"player_name\": \"Alice\"
}'
```

**Joining a Room**
```bash
curl -X POST http://localhost:8080/join-room \\
-H \"Content-Type: application/json\" \\
-d '{
    \"room_id\": \"12345\",
    \"player_name\": \"Bob\"
}'
```

### WebSocket Usage

Connect to WebSocket with: `ws://localhost:8080/ws?room_id={room_id}\u0026player_name={player_name}`

- **SUBMIT_LINE**: Submit a line for your turn.
- **START_GAME**: Start the game (only host can initiate).

## Contributing

Contributions are welcome! Feel free to open issues or submit pull requests.

1. Fork the repository.
2. Create a new branch for your feature.
3. Commit and push your changes.
4. Open a pull request.

## License

This project is licensed under the MIT License. See the `LICENSE` file for more details.

---

Enjoy storytelling with friends and watch the story unfold!
```