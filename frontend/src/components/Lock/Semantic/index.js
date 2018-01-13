/**
 * 
 * <Lock to={"normalUser"}>
 *  <Component />
 * </Lock>
 * 
 */

import { Component } from 'react'

class Lock extends Component {
    /* if provided props.user role is at least props.minRole, render children */
    validRoles = [
        'anonUser',
        'normalUser',
        'unverifiedUser',
        'admin'
    ]

    render() {
        const userRoles = this.props.user.roles.map(role => role.name)
        const minRole = this.props.minRole
        if (this.validRoles.indexOf(minRole) !== -1 && userRoles.indexOf(minRole) !== -1) {
            return this.props.children
        } else {
            return null
        }
    }
}

export default Lock