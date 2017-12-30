import React, { Component } from 'react';
import './index.css';

class Header extends Component {
    render() {
        return (
            <header className="Header">
                <button className="back convex">&lt;</button>
                <span className="room">@{this.props.room}</span>
                <button className="profile convex"></button>
            </header>
        )
    }
}

export default Header