import React, { Component } from 'react'
import { observer, inject } from 'mobx-react'

import UserForm from '../../components/UserForm'

import './index.css'
import logo from '../../logo.png'

@inject('commonStore')
@observer
class Login extends Component {
    componentWillMount() {
        document.getElementsByTagName('body')[0].classList.add('no-padding')
    }

    render() {
        return (
            <div className="Login">
                <div className="branding">
                    <h1 className="brand-title">Chat</h1>
                    <h2 className="subtitle">a <img className="logo" src={logo} alt="connor.fun" /> product.</h2>
                </div>
                <div className="account-container">
                    <div className="login-box">
                        <UserForm />
                    </div>
                </div>
            </div>
        )
    }
}

export default Login