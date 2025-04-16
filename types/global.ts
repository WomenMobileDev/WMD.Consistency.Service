import { Request, Response } from "express";
import { Boom } from "express-boom";
export type userData = {
    id: string,
    firstName: string,
    lastName: string,
    email: string,
    phone?: number | null,
    logInProvider: string,
};
export type CustomRequest = Request & { userData };
export type CustomResponse = Response & { boom: Boom };