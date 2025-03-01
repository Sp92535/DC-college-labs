// model.js
const mongoose = require('mongoose');

const Chat  = mongoose.Schema({

    room_id : {type: String, required: true, unique : true},
    messages : [

        {
            username : {type: String , required: true},
            msg: {type: String, required: true},
            timestamp: {type: Date, default: new Date()}
        }

    ]
});

module.exports = mongoose.model("Chat", Chat);