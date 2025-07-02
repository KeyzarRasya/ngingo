require('dotenv').config();
const express = require('express')

const app = express()

const PORT = process.env.PORT;

app.get("/", (req, res) => {
    for (let i = 0; i < 100000; i++) {

    }

    res.send(`Hello From ${PORT}`)
})

app.get("/say/:name", (req,res) => {
    const {name} = req.params
    res.send(`Hello, ${name}`)
})

app.listen(PORT, () => {
    console.log(`Running from ${PORT}`);
})