import { NextFuncion } from "express";
import { CustomRequest, CustomResponse } from "../types/global";
import { ROUTE_NOT_FOUND } from "../constants/constants";
import { logger } from "../utils/logger";

export const devFlagMiddleware = (req: CustomRequest, res: CustomResponse, next: NextFuncion) => {
    try {
        const dev = req.query.dev === "true";

        if (!dev) {
            return res.boom.notFound(ROUTE_NOT_FOUND);
        }
        next();
    } catch (error) {
        logger.error("Error occurred in devFlagMiddleware:", error.message);
        next(error);
    }
};