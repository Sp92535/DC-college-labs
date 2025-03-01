// server.js

const app = require('./app');
const connectDB = require('./connectDB');
const { Server } = require("socket.io");
const http = require("http");
const socketHandler = require('./socketHandler');
require("dotenv").config();

const server = http.createServer(app);

const io = new Server(server, {
    cors: { origin: "*" }
});

app.set("io", io);

socketHandler(io);

PORT = process.env.SERVER_PORT || 5000
HOST = "0.0.0.0"

server.listen(5000, HOST, () => {
    console.log(`SERVER LISTENING ON ${HOST}:${PORT}`);
});

connectDB();