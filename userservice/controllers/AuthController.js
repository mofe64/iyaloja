import User from "../../userservice/models/User.js";
import jwt from "jsonwebtoken";
import catchAsync from "../../userservice/util/catchAsync.js";
import AppError from '../../userservice/util/AppError.js';
import bcrypt from 'bcryptjs';
import {promisify} from 'util';
import {createAndSendToken} from "../util/util.js";
import Role from "../models/Role.js";
import {preparePermissions} from "./RoleController.js";


export const register = catchAsync(async (req, res, next) => {
    const {firstName, lastName, email, password} = req.body;
    const newUser = {firstName, lastName, email, password};
    let baseUserRole = await Role.findOne({name: 'USER'})
    if (!baseUserRole) {
        const permissions = [
            'ORDER_READ',
            'ORDER_WRITE',
            'INVENTORY_READ',
            'INVENTORY_WRITE',
            'STATS_READ',
            'STATS_WRITE',
            'VENDOR_READ',
        ]
        const permIds = await preparePermissions(permissions)
        baseUserRole = await Role.create({
            name: 'ROLE_USER',
            permissions: permIds
        })
    }
    newUser['roles'] = [baseUserRole['_id']]
    const user = await User.create(newUser)
    createAndSendToken(user, 201, res);
})


export const login = catchAsync(async (req, res, next) => {
    const {email, password} = req.body;
    if (!email || !password) {
        return next(new AppError('Please provide email and password', 400));
    }
    const user = await User.findOne({email});
    if (!user) {
        return next(new AppError('No user found with that email', 401));
    }
    const verified = await bcrypt.compare(password, user.password);
    if (!verified) {
        return next(new AppError('Incorrect password given', 401));
    }
    createAndSendToken(user, 200, res);
});

export const authenticate = catchAsync(async (req, res, next) => {
    let token;
    if (req.headers.authorization && req.headers.authorization.startsWith("Bearer")) {
        token = req.headers.authorization.split(' ')[1];
    }
    if (!token) {
        return next(new AppError('You are not logged in, please log in to gain access', 401));
    }
    const decodedToken = await promisify(jwt.verify)(
        token,
        process.env.JWT_SECRET
    );
    const embeddedId = decodedToken.id;
    const currentUser = await User.findById(embeddedId);
    if (!currentUser) {
        return next(new AppError('User on this token no longer exists', 401));
    }
    req.user = currentUser;
    next();
});