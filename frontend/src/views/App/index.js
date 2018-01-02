import React, { Component } from 'react'
import { observer, inject } from 'mobx-react'
import { Switch, Route, Redirect } from 'react-router-dom'

import Chat from '../Chat'
import Landing from '../Landing'

@inject('commonStore')
@observer
class App extends Component {
    render() {
        return (       
            <Switch>
                <Route exact path="/" component={Landing} />
                <Route path="/at/:room" component={Chat} />
            </Switch>
        )
    }
}

export default App