import React, { Component } from 'react'
import { observer } from 'mobx-react'
import Message from './Message'
import './index.css'

@observer
class Messages extends Component {
    componentDidUpdate(prevProps, prevState) {
        this.updateScroll()
    }
    render() {
        return (
            <div className="Messages">
                {
                    this.props.messages.map((obj) => (
                        <Message key={obj.id} date={new Date(obj.createTime * 1000)} message={obj.text} sender={obj.sender.username}/>
                    ))
                }
            </div>
        )
    }

    updateScroll() {
        if (this.props.messages.length > 0) {
            let messageEls = document.getElementsByClassName('Message')
            messageEls[messageEls.length - 1].scrollIntoView()
        }
    }
}

export default Messages