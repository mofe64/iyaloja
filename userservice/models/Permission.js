import mongoose from 'mongoose';

const PermissionSchema = new mongoose.Schema(
    {
        name: {
            type: String,
            required: [true, "Please provide the permission name"],
            unique:[true, "Permission already exists"]
        }
    }
)

const Permission  = mongoose.model("Permission", PermissionSchema);

export  default Permission