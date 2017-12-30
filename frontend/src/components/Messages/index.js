import React, { Component } from 'react'
import Message from './Message'
import './index.css'

class Messages extends Component {
    render() {
        return (
            <div className="Messages">
                {
                    this.props.messages.map((value) => (
                        <Message />
                    ))
                }
            </div>
        )
    }
}

export default Messages