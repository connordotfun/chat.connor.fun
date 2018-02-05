import React, { Component } from 'react'
import { observer, inject } from 'mobx-react'

import './index.css'

@inject('authStore', 'commonStore')
@observer
class VerifyAccount extends Component {
    componentWillMount() {
        document.getElementsByTagName('body')[0].classList = []
        this.props.authStore.verifyAccount(this.props.match.params.code)
    }
    
    render() {
        return (
            <div className="VerifyAccount convex">
                <h1>Welcome back, {this.props.commonStore.user.username}!</h1>
            </div>
        )
    }
}

export default VerifyAccount