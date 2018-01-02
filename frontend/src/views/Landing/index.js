import React, { Component } from 'react'
import { observer, inject } from 'mobx-react'

import UserForm from '../../components/UserForm'
import './index.css'

@inject('commonStore')
@observer
class Landing extends Component {
    render() {
        this.props.history.replace('/')
        return (
            <div className="Landing convex">
                <h1>Welcome!</h1>
                <p><em>chat.connor.fun</em> is a cool hangout place.</p>
                <UserForm />
                <p>If we gotta token, it's {this.props.commonStore.token}</p>
            </div>
        )
    }
}

export default Landing