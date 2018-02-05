import React, { Component } from 'react'
import { observer, inject } from 'mobx-react'
import { Link } from 'react-router-dom'

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
                <Link to="/login"><button className="convex">Log In</button></Link>
                <Link to="/register"><button className="convex">Register</button></Link>
            </div>
        )
    }
}

export default Landing