// server.js

const app = require('./app');
const connectDB = require('./connectDB');
const { Server } = require("socket.io");
const http = require("http");
const socketHandler = require('./socketHandler');

const server = http.createServer(app);

const io = new Server(server,{
    cors: {origin:"*"}
});

app.set("io",io);

socketHandler(io);

server.listen(5000,()=>{
    console.log("Listening on 5000");
});

connectDB();