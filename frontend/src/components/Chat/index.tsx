import * as React from 'react';
import { observer } from 'mobx-react';
import { observable } from 'mobx';
import './index.css';

interface Props {
    send: (message: string) => void;
    messages: string[];
}

@observer
class Chat extends React.Component<Props> {
    @observable private _$message: string;
    public render() {
        return (
            <div className="chat">
                <div className="messages">
                    {this.props.messages.map((message: string) => (
                        <p>{message}</p>
                    ))}
                </div>
                <input 
                    type="text"
                    id="message"
                    onKeyPress={this._handleKeyPress.bind(this)}
                    onChange={event => {this._$message = event.currentTarget.value; }}
                />
            </div>
        );
    }

    private _handleKeyPress(e: React.KeyboardEvent<HTMLInputElement>): void {
        console.log(e.key);
        if (e.key === 'Enter') {
            console.log(this._$message);
            this.props.send(this._$message);
            this._$message = '';
            e.currentTarget.value = '';
        }
    }
}

export default Chat;