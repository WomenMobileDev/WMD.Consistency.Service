import { userData } from "../types/global";
import { createUser } from "../models/users";

export const addOrUpdateUser = async (
    user: userData
) => {
    await createUser(user);
};