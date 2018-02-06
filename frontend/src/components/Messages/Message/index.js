import React, { Component } from 'react'
import './index.css'
import IntlRelativeFormat from 'intl-relativeformat'
// require('intl-relativeformat/dist/locale-data/en.js')

class Message extends Component {
    render() {
        let rf = new IntlRelativeFormat('en-US')
        return (
            <div className="Message">
                <div className="avatar" >
                    <img className="concave-small" src={"https://sigil.cupcake.io/" + this.props.sender + ".png?inverted=1"} alt={this.props.sender} />
                </div>
                <span className="handle">{this.props.sender}</span>
                <span className="timestamp">{rf.format(this.props.date)}</span>
                <span className="content">{this.props.message}</span>
            </div>
        )
    }
}

export default Message