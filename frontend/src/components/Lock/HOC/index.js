function Lock(wrappedComponent, minRole, user) {
    /* if provided user role is at least minRole, render children */
    const validRoles = [
        'anonUser',
        'normalUser',
        'unverifiedUser',
        'admin'
    ]

    const userRoles = user.roles.map(role => role.name)

    if (validRoles.indexOf(minRole) !== -1 && userRoles.indexOf(minRole) !== -1) {
        return wrappedComponent
    } else {
        return null
    }

}

export default Lock