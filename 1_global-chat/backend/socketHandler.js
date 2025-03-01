const Chat = require('./model');

module.exports = (io) => {
    io.on("connection", (socket) => {
        console.log("A user connected");

        // Join a room
        socket.on("joinRoom", async (room_id) => {

            const room = await Chat.findOne({ room_id });

            if (!room) {
                socket.emit("error", "Room not found");
                return;
            }

            socket.join(room_id);
            console.log(`User joined room: ${room_id}`);
        });

        // Handle new message
        socket.on("newMessage", async (data) => {
            const { room_id, username, msg } = data;

            try {
                console.log("MESSAGE");
                const room = await Chat.findOne({ room_id });

                if (!room) {
                    socket.emit("error", "Room not found");
                    return;
                }
                
                room.messages.push({ username, msg });
                await room.save();
                
                io.to(room_id).emit("message", { username, msg, timestamp: new Date()});
            } catch (error) {
                socket.emit("error", "Failed to save message");
                console.error(error);
            }
        });

        // Handle user disconnect
        socket.on("disconnect", () => {
            console.log("A user disconnected");
        });
    });
};
