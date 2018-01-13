import React, { Component } from 'react'
import { observer, inject } from 'mobx-react'

import UserForm from '../../components/UserForm'

import './index.css'

@inject('commonStore')
@observer
class Login extends Component {
    render() {
        return (
            <div className="Login convex">
                <UserForm />
            </div>
        )
    }
}

export default Login