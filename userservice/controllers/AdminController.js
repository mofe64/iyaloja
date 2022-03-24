import User from "../../userservice/models/User.js";
import Role from "../models/Role.js";
import catchAsync from "../../userservice/util/catchAsync.js";
import AppError from '../../userservice/util/AppError.js';
import {createAndSendToken} from "../util/util.js";
import {preparePermissions} from "./RoleController.js";



export const registerAdmin = catchAsync(async (req, res, next) => {
    const {firstName, lastName, email, password} = req.body;
    const newUser = {firstName, lastName, email, password};
    let adminRole = await Role.findOne({name: 'admin'});
    if (!adminRole) {
        const permissions = ["ADMIN_READ", "ADMIN_WRITE"]
        const permIds = await preparePermissions(permissions)
        adminRole = await Role.create({
            name: 'ROLE_ADMIN',
            permissions: permIds
        });
    }
    newUser['roles'] = [adminRole['_id']];
    const admin = await User.create(newUser);
    createAndSendToken(admin, 201, res)
});


export const addRoleToUser = catchAsync(async (req, res, next) => {
    const {roleName, userId} = req.body
    const role = await Role.findOne({name: roleName})
    if (!role) {
        return next(new AppError("No role found with that name", 404))
    }
    const user = await User.findById(userId);
    if (!user) {
        return next(new AppError("No user found with that id", 404))
    }
    user['roles'].push(role)
    await user.save()
    res.status(200).json({
        status: "success",
    })
})
export const  removeRoleFromUser = catchAsync(async (req,res,next)=> {
    const {roleName, userId} = req.body;
    const role = await Role.findOne({name: roleName})
    if (!role) {
        return next(new AppError("No role found with that name", 404))
    }
    const user = await User.findById(userId);
    if (!user) {
        return next(new AppError("No user found with that id", 404))
    }
    user['roles'] = user['roles'].filter(role => role.name !== roleName)
    await user.save()
    res.status(200).json({
        status: "success",
    })

})