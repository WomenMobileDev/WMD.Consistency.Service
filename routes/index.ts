import express from "express";

const app = express.Router();

app.use("/auth", require("./auth"));

module.exports = app;