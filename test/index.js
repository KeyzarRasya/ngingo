require('dotenv').config();
const express = require('express')

const app = express()

const PORT = process.env.PORT;

app.get("/", (req, res) => {
    res.send(`Hello from port ${PORT}`)
})

app.listen(PORT, () => {
    console.log(`Running from ${PORT}`);
})