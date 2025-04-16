const config = require("config");
const winston = require("winston");

const options = {
    file: {
        level: "info",
        filename: "logs/app.log",
        handleExceptions: true,
        json: true,
        maxsize: 5242880, // 5MB
        maxFiles: 5,
        colorize: false,
    },
    console: {
        level: "info",
        handleExceptions: true,
        json: false,
        colorize: true,
        silent: process.env.NODE_ENV === "test", // Disable logs in test env
    },
};

export const logger = new winston.createLogger({
    transports: [
        ...(config.get("enableFileLogs")? [new winston.transports.File(options.file)]: []),
        ...(config.get("enableConsoleLogs")? [new winston.transports.Console(options.file)]: []),
    ],

    exitOnError: false,
});

logger.stream = {
    write: function (message, encoding) {
        logger.info(message);
    },
};
