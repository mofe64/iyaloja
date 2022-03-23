import express from 'express';
import dotenv from 'dotenv';
import bodyParser from 'body-parser';
import AppError from './util/AppError.js';
import globalErrorHandler from './controllers/ErrorController.js';
import authRouter from './routers/authRouter.js';

const app = express();
dotenv.config({path: './config.env'});


app.use(express.json({limit: '10kb'}));
app.use(bodyParser.urlencoded({limit: '10kb', extended: false}));

app.use('/api/v1/auth', authRouter);

app.all('*', (req, res, next) => {
    next(new AppError(`Can't find ${req.originalUrl} on this server`, 404));
});

app.use(globalErrorHandler);


export default app