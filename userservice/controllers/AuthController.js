import User from "../../userservice/models/User.js";
import jwt from "jsonwebtoken";
import catchAsync from "../../userservice/util/catchAsync.js";
import AppError from '../../userservice/util/AppError.js';
import bcrypt from 'bcryptjs';
import {promisify} from 'util';

const signToken = (id) => {
    return jwt.sign({id}, process.env.JWT_SECRET, {
        expiresIn: process.env.JWT_EXPIRES_IN,
    });
};

const createAndSendToken = (user, statusCode, res) => {
    const token = signToken(user._id);

    res.status(statusCode).json({
        status: 'success',
        token,
        data: {
            user,
        },
    });
};

export const register = catchAsync(async (req, res, next) => {
    const {email, password} = req.body;
    const newUser = await User.create({email, password});
    createAndSendToken(newUser, 201, res);
});

export const login = catchAsync(async (req, res, next) => {
    const {email, password} = req.body;
    if (!email || !password) {
        return next(new AppError('Please provide email and password', 400));
    }
    const user = await User.findOne({email});
    if (!user) {
        return next(new AppError('No user found with that email', 401));
    }
    ;
    const verified = await bcrypt.compare(password, user.password);
    if (verified) {
        return next(new AppError('Incorrect password given', 401));
    }
    ;
    createSendToken(user, 200, res);
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