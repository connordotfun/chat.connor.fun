import React, { Component } from 'react'
import './index.css'

class Header extends Component {
    render() {
        return (
            <header className="Header">
                <button className="back convex" onClick={this.props.handleLeave}>&lt;</button>
                <span className="room">@{this.props.room}</span>
                <button className="logout convex" onClick={this.props.handleExit}>&times;</button>
            </header>
        )
    }
}

export default Header