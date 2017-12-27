import * as React from 'react';
import Buddy from './models/Buddy';
import BuddyList from './components/BuddyList';
import Chat from './components/Chat';
import { observer } from 'mobx-react';
import { observable, action } from 'mobx';
import './App.css';

@observer
class App extends React.Component {
  @observable public $messages: string[];

  @action
  public componentWillMount() {
    if (history.state) {
      this.$messages = history.state.messages;
    } else {
      this.$messages = [];
    }
  }

  @action
  public render() {
    return (
      <div className="App">
        <BuddyList users={this._getUsers()}/>
        <Chat send={this._pushState.bind(this)} messages={this.$messages}/>
      </div>
    );
  }

  private _getUsers(): Buddy[] {
    return([
      {
        name: 'Connor Hudson',
        avatar: '',
        link: ''
      },
      {
        name: 'Bob Bobberson',
        avatar: '',
        link: ''
      }
    ]);
  }

  @action
  private _pushState(s: string): void {
    const titles: string[] = [
      'chat.connor.fun',
      'I bet polygon.com doesn\'t have a chat site',
      'the yin to flatiron.live\'s yang',
      'SEIZE THE MEMES OF PRODUCTION'
    ];
    let currentState = history.state || {};
    if (currentState.messages) {
      currentState.messages.push(s);
    } else {
      currentState.messages = [s];
    }
    this.$messages.push(s);
    history.pushState(currentState, titles[Math.floor(Math.random() * titles.length)]);
  }
}

export default App;
