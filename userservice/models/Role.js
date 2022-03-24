import mongoose from 'mongoose';


const RoleSchema = new mongoose.Schema(
    {
        name: {
            type: String,
            required: [true, "Role must have a name"],
            unique: [true, "Role already exists"]
        },
        permissions: {
            type: [
                {
                    type: mongoose.Schema.Types.ObjectId,
                    ref: "Permission"
                }
            ],
            required: [true, "Role must have permissions"],
        }
    }
)
RoleSchema.pre(/^find/, async function(next){
    this.populate({
        path: 'permissions',
        select: 'name _id'
    })
    next()
})

const Role = mongoose.model("Role", RoleSchema);
export default Role;