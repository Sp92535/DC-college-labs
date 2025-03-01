// routes.js
const express = require('express');
const router = express.Router();
const crypto = require('crypto');
const Chat = require('./model');


// Create a room
router.get("/create-room", async (req, res) => {

    try {
        const code = crypto.randomBytes(4).toString("hex").toUpperCase();

        const room = new Chat({ room_id: code, messages: [] });
        await room.save();

        res.json({ room_id: code }).status(200);
    } catch (error) {
        res.json({ error: "Internal Server Error." }).status(500);
    }

});

router.get("/messages/:roomId", async (req, res) => {
    try {
        const { roomId } = req.params;

        const room = await Chat.findOne({ room_id: roomId });
        
        if (!room) {
            res.json({ error: "Room not found." }).status(404);
        }

        res.send({ messages: room.messages }).status(200);

    } catch (error) {
        res.json({ error: "Internal Server Error." }).status(500);
    }
});

router.get("/",(req,res)=>{
    res.send("<h1>ALL WORKING GOOD</h1>")
})

module.exports = router;
