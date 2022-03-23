import AppError from "../util/AppError.js";

const handleCastErrorDB = (err) => {
    const message = `Invalid ${err.path}: ${err.value}`;
    return new AppError(message, 400);
};

const handleDuplicateFieldDB = (err) => {
    const value = err.errmsg.match(/(["'])(\\?.)*?\1/);
    const message = `Duplicate field value:${value}. Please use another value `;
    return new AppError(message, 400);
};

const handleValidationErrorDB = (err) => {
    const error = Object.values(err.errors).map((el) => el.message);
    const message = `Invalid input Data. ${error.join('. ')}`;
    return new AppError(message, 400);
};

const handleJWTError = () => {
    new AppError('Invalid token, Please log in again', 401);
};

const handleJWTExpiredError = () => {
    new AppError('Your token has expired, Please log in again', 401);
};

const sendErrorInDevelopment = (err, req, res) => {
    return res.status(err.statusCode).json({
        status: err.status,
        error: err,
        message: err.message,
        stack: err.stack,
    });
}

const sendErrorInProduction = (err, req, res) => {
    if (err.isOperational) {
        return res.status(err.statusCode).json({
            status: err.status,
            message: err.message,
        });
    };

    //internal errors
    console.log(err);
    return res.status(500).json({
        status: 'error',
        message: `Looks like Something went wrong, Please don't hate us  `,
    });
};

const globalErrorHandler = (err, req, res, next) => {
    err.statusCode = err.statusCode || 500;
    err.status = err.status || 'error';

    if (process.env.NODE_ENV === 'development') {
        sendErrorInDevelopment(err, req, res);
    } else if (process.env.NODE_ENV === 'production') {
        let error = { ...err };
        error.message = err.message;
        if (error.name === 'CastError') error = handleCastErrorDB(error);
        if (error.code === 11000) error = handleDuplicateFieldDB(error);
        if (error.name === 'ValidationError')
            error = handleValidationErrorDB(error);
        if (error.name === 'JsonWebTokenError') error = handleJWTError();
        if (error.name === 'TokenExpiredError') error = handleJWTExpiredError();

        sendErrorInProduction(error, req, res);
    }
};

export default globalErrorHandler;