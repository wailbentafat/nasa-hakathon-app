# Simple Multiplayer FPS Game WebSocket Server

## Overview
This project is a WebSocket server built with Go for a simple multiplayer first-person shooter (FPS) game. It supports real-time player movement, shooting mechanics, and damage handling, allowing for interactive gameplay and stor every user information and mouvemnt inside blckchain operateur 

## Key Features
- Real-time player movement synchronization
- Shooting mechanics with damage calculation
- Efficient WebSocket communication
- blackchian relatime mechanisme 

## Technologies Used
- Go
- Gorilla WebSocket package

## Getting Started

### Prerequisites
- Go (v1.16 or higher)

### Installation

1. Clone the repository.
2. Build the server using Go.
3. Run the server.

The server will run on `http://localhost:8080` by default.

## Usage
Players can connect via WebSocket to send movement and shooting data, which the server processes and broadcasts to all clients.

## Project Structure
```
multiplayer-fps-server/
├── main.go             # Main server file
├── go.mod              # Go module file
└── README.md           # Project documentation
```
 ##testing:
![image](https://github.com/user-attachments/assets/880b1f16-2233-47e0-93ad-30e0287cbbfe)

## License
This project is licensed under the MIT License.

## Acknowledgments
Inspired by various multiplayer game tutorials and Go documentation.

Feel free to contribute or report issues!
