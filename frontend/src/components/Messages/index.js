import React, { Component } from 'react'
import { observer } from 'mobx-react'
import Message from './Message'
import './index.css'

@observer
class Messages extends Component {
    render() {
        return (
            <div className="Messages">
                {
                    this.props.messages.map((obj) => (
                        <Message key={obj.createTime} message={obj.text} sender={obj.sender.username}/>
                    ))
                }
            </div>
        )
    }
}

export default Messages