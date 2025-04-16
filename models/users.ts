import config from "config";
import { pool } from "../config/db";
import { userData } from "../types/global";
import { logger } from "../utils/logger";
import { ResultSetHeader } from "mysql2";

export const createUser = async (user: userData): Promise<void> => {
    const query = `
        INSERT INTO ${config.get("mysqlRemote.db")}.users 
        (id, first_name, last_name, email, phone, log_in_provider)
        VALUES (?, ?, ?, ?, ?, ?)
    `;

    const values = [
        user.id,
        user.firstName,
        user.lastName,
        user.email,
        user.phone,
        user.logInProvider,
    ];

    try {
        const [result] = await pool.query<ResultSetHeader>(query, values);
        logger.info(`User added successfully. Insert ID: ${result.insertId}`);
    } catch (error) {
        logger.error("Error inserting user:", error);
        throw error;
    }
};

export const isUserAlreadyExist = async (email: string): Promise<boolean> => {
    const query = `
        SELECT 1
        FROM ${config.get("mysqlRemote.db")}.users
        WHERE email = ?
        LIMIT 1
    `;

    try {
        const [rows] = await pool.query(query, [email]);

        const exists = Array.isArray(rows) && rows.length > 0;

        if (exists) {
            logger.info(`User already exists with email: ${email}`);
        }

        return exists;
    } catch (error) {
        logger.error("Error checking if user exists:", error);
        throw error;
    }
};