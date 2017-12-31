import React, { Component } from 'react'
import Header from '../../components/Header'
import Messages from '../../components/Messages'
import Input from '../../components/Input'
import './index.css';

class Landing extends Component {
    render() {
        return (
            <div className="Landing convex">
                <h1>Welcome!</h1>
                <p><em>chat.connor.fun</em> is a cool hangout place.</p>
            </div>
        )
    }
}

export default Landing