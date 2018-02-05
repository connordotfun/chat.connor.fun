import React, { Component } from 'react'
import { observer, inject } from 'mobx-react'

import './index.css'

@inject('commonStore')
@observer
class VerifyLanding extends Component {
    componentWillMount() {
        this.props.history.replace('/')
        document.getElementsByTagName('body')[0].classList = []
    }
    
    render() {
        return (
            <div className="VerifyLanding convex">
                <h1>Welcome back, {this.props.commonStore.user.username}!</h1>
                <p>You need to verify your account before using <em>chat.connor.fun</em>.</p>
                <p>Check your email or click here to resend the verification message.</p>
            </div>
        )
    }
}

export default VerifyLanding