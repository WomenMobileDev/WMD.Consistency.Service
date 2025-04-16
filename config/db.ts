import mysql from "mysql2/promise";
import config from "config";
import { logger } from "../utils/logger";

// export const pool = mysql.createPool({
//     host     : config.get("mysqlRemote.host"),
//     user     : config.get("mysqlRemote.username"),
//     password : config.get("mysqlRemote.password"),
//     port: config.get("mysqlRemote.port"),
//     database : config.get("mysqlRemote.db"),
//     waitForConnections: true,
//     connectionLimit: 10,
//     queueLimit: 0,
// });

export const pool = mysql.createPool({
    host: config.get("mysqlLocal.host"),
    user: config.get("mysqlLocal.username"),
    password: config.get("mysqlLocal.password"),
    database: config.get("mysqlLocal.db"),
    waitForConnections: true,
    connectionLimit: 10,
    queueLimit: 0,
});

export const testDBConnection = async () => {
    try {
        const connection = await pool.getConnection();
        const [rows] = await connection.query('SELECT 1 + 1 AS solution');

        if ((rows as any)[0].solution === 2) {
            logger.info("MySQL DB connected successfully 🎉");
        }

        connection.release();
    } catch (error) {
        logger.error("Error testing DB connection:", error);
        throw error;
    }
};

testDBConnection();