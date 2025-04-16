import express, { Request, Response } from "express";
const app = express();
const config = require("config");
import { logger } from "./utils/logger";
const indexRouter = require("./routes/index");
const appMiddleware = require("./middlewares");

const PORT = 3000;

global.config = config;
global.logger = logger;

app.get("/", (req: Request, res: Response) => {
  res.json({ message: `Server is up and running 🎉` });
});

appMiddleware(app);
app.use("/", indexRouter);

app.listen(PORT, () => {
  logger.info(`Server is up and running 🎉 on port http://localhost:${PORT}`);
});