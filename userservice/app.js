import express from 'express';
import dotenv from 'dotenv';
import bodyParser from 'body-parser';
const app = express();
dotenv.config({ path: './config.env' });


app.use(express.json({ limit: '10kb' }));
app.use(bodyParser.urlencoded({ limit: '10kb', extended: false }));


export default  app