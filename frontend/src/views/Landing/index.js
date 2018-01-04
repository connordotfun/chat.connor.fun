import React, { Component } from 'react'
import { observer, inject } from 'mobx-react'

import UserForm from '../../components/UserForm'
import RegisterForm from '../../components/RegisterForm'

import './index.css'

@inject('commonStore')
@observer
class Landing extends Component {
    componentWillMount() {
        this.props.history.replace('/')
    }
    
    render() {
        return (
            <div className="Landing convex">
                <h1>Welcome!</h1>
                <p><em>chat.connor.fun</em> is a cool hangout place.</p>
                <UserForm />
                <RegisterForm />
            </div>
        )
    }
}

export default Landing