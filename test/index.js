require('dotenv').config();
const express = require('express')

const app = express()

const PORT = process.env.PORT;

function multiplyMatrix(size) {
  const a = [...Array(size)].map(() => Array(size).fill(1));
  const b = [...Array(size)].map(() => Array(size).fill(1));
  const result = Array(size).fill(0).map(() => Array(size).fill(0));

  for (let i = 0; i < size; i++) {
    for (let j = 0; j < size; j++) {
      for (let k = 0; k < size; k++) {
        result[i][j] += a[i][k] * b[k][j];
      }
    }
  }
  return result;
}

app.get("/", (req, res) => {
    for (let i = 0; i < 100000; i++) {

    }

    res.send(`Hello From ${PORT}`)
})

app.get("/say/:name", (req,res) => {
    const {name} = req.params
    res.send(`Hello, ${name}`)
})

app.get("/matrix", (req, res) => {
    multiplyMatrix(100);
    res.send("OKE");
})

app.listen(PORT, () => {
    multiplyMatrix(300);
    console.log(`Running from ${PORT}`);
})