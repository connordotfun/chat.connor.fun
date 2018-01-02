import React from 'react'
import ReactDOM from 'react-dom'
import promiseFinally from 'promise.prototype.finally'
import { useStrict } from 'mobx'
import { Provider } from 'mobx-react'

import { BrowserRouter, Switch, Route, Redirect } from 'react-router-dom'
import './index.css'
import registerServiceWorker from './registerServiceWorker'

import App from './views/App'

import authStore from './stores/authStore'
import commonStore from './stores/commonStore'
import socketStore from './stores/socketStore'

promiseFinally.shim()
useStrict(true)

const stores = {
    authStore,
    commonStore,
    socketStore
}

ReactDOM.render(
    <Provider {...stores}>
        <BrowserRouter>
            <App />
        </BrowserRouter>
    </Provider>,
    document.getElementById('root'))
registerServiceWorker()