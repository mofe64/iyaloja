import express from 'express';
import {
    addPermissionToRole,
    createRole,
    getRole,
    getRoles,
    removePermissionFromRole
} from "../controllers/RoleController.js";

const roleRouter = express()


roleRouter.route("")
    .post(createRole)
    .get(getRoles)

roleRouter.route("/permissions")
    .post(addPermissionToRole)
    .patch(removePermissionFromRole)

roleRouter.route("/:name")
    .get(getRole)


export default roleRouter