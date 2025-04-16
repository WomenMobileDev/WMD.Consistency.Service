import express from "express";
import { devFlagMiddleware } from "../middlewares/devFlag";
import { googleAuthCallback, googleAuthLogin } from "../controllers/auth";

const router = express.Router();

router.get("/google/login", devFlagMiddleware, googleAuthLogin);
router.get("/google/callback", googleAuthCallback);

module.exports = router;