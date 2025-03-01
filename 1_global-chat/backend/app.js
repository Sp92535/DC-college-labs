// app.js

const express = require('express');
const app = express();
const router = require("./routes")
const cors = require('cors');

app.use(express.json());
app.use(cors())
app.use(router);

app.get('/', (req, res) => {
    res.send("YO").status(200);
})


module.exports = app