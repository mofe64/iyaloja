import mongoose from 'mongoose';
import bcrypt from 'bcryptjs';
import validator from 'validator';

const userSchema = new mongoose.Schema(
    {
        firstName: {
            type: String,
            required: [true, 'Please enter a firstname'],
        },
        lastName: {
            type: String,
            required: [true, 'Please enter a lastname'],
        },
        email: {
            type: String,
            required: [true, 'Please enter an email'],
            unique: true,
            lowercase: true,
            validate: [validator.isEmail, 'please provide a valid email'],
        },
        password: {
            type: String,
            required: [true, 'please provide a password'],
            minlength: 8,
        },
        passwordChangedAt: {
            type: Date,
        },
        passwordResetToken: String,
        passwordResetExpires: Date,
        active: {
            type: Boolean,
            default: true,
        },
        roles: {
            type: [
                {
                    type: mongoose.Schema.Types.ObjectId,
                    ref: "Role"
                }
            ],
            required: [true, "User must have role"],
        },
    }
);
userSchema.pre(/^find/, async function (next) {
    this.populate({
        path: 'roles',
        select: 'name permissions'
    })
    next()
})
userSchema.pre('save', async function (next) {
    if (!this.isModified('password')) return next();

    this.password = await bcrypt.hash(this.password, 12);
    next();
});

const User = mongoose.model('User', userSchema);

export default User;