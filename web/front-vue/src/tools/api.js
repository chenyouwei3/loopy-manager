import { axios as request } from './request'; 

/*----------------------------首页---------------------------------*/
//用户登录
export function login(parameter){
    return request({
        url: "/sign_in",
        method: 'post',
        data:parameter
    })
}

//用户注册
export function register(parameter){
    return request({
        url: "/sign_up",
        method: 'post',
        data:parameter
    })
}
/*----------------------------权限中心（user）---------------------------------*/
//删除用户
export function userRemove (parameter) {
    return request({
        url: "/user/delete",
        method: 'post',
        data: parameter
    })
}

//修改用户
export function userEdit(parameter){
    return request({
        url: "/user/update",
        method: 'post',
        data: parameter
    })
}

//查询用户列表
export function userList(parameter){
    const queryString = new URLSearchParams(parameter).toString(); 
    return request({
        url: `/user/getList?${queryString}`, 
        method: 'get', 
    });
}
/*----------------------------权限中心（role）---------------------------------*/
//新增角色
export function roleInsert(parameter){
    return request({
        url: "/role/insert",
        method: 'post',
        data: parameter
    })
}

//删除角色
export function roleRemove (parameter) {
    return request({
        url: "/role/delete",
        method: 'post',
        data: parameter
    })
}

//修改角色
export function roleEdit(parameter){
    return request({
        url: "/role/update",
        method: 'post',
        data: parameter
    })
}

//查询角色列表
export function roleList(parameter){
    const queryString = new URLSearchParams(parameter).toString(); 
    return request({
        url: `/role/getList?${queryString}`, 
        method: 'get', 
    });
}

//根据用户id查询对应的角色
export function userByRolesList(userId) {
    return request({
        url: `/user/getRolesByUserId?id=${userId}`,
        method: 'get'
    })
}
/*----------------------------日志中心---------------------------------*/
//查询操作日志列表
export function logByOperationList(parameter){
    const queryString = new URLSearchParams(parameter).toString();    
    return request({
        url: `/log/operation/getList?${queryString}`,
        method: 'get'
    })
}

//查询运行日志列表
export function logByRunningList(parameter){
    const queryString = new URLSearchParams(parameter).toString();    
    return request({
        url: `/log/running/getList?${queryString}`,
        method: 'get'
    })
}

//查询单个日志
export function logByRunningOne(path) {
    return request({
        url: `/log/running/getOne?path=${path}`,
        method: 'get'
    })
}