import React, { Component } from 'react'
import { observer, inject } from 'mobx-react'
import RegisterForm from '../../components/RegisterForm'

import './index.css'
import logo from '../../logo.png'

@inject('commonStore')
@observer
class Landing extends Component {
    componentWillMount() {
        this.props.history.replace('/')
        document.getElementsByTagName('body')[0].classList.add('no-padding')
    }
    
    render() {
        return (
            <div className="Landing">
                <div className="branding">
                    <h1 className="brand-title">Chat</h1>
                    <h2 className="subtitle">a <img className="logo" src={logo} alt="connor.fun" /> product.</h2>
                </div>
                <div className="account-container">
                    <div className="register-box">
                        <RegisterForm />
                    </div>
                </div>
                {/* <Link to="/login"><button className="convex">Log In</button></Link>
                <Link to="/register"><button className="convex">Register</button></Link> */}
            </div>
        )
    }
}

export default Landing