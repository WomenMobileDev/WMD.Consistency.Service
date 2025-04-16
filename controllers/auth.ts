const passport = require("passport");
import { isUserAlreadyExist } from "../models/users";
import { addOrUpdateUser } from "../services/users";
import { userData } from "../types/global";
import { logger } from "../utils/logger";
const config = require("config");

export const googleAuthLogin = (req, res, next) => {
    const { redirectURL } = req.query;

    return passport.authenticate("google", {
        scope: ["profile", "email"],
        state: redirectURL,
    })(req, res, next);
};

export const googleAuthCallback = (req, res, next) => {
    const { authRedirectionUrl } = handleRedirectUrl(req);

    return passport.authenticate(
        "google",
        { 
            scope: ["profile", "email"],
            session: false,
        },
        async (error, accessToken, user) => {
            if (error) {
                logger.error(error);
                return res.boom.unauthorized("User cannot be authenticated");
            }

            const {id: userId, name, emails, provider} = user;

            console.log(`############################ name ############################ `);
            console.log(name);

            const newUser: userData = {
                id: userId,
                firstName: name.givenName,
                lastName: name.familyName,
                email: emails[0].value,
                logInProvider: provider
            };

            const userAlreadyExist = await isUserAlreadyExist(newUser.email);

            if (!userAlreadyExist) {
                await addOrUpdateUser(newUser);
            }

            return await handleGoogleLogin(req, res, user, authRedirectionUrl);
        }
    )(req, res, next);
};

function handleRedirectUrl(req) {
    const mtUiUrl = new URL(config.get("services.mtUi.baseUrl"));
    let authRedirectionUrl = mtUiUrl;
    return {
        authRedirectionUrl,
    };
}

function handleGoogleLogin(req, res, user, authRedirectionUrl) {
    try {
        return res.redirect(authRedirectionUrl);
    } catch (error) {
        logger.error("Unexpected error during Google login", error);
        return res.boom.unauthorized("User cannot be authenticated");
    }
}