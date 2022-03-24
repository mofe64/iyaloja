import express from 'express';
import {
    registerAdmin,
    addRoleToUser,
    removeRoleFromUser
} from "../controllers/AdminController.js";

const adminRouter = express()

adminRouter
    .route("")
    .post(registerAdmin)

adminRouter
    .route("/users/role")
    .post(addRoleToUser)
    .patch(removeRoleFromUser)


export default adminRouter;