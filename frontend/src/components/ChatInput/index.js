import React, { Component } from 'react'
import { observer, inject } from 'mobx-react'
import './index.css'

@inject('socketStore') @observer
class ChatInput extends Component {
    handleChange = e => this.props.socketStore.setMessage(e.target.value)
    handleKeyPress = e => {
        var keyCode = e.keyCode || e.which
        if ((keyCode === '13' || keyCode === 13) && this.props.socketStore.connected) {
            this.props.socketStore.sendMessage()
            this.props.socketStore.setMessage("")
            e.target.value = ""
        }
    }

    render() {
        return (
            <div className="ChatInput">
                <input disabled={!this.props.socketStore.connected} className="concave" type="text" onKeyPress={this.handleKeyPress} onChange={this.handleChange} />
            </div>
        )
    }
}

export default ChatInput