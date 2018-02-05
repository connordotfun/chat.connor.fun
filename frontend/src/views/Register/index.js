import React, { Component } from 'react'
import { observer, inject } from 'mobx-react'

import RegisterForm from '../../components/RegisterForm'

import './index.css'

@inject('commonStore')
@observer
class Register extends Component {
    render() {
        return (
            <div className="Register convex">
                <RegisterForm />
            </div>
        )
    }
}

export default Register