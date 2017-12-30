import React, { Component } from 'react'
import Header from '../../components/Header'
import Messages from '../../components/Messages'
import Input from '../../components/Input'
import './index.css';

class Chat extends Component {
    render() {
        return (
            <div className="Chat convex">
                <Header room={this.props.match.params.room}/>
                <Messages messages={[1,2,3]} />
                <Input />
            </div>
        )
    }
}

export default Chat