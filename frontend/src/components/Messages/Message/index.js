import React, { Component } from 'react';
import './index.css';

class Message extends Component {
    render() {
        return (
            <div className="Message">
                <div className="avatar" >
                    <img src="https://avatars2.githubusercontent.com/u/3019167?s=460&v=4" />
                </div>
                <span className="handle">{this.props.sender}</span>
                <span className="content">{this.props.message}</span>
            </div>
        )
    }
}

export default Message