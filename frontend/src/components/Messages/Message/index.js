import React, { Component } from 'react'
import './index.css'

class Message extends Component {
    render() {
        return (
            <div className="Message">
                <div className="avatar" >
                    <img src={"https://sigil.cupcake.io/" + this.props.sender + ".png?inverted=1"} alt={this.props.sender} />
                </div>
                <span className="handle">{this.props.sender}</span>
                <span className="content">{this.props.message}</span>
            </div>
        )
    }
}

export default Message