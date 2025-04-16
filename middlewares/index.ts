const express = require("express");
const boom = require("express-boom");
import passport from "passport";
import morgan from "morgan";
import { logger } from "../utils/logger";

// auth middleware
require("./passport");

const middleware = (app) => {
    app.use(boom());
    app.use(express.json());
    app.use(morgan("combined", { stream: logger.stream }));
    app.use(passport.initialize());
};

module.exports = middleware;
