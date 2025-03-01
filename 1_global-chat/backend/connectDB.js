const mongoose = require('mongoose');
const { DB } = require('./config');

const MONGO_URI = DB

const connectDB = () =>{
    mongoose
    .connect(MONGO_URI)
    .then(()=>{
        console.log("CONNECTED TO MONGO DB");
    })
    .catch((err)=>{
        console.log("Error Connecting DB",err);
    })
}

module.exports = connectDB