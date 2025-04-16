const passport = require("passport");
const GoogleStrategy = require("passport-google-oauth20").Strategy;
import { logger } from "../utils/logger";
const config = require("config");

try {
    passport.use(
        new GoogleStrategy(
            {
                clientID: config.get("googleOauth.clientId"),
                clientSecret: config.get("googleOauth.clientSecret"),
                callbackURL: `${config.get("services.mtApi.baseUrl")}${config.get("googleOauth.callbackURL")}`
            },
        (accessToken, refreshToken, profile, done) => {
            return done(null, accessToken, profile);
        }
    ));
} catch (error) {
    logger.info("Error Initiating passport", error);
}
