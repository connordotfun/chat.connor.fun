import React, { Component } from 'react'
import { observer, inject } from 'mobx-react'
import { Link, withRouter } from 'react-router-dom'

import './index.css'

@inject('commonStore')
@withRouter
@observer
class Home extends Component {
    render() {
        this.props.history.replace('/')
        return (
            <div className="Home convex">
                <h1>Welcome back, {this.props.commonStore.username}!</h1>
                <p>Here's a list of channels you can join.</p>
                <ul>
                    <li><Link to='/at/farrand'>Farrand</Link></li>
                    <li><Link to='/at/secretlobby'>Engineering Center (South Entrance)</Link></li>
                    <li><Link to='/at/atlas'>Roser ATLAS</Link></li>
                </ul>
            </div>
        )
    }
}

export default Home