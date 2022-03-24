import catchAsync from "../util/catchAsync.js";
import Role from "../models/Role.js";
import AppError from "../util/AppError.js";
import Permission from "../models/Permission.js";


/**
 export const createPermission = catchAsync(async (req, res, next) => {
    const {name} = req.body;
    const newPerm = await Permission.create({name})
    res.status(201).json({
        status: 'success',
        permission: newPerm
    })
})
 **/


export const addPermissionToRole = catchAsync(async (req, res, next) => {
    const {roleName, permission} = req.body;
    const roleObj = await Role.findOne({name: roleName})
    if (!roleObj) {
        return next(new AppError("No role with that name", 404))
    }
    const rolePerms = roleObj['permissions']
    let permExists = false
    rolePerms.forEach(rp => {
        if (rp.name === permission) {
            permExists = true
        }
    })
    if (permExists) {
        return next(new AppError("Role already has this permission", 400))
    }
    let perm = await Permission.findOne({name: permission});
    if (!perm) {
        perm = await Permission.create({
            name: permission
        })
    }
    roleObj['permissions'].push(perm['_id']);
    await roleObj.save();
    res.status(200).json({
        status: 'success'
    })
})

export const removePermissionFromRole = catchAsync(async (req, res, next) => {
    const {roleName, permission} = req.body;
    const role = await Role.findOne({name: roleName})
    if (!role) {
        return next(new AppError("No role with that name", 404))
    }
    let rolePerms = role['permissions']
    let permExists = false
    rolePerms.forEach(rp => {
        if (rp.name === permission) {
            permExists = true
        }
    })
    if (!permExists) {
        return next(new AppError("Role does not have this permission", 400))
    }
    rolePerms = rolePerms.filter(rp => rp.name !== permission)
    const update = await Role.updateOne(
        {_id: role['_id']},
        {
            permissions: rolePerms
        },
        {
            safe: true,
            upsert: true,
            multi: false
        })
    res.status(200).json({
        status: 'success',
        operationDetails: update
    })
})

export const getRoles = catchAsync(async (req, res, next) => {
    const roles = await Role.find();
    res.status(200).json({
        status: 'success',
        roles
    })
})

const roleExists = async (roleName = '') => {
    const role = await Role.findOne({name: roleName})
    return !!role;

}

export const getRole = catchAsync(async (req, res, next) => {
    const {name} = req.params
    const role = await Role.findOne({name})
    if (!role) {
        return next(new AppError("No role with that name", 404))
    }
    res.status(200).json({
        status: 'success',
        role
    })
})

const permissionCheck = async (perms = []) => {
    const foundPerms = []
    const foundPermNames = []
    let newPerms = []

    const cursor = Permission.find().sort({name: 1}).cursor()
    await cursor.eachAsync(async function (doc) {
        if (perms.includes(doc.name)) {
            foundPerms.push(doc._id)
            foundPermNames.push(doc.name)
        }
    })
    newPerms = perms.filter((perm) => !foundPermNames.includes(perm))
    return {foundPerms, newPerms}
}


export const preparePermissions = async (permissions =[])=> {
    /**
     Goes through permission array provided and
     return an obj containing two arrays
     one which holds the list of the permissions which already exist
     and another which holds a list of new permissions
     **/
    const permissionCheckObj = await permissionCheck(permissions)

    // Create Permission obj for the new permissions provided
    const newPermissionObjs = permissionCheckObj['newPerms'].map((perm) => {
        return {name: perm.toUpperCase()}
    })

    // save new permissions
    const savedPermissions = await Permission.insertMany(newPermissionObjs);

    // extract ids from the newly saved permissions
    let permissionIds = savedPermissions.map((permission) => permission._id.toString())

    // Add already existing permission ids to list of permissions so that we create the ref with our role
    permissionIds = [...permissionIds, ...permissionCheckObj['foundPerms']]
    return permissionIds;
}

export const createRole = catchAsync(async (req, res, next) => {
    // Extract role and permissions array from body
    let {role, permissions} = req.body;
    role = role.toUpperCase()

    // Check if the role provided already exists
    const roleAlreadyExisting = await roleExists(role)
    if (roleAlreadyExisting) {
        return next(new AppError("Role already exists", 400))
    }
    // Prepare ids for permission refs
    const permissionIds = await preparePermissions(permissions);

    const savedRole = await Role.create({
        name: role,
        permissions: permissionIds
    })

    res.status(201).json({
        status: 'success',
        role: savedRole
    })
})